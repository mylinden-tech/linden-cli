// Package client provides the HTTP client for the Linden API.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/linden-family/linden-cli/internal/output"
)

// TokenProvider supplies access tokens.
type TokenProvider interface {
	AccessToken(ctx context.Context) (string, error)
}

// Client wraps the Linden REST API.
type Client struct {
	baseURL  string
	tokens   TokenProvider
	http     *http.Client
}

// New creates a new API client.
func New(baseURL string, tokens TokenProvider) *Client {
	return &Client{
		baseURL: baseURL,
		tokens:  tokens,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// --- Model types ---

// Account represents a Linden account.
type Account struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsPrimary bool   `json:"is_primary"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Person represents a person within an account.
type Person struct {
	ID               string  `json:"id"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Email            *string `json:"email,omitempty"`
	Phone            *string `json:"phone,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	Birthday         *string `json:"birthday,omitempty"`
	IsActive         bool    `json:"is_active"`
	AddressStreet    *string `json:"address_street,omitempty"`
	AddressApt       *string `json:"address_apt,omitempty"`
	AddressCity      *string `json:"address_city,omitempty"`
	AddressState     *string `json:"address_state,omitempty"`
	AddressZip       *string `json:"address_zip,omitempty"`
	AddressCountry   *string `json:"address_country,omitempty"`
	RelationshipType *string `json:"relationship_type,omitempty"`
	AvatarURL        *string `json:"avatar_url,omitempty"`
	AccountID        string  `json:"account_id"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
}

// PersonCreateRequest is the payload for creating a person.
type PersonCreateRequest struct {
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Email            *string `json:"email,omitempty"`
	Phone            *string `json:"phone,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	Birthday         *string `json:"birthday,omitempty"`
	IsActive         bool    `json:"is_active"`
	AddressStreet    *string `json:"address_street,omitempty"`
	AddressApt       *string `json:"address_apt,omitempty"`
	AddressCity      *string `json:"address_city,omitempty"`
	AddressState     *string `json:"address_state,omitempty"`
	AddressZip       *string `json:"address_zip,omitempty"`
	AddressCountry   *string `json:"address_country,omitempty"`
	RelationshipType *string `json:"relationship_type,omitempty"`
}

// PersonUpdateRequest is the payload for updating a person (all optional).
type PersonUpdateRequest struct {
	FirstName        *string `json:"first_name,omitempty"`
	LastName         *string `json:"last_name,omitempty"`
	Email            *string `json:"email,omitempty"`
	Phone            *string `json:"phone,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	Birthday         *string `json:"birthday,omitempty"`
	IsActive         *bool   `json:"is_active,omitempty"`
	AddressStreet    *string `json:"address_street,omitempty"`
	AddressApt       *string `json:"address_apt,omitempty"`
	AddressCity      *string `json:"address_city,omitempty"`
	AddressState     *string `json:"address_state,omitempty"`
	AddressZip       *string `json:"address_zip,omitempty"`
	AddressCountry   *string `json:"address_country,omitempty"`
	RelationshipType *string `json:"relationship_type,omitempty"`
}

// ListPersonsParams holds pagination options for listing persons.
type ListPersonsParams struct {
	Page int
	Size int
}

// RelationshipType is an item from the relationship-types lookup.
type RelationshipType struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// --- Account methods ---

// ListAccounts returns all accounts for the current user.
func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	var resp struct {
		Items []Account `json:"items"`
	}
	if err := c.get(ctx, "/accounts", &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// --- Person methods ---

// ListPersons returns persons under an account with optional pagination.
func (c *Client) ListPersons(ctx context.Context, accountID string, p ListPersonsParams) ([]Person, error) {
	path := fmt.Sprintf("/accounts/%s/persons", url.PathEscape(accountID))
	params := url.Values{}
	if p.Page > 0 {
		params.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.Size > 0 {
		params.Set("page_size", fmt.Sprintf("%d", p.Size))
	}
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}

	var persons []Person
	if err := c.get(ctx, path, &persons); err != nil {
		return nil, err
	}
	return persons, nil
}

// GetPerson fetches a single person by ID.
func (c *Client) GetPerson(ctx context.Context, personID string) (*Person, error) {
	path := fmt.Sprintf("/persons/%s", url.PathEscape(personID))
	var p Person
	if err := c.get(ctx, path, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// CreatePerson creates a new person under an account.
func (c *Client) CreatePerson(ctx context.Context, accountID string, req PersonCreateRequest) (*Person, error) {
	path := fmt.Sprintf("/accounts/%s/persons", url.PathEscape(accountID))
	var p Person
	if err := c.post(ctx, path, req, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdatePerson updates a person by ID.
func (c *Client) UpdatePerson(ctx context.Context, personID string, req PersonUpdateRequest) (*Person, error) {
	path := fmt.Sprintf("/persons/%s", url.PathEscape(personID))
	var p Person
	if err := c.put(ctx, path, req, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// DeletePerson deletes a person by ID.
func (c *Client) DeletePerson(ctx context.Context, personID string) error {
	path := fmt.Sprintf("/persons/%s", url.PathEscape(personID))
	return c.delete(ctx, path)
}

// ListRelationshipTypes returns available relationship types.
func (c *Client) ListRelationshipTypes(ctx context.Context) ([]RelationshipType, error) {
	var resp struct {
		Items []RelationshipType `json:"items"`
	}
	if err := c.get(ctx, "/persons/relationship-types", &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// --- HTTP helpers ---

func (c *Client) get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, body, out)
}

func (c *Client) put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, body, out)
}

func (c *Client) delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	token, err := c.tokens.AccessToken(ctx)
	if err != nil {
		return output.ErrAuth(err.Error())
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshalling request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return output.ErrNetwork(err)
	}
	defer resp.Body.Close()

	if err := checkStatus(resp); err != nil {
		return err
	}

	if out != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

// checkStatus maps HTTP error status codes to typed CLI errors.
func checkStatus(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Try to read an error body
	var apiErr struct {
		Detail  any    `json:"detail"`
		Message string `json:"message"`
	}
	rawBody, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(rawBody, &apiErr)

	msg := apiErr.Message
	if msg == "" {
		if detail, ok := apiErr.Detail.(string); ok {
			msg = detail
		}
	}
	if msg == "" {
		msg = http.StatusText(resp.StatusCode)
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return output.ErrAuth(msg)
	case http.StatusForbidden:
		return output.ErrForbidden(msg)
	case http.StatusNotFound:
		return &output.Error{Code: output.CodeNotFound, Message: msg, HTTPStatus: 404}
	case http.StatusTooManyRequests:
		return &output.Error{Code: output.CodeRateLimit, Message: "rate limited", Retryable: true}
	default:
		return output.ErrAPI(resp.StatusCode, msg)
	}
}
