package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func cookieSameSiteFiber(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "strict":
		return fiber.CookieSameSiteStrictMode
	case "none":
		return fiber.CookieSameSiteNoneMode
	default:
		return fiber.CookieSameSiteLaxMode
	}
}

// SetAuthCookies sets httpOnly access and refresh cookies.
func SetAuthCookies(c *fiber.Ctx, cfg *Config, tr *TokenResponse) {
	maxAgeAccess := tr.ExpiresIn
	if maxAgeAccess <= 0 {
		maxAgeAccess = 3600
	}
	ss := cookieSameSiteFiber(cfg.CookieSameSite)
	c.Cookie(&fiber.Cookie{
		Name:     cfg.AccessCookieName,
		Value:    tr.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: ss,
		MaxAge:   maxAgeAccess,
	})
	if tr.RefreshToken != "" {
		refreshMax := tr.RefreshExpiresIn
		if refreshMax <= 0 {
			refreshMax = 60 * 60 * 24 * 30
		}
		c.Cookie(&fiber.Cookie{
			Name:     cfg.RefreshCookieName,
			Value:    tr.RefreshToken,
			Path:     "/",
			HTTPOnly: true,
			Secure:   cfg.CookieSecure,
			SameSite: ss,
			MaxAge:   refreshMax,
		})
	}
}
