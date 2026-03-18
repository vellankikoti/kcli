# kcli Architecture

## Overview

`kcli` is a thin, safety-first command layer around Kubernetes operations:

- Built-in workflows (`ctx`, `ns`, `incident`, `ui`, `ai`, `plugin`)
- kubectl-parity pass-through for remaining commands
- Guardrails for mutating operations
- Completion and performance caches for fast CLI UX

## Runtime flow

1. Entry point: `cmd/kcli/main.go`
2. Dispatch order:
   - Plugin command resolution (`kcli-<name>`)
   - Built-in commands (Cobra)
   - kubectl fallback for non-built-in verbs
3. Built-ins delegate to:
   - `internal/cli`: command wiring and UX behavior
   - `internal/runner`: guarded kubectl execution
   - `internal/k8sclient`: client-go access (auth/context checks)
   - `internal/state`: persisted local state (`~/.kcli/state.json`)
   - `internal/ui`: Bubble Tea TUI

## Key components

- `internal/cli`
  - Root command graph and flags
  - Context/namespace ergonomics
  - Observability/incident/AI commands
  - Completion engine and cache
- `internal/runner`
  - Command classification (mutating vs non-mutating)
  - Interactive confirmations and `--force` bypass
- `internal/k8sclient`
  - kubeconfig parsing
  - context/auth inspection
  - short-TTL in-process cache
  - parallel clientset/dynamic initialization
- `internal/ui`
  - Bubble Tea model/update/view loop
  - pod table + detail mode (Overview/Events/YAML)
- `internal/plugin`
  - plugin discovery, manifest validation, policy permissions

## Safety model

- Mutating verbs require explicit confirmation by default.
- CI/automation paths use `--force`.
- Plugin permissions are explicitly allowed/revoked and persisted.

## Performance model

- Startup path avoids unnecessary heavy initialization.
- Completion uses short-lived in-memory cache.
- kubeconfig and client-go setup use short-TTL in-process caching.
- `scripts/perf-check.sh` enforces startup/get/ctx/memory gates.
