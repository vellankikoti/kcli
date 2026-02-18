# Troubleshooting

Common issues and how to fix them.

## kcli fails with "invalid config"

**Cause:** `~/.kcli/config.yaml` has invalid YAML or a value that fails validation (e.g. theme not in allowed list, numeric out of range).

**Fix:**

1. Run `kcli config view` — if it fails, the file is invalid.
2. Edit: `kcli config edit` and fix the reported key, or run `kcli config set <key> <value>` with a valid value.
3. To start fresh: `kcli config reset --yes` (resets to defaults).

## "kubectl" not found or command fails

**Cause:** kcli delegates cluster operations to the `kubectl` binary. If `kubectl` is missing or not on `PATH`, or returns an error, kcli will fail.

**Fix:**

1. Install kubectl and ensure it’s on your `PATH`: `kubectl version --client`.
2. Use the same kubeconfig as kubectl: `kubectl config get-contexts` and `kcli ctx` should list the same contexts.
3. If you use a custom kubeconfig: `kcli --kubeconfig /path/to/config get pods` or set `KUBECONFIG`.

## Mutating command prompts "Proceed? [y/N]"

**Cause:** Safety prompt for commands that change cluster state (e.g. delete, apply, patch).

**Fix:**

- Type `y` and Enter to proceed, or `n` / Enter to abort.
- In scripts or CI, use `kcli --force delete ...` to skip the prompt (only when intended).

## AI commands say "AI disabled"

**Cause:** No provider configured or AI disabled in config.

**Fix:**

1. Set a provider and credentials: `kcli ai config --provider openai --key sk-... --enable` (or use env vars; see [AI guide](ai-guide.md)).
2. Check: `kcli ai status`. Ensure “Enabled” is true and provider/model are set.
3. If you use only env vars, ensure `KCLI_AI_PROVIDER` and the provider-specific env vars are set in the shell that runs kcli.

## AI "hard limit reached" or "soft limit reached"

**Cause:** Monthly AI budget in config is exceeded (100%) or at soft limit (default 80%).

**Fix:**

- Raise budget: `kcli ai config --budget 100`.
- Or wait until the next month (usage resets per month).
- Check usage: `kcli ai usage` and `kcli ai cost`.

## Completion is slow or times out

**Cause:** Completion calls kubectl (e.g. `api-resources`, `get ... -o name`). Slow cluster or network can exceed the completion timeout.

**Fix:**

1. Increase timeout: `kcli --completion-timeout 500ms get pods <TAB>` (or set in a wrapper).
2. Ensure cluster is reachable and kubeconfig context is correct.
3. On very large clusters, completion may remain slow; consider filtering by namespace (`-n`) to reduce list size.

## Plugin not found or not executable

**Cause:** Plugin must live under `~/.kcli/plugins/` and be named `kcli-<name>`, or be on `PATH` with `KCLI_PLUGIN_ALLOW_PATH=1`. It must be executable.

**Fix:**

1. Install into plugins dir: `kcli plugin install <source>`.
2. If you use a custom path: `chmod +x ~/.kcli/plugins/kcli-<name>`.
3. To allow PATH plugins: `export KCLI_PLUGIN_ALLOW_PATH=1` (use only if you trust PATH).

## Plugin "permission denied" or "pending-approval"

**Cause:** Plugin declares permissions in `plugin.yaml`; they must be approved before use.

**Fix:**

- `kcli plugin allow <name>` — approve all declared permissions.
- `kcli plugin inspect <name>` — see which permissions are pending.

## TUI is slow or flickers

**Cause:** Short refresh interval or heavy cluster.

**Fix:**

1. Increase refresh interval: `kcli config set tui.refresh_interval 5s`.
2. Disable animations: `kcli config set tui.animations false`.
3. Use a smaller namespace scope (e.g. `-n default`) when launching: `kcli -n default ui`.

## "no previous context recorded" for `kcli ctx -`

**Cause:** You haven’t switched context at least once in this environment, so there’s no “previous” context stored.

**Fix:** Switch at least once: `kcli ctx some-context`, then `kcli ctx -` will switch back.

## Multi-cluster get/search fails on some contexts

**Cause:** One or more contexts in `--all-contexts` or in the active context group are unreachable (wrong kubeconfig, network, or auth).

**Fix:**

1. Test each context: `kubectl --context <name> get ns`.
2. Remove bad contexts from the group: `kcli ctx group remove <group> <context>`.
3. Fix or remove the context from your kubeconfig.

## Logs --ai-summarize / --ai-errors fail

**Cause:** AI is disabled, over budget, or the provider request failed.

**Fix:** Same as “AI disabled” and “hard limit” above. Ensure AI works first: `kcli ai explain pod/some-pod`.

## Build from source fails

**Cause:** Wrong Go version or missing dependencies.

**Fix:**

1. Use Go 1.21 or later: `go version`.
2. In repo: `go mod download` then `go build -o bin/kcli ./cmd/kcli`.
3. If the repo uses CGO or external tools, check the project’s [Building from Source](../developer/building.md) and README.

## Getting help

- **Command help:** `kcli <command> --help` (e.g. `kcli plugin --help`, `kcli ai config --help`).
- **Version:** `kcli version`.
- **Config path:** Shown in errors when config is invalid; default is `~/.kcli/config.yaml`.
- For bugs or feature requests, open an issue in the project repository with your `kcli version` and steps to reproduce.
