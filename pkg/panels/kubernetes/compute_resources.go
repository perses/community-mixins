package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	statPanel "github.com/perses/perses/go-sdk/panel/stat"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func KubernetesCPUUtilizationStat(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		panelName = "CPU Utilization"
		description = "Shows the CPU utilization of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(cluster:node_cpu:ratio_rate5m) / count(cluster:node_cpu:ratio_rate5m)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"cluster:node_cpu:ratio_rate5m{cluster=\"$cluster\"}",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}) / sum(kube_pod_container_resource_requests{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", resource=\"cpu\"})",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}) / sum(kube_pod_container_resource_limits{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", resource=\"cpu\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesCPURequestsCommitmentStat(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the CPU requests commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(kube_pod_container_resource_requests{"+GetKubeStateMetricsMatcher()+", resource=\"cpu\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+", resource=\"cpu\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the CPU requests commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(namespace_cpu:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+",resource=\"cpu\",cluster=\"$cluster\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("CPU Requests Commitment", panelOpts...)
}

func KubernetesCPULimitsCommitmentStat(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the CPU limits commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(kube_pod_container_resource_limits{"+GetKubeStateMetricsMatcher()+", resource=\"cpu\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+", resource=\"cpu\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the CPU limits commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(namespace_cpu:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+",resource=\"cpu\",cluster=\"$cluster\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("CPU Limits Commitment", panelOpts...)
}

func KubernetesMemoryUtilizationStat(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		panelName = "Memory Utilization"
		description = "Shows the Memory utilization of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"1 - sum(:node_memory_MemAvailable_bytes:sum) / sum(node_memory_MemTotal_bytes{"+GetNodeExporterMatcher()+"})",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"1 - sum(:node_memory_MemAvailable_bytes:sum{cluster=\"$cluster\"}) / sum(node_memory_MemTotal_bytes{"+GetNodeExporterMatcher()+",cluster=\"$cluster\"})",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\", image!=\"\"}) / sum(kube_pod_container_resource_requests{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", resource=\"memory\"})",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\",container!=\"\", image!=\"\"}) / sum(kube_pod_container_resource_limits{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", resource=\"memory\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesMemoryRequestsCommitmentStat(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the Memory requests commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(kube_pod_container_resource_requests{"+GetKubeStateMetricsMatcher()+", resource=\"memory\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+", resource=\"memory\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the Memory requests commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(namespace_memory:kube_pod_container_resource_requests:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+",resource=\"memory\",cluster=\"$cluster\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Memory Requests Commitment", panelOpts...)
}

func KubernetesMemoryLimitsCommitmentStat(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "multicluster":
		description = "Shows the Memory limits commitment of all clusters."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(kube_pod_container_resource_limits{"+GetKubeStateMetricsMatcher()+", resource=\"memory\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+", resource=\"memory\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	case "cluster":
		description = "Shows the Memory limits commitment of the cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(namespace_memory:kube_pod_container_resource_limits:sum{cluster=\"$cluster\"}) / sum(kube_node_status_allocatable{"+GetKubeStateMetricsMatcher()+",resource=\"memory\",cluster=\"$cluster\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		statPanel.Chart(
			statPanel.Format(commonSdk.Format{
				Unit:          string(commonSdk.PercentUnit),
				DecimalPlaces: 2,
			}),
			statPanel.ValueFontSize(50),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Memory Limits Commitment", panelOpts...)
}

func KubernetesCPUUsage(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "multicluster":
		description = "Shows the CPU usage of each cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+") by (cluster)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\"}) by (namespace)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(kube_node_status_capacity{cluster=\"$cluster\", "+GetKubeStateMetricsMatcher()+", node=~\"$node\", resource=\"cpu\"})",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("max capacity"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", node=~\"$node\"}) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=\"requests.cpu\"}))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=\"limits.cpu\"}))",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(\n  "+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}\n* on(namespace,pod)\n  group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload_type=~\"$type\"}\n) by (workload, workload_type)\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}} - {{workload_type}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=\"requests.cpu\"}))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=\"limits.cpu\"}))",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(\n    "+GetNodeNSCPUSecondsRecordingRule()+"{cluster=\"$cluster\", namespace=\"$namespace\"}\n  * on(namespace,pod)\n    group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=\"$workload\", workload_type=~\"$type\"}\n) by (pod)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum("+GetNodeNSCPUSecondsRecordingRule()+"{namespace=\"$namespace\", pod=\"$pod\", cluster=\"$cluster\", container!=\"\"}) by (container)",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_requests{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"cpu\"}\n)\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_limits{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"cpu\"}\n)\n",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("CPU Usage", panelOpts...)
}

func KubernetesMemoryUsage(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "multicluster":
		description = "Shows memory usage w/o cache, for each cluster."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(container_memory_rss{"+GetCAdvisorMatcher()+", container!=\"\"}) by (cluster)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(container_memory_rss{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", container!=\"\"}) by (namespace)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(kube_node_status_capacity{cluster=\"$cluster\", "+GetKubeStateMetricsMatcher()+", node=~\"$node\", resource=\"memory\"})",
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
		}
	case "node-without-cache":
		description = "Shows the memory usage (RSS) of the node by pod, and the memory capacity of the node."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(kube_node_status_capacity{cluster=\"$cluster\", "+GetKubeStateMetricsMatcher()+", node=~\"$node\", resource=\"memory\"})",
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
		}
	case "namespace-pod":
		description = "Shows the memory usage of the namespace by pod, and the memory resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", container!=\"\", image!=\"\"}) by (pod)",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=\"requests.memory\"}))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=\"limits.memory\"}))",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(\n    container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", container!=\"\", image!=\"\"}\n  * on(namespace,pod)\n    group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload_type=~\"$type\"}\n) by (workload, workload_type)\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}} - {{workload_type}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=~\"requests.memory|memory\"}))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("quota - requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"scalar(max(kube_resourcequota{cluster=\"$cluster\", namespace=\"$namespace\", type=\"hard\",resource=~\"limits.memory\"}))",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(\n    container_memory_working_set_bytes{cluster=\"$cluster\", namespace=\"$namespace\", container!=\"\", image!=\"\"}\n  * on(namespace,pod)\n    group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=\"$workload\", workload_type=~\"$type\"}\n) by (pod)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(container_memory_working_set_bytes{"+GetCAdvisorMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", container!=\"\", image!=\"\"}) by (container)",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_requests{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"memory\"}\n)\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_limits{"+GetKubeStateMetricsMatcher()+", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"memory\"}\n)\n",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Memory Usage", panelOpts...)
}
