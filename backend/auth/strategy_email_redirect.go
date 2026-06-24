package auth

import (
	"context"
	"fmt"

	"github.com/auth0/go-jwt-middleware/v2/validator"
)

// EmailRedirectName is the strategy identifier for Auth0 Universal Login with email redirect.
const EmailRedirectName = "email_redirect"

// EmailRedirect authenticates users via Auth0's hosted login page (authorization code + redirect).
// To add another mode later (e.g. OTP), implement Strategy in a new file and register it in Init.
type EmailRedirect struct {
	cfg       *Config
	validator *validator.Validator
}

func NewEmailRedirect(cfg *Config, v *validator.Validator) *EmailRedirect {
	return &EmailRedirect{cfg: cfg, validator: v}
}

func (s *EmailRedirect) Name() string { return EmailRedirectName }

func (s *EmailRedirect) Config() *Config { return s.cfg }

func (s *EmailRedirect) AuthorizeURL(state string) string {
	return AuthorizeURL(s.cfg, state)
}

func (s *EmailRedirect) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	return ExchangeAuthorizationCode(ctx, s.cfg, code)
}

func (s *EmailRedirect) RefreshTokens(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return RefreshTokens(ctx, s.cfg, refreshToken)
}

func (s *EmailRedirect) LogoutURL() string {
	return LogoutURL(s.cfg)
}

func (s *EmailRedirect) ValidateAccessToken(ctx context.Context, raw string) (*validator.ValidatedClaims, error) {
	claims, err := s.validator.ValidateToken(ctx, raw)
	if err != nil {
		return nil, err
	}
	validated, ok := claims.(*validator.ValidatedClaims)
	if !ok || validated == nil {
		return nil, fmt.Errorf("unexpected claims type")
	}
	return validated, nil
}
