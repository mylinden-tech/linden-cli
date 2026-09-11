# Publishing to mylinden-tech/skills

Skill bodies are authored in this monorepo under `skills/`. The public install target is `npx skills add mylinden-tech/skills`.

## Autopublish (preferred)

On every push to `main` that touches `skills/`, `skills-publish/`, or `scripts/publish-skills.sh`, the [Publish skills](../.github/workflows/publish-skills.yml) workflow:

1. Runs `go test ./skills`
2. Checks out `mylinden-tech/skills`
3. Runs `scripts/publish-skills.sh`
4. Commits and pushes only when the sync produces a diff

You can also run it manually via **Actions → Publish skills → Run workflow**.

### Required secret

| Secret | Purpose |
|---|---|
| `SKILLS_REPO_TOKEN` | PAT (or fine-grained token) with **Contents: Read and write** on `mylinden-tech/skills` |

Create the token in GitHub, then add it under **linden-cli → Settings → Secrets and variables → Actions**.

Do not use the default `GITHUB_TOKEN` — it cannot push to a different repository.

## Manual publish

Still useful for first-time repo setup or local dry runs.

### First-time setup

```sh
gh repo create mylinden-tech/skills --public --description "AI agent skills for Linden" --clone
cd skills   # or wherever gh cloned
# from linden-cli:
../linden-cli/scripts/publish-skills.sh "$(pwd)"
# or with absolute path:
# /path/to/linden-cli/scripts/publish-skills.sh /path/to/skills
git add -A
git commit -m "Initial Linden Agent Skills publish"
git push -u origin HEAD
```

### Updating after skill changes

```sh
/path/to/linden-cli/scripts/publish-skills.sh /path/to/mylinden-tech/skills
cd /path/to/mylinden-tech/skills
git add -A
git commit -m "Sync skills from linden-cli"
git push
```

### Local dry run

```sh
mkdir -p .publish-staging
./scripts/publish-skills.sh .publish-staging
```

Do not edit skill markdown in the publish repo by hand.
