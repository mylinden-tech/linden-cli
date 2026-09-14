package commands

import (
	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func NewVehiclesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "vehicles", Aliases: []string{"vehicle"}, Short: "View vehicles", Long: "List and view vehicles in the active account."}
	cmd.AddCommand(newVehiclesListCmd(), newVehiclesShowCmd())
	return cmd
}

func newVehiclesListCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "list", Short: "List vehicles in the active account", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, total, err := app.Client.ListVehicles(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, items, output.VehicleAgentOmit,
			output.WithSummary(countSummary(len(items), total, "vehicle", "vehicles")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: "linden vehicles show <id>", Description: "View vehicle details"}))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newVehiclesShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show <id>", Short: "Show a vehicle by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		app := appctx.FromContext(cmd.Context())
		item, err := app.Client.GetVehicle(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, item, output.VehicleAgentOmit, output.WithSummary(item.Name))
	}}
}
