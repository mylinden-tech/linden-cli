#!/usr/bin/env bash
set -euo pipefail
# Usage: ./scripts/publish-skills.sh /path/to/mylinden-tech/skills
dest="${1:-}"
if [[ -z "$dest" || ! -d "$dest" ]]; then
  echo "usage: $0 /path/to/mylinden-tech/skills" >&2
  exit 1
fi
root="$(cd "$(dirname "$0")/.." && pwd)"
rm -rf "$dest/skills"
mkdir -p "$dest/skills"
cp -R "$root/skills/linden" "$dest/skills/linden"
cp -R "$root/skills/linden-doctor" "$dest/skills/linden-doctor"
cp -R "$root/skills/linden-persons" "$dest/skills/linden-persons"
rm -f "$dest/skills/linden"/*.go "$dest/skills/"*.go 2>/dev/null || true
find "$dest/skills" -name '*_test.go' -delete
cp "$root/skills-publish/README.md" "$dest/README.md"
cp "$root/skills-publish/install.md" "$dest/install.md"
cp "$root/skills-publish/LICENSE" "$dest/LICENSE"
mkdir -p "$dest/.claude-plugin"
cp "$root/skills-publish/.claude-plugin/plugin.json" "$dest/.claude-plugin/plugin.json"
if [[ -f "$root/LICENSE" ]]; then
  cp "$root/LICENSE" "$dest/LICENSE"
fi
