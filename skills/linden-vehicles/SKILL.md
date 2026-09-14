---
name: linden-vehicles
description: |
  This skill should be used when listing or viewing vehicles in a Linden account.
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Vehicles

Read vehicle records in the active account through the CLI.

## Commands

- List: `linden vehicles list --agent` (optional `--page`, `--size`)
- Show: `linden vehicles show <uuid> --agent`

Use IDs returned by list; never guess them. Agent output omits VIN and license plate. Use `--json` only when the user explicitly needs those fields. This skill is read-only.
