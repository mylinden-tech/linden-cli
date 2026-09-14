---
name: linden-memberships
description: |
  This skill should be used when listing or viewing Linden account memberships,
  or when inspecting the available membership roles.
disable-model-invocation: true
---

# Linden Memberships

Inspect memberships in the active Linden account via the CLI. Requires authentication and an active account for listing. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List | `linden memberships list --json` (optional `--page`, `--size`) |
| Show | `linden memberships show <uuid> --json` |
| Roles | `linden memberships roles --json` |

Use **`--json`** when summaries or breadcrumbs are needed. Use **`--agent`** for data-only JSON.

## Rules

- Resolve a membership by listing first, then use the returned UUID. Do not guess IDs.
- Use `memberships roles` to inspect role names and descriptions.
- On **not found**, list again; do not invent UUIDs.
- Memberships are read-only in this CLI. Do not invent update or delete commands.

## Safety

Membership records can include user details. Prefer `--agent` and summarize personal information for the human unless full details are requested.
