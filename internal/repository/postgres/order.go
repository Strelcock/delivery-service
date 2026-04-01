package postgres

import (
	"delivery-service/internal/domain"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *domain.Order) error {
	query := `
		INSERT INTO orders (address, loc_lat, loc_lon, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRowx(query, order.Address, order.LocLat, order.LocLon, order.Status).
		Scan(&order.ID, &order.CreatedAt)
}

func (r *OrderRepo) GetByID(id int64) (*domain.Order, error) {
	var order domain.Order
	err := r.db.Get(&order, `SELECT * FROM orders WHERE id = $1`, id)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}
	return &order, nil
}

func (r *OrderRepo) List() ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Select(&orders, `SELECT * FROM orders ORDER BY created_at DESC`)
	return orders, err
}

func (r *OrderRepo) Update(id int64, input domain.UpdateOrderInput) (*domain.Order, error) {
	order, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Address != nil {
		order.Address = *input.Address
	}
	if input.LocLat != nil {
		order.LocLat = *input.LocLat
	}
	if input.LocLon != nil {
		order.LocLon = *input.LocLon
	}

	_, err = r.db.Exec(
		`UPDATE orders SET address=$1, loc_lat=$2, loc_lon=$3 WHERE id=$4`,
		order.Address, order.LocLat, order.LocLon, id,
	)
	return order, err
}

func (r *OrderRepo) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (r *OrderRepo) AssignCourier(orderID, courierID int64) error {
	res, err := r.db.Exec(
		`UPDATE orders SET status=$1, courier_id=$2 WHERE id=$3 AND status=$4`,
		domain.OrderStatusAssigned, courierID, orderID, domain.OrderStatusPending,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("order %d not found or already assigned", orderID)
	}
	return nil
}

func (r *OrderRepo) ListPending() ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Select(&orders, `SELECT * FROM orders WHERE status=$1`, domain.OrderStatusPending)
	if errors.Is(err, nil) && orders == nil {
		return []*domain.Order{}, nil
	}
	return orders, err
}
