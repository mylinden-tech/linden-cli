package client

import (
	"context"
	"fmt"
	"net/url"
)

// Membership represents a user's membership in an account.
type Membership struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	AccountID       string `json:"account_id"`
	MembershipType  string `json:"membership_type"`
	IsActive        bool   `json:"is_active"`
	EmergencyAccess bool   `json:"emergency_access"`
	CreatedByID     string `json:"created_by_id,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	User            any    `json:"user,omitempty"`
	CreatedBy       any    `json:"created_by,omitempty"`
}

// MembershipRole describes an available account membership role.
type MembershipRole struct {
	Role        string `json:"role"`
	Description string `json:"description"`
}

// ListMemberships returns a page of account memberships and the total count.
func (c *Client) ListMemberships(
	ctx context.Context,
	accountID string,
	page, size int,
) ([]Membership, int, error) {
	path := fmt.Sprintf("/accounts/%s/memberships", url.PathEscape(accountID))
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
		Items []Membership `json:"items"`
		Total int          `json:"total"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// GetMembership fetches a membership by ID.
func (c *Client) GetMembership(ctx context.Context, membershipID string) (*Membership, error) {
	path := fmt.Sprintf("/memberships/%s", url.PathEscape(membershipID))
	var membership Membership
	if err := c.get(ctx, path, &membership); err != nil {
		return nil, err
	}
	return &membership, nil
}

// ListMembershipRoles returns all available membership roles.
func (c *Client) ListMembershipRoles(ctx context.Context) ([]MembershipRole, error) {
	var resp struct {
		Roles []MembershipRole `json:"roles"`
	}
	if err := c.get(ctx, "/memberships/roles", &resp); err != nil {
		return nil, err
	}
	return resp.Roles, nil
}
