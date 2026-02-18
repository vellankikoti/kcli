# Configuration Reference

kcli reads configuration from `~/.kcli/config.yaml`. The file is created with defaults on first use of `kcli config set` or `kcli config edit`. All keys are optional; missing values fall back to the defaults below.

## Config file location

- **Path:** `~/.kcli/config.yaml` (or `$HOME/.kcli/config.yaml`).
- **Create default:** `kcli config edit` or `kcli config set <key> <value>`.
- **View:** `kcli config view` (YAML) or `kcli config view --output json`.

## Configuration keys

Keys use dot notation. Case-insensitive for the CLI; in YAML use the exact key names below.

### General

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `general.theme` | string | `ocean` | Global theme: `ocean`, `forest`, or `amber`. |
| `general.startupTimeBudget` | duration | `250ms` | Target startup time (informational). |

### Context

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `context.recentLimit` | int | 10 | Max recent contexts to keep (1–1000). |
| `context.favorites` | list of strings | `[]` | Favorite context names. |
| `context.groups.<name>` | list of strings | — | Context group: list of context names. |

### TUI

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `tui.refreshInterval` | duration | `2s` | Auto-refresh interval for resource lists. |
| `tui.theme` | string | (inherits general) | TUI theme: `ocean`, `forest`, `amber`. |
| `tui.colors` | bool | true | Use colors in TUI. |
| `tui.animations` | bool | true | Enable animations. |

### Logs

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `logs.followNewPods` | bool | true | In multi-pod logs, follow new pods when scale increases. |
| `logs.maxPods` | int | 20 | Max pods to tail in one logs session (1–500). |
| `logs.colors` | bool | true | Color-code log lines by pod. |

### Performance

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `performance.cacheTTL` | duration | `60s` | TTL for completion/cache. |
| `performance.memoryLimitMB` | int | 256 | Soft memory limit in MB (64–65536). |

### Shell

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `shell.promptFormat` | string | `[{{.context}}/{{.namespace}}]$ ` | Prompt template (context, namespace). |
| `shell.aliases.<name>` | string | — | Custom alias: value is the expansion. |

### AI

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `ai.enabled` | bool | true | Enable AI commands when provider is set. |
| `ai.provider` | string | — | One of: `openai`, `anthropic`, `azure-openai`, `ollama`, `custom`. |
| `ai.model` | string | — | Default model name. |
| `ai.apiKey` | string | — | API key (stored in config; prefer env for secrets). |
| `ai.endpoint` | string | — | Custom endpoint (for custom/azure). |
| `ai.budgetMonthlyUSD` | float | 50 | Monthly budget in USD (0 = no limit). |
| `ai.softLimitPercent` | float | 80 | Warn when usage reaches this % of budget. |

## Setting and getting values

```bash
kcli config set general.theme forest
kcli config set tui.refresh_interval 3s
kcli config set context.recentLimit 20
kcli config set ai.budgetMonthlyUSD 100
kcli config get tui.refresh_interval
kcli config view
```

- **set** — Updates the key and writes `~/.kcli/config.yaml`. Invalid values are rejected.
- **get** — Prints the value (or empty for unset). For list values (e.g. `context.favorites`), output is comma-separated.
- **view** — Prints the full effective config (YAML or JSON with `--output json`).

## Resetting configuration

```bash
kcli config reset --yes
```

Restores all keys to defaults and overwrites `~/.kcli/config.yaml`. Requires `--yes`.

## Editing the file directly

```bash
kcli config edit
```

Opens the config file in `$VISUAL` or `$EDITOR`. Invalid YAML or invalid values may cause kcli to report an error on next load; fix the file or use `kcli config set` to correct keys.

## Validation rules

- **general.theme**, **tui.theme:** Must be one of `ocean`, `forest`, `amber`.
- **context.recentLimit:** 1–1000.
- **logs.maxPods:** 1–500.
- **performance.memoryLimitMB:** 64–65536.
- **ai.budgetMonthlyUSD:** ≥ 0.
- **ai.softLimitPercent:** &gt; 0 and &lt; 100.
- Duration fields accept values like `2s`, `30m`, `1h`.

## State vs config

- **Config** (`~/.kcli/config.yaml`) — User preferences (theme, intervals, AI, etc.). Edited via `kcli config` or `kcli config edit`.
- **State** (`~/.kcli/state.json`) — Runtime state: last context, recent contexts, favorites, context groups. Managed by kcli when you switch context or edit groups; do not edit by hand unless you know the format.
