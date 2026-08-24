# linden-cli

`linden` is a command-line interface for [Linden](https://www.mylinden.family/) — the AI-powered hub families use to organize sensitive life information (legal documents, wills, insurance, digital assets) in one secure place. This CLI lets you manage your Linden accounts and the people in them from the terminal, or drive them programmatically through AI agents.

Output is designed for both humans (styled terminal tables) and machines (a consistent JSON envelope with summaries and breadcrumbs), so it works equally well typed by hand or called by an agent.

## Key features

- **OAuth 2.1 / PKCE login** against Auth0, with tokens stored in your OS keychain and refreshed automatically
- **Structured output** — styled by default in a TTY, JSON envelope otherwise, plus `--md` and `--jq <expr>` for filtering
- **Multi-account support** — switch the active account per-invocation (`--account`) or persist it globally/locally
- **Interactive TUI** for browsing persons (Bubble Tea), with forms and confirmations for creates/deletes
- **`linden doctor`** to diagnose auth, connectivity, and account configuration in one shot

## Requirements

- Go 1.25 or later (no prebuilt binaries yet — see [Installation](#installation))

## Installation

There's no published install script or package yet, so build from source:

```sh
git clone https://github.com/linden-family/linden-cli.git
cd linden-cli
make build       # builds ./bin/linden
```

Put it on your `$PATH`:

```sh
make install     # go install ./cmd/linden -> $GOBIN or $GOPATH/bin
```

Or run it directly during development without building:

```sh
go run ./cmd/linden --help
```

## Quick start

```sh
linden auth login              # opens a browser, authenticates via Auth0 (PKCE)
linden accounts list           # see which accounts you have access to
linden accounts use <id>       # set the active account
linden persons list            # browse persons in that account
linden doctor                  # sanity-check auth, API, and account config
```

## Commands

### Authentication

```sh
linden auth login     # authenticate via Auth0 (PKCE) in your browser
linden auth status    # show whether you're authenticated
linden auth logout    # remove stored credentials
```

### Accounts

```sh
linden accounts list                             # list accounts you have access to
linden accounts use <account-id>                  # set the active account (global scope)
linden accounts use <account-id> --scope local     # scope to the current directory instead
```

### Persons

```sh
linden persons list                        # list persons (opens an interactive TUI in a TTY)
linden persons show <id>                    # show a person by ID
linden persons create --first-name Jane --last-name Doe --email jane@example.com
linden persons update <id> --email new@example.com
linden persons delete <id> --yes
```

`persons create` opens an interactive form for required fields when run in a TTY without them; pass flags directly for scripting.

### Diagnostics

```sh
linden doctor   # checks auth, API connectivity, and active account configuration
```

## Output formats & flags

Every command supports these global flags:

| Flag | Description |
| --- | --- |
| `--json` | Output as JSON (envelope with `ok`/`data`/`summary`/`breadcrumbs`) |
| `--md` | Output as Markdown |
| `--styled` | Force ANSI styled output (default when stdout is a TTY) |
| `--quiet` | Output data only, no envelope |
| `--agent` | Agent mode — quiet JSON, suited for AI agents/scripts |
| `--jq <expr>` | Filter JSON output with a [jq](https://jqlang.org/) expression |
| `--account <id>` | Override the active account for this invocation only |
| `--page`, `--size` | Pagination for list commands |

By default, `linden` auto-detects: styled output in an interactive terminal, JSON envelope otherwise — so piping to another tool or an agent "just works" without extra flags.

## Configuration

Config is layered, lowest to highest precedence:

1. Built-in defaults
2. Global config — `$XDG_CONFIG_HOME/linden/config.json` (or `~/.config/linden/config.json`)
3. Local config — `.linden/config.json` in the current directory
4. Environment variables
5. Command-line flags

### Environment variables

| Variable | Purpose |
| --- | --- |
| `LINDEN_BASE_URL` | Override the API base URL |
| `LINDEN_ACCOUNT` | Override the active account ID |
| `LINDEN_TOKEN` | Bypass interactive login with a static token (useful in CI) |
| `LINDEN_NO_TUI` | Disable interactive TUI prompts/browsers (force plain output) |
| `LINDEN_NO_KEYRING` | Disable OS keyring usage for credential storage |
| `XDG_CONFIG_HOME` | Override the global config directory location |

## AI agent integration

`linden` is built for agent use as much as human use: pass `--agent` (or just pipe stdout) to get quiet, structured JSON with a `summary` field and `breadcrumbs` suggesting relevant follow-up commands — enough for an agent to navigate the CLI without needing `--help` on every step.

## Troubleshooting

Run `linden doctor` first — it checks authentication, API reachability, and whether an active account is configured, and points at the fix for whichever check fails.

## Development

```sh
make build   # build ./bin/linden
make test    # run the test suite
make tidy    # tidy go.mod/go.sum
```

### Project layout

- `cmd/linden` — entrypoint (`main.go`)
- `internal/cli` — root Cobra command and global flags
- `internal/commands` — subcommands: `auth`, `accounts`, `persons`, `doctor`
- `internal/client` — HTTP client for the Linden API
- `internal/auth` — Auth0 PKCE login, keyring storage
- `internal/config` — layered config loading (default → global → local → env → flags)
- `internal/tui` / `internal/tui/persons` — interactive TUI (forms, confirmations, persons browser)
- `internal/output` — response envelope and error formatting (JSON/Markdown/styled)
- `internal/appctx` — request-scoped app context

### Tech stack

- **Go 1.25**
- **[Cobra](https://github.com/spf13/cobra)** — command routing and flag parsing
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 / [Lip Gloss](https://github.com/charmbracelet/lipgloss) v2 / [Bubbles](https://github.com/charmbracelet/bubbles)** — terminal UI
- **[Huh](https://github.com/charmbracelet/huh)** — interactive forms
- **[gojq](https://github.com/itchyny/gojq)** — powers `--jq`
- **[go-keyring](https://github.com/zalando/go-keyring)** — OS keychain credential storage
- **Auth0** (OAuth 2.1 PKCE) — authentication against the Linden API
