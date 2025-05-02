package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func NodeCPUUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		panel.Description("Shows the CPU usage of the node by pod, and the CPU capacity of the node."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(kube_node_status_capacity{cluster=\"$cluster\", job=\"kube-state-metrics\", node=~\"$node\", resource=\"cpu\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("max capacity"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}}"),
			),
		),
	)
}

func CPUUsageQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Quota",
		panel.Description("Shows the CPU requests, limits, and usage of pods on the node in tabular format."),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "pod",
					Header: "Pod",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "CPU Usage",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #2",
					Header: "CPU Requests",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #3",
					Header: "CPU Requests %",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #4",
					Header: "CPU Limits",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #5",
					Header: "CPU Limits %",
					Align:  tablePanel.RightAlign,
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_requests{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", node=~\"$node\"}) by (pod) / sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_requests{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_limits{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", node=~\"$node\"}) by (pod) / sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_limits{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func NodeMemoryUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage with Cache",
		panel.Description("Shows the memory usage of the node by pod, and the memory capacity of the node."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(kube_node_status_capacity{cluster=\"$cluster\", job=\"kube-state-metrics\", node=~\"$node\", resource=\"memory\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("max capacity"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_working_set_bytes{cluster=\"$cluster\", node=~\"$node\", container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}}"),
			),
		),
	)
}

func NodeMemoryUsageWithoutCache(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage Without Cache",
		panel.Description("Shows the memory usage (RSS) of the node by pod, and the memory capacity of the node."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(kube_node_status_capacity{cluster=\"$cluster\", job=\"kube-state-metrics\", node=~\"$node\", resource=\"memory\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("max capacity"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_rss{cluster=\"$cluster\", node=~\"$node\", container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}}"),
			),
		),
	)
}

func MemoryQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Quota",
		panel.Description("Shows the Memory requests, limits, and usage of pods on the node in tabular format."),
		tablePanel.Table(
			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "pod",
					Header: "Pod",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Memory Usage",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #2",
					Header: "Memory Requests",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #3",
					Header: "Memory Requests %",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #4",
					Header: "Memory Limits",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #5",
					Header: "Memory Limits %",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #6",
					Header: "Memory Usage (RSS)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #7",
					Header: "Memory Usage (Cache)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #8",
					Header: "Memory Usage (Swap)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_working_set_bytes{cluster=\"$cluster\", node=~\"$node\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_working_set_bytes{cluster=\"$cluster\", node=~\"$node\",container!=\"\"}) by (pod) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_limits{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_working_set_bytes{cluster=\"$cluster\", node=~\"$node\",container!=\"\"}) by (pod) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_limits{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_rss{cluster=\"$cluster\", node=~\"$node\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_cache{cluster=\"$cluster\", node=~\"$node\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_memory_swap{cluster=\"$cluster\", node=~\"$node\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
