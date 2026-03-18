# Architecture

This document describes the high-level layout and data flow of kcli.

## Overview

kcli is a Go CLI that:

1. Parses commands with **Cobra** and applies global flags (context, namespace, kubeconfig, force).
2. For **kubectl-style verbs**, builds the full argument list and runs the **kubectl** binary via `internal/runner`.
3. For **native features** (context, namespace, search, observability, incident, AI, TUI, plugins), implements logic in `internal/cli` and related packages.
4. Reads **config** from `~/.kcli/config.yaml` (`internal/config`) and **state** from `~/.kcli/state.json` (`internal/state`).
5. Optionally invokes **plugins** from `~/.kcli/plugins/` (`internal/plugin`) when the first argument is not a builtin.

## Directory layout

```
kcli/
├── cmd/kcli/          # Main entrypoint; plugin dispatch, then NewRootCommand()
├── pkg/api/           # Public API for embedding (Execute, ExecuteStream)
├── internal/
│   ├── ai/            # AI client, providers, prompts, usage/cost tracking
│   ├── cli/           # All Cobra commands (root, context, ns, search, observability, incident, ai, config, plugin, ui, logs, completion)
│   ├── config/        # Config load/save/validate (~/.kcli/config.yaml)
│   ├── k8sclient/     # Kubernetes client and caching (when used)
│   ├── plugin/        # Plugin discovery, manifest, install/update, policy (allow/revoke)
│   ├── runner/         # Kubectl execution (RunKubectl, CaptureKubectl, confirmation prompt)
│   ├── state/         # Context history, favorites, groups (~/.kcli/state.json)
│   ├── ui/            # Bubble Tea TUI (list, detail, XRay, bulk, themes)
│   └── version/       # Version, commit, build date
├── docs/              # User and developer documentation
└── scripts/           # Alpha smoke, perf check, release gate, official plugins
```

## Command flow

1. **Entry:** `cmd/kcli/main.go` parses `os.Args`. If the first argument is not a builtin, it tries to run it as a plugin (`plugin.TryRunForArgs` with `cli.IsBuiltinFirstArg`). If that fails or the first arg is builtin, it runs `cli.NewRootCommand().Execute()`.
2. **Root:** `internal/cli/root.go` builds the Cobra root with persistent flags (force, context, namespace, kubeconfig, ai-timeout, completion-timeout) and adds all subcommands. `PersistentPreRunE` validates config and rejects `--context -`.
3. **Kubectl verbs:** Commands like `get`, `describe`, `apply`, `delete`, `logs`, `exec`, etc., are implemented via `newKubectlVerbCmd` or specialized commands (e.g. `newLogsCmd`). They call `app.runKubectl(args)` or `app.captureKubectl(args)`. The `app` struct holds config, cache, and AI client; it prepends context/namespace/kubeconfig via `scopeArgsFor` and passes the full argv to `runner.RunKubectl` or `runner.CaptureKubectl`.
4. **Runner:** `internal/runner/kubectl.go` runs `exec.Command("kubectl", args...)`. For mutating verbs it may prompt for confirmation unless `ExecOptions.Force` is true. It supports stdin/stdout/stderr and optional timeout for capture.
5. **State and config:** Context switches and group operations load/save `state.Store` (`internal/state`). Config is loaded once at root init from `config.Load()` and can be updated by `kcli config set` and then re-read in the same process where the app is updated.

## Data flow

- **Kubeconfig:** Read by kubectl; kcli only passes `--context`, `--namespace`, `--kubeconfig` in argv. No in-process Kubernetes client is required for basic verbs (kubectl does the work).
- **Config:** YAML at `~/.kcli/config.yaml` → `config.Load()` → `Config` struct. Used for TUI (refresh, theme), logs (maxPods, colors), AI (provider, model, budget), completion timeout, etc.
- **State:** JSON at `~/.kcli/state.json` → `state.Load()` / `state.Save()`. Tracks last context, recent contexts, favorites, context groups.
- **AI:** `internal/ai` provides a client that calls configured providers (OpenAI, Anthropic, Azure, Ollama, custom). Usage is persisted under the kcli home directory and used for budget enforcement.
- **Plugins:** Discovered under `~/.kcli/plugins/` (executables named `kcli-<name>`). Optional `plugin.yaml` (manifest) and `plugin-policy.json` (allowed permissions). Invoked via `exec.Command` with args; stdin/stdout/stderr forwarded.

## TUI

- **Stack:** Bubble Tea + Lip Gloss (`internal/ui/tui.go`). Single program; no separate process.
- **Data:** Options (context, namespace, kubeconfig, AI callback, refresh interval, theme) are passed in. Resource lists and details are fetched by running kubectl (or equivalent) in the background and parsing output. XRay builds a tree from deployment/service/pod relationships using kubectl get/describe.
- **Modes:** List view (filter, sort, multi-select), detail view (tabs: overview, events, YAML, AI), XRay view, bulk action (delete, scale), confirm mode.

## Public API (embedding)

`pkg/api` exposes:

- `NewKCLI(cfg)` — create a client with optional stdin and env.
- `Execute(command string)` — run a single command string (e.g. `"get pods -n default"`) and return combined stdout/stderr.
- `ExecuteStream(command string)` — run and stream stdout/stderr via a channel of `StreamChunk`.

This allows servers (e.g. Kubilitics backend) to run kcli commands programmatically. Parsing is done by splitting the command string into words and passing them to the same Cobra root (with I/O redirected).

## Performance considerations

- **Startup:** Config and state are read once; AI client is lazy-initialized. Minimize work in init.
- **Completion:** Completion handlers call `captureKubectlWithTimeout` with a short timeout (e.g. 250ms). Results are cached in the app’s cache map with TTL to avoid repeated kubectl calls.
- **TUI:** Refresh is on a timer (configurable). Large lists can be paginated or limited; XRay and detail views issue kubectl on demand.

## Dependencies

- **Cobra** — CLI framework.
- **Bubble Tea / Bubbles / Lip Gloss** — TUI.
- **gopkg.in/yaml.v3** — Config and plugin manifest.
- **Go standard library** — exec, os, path, json, etc.
- **kubectl** — External binary; not a Go dependency.

Optional: Kubernetes client-go is used only where the codebase uses an in-process client (e.g. some discovery or caching); the primary path uses kubectl.
