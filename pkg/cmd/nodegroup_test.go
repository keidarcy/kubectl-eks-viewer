package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestNewNodeGroupPrinter(t *testing.T) {
	tests := []struct {
		name           string
		nodegroups     []types.Nodegroup
		expectedOutput []string
	}{
		{
			name: "single nodegroup with managed nodes",
			nodegroups: []types.Nodegroup{
				{
					NodegroupName: stringPtr("managed-ng-1"),
					Status:        "ACTIVE",
					InstanceTypes: []string{"t3.medium"},
					ScalingConfig: &types.NodegroupScalingConfig{
						DesiredSize: int32Ptr(2),
						MinSize:     int32Ptr(1),
						MaxSize:     int32Ptr(4),
					},
					Version:      stringPtr("1.24"),
					AmiType:      "AL2_x86_64",
					CapacityType: "ON_DEMAND",
					DiskSize:     int32Ptr(20),
				},
			},
			expectedOutput: []string{
				"NAME",
				"STATUS",
				"INSTANCE TYPE",
				"DESIRED SIZE",
				"MIN SIZE",
				"MAX SIZE",
				"VERSION",
				"AMI TYPE",
				"CAPACITY TYPE",
				"managed-ng-1",
				"ACTIVE",
				"t3.medium",
				"2",
				"1",
				"4",
				"1.24",
				"AL2_x86_64",
				"ON_DEMAND",
			},
		},
		{
			name: "spot nodegroup with multiple instance types",
			nodegroups: []types.Nodegroup{
				{
					NodegroupName: stringPtr("spot-ng-1"),
					Status:        "CREATING",
					InstanceTypes: []string{"t3.large", "t3.xlarge"},
					ScalingConfig: &types.NodegroupScalingConfig{
						DesiredSize: int32Ptr(3),
						MinSize:     int32Ptr(1),
						MaxSize:     int32Ptr(10),
					},
					Version:      stringPtr("1.25"),
					AmiType:      "AL2_ARM_64",
					CapacityType: "SPOT",
					DiskSize:     int32Ptr(50),
				},
			},
			expectedOutput: []string{
				"NAME",
				"STATUS",
				"INSTANCE TYPE",
				"DESIRED SIZE",
				"MIN SIZE",
				"MAX SIZE",
				"VERSION",
				"AMI TYPE",
				"CAPACITY TYPE",
				"spot-ng-1",
				"CREATING",
				"t3.large,t3.xlarge",
				"3",
				"1",
				"10",
				"1.25",
				"AL2_ARM_64",
				"SPOT",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create printer
			printer := NewNodeGroupPrinter()

			// Create buffer to capture output
			buf := &bytes.Buffer{}

			// Create nodegroup list
			list := &NodeGroupList{
				Items: tt.nodegroups,
			}

			// Print nodegroups
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

func int32Ptr(i int32) *int32 {
	return &i
}
