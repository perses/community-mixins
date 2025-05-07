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

func KubernetesReceiveBandwidth(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "cluster":
		description = "Shows the network receive bandwidth of the cluster by namespace."
		queries = []panel.Option{
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
		}
	case "cluster-networking":
		description = "Shows the network receive bandwidth of the cluster highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by (namespace) (\n    rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_receive_bytes_total{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (pod) (\n    rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the network receive bandwidth of the workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=~\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(irate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace=~\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
	}

	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Receive Bandwidth", panelOpts...)
}

func KubernetesTransmitBandwidth(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var queries []panel.Option
	var description string

	switch granularity {
	case "cluster":
		description = "Shows the network transmit bandwidth of the cluster by namespace."
		queries = []panel.Option{
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
		}
	case "cluster-networking":
		description = "Shows the network transmit bandwidth of the cluster highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by (namespace) (\n    rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_bytes_total{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (pod) (\n    rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{workload}}"),
				),
			),
		}
	case "workload-networking":
		description = "Shows the network transmit bandwidth of the workload highlighting top pods."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=~\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace=~\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
	}

	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Transmit Bandwidth", panelOpts...)
}

func KubernetesAvgContainerBandwidthTransmitted(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		panelName = "Average Container Bandwidth by Namespace: Transmitted"
		description = "Shows the average network bandwidth transmitted in container by namespace."
		queries = []panel.Option{
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
		}
	case "namespace-workload":
		panelName = "Average Container Bandwidth by Workload: Transmitted"
		description = "Shows the average network bandwidth transmitted in container by workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"(avg(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(avg(rate(container_network_transmit_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(avg(rate(container_network_transmit_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesAvgContainerBandwidthReceived(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		panelName = "Average Container Bandwidth by Namespace: Received"
		description = "Shows the average network bandwidth received in container by namespace."
		queries = []panel.Option{
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
		}
	case "namespace-workload":
		panelName = "Average Container Bandwidth by Workload: Received"
		description = "Shows the average network bandwidth received in container by workload."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"(avg(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(avg(rate(container_network_receive_bytes_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(avg(rate(container_network_receive_bytes_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesReceivedPackets(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of received packets by namespace."
		queries = []panel.Option{
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
		}
	case "namespace-pod":
		description = "Shows the rate of received packets by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(irate(container_network_receive_packets_total{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_receive_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_receive_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_receive_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (namespace) (\n    rate(container_network_receive_packets_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (pod) (\n    rate(container_network_receive_packets_total{cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_receive_packets_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_receive_packets_total{job=\"cadvisor\", cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=~\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_receive_packets_total{cluster=\"$cluster\",namespace=~\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Received Packets", panelOpts...)
}

func KubernetesReceivedPacketsDropped(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of received packets dropped by namespace."
		queries = []panel.Option{
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
		}
	case "namespace-pod":
		description = "Shows the rate of received packets dropped by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(irate(container_network_receive_packets_dropped_total{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_receive_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_receive_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_receive_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (namespace) (\n    rate(container_network_receive_packets_dropped_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (pod) (\n    rate(container_network_receive_packets_dropped_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_receive_packets_dropped_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_receive_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=~\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_receive_packets_dropped_total{cluster=\"$cluster\",namespace=~\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Received Packets Dropped", panelOpts...)
}

func KubernetesTransmittedPackets(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of transmitted packets by namespace."
		queries = []panel.Option{
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
		}
	case "namespace-pod":
		description = "Shows the rate of transmitted packets by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(irate(container_network_transmit_packets_total{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_transmit_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_transmit_packets_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (namespace) (\n    rate(container_network_transmit_packets_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (pod) (\n    rate(container_network_transmit_packets_total{cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_transmit_packets_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_transmit_packets_total{job=\"cadvisor\", cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=~\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_packets_total{cluster=\"$cluster\",namespace=~\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Transmitted Packets", panelOpts...)
}

func KubernetesTransmittedPacketsDropped(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		description = "Shows the rate of transmitted packets dropped by namespace."
		queries = []panel.Option{
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
		}
	case "namespace-pod":
		description = "Shows the rate of transmitted packets dropped by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum(irate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\", namespace=\"$namespace\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_transmit_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"(sum(rate(container_network_transmit_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\", namespace=\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\", namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (namespace) (\n    rate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\",namespace!=\"\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum by (pod) (\n    rate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n  * on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n)\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\",namespace=\"$namespace\"}[5m])\n* on (cluster,namespace,pod) group_left ()\n    topk by (cluster,namespace,pod) (\n      1,\n      max by (cluster,namespace,pod) (kube_pod_info{host_network=\"false\"})\n    )\n* on (cluster,namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=\"$namespace\", workload=~\".+\", workload_type=~\"$type\"}) by (workload))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sort_desc(sum(rate(container_network_transmit_packets_dropped_total{job=\"cadvisor\", cluster=\"$cluster\",namespace=~\"$namespace\"}[5m])\n* on (namespace,pod)\ngroup_left(workload,workload_type) namespace_workload_pod:kube_pod_owner:relabel{cluster=\"$cluster\",namespace=~\"$namespace\", workload=~\"$workload\", workload_type=~\"$type\"}) by (pod))\n",
						labelMatchers,
					),
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
					promql.SetLabelMatchers(
						"sum(rate(container_network_transmit_packets_dropped_total{cluster=\"$cluster\",namespace=~\"$namespace\", pod=~\"$pod\"}[5m])) by (pod)",
						labelMatchers,
					),
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
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel("Rate of Transmitted Packets Dropped", panelOpts...)
}
