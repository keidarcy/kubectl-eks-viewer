package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewAddonPrinter(t *testing.T) {
	tests := []struct {
		name           string
		addons         []types.Addon
		expectedOutput []string
	}{
		{
			name: "single addon with no issues",
			addons: []types.Addon{
				{
					AddonName:    stringPtr("vpc-cni"),
					AddonVersion: stringPtr("v1.12.0"),
					Status:       types.AddonStatusActive,
					Health: &types.AddonHealth{
						Issues: []types.AddonIssue{},
					},
				},
			},
			expectedOutput: []string{
				"NAME",
				"VERSION",
				"STATUS",
				"ISSUES",
				"vpc-cni",
				"v1.12.0",
				"ACTIVE",
				"0",
			},
		},
		{
			name: "multiple addons with issues",
			addons: []types.Addon{
				{
					AddonName:    stringPtr("vpc-cni"),
					AddonVersion: stringPtr("v1.12.0"),
					Status:       types.AddonStatusActive,
					Health: &types.AddonHealth{
						Issues: []types.AddonIssue{
							{
								Code:        types.AddonIssueCodeConfigurationConflict,
								Message:     stringPtr("Configuration conflict detected"),
								ResourceIds: []string{"resource1"},
							},
						},
					},
				},
				{
					AddonName:    stringPtr("kube-proxy"),
					AddonVersion: stringPtr("v1.23.0"),
					Status:       types.AddonStatusDegraded,
					Health: &types.AddonHealth{
						Issues: []types.AddonIssue{
							{
								Code:        types.AddonIssueCodeAccessDenied,
								Message:     stringPtr("Health check failed"),
								ResourceIds: []string{"resource2"},
							},
							{
								Code:        types.AddonIssueCodeConfigurationConflict,
								Message:     stringPtr("Version mismatch detected"),
								ResourceIds: []string{"resource3"},
							},
						},
					},
				},
			},
			expectedOutput: []string{
				"NAME",
				"VERSION",
				"STATUS",
				"ISSUES",
				"vpc-cni",
				"v1.12.0",
				"ACTIVE",
				"1",
				"kube-proxy",
				"v1.23.0",
				"DEGRADED",
				"2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create printer
			printer := NewAddonPrinter()

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create addon list
			list := &AddonList{
				Items: tt.addons,
			}

			// Print addons
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
