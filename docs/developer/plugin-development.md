# Plugin Development

This guide explains how to build a kcli plugin that can be installed and run with `kcli plugin run <name>` or `kcli <name>`.

## Requirements

- Plugin **executable** name must be `kcli-<name>`, where `<name>` matches the pattern `[a-z0-9][a-z0-9-]*`.
- Plugin must live under **`~/.kcli/plugins/`** (or be on `PATH` only if `KCLI_PLUGIN_ALLOW_PATH=1` is set).
- Optional: **manifest** file `plugin.yaml` in the same directory as the executable for metadata and permissions.

## Manifest: plugin.yaml

Place `plugin.yaml` next to the executable (e.g. `~/.kcli/plugins/kcli-mytool` and `~/.kcli/plugins/plugin.yaml` for a single-plugin dir, or use a subdir per plugin if your install process does). Example:

```yaml
name: mytool
version: 1.0.0
author: Your Name
description: Short description of what the plugin does
commands:
  - mt
permissions:
  - read:pods
  - read:deployments
  - write:deployments
```

- **name** — Must match the executable name without the `kcli-` prefix (e.g. `mytool` for `kcli-mytool`).
- **version** — Semantic version (e.g. `1.0.0`).
- **author**, **description** — Shown in `kcli plugin inspect` and search.
- **commands** — Alternate names users can type (e.g. `kcli mt`).
- **permissions** — List of permission identifiers. Users must run `kcli plugin allow <name>` before the plugin can use them (enforcement is policy-based; the plugin binary is still executed, but kcli tracks approvals).

## Implementing the binary

- The plugin is invoked as: `kcli-mytool [args...]` where args are whatever the user passed after `kcli mytool` or `kcli plugin run mytool -- ...`.
- **stdin / stdout / stderr** are inherited from kcli. Read from stdin and write to stdout/stderr as needed.
- **Environment** — kcli does not inject extra env vars by default. Use `KUBECONFIG` or `KCLI_*` if you need to detect kcli context; otherwise rely on flags or args.
- **Exit code** — Use `os.Exit(0)` for success and non-zero for failure so scripts can detect errors.

Example (Go):

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("mytool 1.0.0")
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "Usage: kcli mytool [options]")
	os.Exit(1)
}
```

Build:

```bash
go build -o kcli-mytool .
```

Install (copy to plugin dir):

```bash
mkdir -p ~/.kcli/plugins
cp kcli-mytool ~/.kcli/plugins/
chmod +x ~/.kcli/plugins/kcli-mytool
# Optional: copy plugin.yaml next to it or in same dir
kcli plugin list
kcli mytool --version
```

## Permissions

- Declare what the plugin needs in `plugin.yaml` under `permissions`. Use a consistent naming scheme (e.g. `read:pods`, `write:deployments`).
- Users approve with `kcli plugin allow mytool` or `kcli plugin allow mytool read:pods write:deployments`.
- Revoke with `kcli plugin revoke mytool`. kcli stores approvals in `~/.kcli/plugin-policy.json`. The plugin binary does not receive permission info via env; permission checks are on the kcli side (e.g. before running the plugin or in future middleware). For now, document what your plugin needs and rely on users to run `plugin allow`.

## Discovery and install from GitHub

- **Install from local path:** `kcli plugin install ./path/to/plugin` — kcli looks for an executable or buildable Go module.
- **Install from GitHub:** `kcli plugin install github.com/org/repo` — kcli clones the repo and builds (e.g. main package in repo root or a documented path). The built binary and optional manifest are placed under `~/.kcli/plugins/` and the source is recorded for `kcli plugin update <name>`.

Implementing install/update is in `internal/plugin` (InstallFromSource, UpdateInstalled, etc.). Your repo should have a build that produces `kcli-<name>` and optionally ship a `plugin.yaml`.

## Testing

- **Locally:** Build `kcli-<name>`, copy to `~/.kcli/plugins/`, run `kcli plugin list`, `kcli plugin inspect <name>`, `kcli <name> [args]`.
- **Allow permissions:** `kcli plugin allow <name>` then run the plugin again if it checks policy.
- **Uninstall:** `kcli plugin remove <name>` and confirm the binary and manifest are removed.

## Best practices

- **Keep names lowercase with hyphens** (e.g. `my-tool`) to match `[a-z0-9][a-z0-9-]*`.
- **Version in manifest** so `kcli plugin inspect` and marketplace show it.
- **Document required permissions** in the plugin’s README so users know what to allow.
- **Handle --help and --version** so `kcli mytool --help` and `kcli mytool --version` are useful.
- **Do not assume kubeconfig path** — respect `KUBECONFIG` or document that users should set context via kcli before running the plugin.
