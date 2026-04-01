package usecase

import "delivery-service/internal/domain"

type OrderUseCase struct {
	orderRepo   domain.OrderRepository
	courierRepo domain.CourierRepository
}

func NewOrderUseCase(o domain.OrderRepository, c domain.CourierRepository) *OrderUseCase {
	return &OrderUseCase{orderRepo: o, courierRepo: c}
}

func (uc *OrderUseCase) Create(input domain.CreateOrderInput) (*domain.Order, error) {
	order := &domain.Order{
		Address: input.Address,
		LocLat:  input.LocLat,
		LocLon:  input.LocLon,
		Status:  domain.OrderStatusPending,
	}
	if err := uc.orderRepo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (uc *OrderUseCase) GetByID(id int64) (*domain.Order, error) {
	return uc.orderRepo.GetByID(id)
}

func (uc *OrderUseCase) List() ([]*domain.Order, error) {
	return uc.orderRepo.List()
}

func (uc *OrderUseCase) Update(id int64, input domain.UpdateOrderInput) (*domain.Order, error) {
	return uc.orderRepo.Update(id, input)
}

func (uc *OrderUseCase) Delete(id int64) error {
	return uc.orderRepo.Delete(id)
}

// AssignOptimal решает задачу назначения венгерским алгоритмом:
// минимизирует суммарное евклидово расстояние курьер→заказ.
// Назначает курьеров на pending-заказы и сохраняет результат в БД.
func (uc *OrderUseCase) AssignOptimal() ([]AssignmentResult, error) {
	orders, err := uc.orderRepo.ListPending()
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, nil
	}

	couriers, err := uc.courierRepo.ListFree()
	if err != nil {
		return nil, err
	}
	if len(couriers) == 0 {
		return nil, domain.ErrNoFreeCouriers
	}

	// Строим матрицу стоимостей: строки — курьеры, столбцы — заказы
	cost := make([][]float64, len(couriers))
	for i, courier := range couriers {
		cost[i] = make([]float64, len(orders))
		for j, order := range orders {
			cost[i][j] = euclidean(courier.LocLat, courier.LocLon, order.LocLat, order.LocLon)
		}
	}

	assignment := hungarian(cost)

	var results []AssignmentResult
	for courierIdx, orderIdx := range assignment {
		if orderIdx == -1 {
			continue
		}
		c := couriers[courierIdx]
		o := orders[orderIdx]

		if err := uc.orderRepo.AssignCourier(o.ID, c.ID); err != nil {
			return nil, err
		}
		if err := uc.courierRepo.SetBusy(c.ID); err != nil {
			return nil, err
		}

		results = append(results, AssignmentResult{
			CourierID:   c.ID,
			CourierName: c.Name,
			OrderID:     o.ID,
			OrderAddr:   o.Address,
			Distance:    cost[courierIdx][orderIdx],
		})
	}

	return results, nil
}

type AssignmentResult struct {
	CourierID   int64   `json:"courier_id"`
	CourierName string  `json:"courier_name"`
	OrderID     int64   `json:"order_id"`
	OrderAddr   string  `json:"order_address"`
	Distance    float64 `json:"distance"`
}
