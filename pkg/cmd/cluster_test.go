package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewClusterPrinter(t *testing.T) {
	tests := []struct {
		name           string
		clusters       []types.Cluster
		expectedOutput []string
	}{
		{
			name: "single cluster with standard auth",
			clusters: []types.Cluster{
				{
					Name:            stringPtr("test-cluster"),
					Version:         stringPtr("1.24"),
					Status:          types.ClusterStatusActive,
					PlatformVersion: stringPtr("eks.1"),
				},
			},
			expectedOutput: []string{
				"NAME",
				"VERSION",
				"STATUS",
				"PLATFORM VERSION",
				"AUTH MODE",
				"test-cluster",
				"1.24",
				"ACTIVE",
				"eks.1",
				"<none>",
			},
		},
		{
			name: "single cluster with config map auth",
			clusters: []types.Cluster{
				{
					Name:            stringPtr("prod-cluster"),
					Version:         stringPtr("1.25"),
					Status:          types.ClusterStatusUpdating,
					PlatformVersion: stringPtr("eks.2"),
				},
			},
			expectedOutput: []string{
				"NAME",
				"VERSION",
				"STATUS",
				"PLATFORM VERSION",
				"AUTH MODE",
				"prod-cluster",
				"1.25",
				"UPDATING",
				"eks.2",
				"<none>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create printer
			printer := NewClusterPrinter()

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create cluster list
			list := &ClusterList{
				Items: tt.clusters,
			}

			// Print clusters
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
