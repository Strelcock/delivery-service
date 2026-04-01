package domain

import (
	"errors"
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusAssigned  OrderStatus = "assigned"
	OrderStatusDelivered OrderStatus = "delivered"
)

var (
	ErrOrderNotFound   = errors.New("order not found")
	ErrCourierNotFound = errors.New("courier not found")
	ErrNoFreeCouriers  = errors.New("no free couriers available")
)

type Order struct {
	ID        int64       `db:"id"         json:"id"`
	Address   string      `db:"address"    json:"address"`
	LocLat    float64     `db:"loc_lat"    json:"loc_lat"`
	LocLon    float64     `db:"loc_lon"    json:"loc_lon"`
	Status    OrderStatus `db:"status"     json:"status"`
	CourierID *int64      `db:"courier_id" json:"courier_id,omitempty"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
}

type CreateOrderInput struct {
	Address string  `json:"address"`
	LocLat  float64 `json:"loc_lat"`
	LocLon  float64 `json:"loc_lon"`
}

type UpdateOrderInput struct {
	Address *string  `json:"address,omitempty"`
	LocLat  *float64 `json:"loc_lat,omitempty"`
	LocLon  *float64 `json:"loc_lon,omitempty"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int64) (*Order, error)
	List() ([]*Order, error)
	Update(id int64, input UpdateOrderInput) (*Order, error)
	Delete(id int64) error
	AssignCourier(orderID, courierID int64) error
	ListPending() ([]*Order, error)
}
