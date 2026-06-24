package auth

import (
	"context"
	"fmt"
)

// Init loads Auth0 configuration, JWT validators, and registers authentication strategies.
func Init(ctx context.Context) error {
	cfg, err := loadAuth0Config("AUTH0_")
	if err != nil {
		return err
	}

	v, err := newJWTValidator(cfg)
	if err != nil {
		return fmt.Errorf("jwt: %w", err)
	}

	Registry.Register(NewEmailRedirect(cfg, v))
	if err := Registry.SetDefault(EmailRedirectName); err != nil {
		return err
	}

	_ = ctx
	return nil
}
