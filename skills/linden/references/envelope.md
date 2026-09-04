# Output envelope

## Success modes

**`--json`** prints the full envelope:

```json
{
  "ok": true,
  "data": {},
  "summary": "human-readable summary",
  "breadcrumbs": [
    {
      "action": "show",
      "cmd": "linden persons show <id>",
      "description": "View person details"
    }
  ],
  "meta": {}
}
```

Use `--json` when chaining commands. Read `breadcrumbs[].cmd` for the suggested next step.

**`--agent`** and **`--quiet`** print JSON of `data` only. There is no `ok`, `summary`, or `breadcrumbs` wrapper. Use `--agent` for scripting and data extraction.

**`--md`** prints Markdown for humans (summary heading plus data).

**`--jq '<expr>'`** filters the `data` field only and prints that filtered JSON. It always skips the envelope (`ok` / `summary` / `breadcrumbs`), even when combined with `--json`. Never pipe to external `jq`.

For breadcrumbs or summary, run with **`--json` alone** (no `--jq`). Filter payload fields in a separate `--jq` (or `--json --jq`) invocation afterward if needed.

**Never combine `--agent --jq`.** Agent mode wins in the CLI; jq is ignored. Filter with `--jq` alone or `--json --jq` (both still filter `data` only).

**`--styled`** forces ANSI output in a TTY. **`--account <id>`** overrides the active account for one invocation.

## Error envelope

When stdout is not a TTY, errors print JSON to stderr or stdout depending on the command path:

```json
{
  "ok": false,
  "error": "message",
  "code": "usage",
  "hint": "Run: linden auth login"
}
```

Error codes: `usage`, `not_found`, `auth`, `forbidden`, `rate_limit`, `network`, `api`.

Execute `hint` only when it names a real `linden` command from the catalog. Otherwise explain the error in plain language.

## Breadcrumb chaining

Breadcrumbs appear only with `--json`, not `--agent`. To chain:

1. Run a command with `--json`.
2. Read `breadcrumbs` from the response.
3. Run the `cmd` from the most relevant breadcrumb.

Do not invent follow-up commands when breadcrumbs are absent.
