# Persons

The hub decision tree already applies. This file covers only what is specific to persons.

## Commands

| Task | Command |
|---|---|
| List | `linden persons list --json` |
| Show | `linden persons show <uuid> --json` |
| Create | `linden persons create --first-name … --last-name … --json` |
| Update | `linden persons update <uuid> --json` |
| Delete | `linden persons delete <uuid> --yes --json` |

## Fields

Before the first create or update in a session, run `linden persons create --agent --help` or `linden persons update --agent --help`.
It is the source of truth for flags and allowed values.

## Domain rules

- Required on create: `--first-name` and `--last-name`.
- When showing a copy-paste create or name-resolution recipe, read references/persons-examples.md before building the command.

## On errors

- Validation error → re-run `--agent --help` and resend only the failing fields.
- Not found → re-list; never invent UUIDs.
