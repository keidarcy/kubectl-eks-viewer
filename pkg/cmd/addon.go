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

type Addon = types.Addon

type AddonList struct {
	metav1.TypeMeta
	Items []types.Addon
}

// Implement runtime.Object interface
func (a *AddonList) GetObjectKind() schema.ObjectKind {
	return &a.TypeMeta
}

func (a *AddonList) DeepCopyObject() runtime.Object {
	return &AddonList{
		TypeMeta: a.TypeMeta,
		Items:    append([]types.Addon(nil), a.Items...),
	}
}

func NewAddonPrinter() printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*AddonList)
		if !ok {
			return fmt.Errorf("expected *AddonList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Addon",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "NAME", Type: "string"},
				{Name: "VERSION", Type: "string"},
				{Name: "STATUS", Type: "string"},
				{Name: "ISSUES", Type: "integer"},
			},
		}

		for _, item := range list.Items {
			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.AddonName,
					*item.AddonVersion,
					string(item.Status),
					len(item.Health.Issues),
				},
			})
		}

		return printTable(w, table, "addons")
	})
}

func (c *EKSClient) ListAddons(ctx context.Context) ([]types.Addon, error) {
	input := &eks.ListAddonsInput{
		ClusterName: c.clusterName,
	}

	result, err := c.client.ListAddons(ctx, input)
	if err != nil {
		return nil, err
	}

	var addons []types.Addon
	for _, addonName := range result.Addons {
		addonOutput, err := c.client.DescribeAddon(ctx, &eks.DescribeAddonInput{
			ClusterName: c.clusterName,
			AddonName:   &addonName,
		})
		if err != nil {
			return nil, err
		}

		addons = append(addons, *addonOutput.Addon)
	}

	return addons, nil
}
