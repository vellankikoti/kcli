# Plugin Guide

kcli supports plugins for extending functionality. Plugins are executables in `~/.kcli/plugins/` with a `kcli-` prefix and an optional `plugin.yaml` manifest.

## Discovering plugins

```bash
kcli plugin list
kcli plugin search <keyword>
kcli plugin marketplace
```

- **list** — Shows installed plugins (name, version, status).
- **search** — Searches installed plugins and the marketplace catalog.
- **marketplace** — Lists plugins in the catalog (official/community, version, downloads, rating).

## Installing plugins

Install from a local path or from a GitHub repo:

```bash
kcli plugin install ./my-plugin
kcli plugin install github.com/org/kcli-mytool
kcli plugin install cert-manager
```

If the source is a directory, kcli looks for an executable (e.g. `kcli-mytool`) or a Go module to build. For `github.com/org/repo`, kcli clones and builds. For catalog names like `cert-manager`, it uses the marketplace source.

After install, the plugin appears in `kcli plugin list` and can be run as `kcli mytool` or `kcli plugin run mytool`.

## Running plugins

- **By name:** `kcli plugin run <name> [args...]`
- **As first-class command:** If the first argument is not a builtin command, kcli treats it as a plugin: `kcli mytool --flag value`

Example:

```bash
kcli mytool --flag value
kcli plugin run mytool --flag value
```

## Plugin manifest (plugin.yaml)

Plugins in `~/.kcli/plugins/` can include a `plugin.yaml` in the same directory:

```yaml
name: mytool
version: 1.0.0
description: Short description
author: Your Name
permissions:
  - read:pods
  - write:deployments
commands:
  - mt
```

- **name** — Plugin name (must match the executable name without the `kcli-` prefix).
- **version** — Semantic version string.
- **description** / **author** — Shown in `kcli plugin inspect` and search.
- **permissions** — List of permission identifiers the plugin may need. Users must approve them (see below).
- **commands** — Alternative command names (e.g. `mt`) so you can run `kcli mt`.

## Permissions and sandboxing

- Plugins declare **permissions** in the manifest (e.g. `read:pods`, `write:deployments`).
- By default, only plugins under `~/.kcli/plugins/` are executable. Plugins found only on `PATH` are **blocked** unless `KCLI_PLUGIN_ALLOW_PATH=1` is set.
- Before a plugin uses a permission, it must be **approved**:
  - `kcli plugin allow <name>` — Approve all declared permissions.
  - `kcli plugin allow <name> read:pods write:deployments` — Approve specific permissions.
- Revoke with `kcli plugin revoke <name>` or `kcli plugin revoke <name> read:pods`.

Inspect manifest and permission status:

```bash
kcli plugin inspect mytool
kcli plugin info mytool
```

## Updating and removing

```bash
kcli plugin update mytool
kcli plugin update --all
kcli plugin update-all
kcli plugin remove mytool
```

- **update &lt;name&gt;** — Update one plugin from its recorded source.
- **update --all** / **update-all** — Update all installed plugins.

## Official plugins

This repository may ship official plugins (e.g. cert-manager). To install and test them:

```bash
./scripts/install-official-plugins.sh
./scripts/test-official-plugins.sh
```

Paths and names may vary; check the repo’s `scripts/` and docs.

## Writing a plugin

See [Plugin Development](../developer/plugin-development.md) for:

- Executable naming (`kcli-<name>`)
- Manifest format and permissions
- Building and packaging
- Testing with `kcli plugin run`

## Troubleshooting

- **“plugin not found”** — Ensure the binary is in `~/.kcli/plugins/` and named `kcli-<name>`, or on `PATH` with `KCLI_PLUGIN_ALLOW_PATH=1`.
- **“not executable”** — `chmod +x ~/.kcli/plugins/kcli-<name>`.
- **“invalid manifest”** — Fix `plugin.yaml` (valid YAML, required fields). Use `kcli plugin inspect <name>` to see errors.
- **Permission denied** — Run `kcli plugin allow <name>` for the required permissions.
