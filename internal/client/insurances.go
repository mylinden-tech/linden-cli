package client

import (
	"context"
	"fmt"
	"net/url"
)

type Insurance struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Description          *string `json:"description,omitempty"`
	PolicyStatus         string  `json:"policy_status"`
	InsuranceType        string  `json:"insurance_type"`
	PolicyNumber         *string `json:"policy_number,omitempty"`
	PolicyType           string  `json:"policy_type"`
	PolicyStartDate      *string `json:"policy_start_date,omitempty"`
	PolicyExpirationDate *string `json:"policy_expiration_date,omitempty"`
	Notes                *string `json:"notes,omitempty"`
	InsuranceProviderID  string  `json:"insurance_provider_id"`
	ContactID            *string `json:"contact_id,omitempty"`
	AccountID            string  `json:"account_id"`
	CreatedByID          string  `json:"created_by_id"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
	InsuranceProvider    any     `json:"insurance_provider,omitempty"`
	Contact              any     `json:"contact,omitempty"`
	CreatedBy            any     `json:"created_by,omitempty"`
}

type InsuranceProvider struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
	CreatedByID *string `json:"created_by_id,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CreatedBy   any     `json:"created_by,omitempty"`
}

func (c *Client) ListInsurances(ctx context.Context, accountID string, page, size int) ([]Insurance, int, error) {
	path := fmt.Sprintf("/accounts/%s/insurances", url.PathEscape(accountID))
	var resp struct {
		Items []Insurance `json:"items"`
		Total int         `json:"total"`
	}
	if err := c.get(ctx, withPageSize(path, page, size), &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

func (c *Client) GetInsurance(ctx context.Context, id string) (*Insurance, error) {
	var item Insurance
	if err := c.get(ctx, "/insurances/"+url.PathEscape(id), &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *Client) ListInsuranceProviders(ctx context.Context, accountID string, page, size int) ([]InsuranceProvider, int, error) {
	path := fmt.Sprintf("/accounts/%s/insurance-providers", url.PathEscape(accountID))
	var resp struct {
		Items []InsuranceProvider `json:"items"`
		Total int                 `json:"total"`
	}
	if err := c.get(ctx, withPageSize(path, page, size), &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

func (c *Client) ListExpiringInsurances(ctx context.Context, accountID string, page, size int) ([]Insurance, error) {
	path := fmt.Sprintf("/accounts/%s/insurances/expiring", url.PathEscape(accountID))
	var items []Insurance
	if err := c.get(ctx, withPageSize(path, page, size), &items); err != nil {
		return nil, err
	}
	return items, nil
}
