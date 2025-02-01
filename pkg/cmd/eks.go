package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EKSClientAPI interface to make testing easier
type EKSClientAPI interface {
	// Access Entry methods
	ListAssociatedAccessPolicies(ctx context.Context, params *eks.ListAssociatedAccessPoliciesInput, optFns ...func(*eks.Options)) (*eks.ListAssociatedAccessPoliciesOutput, error)
	ListAccessEntries(ctx context.Context, params *eks.ListAccessEntriesInput, optFns ...func(*eks.Options)) (*eks.ListAccessEntriesOutput, error)
	DescribeAccessEntry(ctx context.Context, params *eks.DescribeAccessEntryInput, optFns ...func(*eks.Options)) (*eks.DescribeAccessEntryOutput, error)

	// Addon methods
	ListAddons(ctx context.Context, params *eks.ListAddonsInput, optFns ...func(*eks.Options)) (*eks.ListAddonsOutput, error)
	DescribeAddon(ctx context.Context, params *eks.DescribeAddonInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonOutput, error)

	// Cluster methods
	DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)

	// Fargate Profile methods
	ListFargateProfiles(ctx context.Context, params *eks.ListFargateProfilesInput, optFns ...func(*eks.Options)) (*eks.ListFargateProfilesOutput, error)
	DescribeFargateProfile(ctx context.Context, params *eks.DescribeFargateProfileInput, optFns ...func(*eks.Options)) (*eks.DescribeFargateProfileOutput, error)

	// Insight methods
	ListInsights(ctx context.Context, params *eks.ListInsightsInput, optFns ...func(*eks.Options)) (*eks.ListInsightsOutput, error)
	DescribeInsight(ctx context.Context, params *eks.DescribeInsightInput, optFns ...func(*eks.Options)) (*eks.DescribeInsightOutput, error)

	// NodeGroup methods
	ListNodegroups(ctx context.Context, params *eks.ListNodegroupsInput, optFns ...func(*eks.Options)) (*eks.ListNodegroupsOutput, error)
	DescribeNodegroup(ctx context.Context, params *eks.DescribeNodegroupInput, optFns ...func(*eks.Options)) (*eks.DescribeNodegroupOutput, error)

	// Pod Identity Association methods
	ListPodIdentityAssociations(ctx context.Context, params *eks.ListPodIdentityAssociationsInput, optFns ...func(*eks.Options)) (*eks.ListPodIdentityAssociationsOutput, error)
	DescribePodIdentityAssociation(ctx context.Context, params *eks.DescribePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.DescribePodIdentityAssociationOutput, error)
}

type EKSClient struct {
	client      EKSClientAPI
	clusterName *string
}

func NewEKSClient(clusterName *string) (*EKSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &EKSClient{
		client:      eks.NewFromConfig(cfg),
		clusterName: clusterName,
	}, nil
}

type ResourceList struct {
	metav1.TypeMeta
	Items struct {
		AccessEntries           []AccessEntry            `json:"access-entries"`
		Addons                  []Addon                  `json:"addons"`
		Nodegroups              []Nodegroup              `json:"nodegroups"`
		FargateProfiles         []FargateProfile         `json:"fargate-profiles"`
		PodIdentityAssociations []PodIdentityAssociation `json:"pod-identity-associations"`
		Insights                []Insight                `json:"insights"`
		Cluster                 []Cluster                `json:"cluster"`
	} `json:"items"`
}

// Implement runtime.Object interface for ResourceList
func (r *ResourceList) GetObjectKind() schema.ObjectKind {
	return &r.TypeMeta
}

func (r *ResourceList) DeepCopyObject() runtime.Object {
	return &ResourceList{
		TypeMeta: r.TypeMeta,
		Items: struct {
			AccessEntries           []AccessEntry            `json:"access-entries"`
			Addons                  []Addon                  `json:"addons"`
			Nodegroups              []Nodegroup              `json:"nodegroups"`
			FargateProfiles         []FargateProfile         `json:"fargate-profiles"`
			PodIdentityAssociations []PodIdentityAssociation `json:"pod-identity-associations"`
			Insights                []Insight                `json:"insights"`
			Cluster                 []Cluster                `json:"cluster"`
		}{
			AccessEntries:           append([]AccessEntry(nil), r.Items.AccessEntries...),
			Addons:                  append([]Addon(nil), r.Items.Addons...),
			Nodegroups:              append([]Nodegroup(nil), r.Items.Nodegroups...),
			FargateProfiles:         append([]FargateProfile(nil), r.Items.FargateProfiles...),
			PodIdentityAssociations: append([]PodIdentityAssociation(nil), r.Items.PodIdentityAssociations...),
			Insights:                append([]Insight(nil), r.Items.Insights...),
			Cluster:                 append([]Cluster(nil), r.Items.Cluster...),
		},
	}
}
