package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/config"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// NewAccountsCmd creates the accounts command group.
func NewAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "accounts",
		Aliases: []string{"account"},
		Short:   "Manage accounts",
		Long:    "List authorized Linden accounts and set the active account.",
	}
	cmd.AddCommand(newAccountsListCmd(), newAccountsUseCmd())
	return cmd
}

func newAccountsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List authorized accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			accounts, err := app.Client.ListAccounts(cmd.Context())
			if err != nil {
				return err
			}

			count := len(accounts)
			noun := "accounts"
			if count == 1 {
				noun = "account"
			}

			return app.OK(accounts,
				output.WithSummary(fmt.Sprintf("%d %s", count, noun)),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "use",
						Cmd:         "linden accounts use <id>",
						Description: "Set active account",
					},
				),
			)
		},
	}
}

func newAccountsUseCmd() *cobra.Command {
	var scope string

	cmd := &cobra.Command{
		Use:   "use <id>",
		Short: "Set the active account",
		Long: `Set the active Linden account for subsequent commands.

  linden accounts use <account-id>`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			if scope != "global" && scope != "local" {
				return output.ErrUsage(`--scope must be "global" or "local"`)
			}

			accountID := args[0]

			// Validate the account exists and is accessible
			accounts, err := app.Client.ListAccounts(cmd.Context())
			if err != nil {
				return err
			}

			var found bool
			var accountName string
			for _, a := range accounts {
				if a.ID == accountID {
					found = true
					accountName = a.Name
					break
				}
			}
			if !found {
				return output.ErrNotFound("account", accountID)
			}

			if err := config.PersistValue("account_id", accountID, scope); err != nil {
				return fmt.Errorf("saving account: %w", err)
			}

			summary := fmt.Sprintf("Active account set to %q (%s scope)", accountName, scope)

			return app.OK(map[string]any{
				"id":    accountID,
				"name":  accountName,
				"scope": scope,
			},
				output.WithSummary(summary),
				output.WithBreadcrumbs(
					output.Breadcrumb{
						Action:      "list",
						Cmd:         "linden persons list",
						Description: "List persons in this account",
					},
				),
			)
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "global", `Config scope: "global" or "local"`)
	return cmd
}
