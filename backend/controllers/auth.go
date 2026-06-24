package controllers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"twitter/backend/auth"
	"twitter/backend/initializers"
	"twitter/backend/models"
	"twitter/backend/repositories"
)

func cookieSameSite(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "strict":
		return fiber.CookieSameSiteStrictMode
	case "none":
		return fiber.CookieSameSiteNoneMode
	default:
		return fiber.CookieSameSiteLaxMode
	}
}

func authCookiePath() string {
	return "/api/auth"
}

func setFlowCookies(c *fiber.Ctx, cfg *auth.Config, mode, state string) {
	ss := cookieSameSite(cfg.CookieSameSite)
	path := authCookiePath()
	c.Cookie(&fiber.Cookie{
		Name:     cfg.StateCookieName,
		Value:    state,
		Path:     path,
		HTTPOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: ss,
		MaxAge:   600,
	})
	c.Cookie(&fiber.Cookie{
		Name:     cfg.ModeCookieName,
		Value:    mode,
		Path:     path,
		HTTPOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: ss,
		MaxAge:   600,
	})
}

func clearFlowCookies(c *fiber.Ctx, cfg *auth.Config) {
	ss := cookieSameSite(cfg.CookieSameSite)
	path := authCookiePath()
	clear := func(name string) {
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			HTTPOnly: true,
			Secure:   cfg.CookieSecure,
			SameSite: ss,
			MaxAge:   -1,
		})
	}
	clear(cfg.StateCookieName)
	clear(cfg.ModeCookieName)
}

func clearAuthCookies(c *fiber.Ctx, cfg *auth.Config) {
	ss := cookieSameSite(cfg.CookieSameSite)
	clear := func(name, path string) {
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			HTTPOnly: true,
			Secure:   cfg.CookieSecure,
			SameSite: ss,
			MaxAge:   -1,
		})
	}
	clear(cfg.AccessCookieName, "/")
	clear(cfg.RefreshCookieName, "/")
	clear(cfg.RefreshCookieName, cfg.RefreshCookiePath)
}

func strategyFromRequest(c *fiber.Ctx) (auth.Strategy, error) {
	mode := strings.TrimSpace(c.Query("mode"))
	if mode == "" {
		if defaultS, err := auth.Registry.Default(); err == nil {
			mode = strings.TrimSpace(c.Cookies(defaultS.Config().ModeCookieName))
		}
	}
	return auth.Registry.Get(mode)
}

// AuthLogin starts login for the requested strategy (?mode=email_redirect; defaults to the registered default).
func AuthLogin(c *fiber.Ctx) error {
	s, err := strategyFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cfg := s.Config()
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create login state"})
	}
	state := hex.EncodeToString(b)
	setFlowCookies(c, cfg, s.Name(), state)

	return c.Redirect(s.AuthorizeURL(state), fiber.StatusFound)
}

// AuthCallback handles the OAuth redirect, exchanges the code, and sets httpOnly cookies.
func AuthCallback(c *fiber.Ctx) error {
	if errMsg := c.Query("error"); errMsg != "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":             errMsg,
			"error_description": c.Query("error_description"),
		})
	}

	s, err := strategyFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	cfg := s.Config()

	code := c.Query("code")
	stateQ := c.Query("state")
	stateC := c.Cookies(cfg.StateCookieName)
	clearFlowCookies(c, cfg)

	if code == "" || stateQ == "" || stateC == "" ||
		subtle.ConstantTimeCompare([]byte(stateQ), []byte(stateC)) != 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid or missing OAuth state or code"})
	}

	tr, err := s.ExchangeCode(c.UserContext(), code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	auth.SetAuthCookies(c, cfg, tr)

	if tr.IDToken != "" {
		if _, err := saveUserFromIDToken(c, tr.IDToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.Redirect(cfg.PostLoginRedirect, fiber.StatusFound)
}

// AuthRefresh rotates the access token using the refresh token cookie.
func AuthRefresh(c *fiber.Ctx) error {
	s, err := strategyFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	cfg := s.Config()

	rt := strings.TrimSpace(c.Cookies(cfg.RefreshCookieName))
	if rt == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing refresh token"})
	}

	tr, err := s.RefreshTokens(c.UserContext(), rt)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	auth.SetAuthCookies(c, cfg, tr)
	return c.JSON(fiber.Map{
		"token_type":         tr.TokenType,
		"expires_in":         tr.ExpiresIn,
		"refresh_expires_in": tr.RefreshExpiresIn,
	})
}

// AuthLogout clears app cookies and returns the IdP logout URL for the browser.
func AuthLogout(c *fiber.Ctx) error {

	// return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 	"message": "Logged out successfully",
	// })

	s, err := strategyFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	cfg := s.Config()
	clearAuthCookies(c, cfg)
	return c.JSON(fiber.Map{
		"logout_url": s.LogoutURL(),
	})
}

// AuthMe returns the authenticated user from the database.
func AuthMe(c *fiber.Ctx) error {
	sub, err := subjectFromAccessToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	u, err := initializers.UserRepo.GetByAuth0ID(c.UserContext(), sub)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch user"})
	}
	return c.JSON(fiber.Map{"user": fiber.Map{
		"id":             u.ID,
		"first_name":     u.FirstName,
		"last_name":      u.LastName,
		"email":          u.Email,
		"phone_number":   u.PhoneNumber,
		"email_verified": u.EmailVerified,
		"image_url":      u.Image,
	}})
}

func subjectFromAccessToken(c *fiber.Ctx) (string, error) {
	token := auth.AccessTokenFromCtx(c)
	if token == "" {
		return "", fmt.Errorf("missing access token")
	}
	validated, _, err := auth.ValidateAccessTokenAny(c.UserContext(), token)
	if err != nil {
		return "", fmt.Errorf("invalid access token")
	}
	sub := strings.TrimSpace(validated.RegisteredClaims.Subject)
	if sub == "" {
		return "", fmt.Errorf("missing subject")
	}
	return sub, nil
}

func saveUserFromIDToken(c *fiber.Ctx, idToken string) (*models.User, error) {
	claims, err := auth.IDTokenClaims(idToken)
	if err != nil {
		return nil, err
	}
	u, err := userFromIDTokenClaims(claims)
	if err != nil {
		return nil, err
	}
	return initializers.UserRepo.UpsertByAuth0ID(c.UserContext(), u)
}

func userFromIDTokenClaims(claims map[string]interface{}) (models.User, error) {
	sub := claimString(claims, "sub")
	if sub == "" {
		return models.User{}, fmt.Errorf("id_token missing sub")
	}

	first := claimString(claims, "given_name")
	last := claimString(claims, "family_name")
	if first == "" {
		first, last = splitName(claimString(claims, "name"))
	}

	email := claimString(claims, "email")
	if email == "" {
		return models.User{}, fmt.Errorf("id_token missing email")
	}

	return models.User{
		FirstName:     first,
		LastName:      last,
		Email:         email,
		PhoneNumber:   claimString(claims, "phone_number"),
		EmailVerified: claimBool(claims, "email_verified"),
		Image:         claimString(claims, "picture"),
		Auth0ID:       sub,
	}, nil
}

func claimString(claims map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		v, ok := claims[key]
		if !ok || v == nil {
			continue
		}
		if s, ok := v.(string); ok {
			if s = strings.TrimSpace(s); s != "" {
				return s
			}
		}
	}
	return ""
}

func claimBool(claims map[string]interface{}, key string) bool {
	v, ok := claims[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}

func splitName(full string) (first, last string) {
	full = strings.TrimSpace(full)
	if full == "" {
		return "", ""
	}
	parts := strings.SplitN(full, " ", 2)
	first = parts[0]
	if len(parts) > 1 {
		last = strings.TrimSpace(parts[1])
	}
	return first, last
}
