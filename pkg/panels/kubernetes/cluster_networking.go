package kubernetes

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func ClusterTCPRetransmitRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of TCP Retransmits out of all sent segments",
		panel.Description("Shows the rate of TCP retransmits out of all sent segments in a cluster."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PercentUnit,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (instance) (\n    rate(node_netstat_Tcp_RetransSegs{cluster=\"$cluster\"}[$__rate_interval]) / rate(node_netstat_Tcp_OutSegs{cluster=\"$cluster\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

func ClusterTCPSYNRetransmitRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of TCP SYN Retransmits out of all sent segments",
		panel.Description("Shows the rate of TCP SYN retransmits out of all sent segments in a cluster."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PercentUnit,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (instance) (\n    rate(node_netstat_TcpExt_TCPSynRetrans{cluster=\"$cluster\"}[$__rate_interval]) / rate(node_netstat_Tcp_RetransSegs{cluster=\"$cluster\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}}"),
			),
		),
	)
}

func ClusterNetworkingCurrentStatus(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Status",
		panel.Description("Shows the current network status of the cluster by namespace."),
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
					Header: "Rx Bytes",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #2",
					Header: "Tx Bytes",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #3",
					Header: "Rx Bytes (Avg)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #4",
					Header: "Tx Bytes (Avg)",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.BytesPerSecondsUnit,
					},
				},
				{
					Name:   "value #5",
					Header: "Rx Packets",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name:   "value #6",
					Header: "Tx Packets",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name:   "value #7",
					Header: "Rx Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
					},
				},
				{
					Name:   "value #8",
					Header: "Tx Packets Dropped",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit: &dashboards.PacketsPerSecondsUnit,
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
					"sum by (namespace) (\n    rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace) (\n    rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg by (namespace) (\n    rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg by (namespace) (\n    rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace) (\n    rate(container_network_receive_packets_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace) (\n    rate(container_network_transmit_packets_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace) (\n    rate(container_network_receive_packets_dropped_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n", labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace) (\n    rate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\",namespace!=\"\"}[$__rate_interval])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n", labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
