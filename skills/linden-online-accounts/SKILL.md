---
name: linden-online-accounts
description: |
  This skill should be used when listing or viewing online-account metadata in Linden.
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Online Accounts

Read online-account metadata through the CLI.

## Commands

- List: `linden online-accounts list --agent` (optional `--page`, `--size`)
- Show: `linden online-accounts show <uuid> --agent`

Resolve IDs from list results. These commands never reveal passwords and must not be used to request, infer, or expose credentials. This skill is read-only.
