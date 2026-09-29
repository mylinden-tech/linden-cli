# Reminders

The hub decision tree already applies. This file covers only what is specific to reminders.

## Commands

| Task | Command |
|---|---|
| List | `linden reminders list --json` |
| Show | `linden reminders show <uuid> --json` |
| Create for a person | `linden reminders create --person <uuid> --name … --due-date YYYY-MM-DD --json` |
| Create for a pet | `linden reminders create --pet <uuid> --name … --due-date YYYY-MM-DD --json` |
| Update | `linden reminders update <uuid> --json` |
| Complete | `linden reminders complete <uuid> --json` |
| Delete | `linden reminders delete <uuid> --yes --json` |

## Fields

Before the first create or update in a session, run `linden reminders create --agent --help` or `linden reminders update --agent --help`.
It is the source of truth for flags and allowed values.

## Domain rules

- Creation targets a person or a pet. Do not invent an account-scoped create.
- Use `complete` to mark done. Reopen with `linden reminders update <uuid> --completed=false --json`.
- Resolve the person or pet id from that domain before create. When the target is a person, read references/persons.md before `linden reminders create --person`. When the target is a pet, read references/pets.md before `linden reminders create --pet`.

## On errors

- Validation error → re-run `--agent --help` and resend only the failing fields.
- Not found → re-list; never invent UUIDs.
