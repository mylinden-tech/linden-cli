package client

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mylinden-tech/linden-cli/internal/output"
)

type ShareLink struct {
	ID         string  `json:"id"`
	URL        *string `json:"url,omitempty"`
	Mode       string  `json:"mode"`
	ExpiresAt  *string `json:"expires_at,omitempty"`
	AccessedAt *string `json:"accessed_at,omitempty"`
	CreatedAt  string  `json:"created_at"`
	Status     string  `json:"status"`
}

type ShareLinkPage struct {
	Items []ShareLink `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
	Pages int         `json:"pages"`
}

func (c *Client) ListShareLinks(ctx context.Context, resourceType, resourceID string, page, size int) (*ShareLinkPage, error) {
	path := fmt.Sprintf("/%s/%s/shares", url.PathEscape(resourceType), url.PathEscape(resourceID))
	var resp ShareLinkPage
	if err := c.get(ctx, withPageSize(path, page, size), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) FindShareLink(ctx context.Context, resourceType, resourceID, shareID string) (*ShareLink, error) {
	const size = 100
	for page := 1; ; page++ {
		result, err := c.ListShareLinks(ctx, resourceType, resourceID, page, size)
		if err != nil {
			return nil, err
		}
		for i := range result.Items {
			if result.Items[i].ID == shareID {
				return &result.Items[i], nil
			}
		}
		if result.Pages == 0 || page >= result.Pages {
			return nil, output.ErrNotFound("share link", shareID)
		}
	}
}
