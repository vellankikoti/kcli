# Getting Started

This guide gets you from zero to running kcli against a Kubernetes cluster.

## Prerequisites

- **Go 1.21+** (if building from source), or a [released kcli binary](installation.md) for your platform
- **kubectl** installed and on your `PATH` (kcli delegates to kubectl for cluster communication)
- A **kubeconfig** (e.g. `~/.kube/config`) with at least one context

## 1. Install kcli

Choose one:

- **From source:** See [Building from Source](../developer/building.md).
- **Pre-built binary:** See [Installation](installation.md) for your OS.

Verify:

```bash
kcli version
```

Example output: `kcli 0.1.0 (commit abc123, built 2026-02-16)`

## 2. Shell completion (recommended)

```bash
# Bash
source <(kcli completion bash)
# Or persist: kcli completion bash > /etc/bash_completion.d/kcli

# Zsh
source <(kcli completion zsh)
# Or: kcli completion zsh > "${fpath[1]}/_kcli"
```

## 3. Use your existing kubeconfig

kcli uses the same kubeconfig as kubectl. No extra setup is required.

```bash
# List contexts
kcli ctx

# Switch context (e.g. to prod)
kcli ctx my-prod-context

# Switch to previous context
kcli ctx -

# Set default namespace for current context
kcli ns kube-system
```

## 4. Core commands (kubectl parity)

All standard kubectl verbs are supported. Examples:

```bash
kcli get pods -A
kcli get pods -n default
kcli describe pod my-pod -n default
kcli logs my-pod -n default --tail=100
kcli exec -it pod/my-pod -n default -- /bin/sh
kcli get deploy,svc -n default
```

Global flags apply to every command:

- `--context <name>` — use a specific context
- `-n, --namespace <name>` — use a specific namespace
- `--kubeconfig <path>` — use a specific kubeconfig file
- `--force` — skip confirmation prompts for mutating commands

## 5. Multi-cluster and search

```bash
# Run get across all contexts
kcli get pods -A --all-contexts

# Use a context group (define with: kcli ctx group set prod ctx1 ctx2)
kcli get deploy -A --context-group prod

# Search resource names across contexts
kcli search my-app
kcli search api --context-group prod --kinds deployments,services,ingresses
```

## 6. Observability shortcuts

```bash
kcli health              # Cluster health summary
kcli health pods         # Pod health breakdown
kcli health nodes        # Node health
kcli restarts            # Pods sorted by restart count
kcli instability         # Restarts + recent warning events
kcli events --recent 30m --output json
kcli metrics             # Same as kubectl top pods -A
```

## 7. Incident mode

For quick triage during incidents:

```bash
kcli incident
kcli incident --recent 2h --restarts-threshold 5 --output json
kcli incident --watch --interval 10s
```

Subcommands: `kcli incident logs <ns/pod>`, `kcli incident describe <ns/pod>`, `kcli incident restart <ns/pod>`.

## 8. Optional AI

If you configure an [AI provider](ai-guide.md), you can use:

```bash
kcli ai explain deployment/my-app
kcli why pod/my-pod
kcli summarize events
kcli suggest fix deployment/my-app
kcli fix pod/my-pod --dry-run
```

## 9. Interactive TUI

```bash
kcli ui
```

See the [TUI guide](tui-guide.md) for keybindings (`/` filter, `:pods`, `:deploy`, `l` logs, `d` describe, `y` yaml, `A` AI, `q` quit).

## 10. Configuration and state

- **Config:** `~/.kcli/config.yaml` — see [Configuration](configuration.md).
- **State:** `~/.kcli/state.json` — context history and favorites (managed by kcli).

View or edit config:

```bash
kcli config view
kcli config set tui.refresh_interval 3s
kcli config get tui.refresh_interval
kcli config edit
```

## Safety model

Mutating commands (e.g. `delete`, `apply`, `patch`) prompt for confirmation unless you pass `--force`. Use `--force` only in scripts or when you intend to skip prompts.

## Next steps

- [Command reference](command-reference.md) — full list of commands and flags
- [Configuration](configuration.md) — all config keys and defaults
- [TUI guide](tui-guide.md) — interactive UI keybindings and views
- [AI guide](ai-guide.md) — enabling and using AI features
