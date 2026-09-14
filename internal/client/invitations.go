package client

import (
	"context"
	"fmt"
	"net/url"
)

// Invitation represents an invitation to join an account.
type Invitation struct {
	ID             string  `json:"id"`
	Email          string  `json:"email"`
	MembershipType string  `json:"membership_type"`
	Message        *string `json:"message,omitempty"`
	PersonID       *string `json:"person_id,omitempty"`
	AccountID      string  `json:"account_id"`
	InviterID      string  `json:"inviter_id"`
	ExpiresAt      string  `json:"expires_at,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
	Account        any     `json:"account,omitempty"`
	Person         any     `json:"person,omitempty"`
	Inviter        any     `json:"inviter,omitempty"`
}

// ListInvitations returns a page of account invitations and the total count.
func (c *Client) ListInvitations(
	ctx context.Context,
	accountID string,
	page, size int,
) ([]Invitation, int, error) {
	path := fmt.Sprintf("/accounts/%s/invitations", url.PathEscape(accountID))
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

	var resp struct {
		Items []Invitation `json:"items"`
		Total int          `json:"total"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// GetInvitation fetches an invitation by ID.
func (c *Client) GetInvitation(ctx context.Context, invitationID string) (*Invitation, error) {
	path := fmt.Sprintf("/invitations/%s", url.PathEscape(invitationID))
	var invitation Invitation
	if err := c.get(ctx, path, &invitation); err != nil {
		return nil, err
	}
	return &invitation, nil
}
