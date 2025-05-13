package controller_manager

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panelsGostats "github.com/perses/community-dashboards/pkg/panels/gostats"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withCMStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Controller Manager Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.ControllerManagerUpStatus(datasource, labelMatcher),
	)
}

func withCMWorkQueueGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Work Queue",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.WorkQueueAddRate(datasource, labelMatcher),
		panels.WorkQueueDepth(datasource, labelMatcher),
		panels.WorkQueueLatency(datasource, labelMatcher),
	)
}

func withCMKubeAPIRequestsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Kube API Requests",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.KubeAPIRequestRate(datasource, labelMatcher),
		panels.PostRequestLatency(datasource, labelMatcher),
		panels.GetRequestLatency(datasource, labelMatcher),
	)
}

func withCMResources(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []promql.LabelMatcher{
		promql.ClusterVar,
		promql.InstanceVar,
		{
			Name:  "job",
			Value: "kube-controller-manager",
			Type:  "=",
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, clusterLabelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panelsGostats.MemoryUsage(datasource, labelMatchersToUse...),
		panels.KubeletCPU(datasource, labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, labelMatchersToUse...),
	)
}

func BuildControllerManagerOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("controller-manager-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Controller Manager"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("up{"+panels.GetControllerManagerMatcher()+"}"),
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
								"up{"+panels.GetControllerManagerMatcher()+"}",
								[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("instance"),
				),
			),
			withCMStatsGroup(datasource, clusterLabelMatcher),
			withCMWorkQueueGroup(datasource, clusterLabelMatcher),
			withCMKubeAPIRequestsGroup(datasource, clusterLabelMatcher),
			withCMResources(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
