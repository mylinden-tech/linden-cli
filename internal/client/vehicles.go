package client

import (
	"context"
	"fmt"
	"net/url"
)

// Vehicle represents a vehicle in a Linden account.
type Vehicle struct {
	ID                     string  `json:"id"`
	Name                   string  `json:"name"`
	VehicleType            string  `json:"vehicle_type"`
	OwnershipType          *string `json:"ownership_type,omitempty"`
	Make                   string  `json:"make"`
	Model                  string  `json:"model"`
	Year                   int     `json:"year"`
	LicensePlate           *string `json:"license_plate,omitempty"`
	VIN                    *string `json:"vin,omitempty"`
	WarrantyExpirationDate *string `json:"warranty_expiration_date,omitempty"`
	PhotoAssetID           *string `json:"photo_asset_id,omitempty"`
	PhotoURL               *string `json:"photo_url,omitempty"`
	Notes                  *string `json:"notes,omitempty"`
	IsActive               bool    `json:"is_active"`
	AccountID              string  `json:"account_id"`
	CreatedByID            string  `json:"created_by_id"`
	CreatedAt              string  `json:"created_at"`
	UpdatedAt              string  `json:"updated_at"`
	CreatedBy              any     `json:"created_by,omitempty"`
}

// ListVehicles returns a page of vehicles and its total count.
func (c *Client) ListVehicles(ctx context.Context, accountID string, page, size int) ([]Vehicle, int, error) {
	path := fmt.Sprintf("/accounts/%s/vehicles", url.PathEscape(accountID))
	path = withPageSize(path, page, size)
	var resp struct {
		Items []Vehicle `json:"items"`
		Total int       `json:"total"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// GetVehicle fetches a vehicle by ID.
func (c *Client) GetVehicle(ctx context.Context, id string) (*Vehicle, error) {
	var item Vehicle
	if err := c.get(ctx, "/vehicles/"+url.PathEscape(id), &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func withPageSize(path string, page, size int) string {
	params := url.Values{}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if size > 0 {
		params.Set("size", fmt.Sprintf("%d", size))
	}
	if len(params) > 0 {
		return path + "?" + params.Encode()
	}
	return path
}
