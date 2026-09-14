package client

import "context"

// UserSettings represents preferences for the current user.
type UserSettings struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	DefaultAccountID *string `json:"default_account_id"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
}

// GetUserSettings fetches settings for the current user.
func (c *Client) GetUserSettings(ctx context.Context) (*UserSettings, error) {
	var settings UserSettings
	if err := c.get(ctx, "/user/settings/", &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}
