# Real estate

The hub decision tree already applies. This file covers only what is specific to real estate.

## Commands

| Task | Command |
|---|---|
| List | `linden real-estates list --agent` |
| Show | `linden real-estates show <uuid> --agent` |

## Domain rules

- Read-only. Do not invent create, update, or delete commands.
- Summarize by name, city, and state. Do not expose the street address unless the user asks for it.

## On errors

- Not found → re-list. Never invent UUIDs.
