package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewInvitationsCmd creates the invitations command group.
func NewInvitationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "invitations",
		Aliases: []string{"invitation"},
		Short:   "View account invitations",
		Long:    "List pending invitations in the active account and view invitation details.",
	}
	cmd.AddCommand(newInvitationsListCmd(), newInvitationsShowCmd())
	return cmd
}

func newInvitationsListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List invitations in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			items, total, err := app.Client.ListInvitations(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}
			return app.OK(items,
				output.WithSummary(countSummary(len(items), total, "invitation", "invitations")),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "show",
					Cmd:         "linden invitations show <id>",
					Description: "View invitation details",
				}),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newInvitationsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show an invitation by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			item, err := app.Client.GetInvitation(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return app.OK(item,
				output.WithSummary(fmt.Sprintf("Invitation for %s", item.Email)),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "list",
					Cmd:         "linden invitations list",
					Description: "List invitations",
				}),
			)
		},
	}
}
