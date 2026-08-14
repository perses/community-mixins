// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func KubernetesReceiveBandwidth(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "cluster":
		description = "Shows the network receive bandwidth of the cluster by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "cluster-networking":
		description = "Shows the network receive bandwidth of the cluster highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the network receive bandwidth of the namespace by pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the network receive bandwidth of the namespace by workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the network receive bandwidth of the namespace by workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the network receive bandwidth of the namespace by pod highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload":
		description = "Shows the network receive bandwidth of the workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the network receive bandwidth of the workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the network receive bandwidth of the pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the network receive bandwidth of the pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceiveBandwidth10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}

	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Receive Bandwidth", panelOpts...)
}

func KubernetesTransmitBandwidth(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "cluster":
		description = "Shows the network transmit bandwidth of the cluster by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "cluster-networking":
		description = "Shows the network transmit bandwidth of the cluster highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the network transmit bandwidth of the namespace by pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the network transmit bandwidth of the namespace by workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the network transmit bandwidth of the namespace by workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the network transmit bandwidth of the namespace by pod highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload":
		description = "Shows the network transmit bandwidth of the workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the network transmit bandwidth of the workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the network transmit bandwidth of the pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the network transmit bandwidth of the pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmitBandwidth10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}

	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Transmit Bandwidth", panelOpts...)
}

func KubernetesAvgContainerBandwidthTransmitted(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		panelName = "Average Container Bandwidth by Namespace: Transmitted"
		description = "Shows the average network bandwidth transmitted in container by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthTransmitted1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-workload":
		panelName = "Average Container Bandwidth by Workload: Transmitted"
		description = "Shows the average network bandwidth transmitted in container by workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthTransmitted2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "namespace-workload-networking":
		panelName = "Average Container Bandwidth by Workload: Transmitted"
		description = "Shows the average network bandwidth transmitted in container by workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthTransmitted3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload":
		panelName = "Average Container Bandwidth by Pod: Transmitted"
		description = "Shows the average network bandwidth transmitted by containers of a pod in a workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthTransmitted4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesAvgContainerBandwidthReceived(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		panelName = "Average Container Bandwidth by Namespace: Received"
		description = "Shows the average network bandwidth received in container by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthReceived1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-workload":
		panelName = "Average Container Bandwidth by Workload: Received"
		description = "Shows the average network bandwidth received in container by workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthReceived2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "namespace-workload-networking":
		panelName = "Average Container Bandwidth by Workload: Received"
		description = "Shows the average network bandwidth received in container by workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthReceived3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload":
		panelName = "Average Container Bandwidth by Pod: Received"
		description = "Shows the average network bandwidth received by containers of a pod in a workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAvgContainerBandwidthReceived4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesReceivedPackets(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of received packets by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the rate of received packets by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the rate of received packets by a pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload":
		description = "Shows the rate of received packets by pods in a workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the rate of received packets by pods in a workload in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "cluster-networking":
		description = "Shows the rate of received packets by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the rate of received packets by pods in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the rate of received packets by top pods in a workload in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the rate of received packets by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the rate of received packets by a pod in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPackets10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PacketsPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Received Packets", panelOpts...)
}

func KubernetesReceivedPacketsDropped(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of received packets dropped by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the rate of received packets dropped by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the rate of received packets dropped by a pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload":
		description = "Shows the rate of received packets dropped by pods in a workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the rate of received packets droppedby pods in a workload in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "cluster-networking":
		description = "Shows the rate of received packets dropped by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the rate of received packets dropped by pods in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the rate of received packets dropped by top pods in a workload in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the rate of received packets dropped by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the rate of received packets dropped by a pod in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesReceivedPacketsDropped10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PacketsPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Received Packets Dropped", panelOpts...)
}

func KubernetesTransmittedPackets(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of transmitted packets by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the rate of transmitted packets by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the rate of transmitted packets by a pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload":
		description = "Shows the rate of transmitted packets by pods in a workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the rate of transmitted packets by pods in a workload in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "cluster-networking":
		description = "Shows the rate of transmitted packets by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the rate of transmitted packets by pods in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the rate of transmitted packets by top pods in a workload in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the rate of transmitted packets by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the rate of transmitted packets by a pod in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPackets10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PacketsPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Transmitted Packets", panelOpts...)
}

func KubernetesTransmittedPacketsDropped(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of transmitted packets dropped by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the rate of transmitted packets dropped by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the rate of transmitted packets droppedby a pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "workload":
		description = "Shows the rate of transmitted packets dropped by pods in a workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the rate of transmitted packets dropped by pods in a workload in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "cluster-networking":
		description = "Shows the rate of transmitted packets dropped by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the rate of transmitted packets dropped by pods in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the rate of transmitted packets dropped by top pods in a workload in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the rate of transmitted packets dropped by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the rate of transmitted packets dropped by a pod in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesTransmittedPacketsDropped10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PacketsPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Transmitted Packets Dropped", panelOpts...)
}

func KubernetesCurrentRateOfBytesReceived(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster-networking":
		description = "Shows the rate of bytes received by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesReceived1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the rate of bytes received by top pods in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesReceived2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the rate of bytes received by top workload in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesReceived3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the rate of bytes received by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesReceived4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the rate of bytes received by a pod in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesReceived5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Current Rate of Bytes Received", panelOpts...)
}

func KubernetesCurrentRateOfBytesTransmitted(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster-networking":
		description = "Shows the rate of bytes transmitted by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesTransmitted1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod-networking":
		description = "Shows the rate of bytes transmitted by top pods in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesTransmitted2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-workload-networking":
		description = "Shows the rate of bytes transmitted by top workload in a namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesTransmitted3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the rate of bytes transmitted by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesTransmitted4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod-networking":
		description = "Shows the rate of bytes transmitted by a pod in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCurrentRateOfBytesTransmitted5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Current Rate of Bytes Transmitted", panelOpts...)
}

func KubernetesAverageRateOfBytesReceived(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster-networking":
		description = "Shows the average rate of bytes received by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAverageRateOfBytesReceived1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}

	case "workload-networking":
		description = "Shows the average rate of bytes received by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAverageRateOfBytesReceived2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Average Rate of Bytes Received", panelOpts...)
}

func KubernetesAverageRateOfBytesTransmitted(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster-networking":
		description = "Shows the average rate of bytes transmitted by namespace in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAverageRateOfBytesTransmitted1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}

	case "workload-networking":
		description = "Shows the average rate of bytes transmitted by top pods in a workload in a cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesAverageRateOfBytesTransmitted2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Average Rate of Bytes Received", panelOpts...)
}
