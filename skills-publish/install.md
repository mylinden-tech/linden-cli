# Install Linden Agent Skills

Install Agent Skills so your coding agent can inspect and manage Linden family account data via the `linden` CLI.

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

One package ships seventeen skills: `linden` (hub), `linden-doctor`, and domain skills `linden-account-settings`, `linden-account-shares`, `linden-contacts`, `linden-invitations`, `linden-memberships`, `linden-persons`, `linden-pets`, `linden-reminders`, `linden-todos`, `linden-vehicles`, `linden-real-estates`, `linden-online-accounts`, `linden-insurances`, `linden-wills`, and `linden-share-links`. All domain skills are hub-routed (`disable-model-invocation: true`)—agents load them via the `linden` routing table, not ambient auto-invocation. Install them all in a single command—do not pick them one-by-one in the prompt:

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

The same package includes `.claude-plugin/plugin.json` (`linden-skills`). In Claude Code you can install that **one plugin** instead of using `npx`; it loads all skills together:

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

In the agent, ask about any supported Linden domain—or run a doctor check. The agent should route through the `linden` hub, load the matching domain skill, and use `linden` commands with `--agent` or `--json`. It should not invent unsupported resource commands (documents, passports, SSNs, files, driver licenses, birth certificates, body memorial wishes, legacy messages, or contact-import).

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
ln -sfn ~/.linden-skills/skills/linden-account-settings ~/.cursor/skills/linden-account-settings
ln -sfn ~/.linden-skills/skills/linden-account-shares ~/.cursor/skills/linden-account-shares
ln -sfn ~/.linden-skills/skills/linden-contacts ~/.cursor/skills/linden-contacts
ln -sfn ~/.linden-skills/skills/linden-invitations ~/.cursor/skills/linden-invitations
ln -sfn ~/.linden-skills/skills/linden-memberships ~/.cursor/skills/linden-memberships
ln -sfn ~/.linden-skills/skills/linden-persons ~/.cursor/skills/linden-persons
ln -sfn ~/.linden-skills/skills/linden-pets ~/.cursor/skills/linden-pets
ln -sfn ~/.linden-skills/skills/linden-reminders ~/.cursor/skills/linden-reminders
ln -sfn ~/.linden-skills/skills/linden-todos ~/.cursor/skills/linden-todos
ln -sfn ~/.linden-skills/skills/linden-vehicles ~/.cursor/skills/linden-vehicles
ln -sfn ~/.linden-skills/skills/linden-real-estates ~/.cursor/skills/linden-real-estates
ln -sfn ~/.linden-skills/skills/linden-online-accounts ~/.cursor/skills/linden-online-accounts
ln -sfn ~/.linden-skills/skills/linden-insurances ~/.cursor/skills/linden-insurances
ln -sfn ~/.linden-skills/skills/linden-wills ~/.cursor/skills/linden-wills
ln -sfn ~/.linden-skills/skills/linden-share-links ~/.cursor/skills/linden-share-links
```

Update with `cd ~/.linden-skills && git pull`. Symlinks pick up changes immediately.
