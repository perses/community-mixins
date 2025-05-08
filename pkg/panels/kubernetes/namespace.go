package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
)

func NamespaceCPUUsageQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Quota",
		panel.Description("Shows the CPU requests, limits, and usage of pods in a namespace in tabular format."),
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
					"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_requests{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod) / sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_requests{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_limits{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod) / sum(cluster:namespace:pod_cpu:active:kube_pod_container_resource_limits{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func NamespaceMemoryUsageQuota(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Quota",
		panel.Description("Shows the memory requests, limits, and usage of pods in a namespace in tabular format."),
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
					"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\", image!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\", image!=\"\"}) by (pod) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_limits{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\", image!=\"\"}) by (pod) / sum(cluster:namespace:pod_memory:active:kube_pod_container_resource_limits{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_rss{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_cache{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(container_memory_swap{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\"}) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func NamespaceCurrentNetworkUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Network Usage",
		panel.Description("Shows the current network usage of the namespace by pods."),
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
					"sum(rate(container_network_receive_bytes_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_bytes_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_packets_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_packets_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_receive_packets_dropped_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(container_network_transmit_packets_dropped_total{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func NamespaceCurrentStorageIO(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Storage IO",
		panel.Description("Shows the current storage IO of the namespace in tabular form, by pod."),
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
					"sum by(pod) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(pod) (rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(pod) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]) + rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(pod) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(pod) (rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by(pod) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]) + rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
