# Testing

How to run kcli’s test suites and add new tests.

## Unit tests

Run all packages:

```bash
go test ./...
```

Run a specific package:

```bash
go test ./internal/cli/...
go test ./internal/config/...
go test ./internal/plugin/...
go test ./internal/state/...
go test ./internal/runner/...
go test ./internal/ai/...
go test ./pkg/api/...
```

Run with coverage:

```bash
go test -cover ./...
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out
```

Tests in `*_test.go` next to the code are unit tests. They should not require a live cluster; use mocks or capture kubectl output where needed.

## Integration-style tests

Some packages may run tests that call `kubectl` or expect a cluster. Those tests are often skipped when no cluster is available (e.g. `if testing.Short() { t.Skip(...) }` or a check for `KUBECONFIG`). Run them with:

```bash
go test -count=1 ./...
# Or without -short if the suite respects it
go test -short ./...
```

Check the package’s `*_test.go` for environment variables or build tags that enable integration tests.

## Scripts

The repo may provide:

- **scripts/alpha-smoke.sh** — Smoke test: core commands, completion generation, and optionally cluster checks. Run after build to validate a release candidate.
- **scripts/perf-check.sh** — Performance gate: startup time, get/ctx latency, memory. Used to enforce TASK-KCLI-018–style targets.
- **scripts/release-gate.sh** — Full gate: tests + alpha smoke + perf. Run before release.

Example:

```bash
./scripts/alpha-smoke.sh
./scripts/perf-check.sh
./scripts/release-gate.sh
```

These scripts assume a built binary (e.g. `./bin/kcli` or `kcli` on PATH) and, for cluster-dependent steps, a valid kubeconfig.

## Writing new tests

### Unit test example

```go
// internal/cli/example_test.go
package cli

import (
	"testing"
)

func TestParseSomething(t *testing.T) {
	out, err := parseSomething("a,b,c")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3, got %d", len(out))
	}
}
```

### Table-driven test

```go
func TestScopeArgsFor(t *testing.T) {
	a := &app{context: "prod", namespace: "default"}
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"adds context and ns", []string{"get", "pods"}, []string{"--context", "prod", "-n", "default", "get", "pods"}},
		{"respects -A", []string{"get", "pods", "-A"}, []string{"--context", "prod", "get", "pods", "-A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.scopeArgsFor(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scopeArgsFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

### Testing commands

Use `NewRootCommandWithIO(stdin, stdout, stderr)` to capture output and inject input:

```go
func TestConfigView(t *testing.T) {
	var out bytes.Buffer
	root := cli.NewRootCommandWithIO(strings.NewReader(""), &out, io.Discard)
	root.SetArgs([]string{"config", "view"})
	err := root.Execute()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "general:") {
		t.Error("expected YAML output with general:")
	}
}
```

### Skipping cluster-dependent tests

```go
func TestGetPodsRealCluster(t *testing.T) {
	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("KUBECONFIG not set, skipping cluster test")
	}
	// ...
}
```

Or use `testing.Short()`:

```go
if testing.Short() {
	t.Skip("skipping in short mode")
}
```

Run with `go test -short ./...` to skip these.

## CI

CI typically runs:

- `go test ./...` (and possibly `go test -race ./...`)
- Linters (e.g. golangci-lint, go vet)
- Optional: alpha smoke and perf checks if a cluster or runner is configured

Check the repo’s workflow files (e.g. `.github/workflows/`) for the exact commands.
