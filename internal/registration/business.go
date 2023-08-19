package registration

import (
	"context"
	"log"

	"github.com/jnka9755/go-05DOMAIN/domain"
)

type (
	Business interface {
		Create(ctx context.Context, register *CreateReq) (*domain.Registration, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Registration, error)
		Update(ctx context.Context, register *UpdateReq) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	business struct {
		log        *log.Logger
		repository Repository
	}

	Filters struct {
		UserID   string
		CourseID string
	}

	UpdateRegister struct {
		ID     string
		Status *string
	}
)

func NewBusiness(log *log.Logger, repository Repository) Business {
	return &business{
		log:        log,
		repository: repository,
	}
}

func (b business) Create(ctx context.Context, request *CreateReq) (*domain.Registration, error) {

	register := domain.Registration{
		UserID:   request.UserID,
		CourseID: request.CourseID,
		Status:   "P",
	}

	if err := b.repository.Create(ctx, &register); err != nil {
		return nil, err
	}

	return &register, nil
}

func (b business) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Registration, error) {

	registers, err := b.repository.GetAll(ctx, filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return registers, nil
}

func (b business) Update(ctx context.Context, request *UpdateReq) error {

	registerUpdate := UpdateRegister{
		ID:     request.ID,
		Status: request.Status,
	}

	if err := b.repository.Update(ctx, &registerUpdate); err != nil {
		return err
	}

	return nil
}

func (b business) Count(ctx context.Context, filters Filters) (int, error) {
	return b.repository.Count(ctx, filters)
}
