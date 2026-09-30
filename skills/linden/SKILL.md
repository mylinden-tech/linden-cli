---
name: linden
description: |
  Use when the user mentions Linden, mylinden, the linden CLI, or a Linden family
  account, or wants to view or change family data stored in Linden: people, contacts,
  pets, vehicles, homes, insurance, wills, reminders, todos, account members, invitations,
  or sharing.
argument-hint: "[domain]"
---

# Linden

Drive the `linden` CLI. One slash command: `/linden`. There is no `/linden-reminders` or other domain command.

## Invocation

- `/linden` with no argument → decision tree. Pick the reference from the routing table.
- `/linden <domain>` → read that reference before any domain CLI command, then start the decision tree at step 1. At step 2, treat the domain as already chosen.
  The argument is the reference filename without `.md`: `/linden reminders` → `references/reminders.md`.
  Aliases: `auth` and `accounts` → `references/auth-and-accounts.md`; `settings` → `references/account-settings.md`; `shares` → `references/account-shares.md`.
  Unknown argument → say so. Do not invent a reference or a CLI group.

## Session rules

- Do not call the Linden API over HTTP. Do not invent command groups.
- Set `LINDEN_NO_TUI=1` for the session. Persons list and create open a TUI when stdout is a TTY and this variable is unset.
- Do not read credential files or the OS keyring, and do not print `LINDEN_TOKEN`.

## Decision tree (every request)

1. CLI healthy in this session?
   No / unknown → `linden doctor --agent`. Failing → read references/doctor.md and follow its remediation. Do not run domain commands until doctor passes.
2. Domain already chosen by `/linden <domain>`?
   Yes → the Invocation section already named the reference. Read it if not read yet.
   No → routing table → read that reference BEFORE running any domain CLI command.
   No match → say the CLI does not support it.
3. Does the selected command require an active account?
   No, when the reference says this command does not → step 4.
   Yes, and none is set → `linden accounts list --json`; one account → use it; several → ask the user which.
4. Read, create, or change an existing record?
   Read → `list` / `show` with `--json` → summarize; never paste full PII unless asked.
   Create → step 7. There is no target UUID. If a same-name record already exists, show it and ask before creating another.
   Update, or any command that needs an existing id → step 5, then step 7.
   Delete, revoke, unshare, or remove a member → step 6. Stop. Do not continue to step 7.
5. Target UUID observed in this session?
   No → `linden <domain> list --json` and match by name.
        0 matches → ask the user. Do not create a record.
        More than 1 match → show the candidates and ask which one.
        Exactly 1 match → use it.
6. Destructive (delete, revoke, unshare, remove member)?
   Yes → do not run the command, and do not pass `--yes`. Tell the user, politely, that deleting records, revoking access, unsharing, and removing members is not supported here.
7. Execute with `--json` so the envelope includes `breadcrumbs`. Follow `breadcrumbs`, then verify with `show` when the breadcrumb says to. Report the result.
   Use `--agent` for `linden doctor` and for payload-only reads that do not need breadcrumbs.

## Routing

| User talks about… | Read |
|---|---|
| login, setup, CLI not working, linden doctor | references/doctor.md |
| switch account, which account, log out | references/auth-and-accounts.md |
| person, family member, relative, sister, son, birthday, relationship | references/persons.md |
| contact, phone book, family doctor, lawyer | references/contacts.md |
| who has access, member, role, admin | references/memberships.md |
| invite, pending invitation | references/invitations.md |
| share with someone, shared resource | references/account-shares.md |
| share link, public link | references/share-links.md |
| account name, statistics, my settings | references/account-settings.md |
| pet, dog, cat, vet | references/pets.md |
| reminder, remind me | references/reminders.md |
| todo, task, checklist, todo list | references/todos.md |
| car, vehicle, plate, VIN | references/vehicles.md |
| house, home, property, real estate | references/real-estates.md |
| online account, website login, subscription | references/online-accounts.md |
| insurance, policy, provider, coverage | references/insurances.md |
| will, testament, executor | references/wills.md |
| documents, passport, SSN, driver license, birth certificate, legacy messages, memorial wishes, contact import | Not in the CLI. Say so. |

A family doctor or lawyer is a contact. The doctor row is the CLI doctor.

Read a second reference only when the command needs an id from that domain (for example, an insurance policy linked to a vehicle). A reminder about car insurance reads references/reminders.md only.

## Output modes

When choosing between `--agent`, `--json`, `--md`, or `--jq`, read references/envelope.md before picking a flag.

## Safety

Summarize PII (names, emails, phones, addresses, birthdays). Do not paste full lists unless asked.
