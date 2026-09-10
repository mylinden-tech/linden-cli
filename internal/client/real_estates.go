package client

import (
	"context"
	"fmt"
	"net/url"
)

// RealEstate represents real estate in a Linden account.
type RealEstate struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`
	AddressStreet          *string `json:"address_street,omitempty"`
	AddressApt             *string `json:"address_apt,omitempty"`
	AddressCity            *string `json:"address_city,omitempty"`
	AddressState           *string `json:"address_state,omitempty"`
	AddressZip             *string `json:"address_zip,omitempty"`
	AddressCountry         *string `json:"address_country,omitempty"`
	PhotoAssetID           *string `json:"photo_asset_id,omitempty"`
	PhotoURL               *string `json:"photo_url,omitempty"`
	OwnershipType          *string `json:"ownership_type,omitempty"`
	Notes                  *string `json:"notes,omitempty"`
	IsActive               bool    `json:"is_active"`
	AccountID              string  `json:"account_id"`
	WarrantyExpirationDate *string `json:"warranty_expiration_date,omitempty"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
}

func (c *Client) ListRealEstates(ctx context.Context, accountID string, page, size int) ([]RealEstate, int, error) {
	path := fmt.Sprintf("/accounts/%s/real-estates", url.PathEscape(accountID))
	var resp struct {
		Items []RealEstate `json:"items"`
		Total int          `json:"total"`
	}
	if err := c.get(ctx, withPageSize(path, page, size), &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

func (c *Client) GetRealEstate(ctx context.Context, id string) (*RealEstate, error) {
	var item RealEstate
	if err := c.get(ctx, "/real-estates/"+url.PathEscape(id), &item); err != nil {
		return nil, err
	}
	return &item, nil
}
