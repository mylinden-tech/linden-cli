# Agent Skills for Linden

[Agent skills](https://agentskills.io) for [Linden](https://www.mylinden.family/) — manage family accounts and persons from coding agents via the `linden` CLI.

```
npx skills add mylinden-tech/skills --skill '*' -y
```

`--skill '*'` installs all skills in the package at once (`-y` skips prompts). See [install.md](install.md) for full setup including CLI authentication and the optional Claude Code plugin.

## Available skills

| Skill | Description |
|-------|-------------|
| [linden](skills/linden/SKILL.md) | Orchestrator — invariants, auth, accounts, routing to domain skills |
| [linden-doctor](skills/linden-doctor/SKILL.md) | Diagnose CLI, auth, API, and active account |
| [linden-persons](skills/linden-persons/SKILL.md) | List, show, create, update, and delete persons |

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
5. "List vehicles in Linden." — should say the CLI has no vehicles command (do not invent one).

Or run yourself:

```sh
linden doctor --agent
linden accounts list --json
linden persons list --json
```

## About

Skill bodies are authored in [mylinden-tech/linden-cli](https://github.com/mylinden-tech/linden-cli) under `skills/` and published here. To contribute, open issues or PRs in the CLI repo — do not edit skill content in this repo directly.
