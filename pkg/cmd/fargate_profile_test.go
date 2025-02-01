package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewFargateProfilePrinter(t *testing.T) {
	tests := []struct {
		name           string
		profiles       []types.FargateProfile
		expectedOutput []string
	}{
		{
			name: "single profile with namespace selector",
			profiles: []types.FargateProfile{
				{
					FargateProfileName:  stringPtr("default"),
					PodExecutionRoleArn: stringPtr("arn:aws:iam::123456789012:role/eks-fargate-pods"),
					Status:              types.FargateProfileStatusActive,
					Subnets:             []string{"subnet-1234", "subnet-5678"},
					Selectors: []types.FargateProfileSelector{
						{
							Namespace: stringPtr("default"),
							Labels: map[string]string{
								"environment": "prod",
							},
						},
					},
				},
			},
			expectedOutput: []string{
				"NAME",
				"SELECTOR NAMESPACE",
				"SELECTOR LABELS",
				"POD EXECUTION ROLE ARN",
				"SUBNETS",
				"STATUS",
				"default",
				"default",
				"environment=prod",
				"arn:aws:iam::123456789012:role/eks-fargate-pods",
				"subnet-1234,subnet-5678",
				"ACTIVE",
			},
		},
		{
			name: "profile without selectors",
			profiles: []types.FargateProfile{
				{
					FargateProfileName:  stringPtr("minimal"),
					PodExecutionRoleArn: stringPtr("arn:aws:iam::123456789012:role/eks-fargate-pods"),
					Status:              types.FargateProfileStatusCreating,
					Subnets:             []string{"subnet-abcd"},
					Selectors:           []types.FargateProfileSelector{},
				},
			},
			expectedOutput: []string{
				"NAME",
				"SELECTOR NAMESPACE",
				"SELECTOR LABELS",
				"POD EXECUTION ROLE ARN",
				"SUBNETS",
				"STATUS",
				"minimal",
				"<none>",
				"<none>",
				"arn:aws:iam::123456789012:role/eks-fargate-pods",
				"subnet-abcd",
				"CREATING",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create printer
			printer := NewFargateProfilePrinter()

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create fargate profile list
			list := &FargateProfileList{
				Items: tt.profiles,
			}

			// Print fargate profiles
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
