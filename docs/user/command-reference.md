# Command Reference

This page lists kcli commands and global flags. kubectl-style verbs pass remaining arguments to kubectl; see `kubectl --help` for verb-specific options.

## Global flags

These apply to the root command and can be used with any subcommand:

| Flag | Short | Description |
|------|-------|-------------|
| `--force` | | Skip confirmation prompts for mutating commands (e.g. delete, apply). |
| `--context` | | Override kubectl context for this run. |
| `--namespace` | `-n` | Override namespace. |
| `--kubeconfig` | | Path to kubeconfig file. |
| `--ai-timeout` | | AI request timeout (default: 30s). |
| `--completion-timeout` | | Timeout for shell completion lookups (default: 250ms). |

Examples:

```bash
kcli --context prod get pods
kcli -n kube-system get pods
kcli --force delete pod my-pod
```

---

## Core Kubernetes (kubectl parity)

These commands forward to `kubectl` with context/namespace applied. All standard kubectl flags are supported after the verb.

| Command | Aliases | Description |
|---------|---------|-------------|
| `kcli get [resource] [args...]` | `g` | Get resources. Supports `--all-contexts` and `--context-group <name>`. |
| `kcli describe <resource> <name> [args...]` | `desc` | Describe a resource. |
| `kcli apply [args...]` | `ap` | Apply configuration from file or stdin. |
| `kcli create [args...]` | `cr` | Create resources from file or stdin. |
| `kcli delete [args...]` | `del` | Delete resources (prompts unless `--force`). |
| `kcli logs [target] [args...]` | | Pod logs; supports selectors, multi-pod, `--grep`, `--save`, `--ai-summarize`, `--ai-errors`, `--ai-explain`. |
| `kcli exec [args...]` | | Execute command in a container. |
| `kcli port-forward [args...]` | | Forward local ports to a pod or service. |
| `kcli top [args...]` | | Resource usage (pods, nodes). |
| `kcli rollout [args...]` | | Rollout status, history, undo, restart. |
| `kcli diff [args...]` | | Diff live vs local. |
| `kcli explain [resource] [args...]` | | Field documentation. |
| `kcli wait [args...]` | | Wait for conditions. |
| `kcli scale [args...]` | | Scale workloads. |
| `kcli patch [args...]` | | Patch resources. |
| `kcli label [args...]` | | Update labels. |
| `kcli annotate [args...]` | | Update annotations. |
| `kcli edit [args...]` | | Edit resource on server. |
| `kcli kgp [args...]` | | Shortcut for `kcli get pods`. |
| `kcli auth can-i [args...]` | | Check whether an action is allowed. |

### Logs-specific behavior

- **Target:** pod name, `pod/name`, `deployment/name`, or label selector (e.g. `app=web`).
- **Flags (examples):** `--follow` / `-f`, `--tail=N`, `--since=DURATION`, `--timestamps`, `-c CONTAINER`, `--grep=REGEX`, `--grep-v=REGEX`, `--save=FILE`, `--ai-summarize`, `--ai-errors`, `--ai-explain`.

---

## Workflow

### Context

| Command | Description |
|---------|-------------|
| `kcli ctx` | List contexts (current marked). |
| `kcli ctx <name>` | Switch to context. |
| `kcli ctx -` | Switch to previous context. |
| `kcli ctx --favorites` / `-f` | List favorite contexts only. |
| `kcli ctx fav add <context>` | Add context to favorites. |
| `kcli ctx fav rm <context>` | Remove from favorites. |
| `kcli ctx fav ls` | List favorites. |
| `kcli ctx group` | List context groups; with `<name>` switch active group. |
| `kcli ctx group set <name> <context...>` | Create/replace a context group. |
| `kcli ctx group add <name> <context...>` | Add contexts to a group. |
| `kcli ctx group remove <name> [context...]` | Remove contexts (or delete group if none). |
| `kcli ctx group export [file]` | Export groups to JSON. |
| `kcli ctx group import <file> [--merge]` | Import groups from JSON. |

### Namespace

| Command | Description |
|---------|-------------|
| `kcli ns` | Print current namespace. |
| `kcli ns <name>` | Set default namespace for current context. |
| `kcli ns --list` / `-l` | List namespaces. |

### Search

| Command | Description |
|---------|-------------|
| `kcli search <query>` | Search resource names across contexts. |
| `kcli search <query> --context-group <name>` | Limit to a context group. |
| `kcli search <query> --kinds <list>` | Comma-separated kinds (e.g. `deployments,services,ingresses`). |

### Config

| Command | Description |
|---------|-------------|
| `kcli config view [--output yaml\|json]` | Show effective config. |
| `kcli config get <key>` | Get value by key (e.g. `tui.refresh_interval`). |
| `kcli config set <key> <value>` | Set value (persists to `~/.kcli/config.yaml`). |
| `kcli config reset --yes` | Reset config to defaults. |
| `kcli config edit` | Open config in `$EDITOR` or `$VISUAL`. |

---

## Observability

| Command | Description |
|---------|-------------|
| `kcli health` | Overall cluster health (pods + nodes summary). |
| `kcli health pods` | Pod health breakdown (Running, Pending, Failed, etc.). |
| `kcli health nodes` | Node health (Ready, pressure conditions). |
| `kcli metrics [resource]` | Same as `kubectl top`; default `top pods -A`. |
| `kcli restarts` | Pods sorted by restart count. |
| `kcli restarts --min-restarts N` | Only pods with at least N restarts. |
| `kcli instability` | Restart leaders + recent warning events. |
| `kcli events` | Cluster events. |
| `kcli events --recent 30m` | Events in last 30 minutes (default 1h without `--all`). |
| `kcli events --all` | All events (no time filter). |
| `kcli events --type Warning` | Filter by type. |
| `kcli events --output json` | JSON output. |
| `kcli events --watch` | Watch event stream. |

---

## Incident response

| Command | Description |
|---------|-------------|
| `kcli incident` | Incident summary (CrashLoop, OOM, high restarts, node pressure, critical events). |
| `kcli incident --recent 2h` | Time window for events (default 2h). |
| `kcli incident --restarts-threshold 5` | Min restarts for “high restarts” (default 5). |
| `kcli incident --output json` | JSON output. |
| `kcli incident --watch [--interval 10s]` | Auto-refresh (default interval 5s). |
| `kcli incident logs <namespace>/<pod> [--tail=N]` | Tail logs for a pod from incident list. |
| `kcli incident describe <namespace>/<pod>` | Describe pod. |
| `kcli incident restart <namespace>/<pod>` | Restart pod (delete; recreates if managed). |

---

## AI (optional)

Requires [AI configuration](ai-guide.md). Commands degrade gracefully if AI is disabled or unavailable.

| Command | Description |
|---------|-------------|
| `kcli ai <question>` | Natural-language query (e.g. “which pods are crashing?”). |
| `kcli ai explain [resource]` | Explain resource or concept. |
| `kcli ai why [resource]` | Explain probable cause of current state. |
| `kcli ai suggest-fix [resource]` | Suggest remediation. |
| `kcli why [resource]` | Same as `kcli ai why`. |
| `kcli summarize events [resource] [--since=6h]` | Summarize events. |
| `kcli suggest fix [resource]` | Same as suggest-fix. |
| `kcli fix [resource] [--dry-run]` | Fix suggestions; `--dry-run` only prints. |
| `kcli ai config` | Show AI config. |
| `kcli ai config --provider openai --model gpt-4o --enable` | Set provider, model, enable. |
| `kcli ai config --key sk-...` | Set API key (stored in config). |
| `kcli ai config --budget 50 --soft-limit 80` | Monthly budget (USD) and soft limit %. |
| `kcli ai status` | Runtime status, provider, budget utilization. |
| `kcli ai usage` | Monthly usage (calls, tokens, cost). |
| `kcli ai cost` | Cost and budget status. |

---

## Plugins

| Command | Description |
|---------|-------------|
| `kcli plugin list` | List installed plugins. |
| `kcli plugin search <keyword>` | Search installed + marketplace. |
| `kcli plugin marketplace` | List marketplace catalog. |
| `kcli plugin inspect <name>` | Show manifest and permissions. |
| `kcli plugin info <name>` | Alias for inspect. |
| `kcli plugin install <path-or-repo>` | Install from local path or `github.com/org/repo`. |
| `kcli plugin update <name>` | Update one plugin. |
| `kcli plugin update --all` | Update all installed. |
| `kcli plugin update-all` | Update all from recorded sources. |
| `kcli plugin remove <name>` | Uninstall plugin. |
| `kcli plugin allow <name> [permission...]` | Approve permissions. |
| `kcli plugin revoke <name> [permission...]` | Revoke permissions. |
| `kcli plugin run <name> [args...]` | Run plugin by name. |

Plugins can also be invoked as first-class commands: `kcli <plugin-name> [args...]` if the name is not a builtin.

---

## UI and version

| Command | Description |
|---------|-------------|
| `kcli ui` | Start interactive TUI (see [TUI guide](tui-guide.md)). |
| `kcli version` | Print version, commit, and build date. |
| `kcli completion bash` | Generate bash completion. |
| `kcli completion zsh` | Generate zsh completion. |
| `kcli completion fish` | Generate fish completion. |

---

## Resource types (get/describe)

kcli uses the same resource types as kubectl (including short names and API groups). Examples: `pods`, `po`, `deployments`, `deploy`, `services`, `svc`, `nodes`, `no`, `namespaces`, `ns`, `events`, `ev`, `configmaps`, `cm`, `secrets`, `ingresses`, `ing`, `persistentvolumeclaims`, `pvc`, `jobs`, `cronjobs`, `cj`, etc. Use `kubectl api-resources` for the full list.
