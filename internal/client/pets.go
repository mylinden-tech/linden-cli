package client

import (
	"context"
	"fmt"
	"net/url"
)

// Pet represents a pet within an account.
type Pet struct {
	Name            string  `json:"name"`
	Species         string  `json:"species"`
	Breed           *string `json:"breed,omitempty"`
	Color           *string `json:"color,omitempty"`
	IsActive        bool    `json:"is_active"`
	Notes           *string `json:"notes,omitempty"`
	BirthDate       *string `json:"birth_date,omitempty"`
	MicrochipNumber *string `json:"microchip_number,omitempty"`
	Weight          *string `json:"weight,omitempty"`
	AccountID       string  `json:"account_id"`
	Avatar          *string `json:"avatar,omitempty"`
	CreatedByID     string  `json:"created_by_id"`
	ID              string  `json:"id"`
	AvatarURL       *string `json:"avatar_url,omitempty"`
	CreatedBy       any     `json:"created_by,omitempty"`
}

// PetCreateRequest is the payload for creating a pet.
type PetCreateRequest struct {
	Name            string  `json:"name"`
	Species         string  `json:"species"`
	Breed           *string `json:"breed,omitempty"`
	Color           *string `json:"color,omitempty"`
	IsActive        bool    `json:"is_active"`
	Notes           *string `json:"notes,omitempty"`
	BirthDate       *string `json:"birth_date,omitempty"`
	MicrochipNumber *string `json:"microchip_number,omitempty"`
	Weight          *string `json:"weight,omitempty"`
	Avatar          *string `json:"avatar,omitempty"`
}

// PetUpdateRequest is the payload for updating a pet.
type PetUpdateRequest struct {
	Name            *string `json:"name,omitempty"`
	Species         *string `json:"species,omitempty"`
	Breed           *string `json:"breed,omitempty"`
	Color           *string `json:"color,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	BirthDate       *string `json:"birth_date,omitempty"`
	MicrochipNumber *string `json:"microchip_number,omitempty"`
	Weight          *string `json:"weight,omitempty"`
	Avatar          *string `json:"avatar,omitempty"`
}

// ListPets returns pets under an account with optional pagination.
func (c *Client) ListPets(ctx context.Context, accountID string, page, size int) ([]Pet, error) {
	path := fmt.Sprintf("/accounts/%s/pets", url.PathEscape(accountID))
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
		Items []Pet `json:"items"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// GetPet fetches a single pet by ID.
func (c *Client) GetPet(ctx context.Context, petID string) (*Pet, error) {
	path := fmt.Sprintf("/pets/%s", url.PathEscape(petID))
	var pet Pet
	if err := c.get(ctx, path, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

// CreatePet creates a pet under an account.
func (c *Client) CreatePet(ctx context.Context, accountID string, req PetCreateRequest) (*Pet, error) {
	path := fmt.Sprintf("/accounts/%s/pets", url.PathEscape(accountID))
	var pet Pet
	if err := c.post(ctx, path, req, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

// UpdatePet updates a pet by ID.
func (c *Client) UpdatePet(ctx context.Context, petID string, req PetUpdateRequest) (*Pet, error) {
	path := fmt.Sprintf("/pets/%s", url.PathEscape(petID))
	var pet Pet
	if err := c.put(ctx, path, req, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

// DeletePet deletes a pet by ID.
func (c *Client) DeletePet(ctx context.Context, petID string) error {
	path := fmt.Sprintf("/pets/%s", url.PathEscape(petID))
	return c.delete(ctx, path)
}
