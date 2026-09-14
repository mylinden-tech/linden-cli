package client

import (
	"context"
	"fmt"
	"net/url"
)

// OnlineAccount represents account metadata. It deliberately has no password field.
type OnlineAccount struct {
	ID                string         `json:"id"`
	Website           string         `json:"website"`
	Platform          *string        `json:"platform,omitempty"`
	Name              string         `json:"name"`
	Username          string         `json:"username"`
	Notes             *string        `json:"notes,omitempty"`
	AccountID         string         `json:"account_id"`
	CreatedByID       string         `json:"created_by_id"`
	AdditionalDetails map[string]any `json:"additional_details,omitempty"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
	CreatedBy         any            `json:"created_by,omitempty"`
}

func (c *Client) ListOnlineAccounts(ctx context.Context, accountID string, page, size int) ([]OnlineAccount, int, error) {
	path := fmt.Sprintf("/accounts/%s/online-accounts", url.PathEscape(accountID))
	var resp struct {
		Items []OnlineAccount `json:"items"`
		Total int             `json:"total"`
	}
	if err := c.get(ctx, withPageSize(path, page, size), &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

func (c *Client) GetOnlineAccount(ctx context.Context, id string) (*OnlineAccount, error) {
	var item OnlineAccount
	if err := c.get(ctx, "/online-accounts/"+url.PathEscape(id), &item); err != nil {
		return nil, err
	}
	return &item, nil
}
