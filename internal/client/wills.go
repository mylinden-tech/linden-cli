package client

import (
	"context"
	"fmt"
	"net/url"
)

type Will struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	FromAttorney         bool    `json:"from_attorney"`
	AttorneyID           *string `json:"attorney_id,omitempty"`
	IssuedDate           *string `json:"issued_date,omitempty"`
	WhereIsIt            *string `json:"where_is_it,omitempty"`
	DocumentFileID       *string `json:"document_file_id,omitempty"`
	DocumentFileURL      *string `json:"document_file_url,omitempty"`
	Notes                *string `json:"notes,omitempty"`
	PlainLanguageSummary *string `json:"plain_language_summary,omitempty"`
	AccountID            string  `json:"account_id"`
	CreatedByID          string  `json:"created_by_id"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	Attorney             any     `json:"attorney,omitempty"`
	CreatedBy            any     `json:"created_by,omitempty"`
}

func (c *Client) ListWills(ctx context.Context, accountID string, page, size int) ([]Will, int, error) {
	path := fmt.Sprintf("/accounts/%s/wills", url.PathEscape(accountID))
	var resp struct {
		Items []Will `json:"items"`
		Total int    `json:"total"`
	}
	if err := c.get(ctx, withPageSize(path, page, size), &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

func (c *Client) GetWill(ctx context.Context, id string) (*Will, error) {
	var item Will
	if err := c.get(ctx, "/wills/"+url.PathEscape(id), &item); err != nil {
		return nil, err
	}
	return &item, nil
}
