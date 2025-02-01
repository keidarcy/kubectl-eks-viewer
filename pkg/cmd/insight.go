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

type Insight = types.Insight

type InsightList struct {
	metav1.TypeMeta
	Items []types.Insight
}

// Implement runtime.Object interface
func (i *InsightList) GetObjectKind() schema.ObjectKind {
	return &i.TypeMeta
}

func (i *InsightList) DeepCopyObject() runtime.Object {
	return &InsightList{
		TypeMeta: i.TypeMeta,
		Items:    append([]types.Insight(nil), i.Items...),
	}
}

func NewInsightPrinter() printers.ResourcePrinter {
	return printers.ResourcePrinterFunc(func(obj runtime.Object, w io.Writer) error {
		list, ok := obj.(*InsightList)
		if !ok {
			return fmt.Errorf("expected *InsightList, got %T", obj)
		}

		table := &metav1.Table{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Insight",
			},
			ColumnDefinitions: []metav1.TableColumnDefinition{
				{Name: "NAME", Type: "string"},
				{Name: "CATEGORY", Type: "string"},
				{Name: "STATUS", Type: "string"},
			},
		}

		for _, item := range list.Items {
			status := "<none>"
			if item.InsightStatus != nil {
				status = string(item.InsightStatus.Status)
				if item.InsightStatus.Reason != nil {
					status = fmt.Sprintf("%s (%s)", status, *item.InsightStatus.Reason)
				}
			}

			table.Rows = append(table.Rows, metav1.TableRow{
				Cells: []interface{}{
					*item.Name,
					string(item.Category),
					status,
				},
			})
		}

		return printTable(w, table, "insights")
	})
}

func (c *EKSClient) ListInsights(ctx context.Context) ([]types.Insight, error) {
	input := &eks.ListInsightsInput{
		ClusterName: c.clusterName,
	}

	result, err := c.client.ListInsights(ctx, input)
	if err != nil {
		return nil, err
	}

	var insights []types.Insight
	for _, summary := range result.Insights {
		// Get detailed information for each insight
		detail, err := c.client.DescribeInsight(ctx, &eks.DescribeInsightInput{
			ClusterName: c.clusterName,
			Id:          summary.Id,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe insight %s: %v", *summary.Id, err)
		}
		insights = append(insights, *detail.Insight)
	}

	return insights, nil
}
