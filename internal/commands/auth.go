package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewAuthCmd creates the auth command group.
func NewAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
		Long:  "Authenticate with the Linden API using Auth0.",
	}
	cmd.AddCommand(newAuthLoginCmd(), newAuthStatusCmd(), newAuthLogoutCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Log in to Linden",
		Long:  "Open a browser and authenticate with Auth0 using PKCE.",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			return app.Auth.Login(cmd.Context())
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			creds, err := app.Auth.Status()
			if err != nil {
				return err
			}
			if creds == nil {
				return app.OK(map[string]any{
					"authenticated": false,
				}, output.WithSummary("Not authenticated — run: linden auth login"))
			}

			data := map[string]any{
				"authenticated": true,
			}
			if !creds.ExpiresAt.IsZero() {
				data["expires_at"] = creds.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
			}
			if creds.TokenType != "" {
				data["token_type"] = creds.TokenType
			}

			return app.OK(data, output.WithSummary("Authenticated"))
		},
	}
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Log out and remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())
			if err := app.Auth.Logout(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "✓ Logged out.")
			return nil
		},
	}
}
