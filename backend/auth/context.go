package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AccessTokenCtxKey is the fiber Locals key for the validated access JWT string.
const AccessTokenCtxKey = "auth_access_token"

// StrategyCtxKey is the fiber Locals key for the Strategy that validated the request.
const StrategyCtxKey = "auth_strategy"

// AccessTokenFromCtx returns the access JWT placed by RequireAuth, or "" if missing.
func AccessTokenFromCtx(c *fiber.Ctx) string {
	v := c.Locals(AccessTokenCtxKey)
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

// StrategyFromCtx returns the Strategy placed by RequireAuth, or nil.
func StrategyFromCtx(c *fiber.Ctx) Strategy {
	v := c.Locals(StrategyCtxKey)
	if v == nil {
		return nil
	}
	s, ok := v.(Strategy)
	if !ok {
		return nil
	}
	return s
}
