---
name: linden-doctor
description: |
  This skill should be used when the user asks to run Linden doctor, diagnose the Linden CLI,
  fix Linden login, or reports that Linden, auth, or the active account is not working.
---

# Linden Doctor

Diagnose CLI, auth, API, and active account. Do not dump config or credential files.

## Procedure

1. Run `linden doctor --agent`.
2. The success payload is the data object (no envelope): `{ "ok": bool, "checks": [ ... ] }`. Read `checks`.
3. Each check has `check`, `ok`, and `detail`. Current names: `auth`, `api`, `account`. The `api` check includes `url`. The `account` check includes `account_id` when set.
4. Failed `auth` → `linden auth login` (human completes the browser flow) → `linden auth status --agent`.
5. Failed `api` → report `url` from the check. Do not retry in a tight loop.
6. Failed `account` → `linden accounts list --json` → `linden accounts use <id>`.
7. If every check has `"ok": true` (and top-level `ok` is true), say Linden is ready.

Re-run `linden doctor --agent` after remediations.
