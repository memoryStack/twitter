package auth

import (
	"context"

	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// Strategy defines one authentication mode (e.g. email redirect, OTP, social).
// Controllers depend on this interface — not on a specific IdP flow — so new modes
// are added by implementing Strategy and registering with Registry.
type Strategy interface {
	// Name is the stable identifier used in API routes and ?mode= query params.
	Name() string

	// Config returns OAuth/OIDC settings for this strategy.
	Config() *Config

	// AuthorizeURL returns the browser redirect URL that starts login.
	AuthorizeURL(state string) string

	// ExchangeCode completes the authorization-code redirect flow.
	ExchangeCode(ctx context.Context, code string) (*TokenResponse, error)

	// RefreshTokens rotates the access token using a refresh token.
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// LogoutURL returns the IdP logout URL for the browser.
	LogoutURL() string

	// ValidateAccessToken validates an access JWT issued for this strategy.
	ValidateAccessToken(ctx context.Context, raw string) (*validator.ValidatedClaims, error)
}
