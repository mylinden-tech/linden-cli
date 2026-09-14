package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccount(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.EscapedPath() != "/accounts/account%2Fid" {
			t.Errorf("request = %s %s, want GET /accounts/account%%2Fid", r.Method, r.URL.EscapedPath())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "account/id", "name": "Family", "description": "Estate",
			"is_primary": true, "onboard": false,
		})
	}))
	defer srv.Close()

	account, err := New(srv.URL, staticToken("t")).GetAccount(
		context.Background(), "account/id",
	)
	if err != nil {
		t.Fatal(err)
	}
	if account.ID != "account/id" || account.Name != "Family" || account.Description == nil {
		t.Fatalf("account = %+v", account)
	}
}

func TestGetAccountStats(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/stats" {
			t.Errorf("request = %s %s, want GET /accounts/acc/stats", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"persons_count": 3, "pets_count": 2,
		})
	}))
	defer srv.Close()

	stats, err := New(srv.URL, staticToken("t")).GetAccountStats(
		context.Background(), "acc",
	)
	if err != nil {
		t.Fatal(err)
	}
	if stats["persons_count"] != float64(3) || stats["pets_count"] != float64(2) {
		t.Fatalf("stats = %#v", stats)
	}
}
