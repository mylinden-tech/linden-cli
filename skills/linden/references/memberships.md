# Memberships

The hub decision tree already applies. This file covers only what is specific to memberships.

## Commands

| Task | Command |
|---|---|
| List | `linden memberships list --json` |
| Show | `linden memberships show <uuid> --json` |
| Roles | `linden memberships roles --json` |

## Domain rules

- `memberships roles` does not require an active account. List and show do.
- When the role name is unclear, run `linden memberships roles --json` before describing it.
- Memberships are read-only. Do not invent update or delete commands.

## On errors

- Not found → re-list. Never invent UUIDs.
