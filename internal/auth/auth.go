// Package auth provides OAuth2/PKCE authentication for linden-cli.
package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mylinden-tech/linden-cli/internal/config"
)

const (
	// Auth0 application credentials (public native client — not a secret).
	auth0Domain   = "auth.mylinden.family"
	auth0ClientID = "GfAn4MFKnC57zlqIissU9x8WuS1z0KTA" //nolint:gosec // G101: public OAuth client ID //gitleaks:allow
	auth0Audience = "https://mylinden.family"

	// OAuth callback server
	callbackAddr = "127.0.0.1:3009"
	redirectURI  = "http://localhost:3009/callback"
)

// TokenEndpoints holds the discovered OIDC endpoints.
type TokenEndpoints struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
}

// Manager handles OAuth2/PKCE authentication.
type Manager struct {
	cfg        *config.Config
	store      *Store
	httpClient *http.Client
	mu         sync.Mutex
}

// NewManager creates an auth Manager.
func NewManager(cfg *config.Config, httpClient *http.Client) *Manager {
	return &Manager{
		cfg:        cfg,
		store:      NewStore(config.GlobalConfigDir()),
		httpClient: httpClient,
	}
}

// AccessToken returns a valid access token, refreshing if needed.
// If LINDEN_TOKEN env var is set it is returned directly (CI bypass).
func (m *Manager) AccessToken(ctx context.Context) (string, error) {
	if tok := os.Getenv("LINDEN_TOKEN"); tok != "" {
		return tok, nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	creds, err := m.store.Load(credentialKey())
	if err != nil {
		return "", fmt.Errorf("not authenticated — run: linden auth login")
	}

	if !creds.ExpiresAt.IsZero() && time.Until(creds.ExpiresAt) < 5*time.Minute {
		refreshed, err := m.refresh(ctx, creds)
		if err != nil {
			return "", fmt.Errorf("token refresh failed — run: linden auth login: %w", err)
		}
		creds = refreshed
	}

	return creds.AccessToken, nil
}

// Login performs the full Auth0 PKCE flow.
func (m *Manager) Login(ctx context.Context) error {
	endpoints, err := discoverEndpoints(m.httpClient)
	if err != nil {
		return fmt.Errorf("OIDC discovery failed: %w", err)
	}

	verifier, challenge, err := pkceChallenge()
	if err != nil {
		return fmt.Errorf("generating PKCE: %w", err)
	}

	state, err := randomState()
	if err != nil {
		return fmt.Errorf("generating state: %w", err)
	}

	authURL := buildAuthURL(endpoints.AuthorizationEndpoint, state, challenge)
	fmt.Printf("Opening browser for authentication...\n")
	fmt.Printf("If the browser does not open, visit:\n  %s\n", authURL)
	openBrowser(authURL)

	code, err := waitForCallback(ctx, callbackAddr, state)
	if err != nil {
		return fmt.Errorf("waiting for callback: %w", err)
	}

	creds, err := exchangeCode(ctx, m.httpClient, endpoints.TokenEndpoint, code, verifier)
	if err != nil {
		return fmt.Errorf("exchanging code: %w", err)
	}

	if err := m.store.Save(credentialKey(), creds); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	fmt.Printf("✓ Authenticated successfully.\n")
	return nil
}

// Status returns the current credentials (nil if not authenticated).
func (m *Manager) Status() (*Credentials, error) {
	creds, err := m.store.Load(credentialKey())
	if err != nil {
		return nil, nil //nolint:nilerr // not authenticated is fine
	}
	return &creds, nil
}

// Logout removes stored credentials.
func (m *Manager) Logout() error {
	return m.store.Delete(credentialKey())
}

// refresh exchanges a refresh token for new credentials.
func (m *Manager) refresh(ctx context.Context, creds Credentials) (Credentials, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {auth0ClientID},
		"refresh_token": {creds.RefreshToken},
	}

	tokenURL := fmt.Sprintf("https://%s/oauth/token", auth0Domain)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return Credentials{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return Credentials{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Credentials{}, fmt.Errorf("token refresh HTTP %d", resp.StatusCode)
	}

	var tok tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return Credentials{}, err
	}

	refreshed := Credentials{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		TokenType:    tok.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second),
	}
	if refreshed.RefreshToken == "" {
		refreshed.RefreshToken = creds.RefreshToken
	}

	if err := m.store.Save(credentialKey(), refreshed); err != nil {
		return Credentials{}, err
	}
	return refreshed, nil
}

// --- OIDC discovery ---

func discoverEndpoints(client *http.Client) (TokenEndpoints, error) {
	wellKnown := fmt.Sprintf("https://%s/.well-known/openid-configuration", auth0Domain)
	resp, err := client.Get(wellKnown) //nolint:noctx
	if err != nil {
		return TokenEndpoints{}, err
	}
	defer resp.Body.Close()

	var doc struct {
		AuthorizationEndpoint string `json:"authorization_endpoint"`
		TokenEndpoint         string `json:"token_endpoint"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return TokenEndpoints{}, err
	}
	return TokenEndpoints{
		AuthorizationEndpoint: doc.AuthorizationEndpoint,
		TokenEndpoint:         doc.TokenEndpoint,
	}, nil
}

// --- Auth URL construction ---

func buildAuthURL(endpoint, state, challenge string) string {
	params := url.Values{
		"response_type":         {"code"},
		"client_id":             {auth0ClientID},
		"redirect_uri":          {redirectURI},
		"scope":                 {"openid profile email offline_access"},
		"audience":              {auth0Audience},
		"state":                 {state},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
	}
	return endpoint + "?" + params.Encode()
}

// --- PKCE helpers ---

func pkceChallenge() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	sum := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(sum[:])
	return
}

func randomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// --- Token exchange ---

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func exchangeCode(ctx context.Context, client *http.Client, tokenURL, code, verifier string) (Credentials, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {auth0ClientID},
		"redirect_uri":  {redirectURI},
		"code":          {code},
		"code_verifier": {verifier},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL,
		strings.NewReader(data.Encode()))
	if err != nil {
		return Credentials{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return Credentials{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		return Credentials{}, fmt.Errorf("token exchange failed (%d): %s — %s",
			resp.StatusCode, errBody.Error, errBody.Description)
	}

	var tok tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return Credentials{}, err
	}

	return Credentials{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		TokenType:    tok.TokenType,
		ExpiresAt:    time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second),
	}, nil
}
