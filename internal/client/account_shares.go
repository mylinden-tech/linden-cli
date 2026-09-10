package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mylinden-tech/linden-cli/internal/output"
)

// AccountShare represents a resource share visible from an account.
type AccountShare struct {
	ID           string  `json:"id"`
	ResourceType string  `json:"resource_type"`
	ResourceID   string  `json:"resource_id"`
	Label        string  `json:"label"`
	Mode         string  `json:"mode"`
	Status       string  `json:"status"`
	ExpiresAt    *string `json:"expires_at"`
	AccessedAt   *string `json:"accessed_at"`
	CreatedAt    string  `json:"created_at"`
	CreatedBy    any     `json:"created_by"`
}

// AccountSharePage is a paginated account-share response.
type AccountSharePage struct {
	Items []AccountShare `json:"items"`
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Pages int            `json:"pages"`
}

// ListAccountShares returns a page of resource shares for an account.
func (c *Client) ListAccountShares(
	ctx context.Context,
	accountID string,
	page, size int,
) (*AccountSharePage, error) {
	path := fmt.Sprintf("/accounts/%s/shares", url.PathEscape(accountID))
	params := url.Values{}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if size > 0 {
		params.Set("size", fmt.Sprintf("%d", size))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var resp AccountSharePage
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindAccountShare finds a share by ID by walking the account's share pages.
func (c *Client) FindAccountShare(
	ctx context.Context,
	accountID, shareID string,
) (*AccountShare, error) {
	const size = 100
	for page := 1; ; page++ {
		result, err := c.ListAccountShares(ctx, accountID, page, size)
		if err != nil {
			return nil, err
		}
		for i := range result.Items {
			if result.Items[i].ID == shareID {
				return &result.Items[i], nil
			}
		}
		if result.Pages == 0 || page >= result.Pages {
			return nil, output.ErrNotFound("account share", shareID)
		}
	}
}
