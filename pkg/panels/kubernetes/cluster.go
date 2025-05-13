package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/plugins/table/sdk/go"
)

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
					Format: commonSdk.Format{
						Unit: string(commonSdk.DecimalUnit),
					},
				},
				{
					Name:   "value #2",
					Header: "Workloads",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.DecimalUnit),
					},
				},
				{
					Name:   "value #3",
					Header: "CPU Usage",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.DecimalUnit),
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #4",
					Header: "CPU Requests",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.DecimalUnit),
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #5",
					Header: "CPU Requests %",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.PercentDecimalUnit),
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #6",
					Header: "CPU Limits",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.DecimalUnit),
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #7",
					Header: "CPU Limits %",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.PercentDecimalUnit),
						DecimalPlaces: 4,
					},
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
					"sum(kube_pod_owner{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\"}) by (namespace)",
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
					"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\"}) by (namespace)",
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
					"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\"}) by (namespace) / sum(namespace_cpu:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) by (namespace)",
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
					"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\"}) by (namespace) / sum(namespace_cpu:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClusterMemoryUsageQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Requests by Namespace",
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
					Format: commonSdk.Format{
						Unit: string(commonSdk.DecimalUnit),
					},
				},
				{
					Name:   "value #2",
					Header: "Workloads",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.DecimalUnit),
					},
				},
				{
					Name:   "value #3",
					Header: "Memory Usage",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesUnit),
					},
				},
				{
					Name:   "value #4",
					Header: "Memory Requests",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesUnit),
					},
				},
				{
					Name:   "value #5",
					Header: "Memory Requests %",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.PercentDecimalUnit),
						DecimalPlaces: 4,
					},
				},
				{
					Name:   "value #6",
					Header: "Memory Limits",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesUnit),
					},
				},
				{
					Name:   "value #7",
					Header: "Memory Limits %",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit:          string(commonSdk.PercentDecimalUnit),
						DecimalPlaces: 4,
					},
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
					"sum(kube_pod_owner{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\"}) by (namespace)",
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
					"sum(container_memory_rss{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", container!=\"\"}) by (namespace)",
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
					"sum(container_memory_rss{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", container!=\"\"}) by (namespace) / sum(namespace_memory:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) by (namespace)",
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
					"sum(container_memory_rss{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", container!=\"\"}) by (namespace) / sum(namespace_memory:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) by (namespace)",
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
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #2",
					Header: "Current Transmit Bandwidth",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #3",
					Header: "Rate of Received Packets",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
				},
				{
					Name:   "value #4",
					Header: "Rate of Transmitted Packets",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
				},
				{
					Name:   "value #5",
					Header: "Rate of Received Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
				},
				{
					Name:   "value #6",
					Header: "Rate of Transmitted Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
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
					"sum(rate(container_network_receive_bytes_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=~\".+\"}[$__rate_interval])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_bytes_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=~\".+\"}[$__rate_interval])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_packets_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=~\".+\"}[$__rate_interval])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_packets_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=~\".+\"}[$__rate_interval])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_packets_dropped_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=~\".+\"}[$__rate_interval])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_packets_dropped_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=~\".+\"}[$__rate_interval])) by (namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
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
					Format: commonSdk.Format{
						Unit: string(commonSdk.OpsPerSecondsUnit),
					},
				},
				{
					Name:   "value #2",
					Header: "IOPS(Writes)",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.OpsPerSecondsUnit),
					},
				},
				{
					Name:   "value #3",
					Header: "IOPS(Reads + Writes)",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.OpsPerSecondsUnit),
					},
				},
				{
					Name:   "value #4",
					Header: "Throughput(Reads)",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #5",
					Header: "Throughput(Writes)",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #6",
					Header: "Throughput(Reads + Writes)",
					Align:  tablePanel.RightAlign,
					Format: commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
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
					"sum by(namespace) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]) + rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(namespace) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]) + rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
