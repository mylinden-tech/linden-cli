---
name: linden-account-shares
description: |
  This skill should be used when listing or viewing resource shares for a
  Linden account.
disable-model-invocation: true
---

# Linden Account Shares

Inspect resource shares in the active Linden account via the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List | `linden account-shares list --json` (optional `--page`, `--size`) |
| Show | `linden account-shares show <uuid> --json` |

Use **`--json`** when summaries or breadcrumbs are needed. Use **`--agent`** for data-only JSON.

## Rules

- Resolve an account share by listing first, then use the returned UUID. Do not guess IDs.
- `account-shares show` searches the active account's share list because the API has no dedicated read endpoint for one share.
- On **not found**, list again and confirm the active account; do not invent UUIDs.
- Account shares are read-only in this CLI. Do not invent revoke or delete commands.

## Safety

Share metadata can reveal sensitive estate resources and recipients. Summarize it for the human unless full details are requested.
