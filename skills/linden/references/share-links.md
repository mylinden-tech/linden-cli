# Share links

The hub decision tree already applies. This file covers resource-specific share links, which are distinct from account shares.

## Commands

| Task | Command |
|---|---|
| List for the active account | `linden share-links list --agent` |
| List for another resource | `linden share-links list --resource-type <type> --resource-id <uuid> --agent` |
| Show | `linden share-links show <uuid> --resource-type <type> --resource-id <uuid> --agent` |

## Domain rules

- The resource type defaults to `accounts`. The resource id defaults to the active account.
- Show filters that resource's list. Creating and revoking links are unsupported.
- When the user means who the account is shared with, read references/account-shares.md before `linden account-shares list`.

## On errors

- Not found → re-list that resource. Never invent UUIDs.
