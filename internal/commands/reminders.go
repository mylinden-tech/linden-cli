package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/client"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewRemindersCmd creates the reminders command group.
func NewRemindersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "reminders",
		Aliases: []string{"reminder"},
		Short:   "Manage reminders",
		Long:    "List, view, create, update, complete, and delete reminders.",
	}
	cmd.AddCommand(
		newRemindersListCmd(),
		newRemindersShowCmd(),
		newRemindersCreateCmd(),
		newRemindersUpdateCmd(),
		newRemindersCompleteCmd(),
		newRemindersDeleteCmd(),
	)
	return cmd
}

func newRemindersListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List reminders in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}

			reminders, total, err := app.Client.ListAccountReminders(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}

			count := len(reminders)
			noun := "reminders"
			if count == 1 {
				noun = "reminder"
			}
			summary := fmt.Sprintf("%d %s", count, noun)
			if total > count {
				summary = fmt.Sprintf("%s (of %d total)", summary, total)
			}
			return app.OK(reminders,
				output.WithSummary(summary),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden reminders show <id>",
						Description: "View reminder details",
					},
					output.Breadcrumb{
						Action:      "create",
						Cmd:         "linden reminders create",
						Description: "Create a reminder",
					},
				),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newRemindersShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a reminder by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			reminder, err := app.Client.GetReminder(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return app.OK(reminder,
				output.WithSummary(reminder.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "update",
						Cmd:         "linden reminders update " + reminder.ID,
						Description: "Update this reminder",
					},
					output.Breadcrumb{
						Action:      "complete",
						Cmd:         "linden reminders complete " + reminder.ID,
						Description: "Mark this reminder complete",
					},
				),
			)
		},
	}
}

func newRemindersCreateCmd() *cobra.Command {
	var (
		personID, petID                      string
		name, notes, dueDate, repeatInterval string
		reminderType                         string
		completed                            bool
		remindBeforeDays                     int
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a reminder for a person or pet",
		Long: `Create a reminder associated with exactly one person or pet.

  linden reminders create --person <uuid> --name "Renew passport" --due-date 2026-12-01`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			if (personID == "") == (petID == "") {
				return output.ErrUsage("exactly one of --person or --pet is required")
			}

			req := client.ReminderCreateRequest{
				Name:             name,
				Notes:            strPtr(notes),
				DueDate:          dueDate,
				Completed:        completed,
				RepeatInterval:   strPtr(repeatInterval),
				RemindBeforeDays: remindBeforeDays,
				ReminderType:     strPtr(reminderType),
			}

			var (
				reminder *client.Reminder
				err      error
			)
			if personID != "" {
				reminder, err = app.Client.CreatePersonReminder(cmd.Context(), personID, req)
			} else {
				reminder, err = app.Client.CreatePetReminder(cmd.Context(), petID, req)
			}
			if err != nil {
				return err
			}

			return app.OK(reminder,
				output.WithSummary("Created "+reminder.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden reminders show " + reminder.ID,
						Description: "View created reminder",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&personID, "person", "", "Person UUID")
	cmd.Flags().StringVar(&petID, "pet", "", "Pet UUID")
	cmd.Flags().StringVar(&name, "name", "", "Reminder name (required)")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD, required)")
	cmd.Flags().BoolVar(&completed, "completed", false, "Create as completed")
	cmd.Flags().StringVar(&repeatInterval, "repeat-interval", "", "Repeat interval")
	cmd.Flags().IntVar(&remindBeforeDays, "remind-before-days", 0, "Days before due date to remind")
	cmd.Flags().StringVar(&reminderType, "reminder-type", "", "Reminder type")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("due-date")
	return cmd
}

func newRemindersUpdateCmd() *cobra.Command {
	var (
		name, notes, dueDate, repeatInterval string
		reminderType                         string
		completed                            bool
		remindBeforeDays                     int
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a reminder",
		Long: `Update one or more reminder fields. Only flags you pass are sent to the API.

  linden reminders update <uuid> --due-date 2026-12-15 --remind-before-days 7`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			req := client.ReminderUpdateRequest{}
			changed := false

			if cmd.Flags().Changed("name") {
				req.Name = &name
				changed = true
			}
			if cmd.Flags().Changed("notes") {
				req.Notes = &notes
				changed = true
			}
			if cmd.Flags().Changed("due-date") {
				req.DueDate = &dueDate
				changed = true
			}
			if cmd.Flags().Changed("completed") {
				req.Completed = &completed
				changed = true
			}
			if cmd.Flags().Changed("repeat-interval") {
				req.RepeatInterval = &repeatInterval
				changed = true
			}
			if cmd.Flags().Changed("remind-before-days") {
				req.RemindBeforeDays = &remindBeforeDays
				changed = true
			}
			if cmd.Flags().Changed("reminder-type") {
				req.ReminderType = &reminderType
				changed = true
			}
			if !changed {
				return noChanges(cmd)
			}

			reminder, err := app.Client.UpdateReminder(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return app.OK(reminder,
				output.WithSummary("Updated "+reminder.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden reminders show " + reminder.ID,
						Description: "View updated reminder",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Reminder name")
	cmd.Flags().StringVar(&notes, "notes", "", "Notes")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&completed, "completed", false, "Set completion status")
	cmd.Flags().StringVar(&repeatInterval, "repeat-interval", "", "Repeat interval")
	cmd.Flags().IntVar(&remindBeforeDays, "remind-before-days", 0, "Days before due date to remind")
	cmd.Flags().StringVar(&reminderType, "reminder-type", "", "Reminder type")
	return cmd
}

func newRemindersCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "complete <id>",
		Short: "Mark a reminder complete",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			completed := true
			reminder, err := app.Client.UpdateReminder(
				cmd.Context(), args[0], client.ReminderUpdateRequest{Completed: &completed},
			)
			if err != nil {
				return err
			}
			return app.OK(reminder,
				output.WithSummary("Completed "+reminder.Name),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden reminders show " + reminder.ID,
						Description: "View completed reminder",
					},
				),
			)
		},
	}
}

func newRemindersDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a reminder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return output.ErrUsage("pass --yes to confirm deletion")
			}
			app := appctx.FromContext(cmd.Context())
			reminderID := args[0]
			if err := app.Client.DeleteReminder(cmd.Context(), reminderID); err != nil {
				return err
			}
			return app.OK(map[string]any{"id": reminderID, "deleted": true},
				output.WithSummary("Reminder deleted"),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden reminders list",
						Description: "List remaining reminders",
					},
				),
			)
		},
	}

	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion")
	return cmd
}
