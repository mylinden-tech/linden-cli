---
name: linden-real-estates
description: |
  This skill should be used when listing or viewing real estate in a Linden account.
  Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Real Estate

Read real-estate records through the CLI.

## Commands

- List: `linden real-estates list --agent` (optional `--page`, `--size`)
- Show: `linden real-estates show <uuid> --agent`

Resolve IDs by listing first. Summarize properties by name, city, and state. Do not expose the street address unless the user asks for it. This skill is read-only.
