package client

import (
	"context"
	"fmt"
	"net/url"
)

// GetAccount fetches an account by ID.
func (c *Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	path := fmt.Sprintf("/accounts/%s", url.PathEscape(accountID))
	var account Account
	if err := c.get(ctx, path, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetAccountStats fetches aggregate resource counts for an account.
func (c *Client) GetAccountStats(
	ctx context.Context,
	accountID string,
) (map[string]any, error) {
	path := fmt.Sprintf("/accounts/%s/stats", url.PathEscape(accountID))
	var stats map[string]any
	if err := c.get(ctx, path, &stats); err != nil {
		return nil, err
	}
	return stats, nil
}
