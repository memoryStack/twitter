package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"twitter/backend/auth"
)

// RequireAuth validates the access JWT from Bearer header or httpOnly cookie.
func RequireAuth(c *fiber.Ctx) error {
	token := auth.AccessTokenFromRequest(c)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authentication"})
	}

	_, strategy, err := auth.ValidateAccessTokenAny(c.UserContext(), token)
	if err != nil {
		rt := auth.RefreshTokenFromRequest(c, nil)
		if rt == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		var tr *auth.TokenResponse
		var refreshStrategy auth.Strategy
		for _, s := range auth.Registry.All() {
			cfg := s.Config()
			if strings.TrimSpace(c.Cookies(cfg.RefreshCookieName)) == "" {
				continue
			}
			tr, err = s.RefreshTokens(c.UserContext(), rt)
			if err == nil {
				refreshStrategy = s
				break
			}
		}
		if tr == nil || refreshStrategy == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		if _, err = refreshStrategy.ValidateAccessToken(c.UserContext(), tr.AccessToken); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}
		strategy = refreshStrategy
		auth.SetAuthCookies(c, strategy.Config(), tr)
		token = strings.TrimSpace(tr.AccessToken)
	}

	c.Locals(auth.AccessTokenCtxKey, token)
	c.Locals(auth.StrategyCtxKey, strategy)
	return c.Next()
}
