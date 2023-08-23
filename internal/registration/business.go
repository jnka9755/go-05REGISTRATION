package registration

import (
	"context"
	"log"

	"github.com/jnka9755/go-05DOMAIN/domain"

	sdkCourse "github.com/jnka9755/go-05SDKCOURSE/course"
	sdkUser "github.com/jnka9755/go-05SDKUSER/user"
)

type (
	Business interface {
		Create(ctx context.Context, register *CreateReq) (*domain.Registration, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Registration, error)
		Update(ctx context.Context, register *UpdateReq) error
		Count(ctx context.Context, filters Filters) (int, error)
	}

	business struct {
		log         *log.Logger
		repository  Repository
		transUser   sdkUser.Transport
		transCourse sdkCourse.Transport
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

func NewBusiness(log *log.Logger, repository Repository, transUser sdkUser.Transport, transCourse sdkCourse.Transport) Business {
	return &business{
		log:         log,
		repository:  repository,
		transUser:   transUser,
		transCourse: transCourse,
	}
}

func (b business) Create(ctx context.Context, request *CreateReq) (*domain.Registration, error) {

	register := domain.Registration{
		UserID:   request.UserID,
		CourseID: request.CourseID,
		Status:   domain.Pending,
	}

	if _, err := b.transUser.Get(request.UserID); err != nil {
		return nil, err
	}

	if _, err := b.transCourse.Get(request.CourseID); err != nil {
		return nil, err
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

	if request.Status != nil {
		switch domain.EnrollStatus(*request.Status) {
		case domain.Pending, domain.Active, domain.Studying, domain.Inactive:
		default:
			return ErrInvalidStatus{*request.Status}
		}
	}

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
