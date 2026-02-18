# FAQ

## General

### What is kcli?

kcli is a unified Kubernetes CLI that provides kubectl parity plus context/namespace ergonomics, observability shortcuts (health, restarts, events, incident mode), an interactive TUI, optional AI (explain, why, fix, summarize), and a plugin system. It delegates cluster API calls to kubectl.

### Do I need kubectl?

Yes. kcli runs `kubectl` under the hood for get, describe, logs, exec, apply, delete, etc. Install kubectl and ensure it’s on your `PATH`.

### Does kcli replace kubectl?

It wraps and extends the kubectl workflow. You can use kcli for day-to-day commands and keep using kubectl directly if you prefer. kcli adds context groups, search, incident mode, TUI, and AI; it does not replace the Kubernetes API.

### Where is config and state stored?

- **Config:** `~/.kcli/config.yaml` (themes, TUI, logs, AI, etc.).
- **State:** `~/.kcli/state.json` (last context, recent contexts, favorites, context groups).
- **Plugins:** `~/.kcli/plugins/` (executables and manifests).

---

## Context and namespace

### How do I switch context quickly?

- `kcli ctx <name>` — switch to a context.
- `kcli ctx -` — switch to the previous context.
- `kcli ctx` — list contexts (current is marked).
- Use `kcli ctx fav add <name>` for favorites and `kcli ctx --favorites` to list only favorites.

### What are context groups?

Context groups let you run the same command across multiple clusters. Define a group (e.g. “prod”) with `kcli ctx group set prod ctx1 ctx2`, then use `kcli get pods -A --context-group prod` or `kcli search my-app --context-group prod`.

### How does kcli know which namespace to use?

It uses the default namespace of the current kubeconfig context (same as kubectl). Override with `-n <namespace>` or set the default for the current context with `kcli ns <namespace>`.

---

## Safety and prompts

### Why does kcli ask "Proceed? [y/N]" for delete/apply?

To avoid accidental mutations. Type `y` to proceed. In scripts or automation, use `kcli --force` to skip the prompt when appropriate.

### Which commands prompt?

Mutating verbs: apply, create, delete, edit, patch, label, annotate, scale, rollout, etc. Read-only commands (get, describe, logs, top, etc.) do not prompt.

---

## AI

### Is AI required?

No. All features work without AI. AI commands (explain, why, fix, summarize) simply print a “disabled” message if no provider is configured.

### Where is AI usage stored?

Usage (calls, tokens, cost) is stored under the kcli config directory and tracked per month. View with `kcli ai usage` and `kcli ai cost`.

### Can I use a local model?

Yes. Use the **ollama** provider: `kcli ai config --provider ollama --model llama3 --enable`. Ensure Ollama is running (e.g. `ollama serve`).

### How do I avoid sending sensitive data to AI?

Use the **ollama** provider for fully local inference, or a **custom** endpoint that stays inside your network. For OpenAI/Anthropic/Azure, avoid pasting secrets into prompts; kcli sends resource metadata and event/log content as configured.

---

## Plugins

### How do I install a plugin?

From a local path: `kcli plugin install ./my-plugin`. From GitHub: `kcli plugin install github.com/org/kcli-mytool`. From the marketplace: `kcli plugin install cert-manager` (if available).

### Can I run a plugin without installing it?

Plugins must be under `~/.kcli/plugins/` (or on PATH with `KCLI_PLUGIN_ALLOW_PATH=1`). There is no “run from URL” without install.

### Why is my plugin "invalid"?

Usually the manifest `plugin.yaml` is missing or has invalid YAML/fields. Run `kcli plugin inspect <name>` to see the error. Ensure `name` matches the executable (without `kcli-`), and that `version` and other fields are valid.

---

## TUI

### How do I exit the TUI?

Press `q` or `Ctrl+C`.

### Can I change the TUI refresh rate?

Yes: `kcli config set tui.refresh_interval 5s` (or any duration like `2s`, `10s`).

### What are the themes?

**ocean**, **forest**, and **amber**. Set with `kcli config set general.theme forest` or `tui.theme`. In the TUI, **Ctrl+T** cycles themes.

---

## Performance and limits

### Why is startup slow?

kcli loads config and may initialize the AI client. For sub-200ms targets, ensure `~/.kcli/config.yaml` is small and avoid heavy shell hooks. Use `--completion-timeout` only if needed.

### How many pods can I tail with `kcli logs`?

Controlled by `logs.maxPods` in config (default 20, max 500). Increase with `kcli config set logs.max_pods 50` if needed.

### Does kcli cache API responses?

Yes, for completion and some internal use. TTL is controlled by `performance.cacheTTL` (default 60s). Context and namespace are applied to cache keys.

---

## Integration

### Can I use kcli in CI/CD?

Yes. Use `kcli --force` for mutating commands so prompts don’t block. Set `KUBECONFIG` or `--kubeconfig` and `--context` as needed. Avoid TUI and interactive AI in headless pipelines.

### Does Kubilitics embed kcli?

The Kubilitics platform can embed kcli (e.g. terminal in the UI). When embedded, context/namespace may sync with the UI. See Kubilitics docs for details.

### How do I get shell completion?

Run `kcli completion bash` or `kcli completion zsh` and source the output (or write it to your completion directory). See [Installation](installation.md#shell-completion).
