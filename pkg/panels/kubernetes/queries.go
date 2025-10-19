package kubernetes

import (
	"maps"

	"github.com/perses/community-mixins/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var KubernetesCommonPanelQueries = map[string]parser.Expr{
	// apiserver panels
	"APIServerAvailability": vector.New(
		vector.WithMetricName("apiserver_request:availability30d"),
		vector.WithLabelMatchers(
			label.New("verb").Equal("all"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerErrorBudget": promqlbuilder.Mul(
		&parser.NumberLiteral{Val: 100},
		promqlbuilder.Sub(
			vector.New(
				vector.WithMetricName("apiserver_request:availability30d"),
				vector.WithLabelMatchers(
					label.New("verb").Equal("all"),
					label.New("cluster").EqualRegexp("$cluster"),
				),
			),
			&parser.NumberLiteral{Val: 0.990000},
		),
	),
	"APIServerReadAvailability": vector.New(
		vector.WithMetricName("apiserver_request:availability30d"),
		vector.WithLabelMatchers(
			label.New("verb").Equal("read"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerReadSLIRequests": promql.SumBy(
		"code_resource:apiserver_request_total:rate5m",
		[]string{"code"},
		label.New("verb").Equal("read"),
		label.New("cluster").EqualRegexp("$cluster"),
	),
	"APIServerReadSLIErrors": promqlbuilder.Div(
		promql.SumBy(
			"code_resource:apiserver_request_total:rate5m",
			[]string{"resource"},
			label.New("verb").Equal("read"),
			label.New("code").EqualRegexp("5.."),
			label.New("cluster").EqualRegexp("$cluster"),
		),
		promql.SumBy(
			"code_resource:apiserver_request_total:rate5m",
			[]string{"resource"},
			label.New("verb").Equal("read"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerReadSLIDuration": vector.New(
		vector.WithMetricName("cluster_quantile:apiserver_request_sli_duration_seconds:histogram_quantile"),
		vector.WithLabelMatchers(
			label.New("verb").Equal("read"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerWriteAvailability": vector.New(
		vector.WithMetricName("apiserver_request:availability30d"),
		vector.WithLabelMatchers(
			label.New("verb").Equal("write"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerWriteSLIRequests": promql.SumBy(
		"code_resource:apiserver_request_total:rate5m",
		[]string{"code"},
		label.New("verb").Equal("write"),
		label.New("cluster").EqualRegexp("$cluster"),
	),
	"APIServerWriteSLIErrors": promqlbuilder.Div(
		promql.SumBy(
			"code_resource:apiserver_request_total:rate5m",
			[]string{"resource"},
			label.New("verb").Equal("write"),
			label.New("code").EqualRegexp("5.."),
			label.New("cluster").EqualRegexp("$cluster"),
		),
		promql.SumBy(
			"code_resource:apiserver_request_total:rate5m",
			[]string{"resource"},
			label.New("verb").Equal("write"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerWriteSLIDuration": vector.New(
		vector.WithMetricName("cluster_quantile:apiserver_request_sli_duration_seconds:histogram_quantile"),
		vector.WithLabelMatchers(
			label.New("verb").Equal("write"),
			label.New("cluster").EqualRegexp("$cluster"),
		),
	),
	"APIServerWorkQueueAddRate": promql.SumByRate(
		"workqueue_adds_total",
		[]string{"instance", "name"},
		label.New("cluster").EqualRegexp("$cluster"),
		label.New("instance").EqualRegexp("$instance'"),
		label.New("job").Equal("kube-apiserver"),
	),
	"APIServerWorkQueueDepth": promql.SumByRate(
		"workqueue_depth",
		[]string{"instance", "name"},
		label.New("cluster").EqualRegexp("$cluster"),
		label.New("instance").EqualRegexp("$instance'"),
		label.New("job").Equal("kube-apiserver"),
	),
	"APIServerWorkQueueLatency": promqlbuilder.HistogramQuantile(
		0.99,
		promql.SumByRate(
			"workqueue_queue_duration_seconds_bucket",
			[]string{"instance", "name", "le"},
			label.New("cluster").EqualRegexp("$cluster"),
			label.New("instance").EqualRegexp("$instance'"),
			label.New("job").Equal("kube-apiserver"),
		),
	),
	// cluster
	"ClusterCPUUsageQuotaPodOwn": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("kube_pod_owner"),
			vector.WithLabelMatchers(
				label.New("job").Equal("kube-state-metrics'"),
				label.New("cluster").Equal("$cluster"),
			),
		),
	).By("namespace"),
	"ClusterCPUUsageQuotaNSwWorkload": promqlbuilder.Count(
		promql.AvgBy(
			"namespace_workload_pod:kube_pod_owner:relabel",
			[]string{"workload", "namespace"},
			label.New("job").Equal("kube-state-metrics'"),
			label.New("cluster").Equal("$cluster"),
		),
	).By("namespace"),
	"ClusterCPUUsageQuotaNodeNSCPU": promql.SumBy(
		"node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate",
		[]string{"namespace"},
		label.New("cluster").Equal("$cluster"),
	),
	"ClusterCPUUsageQuotaNSKubePodContainer": promql.SumBy(
		"namespace_cpu:kube_pod_container_resource_requests:sum",
		[]string{"namespace"},
		label.New("cluster").Equal("$cluster"),
	),
	"ClusterCPUUsageQuotaNSKubePodContainerDiv": promqlbuilder.Div(
		promql.SumBy(
			"namespace_cpu:kube_pod_container_resource_requests:sum",
			[]string{"namespace"},
			label.New("cluster").Equal("$cluster"),
		),
		promql.SumBy(
			"namespace_cpu:kube_pod_container_resource_requests:sum",
			[]string{"namespace"},
			label.New("cluster").Equal("$cluster"),
		),
	),
	"ClusterCPUUsageQuotaNSCPUKubePodResources": promql.SumBy(
		"namespace_cpu:kube_pod_container_resource_limits:sum",
		[]string{"namespace"},
		label.New("cluster").Equal("$cluster"),
	),
	"ClusterCPUUsageQuotaNSCPUKubePodResourcesDiv": promqlbuilder.Div(
		promql.SumBy(
			"node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m",
			[]string{"namespace"},
			label.New("cluster").Equal("$cluster"),
		),
		promql.SumBy(
			"namespace_cpu:kube_pod_container_resource_limits:sum",
			[]string{"namespace"},
			label.New("cluster").Equal("$cluster"),
		),
	),
	"ClusterMemoryUsageQuotaPodOwner": promql.SumBy(
		"kube_pod_owner",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
	),
	"ClusterMemoryUsageQuotaNSWorkloadPodOwner": promqlbuilder.Count(
		promql.AvgBy(
			"namespace_workload_pod:kube_pod_owner:relabel",
			[]string{"workload", "namespace"},
			label.New("cluster").Equal("$cluster"),
		),
	).By("namespace"),
	"ClusterMemoryUsageQuotaContainerMem": promql.SumBy(
		"container_memory_rss",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("cluster").NotEqualRegexp(""),
	),
	"ClusterMemoryUsageQuotaContainerResourceReqSum": promql.SumBy(
		"namespace_memory:kube_pod_container_resource_requests:sum",
		[]string{"namespace"},
		label.New("cluster").Equal("$cluster"),
	),
	"ClusterMemoryUsageQuotaContainerResourceReqSumDiv": promqlbuilder.Div(
		promql.SumBy(
			"container_memory_rss",
			[]string{"namespace"},
			label.New("job").Equal("kube-state-metrics"),
			label.New("cluster").Equal("$cluster"),
			label.New("cluster").NotEqualRegexp(""),
		),
		promql.SumBy(
			"namespace_memory:kube_pod_container_resource_requests:sum",
			[]string{"namespace"},
			label.New("cluster").Equal("$cluster"),
		),
	),
	"ClusterMemoryUsageQuotaContainerReqLimits": promql.SumBy(
		"namespace_memory:kube_pod_container_resource_limits:sum",
		[]string{"namespace"},
		label.New("cluster").Equal("$cluster"),
	),
	"ClusterMemoryUsageQuotaContainerReqLimitsDiv": promqlbuilder.Div(
		promql.SumBy(
			"container_memory_rss",
			[]string{"namespace"},
			label.New("job").Equal("kube-state-metrics"),
			label.New("cluster").Equal("$cluster"),
			label.New("cluster").NotEqualRegexp(""),
		),
		promql.SumBy(
			"namespace_memory:kube_pod_container_resource_limits:sum",
			[]string{"namespace"},
			label.New("cluster").Equal("$cluster"),
		),
	),
	"ClusterCurrentNetworkUsageBytesTotal": promql.SumByRate(
		"container_network_receive_bytes_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").EqualRegexp(".+"),
	),
	"ClusterCurrentNetworkTrasmitBytesTotal": promql.SumByRate(
		"container_network_transmit_bytes_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").EqualRegexp(".+"),
	),
	"ClusterCurrentNetworkReceivedTotal": promql.SumByRate(
		"container_network_receive_packets_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").EqualRegexp(".+"),
	),
	"ClusterCurrentNetworkTransmitPacketsTotal": promql.SumByRate(
		"container_network_transmit_packets_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").EqualRegexp(".+"),
	),
	"ClusterCurrentNetworkReceivedPacketsDroppedTotal": promql.SumByRate(
		"container_network_receive_packets_dropped_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").EqualRegexp(".+"),
	),
	"ClusterCurrentNetworkTransmitPacketsDroppedTotal": promql.SumByRate(
		"container_network_transmit_packets_dropped_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").EqualRegexp(".+"),
	),
	"ClusterCurrentStorageIOFsReadsTotal": promql.SumByRate(
		"container_fs_reads_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
		label.New("container").NotEqualRegexp(""),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").NotEqual(""),
	),
	"ClusterCurrentStorageIOFsWritesTotal": promql.SumByRate(
		"container_fs_writes_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
		label.New("container").NotEqualRegexp(""),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").NotEqual(""),
	),
	"ClusterCurrentStorageIOFsReadsWritesTotal": promqlbuilder.Add(
		promql.SumByRate(
			"container_fs_reads_total",
			[]string{"namespace"},
			label.New("job").Equal("kube-state-metrics"),
			label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
			label.New("container").NotEqualRegexp(""),
			label.New("cluster").Equal("$cluster"),
			label.New("namespace").NotEqual(""),
		),
		promql.SumByRate(
			"container_fs_writes_total",
			[]string{"namespace"},
			label.New("job").Equal("kube-state-metrics"),
			label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
			label.New("container").NotEqualRegexp(""),
			label.New("cluster").Equal("$cluster"),
			label.New("namespace").NotEqual(""),
		),
	),
	"ClusterCurrentStorageIOFsReadsBytesTotal": promql.SumByRate(
		"container_fs_reads_bytes_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
		label.New("container").NotEqualRegexp(""),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").NotEqual(""),
	),
	"ClusterCurrentStorageIOFsWritesBytesTotal": promql.SumByRate(
		"container_fs_writes_bytes_total",
		[]string{"namespace"},
		label.New("job").Equal("kube-state-metrics"),
		label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
		label.New("container").NotEqualRegexp(""),
		label.New("cluster").Equal("$cluster"),
		label.New("namespace").NotEqual(""),
	),
	"ClusterCurrentStorageIOFsReadsWritesBytesTotal": promqlbuilder.Add(
		promql.SumByRate(
			"container_fs_reads_bytes_total",
			[]string{"namespace"},
			label.New("job").Equal("kube-state-metrics"),
			label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
			label.New("container").NotEqualRegexp(""),
			label.New("cluster").Equal("$cluster"),
			label.New("namespace").NotEqual(""),
		),
		promql.SumByRate(
			"container_fs_writes_bytes_total",
			[]string{"namespace"},
			label.New("job").Equal("kube-state-metrics"),
			label.New("device").EqualRegexp("(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+"),
			label.New("container").NotEqualRegexp(""),
			label.New("cluster").Equal("$cluster"),
			label.New("namespace").NotEqual(""),
		),
	),
	// compute resources panels
	"KubernetesCPUUtilizationStatAll": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("cluster:node_cpu:ratio_rate5m"),
			),
		),
		promqlbuilder.Count(
			vector.New(
				vector.WithMetricName("cluster:node_cpu:ratio_rate5m"),
			),
		),
	),
	"KubernetesCPUUtilizationStatCluster": vector.New(
		vector.WithMetricName("cluster:node_cpu:ratio_rate5m{"),
		vector.WithLabelMatchers(
			label.New("cluster").Equal("$cluster"),
		),
	),
	"KubernetesCPUUtilizationStatNSPod": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m"),
				vector.WithLabelMatchers(
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namespace"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_pod_container_resource_requests"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namespace"),
					label.New("resource").Equal("cpu"),
				),
			),
		),
	),
	"KubernetesCPUUtilizationStatNSPodLimits": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m"),
				vector.WithLabelMatchers(
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namespace"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_pod_container_resource_limits"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namespace"),
					label.New("resource").Equal("cpu"),
				),
			),
		),
	),
	"KubernetesCPURequestsCommitmentStatsAllClusters": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_pod_container_resource_requests"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("resource").Equal("cpu"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_node_status_allocatable"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("resource").Equal("cpu"),
				),
			),
		),
	),
	"KubernetesCPURequestsCommitmentStatsReqClusters": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("namespace_cpu:kube_pod_container_resource_requests:sum"),
				vector.WithLabelMatchers(
					label.New("cluster").Equal("$cluster"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_node_status_allocatable"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("resource").Equal("cpu"),
					label.New("cluster").Equal("$cluster"),
				),
			),
		),
	),
	"KubernetesCPULimitsCommitmentStatAllClusters": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_pod_container_resource_limits"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("resource").Equal("cpu"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_node_status_allocatable"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("resource").Equal("cpu"),
				),
			),
		),
	),
	"KubernetesCPULimitsCommitmentStatReqClusters": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("namespace_cpu:kube_pod_container_resource_limits:sum"),
				vector.WithLabelMatchers(
					label.New("cluster").Equal("$cluster"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_node_status_allocatable"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("resource").Equal("cpu"),
					label.New("cluster").Equal("$cluster"),
				),
			),
		),
	),
	"KubernetesMemoryUtilizationStatMultiCluster": promqlbuilder.Div(
		promqlbuilder.Sub(
			&parser.NumberLiteral{Val: 1},
			promqlbuilder.Sum(
				vector.New(
					vector.WithMetricName(":node_memory_MemAvailable_bytes:sum"),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("node_memory_MemTotal_bytes"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node-exporter"),
				),
			),
		),
	),
	"KubernetesMemoryUtilizationStatCluster": promqlbuilder.Div(
		promqlbuilder.Sub(
			&parser.NumberLiteral{Val: 1},
			promqlbuilder.Sum(
				vector.New(
					vector.WithMetricName(":node_memory_MemAvailable_bytes:sum"),
					vector.WithLabelMatchers(
						label.New("cluster").Equal("$cluster"),
					),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("node_memory_MemTotal_bytes"),
				vector.WithLabelMatchers(
					label.New("job").Equal("node-exporter"),
					label.New("cluster").Equal("$cluster"),
				),
			),
		),
	),
	"KubernetesMemoryUtilizationNSRequests": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("container_memory_working_set_bytes"),
				vector.WithLabelMatchers(
					label.New("job").Equal("cadvisor"),
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namesapce"),
					label.New("container").NotEqual(""),
					label.New("image").NotEqual(""),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_pod_container_resource_requests"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namespace"),
					label.New("resource").Equal("$memory"),
				),
			),
		),
	),
	"KubernetesMemoryUtilizationNSLimits": promqlbuilder.Div(
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("container_memory_working_set_bytes"),
				vector.WithLabelMatchers(
					label.New("job").Equal("cadvisor"),
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namesapce"),
					label.New("container").NotEqual(""),
					label.New("image").NotEqual(""),
				),
			),
		),
		promqlbuilder.Sum(
			vector.New(
				vector.WithMetricName("kube_pod_container_resource_limits"),
				vector.WithLabelMatchers(
					label.New("job").Equal("kube-state-metrics"),
					label.New("cluster").Equal("$cluster"),
					label.New("namespace").Equal("$namespace"),
					label.New("resource").Equal("$memory"),
				),
			),
		),
	),
}

// OverrideKubernetesPanelQueries overrides the KubernetesCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverrideKubernetesPanelQueries(queries map[string]parser.Expr) {
	maps.Copy(KubernetesCommonPanelQueries, queries)
}
