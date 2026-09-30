# Persons examples

## List and filter

```bash
linden persons list --json
linden persons list --jq '.[].id'
linden persons list --page 1 --size 50 --json
```

## Show by UUID

```bash
linden persons show <uuid> --json
```

## Resolve a name, then show

```bash
linden persons list --json
# Read id from the matching entry, then:
linden persons show <uuid> --json
```

Filter with jq when the list is long:

```bash
linden persons list --json --jq '.[] | select(.first_name == "Jane") | .id'
```

## Create

```bash
linden persons create \
  --first-name Jane \
  --last-name Doe \
  --email jane@example.com \
  --json
```

Follow the response breadcrumb to show the new person.

## Update

```bash
linden persons update <uuid> --phone "+1 555 111 2222" --json
```

## Chaining with breadcrumbs

```bash
linden persons list --json
# breadcrumbs suggest: linden persons show <id>
linden persons show <uuid> --json
# breadcrumbs suggest update commands
```

Breadcrumbs require `--json`, not `--agent`.
