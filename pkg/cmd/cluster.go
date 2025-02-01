package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/printers"
)

type Cluster = types.Cluster

type ClusterList struct {
	metav1.TypeMeta
	Items []types.Cluster
}

// Implement runtime.Object interface
func (c *ClusterList) GetObjectKind() schema.ObjectKind {
	return &c.TypeMeta
}

func (c *ClusterList) DeepCopyObject() runtime.Object {
	return &ClusterList{
		TypeMeta: c.TypeMeta,
		Items:    append([]types.Cluster(nil), c.Items...),
	}
}

func NewClusterPrinter() printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*ClusterList)
		if !ok {
			return fmt.Errorf("expected *ClusterList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Cluster",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "NAME", Type: "string"},
				{Name: "VERSION", Type: "string"},
				{Name: "STATUS", Type: "string"},
				{Name: "PLATFORM VERSION", Type: "string"},
				{Name: "AUTH MODE", Type: "string"},
			},
		}

		for _, item := range list.Items {
			authMode := "<none>"
			if item.AccessConfig != nil {
				authMode = string(item.AccessConfig.AuthenticationMode)
			}

			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.Name,
					*item.Version,
					string(item.Status),
					*item.PlatformVersion,
					authMode,
				},
			})
		}

		return printTable(w, table, "cluster")
	})
}

func (c *EKSClient) DescribeCluster(ctx context.Context) ([]types.Cluster, error) {
	input := &eks.DescribeClusterInput{
		Name: c.clusterName,
	}

	result, err := c.client.DescribeCluster(ctx, input)
	if err != nil {
		return nil, err
	}

	return []types.Cluster{*result.Cluster}, nil
}

func printTable(w io.Writer, table *metav1.Table, resourceType string) error {
	fmt.Fprintf(w, "=== %s ===\n", resourceType)

	if len(table.Rows) == 0 {
		fmt.Fprintf(w, "<none>\n")
		return nil
	}

	printer := printers.NewTablePrinter(printers.PrintOptions{})
	return printer.PrintObj(table, w)
}
