package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/zalando/go-keyring"

	"github.com/mylinden-tech/linden-cli/internal/config"
)

const keyringService = "linden"

// Credentials holds OAuth tokens.
type Credentials struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// Store persists credentials to the system keyring with a file fallback.
type Store struct {
	dir      string
	warnOnce sync.Once
}

// NewStore creates a credential store backed by dir as fallback.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Save persists credentials for the given origin key.
func (s *Store) Save(origin string, creds Credentials) error {
	data, err := json.Marshal(creds)
	if err != nil {
		return fmt.Errorf("marshalling credentials: %w", err)
	}

	if os.Getenv("LINDEN_NO_KEYRING") == "" {
		if err := keyring.Set(keyringService, origin, string(data)); err == nil {
			return nil
		}
	}

	// Fallback: file
	s.warnOnce.Do(func() {
		fmt.Fprintf(os.Stderr, "Warning: system keyring unavailable — storing credentials in %s\n", s.dir)
	})
	return s.saveToFile(origin, data)
}

// Load retrieves credentials for the given origin key.
func (s *Store) Load(origin string) (Credentials, error) {
	if os.Getenv("LINDEN_NO_KEYRING") == "" {
		if raw, err := keyring.Get(keyringService, origin); err == nil {
			var creds Credentials
			if err := json.Unmarshal([]byte(raw), &creds); err == nil {
				return creds, nil
			}
		}
	}

	return s.loadFromFile(origin)
}

// Delete removes credentials for the given origin key.
func (s *Store) Delete(origin string) error {
	_ = keyring.Delete(keyringService, origin) // best effort
	_ = os.Remove(s.credFilePath(origin))
	return nil
}

func (s *Store) saveToFile(origin string, data []byte) error {
	if err := os.MkdirAll(s.dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(s.credFilePath(origin), data, 0600)
}

func (s *Store) loadFromFile(origin string) (Credentials, error) {
	raw, err := os.ReadFile(s.credFilePath(origin)) //nolint:gosec
	if err != nil {
		return Credentials{}, fmt.Errorf("no credentials found — run: linden auth login")
	}
	var creds Credentials
	if err := json.Unmarshal(raw, &creds); err != nil {
		return Credentials{}, fmt.Errorf("corrupt credentials: %w", err)
	}
	return creds, nil
}

func (s *Store) credFilePath(origin string) string {
	// Sanitize origin to a safe filename
	safe := filepath.Clean(origin)
	for _, r := range []string{"/", ":", "\\", "?"} {
		safe = replaceAll(safe, r, "_")
	}
	return filepath.Join(s.dir, "credentials_"+safe+".json")
}

func replaceAll(s, old, new string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			out = append(out, new...)
			i += len(old)
		} else {
			out = append(out, s[i])
			i++
		}
	}
	return string(out)
}

// credentialKey returns the storage key used for this CLI.
func credentialKey() string {
	return config.GlobalConfigDir()
}
