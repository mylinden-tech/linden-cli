---
name: linden-share-links
description: |
  This skill should be used when listing or viewing share links for a Linden resource.
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Share Links

Read resource-specific share links through the CLI. Share links are distinct from account-shares.

## Commands

- Active account: `linden share-links list --agent`
- Other resource: `linden share-links list --resource-type <type> --resource-id <uuid> --agent`
- Show: `linden share-links show <uuid> --resource-type <type> --resource-id <uuid> --agent`

The resource type defaults to `accounts`; the resource ID defaults to the active account. Show filters the resource's list. Creating and revoking links are unsupported.
