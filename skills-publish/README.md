# Agent Skills for Linden

[Agent skills](https://agentskills.io) for [Linden](https://www.mylinden.family/) — inspect and manage family account data from coding agents via the `linden` CLI.

```
npx skills add mylinden-tech/skills --skill '*' -y
```

`--skill '*'` installs all skills in the package at once (`-y` skips prompts). See [install.md](install.md) for full setup including CLI authentication and the optional Claude Code plugin.

## Available skill

One skill, `linden`. Domain behavior lives in references. `/linden reminders` reads that reference. There is no `/linden-reminders` command.

| Reference | Use when the user talks about |
|---|---|
| [SKILL.md](skills/linden/SKILL.md) | Routing, decision tree, and `/linden [domain]` |
| [doctor](skills/linden/references/doctor.md) | Login, setup, or a CLI that is not working |
| [auth-and-accounts](skills/linden/references/auth-and-accounts.md) | Login, logout, or switching accounts |
| [account-settings](skills/linden/references/account-settings.md) | Account details, statistics, and user settings |
| [account-shares](skills/linden/references/account-shares.md) | Who a resource is shared with |
| [share-links](skills/linden/references/share-links.md) | Public or resource share links |
| [contacts](skills/linden/references/contacts.md) | Contacts, including a family doctor or lawyer |
| [invitations](skills/linden/references/invitations.md) | Pending invitations |
| [memberships](skills/linden/references/memberships.md) | Members and roles |
| [persons](skills/linden/references/persons.md) | People and family members |
| [pets](skills/linden/references/pets.md) | Pets |
| [reminders](skills/linden/references/reminders.md) | Reminders |
| [todos](skills/linden/references/todos.md) | Todos and todo lists |
| [vehicles](skills/linden/references/vehicles.md) | Vehicles |
| [real-estates](skills/linden/references/real-estates.md) | Homes and property |
| [online-accounts](skills/linden/references/online-accounts.md) | Website logins and subscriptions |
| [insurances](skills/linden/references/insurances.md) | Policies, providers, and expiring coverage |
| [wills](skills/linden/references/wills.md) | Will metadata |

`--skill '*'` installs the one skill in the package.

## Requires

- [linden CLI](https://github.com/mylinden-tech/linden-cli) (`brew install linden` or see the CLI README)
- An authenticated Linden account (`linden auth login`)

## Practical rule for agents

- Need next-step hints → `--json`
- Just need the payload (status, doctor checks, delete result) → `--agent`
- You're typing as a human → no flag is fine

Do not combine `--agent --jq` (agent wins; jq is ignored). Use `--jq` or `--json --jq` to filter `data` only — both skip the envelope. For breadcrumbs/summary, use `--json` alone.

## Try it

After install and `linden auth login`, ask your agent:

1. "Run Linden doctor and tell me if everything is set up."
2. "List my Linden accounts."
3. "List people in my Linden account."
4. "Show details for [a person from the list]."
5. "List pets in my Linden account."
6. "List vehicles in my Linden account."
7. "List my Linden documents." — should say the CLI has no documents command (do not invent one).

Or run yourself:

```sh
linden doctor --agent
linden accounts list --json
linden persons list --json
linden pets list --json
linden vehicles list --agent
```

## About

Skill bodies are authored in [mylinden-tech/linden-cli](https://github.com/mylinden-tech/linden-cli) under `skills/` and published here. To contribute, open issues or PRs in the CLI repo — do not edit skill content in this repo directly.
