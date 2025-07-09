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

func WorkloadNamespaceCurrentNetworkStatus(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Network Status",
		panel.Description("Shows the current network status of the namespace by workloads."),
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
					Name:   "workload",
					Header: "Workload",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "workload_type",
					Header: "Type",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Rx Bytes",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #2",
					Header: "Tx Bytes",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #3",
					Header: "Rx Bytes (Avg)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #4",
					Header: "Tx Bytes (Avg)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.BytesPerSecondsUnit),
					},
				},
				{
					Name:   "value #5",
					Header: "Rx Packets",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
				},
				{
					Name:   "value #6",
					Header: "Tx Packets",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
				},
				{
					Name:   "value #7",
					Header: "Rx Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: string(commonSdk.PacketsPerSecondsUnit),
					},
				},
				{
					Name:   "value #8",
					Header: "Tx Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
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
					"sort_desc(sum(rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(sum(rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(avg(rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(avg(rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(sum(rate(container_network_receive_packets_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(sum(rate(container_network_transmit_packets_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(sum(rate(container_network_receive_packets_dropped_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sort_desc(sum(rate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\",namespace=\"$namespace\"}[$__rate_interval])\n* on (namespace,pod) kube_pod_info{cluster=\"$cluster\",namespace=\"$namespace\",host_network=\"false\"}\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload, workload_type))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
