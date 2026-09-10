package client

import (
	"context"
	"fmt"
	"net/url"
)

// TodoList represents a named collection of todos in an account.
type TodoList struct {
	ID           string  `json:"id"`
	AccountID    string  `json:"account_id"`
	CreatedByID  string  `json:"created_by_id"`
	CreatedBy    any     `json:"created_by,omitempty"`
	Title        string  `json:"title"`
	IsProtected  bool    `json:"is_protected"`
	ResourceType *string `json:"resource_type,omitempty"`
	ResourceID   *string `json:"resource_id,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

// Todo represents an actionable item in a todo list.
type Todo struct {
	ID           string    `json:"id"`
	AccountID    string    `json:"account_id"`
	ListID       string    `json:"list_id"`
	Status       string    `json:"status"`
	IsStarred    bool      `json:"is_starred"`
	CreatedByID  string    `json:"created_by_id"`
	CreatedBy    any       `json:"created_by,omitempty"`
	AssignedToID *string   `json:"assigned_to_id,omitempty"`
	AssignedTo   any       `json:"assigned_to,omitempty"`
	TodoList     *TodoList `json:"todo_list,omitempty"`
	Title        string    `json:"title"`
	Description  *string   `json:"description,omitempty"`
	DueDate      *string   `json:"due_date,omitempty"`
	CompletedAt  *string   `json:"completed_at,omitempty"`
	CreatedAt    string    `json:"created_at,omitempty"`
	UpdatedAt    string    `json:"updated_at,omitempty"`
}

// TodoListCreateRequest is the payload for creating a todo list.
type TodoListCreateRequest struct {
	Title        string  `json:"title"`
	ResourceType *string `json:"resource_type,omitempty"`
	ResourceID   *string `json:"resource_id,omitempty"`
}

// TodoCreateRequest is the payload for creating a todo.
type TodoCreateRequest struct {
	Title        string  `json:"title"`
	Description  *string `json:"description,omitempty"`
	DueDate      *string `json:"due_date,omitempty"`
	AssignedToID *string `json:"assigned_to_id,omitempty"`
}

// TodoUpdateRequest is the payload for updating a todo.
type TodoUpdateRequest struct {
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	DueDate      *string `json:"due_date,omitempty"`
	AssignedToID *string `json:"assigned_to_id,omitempty"`
	Status       *string `json:"status,omitempty"`
}

// ListAccountTodos returns a page of todos for an account and the total count.
func (c *Client) ListAccountTodos(
	ctx context.Context,
	accountID string,
	page, size int,
) ([]Todo, int, error) {
	path := fmt.Sprintf("/accounts/%s/todos", url.PathEscape(accountID))
	params := url.Values{}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if size > 0 {
		params.Set("size", fmt.Sprintf("%d", size))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var resp struct {
		Items []Todo `json:"items"`
		Total int    `json:"total"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// ListTodoLists returns a page of todo lists for an account and the total count.
func (c *Client) ListTodoLists(
	ctx context.Context,
	accountID string,
	page, size int,
) ([]TodoList, int, error) {
	path := fmt.Sprintf("/accounts/%s/todo-lists", url.PathEscape(accountID))
	params := url.Values{}
	if page > 0 {
		params.Set("page", fmt.Sprintf("%d", page))
	}
	if size > 0 {
		params.Set("size", fmt.Sprintf("%d", size))
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	var resp struct {
		Items []TodoList `json:"items"`
		Total int        `json:"total"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// CreateTodoList creates a todo list in an account.
func (c *Client) CreateTodoList(
	ctx context.Context,
	accountID string,
	req TodoListCreateRequest,
) (*TodoList, error) {
	path := fmt.Sprintf("/accounts/%s/todo-lists", url.PathEscape(accountID))
	var todoList TodoList
	if err := c.post(ctx, path, req, &todoList); err != nil {
		return nil, err
	}
	return &todoList, nil
}

// GetTodo fetches a todo by ID.
func (c *Client) GetTodo(ctx context.Context, todoID string) (*Todo, error) {
	path := fmt.Sprintf("/todos/%s", url.PathEscape(todoID))
	var todo Todo
	if err := c.get(ctx, path, &todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

// CreateTodo creates a todo in a list.
func (c *Client) CreateTodo(
	ctx context.Context,
	listID string,
	req TodoCreateRequest,
) (*Todo, error) {
	path := fmt.Sprintf("/todo-lists/%s/todos", url.PathEscape(listID))
	var todo Todo
	if err := c.post(ctx, path, req, &todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo updates a todo by ID.
func (c *Client) UpdateTodo(
	ctx context.Context,
	todoID string,
	req TodoUpdateRequest,
) (*Todo, error) {
	path := fmt.Sprintf("/todos/%s", url.PathEscape(todoID))
	var todo Todo
	if err := c.put(ctx, path, req, &todo); err != nil {
		return nil, err
	}
	return &todo, nil
}

// DeleteTodo deletes a todo by ID.
func (c *Client) DeleteTodo(ctx context.Context, todoID string) error {
	path := fmt.Sprintf("/todos/%s", url.PathEscape(todoID))
	return c.delete(ctx, path)
}
