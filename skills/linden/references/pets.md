# Pets

The hub decision tree already applies. This file covers only what is specific to pets.

## Commands

| Task | Command |
|---|---|
| List | `linden pets list --json` |
| Show | `linden pets show <uuid> --json` |
| Create | `linden pets create --name … --species … --json` |
| Update | `linden pets update <uuid> --json` |

## Fields

Before the first create or update in a session, run `linden pets create --agent --help` or `linden pets update --agent --help`.
It is the source of truth for flags and allowed values.

## Domain rules

- Required on create: `--name` and `--species`.
- Deleting a pet is not supported. If asked, say so politely and do not run a delete command.
- Summarize microchip numbers unless the user asks for them.

## On errors

- Validation error → re-run `--agent --help` and resend only the failing fields.
- Not found → re-list; never invent UUIDs.
