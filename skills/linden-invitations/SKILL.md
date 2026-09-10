---
name: linden-invitations
description: |
  This skill should be used when listing or viewing pending invitations for a
  Linden account.
disable-model-invocation: true
---

# Linden Invitations

Inspect pending invitations in the active Linden account via the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List | `linden invitations list --json` (optional `--page`, `--size`) |
| Show | `linden invitations show <uuid> --json` |

Use **`--json`** when summaries or breadcrumbs are needed. Use **`--agent`** for data-only JSON.

## Rules

- Resolve an invitation by listing first, then use the returned UUID. Do not guess IDs.
- On **not found**, list again; do not invent UUIDs.
- Invitations are read-only in this CLI. Do not invent create, accept, decline, resend, update, or delete commands.

## Safety

Invitations contain email addresses and may contain a personal message. Prefer `--agent` and summarize personal information unless full details are requested.
