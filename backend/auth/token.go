package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AccessTokenFromRequest reads Bearer token or access-token cookies from any registered strategy.
func AccessTokenFromRequest(c *fiber.Ctx) string {
	if h := strings.TrimSpace(c.Get("Authorization")); len(h) > 7 && strings.EqualFold(h[:7], "bearer ") {
		if t := strings.TrimSpace(h[7:]); t != "" {
			return t
		}
	}
	for _, s := range Registry.All() {
		cfg := s.Config()
		if t := strings.TrimSpace(c.Cookies(cfg.AccessCookieName)); t != "" {
			return t
		}
	}
	return ""
}

// RefreshTokenFromRequest reads the refresh cookie for a strategy, or scans all strategies when s is nil.
func RefreshTokenFromRequest(c *fiber.Ctx, s Strategy) string {
	if s != nil {
		return strings.TrimSpace(c.Cookies(s.Config().RefreshCookieName))
	}
	for _, strategy := range Registry.All() {
		if t := strings.TrimSpace(c.Cookies(strategy.Config().RefreshCookieName)); t != "" {
			return t
		}
	}
	return ""
}
