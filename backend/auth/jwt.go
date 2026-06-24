package auth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	gojwt "gopkg.in/go-jose/go-jose.v2/jwt"
)

func newJWTValidator(cfg *Config) (*validator.Validator, error) {
	issuer := "https://" + cfg.Domain + "/"
	issuerURL, err := url.Parse(issuer)
	if err != nil {
		return nil, fmt.Errorf("issuer url: %w", err)
	}
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	return validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuer,
		[]string{cfg.Audience},
		validator.WithAllowedClockSkew(30*time.Second),
	)
}

func validateWith(v *validator.Validator, raw string, ctx context.Context) (*validator.ValidatedClaims, error) {
	if v == nil {
		return nil, fmt.Errorf("jwt validator not initialized")
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("empty token")
	}
	claims, err := v.ValidateToken(ctx, raw)
	if err != nil {
		return nil, err
	}
	validated, ok := claims.(*validator.ValidatedClaims)
	if !ok || validated == nil {
		return nil, fmt.Errorf("unexpected claims type")
	}
	return validated, nil
}

// ValidateAccessTokenAny tries each registered strategy until one validates the token.
func ValidateAccessTokenAny(ctx context.Context, raw string) (*validator.ValidatedClaims, Strategy, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil, fmt.Errorf("empty token")
	}

	var lastErr error
	for _, s := range Registry.All() {
		claims, err := s.ValidateAccessToken(ctx, raw)
		if err == nil {
			return claims, s, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return nil, nil, lastErr
	}
	return nil, nil, fmt.Errorf("jwt validator not initialized")
}

// IDTokenClaims returns all id_token claims as map[string]interface{}.
// This only decodes claims; add signature/iss/aud validation if used for auth decisions.
func IDTokenClaims(raw string) (map[string]interface{}, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]interface{}{}, fmt.Errorf("empty id_token")
	}
	tok, err := gojwt.ParseSigned(raw)
	if err != nil {
		return nil, fmt.Errorf("parse id_token: %w", err)
	}
	claims := map[string]interface{}{}
	if err := tok.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return nil, fmt.Errorf("decode id_token claims: %w", err)
	}
	return claims, nil
}
