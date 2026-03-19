# kcli Command Reference

kcli wraps kubectl transparently while adding intelligent enhancements: colored tables, status icons, "with" modifiers, natural language aliases, and observability commands.

---

## Quick Start

```bash
kcli status                              # cluster health at a glance
kcli get pods                            # same as kubectl, with crash hints
kcli get pods with ip,node               # extra columns — the kcli superpower
kcli health                              # health score with issues table
kcli find nginx                          # search resources by name pattern
kcli who pod/my-pod                      # trace ownership chain
```

---

## Core kubectl Commands (Full Parity)

All kubectl verbs work through kcli transparently. kcli adds colored output, safety confirmations, and crash hints.

```bash
# Standard CRUD
kcli get pods -A                         # list pods across all namespaces
kcli get pods -n prod -o yaml            # YAML output (unchanged from kubectl)
kcli describe pod my-pod -n prod         # detailed resource view
kcli apply -f manifest.yaml              # apply configuration
kcli create deployment my-app --image=nginx
kcli delete pod my-pod -n staging        # safety confirmation prompt
kcli delete pod my-pod --yes             # skip confirmation for CI/CD

# Workload management
kcli scale deployment/api --replicas=3
kcli rollout restart deployment/api
kcli rollout status deployment/api
kcli edit deployment/api

# Debugging
kcli logs pod/my-pod -f --timestamps
kcli logs deployment/api --tail=100
kcli exec pod/my-pod -- bash
kcli debug pod/my-pod --image=busybox
kcli top pods -A
kcli top nodes

# Node operations
kcli drain node-1 --ignore-daemonsets
kcli cordon node-1
kcli uncordon node-1
kcli taint nodes node-1 key=value:NoSchedule

# Other
kcli cp my-pod:/tmp/file ./local-file
kcli port-forward svc/api 8080:80
kcli explain deployment.spec.strategy
kcli annotate pod/my-pod note="updated"
kcli label pod/my-pod env=prod
kcli expose deployment/api --port=80
kcli patch deployment/api -p '{"spec":{"replicas":5}}'
kcli wait --for=condition=ready pod/my-pod
```

---

## "with" Modifiers — The kcli Superpower

Add extra columns to any `get` command. No more `-o wide` or custom-columns.

### Pods
```bash
kcli get pods with ip                    # pod IP addresses
kcli get pods with node                  # which node each pod runs on
kcli get pods with ip,node               # combine modifiers
kcli get pods with containers            # container names
kcli get pods with requests              # cpu/memory requests
kcli get pods with limits                # cpu/memory limits
kcli get pods with containers,requests,limits  # full resource view
kcli get pods with labels                # pod labels
kcli get pods with qos                   # QoS class (Guaranteed/Burstable/BestEffort)
kcli get pods with images                # container images
kcli get pods with ports                 # exposed ports
kcli get pods with sa                    # service account
kcli get pods with all                   # every available column
kcli get pods -A with ip,node            # all namespaces + extra columns
```

### Deployments
```bash
kcli get deployments with replicas       # desired/current replica counts
kcli get deployments with images         # container images
kcli get deployments with strategy       # RollingUpdate/Recreate
kcli get deployments with labels         # deployment labels
kcli get deployments with selectors      # pod selectors
kcli get deployments with conditions     # deployment conditions
kcli get deployments with all            # everything
```

### Services
```bash
kcli get services with endpoints         # ready/total endpoint counts
kcli get services with ports             # detailed port mappings
kcli get services with selectors         # pod selectors
kcli get services with all
```

### Nodes
```bash
kcli get nodes with capacity             # CPU and memory capacity
kcli get nodes with zone                 # availability zone
kcli get nodes with taints               # node taints
kcli get nodes with pods                 # allocatable pod count
kcli get nodes with conditions           # node conditions
kcli get nodes with labels               # all labels
kcli get nodes with all
```

### StatefulSets
```bash
kcli get statefulsets with replicas      # desired/current
kcli get statefulsets with images
kcli get statefulsets with labels
kcli get statefulsets with conditions
kcli get statefulsets with all
```

### DaemonSets
```bash
kcli get daemonsets with images
kcli get daemonsets with labels
kcli get daemonsets with node-selector   # node selector labels
kcli get daemonsets with conditions
kcli get daemonsets with all
```

### Other Resources
```bash
kcli get pvc with all                    # PersistentVolumeClaims
kcli get pv with all                     # PersistentVolumes (cluster-scoped)
kcli get ingresses with all              # Ingresses
kcli get events with all                 # Events
```

---

## Natural Language Aliases

Speak to your cluster, don't memorize syntax.

```bash
# Health & Status
kcli status                              # cluster health table (nodes, pods, deployments)
kcli status pod/my-pod                   # detailed pod status with containers
kcli status deployment/api               # deployment status with conditions
kcli status service/api                  # service details with ports and endpoints
kcli health                              # health score (0-100) with issues
kcli health pods                         # pod health breakdown
kcli health nodes                        # node health breakdown
kcli health -o json                      # machine-readable output

# Discovery
kcli find nginx                          # search all resources matching "nginx"
kcli find pod payment                    # search only pods
kcli find svc api -A                     # search services across all namespaces
kcli show pods                           # enhanced colored pod listing
kcli show deployments                    # deployments with ready/status
kcli show nodes                          # nodes with status/roles

# Observability
kcli age pods                            # pods sorted by creation time (newest first)
kcli age pods --oldest                   # oldest first
kcli age deployments -A                  # deployments across all namespaces
kcli count pods                          # pod counts by status
kcli count deployments -A                # deployment counts across namespaces
kcli count all                           # count all resource types
kcli restarts                            # pods sorted by restart count
kcli restarts --threshold=5              # only pods with 5+ restarts
kcli restarts --recent=1h                # restarts in the last hour
kcli events                              # recent cluster events
kcli events --type=Warning               # only warning events
kcli events --recent=30m                 # events from last 30 minutes
kcli events --watch                      # live event stream
kcli events --resource=pod/nginx         # events for specific resource
kcli metrics                             # node and pod resource usage
kcli instability                         # restart leaders + warning events

# Tracing & Attribution
kcli who pod/my-pod                      # ownership chain + related resources
kcli who deployment/api                  # deployment → pods + services
kcli who service/api                     # service → selected pods
kcli where pod/my-pod                    # node, zone, region, IP
kcli where pods -n prod                  # all pods with their locations
kcli where deployment/api                # zone distribution
kcli blame deployment/api                # who changed it, when, from what system

# Incident Mode
kcli incident                            # CrashLoop/OOM/restarts/pressure/events
```

---

## Context & Namespace Management

```bash
kcli ctx                                 # list contexts (* marks active)
kcli ctx production                      # switch context
kcli ctx -                               # switch to previous context
kcli ns                                  # list namespaces (* marks active)
kcli ns prod                             # switch namespace
kcli ns -                                # switch to previous namespace
```

---

## Configuration

```bash
kcli config view                         # show effective config
kcli config get tui.refresh_interval     # get a specific value
kcli config set tui.refresh_interval 3s  # set a value
kcli config reset --yes                  # reset to defaults
kcli config edit                         # open in editor
kcli config profile list                 # manage config profiles
```

---

## Diagnostics

```bash
kcli doctor                              # validate environment
kcli audit                               # view command audit trail
kcli audit enable                        # enable audit logging
kcli audit disable                       # disable audit logging
kcli rbac                                # RBAC privilege analysis
kcli rbac who-can get pods               # who can perform an action
```

---

## Shell Completion

```bash
kcli completion bash                     # generate bash completions
kcli completion zsh                      # generate zsh completions
kcli completion fish                     # generate fish completions
kcli completion powershell               # generate powershell completions

# Install (zsh example):
echo 'source <(kcli completion zsh)' >> ~/.zshrc
```

---

## TUI Mode

```bash
kcli ui                                  # launch interactive terminal UI
```

Keys: `/` filter, `j/k` navigate, `Enter` detail, `1/2/3` tabs, `Esc` back, `q` quit

---

## Plugin System

```bash
kcli plugin list                         # installed plugins
kcli plugin install <name>               # install from marketplace
kcli plugin install ./local-plugin       # install local plugin
kcli plugin remove <name>                # uninstall
kcli plugin run <name> [args...]         # run a plugin
```

---

## Global Flags

```bash
--context <name>                         # override kubectl context
-n, --namespace <name>                   # override namespace
--kubeconfig <path>                      # custom kubeconfig path
--yes                                    # skip safety confirmations (CI/CD)
-o, --output <format>                    # output format (table/json/yaml)
```

---

## Environment Variables

```bash
KCLI_CI=true                             # CI mode: auto --yes, no animations
KCLI_HINTS=0                             # disable crash hints
KCLI_NO_AUDIT=1                          # disable audit trail
KCLI_DEBUG_STARTUP=1                     # show startup timing
KCLI_KUBECTL=/path/to/kubectl            # custom kubectl binary
NO_COLOR=1                               # disable all colors
```
