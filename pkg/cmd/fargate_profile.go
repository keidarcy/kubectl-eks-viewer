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

type FargateProfile = types.FargateProfile

type FargateProfileList struct {
	metav1.TypeMeta
	Items []types.FargateProfile
}

// Implement runtime.Object interface
func (f *FargateProfileList) GetObjectKind() schema.ObjectKind {
	return &f.TypeMeta
}

func (f *FargateProfileList) DeepCopyObject() runtime.Object {
	return &FargateProfileList{
		TypeMeta: f.TypeMeta,
		Items:    append([]types.FargateProfile(nil), f.Items...),
	}
}

func NewFargateProfilePrinter() printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*FargateProfileList)
		if !ok {
			return fmt.Errorf("expected *FargateProfileList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "FargateProfile",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "NAME", Type: "string"},
				{Name: "SELECTOR NAMESPACE", Type: "string"},
				{Name: "SELECTOR LABELS", Type: "string"},
				{Name: "POD EXECUTION ROLE ARN", Type: "string"},
				{Name: "SUBNETS", Type: "string"},
				{Name: "STATUS", Type: "string"},
			},
		}

		for _, item := range list.Items {
			selectorNamespace := "<none>"
			selectorLabels := "<none>"
			if len(item.Selectors) > 0 {
				selector := item.Selectors[0]
				if selector.Namespace != nil {
					selectorNamespace = *selector.Namespace
				}
				if len(selector.Labels) > 0 {
					var labels []string
					for k, v := range selector.Labels {
						labels = append(labels, fmt.Sprintf("%s=%s", k, v))
					}
					selectorLabels = strings.Join(labels, ",")
				}
			}

			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.FargateProfileName,
					selectorNamespace,
					selectorLabels,
					*item.PodExecutionRoleArn,
					strings.Join(item.Subnets, ","),
					string(item.Status),
				},
			})
		}

		return printTable(w, table, "fargate-profiles")
	})
}

func (c *EKSClient) ListFargateProfiles(ctx context.Context) ([]types.FargateProfile, error) {
	input := &eks.ListFargateProfilesInput{
		ClusterName: c.clusterName,
	}

	result, err := c.client.ListFargateProfiles(ctx, input)
	if err != nil {
		return nil, err
	}

	var profiles []types.FargateProfile
	for _, profileName := range result.FargateProfileNames {
		profile, err := c.client.DescribeFargateProfile(ctx, &eks.DescribeFargateProfileInput{
			ClusterName:        c.clusterName,
			FargateProfileName: &profileName,
		})
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, *profile.FargateProfile)
	}

	return profiles, nil
}
