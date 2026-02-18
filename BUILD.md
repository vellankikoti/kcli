# kcli Build Guide

This guide covers building kcli binary for different platforms and deployment scenarios.

## Prerequisites

- Go 1.24 or later
- kubectl (for testing, not required for build)

## Basic Build

Build kcli binary for current platform:

```bash
cd kcli
go build -ldflags="-s -w" -o bin/kcli ./cmd/kcli
```

The binary will be created at `kcli/bin/kcli`.

## Cross-Platform Builds

### Linux (amd64)
```bash
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/kcli-linux-amd64 ./cmd/kcli
```

### Linux (arm64)
```bash
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/kcli-linux-arm64 ./cmd/kcli
```

### macOS (amd64)
```bash
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/kcli-darwin-amd64 ./cmd/kcli
```

### macOS (arm64)
```bash
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/kcli-darwin-arm64 ./cmd/kcli
```

### Windows (amd64)
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/kcli-windows-amd64.exe ./cmd/kcli
```

## Building for Desktop App

Use the provided script:

```bash
./scripts/build-kcli-for-desktop.sh
```

This script:
1. Builds kcli for current platform
2. Copies binary to `kubilitics-desktop/binaries/` with target triple suffix
3. Creates executable permissions

**Target Triple Suffixes:**
- macOS: `kcli-x86_64-apple-darwin`, `kcli-aarch64-apple-darwin`
- Linux: `kcli-x86_64-unknown-linux-gnu`, `kcli-aarch64-unknown-linux-gnu`
- Windows: `kcli-x86_64-pc-windows-msvc.exe`

**Manual Build:**
```bash
cd kcli
# Build for macOS universal (both architectures)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../kubilitics-desktop/binaries/kcli-x86_64-apple-darwin ./cmd/kcli
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../kubilitics-desktop/binaries/kcli-aarch64-apple-darwin ./cmd/kcli

# Build for Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../kubilitics-desktop/binaries/kcli-x86_64-unknown-linux-gnu ./cmd/kcli
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ../kubilitics-desktop/binaries/kcli-aarch64-unknown-linux-gnu ./cmd/kcli

# Build for Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../kubilitics-desktop/binaries/kcli-x86_64-pc-windows-msvc.exe ./cmd/kcli

# Make executable (Unix)
chmod +x ../kubilitics-desktop/binaries/kcli-*
```

## Building for Docker

kcli is built during Docker image build (multi-stage build).

**Manual Build (for testing):**
```bash
cd kcli
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../kubilitics-backend/bin/kcli ./cmd/kcli
```

**Docker Build:**
```bash
docker build -f kubilitics-backend/Dockerfile -t kubilitics-backend .
```

The Dockerfile includes:
- Multi-stage build for kcli
- Installation to `/usr/local/bin/kcli`
- `KCLI_BIN` environment variable

## Build Flags

### -ldflags="-s -w"
- `-s`: Omit symbol table and debug information
- `-w`: Omit DWARF symbol table

These flags reduce binary size significantly.

### Version Information
To include version information:
```bash
go build -ldflags="-s -w -X main.version=$(git describe --tags)" -o bin/kcli ./cmd/kcli
```

## Verification

After building, verify the binary:

```bash
# Test version command
./bin/kcli version

# Test help command
./bin/kcli --help

# Test basic command
./bin/kcli get pods --help
```

## CI/CD Builds

### GitHub Actions
kcli binaries are built automatically in:
- `.github/workflows/desktop-ci.yml` - For desktop app
- `.github/workflows/release.yml` - For releases

### Release Builds
All platform binaries are built and uploaded as release artifacts:
- Linux amd64/arm64
- macOS amd64/arm64
- Windows amd64

## Troubleshooting

### Build Fails
- Check Go version: `go version` (requires 1.24+)
- Verify module dependencies: `go mod download`
- Check for syntax errors: `go build ./...`

### Binary Too Large
- Use `-ldflags="-s -w"` to strip debug info
- Consider UPX compression (not recommended for production)

### Cross-Compilation Issues
- Ensure CGO is disabled: `CGO_ENABLED=0`
- Check target platform support
- Verify Go cross-compilation support

## Development Build

For development (with debug info):

```bash
go build -o bin/kcli ./cmd/kcli
```

This includes:
- Symbol table
- Debug information
- Larger binary size

## Production Build

For production:

```bash
go build -ldflags="-s -w" -o bin/kcli ./cmd/kcli
```

This produces:
- Smaller binary
- No debug symbols
- Optimized for distribution
