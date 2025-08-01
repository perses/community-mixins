package apiserver

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

func withMarkdown(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Notice",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(3),
		panels.APIServerSLONotice(datasource, labelMatcher),
	)
}

func withAllAvailabilityAndErrorBudget(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("All Availability And Error Budget",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.APIServerAvailability(datasource, labelMatcher),
		panels.APIServerErrorBudget(datasource, labelMatcher),
	)
}

func withReadStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("API Server Read",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.APIServerReadAvailability(datasource, labelMatcher),
		panels.APIServerReadSLIRequests(datasource, labelMatcher),
		panels.APIServerReadSLIErrors(datasource, labelMatcher),
		panels.APIServerReadSLIDuration(datasource, labelMatcher),
	)
}

func withWriteStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("API Server Write",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.APIServerWriteAvailability(datasource, labelMatcher),
		panels.APIServerWriteSLIRequests(datasource, labelMatcher),
		panels.APIServerWriteSLIErrors(datasource, labelMatcher),
		panels.APIServerWriteSLIDuration(datasource, labelMatcher),
	)
}

func withWorkQueueGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Work Queue",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.APIServerWorkQueueAddRate(datasource, labelMatcher),
		panels.APIServerWorkQueueDepth(datasource, labelMatcher),
		panels.APIServerWorkQueueLatency(datasource, labelMatcher),
	)
}

func withAPIServerResources(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	labelMatchersToUse := []*labels.Matcher{
		promql.ClusterVarV2,
		promql.InstanceVarV2,
		{
			Name:  "job",
			Value: "kube-apiserver",
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

func BuildAPIServerOverview(project string, datasource string, clusterLabelName string, variableOverrides ...dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{"+panels.GetAPIServerMatcher()+"}"),
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
							"up{"+panels.GetAPIServerMatcher()+"}",
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
		dashboard.Name("Kubernetes / API server"),
	}, vars...)
	options = append(options,
		withMarkdown(datasource, clusterLabelMatcher),
		withAllAvailabilityAndErrorBudget(datasource, clusterLabelMatcher),
		withReadStats(datasource, clusterLabelMatcher),
		withWriteStats(datasource, clusterLabelMatcher),
		withWorkQueueGroup(datasource, clusterLabelMatcher),
		withAPIServerResources(datasource, clusterLabelMatcherV2),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("api-server-overview", options...),
	).Component("kubernetes")
}
