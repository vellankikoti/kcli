# TUI Guide

`kcli ui` starts the interactive terminal UI (Bubble Tea). Use it to browse resources, view details, tail logs, and run AI analysis without leaving the terminal.

## Launching

```bash
kcli ui
```

The TUI uses your current context and namespace (or `--context` / `-n` if set). Config from `~/.kcli/config.yaml` applies (e.g. `tui.refresh_interval`, `tui.theme`, `tui.animations`).

## Main view: resource list

- **Navigation:** `j` / `down` or `k` / `up` to move; `pgup` / `pgdown` to page.
- **Filter:** Press `/` to type a filter; reduce the list by name, namespace, or status. `Esc` or `Enter` to apply and exit filter mode.
- **Refresh:** `r` to reload the current resource list.
- **Quit:** `q` or `Ctrl+C`.

## Switching resource types

Press `:` then type a resource kind and `Enter`:

| Input | Resource |
|-------|----------|
| `pods`, `po` | Pods |
| `deploy`, `deployments` | Deployments |
| `svc`, `services` | Services |
| `nodes`, `no` | Nodes |
| `events`, `ev` | Events |
| `ns`, `namespaces` | Namespaces |
| `ing`, `ingresses` | Ingresses |
| `cm`, `configmaps` | ConfigMaps |
| `secrets` | Secrets |
| `pvc` | PersistentVolumeClaims |
| `jobs` | Jobs |
| `cronjobs`, `cj` | CronJobs |
| `ep`, `endpoints` | Endpoints |

## Actions on selected row

With one row selected:

| Key | Action |
|-----|--------|
| `Enter` | Open detail view for the resource. |
| `y` | Show YAML (same as `kubectl get ... -o yaml`). |
| `d` | Describe resource. |
| `l` | Logs (pods only; tail). |
| `A` | AI analyze (why) — only when AI is enabled. |
| `S` | Sort by column (cycle). |

In **detail view**, tabs are available:

- **1** or **h** / **Shift+Tab** / **left**: Overview.
- **2** or **l** / **Tab** / **right**: Events.
- **3**: YAML.
- **4**: AI Analysis (when AI enabled).

Press `Esc` to close detail and return to the list.

## Multi-select and bulk actions

- **Space:** Toggle selection of the current row.
- **Ctrl+A:** Select all visible rows.
- **Ctrl+B:** Enter bulk action mode. You can then type:
  - `delete` — delete selected resources (confirmation required; type `yes`).
  - `scale=<n>` — scale selected deployments to `n` replicas (confirmation; type `yes`).

## XRay (resource relationships)

- **Ctrl+X:** Open XRay view for the selected resource. Shows a tree of related resources (e.g. Deployment → ReplicaSets → Pods, Services, Ingress).
- Navigate with `j`/`k`; `Enter` to jump to that resource in the main list. `Esc` to close XRay.

## Other shortcuts

| Key | Action |
|-----|--------|
| **Ctrl+W** | Toggle wide mode (more columns). |
| **Ctrl+S** | Save current table snapshot to a file (prompted for path). |
| **Ctrl+T** | Cycle TUI theme (ocean, forest, amber). |
| **Ctrl+A** (in list) | Toggle AI on/off in TUI (when AI is configured). |

## Sorting

- **S** on the list view cycles sort by the current column (e.g. name, status, restarts).
- In filter mode, you can also sort by column: type a column name (e.g. `NAME`, `STATUS`, `RESTARTS`, `AGE`, `NODE`, `TYPE`) or a value to filter.

## Detail view

After pressing `Enter` on a resource:

- **Overview:** Status, conditions, and key fields.
- **Events:** Recent events for the resource.
- **YAML:** Full resource YAML.
- **AI Analysis:** Output of “why” analysis (if AI enabled).

Use **1**–**4** or **Tab** / **Shift+Tab** to switch tabs; **r** to refresh; **Esc** to go back to the list.

## Screenshots

To capture a screenshot for docs or reports:

1. Run `kcli ui` and navigate to the view you want (e.g. pod list, XRay, detail).
2. Use your terminal or OS screenshot tool (e.g. Cmd+Shift+4 on macOS, or a terminal capture tool).
3. For GIFs, use a tool like LICEcap or terminal recorder (e.g. asciinema) while using the keybindings above.

## Performance

- Refresh interval is controlled by `tui.refresh_interval` in config (default `2s`). Increase it on large clusters to reduce API load.
- Disable animations with `tui.animations: false` in `~/.kcli/config.yaml` if the UI feels slow.
