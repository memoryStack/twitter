package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"twitter/backend/auth"
	"twitter/backend/initializers"
	"twitter/backend/models"
	"twitter/backend/repositories"
)

type createTweetRequest struct {
	Text     string `json:"text"`
	MediaURL string `json:"mediaUrl"`
}

type updateTweetRequest struct {
	Text     *string `json:"text"`
	MediaURL *string `json:"mediaUrl"`
	Likes    *int64  `json:"likes"`
}

func CreateTweet(c *fiber.Ctx) error {
	user, err := currentUserFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	var req createTweetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req.Text = strings.TrimSpace(req.Text)
	req.MediaURL = strings.TrimSpace(req.MediaURL)
	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "text is required"})
	}

	tweet := models.Tweet{
		Text:     req.Text,
		MediaURL: req.MediaURL,
		Likes:    0,
		UserID:   user.ID,
	}
	if err := initializers.TweetRepo.Create(c.UserContext(), &tweet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create tweet"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"tweet": tweet})
}

func GetMyTweets(c *fiber.Ctx) error {
	user, err := currentUserFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	tweets, err := initializers.TweetRepo.ListByUserID(c.UserContext(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch tweets"})
	}

	return c.JSON(fiber.Map{"tweets": tweets})
}

func GetTweetByID(c *fiber.Ctx) error {
	id, err := parseTweetID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tweet id"})
	}

	tweet, err := initializers.TweetRepo.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tweet not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch tweet"})
	}

	return c.JSON(fiber.Map{"tweet": tweet})
}

func UpdateMyTweet(c *fiber.Ctx) error {
	user, err := currentUserFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := parseTweetID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tweet id"})
	}

	tweet, err := initializers.TweetRepo.GetByIDAndUserID(c.UserContext(), id, user.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tweet not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch tweet"})
	}

	var req updateTweetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Text != nil {
		trimmed := strings.TrimSpace(*req.Text)
		if trimmed == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "text cannot be empty"})
		}
		tweet.Text = trimmed
	}
	if req.MediaURL != nil {
		tweet.MediaURL = strings.TrimSpace(*req.MediaURL)
	}
	if req.Likes != nil {
		if *req.Likes < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "likes cannot be negative"})
		}
		tweet.Likes = *req.Likes
	}

	if err := initializers.TweetRepo.Update(c.UserContext(), tweet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update tweet"})
	}

	return c.JSON(fiber.Map{"tweet": tweet})
}

func DeleteMyTweet(c *fiber.Ctx) error {
	user, err := currentUserFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	id, err := parseTweetID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tweet id"})
	}

	deleted, err := initializers.TweetRepo.DeleteByIDAndUserID(c.UserContext(), id, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete tweet"})
	}
	if !deleted {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tweet not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func LikeTweet(c *fiber.Ctx) error {
	id, err := parseTweetID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tweet id"})
	}

	tweet, err := initializers.TweetRepo.GetByID(c.UserContext(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tweet not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch tweet"})
	}

	tweet.Likes++
	if err := initializers.TweetRepo.Update(c.UserContext(), tweet); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to like tweet"})
	}

	return c.JSON(fiber.Map{"tweet": tweet})
}

func parseTweetID(raw string) (uint, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func currentUserFromToken(c *fiber.Ctx) (*models.User, error) {
	token := auth.AccessTokenFromCtx(c)
	if token == "" {
		return nil, errors.New("missing access token")
	}

	validated, _, err := auth.ValidateAccessTokenAny(c.UserContext(), token)
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	sub := strings.TrimSpace(validated.RegisteredClaims.Subject)
	if sub == "" {
		return nil, errors.New("missing subject in access token")
	}

	user, err := initializers.UserRepo.GetByAuth0ID(c.UserContext(), sub)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
