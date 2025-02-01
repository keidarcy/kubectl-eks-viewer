package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewPodIdentityAssociationPrinter(t *testing.T) {
	tests := []struct {
		name           string
		associations   []types.PodIdentityAssociation
		expectedOutput []string
	}{
		{
			name: "single association with owner",
			associations: []types.PodIdentityAssociation{
				{
					AssociationArn: stringPtr("arn:aws:eks:us-west-2:123456789012:podidentityassociation/test-cluster/assoc-1"),
					Namespace:      stringPtr("default"),
					ServiceAccount: stringPtr("app-sa"),
					RoleArn:        stringPtr("arn:aws:iam::123456789012:role/app-role"),
					OwnerArn:       stringPtr("arn:aws:iam::123456789012:role/owner-role"),
				},
			},
			expectedOutput: []string{
				"ARN",
				"NAMESPACE",
				"SERVICE ACCOUNT NAME",
				"IAM ROLE ARN",
				"OWNER ARN",
				"arn:aws:eks:us-west-2:123456789012:podidentityassociation/test-cluster/assoc-1",
				"default",
				"app-sa",
				"arn:aws:iam::123456789012:role/app-role",
				"arn:aws:iam::123456789012:role/owner-role",
			},
		},
		{
			name: "association without owner",
			associations: []types.PodIdentityAssociation{
				{
					AssociationArn: stringPtr("arn:aws:eks:us-west-2:123456789012:podidentityassociation/test-cluster/assoc-2"),
					Namespace:      stringPtr("kube-system"),
					ServiceAccount: stringPtr("system-sa"),
					RoleArn:        stringPtr("arn:aws:iam::123456789012:role/system-role"),
				},
			},
			expectedOutput: []string{
				"ARN",
				"NAMESPACE",
				"SERVICE ACCOUNT NAME",
				"IAM ROLE ARN",
				"OWNER ARN",
				"arn:aws:eks:us-west-2:123456789012:podidentityassociation/test-cluster/assoc-2",
				"kube-system",
				"system-sa",
				"arn:aws:iam::123456789012:role/system-role",
				"<none>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create printer
			printer := NewPodIdentityAssociationPrinter()

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create pod identity association list
			list := &PodIdentityAssociationList{
				Items: tt.associations,
			}

			// Print pod identity associations
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
