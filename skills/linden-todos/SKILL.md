---
name: linden-todos
description: |
  This skill should be used when managing Linden todos or todo lists (list,
  show, create, update, or delete). Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Todos

Manage todos and todo lists in the active Linden account through the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List account todos | `linden todos list --json` |
| List todo lists | `linden todos lists --json` |
| Create a todo list | `linden todos create-list --title … --json` |
| Show a todo | `linden todos show <uuid> --json` |
| Create a todo | `linden todos create --list <uuid> --title … --json` |
| Update a todo | `linden todos update <uuid> --<flag> … --json` |
| Delete a todo | `linden todos delete <uuid> --yes --agent` |

Use `--json` when chaining through breadcrumbs. Use `--agent` for data-only output.

## Fields

- Todo creation requires `--list` and `--title`. Optional flags are `--description`, `--due-date`, and `--assigned-to`.
- Todo updates accept `--title`, `--description`, `--due-date`, `--assigned-to`, and `--status`. Valid statuses are `open` and `completed`; only explicitly supplied flags are sent.
- Todo-list creation requires `--title`. To link the list to a resource, supply `--resource-type` and `--resource-id` together.
- Dates use `YYYY-MM-DD`.

## Rules

- Run `linden todos lists --json` before creating a todo and use an observed list UUID. Never guess UUIDs.
- If no suitable list exists, create one with `linden todos create-list --title … --json`, then follow its create breadcrumb.
- Resolve todo IDs with list or show before updating or deleting.
- Star and unstar operations are not supported by this CLI.
- Delete requires `--yes` and a todo UUID observed in this session.
- On not found, re-list todos. On validation errors, correct only the failing fields.

## Safety

Summarize todo descriptions, assignments, and linked list details for the human. Do not paste full todo lists unless asked.
