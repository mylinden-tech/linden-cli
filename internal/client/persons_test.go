package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPersonsUsesSizeQueryParameter(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/persons" {
			t.Errorf("request = %s %s, want GET /accounts/acc/persons", r.Method, r.URL.Path)
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
				"id": "person1", "first_name": "Ada", "last_name": "Lovelace",
			}},
		})
	}))
	defer srv.Close()

	persons, err := New(srv.URL, staticToken("t")).ListPersons(
		context.Background(),
		"acc",
		ListPersonsParams{Page: 2, Size: 25},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(persons) != 1 || persons[0].ID != "person1" {
		t.Fatalf("persons = %+v", persons)
	}
}
