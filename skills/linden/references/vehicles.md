# Vehicles

The hub decision tree already applies. This file covers only what is specific to vehicles.

## Commands

| Task | Command |
|---|---|
| List | `linden vehicles list --agent` |
| Show | `linden vehicles show <uuid> --agent` |

## Domain rules

- Read-only. Do not invent create, update, or delete commands.
- `--agent` omits VIN and license plate. Use `--json` only when the user explicitly needs those fields.

## On errors

- Not found → re-list. Never invent UUIDs.
