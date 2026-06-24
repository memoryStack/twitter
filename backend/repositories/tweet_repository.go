package repositories

import (
	"context"

	"twitter/backend/models"
)

type TweetRepository interface {
	Create(ctx context.Context, t *models.Tweet) error
	ListByUserID(ctx context.Context, userID uint) ([]models.Tweet, error)
	GetByID(ctx context.Context, id uint) (*models.Tweet, error)
	GetByIDAndUserID(ctx context.Context, id uint, userID uint) (*models.Tweet, error)
	Update(ctx context.Context, t *models.Tweet) error
	DeleteByIDAndUserID(ctx context.Context, id uint, userID uint) (bool, error)
}
