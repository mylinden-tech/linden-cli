---
name: linden-pets
description: |
  This skill should be used when managing Linden pets (list, show, create, update, delete).
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Pets

Manage pets in the active Linden account via the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List | `linden pets list --json` (optional `--page`, `--size`) |
| Show | `linden pets show <uuid> --json` |
| Create | `linden pets create --name … --species … [flags] --json` |
| Update | `linden pets update <uuid> --<flag> … --json` |
| Delete | `linden pets delete <uuid> --yes --agent` |

Use **`--json`** when chaining via breadcrumbs. Use **`--agent`** for data-only output on delete.

## Rules

- Resolve names by listing first (`linden pets list --json`), then use the UUID from the response. Do not guess IDs.
- **`delete`** requires `--yes` and a UUID observed via list or show in this session. If multiple pets match a name, confirm the UUID before deleting.
- After **create**, follow the breadcrumb to `linden pets show <id>`.
- Required create flags: **`--name`** and **`--species`**. All other fields are optional — see [references/pet-fields.md](references/pet-fields.md).
- On **not found**, re-list; do not invent UUIDs.
- On **validation errors**, re-read [references/pet-fields.md](references/pet-fields.md) and resend only the failing fields.

## Safety

Summarize pet records and potentially sensitive identifiers such as microchip numbers for the human. Do not paste full pet lists unless asked.
