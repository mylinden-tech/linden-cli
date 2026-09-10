package client

import (
	"context"
	"fmt"
	"net/url"
)

// Reminder represents a reminder associated with an account entity.
type Reminder struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Notes            *string `json:"notes,omitempty"`
	DueDate          string  `json:"due_date"`
	Completed        bool    `json:"completed"`
	RepeatInterval   *string `json:"repeat_interval,omitempty"`
	RemindBeforeDays int     `json:"remind_before_days"`
	ReminderType     *string `json:"reminder_type,omitempty"`
	EntityID         string  `json:"entity_id,omitempty"`
	EntityType       string  `json:"entity_type,omitempty"`
	CreatedByID      string  `json:"created_by_id,omitempty"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
	CreatedBy        any     `json:"created_by,omitempty"`
	Entity           any     `json:"entity,omitempty"`
}

// ReminderCreateRequest is the shared payload for person and pet reminders.
type ReminderCreateRequest struct {
	Name             string  `json:"name"`
	Notes            *string `json:"notes,omitempty"`
	DueDate          string  `json:"due_date"`
	Completed        bool    `json:"completed"`
	RepeatInterval   *string `json:"repeat_interval,omitempty"`
	RemindBeforeDays int     `json:"remind_before_days"`
	ReminderType     *string `json:"reminder_type,omitempty"`
}

// PersonReminderCreateRequest is the API payload for a person reminder.
type PersonReminderCreateRequest = ReminderCreateRequest

// PetReminderCreate is the API payload for a pet reminder.
type PetReminderCreate = ReminderCreateRequest

// ReminderUpdateRequest is the payload for updating a reminder.
type ReminderUpdateRequest struct {
	Name             *string `json:"name,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	DueDate          *string `json:"due_date,omitempty"`
	Completed        *bool   `json:"completed,omitempty"`
	RepeatInterval   *string `json:"repeat_interval,omitempty"`
	RemindBeforeDays *int    `json:"remind_before_days,omitempty"`
	ReminderType     *string `json:"reminder_type,omitempty"`
}

// ListAccountReminders returns a page of reminders for an account and the total count.
func (c *Client) ListAccountReminders(
	ctx context.Context,
	accountID string,
	page, size int,
) ([]Reminder, int, error) {
	path := fmt.Sprintf("/accounts/%s/reminders", url.PathEscape(accountID))
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
		Items []Reminder `json:"items"`
		Total int        `json:"total"`
	}
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Items, resp.Total, nil
}

// GetReminder fetches a reminder by ID.
func (c *Client) GetReminder(ctx context.Context, reminderID string) (*Reminder, error) {
	path := fmt.Sprintf("/reminders/%s", url.PathEscape(reminderID))
	var reminder Reminder
	if err := c.get(ctx, path, &reminder); err != nil {
		return nil, err
	}
	return &reminder, nil
}

// CreatePersonReminder creates a reminder associated with a person.
func (c *Client) CreatePersonReminder(
	ctx context.Context,
	personID string,
	req ReminderCreateRequest,
) (*Reminder, error) {
	path := fmt.Sprintf("/persons/%s/reminders", url.PathEscape(personID))
	return c.createReminder(ctx, path, req)
}

// CreatePetReminder creates a reminder associated with a pet.
func (c *Client) CreatePetReminder(
	ctx context.Context,
	petID string,
	req ReminderCreateRequest,
) (*Reminder, error) {
	path := fmt.Sprintf("/pets/%s/reminders", url.PathEscape(petID))
	return c.createReminder(ctx, path, req)
}

func (c *Client) createReminder(
	ctx context.Context,
	path string,
	req ReminderCreateRequest,
) (*Reminder, error) {
	var reminder Reminder
	if err := c.post(ctx, path, req, &reminder); err != nil {
		return nil, err
	}
	return &reminder, nil
}

// UpdateReminder updates a reminder by ID.
func (c *Client) UpdateReminder(
	ctx context.Context,
	reminderID string,
	req ReminderUpdateRequest,
) (*Reminder, error) {
	path := fmt.Sprintf("/reminders/%s", url.PathEscape(reminderID))
	var reminder Reminder
	if err := c.put(ctx, path, req, &reminder); err != nil {
		return nil, err
	}
	return &reminder, nil
}

// DeleteReminder deletes a reminder by ID.
func (c *Client) DeleteReminder(ctx context.Context, reminderID string) error {
	path := fmt.Sprintf("/reminders/%s", url.PathEscape(reminderID))
	return c.delete(ctx, path)
}
