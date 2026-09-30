# Insurances

The hub decision tree already applies. This file covers only what is specific to insurance.

## Commands

| Task | Command |
|---|---|
| List policies | `linden insurances list --agent` |
| Show a policy | `linden insurances show <uuid> --agent` |
| List providers | `linden insurances providers --agent` |
| List expiring policies | `linden insurances expiring --agent` |

## Domain rules

- Read-only. Do not invent create, update, or delete commands.
- `--agent` omits policy numbers. Summarize status, type, provider, and expiration. Use `--json` only when the user explicitly needs the policy number.
- When the command needs a vehicle id, read references/vehicles.md before using that id.

## On errors

- Not found → re-list. Never invent UUIDs.
