# Changelog

## [0.3.0] - 2026-09-29

### Breaking

- Domain skills (`linden-persons`, `linden-doctor`, and the rest) are merged into `linden` as references. `/linden-reminders` and the other domain slash commands no longer exist. Use `/linden reminders` (and the same form for the other references).

### Added

- Shared decision tree and `/linden [domain]` argument.
- User-language routing table.
- `--agent --help` JSON for flags, required markers, enum values, and agent notes.
- Routing fixture at `evals/routing.json` (reviewed, not a CI gate).
- Contract test in `skills/check_test.go`.
