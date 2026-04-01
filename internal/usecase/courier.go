package usecase

import "delivery-service/internal/domain"

type CourierUseCase struct {
	repo domain.CourierRepository
}

func NewCourierUseCase(r domain.CourierRepository) *CourierUseCase {
	return &CourierUseCase{repo: r}
}

func (uc *CourierUseCase) Create(input domain.CreateCourierInput) (*domain.Courier, error) {
	courier := &domain.Courier{
		Name:   input.Name,
		LocLat: input.LocLat,
		LocLon: input.LocLon,
		Status: domain.CourierStatusFree,
	}
	if err := uc.repo.Create(courier); err != nil {
		return nil, err
	}
	return courier, nil
}

func (uc *CourierUseCase) GetByID(id int64) (*domain.Courier, error) {
	return uc.repo.GetByID(id)
}

func (uc *CourierUseCase) List() ([]*domain.Courier, error) {
	return uc.repo.List()
}

func (uc *CourierUseCase) Update(id int64, input domain.UpdateCourierInput) (*domain.Courier, error) {
	return uc.repo.Update(id, input)
}

func (uc *CourierUseCase) Delete(id int64) error {
	return uc.repo.Delete(id)
}
