package usecase

import (
	"errors"
	"math"
	"testing"
	"time"

	"delivery-service/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─── Моки ────────────────────────────────────────────────────────────────────

type mockOrderRepo struct{ mock.Mock }

func (m *mockOrderRepo) Create(o *domain.Order) error { return m.Called(o).Error(0) }
func (m *mockOrderRepo) GetByID(id int64) (*domain.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *mockOrderRepo) List() ([]*domain.Order, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *mockOrderRepo) Update(id int64, inp domain.UpdateOrderInput) (*domain.Order, error) {
	args := m.Called(id, inp)
	return args.Get(0).(*domain.Order), args.Error(1)
}
func (m *mockOrderRepo) Delete(id int64) error              { return m.Called(id).Error(0) }
func (m *mockOrderRepo) AssignCourier(oID, cID int64) error { return m.Called(oID, cID).Error(0) }
func (m *mockOrderRepo) ListPending() ([]*domain.Order, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Order), args.Error(1)
}

type mockCourierRepo struct{ mock.Mock }

func (m *mockCourierRepo) Create(c *domain.Courier) error { return m.Called(c).Error(0) }
func (m *mockCourierRepo) GetByID(id int64) (*domain.Courier, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Courier), args.Error(1)
}

func (m *mockCourierRepo) List() ([]*domain.Courier, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Courier), args.Error(1)
}

func (m *mockCourierRepo) Update(id int64, inp domain.UpdateCourierInput) (*domain.Courier, error) {
	args := m.Called(id, inp)
	return args.Get(0).(*domain.Courier), args.Error(1)
}
func (m *mockCourierRepo) Delete(id int64) error { return m.Called(id).Error(0) }
func (m *mockCourierRepo) ListFree() ([]*domain.Courier, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Courier), args.Error(1)
}
func (m *mockCourierRepo) SetBusy(id int64) error { return m.Called(id).Error(0) }

// ─── Хелперы ─────────────────────────────────────────────────────────────────

func mkOrder(id int64, lat, lon float64) *domain.Order {
	return &domain.Order{
		ID: id, Address: "ул. Тестовая",
		LocLat: lat, LocLon: lon,
		Status: domain.OrderStatusPending, CreatedAt: time.Now(),
	}
}

func mkCourier(id int64, lat, lon float64) *domain.Courier {
	return &domain.Courier{
		ID: id, Name: "Курьер",
		LocLat: lat, LocLon: lon,
		Status: domain.CourierStatusFree,
	}
}

func newUC(or domain.OrderRepository, cr domain.CourierRepository) *OrderUseCase {
	return NewOrderUseCase(or, cr)
}

// ─── hungarian ───────────────────────────────────────────────────────────────

func TestHungarian_Nil(t *testing.T) {
	assert.Nil(t, hungarian(nil))
}

func TestHungarian_Empty(t *testing.T) {
	assert.Nil(t, hungarian([][]float64{}))
}

func TestHungarian_OneByOne(t *testing.T) {
	assert.Equal(t, []int{0}, hungarian([][]float64{{5}}))
}

func TestHungarian_TwoByTwo_Diagonal(t *testing.T) {
	cost := [][]float64{
		{1, 3},
		{2, 1},
	}
	res := hungarian(cost)
	assert.Equal(t, []int{0, 1}, res)
}

func TestHungarian_TwoByTwo_SwapNeeded(t *testing.T) {
	cost := [][]float64{
		{4, 1},
		{2, 3},
	}
	res := hungarian(cost)
	assert.Equal(t, []int{1, 0}, res)
}

func TestHungarian_ThreeByThree(t *testing.T) {
	cost := [][]float64{
		{9, 2, 7},
		{3, 6, 2},
		{7, 8, 3},
	}
	res := hungarian(cost)
	require.Len(t, res, 3)

	// уникальность
	seen := map[int]bool{}
	for _, j := range res {
		assert.False(t, seen[j], "дублирующееся назначение на заказ %d", j)
		seen[j] = true
	}

	total := cost[0][res[0]] + cost[1][res[1]] + cost[2][res[2]]
	assert.Equal(t, 8.0, total)
}

func TestHungarian_MoreRowsThanCols(t *testing.T) {
	// 3 курьера, 2 заказа — один курьер не назначается
	cost := [][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}
	res := hungarian(cost)
	require.Len(t, res, 3)

	assigned := 0
	seen := map[int]bool{}
	for _, j := range res {
		if j == -1 {
			continue
		}
		assigned++
		assert.False(t, seen[j], "заказ %d назначен дважды", j)
		seen[j] = true
	}
	assert.Equal(t, 2, assigned)
}

func TestHungarian_OneRowManyCols(t *testing.T) {
	// один курьер берёт минимальный заказ (индекс 1, стоимость 1)
	res := hungarian([][]float64{{5, 1, 3, 2}})
	require.Len(t, res, 1)
	assert.Equal(t, 1, res[0])
}

func TestHungarian_AllZeros(t *testing.T) {
	cost := [][]float64{{0, 0}, {0, 0}}
	res := hungarian(cost)
	require.Len(t, res, 2)
	assert.NotEqual(t, res[0], res[1])
}

// ─── AssignOptimal ────────────────────────────────────────────────────────────

func TestAssignOptimal_NoPendingOrders(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	or.On("ListPending").Return([]*domain.Order{}, nil)

	res, err := newUC(or, cr).AssignOptimal()

	require.NoError(t, err)
	assert.Nil(t, res)
	cr.AssertNotCalled(t, "ListFree")
}

func TestAssignOptimal_NoFreeCouriers(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	or.On("ListPending").Return([]*domain.Order{mkOrder(1, 55.75, 37.61)}, nil)
	cr.On("ListFree").Return([]*domain.Courier{}, nil)

	_, err := newUC(or, cr).AssignOptimal()
	require.ErrorIs(t, err, domain.ErrNoFreeCouriers)
}

func TestAssignOptimal_OneToOne(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	o := mkOrder(10, 55.75, 37.61)
	c := mkCourier(20, 55.76, 37.62)

	or.On("ListPending").Return([]*domain.Order{o}, nil)
	cr.On("ListFree").Return([]*domain.Courier{c}, nil)
	or.On("AssignCourier", o.ID, c.ID).Return(nil)
	cr.On("SetBusy", c.ID).Return(nil)

	res, err := newUC(or, cr).AssignOptimal()

	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, c.ID, res[0].CourierID)
	assert.Equal(t, o.ID, res[0].OrderID)
	assert.Greater(t, res[0].Distance, 0.0)
	or.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestAssignOptimal_OptimalAssignment(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}

	o1 := mkOrder(1, 0.0, 0.0)
	o2 := mkOrder(2, 10.0, 10.0)
	cA := mkCourier(100, 0.1, 0.1) // близко к o1
	cB := mkCourier(200, 9.9, 9.9) // близко к o2

	or.On("ListPending").Return([]*domain.Order{o1, o2}, nil)
	cr.On("ListFree").Return([]*domain.Courier{cA, cB}, nil)
	or.On("AssignCourier", mock.AnythingOfType("int64"), mock.AnythingOfType("int64")).Return(nil)
	cr.On("SetBusy", mock.AnythingOfType("int64")).Return(nil)

	res, err := newUC(or, cr).AssignOptimal()

	require.NoError(t, err)
	require.Len(t, res, 2)

	seen := map[int64]bool{}
	for _, r := range res {
		assert.False(t, seen[r.OrderID], "заказ %d назначен дважды", r.OrderID)
		seen[r.OrderID] = true
	}

	totalDist := res[0].Distance + res[1].Distance
	naiveDist := math.Sqrt(math.Pow(cA.LocLat-o2.LocLat, 2)+math.Pow(cA.LocLon-o2.LocLon, 2)) +
		math.Sqrt(math.Pow(cB.LocLat-o1.LocLat, 2)+math.Pow(cB.LocLon-o1.LocLon, 2))
	assert.Less(t, totalDist, naiveDist)
}

func TestAssignOptimal_MoreCouriersThanOrders(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	o1 := mkOrder(1, 0.0, 0.0)
	cA := mkCourier(100, 0.1, 0.1)
	cB := mkCourier(200, 5.0, 5.0)

	or.On("ListPending").Return([]*domain.Order{o1}, nil)
	cr.On("ListFree").Return([]*domain.Courier{cA, cB}, nil)
	or.On("AssignCourier", o1.ID, mock.AnythingOfType("int64")).Return(nil)
	cr.On("SetBusy", mock.AnythingOfType("int64")).Return(nil)

	res, err := newUC(or, cr).AssignOptimal()
	require.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestAssignOptimal_ListPendingError(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	dbErr := errors.New("db dead")
	or.On("ListPending").Return([]*domain.Order(nil), dbErr)

	_, err := newUC(or, cr).AssignOptimal()
	require.ErrorIs(t, err, dbErr)
}

func TestAssignOptimal_ListFreeError(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	dbErr := errors.New("db dead")
	or.On("ListPending").Return([]*domain.Order{mkOrder(1, 0, 0)}, nil)
	cr.On("ListFree").Return([]*domain.Courier(nil), dbErr)

	_, err := newUC(or, cr).AssignOptimal()
	require.ErrorIs(t, err, dbErr)
}

func TestAssignOptimal_AssignCourierError(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	dbErr := errors.New("assign failed")
	o := mkOrder(1, 0, 0)
	c := mkCourier(10, 0.1, 0.1)

	or.On("ListPending").Return([]*domain.Order{o}, nil)
	cr.On("ListFree").Return([]*domain.Courier{c}, nil)
	or.On("AssignCourier", o.ID, c.ID).Return(dbErr)

	_, err := newUC(or, cr).AssignOptimal()
	require.ErrorIs(t, err, dbErr)
}

func TestAssignOptimal_SetBusyError(t *testing.T) {
	or := &mockOrderRepo{}
	cr := &mockCourierRepo{}
	dbErr := errors.New("set busy failed")
	o := mkOrder(1, 0, 0)
	c := mkCourier(10, 0.1, 0.1)

	or.On("ListPending").Return([]*domain.Order{o}, nil)
	cr.On("ListFree").Return([]*domain.Courier{c}, nil)
	or.On("AssignCourier", o.ID, c.ID).Return(nil)
	cr.On("SetBusy", c.ID).Return(dbErr)

	_, err := newUC(or, cr).AssignOptimal()
	require.ErrorIs(t, err, dbErr)
}
