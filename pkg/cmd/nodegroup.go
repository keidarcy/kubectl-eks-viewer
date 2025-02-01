package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
)

type Nodegroup = types.Nodegroup

type NodeGroupList struct {
	metav1.TypeMeta
	Items []types.Nodegroup
}

// Implement runtime.Object interface
func (n *NodeGroupList) GetObjectKind() schema.ObjectKind {
	return &n.TypeMeta
}

func (n *NodeGroupList) DeepCopyObject() runtime.Object {
	return &NodeGroupList{
		TypeMeta: n.TypeMeta,
		Items:    append([]types.Nodegroup(nil), n.Items...),
	}
}

func NewNodeGroupPrinter() printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*NodeGroupList)
		if !ok {
			return fmt.Errorf("expected *NodeGroupList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "NodeGroup",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "NAME", Type: "string"},
				{Name: "STATUS", Type: "string"},
				{Name: "INSTANCE TYPE", Type: "string"},
				{Name: "DESIRED SIZE", Type: "integer"},
				{Name: "MIN SIZE", Type: "integer"},
				{Name: "MAX SIZE", Type: "integer"},
				{Name: "VERSION", Type: "string"},
				{Name: "AMI TYPE", Type: "string"},
				{Name: "CAPACITY TYPE", Type: "string"},
			},
		}

		for _, item := range list.Items {
			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.NodegroupName,
					string(item.Status),
					strings.Join(item.InstanceTypes, ","),
					int(*item.ScalingConfig.DesiredSize),
					int(*item.ScalingConfig.MinSize),
					int(*item.ScalingConfig.MaxSize),
					*item.Version,
					string(item.AmiType),
					string(item.CapacityType),
				},
			})
		}

		return printTable(w, table, "nodegroups")
	})
}

func (c *EKSClient) ListNodeGroups(ctx context.Context) ([]types.Nodegroup, error) {
	input := &eks.ListNodegroupsInput{
		ClusterName: c.clusterName,
	}

	result, err := c.client.ListNodegroups(ctx, input)
	if err != nil {
		return nil, err
	}

	var nodeGroups []types.Nodegroup
	for _, ngName := range result.Nodegroups {
		ngOutput, err := c.client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
			ClusterName:   c.clusterName,
			NodegroupName: &ngName,
		})
		if err != nil {
			return nil, err
		}

		nodeGroups = append(nodeGroups, *ngOutput.Nodegroup)
	}

	return nodeGroups, nil
}
