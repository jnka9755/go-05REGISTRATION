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
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Registration, error)
		Update(ctx context.Context, register *UpdateRegister) error
		Count(ctx context.Context, filters Filters) (int, error)
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

func (r *repository) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Registration, error) {

	var registrations []domain.Registration

	tx := r.db.WithContext(ctx).Model(&registrations)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&registrations)

	if result.Error != nil {
		r.log.Println("Error-Repository GetAllRegistration ->", result.Error)
		return nil, result.Error
	}

	return registrations, nil
}

func (r *repository) Update(ctx context.Context, register *UpdateRegister) error {

	values := make(map[string]interface{})

	if register.Status != nil {
		values["status"] = *register.Status
	}

	result := r.db.WithContext(ctx).Model(&domain.Registration{}).Where("id = ?", register.ID).Updates(values)

	if result.Error != nil {
		r.log.Println("Error-Repository UdateRegister ->", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.log.Printf("Register with ID -> '%s' doesn't exist", register.ID)
		return ErrNotFound{register.ID}
	}

	return nil
}

func (r *repository) Count(ctx context.Context, filters Filters) (int, error) {

	var count int64
	tx := r.db.WithContext(ctx).Model(domain.Registration{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		r.log.Println("Error-Repository CountRegistration ->", err)
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {

	if filters.UserID != "" {
		tx = tx.Where("user_id = ?", filters.UserID)
	}

	if filters.CourseID != "" {
		tx = tx.Where("course_id = ?", filters.CourseID)
	}

	return tx
}
