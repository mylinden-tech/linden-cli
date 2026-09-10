package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/client"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewTodosCmd creates the todos command group.
func NewTodosCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "todos",
		Aliases: []string{"todo"},
		Short:   "Manage todos",
		Long:    "List, view, create, update, and delete todos and create todo lists.",
	}
	cmd.AddCommand(
		newTodosListCmd(),
		newTodosListsCmd(),
		newTodosCreateListCmd(),
		newTodosShowCmd(),
		newTodosCreateCmd(),
		newTodosUpdateCmd(),
		newTodosDeleteCmd(),
	)
	return cmd
}

func newTodosListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List todos in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			todos, total, err := app.Client.ListAccountTodos(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}
			count := len(todos)
			noun := "todos"
			if count == 1 {
				noun = "todo"
			}
			summary := fmt.Sprintf("%d %s", count, noun)
			if total > count {
				summary = fmt.Sprintf("%s (of %d total)", summary, total)
			}
			return app.OK(todos,
				output.WithSummary(summary),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden todos show <id>",
						Description: "View todo details",
					},
					output.Breadcrumb{
						Action:      "lists",
						Cmd:         "linden todos lists",
						Description: "List available todo lists",
					},
				),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newTodosListsCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "lists",
		Short: "List todo lists in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			lists, total, err := app.Client.ListTodoLists(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}
			count := len(lists)
			noun := "todo lists"
			if count == 1 {
				noun = "todo list"
			}
			summary := fmt.Sprintf("%d %s", count, noun)
			if total > count {
				summary = fmt.Sprintf("%s (of %d total)", summary, total)
			}
			return app.OK(lists,
				output.WithSummary(summary),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "create-list",
						Cmd:         "linden todos create-list --title <title>",
						Description: "Create a todo list",
					},
					output.Breadcrumb{
						Action:      "create",
						Cmd:         "linden todos create --list <id> --title <title>",
						Description: "Create a todo in a list",
					},
				),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newTodosCreateListCmd() *cobra.Command {
	var title, resourceType, resourceID string
	cmd := &cobra.Command{
		Use:   "create-list",
		Short: "Create a todo list in the active account",
		Long: `Create a todo list. Resource type and resource ID must be supplied together.

  linden todos create-list --title "Estate tasks"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			if (resourceType == "") != (resourceID == "") {
				return output.ErrUsage("--resource-type and --resource-id must be supplied together")
			}
			todoList, err := app.Client.CreateTodoList(
				cmd.Context(),
				app.Config.AccountID,
				client.TodoListCreateRequest{
					Title:        title,
					ResourceType: strPtr(resourceType),
					ResourceID:   strPtr(resourceID),
				},
			)
			if err != nil {
				return err
			}
			return app.OK(todoList,
				output.WithSummary("Created "+todoList.Title),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "create",
						Cmd:         "linden todos create --list " + todoList.ID + " --title <title>",
						Description: "Create a todo in this list",
					},
				),
			)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Todo list title (required)")
	cmd.Flags().StringVar(&resourceType, "resource-type", "", "Linked resource type")
	cmd.Flags().StringVar(&resourceID, "resource-id", "", "Linked resource UUID")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newTodosShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a todo by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			todo, err := app.Client.GetTodo(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return app.OK(todo,
				output.WithSummary(todo.Title),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "update",
						Cmd:         "linden todos update " + todo.ID,
						Description: "Update this todo",
					},
				),
			)
		},
	}
}

func newTodosCreateCmd() *cobra.Command {
	var listID, title, description, dueDate, assignedToID string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a todo in a list",
		Long: `Create a todo in an existing todo list.

  linden todos create --list <uuid> --title "Call solicitor"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			todo, err := app.Client.CreateTodo(
				cmd.Context(),
				listID,
				client.TodoCreateRequest{
					Title:        title,
					Description:  strPtr(description),
					DueDate:      strPtr(dueDate),
					AssignedToID: strPtr(assignedToID),
				},
			)
			if err != nil {
				return err
			}
			return app.OK(todo,
				output.WithSummary("Created "+todo.Title),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden todos show " + todo.ID,
						Description: "View created todo",
					},
				),
			)
		},
	}
	cmd.Flags().StringVar(&listID, "list", "", "Todo list UUID (required)")
	cmd.Flags().StringVar(&title, "title", "", "Todo title (required)")
	cmd.Flags().StringVar(&description, "description", "", "Description")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&assignedToID, "assigned-to", "", "Assigned user UUID")
	_ = cmd.MarkFlagRequired("list")
	_ = cmd.MarkFlagRequired("title")
	return cmd
}

func newTodosUpdateCmd() *cobra.Command {
	var title, description, dueDate, assignedToID, status string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a todo",
		Long: `Update one or more todo fields. Only flags you pass are sent to the API.

  linden todos update <uuid> --status completed`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			req := client.TodoUpdateRequest{}
			changed := false
			if cmd.Flags().Changed("title") {
				req.Title = &title
				changed = true
			}
			if cmd.Flags().Changed("description") {
				req.Description = &description
				changed = true
			}
			if cmd.Flags().Changed("due-date") {
				req.DueDate = &dueDate
				changed = true
			}
			if cmd.Flags().Changed("assigned-to") {
				req.AssignedToID = &assignedToID
				changed = true
			}
			if cmd.Flags().Changed("status") {
				req.Status = &status
				changed = true
			}
			if !changed {
				return noChanges(cmd)
			}
			todo, err := app.Client.UpdateTodo(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return app.OK(todo,
				output.WithSummary("Updated "+todo.Title),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden todos show " + todo.ID,
						Description: "View updated todo",
					},
				),
			)
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Todo title")
	cmd.Flags().StringVar(&description, "description", "", "Description")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&assignedToID, "assigned-to", "", "Assigned user UUID")
	cmd.Flags().StringVar(&status, "status", "", "Status (open or completed)")
	return cmd
}

func newTodosDeleteCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a todo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return output.ErrUsage("pass --yes to confirm deletion")
			}
			app := appctx.FromContext(cmd.Context())
			todoID := args[0]
			if err := app.Client.DeleteTodo(cmd.Context(), todoID); err != nil {
				return err
			}
			return app.OK(map[string]any{"id": todoID, "deleted": true},
				output.WithSummary("Todo deleted"),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden todos list",
						Description: "List remaining todos",
					},
				),
			)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm deletion")
	return cmd
}
