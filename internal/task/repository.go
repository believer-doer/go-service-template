package task

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("task not found")

// Repository separates business behavior from the storage implementation.
// In a test, provide an in-memory fake; in production, use PostgresRepository.
type Repository interface {
	Create(context.Context, *Task) error
	List(context.Context) ([]Task, error)
	Get(context.Context, uint) (Task, error)
	Update(context.Context, *Task) error
	Delete(context.Context, uint) error
}

type PostgresRepository struct{ db *gorm.DB }

func NewPostgresRepository(db *gorm.DB) *PostgresRepository { return &PostgresRepository{db: db} }

func (r *PostgresRepository) Create(ctx context.Context, item *Task) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PostgresRepository) List(ctx context.Context) ([]Task, error) {
	var items []Task
	err := r.db.WithContext(ctx).Order("id ASC").Find(&items).Error
	return items, err
}

func (r *PostgresRepository) Get(ctx context.Context, id uint) (Task, error) {
	var item Task
	err := r.db.WithContext(ctx).First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Task{}, ErrNotFound
	}
	return item, err
}

func (r *PostgresRepository) Update(ctx context.Context, item *Task) error {
	result := r.db.WithContext(ctx).Save(item)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&Task{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete task: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
