package kubelet

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panelsGostats "github.com/perses/community-dashboards/pkg/panels/gostats"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"
)

func withKubeletStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Kubelet Stats",
		panelgroup.PanelsPerLine(6),
		panels.RunningKubeletStat(datasource, labelMatcher),
		panels.RunningPodStat(datasource, labelMatcher),
		panels.RunningContainersStat(datasource, labelMatcher),
		panels.ActVolumeCountStat(datasource, labelMatcher),
		panels.DesiredVolumeCountStat(datasource, labelMatcher),
		panels.ConfigErrorCountStat(datasource, labelMatcher),
	)
}

func withKubeletOperations(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Operation Rate and Errors",
		panelgroup.PanelsPerLine(2),
		panels.OperationRate(datasource, labelMatcher),
		panels.OperationErrorRate(datasource, labelMatcher),
	)
}

func withKubeletOperationsQuantile(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Operation Duration 99th quantile",
		panelgroup.PanelsPerLine(1),
		panels.OperationDurationQuantile(datasource, labelMatcher),
	)
}

func withPodStartRateAndDuration(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Pod Start Rate and Duration",
		panelgroup.PanelsPerLine(2),
		panels.PodStartRate(datasource, labelMatcher),
		panels.PodStartDuration(datasource, labelMatcher),
	)
}

func withStorageOperationsAndErrors(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage Operations Rate and Errors",
		panelgroup.PanelsPerLine(2),
		panels.StorageOperationRate(datasource, labelMatcher),
		panels.StorageOperationErrorRate(datasource, labelMatcher),
	)
}

func withStorageOperationsQuantile(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Storage Operation Duration 99th quantile",
		panelgroup.PanelsPerLine(1),
		panels.StorageOperationDuration(datasource, labelMatcher),
	)
}

func withCgroupManager(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Cgroup manager",
		panelgroup.PanelsPerLine(2),
		panels.CgroupManagerOperationRate(datasource, labelMatcher),
		panels.CgroupManagerQuantile(datasource, labelMatcher),
	)
}

func withPLEGRelist(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("PLEG relist",
		panelgroup.PanelsPerLine(2),
		panels.PLEGRelistRate(datasource, labelMatcher),
		panels.PLEGRelistInterval(datasource, labelMatcher),
	)
}

func withPLEGRelistDuration(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("PLEG relist duration",
		panelgroup.PanelsPerLine(1),
		panels.PLEGRelistDuration(datasource, labelMatcher),
	)
}

func withRPCRate(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("RPC rate",
		panelgroup.PanelsPerLine(1),
		panels.RPCRate(datasource, labelMatcher),
	)
}

func withRequestDurationQuantile(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Request duration 99th quantile",
		panelgroup.PanelsPerLine(1),
		panels.RequestDurationQuantile(datasource, labelMatcher),
	)
}

func withKubeletResources(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []*labels.Matcher{
		promql.ClusterVarV2,
		promql.InstanceVarV2,
		{
			Name:  "job",
			Value: "kubelet",
			Type:  labels.MatchEqual,
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, clusterLabelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(10),
		panelsGostats.MemoryUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.CPUUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "instance", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "instance", labelMatchersToUse...),
	)
}

func BuildKubeletOverview(project string, datasource string, clusterLabelName string, variableOverrides []dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{"+panels.GetKubeletMatcher()+"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("cluster"),
			),
		),
		dashboard.AddVariable("instance",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("instance",
					labelValuesVar.Matchers(
						promql.SetLabelMatchers(
							"up{"+panels.GetKubeletMatcher()+"}",
							[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
						),
					),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("instance"),
			),
		),
	}

	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	clusterLabelMatcherV2 := dashboards.GetClusterLabelMatcherV2(clusterLabelName)

	vars := defaultVars
	if len(variableOverrides) > 0 {
		vars = variableOverrides
	}
	options := append([]dashboard.Option{
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Kubelet"),
	}, vars...)
	options = append(options,
		withKubeletStats(datasource, clusterLabelMatcher),
		withKubeletOperations(datasource, clusterLabelMatcher),
		withKubeletOperationsQuantile(datasource, clusterLabelMatcher),
		withPodStartRateAndDuration(datasource, clusterLabelMatcher),
		withStorageOperationsAndErrors(datasource, clusterLabelMatcher),
		withStorageOperationsQuantile(datasource, clusterLabelMatcher),
		withCgroupManager(datasource, clusterLabelMatcher),
		withPLEGRelist(datasource, clusterLabelMatcher),
		withPLEGRelistDuration(datasource, clusterLabelMatcher),
		withRPCRate(datasource, clusterLabelMatcher),
		withRequestDurationQuantile(datasource, clusterLabelMatcher),
		withKubeletResources(datasource, clusterLabelMatcherV2),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("kubelet-overview", options...),
	).Component("kubernetes")
}
