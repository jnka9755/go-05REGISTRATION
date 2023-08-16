package registration

import (
	"context"
	"log"

	"github.com/jnka9755/go-05DOMAIN/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(ctx context.Context, register *domain.Registration) error
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(l *log.Logger, db *gorm.DB) Repository {

	return &repository{
		db:  db,
		log: l,
	}
}

func (r *repository) Create(ctx context.Context, register *domain.Registration) error {

	if err := r.db.WithContext(ctx).Create(register).Error; err != nil {
		r.log.Println("Repository ->", err)
		return err
	}

	r.log.Println("Repository -> Create register with id: ", register.ID)
	return nil
}
