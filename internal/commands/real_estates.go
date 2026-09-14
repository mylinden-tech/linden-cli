package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func NewRealEstatesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "real-estates", Aliases: []string{"real-estate"}, Short: "View real estate", Long: "List and view real estate in the active account."}
	cmd.AddCommand(newRealEstatesListCmd(), newRealEstatesShowCmd())
	return cmd
}

func newRealEstatesListCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "list", Short: "List real estate in the active account", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, total, err := app.Client.ListRealEstates(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return app.OK(items, output.WithSummary(countSummary(len(items), total, "property", "properties")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: "linden real-estates show <id>", Description: "View real estate details"}))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newRealEstatesShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show <id>", Short: "Show real estate by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		app := appctx.FromContext(cmd.Context())
		item, err := app.Client.GetRealEstate(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		location := strings.Trim(strings.Join([]string{valueOrEmpty(item.AddressCity), valueOrEmpty(item.AddressState)}, ", "), ", ")
		summary := item.Name
		if location != "" {
			summary += " — " + location
		}
		return app.OK(item, output.WithSummary(summary))
	}}
}
