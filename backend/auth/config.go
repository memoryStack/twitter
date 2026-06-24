package auth

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds Auth0 OAuth/OIDC settings loaded from the environment.
type Config struct {
	Domain            string
	ClientID          string
	ClientSecret      string
	CallbackURL       string
	Audience          string
	PostLoginRedirect string
	LogoutReturnURL   string
	Connection        string

	CookieSecure      bool
	CookieSameSite    string
	AccessCookieName  string
	RefreshCookieName string
	StateCookieName   string
	ModeCookieName    string
	RefreshCookiePath string
}

func loadAuth0Config(prefix string) (*Config, error) {
	key := func(suffix string) string {
		return prefix + suffix
	}
	get := func(suffix string) (string, error) {
		k := key(suffix)
		v := strings.TrimSpace(os.Getenv(k))
		if v == "" {
			return "", fmt.Errorf("missing required env %s", k)
		}
		return v, nil
	}

	domain, err := get("DOMAIN")
	if err != nil {
		return nil, err
	}
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")

	clientID, err := get("CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := get("CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	callbackURL, err := get("CALLBACK_URL")
	if err != nil {
		return nil, err
	}
	audience, err := get("AUDIENCE")
	if err != nil {
		return nil, err
	}

	postLogin := strings.TrimSpace(os.Getenv(key("POST_LOGIN_REDIRECT")))
	if postLogin == "" {
		postLogin = "http://localhost:5173/"
	}
	logoutReturn := strings.TrimSpace(os.Getenv(key("LOGOUT_RETURN_URL")))
	if logoutReturn == "" {
		logoutReturn = postLogin
	}

	connection := strings.TrimSpace(os.Getenv(key("CONNECTION")))

	cookieSecure := false
	if v := strings.TrimSpace(os.Getenv(key("COOKIE_SECURE"))); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key("COOKIE_SECURE"), err)
		}
		cookieSecure = b
	}

	sameSite := strings.ToLower(strings.TrimSpace(os.Getenv(key("COOKIE_SAMESITE"))))
	if sameSite == "" {
		sameSite = "lax"
	}
	if sameSite == "none" && !cookieSecure {
		return nil, fmt.Errorf("%s=none requires %s=true", key("COOKIE_SAMESITE"), key("COOKIE_SECURE"))
	}

	accessName := strings.TrimSpace(os.Getenv(key("ACCESS_COOKIE_NAME")))
	if accessName == "" {
		accessName = "access_token"
	}
	refreshName := strings.TrimSpace(os.Getenv(key("REFRESH_COOKIE_NAME")))
	if refreshName == "" {
		refreshName = "refresh_token"
	}
	stateName := strings.TrimSpace(os.Getenv(key("STATE_COOKIE_NAME")))
	if stateName == "" {
		stateName = "oauth_state"
	}
	modeName := strings.TrimSpace(os.Getenv(key("MODE_COOKIE_NAME")))
	if modeName == "" {
		modeName = "auth_mode"
	}

	refreshPath := strings.TrimSpace(os.Getenv(key("REFRESH_COOKIE_PATH")))
	if refreshPath == "" {
		refreshPath = "/api/auth"
	}

	return &Config{
		Domain:            domain,
		ClientID:          clientID,
		ClientSecret:      clientSecret,
		CallbackURL:       callbackURL,
		Audience:          audience,
		PostLoginRedirect: postLogin,
		LogoutReturnURL:   logoutReturn,
		Connection:        connection,
		CookieSecure:      cookieSecure,
		CookieSameSite:    sameSite,
		AccessCookieName:  accessName,
		RefreshCookieName: refreshName,
		StateCookieName:   stateName,
		ModeCookieName:    modeName,
		RefreshCookiePath: refreshPath,
	}, nil
}
