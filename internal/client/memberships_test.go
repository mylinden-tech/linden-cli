package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMemberships(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/memberships" {
			t.Errorf("request = %s %s, want GET /accounts/acc/memberships", r.Method, r.URL.Path)
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
				"id": "m1", "user_id": "u1", "account_id": "acc",
				"membership_type": "collaborator", "is_active": true,
			}},
			"total": 3,
		})
	}))
	defer srv.Close()

	items, total, err := New(srv.URL, staticToken("t")).ListMemberships(
		context.Background(), "acc", 2, 25,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != "m1" || total != 3 {
		t.Fatalf("memberships = %+v, total = %d", items, total)
	}
}

func TestGetMembership(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/memberships/m1" {
			t.Errorf("request = %s %s, want GET /memberships/m1", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "m1", "membership_type": "owner", "emergency_access": true,
		})
	}))
	defer srv.Close()

	item, err := New(srv.URL, staticToken("t")).GetMembership(context.Background(), "m1")
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != "m1" || item.MembershipType != "owner" || !item.EmergencyAccess {
		t.Fatalf("membership = %+v", item)
	}
}

func TestListMembershipRoles(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/memberships/roles" {
			t.Errorf("request = %s %s, want GET /memberships/roles", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"roles": []map[string]any{{
				"role": "owner", "description": "Full account access.",
			}},
		})
	}))
	defer srv.Close()

	roles, err := New(srv.URL, staticToken("t")).ListMembershipRoles(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 1 || roles[0].Role != "owner" {
		t.Fatalf("roles = %+v", roles)
	}
}
