package commands

import (
	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func NewOnlineAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "online-accounts", Aliases: []string{"online-account"}, Short: "View online accounts", Long: "List and view online account metadata in the active account."}
	cmd.AddCommand(newOnlineAccountsListCmd(), newOnlineAccountsShowCmd())
	return cmd
}

func newOnlineAccountsListCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "list", Short: "List online accounts", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, total, err := app.Client.ListOnlineAccounts(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, items, output.OnlineAccountAgentOmit,
			output.WithSummary(countSummary(len(items), total, "online account", "online accounts")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: "linden online-accounts show <id>", Description: "View online account metadata"}))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newOnlineAccountsShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show <id>", Short: "Show online account metadata by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		app := appctx.FromContext(cmd.Context())
		item, err := app.Client.GetOnlineAccount(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, item, output.OnlineAccountAgentOmit, output.WithSummary(item.Name))
	}}
}
