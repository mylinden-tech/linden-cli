package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserSettings(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/user/settings/" {
			t.Errorf("request = %s %s, want GET /user/settings/", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":                 "settings-1",
			"user_id":            "user-1",
			"default_account_id": "account-1",
			"created_at":         "2026-09-10T12:00:00Z",
		})
	}))
	defer srv.Close()

	settings, err := New(srv.URL, staticToken("t")).GetUserSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if settings.ID != "settings-1" ||
		settings.DefaultAccountID == nil ||
		*settings.DefaultAccountID != "account-1" {
		t.Fatalf("settings = %+v", settings)
	}
}
