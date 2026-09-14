package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAccountShares(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/shares" {
			t.Errorf("request = %s %s, want GET /accounts/acc/shares", r.Method, r.URL.Path)
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
				"id": "s1", "resource_type": "will", "resource_id": "r1",
				"label": "Executor copy", "mode": "view", "status": "active",
				"created_at": "2026-09-10T12:00:00Z",
			}},
			"total": 5, "page": 2, "size": 25, "pages": 1,
		})
	}))
	defer srv.Close()

	result, err := New(srv.URL, staticToken("t")).ListAccountShares(
		context.Background(), "acc", 2, 25,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].ID != "s1" ||
		result.Total != 5 || result.Page != 2 || result.Size != 25 || result.Pages != 1 {
		t.Fatalf("page = %+v", result)
	}
}

func TestFindAccountShare(t *testing.T) {
	t.Parallel()

	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/shares" {
			t.Errorf("request = %s %s, want GET /accounts/acc/shares", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("size"); got != "100" {
			t.Errorf("size = %q, want 100", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"id": "s1", "resource_type": "will", "resource_id": "r1",
				"label": "Executor copy", "mode": "view", "status": "active",
				"created_at": "2026-09-10T12:00:00Z",
			}},
			"total": 1, "page": 1, "size": 100, "pages": 1,
		})
	}))
	defer srv.Close()

	item, err := New(srv.URL, staticToken("t")).FindAccountShare(
		context.Background(), "acc", "s1",
	)
	if err != nil {
		t.Fatal(err)
	}
	if requests != 1 || item.ID != "s1" || item.Label != "Executor copy" {
		t.Fatalf("requests = %d, share = %+v", requests, item)
	}
}

func TestFindAccountShareNotFound(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []any{}, "total": 0, "page": 1, "size": 100, "pages": 0,
		})
	}))
	defer srv.Close()

	_, err := New(srv.URL, staticToken("t")).FindAccountShare(
		context.Background(), "acc", "missing",
	)
	if err == nil {
		t.Fatal("expected not found error")
	}
}
