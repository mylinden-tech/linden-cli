# Agent Skills for Linden

[Agent skills](https://agentskills.io) for [Linden](https://www.mylinden.family/) — inspect and manage family account data from coding agents via the `linden` CLI.

```
npx skills add mylinden-tech/skills --skill '*' -y
```

`--skill '*'` installs all skills in the package at once (`-y` skips prompts). See [install.md](install.md) for full setup including CLI authentication and the optional Claude Code plugin.

## Available skills

| Skill | Description |
|-------|-------------|
| [linden](skills/linden/SKILL.md) | Hub — invariants, auth, accounts, routing to domain skills |
| [linden-doctor](skills/linden-doctor/SKILL.md) | Diagnose CLI, auth, API, and active account |
| [linden-account-settings](skills/linden-account-settings/SKILL.md) | View account details, statistics, and current user settings |
| [linden-account-shares](skills/linden-account-shares/SKILL.md) | List and show account resource shares |
| [linden-contacts](skills/linden-contacts/SKILL.md) | List, show, and search contacts and contact types |
| [linden-invitations](skills/linden-invitations/SKILL.md) | List and show account invitations |
| [linden-memberships](skills/linden-memberships/SKILL.md) | List and show memberships and available roles |
| [linden-persons](skills/linden-persons/SKILL.md) | List, show, create, update, and delete persons |
| [linden-pets](skills/linden-pets/SKILL.md) | List, show, create, update, and delete pets |
| [linden-reminders](skills/linden-reminders/SKILL.md) | List, show, create, update, complete, and delete reminders |
| [linden-todos](skills/linden-todos/SKILL.md) | List, show, create, update, and delete todos and todo lists |
| [linden-vehicles](skills/linden-vehicles/SKILL.md) | List and show vehicles |
| [linden-real-estates](skills/linden-real-estates/SKILL.md) | List and show real estate |
| [linden-online-accounts](skills/linden-online-accounts/SKILL.md) | List and show online accounts |
| [linden-insurances](skills/linden-insurances/SKILL.md) | List and show insurance policies, providers, and expiring coverage |
| [linden-wills](skills/linden-wills/SKILL.md) | List and show will metadata |
| [linden-share-links](skills/linden-share-links/SKILL.md) | List and show resource-specific share links |

All domain skills set `disable-model-invocation: true`. Agents discover domain skills through the `linden` hub routing table—not via ambient auto-invocation. `--skill '*'` still installs every skill in the package.

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
