# Account shares

The hub decision tree already applies. This file covers only what is specific to account shares.

## Commands

| Task | Command |
|---|---|
| List | `linden account-shares list --json` |
| Show | `linden account-shares show <uuid> --json` |

## Domain rules

- `account-shares show` searches the active account's share list. The API has no single-share read.
- Account shares are read-only. Do not invent revoke or delete commands.
- Summarize recipients and linked resources unless the user asks for the full record.

## On errors

- Not found → re-list and confirm the active account. Never invent UUIDs.
