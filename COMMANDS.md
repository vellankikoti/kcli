# kcli Command Reference (Basic)

## Core kubectl parity

Examples:

```bash
kcli get pods -A
kcli describe pod <name> -n <ns>
kcli apply -f manifest.yaml
kcli delete pod <name>
kcli cp <ns>/<pod>:/tmp/x ./x
```

Non-built-in verbs automatically fall back to kubectl while preserving kcli global flags where relevant.

## Workflow

```bash
kcli ctx
kcli ctx <name>
kcli ctx -
kcli ns
kcli ns <name>
kcli ns -
kcli search <pattern>
kcli config view
kcli config get tui.refresh_interval
kcli config set tui.refresh_interval 3s
kcli config reset --yes
kcli config edit
```

## Observability

```bash
kcli health
kcli health pods
kcli health nodes
kcli metrics
kcli restarts
kcli restarts --recent=1h --threshold=5
kcli events --recent 1h
kcli events --type=Warning --all
kcli events --watch
kcli instability
kcli instability pods
kcli incident
kcli incident logs <ns>/<pod> --tail=200
kcli incident describe <ns>/<pod>
kcli incident restart <ns>/<pod>
kcli logs <pod> --ai-summarize
kcli logs <selector> --ai-errors
kcli logs <pod> --ai-explain
```

## TUI

```bash
kcli ui
```

Main keys:

- `/` filter
- `j/k` navigate
- `Enter` detail mode
- `1/2/3` tabs in detail mode
- `Esc` back
- `q` quit

## AI (optional)

```bash
export KCLI_AI_PROVIDER=openai
export KCLI_OPENAI_API_KEY=sk-...
kcli why pod/<name>
kcli summarize events
kcli suggest fix deployment/<name>
kcli fix deployment/<name>
kcli ai "which pods are crashing?"
kcli ai summarize events --since=6h
kcli ai config --provider=openai --model=gpt-4o-mini --enable
kcli ai status
kcli ai usage
kcli ai cost
```

## Plugin management

```bash
kcli plugin list
kcli plugin search <keyword>
kcli plugin marketplace
kcli plugin info <name>
kcli plugin install <marketplace-name>
kcli plugin install ./local-plugin
kcli plugin install github.com/user/plugin
kcli plugin update <name>
kcli plugin update --all
kcli plugin remove <name>
kcli plugin inspect <name>
kcli plugin allow <name>
kcli plugin revoke <name>
kcli plugin run <name> [args...]
kcli <manifest-command-alias> [args...]
```
