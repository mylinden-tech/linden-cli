package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAccountReminders(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/reminders" {
			t.Errorf("request = %s %s, want GET /accounts/acc/reminders", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("size"); got != "25" {
			t.Errorf("size = %q, want 25", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{
				"id": "r1", "name": "Vet visit", "due_date": "2026-10-01",
				"entity_id": "pet1", "entity_type": "pet",
			}},
			"total": 51,
		})
	}))
	defer srv.Close()

	reminders, total, err := New(srv.URL, staticToken("t")).ListAccountReminders(
		context.Background(), "acc", 2, 25,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(reminders) != 1 || reminders[0].ID != "r1" {
		t.Fatalf("reminders = %+v", reminders)
	}
	if total != 51 {
		t.Fatalf("total = %d, want 51", total)
	}
}

func TestGetReminder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/reminders/r1" {
			t.Errorf("request = %s %s, want GET /reminders/r1", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "r1", "name": "Vet visit", "due_date": "2026-10-01",
		})
	}))
	defer srv.Close()

	reminder, err := New(srv.URL, staticToken("t")).GetReminder(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	if reminder.ID != "r1" || reminder.Name != "Vet visit" {
		t.Fatalf("reminder = %+v", reminder)
	}
}

func TestCreatePersonReminder(t *testing.T) {
	t.Parallel()

	testCreateReminder(t, "/persons/person1/reminders", func(c *Client, req ReminderCreateRequest) (*Reminder, error) {
		return c.CreatePersonReminder(context.Background(), "person1", req)
	})
}

func TestCreatePetReminder(t *testing.T) {
	t.Parallel()

	testCreateReminder(t, "/pets/pet1/reminders", func(c *Client, req ReminderCreateRequest) (*Reminder, error) {
		return c.CreatePetReminder(context.Background(), "pet1", req)
	})
}

func testCreateReminder(
	t *testing.T,
	wantPath string,
	create func(*Client, ReminderCreateRequest) (*Reminder, error),
) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != wantPath {
			t.Errorf("request = %s %s, want POST %s", r.Method, r.URL.Path, wantPath)
		}
		var body ReminderCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.Name != "Vet visit" || body.DueDate != "2026-10-01" ||
			body.RemindBeforeDays != 3 || body.ReminderType == nil ||
			*body.ReminderType != "custom" {
			t.Errorf("body = %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "r1", "name": body.Name, "due_date": body.DueDate,
		})
	}))
	defer srv.Close()

	reminderType := "custom"
	reminder, err := create(New(srv.URL, staticToken("t")), ReminderCreateRequest{
		Name:             "Vet visit",
		DueDate:          "2026-10-01",
		RemindBeforeDays: 3,
		ReminderType:     &reminderType,
	})
	if err != nil {
		t.Fatal(err)
	}
	if reminder.ID != "r1" {
		t.Fatalf("reminder = %+v", reminder)
	}
}

func TestUpdateReminder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/reminders/r1" {
			t.Errorf("request = %s %s, want PUT /reminders/r1", r.Method, r.URL.Path)
		}
		var body ReminderUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.Completed == nil || !*body.Completed {
			t.Errorf("body = %+v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "r1", "name": "Vet visit", "due_date": "2026-10-01", "completed": true,
		})
	}))
	defer srv.Close()

	completed := true
	reminder, err := New(srv.URL, staticToken("t")).UpdateReminder(
		context.Background(), "r1", ReminderUpdateRequest{Completed: &completed},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reminder.Completed {
		t.Fatalf("reminder = %+v", reminder)
	}
}

func TestDeleteReminder(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/reminders/r1" {
			t.Errorf("request = %s %s, want DELETE /reminders/r1", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
	}))
	defer srv.Close()

	if err := New(srv.URL, staticToken("t")).DeleteReminder(context.Background(), "r1"); err != nil {
		t.Fatal(err)
	}
}
