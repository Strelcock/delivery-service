package postgres

import (
	"delivery-service/internal/domain"

	"github.com/jmoiron/sqlx"
)

type CourierRepo struct {
	db *sqlx.DB
}

func NewCourierRepo(db *sqlx.DB) *CourierRepo {
	return &CourierRepo{db: db}
}

func (r *CourierRepo) Create(courier *domain.Courier) error {
	return r.db.QueryRowx(
		`INSERT INTO couriers (name, loc_lat, loc_lon, status) VALUES ($1,$2,$3,$4) RETURNING id`,
		courier.Name, courier.LocLat, courier.LocLon, courier.Status,
	).Scan(&courier.ID)
}

func (r *CourierRepo) GetByID(id int64) (*domain.Courier, error) {
	var c domain.Courier
	if err := r.db.Get(&c, `SELECT * FROM couriers WHERE id=$1`, id); err != nil {
		return nil, domain.ErrCourierNotFound
	}
	return &c, nil
}

func (r *CourierRepo) List() ([]*domain.Courier, error) {
	var couriers []*domain.Courier
	err := r.db.Select(&couriers, `SELECT * FROM couriers ORDER BY id`)
	return couriers, err
}

func (r *CourierRepo) Update(id int64, input domain.UpdateCourierInput) (*domain.Courier, error) {
	courier, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if input.Name != nil {
		courier.Name = *input.Name
	}
	if input.LocLat != nil {
		courier.LocLat = *input.LocLat
	}
	if input.LocLon != nil {
		courier.LocLon = *input.LocLon
	}
	_, err = r.db.Exec(
		`UPDATE couriers SET name=$1, loc_lat=$2, loc_lon=$3 WHERE id=$4`,
		courier.Name, courier.LocLat, courier.LocLon, id,
	)
	return courier, err
}

func (r *CourierRepo) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM couriers WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrCourierNotFound
	}
	return nil
}

func (r *CourierRepo) ListFree() ([]*domain.Courier, error) {
	var couriers []*domain.Courier
	err := r.db.Select(&couriers, `SELECT * FROM couriers WHERE status=$1`, domain.CourierStatusFree)
	return couriers, err
}

func (r *CourierRepo) SetBusy(id int64) error {
	_, err := r.db.Exec(`UPDATE couriers SET status=$1 WHERE id=$2`, domain.CourierStatusBusy, id)
	return err
}
