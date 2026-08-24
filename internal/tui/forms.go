// Package tui provides interactive prompts using huh and bubbletea.
package tui

import (
	"errors"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func escKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(key.WithKeys("ctrl+c", "esc"))
	return km
}

// Confirm shows a yes/no confirmation prompt.
func Confirm(message string, defaultValue bool) (bool, error) {
	var result bool
	field := huh.NewConfirm().
		Title(message).
		Affirmative("Yes").
		Negative("No").
		Value(&result)

	err := huh.NewForm(huh.NewGroup(field)).
		WithKeyMap(escKeyMap()).
		Run()
	if err != nil {
		return defaultValue, err
	}
	return result, nil
}

// ConfirmDangerous shows a confirmation for destructive actions.
func ConfirmDangerous(message string) (bool, error) {
	var result bool
	field := huh.NewConfirm().
		Title(message).
		Description("This action cannot be undone.").
		Affirmative("Yes, I'm sure").
		Negative("Cancel").
		Value(&result)

	err := huh.NewForm(huh.NewGroup(field)).
		WithKeyMap(escKeyMap()).
		Run()
	if err != nil {
		return false, err
	}
	return result, nil
}

// Input shows a text input prompt.
func Input(title, placeholder string) (string, error) {
	var result string
	field := huh.NewInput().
		Title(title).
		Placeholder(placeholder).
		Value(&result)

	err := huh.NewForm(huh.NewGroup(field)).
		WithKeyMap(escKeyMap()).
		Run()
	return result, err
}

// InputRequired shows a required text input prompt.
func InputRequired(title, placeholder string) (string, error) {
	var result string
	field := huh.NewInput().
		Title(title).
		Placeholder(placeholder).
		Value(&result).
		Validate(func(s string) error {
			if s == "" {
				return errors.New("this field is required")
			}
			return nil
		})

	err := huh.NewForm(huh.NewGroup(field)).
		WithKeyMap(escKeyMap()).
		Run()
	return result, err
}

// SelectOption is an item in a Select prompt.
type SelectOption struct {
	Value string
	Label string
}

// Select shows a single-select prompt.
func Select(title string, options []SelectOption) (string, error) {
	huhOpts := make([]huh.Option[string], len(options))
	for i, o := range options {
		huhOpts[i] = huh.NewOption(o.Label, o.Value)
	}
	var result string
	field := huh.NewSelect[string]().
		Title(title).
		Options(huhOpts...).
		Value(&result)

	err := huh.NewForm(huh.NewGroup(field)).
		WithKeyMap(escKeyMap()).
		Run()
	return result, err
}

// FormField is a field descriptor for the multi-field Form.
type FormField struct {
	Key         string
	Title       string
	Placeholder string
	Required    bool
	Default     string
}

// Form shows a multi-field input form and returns a map of key → value.
func Form(title string, fields []FormField) (map[string]string, error) {
	results := make(map[string]string)
	values := make([]*string, len(fields))
	huhFields := make([]huh.Field, len(fields))

	for i, f := range fields {
		val := f.Default
		values[i] = &val

		inp := huh.NewInput().
			Title(f.Title).
			Placeholder(f.Placeholder).
			Value(values[i])

		if f.Required {
			inp = inp.Validate(func(s string) error {
				if s == "" {
					return errors.New("this field is required")
				}
				return nil
			})
		}
		huhFields[i] = inp
	}

	form := huh.NewForm(
		huh.NewGroup(huhFields...).Title(title),
	).WithKeyMap(escKeyMap())

	if err := form.Run(); err != nil {
		return nil, err
	}

	for i, f := range fields {
		results[f.Key] = *values[i]
	}
	return results, nil
}
