package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewInsightPrinter(t *testing.T) {
	tests := []struct {
		name           string
		insights       []types.Insight
		expectedOutput []string
	}{
		{
			name: "single insight with no compatibility details",
			insights: []types.Insight{
				{
					Name:     stringPtr("test-insight"),
					Category: "ADDON",
					InsightStatus: &types.InsightStatus{
						Status: "ACTIVE",
						Reason: stringPtr("All good"),
					},
				},
			},
			expectedOutput: []string{
				"NAME",
				"CATEGORY",
				"STATUS",
				"test-insight",
				"ADDON",
				"ACTIVE (All good)",
			},
		},
		{
			name: "insight with addon compatibility details",
			insights: []types.Insight{
				{
					Name:     stringPtr("addon-insight"),
					Category: "ADDON",
					InsightStatus: &types.InsightStatus{
						Status: "DEGRADED",
						Reason: stringPtr("Version mismatch"),
					},
				},
			},
			expectedOutput: []string{
				"NAME",
				"CATEGORY",
				"STATUS",
				"addon-insight",
				"ADDON",
				"DEGRADED (Version mismatch)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create printer
			printer := NewInsightPrinter()

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create insight list
			list := &InsightList{
				Items: tt.insights,
			}

			// Print insights
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
