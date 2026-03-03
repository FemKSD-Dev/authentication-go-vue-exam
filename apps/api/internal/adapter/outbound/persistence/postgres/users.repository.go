package postgres

import (
	"authentication-project-exam/internal/core/model"
	"authentication-project-exam/internal/core/port/outbound"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) outbound.UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Save(ctx context.Context, user *model.User) error {
	record := toRecord(user)
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *UsersRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var record UserRecord

	err := r.db.WithContext(ctx).Where("username = ?", username).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return toModel(&record), nil
}

func (r *UsersRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var record UserRecord

	err := r.db.WithContext(ctx).
		Where("id = ? AND delete_at IS NULL", id).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return toModel(&record), nil
}
