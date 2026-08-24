// Package appctx provides the shared application context threaded through Cobra commands.
package appctx

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/linden-family/linden-cli/internal/auth"
	"github.com/linden-family/linden-cli/internal/client"
	"github.com/linden-family/linden-cli/internal/config"
	"github.com/linden-family/linden-cli/internal/output"
)

type contextKey string

const appKey contextKey = "app"

// GlobalFlags holds values for global CLI flags.
type GlobalFlags struct {
	JSON     bool
	Quiet    bool
	MD       bool
	Styled   bool
	Agent    bool
	JQFilter string
	Account  string
	Page     int
	Size     int
}

// App holds the shared application context for all commands.
type App struct {
	Config *config.Config
	Auth   *auth.Manager
	Client *client.Client
	Output *output.Writer
	Flags  GlobalFlags
}

// authAdapter wraps auth.Manager to implement client.TokenProvider.
type authAdapter struct {
	mgr *auth.Manager
}

func (a *authAdapter) AccessToken(ctx context.Context) (string, error) {
	return a.mgr.AccessToken(ctx)
}

// NewApp creates an App from the given config.
func NewApp(cfg *config.Config) *App {
	httpClient := &http.Client{Timeout: 30 * time.Second}
	authMgr := auth.NewManager(cfg, httpClient)
	apiClient := client.New(cfg.BaseURL, &authAdapter{mgr: authMgr})

	outWriter, _ := output.New(output.DefaultOptions())

	return &App{
		Config: cfg,
		Auth:   authMgr,
		Client: apiClient,
		Output: outWriter,
	}
}

// ApplyFlags updates the output writer based on global flag values.
func (a *App) ApplyFlags() {
	opts := output.DefaultOptions()

	switch {
	case a.Flags.Agent || a.Flags.Quiet:
		opts.Format = output.FormatQuiet
	case a.Flags.JSON || a.Flags.JQFilter != "":
		opts.Format = output.FormatJSON
		opts.JQFilter = a.Flags.JQFilter
	case a.Flags.MD:
		opts.Format = output.FormatMarkdown
	case a.Flags.Styled:
		opts.Format = output.FormatStyled
	}

	// Rebuild writer with updated options.
	// Ignore error here since JQFilter was already validated in PersistentPreRunE.
	if w, err := output.New(opts); err == nil {
		a.Output = w
	}
}

// OK writes a success response.
func (a *App) OK(data any, opts ...output.ResponseOption) error {
	return a.Output.OK(data, opts...)
}

// Err writes an error response.
func (a *App) Err(err error) error {
	return a.Output.Err(err)
}

// IsInteractive returns true when running in a TTY without machine-output flags.
func (a *App) IsInteractive() bool {
	if a.Flags.Agent || a.Flags.JSON || a.Flags.Quiet || a.Flags.JQFilter != "" {
		return false
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

// IsMachineOutput returns true when output should be machine-readable.
func (a *App) IsMachineOutput() bool {
	return !a.IsInteractive()
}

// RequireAccount returns an error if no account_id is configured.
func (a *App) RequireAccount() error {
	if a.Config.AccountID == "" {
		return output.ErrUsageHint(
			"No active account",
			"Run: linden accounts list\nThen: linden accounts use <id>",
		)
	}
	return nil
}

// WithApp stores the App in the context.
func WithApp(ctx context.Context, app *App) context.Context {
	return context.WithValue(ctx, appKey, app)
}

// FromContext retrieves the App from the context.
func FromContext(ctx context.Context) *App {
	app, _ := ctx.Value(appKey).(*App)
	return app
}
