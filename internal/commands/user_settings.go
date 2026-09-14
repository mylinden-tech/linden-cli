package commands

import (
	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewUserSettingsCmd creates the user-settings command group.
func NewUserSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-settings",
		Short: "View current user settings",
		Long:  "View settings for the authenticated Linden user.",
	}
	cmd.AddCommand(newUserSettingsShowCmd())
	return cmd
}

func newUserSettingsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current user settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			settings, err := app.Client.GetUserSettings(cmd.Context())
			if err != nil {
				return err
			}
			return app.OK(settings,
				output.WithSummary("User settings"),
				output.WithBreadcrumbs(output.Breadcrumb{
					Action:      "accounts",
					Cmd:         "linden accounts list",
					Description: "List authorized accounts",
				}),
			)
		},
	}
}
