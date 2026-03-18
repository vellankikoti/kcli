# Building from Source

Build kcli from source on any platform supported by Go.

## Requirements

- **Go 1.21 or later** — [Install Go](https://go.dev/doc/install).
- **kubectl** — Required at runtime; not needed to build. Install separately for use with kcli.

## Clone and build

```bash
git clone https://github.com/kubilitics/kcli.git
cd kcli
go mod download
go build -o bin/kcli ./cmd/kcli
```

The binary is `bin/kcli`. Run it:

```bash
./bin/kcli version
./bin/kcli --help
```

## Install locally

Copy the binary to a directory on your `PATH`:

```bash
# Linux / macOS
sudo cp bin/kcli /usr/local/bin/

# Or user-local
mkdir -p ~/.local/bin
cp bin/kcli ~/.local/bin/
# Add ~/.local/bin to PATH in ~/.bashrc or ~/.zshrc
```

## Cross-compilation

Build for another OS or architecture:

```bash
# Linux AMD64 from macOS
GOOS=linux GOARCH=amd64 go build -o bin/kcli-linux-amd64 ./cmd/kcli

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o bin/kcli-windows-amd64.exe ./cmd/kcli

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o bin/kcli-darwin-arm64 ./cmd/kcli
```

Use the appropriate binary name for your distribution (e.g. `kcli-darwin-arm64` for Apple Silicon).

## Build tags and CGO

The default build uses no build tags and CGO is typically disabled, so the binary is static and portable. If the project introduces optional features behind build tags (e.g. `plugins` or `k8s_native`), they will be documented in the repo; for a standard build you do not need to pass any tags.

## Version and commit

Version, commit hash, and build date are set at build time via `internal/version`. They may be injected by the project’s Makefile or CI (e.g. `-ldflags "-X ...version.Version=..." -X ...version.Commit=..." -X ...version.BuildDate=..."`). A plain `go build` uses defaults (e.g. "dev", "none", build time). Check with:

```bash
./bin/kcli version
```

## Dependencies

All dependencies are in `go.mod`. After clone:

```bash
go mod download
go mod verify
```

To update dependencies:

```bash
go get -u ./...
go mod tidy
```

## Makefile (if present)

If the repo includes a Makefile:

```bash
make build    # build binary
make test     # run tests
make install  # install to $GOPATH/bin or PREFIX
```

Use `make help` or read the Makefile for targets.

## Docker (if Dockerfile present)

To build a minimal image:

```bash
docker build -t kcli:latest .
docker run --rm -v "$HOME/.kube:/root/.kube" -v "$HOME/.kcli:/root/.kcli" kcli:latest get pods -A
```

The Dockerfile, if any, will be in the repo root or under `scripts/`; ensure kubectl is available inside the image if you need cluster access from the container.
