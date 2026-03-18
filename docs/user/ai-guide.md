# AI Features Guide

kcli can use an optional AI backend to explain resources, analyze failures, summarize events, and suggest fixes. AI is **optional**: all other features work without it.

## Enabling AI

Two ways to configure:

1. **Config file** — Edit `~/.kcli/config.yaml` (or use `kcli ai config`). See [Configuration](configuration.md).
2. **Environment variables** — Override or set provider and credentials via env (see below).

### Supported providers

| Provider | Config / env | Required | Optional |
|----------|--------------|----------|----------|
| **openai** | `ai.provider: openai` | API key | `ai.model`, `ai.endpoint` |
| **anthropic** | `ai.provider: anthropic` | API key | `ai.model` |
| **azure-openai** | `ai.provider: azure-openai` | API key, endpoint, deployment | `KCLI_AZURE_OPENAI_API_VERSION` |
| **ollama** | `ai.provider: ollama` | — | `ai.model`, `KCLI_OLLAMA_ENDPOINT` |
| **custom** | `ai.provider: custom` | `ai.endpoint` | `ai.apiKey` |

Environment overrides (examples):

- **OpenAI:** `KCLI_AI_PROVIDER=openai`, `KCLI_OPENAI_API_KEY` or `KCLI_AI_API_KEY`, optional `KCLI_AI_MODEL`, `KCLI_AI_ENDPOINT`.
- **Anthropic:** `KCLI_AI_PROVIDER=anthropic`, `KCLI_ANTHROPIC_API_KEY` or `KCLI_AI_API_KEY`, optional `KCLI_AI_MODEL`.
- **Azure OpenAI:** `KCLI_AI_PROVIDER=azure-openai`, `KCLI_AZURE_OPENAI_API_KEY`, `KCLI_AZURE_OPENAI_ENDPOINT`, `KCLI_AZURE_OPENAI_DEPLOYMENT`, optional `KCLI_AZURE_OPENAI_API_VERSION`.
- **Ollama:** `KCLI_AI_PROVIDER=ollama`, optional `KCLI_OLLAMA_ENDPOINT`, `KCLI_AI_MODEL`.
- **Custom:** `KCLI_AI_PROVIDER=custom`, `KCLI_AI_ENDPOINT`, optional `KCLI_AI_API_KEY`.

### Quick setup (OpenAI)

```bash
kcli ai config --provider=openai --model=gpt-4o-mini --key=sk-... --enable
kcli ai status
```

### Quick setup (Ollama, local)

```bash
# Ensure Ollama is running (e.g. ollama serve)
kcli ai config --provider=ollama --model=llama3 --enable
kcli ai status
```

## AI commands

| Command | Description |
|---------|-------------|
| `kcli ai explain <resource>` | Explain the resource or concept. |
| `kcli ai why <resource>` | Explain probable cause of current state (e.g. failure). |
| `kcli why <resource>` | Same as `kcli ai why`. |
| `kcli ai suggest-fix <resource>` | Suggest remediation. |
| `kcli suggest fix <resource>` | Same as suggest-fix. |
| `kcli fix <resource>` | Same; use `--dry-run` to only print suggestions. |
| `kcli summarize events [--since=6h]` | Summarize cluster events. |
| `kcli ai <question>` | Natural-language query (e.g. “which pods are crashing?”). |
| `kcli ai config` | Show current AI config. |
| `kcli ai status` | Show enabled state, provider, model, budget, usage. |
| `kcli ai usage` | Monthly usage (calls, tokens, cost). |
| `kcli ai cost` | Cost and budget status. |

Resource can be a short name (e.g. `my-pod`) or qualified (e.g. `pod/my-pod`, `deployment/my-app`). Use `-n <namespace>` if needed.

## Logs + AI

When AI is enabled, logs can be piped to AI for analysis:

- `kcli logs <pod-or-selector> --ai-summarize` — Summarize log content.
- `kcli logs <pod-or-selector> --ai-errors` — Extract and explain errors.
- `kcli logs <pod> --ai-explain` — Explain what the logs indicate.

## Budget and limits

- **Monthly budget:** Set in config as `ai.budgetMonthlyUSD` (default 50). Use `kcli ai config --budget 100`.
- **Soft limit:** `ai.softLimitPercent` (default 80). When usage crosses this percentage of the budget, kcli warns but still allows AI calls.
- **Hard limit:** At 100% of the monthly budget, AI requests are blocked until the next month or you raise the budget.

Usage is stored under the kcli config directory and tracked per month. View with `kcli ai usage` and `kcli ai cost`.

## Caching

AI responses are cached for a period to reduce cost and latency. Repeated identical requests (e.g. same “why” for the same resource) may return cached results. Cache hits are included in `kcli ai usage`.

## Graceful degradation

- If AI is disabled or not configured, AI commands print a short message and exit successfully (no error).
- If the provider is unreachable or returns an error, kcli prints the error and continues; it does not fail the whole command.
- In the TUI, the “A” key and AI tab show a message when AI is disabled or unavailable.

## Privacy and data

- For **OpenAI, Anthropic, Azure:** Resource metadata and event/log content may be sent to the provider’s API. Do not enable on clusters with sensitive data unless your provider and policies allow it.
- **Ollama** runs locally; data stays on your machine.
- **Custom** endpoint: data is sent to whatever URL you configure; ensure you trust that service.
