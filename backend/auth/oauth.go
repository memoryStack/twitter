package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TokenResponse is the subset of Auth0 /oauth/token JSON we care about.
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	IDToken          string `json:"id_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
}

func tokenEndpoint(cfg *Config) string {
	return "https://" + cfg.Domain + "/oauth/token"
}

var httpOAuth = &http.Client{Timeout: 20 * time.Second}

func doTokenForm(ctx context.Context, cfg *Config, form url.Values) (*TokenResponse, error) {
	body := strings.NewReader(form.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint(cfg), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpOAuth.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("auth0 token endpoint %s: %s", resp.Status, strings.TrimSpace(string(b)))
	}

	var tr TokenResponse
	if err := json.Unmarshal(b, &tr); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if tr.AccessToken == "" {
		return nil, fmt.Errorf("auth0 token response missing access_token")
	}
	return &tr, nil
}

// ExchangeAuthorizationCode swaps the auth code for tokens (Authorization Code flow).
func ExchangeAuthorizationCode(ctx context.Context, cfg *Config, code string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("code", code)
	form.Set("redirect_uri", cfg.CallbackURL)
	return doTokenForm(ctx, cfg, form)
}

// RefreshTokens uses a refresh token to obtain a new access (and optionally refresh) token.
func RefreshTokens(ctx context.Context, cfg *Config, refreshToken string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("audience", cfg.Audience)
	return doTokenForm(ctx, cfg, form)
}

// AuthorizeURL builds the browser redirect URL to start login (Universal Login / email redirect).
func AuthorizeURL(cfg *Config, state string) string {
	q := url.Values{}
	q.Set("client_id", cfg.ClientID)
	q.Set("response_type", "code")
	q.Set("redirect_uri", cfg.CallbackURL)
	q.Set("scope", "openid profile email offline_access")
	q.Set("audience", cfg.Audience)
	q.Set("state", state)

	if conn := strings.TrimSpace(cfg.Connection); conn != "" {
		q.Set("connection", conn)
	}

	return "https://" + cfg.Domain + "/authorize?" + q.Encode()
}

// LogoutURL returns the Auth0 logout URL that clears the IdP session in the browser.
func LogoutURL(cfg *Config) string {
	q := url.Values{}
	q.Set("client_id", cfg.ClientID)
	q.Set("returnTo", cfg.LogoutReturnURL)
	return "https://" + cfg.Domain + "/v2/logout?" + q.Encode()
}
