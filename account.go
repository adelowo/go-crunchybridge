package gocrunchybridge

import (
	"context"
	"net/http"
	"time"
)

type AccountService service

type AuthenticatedAccount struct {
	CreatedAt            time.Time `json:"created_at"`
	DefaultTeamID        string    `json:"default_team_id"`
	Email                string    `json:"email"`
	FirstName            string    `json:"first_name"`
	HasPassword          bool      `json:"has_password"`
	HasSso               bool      `json:"has_sso"`
	ID                   string    `json:"id"`
	LastSeenAt           time.Time `json:"last_seen_at"`
	MultiFactorEnabled   bool      `json:"multi_factor_enabled"`
	Name                 string    `json:"name"`
	NotificationsEnabled bool      `json:"notifications_enabled"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (a *AccountService) User(ctx context.Context) (AuthenticatedAccount, error) {
	var resp AuthenticatedAccount

	body, err := ToReader(NoopRequestBody{})
	if err != nil {
		return resp, err
	}

	req, err := a.client.newRequest(http.MethodGet, "/account", body)
	if err != nil {
		return resp, err
	}

	_, err = a.client.Do(ctx, req, &resp)
	return resp, err
}
