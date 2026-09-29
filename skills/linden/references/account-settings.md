# Account settings

The hub decision tree already applies. This file covers account details, statistics, and user settings.

## Commands

| Task | Command |
|---|---|
| Show active account | `linden accounts show --json` |
| Show an account by ID | `linden accounts show <uuid> --json` |
| Statistics | `linden accounts stats --json` |
| User settings | `linden user-settings show --json` |

## Domain rules

- `accounts show` without an ID and `accounts stats` require an active account.
- User settings belong to the authenticated user and do not require an active account.
- These commands are read-only. Do not invent create, update, or delete commands.

## On errors

- Not found → `linden accounts list --json`. Never invent UUIDs.
