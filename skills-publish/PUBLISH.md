# Publishing to mylinden-tech/skills

Skill bodies are authored in this monorepo under `skills/`. The public install target is `npx skills add mylinden-tech/skills`.

## First-time setup

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

## Updating after skill changes

```sh
/path/to/linden-cli/scripts/publish-skills.sh /path/to/mylinden-tech/skills
cd /path/to/mylinden-tech/skills
git add -A
git commit -m "Sync skills from linden-cli"
git push
```

Do not edit skill markdown in the publish repo by hand.
