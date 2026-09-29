# Auth and accounts

The hub decision tree already applies. This file covers login and choosing an account.

## Commands

| Task | Command |
|---|---|
| Log in | `linden auth login` |
| Auth status | `linden auth status --agent` |
| Log out | `linden auth logout` |
| List accounts | `linden accounts list --json` |
| Use account | `linden accounts use <uuid>` |
| Use account in this directory | `linden accounts use <uuid> --scope local` |

## Domain rules

- `linden auth login` opens a browser. Wait for the human. Do not scrape the callback URL or read stored tokens.
- Default `accounts use` writes `~/.config/linden/config.json`. `--scope local` writes `.linden/config.json` in the current directory.
- One-shot override without persisting: `--account <uuid>` on any command, or `LINDEN_ACCOUNT`.
- When showing account details or statistics, read references/account-settings.md before `linden accounts show` or `linden accounts stats`.
- `LINDEN_NO_KEYRING` disables the OS keyring. Do not read `credentials.json`.

## On errors

- When auth fails, read references/doctor.md before running another domain command.
