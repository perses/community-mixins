package scheduler

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panelsGostats "github.com/perses/community-mixins/pkg/panels/gostats"
	panels "github.com/perses/community-mixins/pkg/panels/kubernetes"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"
)

func withSchedulerStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Scheduler Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.SchedulerUpStatus(datasource, labelMatcher),
	)
}

func withSchedulingRateGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Scheduling Rate",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.SchedulingRate(datasource, labelMatcher),
		panels.SchedulingLatency(datasource, labelMatcher),
	)
}

func withSchedulerKubeAPIRequestsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	labelMatchersToUse := []promql.LabelMatcher{
		promql.ClusterVar,
		promql.InstanceVar,
		{
			Name:  "job",
			Value: panels.KUBE_SCHEDULER_LABEL_VALUE,
			Type:  "=",
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, labelMatcher)

	return dashboard.AddPanelGroup("Kube API Requests",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.KubeAPIRequestRate(datasource, labelMatchersToUse...),
		panels.PostRequestLatency(datasource, labelMatchersToUse...),
		panels.GetRequestLatency(datasource, labelMatchersToUse...),
	)
}

func withSchedulerResources(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []*labels.Matcher{
		promql.ClusterVarV2,
		promql.InstanceVarV2,
		{
			Name:  "job",
			Value: panels.KUBE_SCHEDULER_LABEL_VALUE,
			Type:  labels.MatchEqual,
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, clusterLabelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panelsGostats.MemoryUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.CPUUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "instance", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "instance", labelMatchersToUse...),
	)
}

func BuildSchedulerOverview(project string, datasource string, clusterLabelName string, variableOverrides ...dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{"+panels.GetSchedulerMatcher()+"}"),
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
							"up{"+panels.GetSchedulerMatcher()+"}",
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
		dashboard.Name("Kubernetes / Scheduler"),
	}, vars...)
	options = append(options,
		withSchedulerStatsGroup(datasource, clusterLabelMatcher),
		withSchedulingRateGroup(datasource, clusterLabelMatcher),
		withSchedulerKubeAPIRequestsGroup(datasource, clusterLabelMatcher),
		withSchedulerResources(datasource, clusterLabelMatcherV2),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("scheduler-overview", options...),
	).Component("kubernetes")
}
