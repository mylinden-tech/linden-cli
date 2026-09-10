package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/client"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewContactsCmd creates the contacts command group.
func NewContactsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "contacts",
		Aliases: []string{"contact"},
		Short:   "View and search contacts",
		Long:    "List, view, and search contacts in the active account.",
	}
	cmd.AddCommand(
		newContactsListCmd(),
		newContactsShowCmd(),
		newContactsSearchCmd(),
		newContactsTypesCmd(),
	)
	return cmd
}

func newContactsListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List contacts in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}

			contacts, total, err := app.Client.ListContacts(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}

			return app.OK(contacts,
				output.WithSummary(contactSummary(len(contacts), total)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden contacts show <id>",
						Description: "View contact details",
					},
					output.Breadcrumb{
						Action:      "search",
						Cmd:         "linden contacts search --q <text>",
						Description: "Search contacts",
					},
				),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newContactsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a contact by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			contact, err := app.Client.GetContact(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			return app.OK(contact,
				output.WithSummary(contactDisplayName(contact)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden contacts list",
						Description: "List contacts",
					},
				),
			)
		},
	}
}

func newContactsSearchCmd() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search contacts in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}

			contacts, total, err := app.Client.SearchContacts(
				cmd.Context(), app.Config.AccountID, query,
			)
			if err != nil {
				return err
			}

			return app.OK(contacts,
				output.WithSummary(contactSummary(len(contacts), total)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "show",
						Cmd:         "linden contacts show <id>",
						Description: "View contact details",
					},
				),
			)
		},
	}
	cmd.Flags().StringVar(&query, "q", "", "Search query (required)")
	_ = cmd.MarkFlagRequired("q")
	return cmd
}

func newContactsTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List available contact types",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			types, err := app.Client.ListContactTypes(cmd.Context())
			if err != nil {
				return err
			}

			count := len(types)
			noun := "contact types"
			if count == 1 {
				noun = "contact type"
			}
			return app.OK(types,
				output.WithSummary(fmt.Sprintf("%d %s", count, noun)),
			)
		},
	}
}

func contactSummary(count, total int) string {
	noun := "contacts"
	if count == 1 {
		noun = "contact"
	}
	summary := fmt.Sprintf("%d %s", count, noun)
	if total > count {
		summary = fmt.Sprintf("%s (of %d total)", summary, total)
	}
	return summary
}

func contactDisplayName(contact *client.Contact) string {
	return strings.Join(
		strings.Fields(contact.FirstName+" "+valueOrEmpty(contact.MiddleName)+" "+contact.LastName),
		" ",
	)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
