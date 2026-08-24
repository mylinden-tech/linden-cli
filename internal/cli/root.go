// Package cli provides the root Cobra command for linden.
package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
	"github.com/spf13/cobra"

	"github.com/linden-family/linden-cli/internal/appctx"
	"github.com/linden-family/linden-cli/internal/commands"
	"github.com/linden-family/linden-cli/internal/config"
	"github.com/linden-family/linden-cli/internal/output"
)

// Execute runs the root command.
func Execute() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		emitError(err)
		os.Exit(exitCode(err))
	}
}

// NewRootCmd builds the root cobra command.
func NewRootCmd() *cobra.Command {
	var flags appctx.GlobalFlags

	cmd := &cobra.Command{
		Use:           "linden",
		Short:         "Command-line interface for Linden Family",
		Long:          "linden is a CLI for managing your Linden Family accounts and persons.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip setup for help/version
			if cmd.Name() == "help" || cmd.Name() == "version" {
				return nil
			}

			// Validate --jq early so errors are caught before side effects
			if flags.JQFilter != "" {
				q, err := gojq.Parse(flags.JQFilter)
				if err != nil {
					return output.ErrJQValidation(err)
				}
				if _, err := gojq.Compile(q); err != nil {
					return output.ErrJQValidation(err)
				}
			}

			cfg, err := config.Load(config.FlagOverrides{
				Account: flags.Account,
			})
			if err != nil {
				return err
			}

			app := appctx.NewApp(cfg)
			app.Flags = flags
			app.ApplyFlags()
			cmd.SetContext(appctx.WithApp(cmd.Context(), app))
			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	// Global output flags
	cmd.PersistentFlags().BoolVar(&flags.JSON, "json", false, "Output as JSON")
	cmd.PersistentFlags().BoolVar(&flags.Quiet, "quiet", false, "Output data only (no envelope)")
	cmd.PersistentFlags().BoolVar(&flags.MD, "md", false, "Output as Markdown")
	cmd.PersistentFlags().BoolVar(&flags.Styled, "styled", false, "Force ANSI styled output")
	cmd.PersistentFlags().BoolVar(&flags.Agent, "agent", false, "Agent mode (quiet JSON output)")
	cmd.PersistentFlags().StringVar(&flags.JQFilter, "jq", "", "Filter JSON output with a jq expression")

	// Context flags
	cmd.PersistentFlags().StringVar(&flags.Account, "account", "", "Override active account ID")

	// Pagination flags
	cmd.PersistentFlags().IntVar(&flags.Page, "page", 0, "Page number for list commands")
	cmd.PersistentFlags().IntVar(&flags.Size, "size", 0, "Page size for list commands")

	// Register subcommands
	cmd.AddCommand(
		commands.NewAuthCmd(),
		commands.NewAccountsCmd(),
		commands.NewPersonsCmd(),
		commands.NewDoctorCmd(),
	)

	return cmd
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	e := output.AsError(err)
	return e.ExitCode()
}

// emitError writes a structured error to stderr, matching the machine-output
// envelope when stdout is not a TTY.
func emitError(err error) {
	e := output.AsError(err)
	if isTTY(os.Stdout) {
		fmt.Fprintf(os.Stderr, "Error: %s\n", e.Message)
		if e.Hint != "" {
			fmt.Fprintf(os.Stderr, "Hint:  %s\n", e.Hint)
		}
		return
	}

	enc := json.NewEncoder(os.Stderr)
	enc.SetIndent("", "  ")
	_ = enc.Encode(output.ErrorResponse{
		OK:    false,
		Error: e.Message,
		Code:  e.Code,
		Hint:  e.Hint,
	})
}

func isTTY(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
