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
	statPanel "github.com/perses/plugins/statchart/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func KubernetesCPUUtilizationStat(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		panelName = "CPU Utilization"
		description = "Shows the CPU utilization of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUtilizationStatAll"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		panelName = "CPU Utilization"
		description = "Shows the CPU utilization of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUtilizationStatCluster"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "namespace-requests":
		panelName = "CPU Utilization (from requests)"
		description = "Shows the CPU utilization of the namespace from pod CPU requests."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUtilizationStatNSPod"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "namespace-limits":
		panelName = "CPU Utilization (from limits)"
		description = "Shows the CPU utilization of the namespace from pod CPU limits."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUtilizationStatNSPodLimits"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          &dashboards.PercentUnit,
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesCPURequestsCommitmentStat(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the CPU requests commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPURequestsCommitmentStatAllClusters"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the CPU requests commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPURequestsCommitmentStatReqClusters"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          &dashboards.PercentUnit,
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("CPU Requests Commitment", panelOpts...)
}

func KubernetesCPULimitsCommitmentStat(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the CPU limits commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPULimitsCommitmentStatAllClusters"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the CPU limits commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPULimitsCommitmentStatReqClusters"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          &dashboards.PercentUnit,
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("CPU Limits Commitment", panelOpts...)
}

func KubernetesMemoryUtilizationStat(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		panelName = "Memory Utilization"
		description = "Shows the Memory utilization of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUtilizationStatMultiCluster"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		panelName = "Memory Utilization"
		description = "Shows the Memory utilization of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUtilizationStatCluster"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "namespace-requests":
		panelName = "Memory Utilization (from requests)"
		description = "Shows the Memory utilization of the namespace from pod memory requests."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUtilizationNSRequests"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "namespace-limits":
		panelName = "Memory Utilization (from limits)"
		description = "Shows the Memory utilization of the namespace from pod memory limits."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUtilizationNSLimits"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          &dashboards.PercentUnit,
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesMemoryRequestsCommitmentStat(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the Memory requests commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryRequestsCommitmentStat1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the Memory requests commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryRequestsCommitmentStat2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          &dashboards.PercentUnit,
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Memory Requests Commitment", panelOpts...)
}

func KubernetesMemoryLimitsCommitmentStat(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the Memory limits commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryLimitsCommitmentStat1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the Memory limits commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryLimitsCommitmentStat2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          &dashboards.PercentUnit,
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Memory Limits Commitment", panelOpts...)
}

func KubernetesCPUUsage(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "multicluster":
		description = "Shows the CPU usage of each cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{cluster}}"),
				),
			),
		}
	case "cluster":
		description = "Shows the CPU usage of the cluster by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "node":
		description = "Shows the CPU usage of the node by pod, and the CPU capacity of the node."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("max capacity"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the CPU usage of the namespace by pod, and the CPU resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - limits"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the CPU usage of the namespace by workload, and the CPU resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}} - {{workload_type}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - limits"),
				),
			),
		}
	case "workload":
		description = "Shows the CPU usage of the workload (deployment, statefulset, job, cronjob, daemonset, etc.) by pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage11"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the CPU usage of the pod by container, alongwith the requests and limits."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage12"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage13"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesCPUUsage14"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("limits"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("CPU Usage", panelOpts...)
}

func KubernetesMemoryUsage(granularity, datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "multicluster":
		description = "Shows memory usage w/o cache, for each cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage1"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{cluster}}"),
				),
			),
		}
	case "cluster":
		description = "Shows the memory usage of the cluster by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage2"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "node-with-cache":
		description = "Shows the memory usage of the node by pod, and the memory capacity of the node."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage3"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("max capacity"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage4"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "node-without-cache":
		description = "Shows the memory usage (RSS) of the node by pod, and the memory capacity of the node."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage5"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("max capacity"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage6"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "namespace-pod":
		description = "Shows the memory usage of the namespace by pod, and the memory resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage7"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage8"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage9"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - limits"),
				),
			),
		}
	case "namespace-workload":
		description = "Shows the memory usage of the namespace by workload, and the memory resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage10"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}} - {{workload_type}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage11"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage12"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - limits"),
				),
			),
		}
	case "workload":
		description = "Shows the memory usage of the workload (deployment, statefulset, job, cronjob, daemonset, etc.) by pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage13"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		description = "Shows the memory usage (WSS) of the pod by container, alongwith the requests and limits."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage14"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage15"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchersV2(
						KubernetesCommonPanelQueries["KubernetesMemoryUsage16"],
						labelMatchers,
					).Pretty(0),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("limits"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesUnit,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Memory Usage", panelOpts...)
}
