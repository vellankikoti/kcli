# Installation Guide

## Build from source

```bash
cd kcli
go build -o bin/kcli ./cmd/kcli
```

Add to your `PATH`:

```bash
export PATH="$(pwd)/bin:${PATH}"
```

## Verify installation

```bash
kcli version
kcli get pods -A
```

## Shell completion

```bash
kcli completion bash > /etc/bash_completion.d/kcli
kcli completion zsh > "${fpath[1]}/_kcli"
```

## Troubleshooting

- `kubeconfig not found or empty`:
  - Set `KUBECONFIG` or pass `--kubeconfig`.
- `cannot reach Kubernetes API endpoint`:
  - Verify current context, VPN/network, and cluster endpoint.
- Slow command concerns:
  - Run `./scripts/perf-check.sh` to validate local baseline.
