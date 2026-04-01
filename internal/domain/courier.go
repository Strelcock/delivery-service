package domain

type CourierStatus string

const (
	CourierStatusFree CourierStatus = "free"
	CourierStatusBusy CourierStatus = "busy"
)

type Courier struct {
	ID      int64         `db:"id"      json:"id"`
	Name    string        `db:"name"    json:"name"`
	LocLat  float64       `db:"loc_lat" json:"loc_lat"`
	LocLon  float64       `db:"loc_lon" json:"loc_lon"`
	Status  CourierStatus `db:"status"  json:"status"`
}

type CreateCourierInput struct {
	Name   string  `json:"name"`
	LocLat float64 `json:"loc_lat"`
	LocLon float64 `json:"loc_lon"`
}

type UpdateCourierInput struct {
	Name   *string  `json:"name,omitempty"`
	LocLat *float64 `json:"loc_lat,omitempty"`
	LocLon *float64 `json:"loc_lon,omitempty"`
}

type CourierRepository interface {
	Create(courier *Courier) error
	GetByID(id int64) (*Courier, error)
	List() ([]*Courier, error)
	Update(id int64, input UpdateCourierInput) (*Courier, error)
	Delete(id int64) error
	ListFree() ([]*Courier, error)
	SetBusy(id int64) error
}
