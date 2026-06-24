package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"twitter/backend/models"
	"twitter/backend/repositories"
)

type tweetRepository struct {
	db *gorm.DB
}

func NewTweetRepository(db *gorm.DB) repositories.TweetRepository {
	return &tweetRepository{db: db}
}

func (r *tweetRepository) Create(ctx context.Context, t *models.Tweet) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *tweetRepository) ListByUserID(ctx context.Context, userID uint) ([]models.Tweet, error) {
	var tweets []models.Tweet
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&tweets).Error
	return tweets, err
}

func (r *tweetRepository) GetByID(ctx context.Context, id uint) (*models.Tweet, error) {
	var tweet models.Tweet
	err := r.db.WithContext(ctx).Preload("User").First(&tweet, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	return &tweet, nil
}

func (r *tweetRepository) GetByIDAndUserID(ctx context.Context, id uint, userID uint) (*models.Tweet, error) {
	var tweet models.Tweet
	err := r.db.WithContext(ctx).Preload("User").Where("id = ? AND user_id = ?", id, userID).First(&tweet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	return &tweet, nil
}

func (r *tweetRepository) Update(ctx context.Context, t *models.Tweet) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *tweetRepository) DeleteByIDAndUserID(ctx context.Context, id uint, userID uint) (bool, error) {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Tweet{})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}
