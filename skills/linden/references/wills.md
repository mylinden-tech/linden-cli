# Wills

The hub decision tree already applies. This file covers only what is specific to wills.

## Commands

| Task | Command |
|---|---|
| List | `linden wills list --agent` |
| Show | `linden wills show <uuid> --agent` |

## Domain rules

- Read-only. Do not invent create, update, or delete commands.
- Summarize name, issued date, and whether the will came from an attorney. `--agent` omits document URLs, document ids, notes, and plain-language summaries.

## On errors

- Not found → re-list. Never invent UUIDs.
