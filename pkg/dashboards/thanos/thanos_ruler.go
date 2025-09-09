package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
)

func withThanosRuleGroupEvaluationGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rule Group Evaluations",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.RuleEvaluationRate(datasource, labelMatcher),
		panels.RuleEvaluationFailureRate(datasource, labelMatcher),
		panels.RuleGroupEvaluationsMissRate(datasource, labelMatcher),
		panels.RuleGroupEvaluationsTooSlow(datasource, labelMatcher),
	)
}

func withThanosAlertsSentGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alerts Sent",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.AlertsSentRate(datasource, labelMatcher),
		panels.AlertsDroppedRate(datasource, labelMatcher),
		panels.AlertSendingErrors(datasource, labelMatcher),
		panels.AlertSendingDurations(datasource, labelMatcher),
	)
}

func withThanosAlertQueueGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alert Queue",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.AlertQueuePushedRate(datasource, labelMatcher),
		panels.AlertQueuePoppedRate(datasource, labelMatcher),
		panels.DroppedRatio(datasource, labelMatcher),
	)
}

func BuildThanosRulerOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcherV2 := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("thanos-ruler-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Thanos / Ruler / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("thanos_build_info{container=\"thanos-rule\"}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("job"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "thanos_build_info"),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers("thanos_status"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
			withThanosRuleGroupEvaluationGroup(datasource, clusterLabelMatcherV2),
			withThanosAlertsSentGroup(datasource, clusterLabelMatcherV2),
			withThanosAlertQueueGroup(datasource, clusterLabelMatcherV2),
			withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcherV2),
			withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcherV2),
			withThanosResourcesGroup(datasource, clusterLabelMatcherV2),
		),
	).Component("thanos")
}
