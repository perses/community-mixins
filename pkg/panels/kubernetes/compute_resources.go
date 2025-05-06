package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

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
						"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m) by (cluster)",
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
						"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\"}) by (namespace)",
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
		}
	case "namespace-pod":
		description = "Shows the CPU usage of the namespace by pod, and the CPU resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", namespace=\"$namespace\"}) by (pod)",
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
						"sum(\n  node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", namespace=\"$namespace\"}\n* on(namespace,pod)\n  group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload_type=~\"$type\"}\n) by (workload, workload_type)\n",
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
						"sum(\n    node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{cluster=\"$cluster\", namespace=\"$namespace\"}\n  * on(namespace,pod)\n    group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=\"$workload\", workload_type=~\"$type\"}\n) by (pod)\n",
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
						"sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m{namespace=\"$namespace\", pod=\"$pod\", cluster=\"$cluster\", container!=\"\"}) by (container)",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_requests{job=\"kube-state-metrics\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"cpu\"}\n)\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_limits{job=\"kube-state-metrics\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"cpu\"}\n)\n",
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
						"sum(container_memory_rss{job=\"cadvisor\", container!=\"\"}) by (cluster)",
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
						"sum(container_memory_rss{job=\"cadvisor\", cluster=\"$cluster\", container!=\"\"}) by (namespace)",
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
		}
	case "node-without-cache":
		description = "Shows the memory usage (RSS) of the node by pod, and the memory capacity of the node."
		queries = []panel.Option{
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
		}
	case "namespace-pod":
		description = "Shows the memory usage of the namespace by pod, and the memory resource quota of the namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(container_memory_working_set_bytes{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", container!=\"\", image!=\"\"}) by (pod)",
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
						"sum(\n    container_memory_working_set_bytes{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", container!=\"\", image!=\"\"}\n  * on(namespace,pod)\n    group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload_type=~\"$type\"}\n) by (workload, workload_type)\n",
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
						"sum(container_memory_working_set_bytes{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", container!=\"\", image!=\"\"}) by (container)",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_requests{job=\"kube-state-metrics\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"memory\"}\n)\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("requests"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(\n    kube_pod_container_resource_limits{job=\"kube-state-metrics\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\", resource=\"memory\"}\n)\n",
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
