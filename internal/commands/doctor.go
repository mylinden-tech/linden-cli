package commands

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/linden-family/linden-cli/internal/appctx"
	"github.com/linden-family/linden-cli/internal/output"
)

// NewDoctorCmd creates the doctor command.
func NewDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check CLI configuration and connectivity",
		Long:  "Verify authentication, API reachability, and active account configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := appctx.FromContext(cmd.Context())

			checks := []map[string]any{}
			allOK := true

			// 1. Auth check
			creds, _ := app.Auth.Status()
			authOK := creds != nil
			authCheck := map[string]any{"check": "auth", "ok": authOK}
			if authOK {
				authCheck["detail"] = "Authenticated"
			} else {
				authCheck["detail"] = "Not authenticated — run: linden auth login"
				allOK = false
			}
			checks = append(checks, authCheck)

			// 2. API connectivity
			apiOK := checkAPIConnectivity(app.Config.BaseURL)
			apiCheck := map[string]any{"check": "api", "ok": apiOK, "url": app.Config.BaseURL}
			if !apiOK {
				apiCheck["detail"] = "Cannot reach API"
				allOK = false
			} else {
				apiCheck["detail"] = "Reachable"
			}
			checks = append(checks, apiCheck)

			// 3. Account
			accountOK := app.Config.AccountID != ""
			accountCheck := map[string]any{"check": "account", "ok": accountOK}
			if accountOK {
				accountCheck["detail"] = fmt.Sprintf("Active account: %s", app.Config.AccountID)
				accountCheck["account_id"] = app.Config.AccountID
			} else {
				accountCheck["detail"] = "No active account — run: linden accounts use <id>"
				allOK = false
			}
			checks = append(checks, accountCheck)

			summary := "All checks passed"
			if !allOK {
				summary = "Some checks failed"
			}

			return app.OK(map[string]any{
				"ok":     allOK,
				"checks": checks,
			}, output.WithSummary(summary))
		},
	}
}

func checkAPIConnectivity(baseURL string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/health", nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode < 500
}
