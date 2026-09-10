---
name: linden-wills
description: |
  This skill should be used when listing or viewing will metadata in a Linden account.
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Wills

Read will metadata through the CLI.

## Commands

- List: `linden wills list --agent` (optional `--page`, `--size`)
- Show: `linden wills show <uuid> --agent`

Focus summaries on name, issued date, and whether the will came from an attorney. Agent output omits document URLs and IDs, notes, and plain-language summaries. This skill is read-only.
