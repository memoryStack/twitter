package initializers

import (
	"twitter/backend/repositories"
	"twitter/backend/repositories/postgres"
)

var (
	UserRepo  repositories.UserRepository
	TweetRepo repositories.TweetRepository
)

func InitRepositories() {
	UserRepo = postgres.NewUserRepository(DB)
	TweetRepo = postgres.NewTweetRepository(DB)
}
