package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAccountTodos(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/todos" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("page = %q, want 2", got)
		}
		if got := r.URL.Query().Get("size"); got != "25" {
			t.Errorf("size = %q, want 25", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{"id": "todo1", "title": "Call solicitor"}},
			"total": 51,
		})
	}))
	defer srv.Close()

	todos, total, err := New(srv.URL, staticToken("t")).ListAccountTodos(
		context.Background(), "acc", 2, 25,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(todos) != 1 || todos[0].ID != "todo1" {
		t.Fatalf("todos = %+v", todos)
	}
	if total != 51 {
		t.Fatalf("total = %d, want 51", total)
	}
}

func TestListTodoLists(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/accounts/acc/todo-lists" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "3" {
			t.Errorf("page = %q, want 3", got)
		}
		if got := r.URL.Query().Get("size"); got != "10" {
			t.Errorf("size = %q, want 10", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{{"id": "list1", "title": "Estate"}},
			"total": 21,
		})
	}))
	defer srv.Close()

	lists, total, err := New(srv.URL, staticToken("t")).ListTodoLists(
		context.Background(), "acc", 3, 10,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(lists) != 1 || lists[0].ID != "list1" {
		t.Fatalf("lists = %+v", lists)
	}
	if total != 21 {
		t.Fatalf("total = %d, want 21", total)
	}
}

func TestCreateTodoList(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/accounts/acc/todo-lists" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		var body TodoListCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Title != "Estate" {
			t.Errorf("body = %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "list1", "title": body.Title})
	}))
	defer srv.Close()

	list, err := New(srv.URL, staticToken("t")).CreateTodoList(
		context.Background(), "acc", TodoListCreateRequest{Title: "Estate"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if list.ID != "list1" {
		t.Fatalf("list = %+v", list)
	}
}

func TestGetTodo(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/todos/todo1" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "todo1", "title": "Call solicitor"})
	}))
	defer srv.Close()

	todo, err := New(srv.URL, staticToken("t")).GetTodo(context.Background(), "todo1")
	if err != nil {
		t.Fatal(err)
	}
	if todo.ID != "todo1" {
		t.Fatalf("todo = %+v", todo)
	}
}

func TestCreateTodo(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/todo-lists/list1/todos" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		var body TodoCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Title != "Call solicitor" {
			t.Errorf("body = %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "todo1", "title": body.Title})
	}))
	defer srv.Close()

	todo, err := New(srv.URL, staticToken("t")).CreateTodo(
		context.Background(), "list1", TodoCreateRequest{Title: "Call solicitor"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if todo.ID != "todo1" {
		t.Fatalf("todo = %+v", todo)
	}
}

func TestUpdateTodo(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/todos/todo1" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		var body TodoUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Status == nil || *body.Status != "completed" {
			t.Errorf("body = %+v", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "todo1", "title": "Call solicitor", "status": "completed",
		})
	}))
	defer srv.Close()

	status := "completed"
	todo, err := New(srv.URL, staticToken("t")).UpdateTodo(
		context.Background(), "todo1", TodoUpdateRequest{Status: &status},
	)
	if err != nil {
		t.Fatal(err)
	}
	if todo.Status != "completed" {
		t.Fatalf("todo = %+v", todo)
	}
}

func TestDeleteTodo(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/todos/todo1" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	if err := New(srv.URL, staticToken("t")).DeleteTodo(context.Background(), "todo1"); err != nil {
		t.Fatal(err)
	}
}
