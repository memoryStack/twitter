package repositories

import (
	"context"

	"twitter/backend/models"
)

type UserRepository interface {
	GetByAuth0ID(ctx context.Context, auth0ID string) (*models.User, error)
	UpsertByAuth0ID(ctx context.Context, u models.User) (*models.User, error)
}
