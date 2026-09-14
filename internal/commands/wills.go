package commands

import (
	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func NewWillsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "wills", Aliases: []string{"will"}, Short: "View will metadata", Long: "List and view will metadata in the active account."}
	cmd.AddCommand(newWillsListCmd(), newWillsShowCmd())
	return cmd
}

func newWillsListCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "list", Short: "List will metadata", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, total, err := app.Client.ListWills(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, items, output.WillAgentOmit,
			output.WithSummary(countSummary(len(items), total, "will", "wills")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: "linden wills show <id>", Description: "View will metadata"}))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newWillsShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show <id>", Short: "Show will metadata by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		app := appctx.FromContext(cmd.Context())
		item, err := app.Client.GetWill(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, item, output.WillAgentOmit, output.WithSummary(item.Name))
	}}
}
