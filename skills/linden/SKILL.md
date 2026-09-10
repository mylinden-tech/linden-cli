---
name: linden
description: |
  This skill should be used when the user mentions Linden, mylinden, the linden CLI,
  a Linden family account, or asks to inspect or change family data through Linden.
  Pair with linden-persons for people or family members.
  Pair with linden-contacts for contacts.
  Pair with linden-memberships for account memberships and roles.
  Pair with linden-invitations for pending account invitations.
  Pair with linden-account-shares for resource shares.
  Pair with linden-account-settings for account details, statistics, or user settings.
  Pair with linden-pets for pets.
  Pair with linden-reminders for reminders.
  Pair with linden-todos for todos and todo lists.
  Pair with linden-vehicles, linden-real-estates, linden-online-accounts,
  linden-insurances, linden-wills, or linden-share-links for those read domains.
  Pair with linden-doctor for login, setup, doctor, or a failing CLI.
---

# Linden

Drive the `linden` CLI. Do not call `https://api.mylinden.family` over HTTP. Do not invent command groups.

## Invariants

1. If this session has not already proven a healthy CLI, run `linden doctor --agent` before listing or mutating domain data. Load `linden-doctor` to interpret the result.
2. Prefer `--agent` for data-only JSON. Use `--json` alone when `summary` or `breadcrumbs` are needed (`--jq` drops the envelope even with `--json`). Use `--md` when presenting to a human. Filter `data` with `--jq '<expr>'` or `--json --jq '<expr>'`. Never pipe to external `jq`. Never combine `--agent --jq` (jq is ignored).
3. Follow `breadcrumbs` from the last `--json` response. Do not guess the next command.
4. Domain commands require an active account: `linden accounts use <id>` or `--account` / `LINDEN_ACCOUNT`. IDs are UUIDs.
5. If the routing table has no match, say the CLI has no such resource. Do not invent resource command groups.
6. Do not read credential files, the OS keyring, or print `LINDEN_TOKEN`.
7. Destructive commands require a UUID observed in this session and `--yes`.
8. After a passing doctor, skip doctor for later domain commands in this thread. Re-run doctor after a setup-shaped error (`auth`, `account`, unreachable API).

See [references/envelope.md](references/envelope.md) for output modes and [references/auth-and-accounts.md](references/auth-and-accounts.md) for login and accounts.

## Routing

- Setup, login, CLI broken → load `linden-doctor`
- Auth status / login / logout → `linden auth status --agent`, `linden auth login`, or `linden auth logout`
- Account details, statistics, or user settings → load `linden-account-settings`
- List or switch accounts → `linden accounts list --json` then `linden accounts use <id>`
- People and family members → load `linden-persons`
- Contacts → load `linden-contacts`
- Account memberships and roles → load `linden-memberships`
- Pending account invitations → load `linden-invitations`
- Account resource shares → load `linden-account-shares`
- Pets → load `linden-pets`
- Reminders → load `linden-reminders`
- Todos and todo lists → load `linden-todos`
- Vehicles → load `linden-vehicles`
- Real estate → load `linden-real-estates`
- Online accounts → load `linden-online-accounts`
- Insurance policies and providers → load `linden-insurances`
- Wills → load `linden-wills`
- Resource-specific share links → load `linden-share-links`
- Documents, passports, SSNs, files, driver licenses, birth certificates, body memorial wishes, legacy messages, and contact-import → not in the CLI; say so

## Cheat sheet

```bash
linden doctor --agent
linden auth login
linden auth status --agent
linden auth logout
linden accounts list --json
linden accounts use <uuid>
linden accounts use <uuid> --scope local
```

Set `LINDEN_NO_TUI=1` if a TTY might open a form or persons browser.

## Safety

Prefer `--agent` and summarize PII (names, emails, phones, addresses, birthdays) for the human. Do not paste full contact dumps unless asked. CLI only.
