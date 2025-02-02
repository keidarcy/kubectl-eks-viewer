package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Options struct {
	configFlags *genericclioptions.ConfigFlags
	printFlags  *genericclioptions.PrintFlags
	eksClient   *EKSClient
	rawConfig   api.Config

	genericclioptions.IOStreams
	resourceType string
}

func NewOptions(streams genericclioptions.IOStreams) *Options {
	return &Options{
		configFlags: genericclioptions.NewConfigFlags(true),
		printFlags:  genericclioptions.NewPrintFlags(""),
		IOStreams:   streams,
	}
}

var validResourceTypes = []string{
	"access-entries",
	"addons",
	"cluster",
	"fargate-profiles",
	"insights",
	"nodegroups",
	"pod-identity-associations",
}

func NewCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewOptions(streams)

	cmd := &cobra.Command{
		Use:   "kubectl eks-viewer [resource-type]",
		Short: "View EKS cluster resources",
		Long: `View EKS cluster resources.
Without arguments, shows all resource types.
Optionally specify a resource type to show only that type.

Valid resource types:
  - access-entries
  - addons
  - cluster
  - fargate-profiles
  - insights
  - nodegroups
  - pod-identity-associations`,
		Example: `  # List all EKS resources 
  kubectl eks-viewer
  kubectl eks-viewer -o json

  # List specific resources
  kubectl eks-viewer addons
  kubectl eks-viewer -o json nodegroups
  kubectl eks-viewer nodegroups --output=jsonpath='{.items.nodegroups[*].NodegroupName}
  
  # Use with a specific context
  kubectl eks-viewer --context=my-context`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				o.resourceType = args[0]
			}

			if err := o.Validate(); err != nil {
				return err
			}

			if err := o.Complete(); err != nil {
				return err
			}

			if err := o.Run(); err != nil {
				return err
			}
			return nil
		},
	}

	o.configFlags.AddFlags(cmd.Flags())
	o.printFlags.AddFlags(cmd)

	return cmd
}

func (o *Options) Complete() error {
	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return fmt.Errorf("failed to load kubeconfig: %v", err)
	}

	// Get context from flag if specified, otherwise use current-context
	currentContext := ""
	if cf := o.configFlags.Context; cf != nil && *cf != "" {
		currentContext = *cf
	} else {
		currentContext = o.rawConfig.CurrentContext
	}

	if currentContext == "" {
		return fmt.Errorf("no context specified and no current-context found in kubeconfig")
	}

	context, exists := o.rawConfig.Contexts[currentContext]
	if !exists {
		return fmt.Errorf("context %q not found in kubeconfig", currentContext)
	}

	// EKS cluster ARN format: arn:aws:eks:<region>:<account>:cluster/<cluster-name>
	// Or cluster name format: <cluster-name>.<region>.eksctl.io
	clusterName := context.Cluster
	if strings.Contains(clusterName, "/") {
		// Extract from ARN
		parts := strings.Split(clusterName, "/")
		clusterName = parts[len(parts)-1]
	} else if strings.Contains(clusterName, ".eksctl.io") {
		// Extract from eksctl format
		clusterName = strings.Split(clusterName, ".")[0]
	}

	o.eksClient, err = NewEKSClient(&clusterName)
	if err != nil {
		return fmt.Errorf("failed to create AWS client: %v", err)
	}

	return nil
}

func (o *Options) Validate() error {
	if o.resourceType == "" {
		return nil
	}

	for _, validType := range validResourceTypes {
		if o.resourceType == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid resource type %q. Valid types are: %s",
		o.resourceType, strings.Join(validResourceTypes, ", "))
}

type resourceFetcher struct {
	resourceType string
	fetch        func(context.Context) error
	printer      func(interface{}, io.Writer) error
}

func (o *Options) fetchResource(ctx context.Context, f resourceFetcher) error {
	name := strings.ReplaceAll(f.resourceType, "-", " ")
	fmt.Printf("\r \033[36mFetching %s...\033[m", name)
	if err := f.fetch(ctx); err != nil {
		return fmt.Errorf("failed to list %s: %v", name, err)
	}
	fmt.Printf("\r%s\r", strings.Repeat(" ", 50)) // Clear the line
	return nil
}

func (o *Options) Run() error {
	ctx := context.Background()
	resourceList := &ResourceList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "EksResourceList",
		},
	}

	printer, err := o.printFlags.ToPrinter()
	if err != nil {
		return err
	}
	isTableFormat := *o.printFlags.OutputFormat == "" || *o.printFlags.OutputFormat == "wide"

	// Define all available resources
	allResources := []resourceFetcher{
		{
			resourceType: "cluster",
			fetch: func(ctx context.Context) error {
				cluster, err := o.eksClient.DescribeCluster(ctx)
				resourceList.Items.Cluster = cluster
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewClusterPrinter().PrintObj(&ClusterList{Items: resourceList.Items.Cluster}, w)
			},
		},
		{
			resourceType: "access-entries",
			fetch: func(ctx context.Context) error {
				entries, err := o.eksClient.ListAccessEntries(ctx)
				resourceList.Items.AccessEntries = entries
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewAccessEntryPrinter(o.eksClient).PrintObj(&AccessEntryList{Items: resourceList.Items.AccessEntries}, w)
			},
		},
		{
			resourceType: "addons",
			fetch: func(ctx context.Context) error {
				addons, err := o.eksClient.ListAddons(ctx)
				resourceList.Items.Addons = addons
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewAddonPrinter().PrintObj(&AddonList{Items: resourceList.Items.Addons}, w)
			},
		},
		{
			resourceType: "nodegroups",
			fetch: func(ctx context.Context) error {
				nodeGroups, err := o.eksClient.ListNodeGroups(ctx)
				resourceList.Items.Nodegroups = nodeGroups
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewNodeGroupPrinter().PrintObj(&NodeGroupList{Items: resourceList.Items.Nodegroups}, w)
			},
		},
		{
			resourceType: "fargate-profiles",
			fetch: func(ctx context.Context) error {
				fargateProfiles, err := o.eksClient.ListFargateProfiles(ctx)
				resourceList.Items.FargateProfiles = fargateProfiles
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewFargateProfilePrinter().PrintObj(&FargateProfileList{Items: resourceList.Items.FargateProfiles}, w)
			},
		},
		{
			resourceType: "pod-identity-associations",
			fetch: func(ctx context.Context) error {
				podIdentityAssociations, err := o.eksClient.ListPodIdentityAssociations(ctx)
				resourceList.Items.PodIdentityAssociations = podIdentityAssociations
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewPodIdentityAssociationPrinter().PrintObj(&PodIdentityAssociationList{Items: resourceList.Items.PodIdentityAssociations}, w)
			},
		},
		{
			resourceType: "insights",
			fetch: func(ctx context.Context) error {
				insights, err := o.eksClient.ListInsights(ctx)
				resourceList.Items.Insights = insights
				return err
			},
			printer: func(obj interface{}, w io.Writer) error {
				return NewInsightPrinter().PrintObj(&InsightList{Items: resourceList.Items.Insights}, w)
			},
		},
	}

	// Select resources to process
	resourcesToFetch := allResources
	if o.resourceType != "" {
		found := false
		for _, res := range allResources {
			if res.resourceType == o.resourceType {
				resourcesToFetch = []resourceFetcher{res}
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("resource type %q not supported. Valid types are: %s",
				o.resourceType, strings.Join(validResourceTypes, ", "))
		}
	}

	if isTableFormat {
		// Fetch and print each resource type individually
		for i, res := range resourcesToFetch {
			if err := o.fetchResource(ctx, res); err != nil {
				return err
			}
			if err := res.printer(nil, o.Out); err != nil {
				return err
			}
			// Add newline between resource types, but not after the last one
			if i < len(resourcesToFetch)-1 {
				fmt.Fprintln(o.Out)
			}
		}
		return nil
	}

	// For non-table formats, fetch all requested resources
	progressMsg := "Fetching EKS resources..."
	if o.resourceType != "" {
		progressMsg = fmt.Sprintf("Fetching %s...", o.resourceType)
	}
	fmt.Printf("\r \033[36m%s\033[m", progressMsg)

	for _, res := range resourcesToFetch {
		if err := res.fetch(ctx); err != nil {
			return err
		}
	}

	// Clear the progress line
	fmt.Printf("\r%s\r", strings.Repeat(" ", 50))
	return printer.PrintObj(resourceList, o.Out)
}
