---
name: linden-insurances
description: |
  This skill should be used when listing or viewing Linden insurance policies, providers, or expirations.
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Insurances

Read insurance records through the CLI.

## Commands

- List policies: `linden insurances list --agent` (optional `--page`, `--size`)
- Show policy: `linden insurances show <uuid> --agent`
- List providers: `linden insurances providers --agent`
- List expiring policies: `linden insurances expiring --agent`

Use IDs returned by list. Agent output omits policy numbers. Summarize status, type, provider, and expiration date. This skill is read-only.
