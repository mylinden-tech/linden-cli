---
name: linden-account-settings
description: |
  This skill should be used when viewing Linden account details or statistics,
  or when inspecting the authenticated user's settings and default account.
disable-model-invocation: true
---

# Linden Account Settings

Inspect account details, aggregate account statistics, and current user settings via the CLI. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| Show active account | `linden accounts show --json` |
| Show an account by ID | `linden accounts show <uuid> --json` |
| Show active account statistics | `linden accounts stats --json` |
| Show current user settings | `linden user-settings show --json` |

Use **`--json`** when summaries or breadcrumbs are needed. Use **`--agent`** for data-only JSON.

## Rules

- `accounts show` without an ID and `accounts stats` require an active account.
- Use `accounts list` before showing a different account by UUID. Do not guess IDs.
- User settings are for the authenticated user and do not require an active account.
- These commands are read-only. Do not invent account or user-settings create, update, or delete commands.

## Safety

Account details and settings can reveal family and user metadata. Summarize them for the human unless full details are requested.
