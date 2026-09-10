---
name: linden-contacts
description: |
  This skill should be used when listing, viewing, searching, or inspecting
  contact types for Linden contacts.
disable-model-invocation: true
---

# Linden Contacts

Inspect contacts in the active Linden account via the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List | `linden contacts list --json` (optional `--page`, `--size`) |
| Show | `linden contacts show <uuid> --json` |
| Search | `linden contacts search --q "Ada" --json` |
| Types | `linden contacts types --json` |

Use **`--json`** when summaries or breadcrumbs are needed. Use **`--agent`** for data-only JSON.

## Rules

- Resolve a contact by listing or searching first, then use the returned UUID. Do not guess IDs.
- Search requires a non-empty **`--q`** value.
- Use `contacts types` to inspect valid contact classifications.
- On **not found**, list or search again; do not invent UUIDs.
- Contacts are read-only in this CLI. Do not invent create, update, or delete commands.

## Safety

Prefer `--agent` and summarize PII such as names, emails, phone numbers, and addresses for the human. Do not paste full contact lists unless asked.
