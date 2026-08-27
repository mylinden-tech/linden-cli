---
name: linden-persons
description: |
  This skill should be used when the user asks to list, show, create, update, or delete
  Linden people, persons, contacts, or family members, or mentions person fields such as
  name, email, phone, birthday, address, or relationship type.
---

# Linden Persons

Manage persons in the active Linden account via the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List | `linden persons list --json` (optional `--page`, `--size`) |
| Show | `linden persons show <uuid> --json` |
| Create | `linden persons create --first-name … --last-name … [flags] --json` |
| Update | `linden persons update <uuid> --<flag> … --json` |
| Delete | `linden persons delete <uuid> --yes --agent` |

Use **`--json`** when chaining via breadcrumbs. Use **`--agent`** for data-only output on delete.

Set **`LINDEN_NO_TUI=1`** if a TTY might launch the interactive persons browser or create form.

## Rules

- Resolve names by listing first (`linden persons list --json`), then use the UUID from the response. Do not guess IDs.
- **`delete`** requires `--yes` and a UUID observed via list or show in this session. If the user says "delete Jane" and multiple people match, confirm the UUID before deleting.
- After **create**, follow the breadcrumb to `linden persons show <id>`.
- Required create flags: **`--first-name`**, **`--last-name`**. All other fields are optional — see [references/person-fields.md](references/person-fields.md).
- On **not found**, re-list; do not invent UUIDs.
- On **validation errors**, re-read [references/person-fields.md](references/person-fields.md) and resend only the failing fields.

## Examples

See [references/examples.md](references/examples.md) for copy-paste recipes and name-resolution workflows.

## Safety

Summarize PII for the human. Do not paste full contact lists unless asked.
