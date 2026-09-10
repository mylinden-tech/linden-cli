package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewAccountSharesCmd creates the account-shares command group.
func NewAccountSharesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account-shares",
		Aliases: []string{"account-share"},
		Short:   "View account resource shares",
		Long:    "List and view resource shares in the active account.",
	}
	cmd.AddCommand(newAccountSharesListCmd(), newAccountSharesShowCmd())
	return cmd
}

func newAccountSharesListCmd() *cobra.Command {
	var page, size int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List resource shares in the active account",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			result, err := app.Client.ListAccountShares(
				cmd.Context(), app.Config.AccountID, page, size,
			)
			if err != nil {
				return err
			}
			return app.OK(result.Items,
				output.WithSummary(countSummary(
					len(result.Items), result.Total, "account share", "account shares",
				)),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "show",
					Cmd:         "linden account-shares show <id>",
					Description: "View account share details",
				}),
			)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "Page number")
	cmd.Flags().IntVar(&size, "size", 0, "Page size")
	return cmd
}

func newAccountSharesShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show an account share by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.RequireAccount(); err != nil {
				return err
			}
			item, err := app.Client.FindAccountShare(
				cmd.Context(), app.Config.AccountID, args[0],
			)
			if err != nil {
				return err
			}
			return app.OK(item,
				output.WithSummary(fmt.Sprintf("%s (%s)", item.Label, item.Status)),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "list",
					Cmd:         "linden account-shares list",
					Description: "List account shares",
				}),
			)
		},
	}
}
