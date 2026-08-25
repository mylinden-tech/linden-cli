# PRD 0001 — Linden CLI (`linden`)

## Problem Statement

Users and agents need a structured, scriptable way to interact with the Linden Family API. The API manages family accounts and the persons within them — but interacting with it today requires constructing raw HTTP requests with OAuth tokens, which is error-prone and not agent-friendly.

There is no CLI tool for the Linden API today. Users working in terminal environments have no ergonomic way to list accounts, switch between them, or manage the persons within an account. This blocks automation workflows and makes it tedious to perform common tasks without opening a browser.

A secondary challenge is that the API is inherently multi-account: all person data is scoped under an `account_id`, so any CLI must give users a clear, low-friction way to declare which account they are working in.

## Solution

Build `linden` — the Linden CLI — a Go command-line tool that wraps the Linden Family API. It enables users and agents to:

- Authenticate once via a browser-based OAuth2/PKCE flow (Auth0) with automatic token refresh
- List available accounts and set a persistent active account (`linden accounts use <id>`)
- Manage persons within the active account: list, show, create, update, and delete
- Get machine-readable JSON output with breadcrumbs for agent chaining
- Verify their setup with a single `linden doctor` command

The CLI is modeled structurally and stylistically after the basecamp-cli, using Go + Cobra with the same `internal/` package layout, output envelope format, and agent integration conventions.

## User Stories

### Auth

1. As a user, I want to run `linden auth login` to open a browser and authenticate with Auth0, so that I can obtain a token without handling OAuth manually.
2. As a user, I want the CLI to automatically refresh my token when it expires, so that long-running workflows are not interrupted.
3. As a developer, I want to run `linden auth status` to verify whether I am authenticated and see token expiry, so that I can diagnose auth issues quickly.
4. As a developer, I want to run `linden auth logout` to revoke and remove stored credentials, so that I can switch accounts or clean up my environment.

### Accounts

5. As a user, I want to run `linden accounts list` to see all accounts I have access to, so that I can find the account ID I need to work with.
6. As a user, I want to run `linden accounts use <id>` to set a default account, so that I do not have to specify `--account` on every subsequent command.
7. As a developer, I want `linden accounts use` to support `--scope global` (default) and `--scope local` so that I can set a project-specific account without affecting my global config.

### Persons (Family members)

8. As a user, I want to run `linden persons list` to list all persons in my active account, so that I can browse or pipe the output to other tools.
9. As a user, I want to run `linden persons show <id>` to retrieve a single person by UUID, so that I can inspect their full profile including address, birthday, and relationship type.
10. As a user, I want to run `linden persons create` and be prompted for first name, last name, and optional fields (email, phone, notes, birthday, address, relationship type) via an interactive form, so that I can add a new person without remembering all flag names.
11. As a developer, I want to run `linden persons create --first-name "Jane" --last-name "Doe" --email jane@example.com` non-interactively, so that I can script person creation without a TTY.
12. As a user, I want to run `linden persons update <id>` with one or more flags (e.g., `--email`, `--phone`, `--notes`) to update specific fields on a person, so that I can make targeted edits without re-supplying unchanged data.
13. As a user, I want to run `linden persons delete <id>` and be prompted to confirm before the record is removed, so that I do not accidentally delete someone.
14. As a developer, I want to pass `--yes` to `linden persons delete` to skip the confirmation prompt, so that I can automate deletions in scripts.
15. As a user, I want pagination flags (`--page`, `--size`) on `linden persons list`, so that I can page through large contact lists.
16. As an agent, I want all commands to support `--json` to get a structured JSON envelope (`ok`, `data`, `breadcrumbs`), so that I can pipe output to other tools reliably.
17. As a developer, I want all commands to support `--jq <expr>` for inline filtering of JSON output, so that I do not need a separate `jq` binary.
18. As a developer, I want all commands to support `--md` to get Markdown-formatted output, so that I can include results in documents or chat messages.
19. As an agent, I want JSON responses to include `breadcrumbs` with suggested follow-up commands, so that I can chain operations without reasoning about the next step from scratch.
20. As an agent, I want to run `linden doctor` to check auth status, API connectivity, and active account config in one command, so that I can quickly diagnose whether the CLI is correctly configured.
21. As a developer, I want to set `LINDEN_BASE_URL` to point the CLI at a local or staging API, so that I can develop and test without hitting production.
22. As an agent, I want the error envelope to include a `code` and `hint` field when a command fails, so that I can understand what went wrong and how to fix it.
23. As a developer, I want to run `linden --help --agent` to get structured JSON describing all available commands and flags, so that agents can self-discover the CLI surface without reading documentation.

## Implementation Decisions

### Module breakdown

**`internal/cli`** — Root Cobra command, global flags (`--json`, `--jq`, `--md`, `--agent`, `--page`, `--size`), help system with `--agent` JSON output mode. Command registration hub.

**`internal/auth`** — OAuth2/PKCE flow against Auth0. Handles browser launch, local loopback callback server on port 3009, token exchange, token storage, and automatic refresh. Exposes a simple `AccessToken() (string, error)` interface consumed by the API client.

**`internal/config`** — Layered configuration: env vars → local config file → global config file → defaults. Reads `LINDEN_BASE_URL` (default `https://api.mylinden.family`), `LINDEN_TOKEN` (bypass auth for scripts), `LINDEN_NO_KEYRING`, and `account_id` (set via `accounts use`). Config file at `~/.config/linden/config.json`. Credentials at `~/.config/linden/credentials.json` (fallback from system keyring).

**`internal/client`** — HTTP client wrapping `net/http`. Attaches Bearer token from auth module, sets base URL from config, handles 4xx/5xx responses and maps them to typed errors. Exposes typed methods: `ListAccounts`, `ListPersons`, `GetPerson`, `CreatePerson`, `UpdatePerson`, `DeletePerson`.

**`internal/commands`** — One file per resource group (`accounts.go`, `persons.go`, `auth.go`, `doctor.go`). Each exposes a `NewXxxCmd() *cobra.Command` constructor. Commands call the client and pass results to the output module.

**`internal/output`** — JSON envelope (`ok`, `data`, `summary`, `breadcrumbs`, `error`, `hint`, `meta`). Format detection (TTY → styled, non-TTY → JSON). `--json`, `--jq` (via gojq), `--md` rendering. Error envelope on failure.

**`internal/presenter`** — Human-readable rendering for `--md` and styled TTY output. Presenters for account rows and person detail/list views, renders tables and detail cards.

**`internal/tui`** — Interactive UX layer. `forms.go` wraps huh for confirmations and input prompts; `persons/` holds the bubbletea model for the interactive persons browser.

### Auth flow detail

- `linden auth login` starts a local HTTP server on `http://localhost:3009/callback`, constructs the Auth0 authorize URL with PKCE challenge, opens the browser, waits for the callback, exchanges the code for tokens, and persists them.
- Auth0 domain: `auth.mylinden.family`
- Auth0 client ID: `GfAn4MFKnC57zlqIissU9x8WuS1z0KTA` <!-- public native client, not a secret; gitleaks:allow -->
- Redirect URI: `http://localhost:3009/callback`
- Token storage: system keyring preferred, `~/.config/linden/credentials.json` as fallback.
- `LINDEN_TOKEN` env var bypasses the entire flow for CI/service account use.

### Account context

- `linden accounts use <id>` validates the ID against `GET /accounts`, then persists `account_id` to the config file at the chosen scope (`--scope global` writes to `~/.config/linden/config.json`; `--scope local` writes to `.linden/config.json` in the current directory).
- All persons commands read `account_id` from the resolved config. If no account is set, commands print a clear error: `No active account. Run: linden accounts use <id>`.
- Config resolution order for `account_id`: `LINDEN_ACCOUNT` env var → local `.linden/config.json` → global `~/.config/linden/config.json`.

### Persons command flags

**`linden persons create`** flags (all optional except `--first-name` and `--last-name`):

| Flag                  | Type                | API field           |
| --------------------- | ------------------- | ------------------- |
| `--first-name`        | string              | `first_name`        |
| `--last-name`         | string              | `last_name`         |
| `--email`             | string              | `email`             |
| `--phone`             | string              | `phone`             |
| `--notes`             | string              | `notes`             |
| `--birthday`          | string (YYYY-MM-DD) | `birthday`          |
| `--address-street`    | string              | `address_street`    |
| `--address-apt`       | string              | `address_apt`       |
| `--address-city`      | string              | `address_city`      |
| `--address-state`     | string              | `address_state`     |
| `--address-zip`       | string              | `address_zip`       |
| `--address-country`   | string              | `address_country`   |
| `--relationship-type` | string              | `relationship_type` |

When running in a TTY without `--json`/`--jq`/`--md`, and required flags are missing, `create` launches a huh form to collect them interactively.

**`linden persons update <id>`** accepts the same flags as `create`, all optional. Only flags explicitly passed are sent in the `PUT` body.

**`linden persons delete <id>`** prompts for confirmation via huh unless `--yes` is passed.

### Output envelope

All commands emit a consistent JSON structure when `--json` is set or output is non-TTY:

```
{ "ok": true, "data": <resource or array>, "summary": "...", "breadcrumbs": [...] }
{ "ok": false, "error": "message", "code": "ERROR_CODE", "hint": "suggestion" }
```

### Breadcrumbs

- `accounts list` → suggests `linden accounts use <id>`
- `accounts use` → suggests `linden persons list`
- `persons list` → suggests `linden persons show <id>` for the first result
- `persons show` → suggests `linden persons update <id>` and `linden persons delete <id>`
- `persons create` → suggests `linden persons show <id>` for the newly created record

### TUI (bubbletea / huh)

The CLI follows the same layered TUI pattern as `basecamp-cli`:

- **huh** handles simple blocking interactions inside Cobra commands: confirmations (`persons delete`), single-field input prompts, and selection lists. Thin wrappers in `internal/tui/forms.go` (`Confirm`, `Input`, `Select`). Commands call them inline within `RunE`.
- **bubbletea v2** drives the interactive persons browser: a `linden persons list` invocation without `--json` / non-TTY flags launches a full-screen TUI instead of printing raw output.

**Interactive persons TUI** (`linden persons list` in a TTY without `--json`):

- Text input field at the top with debounced filtering (400 ms) as the user types.
- Scrollable results list below; each row shows name, relationship type, email, and birthday if set.
- Arrow keys / `j` / `k` to navigate; `Enter` to inspect the selected person in a detail pane; `n` to create a new person (launches huh form); `e` to edit the selected person; `d` to delete (with confirmation); `Esc` / `q` to quit.
- Selecting a person and pressing `Enter` prints their JSON to stdout (for agent chaining) after the program exits.

**Package layout additions:**

- `internal/tui/forms.go` — huh wrappers: `Confirm`, `Input`, `Select`.
- `internal/tui/persons/` — bubbletea model for the interactive persons browser (input, list, detail pane).

**Integration rules:**

- If stdout is not a TTY, or `--json` / `--jq` / `--md` is set, the TUI is suppressed and output falls back to the plain JSON envelope.
- `LINDEN_NO_TUI=1` also suppresses the TUI unconditionally.
- The bubbletea program is launched with `tea.NewProgram(model, tea.WithAltScreen()).Run()`; after it exits, any selected record is printed to stdout as a JSON envelope.

**Dependency additions (go.mod):**

- `charm.land/bubbletea/v2`
- `github.com/charmbracelet/bubbles`
- `github.com/charmbracelet/huh`
- `github.com/charmbracelet/lipgloss`
- `github.com/charmbracelet/glamour`

### Configuration precedence

1. CLI flags
2. `LINDEN_TOKEN`, `LINDEN_BASE_URL`, `LINDEN_ACCOUNT`, `LINDEN_NO_KEYRING`, `LINDEN_NO_TUI` env vars
3. `.linden/config.json` (local, repo-level)
4. `~/.config/linden/config.json` (global)
5. Hardcoded defaults

### Agent integration

- `--agent` flag on `--help` emits JSON command catalog.
- `skills/linden/SKILL.md` contains trigger keywords, agent invariants, and usage examples.

## Testing Decisions

### What makes a good test

Tests should exercise the observable behavior of a module through its public interface, not its internal implementation. A test should break only when behavior changes, not when internal structure is refactored. Prefer table-driven tests with named cases.

### Modules to test

**`internal/auth`** — Unit test token storage, refresh logic, and PKCE verifier/challenge generation. Mock the HTTP token endpoint; do not test the browser-open flow.

**`internal/config`** — Unit test the precedence chain: env var overrides local config overrides global config overrides default. Test that `LINDEN_BASE_URL`, `LINDEN_TOKEN`, and `LINDEN_ACCOUNT` are correctly picked up.

**`internal/client`** — Unit test each API method by injecting a `*httptest.Server` that returns canned JSON responses. Verify that request headers (Authorization, Content-Type), query parameters (page, size), and JSON body (for create/update) are constructed correctly. Verify error mapping for 401, 403, 404, 500 responses.

**`internal/output`** — Unit test the JSON envelope structure for success and error cases. Test `--jq` expression evaluation with gojq. Test format auto-detection.

**`internal/commands`** — Integration-style unit tests that wire a real Cobra command against a mock client (via interface). Verify that flags are correctly passed through to the client and that output is rendered in the expected format. Model after `basecamp-cli/internal/commands/todos_test.go`.

**`internal/tui/persons`** — Unit test the bubbletea model by sending `tea.Msg` values and asserting on the resulting model state (selected item, filter state, error display). Do not test rendering output. For huh forms in `internal/tui/forms.go`, test the wrapper logic only.

**e2e (BATS)** — Shell script tests that run the compiled binary against a mock HTTP server or a live `LINDEN_BASE_URL` test environment. Cover happy paths for `accounts list`, `persons list`, `persons show`, `persons create`, `persons update`, `persons delete`, `auth status`, and `doctor`. TUI mode is suppressed in e2e tests via `LINDEN_NO_TUI=1`.

### Prior art

`basecamp-cli/internal/commands/todos_test.go` — demonstrates the pattern of setting up a test app context with a disabled HTTP transport, creating a test server, and asserting on command output.

## Out of Scope

- Account create, update, or delete — managed via the web app; CLI is read-only for accounts
- Persons search endpoint (`GET /accounts/{id}/persons/search`) — agents can use `--jq` on list output to filter client-side
- Memberships, invitations, online-accounts endpoints
- Document management (birth certificates, passports, driver licenses, SSNs)
- Reminders management
- Shell completion (can be added later)
- Windows support (darwin/linux only for now)
- Multi-profile support (single identity per config; no `--profile` flag)

## Further Notes

- The Linden API is live at `https://api.mylinden.family`. `LINDEN_BASE_URL` overrides this for local/staging development.
- The Auth0 application must be registered as a **public/native client** with PKCE enabled and `http://localhost:3009/callback` in its allowed redirect URIs. Confirm this with `auth.mylinden.family` before implementing the auth flow.
- The basecamp-cli at `/Users/emiliano.jankowski/sites/equinix/tmp/basecamp-cli` is the structural reference implementation. When in doubt about package layout, output format, or command conventions, defer to how basecamp-cli handles it.
- All person IDs and account IDs are UUIDs (not integers).
- The `relationship_type` field is a free-form string; the API provides a lookup at `GET /persons/relationship-types` that can be used to populate a selection list in the huh create form.
