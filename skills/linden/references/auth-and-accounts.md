# Auth and accounts

## Authentication

**`linden auth login`** opens a browser for Auth0 PKCE authentication (local loopback callback on port 3009). Start the command and wait for the human to complete the browser flow. Do not scrape the callback URL or read stored tokens.

**`linden auth status --agent`** returns data such as:

```json
{ "authenticated": true, "expires_at": "...", "token_type": "Bearer" }
```

When not authenticated:

```json
{ "authenticated": false }
```

**`linden auth logout`** removes stored credentials.

## Accounts

**`linden accounts list --json`** lists accounts the authenticated user can access. Use `--json` to read breadcrumbs suggesting `linden accounts use <id>`.

**`linden accounts use <id>`** sets the active account. Default scope is global (`~/.config/linden/config.json`). Pass **`--scope local`** to write `.linden/config.json` in the current directory instead.

One-shot override without persisting: **`--account <id>`** on any command, or set **`LINDEN_ACCOUNT`** in the environment.

All persons commands require an active account. Without one, the CLI returns a usage error with a hint to list and select an account.

## Configuration

Precedence (lowest to highest): built-in defaults → global config → local config → environment variables → CLI flags.

| Location | Path |
|---|---|
| Global config | `~/.config/linden/config.json` (or `$XDG_CONFIG_HOME/linden/config.json`) |
| Local config | `.linden/config.json` in the current working directory |

| Variable | Purpose |
|---|---|
| `LINDEN_BASE_URL` | API base URL (default `https://api.mylinden.family`) |
| `LINDEN_ACCOUNT` | Override active account ID |
| `LINDEN_NO_TUI` | Disable interactive TUI prompts and browser |
| `LINDEN_NO_KEYRING` | Disable OS keyring for credential storage |
| `LINDEN_TOKEN` | Static token bypass for CI only — do not use for interactive agents |

Do not read `credentials.json` or keyring contents directly.
