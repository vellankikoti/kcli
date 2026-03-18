# Installation

Install kcli on macOS, Linux, or Windows. You need **kubectl** on your `PATH`; kcli uses it to talk to clusters.

## Building from source (all platforms)

Requires **Go 1.21+**.

```bash
git clone https://github.com/kubilitics/kcli.git
cd kcli
go build -o bin/kcli ./cmd/kcli
```

Add `bin/` to your `PATH`, or copy the binary:

```bash
# Linux/macOS
sudo cp bin/kcli /usr/local/bin/

# Or user-local
mkdir -p ~/.local/bin
cp bin/kcli ~/.local/bin/
# Add ~/.local/bin to PATH in ~/.bashrc or ~/.zshrc
```

Verify:

```bash
kcli version
```

## macOS

### Homebrew (when available)

```bash
brew install kubilitics/tap/kcli
```

### Direct download

1. Download the latest `kcli-darwin-amd64` or `kcli-darwin-arm64` from [GitHub Releases](https://github.com/kubilitics/kcli/releases).
2. Rename and make executable:
   ```bash
   mv kcli-darwin-arm64 kcli
   chmod +x kcli
   ```
3. Move to a directory on your `PATH` (e.g. `/usr/local/bin` or `~/.local/bin`).

## Linux

### From release tarball

1. Download the appropriate asset from [GitHub Releases](https://github.com/kubilitics/kcli/releases) (e.g. `kcli-linux-amd64` or `kcli-linux-arm64`).
2. Install:
   ```bash
   chmod +x kcli-linux-amd64
   sudo mv kcli-linux-amd64 /usr/local/bin/kcli
   ```

### Package managers (when available)

- **Debian/Ubuntu (apt):** Add Kubilitics repo and install `kcli` (see project release notes).
- **RHEL/CentOS (yum/dnf):** Add Kubilitics repo and install `kcli`.
- **Snap:** `sudo snap install kcli` (if published).
- **Arch (AUR):** Use the `kcli` AUR package if available.

## Windows

### From release

1. Download `kcli-windows-amd64.exe` from [GitHub Releases](https://github.com/kubilitics/kcli/releases).
2. Rename to `kcli.exe` and place in a directory on your `PATH` (e.g. `C:\Program Files\kcli\` or a folder listed in your user `PATH`).

### Chocolatey (when available)

```powershell
choco install kcli
```

### Scoop (when available)

```powershell
scoop install kcli
```

## Docker

A minimal image can be built from the repo:

```bash
cd kcli
docker build -t kcli:latest .
```

Run with your kubeconfig mounted:

```bash
docker run --rm -v "$HOME/.kube:/root/.kube" -v "$HOME/.kcli:/root/.kcli" kcli:latest get pods -A
```

## Shell completion

After installing, enable completion so that `kcli <TAB>` suggests commands and resources.

**Bash:**

```bash
kcli completion bash > /etc/bash_completion.d/kcli
# Or for current user only:
kcli completion bash > ~/.kcli-completion.bash
echo 'source ~/.kcli-completion.bash' >> ~/.bashrc
```

**Zsh:**

```bash
kcli completion zsh > "${fpath[1]}/_kcli"
# Or if fpath[1] is not set:
mkdir -p ~/.zsh/completions
kcli completion zsh > ~/.zsh/completions/_kcli
echo 'fpath=(~/.zsh/completions $fpath)' >> ~/.zshrc
echo 'autoload -Uz compinit && compinit' >> ~/.zshrc
```

**Fish:**

```bash
kcli completion fish > ~/.config/fish/completions/kcli.fish
```

## Verifying kubectl

kcli delegates cluster calls to `kubectl`. Ensure kubectl is installed and working:

```bash
kubectl version --client
kubectl config get-contexts
```

If `kubectl` is not on your `PATH`, kcli commands that hit the cluster will fail.

## Upgrade

- **From source:** Pull the latest tag or branch and rebuild: `go build -o bin/kcli ./cmd/kcli`.
- **From package/installer:** Use your package manager’s upgrade command or replace the binary with the new release.

## Uninstall

- Remove the `kcli` binary from your `PATH`.
- Optionally remove user data:
  - `~/.kcli/config.yaml`
  - `~/.kcli/state.json`
  - `~/.kcli/plugins/`
- Remove shell completion files you added (e.g. `/etc/bash_completion.d/kcli`, `_kcli` in your zsh path).
