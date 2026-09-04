# Person fields

Flags match `linden persons create` and `linden persons update` in the CLI.

## Create flags

| Flag | Required | Notes |
|---|---|---|
| `--first-name` | yes | |
| `--last-name` | yes | |
| `--email` | no | |
| `--phone` | no | |
| `--notes` | no | |
| `--birthday` | no | Format `YYYY-MM-DD` |
| `--address-street` | no | |
| `--address-apt` | no | Apt/suite |
| `--address-city` | no | |
| `--address-state` | no | State/province |
| `--address-zip` | no | ZIP/postal code |
| `--address-country` | no | |
| `--relationship-type` | no | Free-form string |
| `--active` | no | Default `true` |

When required name flags are missing and stdout is not a TTY, the CLI returns a usage error instead of opening a form.

## Update flags

Same field flags as create, all optional. Only pass flags being changed — omitted flags are not sent to the API.

Additional update-only flags:

| Flag | Notes |
|---|---|
| `--active` | Mark person as active |
| `--inactive` | Mark person as inactive |

## Example create

```bash
linden persons create \
  --first-name Jane \
  --last-name Doe \
  --email jane@example.com \
  --phone "+1 555 000 0000" \
  --birthday 1990-01-15 \
  --relationship-type spouse \
  --json
```

## Example partial update

```bash
linden persons update <uuid> --email new@example.com --phone "+1 555 111 2222" --json
```
