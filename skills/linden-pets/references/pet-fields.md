# Pet fields

Flags match `linden pets create` and `linden pets update` in the CLI.

## Create flags

| Flag | Required | Notes |
|---|---|---|
| `--name` | yes | Pet name |
| `--species` | yes | Species such as dog or cat |
| `--breed` | no | |
| `--color` | no | |
| `--notes` | no | |
| `--birth-date` | no | Format `YYYY-MM-DD` |
| `--microchip` | no | Microchip number |
| `--weight` | no | Include units when useful |
| `--is-active` | no | Default `true`; pass `--is-active=false` to create an inactive pet |

The CLI returns a usage error when `--name` or `--species` is missing.

## Update flags

The same field flags are available for update, and all are optional. Only pass flags being changed; omitted flags are not sent to the API.

Use `--is-active=false` to mark a pet inactive and `--is-active=true` to mark it active.

## Example create

```bash
linden pets create \
  --name Rex \
  --species dog \
  --breed "Golden Retriever" \
  --birth-date 2020-05-12 \
  --weight "30 kg" \
  --json
```

## Example partial update

```bash
linden pets update <uuid> --weight "31 kg" --color golden --json
```
