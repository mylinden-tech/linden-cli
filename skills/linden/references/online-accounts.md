# Online accounts

The hub decision tree already applies. This file covers only what is specific to online accounts.

## Commands

| Task | Command |
|---|---|
| List | `linden online-accounts list --agent` |
| Show | `linden online-accounts show <uuid> --agent` |

## Domain rules

- Read-only. Do not invent create, update, or delete commands.
- These commands never reveal passwords. Do not request, infer, or expose credentials.

## On errors

- Not found → re-list. Never invent UUIDs.
