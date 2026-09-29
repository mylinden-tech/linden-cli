# Contacts

The hub decision tree already applies. This file covers only what is specific to contacts.

## Commands

| Task | Command |
|---|---|
| List | `linden contacts list --json` |
| Show | `linden contacts show <uuid> --json` |
| Search | `linden contacts search --q "Ada" --json` |
| Types | `linden contacts types --json` |

## Domain rules

- `contacts types` does not require an active account. List, show, and search do.
- Search requires a non-empty `--q`.
- When the classification is unclear, run `linden contacts types --json` before describing the contact.
- Contacts are read-only. Do not invent create, update, or delete commands.

## On errors

- Not found → list or search again. Never invent UUIDs.
