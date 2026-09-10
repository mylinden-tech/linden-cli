package commands

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mylinden-tech/linden-cli/internal/appctx"
	"github.com/mylinden-tech/linden-cli/internal/output"
)

// isMachineOutput returns true when output should be machine-readable.
func isMachineOutput(cmd *cobra.Command) bool {
	if app := appctx.FromContext(cmd.Context()); app != nil {
		return app.IsMachineOutput()
	}
	pf := cmd.Root().PersistentFlags()
	for _, flag := range []string{"agent", "json", "quiet"} {
		if v, _ := pf.GetBool(flag); v {
			return true
		}
	}
	if f, ok := cmd.OutOrStdout().(*os.File); ok {
		fi, err := f.Stat()
		if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
			return true
		}
	}
	return false
}

// missingArg returns a usage error for a missing required argument.
func missingArg(cmd *cobra.Command, arg string) error {
	if isMachineOutput(cmd) {
		hint := "Usage: " + cmd.UseLine()
		if cmd.Example != "" {
			if first, _, ok := strings.Cut(cmd.Example, "\n"); ok {
				hint += "\nExample: " + strings.TrimSpace(first)
			} else {
				hint += "\nExample: " + strings.TrimSpace(cmd.Example)
			}
		}
		return output.ErrUsageHint(arg+" required", hint)
	}
	return cmd.Help()
}

// noChanges returns a usage error when an update command has no fields to change.
func noChanges(cmd *cobra.Command) error {
	if isMachineOutput(cmd) {
		return output.ErrUsageHint("No update fields specified", "Usage: "+cmd.UseLine())
	}
	return cmd.Help()
}

// strPtr returns a pointer to s, or nil if s is empty.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// boolPtr returns a pointer to b.
func boolPtr(b bool) *bool { return &b }

func addPageSizeFlags(cmd *cobra.Command, page, size *int) {
	cmd.Flags().IntVar(page, "page", 0, "Page number")
	cmd.Flags().IntVar(size, "size", 0, "Page size")
}

func addShareResourceFlags(cmd *cobra.Command, resourceType, resourceID *string) {
	cmd.Flags().StringVar(resourceType, "resource-type", "accounts", "Resource type")
	cmd.Flags().StringVar(resourceID, "resource-id", "", "Resource ID (defaults to active account)")
}

func resolveShareResourceID(app *appctx.App, resourceID string) (string, error) {
	if resourceID != "" {
		return resourceID, nil
	}
	if err := app.RequireAccount(); err != nil {
		return "", err
	}
	return app.Config.AccountID, nil
}

// okMaybeRedacted writes a success response, redacting fields when agent mode is on.
func okMaybeRedacted(app *appctx.App, data any, fields []string, opts ...output.ResponseOption) error {
	if app.Flags.Agent && len(fields) > 0 {
		data = output.RedactForAgent(data, fields...)
	}
	return app.OK(data, opts...)
}
