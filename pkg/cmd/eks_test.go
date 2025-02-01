package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

// mockEKSClient implements the necessary methods for testing
type mockEKSClient struct {
	// Access Entry methods
	listAssociatedAccessPoliciesFunc func(ctx context.Context, params *eks.ListAssociatedAccessPoliciesInput) (*eks.ListAssociatedAccessPoliciesOutput, error)
	listAccessEntriesFunc            func(ctx context.Context, params *eks.ListAccessEntriesInput) (*eks.ListAccessEntriesOutput, error)
	describeAccessEntryFunc          func(ctx context.Context, params *eks.DescribeAccessEntryInput) (*eks.DescribeAccessEntryOutput, error)

	// Addon methods
	listAddonsFunc    func(ctx context.Context, params *eks.ListAddonsInput) (*eks.ListAddonsOutput, error)
	describeAddonFunc func(ctx context.Context, params *eks.DescribeAddonInput) (*eks.DescribeAddonOutput, error)

	// Cluster methods
	describeClusterFunc func(ctx context.Context, params *eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error)

	// Fargate Profile methods
	listFargateProfilesFunc    func(ctx context.Context, params *eks.ListFargateProfilesInput) (*eks.ListFargateProfilesOutput, error)
	describeFargateProfileFunc func(ctx context.Context, params *eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error)

	// Insight methods
	listInsightsFunc    func(ctx context.Context, params *eks.ListInsightsInput) (*eks.ListInsightsOutput, error)
	describeInsightFunc func(ctx context.Context, params *eks.DescribeInsightInput) (*eks.DescribeInsightOutput, error)

	// NodeGroup methods
	listNodegroupsFunc    func(ctx context.Context, params *eks.ListNodegroupsInput) (*eks.ListNodegroupsOutput, error)
	describeNodegroupFunc func(ctx context.Context, params *eks.DescribeNodegroupInput) (*eks.DescribeNodegroupOutput, error)

	// Pod Identity Association methods
	listPodIdentityAssociationsFunc    func(ctx context.Context, params *eks.ListPodIdentityAssociationsInput) (*eks.ListPodIdentityAssociationsOutput, error)
	describePodIdentityAssociationFunc func(ctx context.Context, params *eks.DescribePodIdentityAssociationInput) (*eks.DescribePodIdentityAssociationOutput, error)
}

func (m *mockEKSClient) ListAssociatedAccessPolicies(ctx context.Context, params *eks.ListAssociatedAccessPoliciesInput, optFns ...func(*eks.Options)) (*eks.ListAssociatedAccessPoliciesOutput, error) {
	return m.listAssociatedAccessPoliciesFunc(ctx, params)
}

func (m *mockEKSClient) ListAccessEntries(ctx context.Context, params *eks.ListAccessEntriesInput, optFns ...func(*eks.Options)) (*eks.ListAccessEntriesOutput, error) {
	return m.listAccessEntriesFunc(ctx, params)
}

func (m *mockEKSClient) DescribeAccessEntry(ctx context.Context, params *eks.DescribeAccessEntryInput, optFns ...func(*eks.Options)) (*eks.DescribeAccessEntryOutput, error) {
	return m.describeAccessEntryFunc(ctx, params)
}

func (m *mockEKSClient) ListAddons(ctx context.Context, params *eks.ListAddonsInput, optFns ...func(*eks.Options)) (*eks.ListAddonsOutput, error) {
	return m.listAddonsFunc(ctx, params)
}

func (m *mockEKSClient) DescribeAddon(ctx context.Context, params *eks.DescribeAddonInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonOutput, error) {
	return m.describeAddonFunc(ctx, params)
}

func (m *mockEKSClient) DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error) {
	return m.describeClusterFunc(ctx, params)
}

func (m *mockEKSClient) ListFargateProfiles(ctx context.Context, params *eks.ListFargateProfilesInput, optFns ...func(*eks.Options)) (*eks.ListFargateProfilesOutput, error) {
	return m.listFargateProfilesFunc(ctx, params)
}

func (m *mockEKSClient) DescribeFargateProfile(ctx context.Context, params *eks.DescribeFargateProfileInput, optFns ...func(*eks.Options)) (*eks.DescribeFargateProfileOutput, error) {
	return m.describeFargateProfileFunc(ctx, params)
}

func (m *mockEKSClient) ListInsights(ctx context.Context, params *eks.ListInsightsInput, optFns ...func(*eks.Options)) (*eks.ListInsightsOutput, error) {
	return m.listInsightsFunc(ctx, params)
}

func (m *mockEKSClient) DescribeInsight(ctx context.Context, params *eks.DescribeInsightInput, optFns ...func(*eks.Options)) (*eks.DescribeInsightOutput, error) {
	return m.describeInsightFunc(ctx, params)
}

func (m *mockEKSClient) ListNodegroups(ctx context.Context, params *eks.ListNodegroupsInput, optFns ...func(*eks.Options)) (*eks.ListNodegroupsOutput, error) {
	return m.listNodegroupsFunc(ctx, params)
}

func (m *mockEKSClient) DescribeNodegroup(ctx context.Context, params *eks.DescribeNodegroupInput, optFns ...func(*eks.Options)) (*eks.DescribeNodegroupOutput, error) {
	return m.describeNodegroupFunc(ctx, params)
}

func (m *mockEKSClient) ListPodIdentityAssociations(ctx context.Context, params *eks.ListPodIdentityAssociationsInput, optFns ...func(*eks.Options)) (*eks.ListPodIdentityAssociationsOutput, error) {
	return m.listPodIdentityAssociationsFunc(ctx, params)
}

func (m *mockEKSClient) DescribePodIdentityAssociation(ctx context.Context, params *eks.DescribePodIdentityAssociationInput, optFns ...func(*eks.Options)) (*eks.DescribePodIdentityAssociationOutput, error) {
	return m.describePodIdentityAssociationFunc(ctx, params)
}
