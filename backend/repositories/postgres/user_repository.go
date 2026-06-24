package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"twitter/backend/models"
	"twitter/backend/repositories"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByAuth0ID(ctx context.Context, auth0ID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("auth0_id = ?", auth0ID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpsertByAuth0ID(ctx context.Context, u models.User) (*models.User, error) {
	var existing models.User
	err := r.db.WithContext(ctx).Where("auth0_id = ?", u.Auth0ID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := r.db.WithContext(ctx).Create(&u).Error; err != nil {
				return nil, err
			}
			return &u, nil
		}
		return nil, err
	}

	existing.FirstName = u.FirstName
	existing.LastName = u.LastName
	existing.Email = u.Email
	existing.PhoneNumber = u.PhoneNumber
	existing.EmailVerified = u.EmailVerified
	if u.Image != "" {
		existing.Image = u.Image
	}

	if err := r.db.WithContext(ctx).Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}
