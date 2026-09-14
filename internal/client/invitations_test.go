package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListInvitations(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/invitations" {
			t.Errorf("request = %s %s, want GET /accounts/acc/invitations", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("size"); got != "25" {
			t.Errorf("size = %q, want 25", got)
		}
		if got := r.URL.Query().Get("page_size"); got != "" {
			t.Errorf("page_size = %q, want empty", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"id": "i1", "email": "ada@example.com", "account_id": "acc",
				"membership_type": "collaborator",
			}},
			"total": 4,
		})
	}))
	defer srv.Close()

	items, total, err := New(srv.URL, staticToken("t")).ListInvitations(
		context.Background(), "acc", 2, 25,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "i1" || total != 4 {
		t.Fatalf("invitations = %+v, total = %d", items, total)
	}
}

func TestGetInvitation(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/invitations/i1" {
			t.Errorf("request = %s %s, want GET /invitations/i1", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "i1", "email": "ada@example.com", "membership_type": "legacy",
		})
	}))
	defer srv.Close()

	item, err := New(srv.URL, staticToken("t")).GetInvitation(context.Background(), "i1")
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "i1" || item.Email != "ada@example.com" || item.MembershipType != "legacy" {
		t.Fatalf("invitation = %+v", item)
	}
}
