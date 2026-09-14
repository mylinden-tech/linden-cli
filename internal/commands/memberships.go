package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewMembershipsCmd creates the memberships command group.
func NewMembershipsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "memberships",
		Aliases: []string{"membership"},
		Short:   "View account memberships",
		Long:    "List and view memberships in the active account, or inspect available roles.",
	}
	cmd.AddCommand(
		newMembershipsListCmd(),
		newMembershipsShowCmd(),
		newMembershipsRolesCmd(),
	)
	return cmd
}

func newMembershipsListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List memberships in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			items, total, err := app.Client.ListMemberships(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}
			return app.OK(items,
				output.WithSummary(countSummary(len(items), total, "membership", "memberships")),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "show",
					Cmd:         "linden memberships show <id>",
					Description: "View membership details",
				}),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newMembershipsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show a membership by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			item, err := app.Client.GetMembership(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return app.OK(item,
				output.WithSummary(fmt.Sprintf("%s membership", item.MembershipType)),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "list",
					Cmd:         "linden memberships list",
					Description: "List memberships",
				}),
			)
		},
	}
}

func newMembershipsRolesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "roles",
		Short: "List available membership roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			roles, err := app.Client.ListMembershipRoles(cmd.Context())
			if err != nil {
				return err
			}
			return app.OK(roles,
				output.WithSummary(countSummary(len(roles), len(roles), "role", "roles")),
			)
		},
	}
}

func countSummary(count, total int, singular, plural string) string {
	noun := plural
	if count == 1 {
		noun = singular
	}
	summary := fmt.Sprintf("%d %s", count, noun)
	if total > count {
		summary = fmt.Sprintf("%s (of %d total)", summary, total)
	}
	return summary
}
