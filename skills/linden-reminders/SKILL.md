---
name: linden-reminders
description: |
  This skill should be used when managing Linden reminders (list, show, create,
  update, complete, or delete). Prefer loading via the linden hub skill routing table.
disable-model-invocation: true
---

# Linden Reminders

Manage reminders in the active Linden account through the CLI. Requires authentication and an active account. See the `linden` hub skill for doctor and account setup.

## Commands

| Task | Command |
|---|---|
| List account reminders | `linden reminders list --json` |
| Show | `linden reminders show <uuid> --json` |
| Create for a person | `linden reminders create --person <uuid> --name … --due-date YYYY-MM-DD --json` |
| Create for a pet | `linden reminders create --pet <uuid> --name … --due-date YYYY-MM-DD --json` |
| Update | `linden reminders update <uuid> --<flag> … --json` |
| Complete | `linden reminders complete <uuid> --json` |
| Delete | `linden reminders delete <uuid> --yes --agent` |

Use `--json` when chaining through breadcrumbs. Use `--agent` for data-only output.

## Fields

- Create requires `--name`, `--due-date`, and exactly one of `--person` or `--pet`.
- Optional create and update flags are `--notes`, `--completed`, `--repeat-interval`, `--remind-before-days`, and `--reminder-type`.
- Update also accepts `--name` and `--due-date`. Only explicitly supplied flags are sent.
- Dates use `YYYY-MM-DD`. `--remind-before-days` must be zero or greater.
- Common repeat intervals are `daily`, `weekly`, `monthly`, `yearly`, and `never`.
- Common reminder types include `birthday`, `document`, `custom`, `warranty`, and `insurance_expiration`.

## Rules

- Resolve entity and reminder IDs by listing or showing the relevant resource first. Never guess UUIDs.
- Creation supports person and pet targets only. Do not invent account-scoped reminder creation.
- Use `complete` to set `completed=true`; use update with `--completed=false` to reopen.
- Delete requires `--yes` and a reminder UUID observed in this session.
- On not found, re-list reminders. On validation errors, correct only the failing fields.

## Safety

Summarize reminder notes and linked entity details for the human. Do not paste full reminder lists unless asked.
