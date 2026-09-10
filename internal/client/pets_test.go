package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type staticToken string

func (s staticToken) AccessToken(context.Context) (string, error) {
	return string(s), nil
}

func TestListPets(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/accounts/acc/pets" {
			t.Errorf("path = %s, want /accounts/acc/pets", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("size"); got != "25" {
			t.Errorf("size = %q, want 25", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer t" {
			t.Errorf("authorization = %q, want Bearer t", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "p1", "name": "Rex", "species": "dog"},
			},
		})
	}))
	defer srv.Close()

	c := New(srv.URL, staticToken("t"))
	pets, err := c.ListPets(context.Background(), "acc", 2, 25)
	if err != nil {
		t.Fatal(err)
	}
	if len(pets) != 1 || pets[0].Name != "Rex" {
		t.Fatalf("pets = %+v", pets)
	}
}

func TestGetPet(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/pets/p1" {
			t.Errorf("request = %s %s, want GET /pets/p1", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "p1", "name": "Rex", "species": "dog",
		})
	}))
	defer srv.Close()

	pet, err := New(srv.URL, staticToken("t")).GetPet(context.Background(), "p1")
	if err != nil {
		t.Fatal(err)
	}
	if pet.ID != "p1" || pet.Name != "Rex" {
		t.Fatalf("pet = %+v", pet)
	}
}

func TestCreatePet(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/accounts/acc/pets" {
			t.Errorf("request = %s %s, want POST /accounts/acc/pets", r.Method, r.URL.Path)
		}
		var body PetCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.Name != "Rex" || body.Species != "dog" || !body.IsActive {
			t.Errorf("body = %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "p1", "name": body.Name, "species": body.Species, "is_active": body.IsActive,
		})
	}))
	defer srv.Close()

	pet, err := New(srv.URL, staticToken("t")).CreatePet(context.Background(), "acc", PetCreateRequest{
		Name: "Rex", Species: "dog", IsActive: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if pet.ID != "p1" || !pet.IsActive {
		t.Fatalf("pet = %+v", pet)
	}
}

func TestUpdatePet(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/pets/p1" {
			t.Errorf("request = %s %s, want PUT /pets/p1", r.Method, r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body["name"] != "Rocket" {
			t.Errorf("body = %+v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "p1", "name": body["name"], "species": "dog",
		})
	}))
	defer srv.Close()

	name := "Rocket"
	pet, err := New(srv.URL, staticToken("t")).UpdatePet(context.Background(), "p1", PetUpdateRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if pet.Name != "Rocket" {
		t.Fatalf("pet = %+v", pet)
	}
}

func TestDeletePet(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/pets/p1" {
			t.Errorf("request = %s %s, want DELETE /pets/p1", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	if err := New(srv.URL, staticToken("t")).DeletePet(context.Background(), "p1"); err != nil {
		t.Fatal(err)
	}
}
