# Install Linden Agent Skills

Install Agent Skills so your coding agent can manage Linden family accounts and persons via the `linden` CLI.

**Done when:** `linden --version` succeeds, `linden auth status` shows you are authenticated (or you complete login), and the `linden` skill is installed in your agent.

## Step 1: Install the CLI

If `linden` is not on PATH, install it:

**Homebrew:**

```sh
brew tap mylinden-tech/tap
brew install linden
```

**Install script:**

```sh
curl -fsSL https://raw.githubusercontent.com/mylinden-tech/linden-cli/main/scripts/install.sh | bash
```

See the [linden-cli README](https://github.com/mylinden-tech/linden-cli) for other options.

## Step 2: Authenticate

```sh
linden auth login
linden auth status
```

## Step 3: Install skills

One package ships three skills (`linden`, `linden-doctor`, `linden-persons`). Install them all in a single command — do not pick them one-by-one in the prompt:

```sh
npx skills add mylinden-tech/skills --skill '*' -y
```

`--skill '*'` selects every skill in the package; `-y` skips confirmation prompts. This uses the [Agent Skills](https://agentskills.io) open standard. The installer auto-detects your agent and places skills in the correct directory.

**Global install (all projects):**

```sh
npx skills add mylinden-tech/skills -g --skill '*' -y
```

**Specific agent:**

```sh
npx skills add mylinden-tech/skills -a cursor --skill '*' -y
npx skills add mylinden-tech/skills -a claude-code --skill '*' -y
```

**Local dry-run (from `linden-cli` before the public skills repo exists):**

```sh
mkdir -p .publish-staging
./scripts/publish-skills.sh .publish-staging
npx skills add ./.publish-staging --skill '*' -y
```

Restart your agent session to pick up the new skills.

### Claude Code plugin (optional)

The same package includes `.claude-plugin/plugin.json` (`linden-skills`). In Claude Code you can install that **one plugin** instead of using `npx`; it loads all three skills together:

```
/plugin install /absolute/path/to/.publish-staging
```

(or the published skills repo path once it exists). Skills-only (`npx`) and the Claude plugin stay independent — use either.

## Verify

```sh
linden doctor --agent
linden accounts list --json
```

**Expected:**

- `linden doctor --agent` prints JSON with `"ok": true` and checks for `auth`, `api`, and `account` all `"ok": true`.
- `linden accounts list --json` prints an envelope with `"ok": true`, a `data` array of accounts, and a `breadcrumbs` entry suggesting `linden accounts use <id>`.

In the agent, ask about Linden persons or run a doctor check. The agent should use `linden` commands with `--agent` or `--json`, not invent unsupported resource commands.

## Practical rule for agents

- Need next-step hints → `--json`
- Just need the payload (status, doctor checks, delete result) → `--agent`
- You're typing as a human → no flag is fine

Do not combine `--agent --jq` (agent wins; jq is ignored). Use `--jq` or `--json --jq` to filter `data` only — both skip the envelope. For breadcrumbs/summary, use `--json` alone.

## Manual installation

Clone this repo and symlink skills into your agent's skill directory:

```sh
git clone https://github.com/mylinden-tech/skills ~/.linden-skills
mkdir -p ~/.cursor/skills
ln -sfn ~/.linden-skills/skills/linden ~/.cursor/skills/linden
ln -sfn ~/.linden-skills/skills/linden-doctor ~/.cursor/skills/linden-doctor
ln -sfn ~/.linden-skills/skills/linden-persons ~/.cursor/skills/linden-persons
```

Update with `cd ~/.linden-skills && git pull`. Symlinks pick up changes immediately.
