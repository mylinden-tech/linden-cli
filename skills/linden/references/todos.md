# Todos

The hub decision tree already applies. This file covers only what is specific to todos.

## Commands

| Task | Command |
|---|---|
| List todos | `linden todos list --json` |
| List todo lists | `linden todos lists --json` |
| Create a list | `linden todos create-list --title … --json` |
| Show a todo | `linden todos show <uuid> --json` |
| Create a todo | `linden todos create --list <uuid> --title … --json` |
| Update a todo | `linden todos update <uuid> --json` |

## Fields

Before the first create or update in a session, run `linden todos create --agent --help` or `linden todos update --agent --help`.
It is the source of truth for flags and allowed values.

## Domain rules

- Run `linden todos lists --json` before creating a todo and use an observed list id.
- If no suitable list exists, create one with `linden todos create-list`, then follow its breadcrumb.
- Deleting a todo is not supported. If asked, say so politely and do not run a delete command.
- Star and unstar are not supported.

## On errors

- Validation error → re-run `--agent --help` and resend only the failing fields.
- Not found → re-list; never invent UUIDs.
