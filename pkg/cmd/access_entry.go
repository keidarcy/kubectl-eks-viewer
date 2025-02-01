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

type AccessEntry = types.AccessEntry

type AccessEntryList struct {
	metav1.TypeMeta
	Items []AccessEntry
}

// Implement runtime.Object interface
func (a *AccessEntryList) GetObjectKind() schema.ObjectKind {
	return &a.TypeMeta
}

func (a *AccessEntryList) DeepCopyObject() runtime.Object {
	return &AccessEntryList{
		TypeMeta: a.TypeMeta,
		Items:    append([]AccessEntry(nil), a.Items...),
	}
}

func NewAccessEntryPrinter(client *EKSClient) printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*AccessEntryList)
		if !ok {
			return fmt.Errorf("expected *AccessEntryList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "AccessEntry",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "ACCESS ENTRY PRINCIPAL ARN", Type: "string"},
				{Name: "KUBERNETES GROUPS", Type: "string"},
				{Name: "ACCESS POLICIES", Type: "string"},
			},
		}

		for _, item := range list.Items {
			// Get associated policies for each access entry
			policies, err := client.client.ListAssociatedAccessPolicies(context.Background(), &eks.ListAssociatedAccessPoliciesInput{
				ClusterName:  client.clusterName,
				PrincipalArn: item.PrincipalArn,
			})
			if err != nil {
				return err
			}

			var policyARNs []string
			for _, policy := range policies.AssociatedAccessPolicies {
				policyARNs = append(policyARNs, *policy.PolicyArn)
			}

			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.PrincipalArn,
					strings.Join(item.KubernetesGroups, ","),
					strings.Join(policyARNs, ","),
				},
			})
		}

		return printTable(w, table, "access-entries")
	})
}

func (c *EKSClient) ListAccessEntries(ctx context.Context) ([]AccessEntry, error) {
	input := &eks.ListAccessEntriesInput{
		ClusterName: c.clusterName,
	}

	result, err := c.client.ListAccessEntries(ctx, input)
	if err != nil {
		return nil, err
	}

	var accessEntries []AccessEntry
	for _, principalARN := range result.AccessEntries {
		// Get Kubernetes groups
		entry, err := c.client.DescribeAccessEntry(ctx, &eks.DescribeAccessEntryInput{
			ClusterName:  c.clusterName,
			PrincipalArn: &principalARN,
		})
		if err != nil {
			return nil, err
		}

		accessEntries = append(accessEntries, *entry.AccessEntry)
	}

	return accessEntries, nil
}
