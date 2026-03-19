package cli

import (
	"strings"

	"github.com/spf13/cobra"

	kubectlpkg "github.com/kubilitics/kcli/internal/kubectl"
)

// newShowCmd creates the `kcli show` command — a natural-language alias for `get`
// that routes through the enhanced get engine for richer table output.
func newShowCmd(a *app) *cobra.Command {
	var sortBy string
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "show <resource> [flags]",
		Short: "Natural language alias for 'get' with enhanced output",
		Long: `Show resources with natural language syntax and enhanced colored table output.

Uses client-go for rich, colored tables with status icons, responsive
column hiding, and optional "with" modifiers for extra detail.

Supported resources:
  pods, deployments, services, nodes, statefulsets, daemonsets,
  pv, pvc, ingresses, events

Examples:
  kcli show pods                        # Pods in current namespace
  kcli show pods -A                     # Pods across all namespaces
  kcli show deployments                 # Deployments with ready/status
  kcli show svc                         # Services with type/IP/ports
  kcli show nodes                       # Nodes with status/roles/version
  kcli show pods --sort age             # Sort by age
  kcli show pods -o json                # JSON output
  kcli show pvc                         # PVCs with capacity/status
  kcli show ingresses                   # Ingresses with hosts/address

Tip: Use "kcli get <resource> with <modifiers>" for extra columns:
  kcli get pods with ip,node            # Add IP and node columns
  kcli get pods with containers,requests,limits  # Resource requests/limits
  kcli get pods with all                # Show all available columns
  kcli get deployments with replicas,images
  kcli get nodes with capacity,zone`,
		GroupID: "core",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			resourceType := mapResourceName(strings.ToLower(args[0]))

			namespace := a.namespace
			allNS, _ := cmd.Flags().GetBool("all-namespaces")

			kubeconfigPath := a.kubeconfig
			if kubeconfigPath == "" {
				kubeconfigPath = kubectlpkg.DefaultKubeconfigPath()
			}

			return kubectlpkg.EnhancedGet(kubeconfigPath, a.context, namespace, resourceType, []string{}, allNS, sortBy, outputFormat)
		},
	}

	cmd.Flags().BoolP("all-namespaces", "A", false, "Show resources from all namespaces")
	cmd.Flags().StringVar(&sortBy, "sort", "", "Sort column (e.g., 'age', 'name')")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json")

	return cmd
}
