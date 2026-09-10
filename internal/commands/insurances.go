package commands

import (
	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

func NewInsurancesCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "insurances", Aliases: []string{"insurance"}, Short: "View insurance", Long: "List and view insurance policies and providers in the active account."}
	cmd.AddCommand(newInsurancesListCmd(), newInsurancesShowCmd(), newInsuranceProvidersCmd(), newExpiringInsurancesCmd())
	return cmd
}

func newInsurancesListCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "list", Short: "List insurance policies", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, total, err := app.Client.ListInsurances(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, items, output.InsuranceAgentOmit,
			output.WithSummary(countSummary(len(items), total, "insurance policy", "insurance policies")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: "linden insurances show <id>", Description: "View policy details"}))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newInsurancesShowCmd() *cobra.Command {
	return &cobra.Command{Use: "show <id>", Short: "Show an insurance policy by ID", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		app := appctx.FromContext(cmd.Context())
		item, err := app.Client.GetInsurance(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, item, output.InsuranceAgentOmit, output.WithSummary(item.Name))
	}}
}

func newInsuranceProvidersCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "providers", Short: "List insurance providers", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, total, err := app.Client.ListInsuranceProviders(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return app.OK(items, output.WithSummary(countSummary(len(items), total, "insurance provider", "insurance providers")))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}

func newExpiringInsurancesCmd() *cobra.Command {
	var page, size int
	cmd := &cobra.Command{Use: "expiring", Short: "List expiring insurance policies", RunE: func(cmd *cobra.Command, _ []string) error {
		app := appctx.FromContext(cmd.Context())
		if err := app.RequireAccount(); err != nil {
			return err
		}
		items, err := app.Client.ListExpiringInsurances(cmd.Context(), app.Config.AccountID, page, size)
		if err != nil {
			return err
		}
		return okMaybeRedacted(app, items, output.InsuranceAgentOmit,
			output.WithSummary(countSummary(len(items), len(items), "expiring policy", "expiring policies")),
			output.WithBreadcrumbs(output.Breadcrumb{Action: "show", Cmd: "linden insurances show <id>", Description: "View policy details"}))
	}}
	addPageSizeFlags(cmd, &page, &size)
	return cmd
}
