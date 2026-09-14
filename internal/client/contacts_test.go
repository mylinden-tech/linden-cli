package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListContacts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/contacts" {
			t.Errorf("request = %s %s, want GET /accounts/acc/contacts", r.Method, r.URL.Path)
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
			"items": []map[string]any{
				{
					"id":           "c1",
					"first_name":   "Ada",
					"last_name":    "Lovelace",
					"contact_type": "friend",
					"phone_type":   "mobile",
					"account_id":   "acc",
				},
			},
			"total": 3,
		})
	}))
	defer srv.Close()

	contacts, total, err := New(srv.URL, staticToken("t")).ListContacts(
		context.Background(), "acc", 2, 25,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(contacts) != 1 || contacts[0].FirstName != "Ada" || total != 3 {
		t.Fatalf("contacts = %+v, total = %d", contacts, total)
	}
}

func TestGetContact(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/contacts/c1" {
			t.Errorf("request = %s %s, want GET /contacts/c1", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "first_name": "Ada", "last_name": "Lovelace",
		})
	}))
	defer srv.Close()

	contact, err := New(srv.URL, staticToken("t")).GetContact(context.Background(), "c1")
	if err != nil {
		t.Fatal(err)
	}
	if contact.ID != "c1" || contact.LastName != "Lovelace" {
		t.Fatalf("contact = %+v", contact)
	}
}

func TestSearchContacts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/accounts/acc/contacts/search" {
			t.Errorf("request = %s %s, want POST /accounts/acc/contacts/search", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if len(body) != 1 || body["q"] != "Ada" {
			t.Errorf("body = %#v, want only q=Ada", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "c1", "first_name": "Ada", "last_name": "Lovelace"},
			},
			"total": 1,
		})
	}))
	defer srv.Close()

	contacts, total, err := New(srv.URL, staticToken("t")).SearchContacts(
		context.Background(), "acc", "Ada",
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(contacts) != 1 || contacts[0].ID != "c1" || total != 1 {
		t.Fatalf("contacts = %+v, total = %d", contacts, total)
	}
}

func TestListContactTypes(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/contacts/types" {
			t.Errorf("request = %s %s, want GET /contacts/types", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "attorney", "name": "Attorney"},
			},
		})
	}))
	defer srv.Close()

	types, err := New(srv.URL, staticToken("t")).ListContactTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(types) != 1 || types[0].ID != "attorney" || types[0].Name != "Attorney" {
		t.Fatalf("types = %+v", types)
	}
}
