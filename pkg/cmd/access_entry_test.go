package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewAccessEntryPrinter(t *testing.T) {
	tests := []struct {
		name           string
		accessEntries  []AccessEntry
		mockPolicies   map[string][]string // principalARN -> policyARNs
		expectedOutput []string
	}{
		{
			name: "single access entry with single policy",
			accessEntries: []AccessEntry{
				{
					PrincipalArn:     stringPtr("arn:aws:iam::123456789012:role/test-role"),
					KubernetesGroups: []string{"system:masters"},
				},
			},
			mockPolicies: map[string][]string{
				"arn:aws:iam::123456789012:role/test-role": {"arn:aws:eks::123456789012:policy/test-policy"},
			},
			expectedOutput: []string{
				"ACCESS ENTRY PRINCIPAL ARN",
				"KUBERNETES GROUPS",
				"ACCESS POLICIES",
				"arn:aws:iam::123456789012:role/test-role",
				"system:masters",
				"arn:aws:eks::123456789012:policy/test-policy",
			},
		},
		{
			name: "multiple access entries with multiple policies",
			accessEntries: []AccessEntry{
				{
					PrincipalArn:     stringPtr("arn:aws:iam::123456789012:role/role1"),
					KubernetesGroups: []string{"system:masters", "group1"},
				},
				{
					PrincipalArn:     stringPtr("arn:aws:iam::123456789012:role/role2"),
					KubernetesGroups: []string{"system:authenticated"},
				},
			},
			mockPolicies: map[string][]string{
				"arn:aws:iam::123456789012:role/role1": {
					"arn:aws:eks::123456789012:policy/policy1",
					"arn:aws:eks::123456789012:policy/policy2",
				},
				"arn:aws:iam::123456789012:role/role2": {
					"arn:aws:eks::123456789012:policy/policy3",
				},
			},
			expectedOutput: []string{
				"ACCESS ENTRY PRINCIPAL ARN",
				"KUBERNETES GROUPS",
				"ACCESS POLICIES",
				"arn:aws:iam::123456789012:role/role1",
				"system:masters,group1",
				"arn:aws:eks::123456789012:policy/policy1,arn:aws:eks::123456789012:policy/policy2",
				"arn:aws:iam::123456789012:role/role2",
				"system:authenticated",
				"arn:aws:eks::123456789012:policy/policy3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &mockEKSClient{
				listAssociatedAccessPoliciesFunc: func(ctx context.Context, params *eks.ListAssociatedAccessPoliciesInput) (*eks.ListAssociatedAccessPoliciesOutput, error) {
					policies := tt.mockPolicies[*params.PrincipalArn]
					var associatedPolicies []types.AssociatedAccessPolicy
					for _, policyARN := range policies {
						associatedPolicies = append(associatedPolicies, types.AssociatedAccessPolicy{
							PolicyArn: stringPtr(policyARN),
						})
					}
					return &eks.ListAssociatedAccessPoliciesOutput{
						AssociatedAccessPolicies: associatedPolicies,
					}, nil
				},
			}

			client := &EKSClient{
				client:      mockClient,
				clusterName: stringPtr("test-cluster"),
			}

			// Create printer
			printer := NewAccessEntryPrinter(client)

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create access entry list
			list := &AccessEntryList{
				Items: tt.accessEntries,
			}

			// Print access entries
			err := printer.PrintObj(list, buf)
			if err != nil {
				t.Fatalf("PrintObj returned error: %v", err)
			}

			// Check output contains expected strings
			output := buf.String()
			for _, expected := range tt.expectedOutput {
				if !strings.Contains(output, expected) {
					t.Errorf("Output does not contain expected string: %s\nGot: %s", expected, output)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
