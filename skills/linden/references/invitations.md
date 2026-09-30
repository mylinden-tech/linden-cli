# Invitations

The hub decision tree already applies. This file covers only what is specific to invitations.

## Commands

| Task | Command |
|---|---|
| List | `linden invitations list --json` |
| Show | `linden invitations show <uuid> --json` |

## Domain rules

- Invitations are read-only. Do not invent create, accept, decline, resend, update, or delete commands.
- Summarize email addresses and personal messages unless the user asks for the full text.

## On errors

- Not found → re-list. Never invent UUIDs.
