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

// Use AWS SDK type directly
type PodIdentityAssociation = types.PodIdentityAssociation

type PodIdentityAssociationList struct {
	metav1.TypeMeta
	Items []types.PodIdentityAssociation
}

// Implement runtime.Object interface
func (p *PodIdentityAssociationList) GetObjectKind() schema.ObjectKind {
	return &p.TypeMeta
}

func (p *PodIdentityAssociationList) DeepCopyObject() runtime.Object {
	return &PodIdentityAssociationList{
		TypeMeta: p.TypeMeta,
		Items:    append([]types.PodIdentityAssociation(nil), p.Items...),
	}
}

func NewPodIdentityAssociationPrinter() printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*PodIdentityAssociationList)
		if !ok {
			return fmt.Errorf("expected *PodIdentityAssociationList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "PodIdentityAssociation",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "ARN", Type: "string"},
				{Name: "NAMESPACE", Type: "string"},
				{Name: "SERVICE ACCOUNT NAME", Type: "string"},
				{Name: "IAM ROLE ARN", Type: "string"},
				{Name: "OWNER ARN", Type: "string"},
			},
		}

		for _, item := range list.Items {
			ownerArn := "<none>"
			if item.OwnerArn != nil {
				ownerArn = *item.OwnerArn
			}

			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.AssociationArn,
					*item.Namespace,
					*item.ServiceAccount,
					*item.RoleArn,
					ownerArn,
				},
			})
		}

		return printTable(w, table, "pod-identity-associations")
	})
}

func (c *EKSClient) ListPodIdentityAssociations(ctx context.Context) ([]types.PodIdentityAssociation, error) {
	input := &eks.ListPodIdentityAssociationsInput{
		ClusterName: c.clusterName,
	}

	result, err := c.client.ListPodIdentityAssociations(ctx, input)
	if err != nil {
		return nil, err
	}

	var associations []types.PodIdentityAssociation
	for _, assoc := range result.Associations {
		describeOut, err := c.client.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
			ClusterName:   c.clusterName,
			AssociationId: assoc.AssociationId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe pod identity association with associationID: %s", *assoc.AssociationId)
		}

		associations = append(associations, *describeOut.Association)
	}

	return associations, nil
}
