package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	statPanel "github.com/perses/perses/go-sdk/panel/stat"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func CPUUtilizationStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Utilization",
		panel.Description("Shows the CPU utilization of the cluster."),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"cluster:node_cpu:ratio_rate5m{cluster=\"$cluster\"}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func CPURequestsCommitmentStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Requests Commitment",
		panel.Description("Shows the CPU requests commitment of the cluster."),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_cpu:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{job=\"kube-state-metrics\",resource=\"cpu\",cluster=\"$cluster\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func CPULimitsCommitmentStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Limits Commitment",
		panel.Description("Shows the CPU limits commitment of the cluster."),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_cpu:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{job=\"kube-state-metrics\",resource=\"cpu\",cluster=\"$cluster\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func MemoryUtilizationStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Utilization",
		panel.Description("Shows the memory utilization of the cluster."),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"1 - sum(:node_memory_MemAvailable_bytes:sum{cluster=\"$cluster\"}) / sum(node_memory_MemTotal_bytes{job=\"node-exporter\",cluster=\"$cluster\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func MemoryRequestsCommitmentStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Requests Commitment",
		panel.Description("Shows the memory requests commitment of the cluster."),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_memory:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{job=\"kube-state-metrics\",resource=\"memory\",cluster=\"$cluster\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func MemoryLimitsCommitmentStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Limits Commitment",
		panel.Description("Shows the memory limits commitment of the cluster."),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_memory:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{job=\"kube-state-metrics\",resource=\"memory\",cluster=\"$cluster\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterCPUUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		panel.Description("Shows the CPU usage of the cluster by namespace."),
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
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterCPUUsageQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Quota",
		panel.Description("Shows the CPU requests, limits, and usage of workloads by namespace in tabular format."),
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
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Pods",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #2",
					Header: "Workloads",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #3",
					Header: "CPU Usage",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #4",
					Header: "CPU Requests",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #5",
					Header: "CPU Requests %",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #6",
					Header: "CPU Limits",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #7",
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
					"sum(kube_pod_owner{job=\"kube-state-metrics\", cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count(avg(namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\"}) by (workload, namespace)) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_cpu:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\"}) by (namespace) / sum(namespace_cpu:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_cpu:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\"}) by (namespace) / sum(namespace_cpu:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterMemoryUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
		panel.Description("Shows the memory usage of the cluster by namespace."),
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
					"sum(container_memory_rss{job=\"cadvisor\", cluster=\"$cluster\", container!=\"\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterMemoryUsageQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Quota",
		panel.Description("Shows the memory requests, limits, and usage of workloads by namespace in tabular format."),
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
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Pods",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #2",
					Header: "Workloads",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #3",
					Header: "Memory Usage",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #4",
					Header: "Memory Requests",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #5",
					Header: "Memory Requests %",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #6",
					Header: "Memory Limits",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #7",
					Header: "Memory Limits %",
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
					"sum(kube_pod_owner{job=\"kube-state-metrics\", cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"count(avg(namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\"}) by (workload, namespace)) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_rss{job=\"cadvisor\", cluster=\"$cluster\", container!=\"\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_memory:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_rss{job=\"cadvisor\", cluster=\"$cluster\", container!=\"\"}) by (namespace) / sum(namespace_memory:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(namespace_memory:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_rss{job=\"cadvisor\", cluster=\"$cluster\", container!=\"\"}) by (namespace) / sum(namespace_memory:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterCurrentNetworkUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Network Usage",
		panel.Description("Shows the current network usage of the cluster by namespace."),
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
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Current Receive Bandwidth",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #2",
					Header: "Current Transmit Bandwidth",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #3",
					Header: "Rate of Received Packets",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #4",
					Header: "Rate of Transmitted Packets",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #5",
					Header: "Rate of Received Packets Dropped",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #6",
					Header: "Rate of Transmitted Packets Dropped",
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
					"sum(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterReceiveBandwidth(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Receive Bandwidth",
		panel.Description("Shows the network receive bandwidth of the cluster by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterTransmitBandwidth(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Transmit Bandwidth",
		panel.Description("Shows the network transmit bandwidth of the cluster by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterAvgContainerBandwidthReceived(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Container Bandwidth by Namespace: Received",
		panel.Description("Shows the average network bandwidth received in container by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg(irate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterAvgContainerBandwidthTransmitted(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Container Bandwidth by Namespace: Transmitted",
		panel.Description("Shows the average network bandwidth transmitted in container by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg(irate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterReceivedPackets(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Received Packets",
		panel.Description("Shows the rate of received packets by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PacketsPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(container_network_receive_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterTransmittedPackets(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Transmitted Packets",
		panel.Description("Shows the rate of transmitted packets by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PacketsPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(container_network_transmit_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterReceivedPacketsDropped(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Received Packets Dropped",
		panel.Description("Shows the rate of received packets dropped by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PacketsPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(container_network_receive_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterTransmittedPacketsDropped(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Transmitted Packets Dropped",
		panel.Description("Shows the rate of transmitted packets dropped by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PacketsPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(container_network_transmit_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=~\".+\"}[5m])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterIOPS(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("IOPS(Reads+Writes)",
		panel.Description("Shows IOPS(Reads+Writes) by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.OpsPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"ceil(sum by(namespace) (rate(container_fs_reads_total{job=\"cadvisor\", container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace!=\"\"}[5m]) + rate(container_fs_writes_total{job=\"cadvisor\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterThroughput(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("ThroughPut(Reads+Writes)",
		panel.Description("Shows ThroughPut(Reads+Writes) by namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_bytes_total{job=\"cadvisor\", container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace!=\"\"}[5m]) + rate(container_fs_writes_bytes_total{job=\"cadvisor\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}}"),
			),
		),
	)
}

func ClusterCurrentStorageIO(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Storage IO",
		panel.Description("Shows the current storage IO of the cluster in tabular form, by namespace."),
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
					Name:   "namespace",
					Header: "Namespace",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "IOPS(Reads)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #2",
					Header: "IOPS(Writes)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #3",
					Header: "IOPS(Reads + Writes)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #4",
					Header: "Throughput(Reads)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #5",
					Header: "Throughput(Writes)",
					Align:  tablePanel.RightAlign,
				},
				{
					Name:   "value #6",
					Header: "Throughput(Reads + Writes)",
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
					"sum by(namespace) (rate(container_fs_reads_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_writes_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]) + rate(container_fs_writes_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_bytes_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_writes_bytes_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_bytes_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]) + rate(container_fs_writes_bytes_total{job=\"cadvisor\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
