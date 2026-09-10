package client

import (
	"context"
	"fmt"
	"net/url"
)

// Contact represents a contact within an account.
type Contact struct {
	ID           string  `json:"id"`
	FirstName    string  `json:"first_name"`
	MiddleName   *string `json:"middle_name,omitempty"`
	LastName     string  `json:"last_name"`
	Company      *string `json:"company,omitempty"`
	Job          *string `json:"job,omitempty"`
	ContactType  string  `json:"contact_type"`
	PhoneType    string  `json:"phone_type"`
	Phone        *string `json:"phone,omitempty"`
	Email        *string `json:"email,omitempty"`
	Website      *string `json:"website,omitempty"`
	AddressLine1 *string `json:"address_line_1,omitempty"`
	AddressLine2 *string `json:"address_line_2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	ZipCode      *string `json:"zip_code,omitempty"`
	Country      *string `json:"country,omitempty"`
	Notes        *string `json:"notes,omitempty"`
	IsActive     bool    `json:"is_active"`
	AccountID    string  `json:"account_id"`
	CreatedByID  string  `json:"created_by_id"`
	Avatar       *string `json:"avatar,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
	AvatarURL    *string `json:"avatar_url,omitempty"`
	CreatedBy    any     `json:"created_by,omitempty"`
}

// ContactType represents an available contact classification.
type ContactType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type contactPage struct {
	Items []Contact `json:"items"`
	Total int       `json:"total"`
}

// ListContacts returns a page of contacts for an account and the total count.
func (c *Client) ListContacts(
	ctx context.Context,
	accountID string,
	page, size int,
) ([]Contact, int, error) {
	path := fmt.Sprintf("/accounts/%s/contacts", url.PathEscape(accountID))
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

	var resp contactPage
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// GetContact fetches a contact by ID.
func (c *Client) GetContact(ctx context.Context, contactID string) (*Contact, error) {
	path := fmt.Sprintf("/contacts/%s", url.PathEscape(contactID))
	var contact Contact
	if err := c.get(ctx, path, &contact); err != nil {
		return nil, err
	}
	return &contact, nil
}

// SearchContacts searches contacts within an account.
func (c *Client) SearchContacts(
	ctx context.Context,
	accountID, query string,
) ([]Contact, int, error) {
	path := fmt.Sprintf("/accounts/%s/contacts/search", url.PathEscape(accountID))
	req := struct {
		Query string `json:"q"`
	}{Query: query}

	var resp contactPage
	if err := c.post(ctx, path, req, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// ListContactTypes returns the available contact classifications.
func (c *Client) ListContactTypes(ctx context.Context) ([]ContactType, error) {
	var resp struct {
		Items []ContactType `json:"items"`
	}
	if err := c.get(ctx, "/contacts/types", &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}
