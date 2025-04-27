package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"

	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/thanos"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withThanosRuleGroupEvaluationGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rule Group Evaluations",
		panelgroup.PanelsPerLine(4),
		panels.RuleEvaluationRate(datasource, labelMatcher),
		panels.RuleEvaluationFailures(datasource, labelMatcher),
		panels.RuleGroupEvaluationsMissed(datasource, labelMatcher),
		panels.RuleGroupEvaluationsTooSlow(datasource, labelMatcher),
	)
}

func withThanosAlertsSentGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alerts Sent",
		panelgroup.PanelsPerLine(4),
		panels.AlertsSent(datasource, labelMatcher),
		panels.AlertsDropped(datasource, labelMatcher),
		panels.AlertSendingErrors(datasource, labelMatcher),
		panels.AlertSendingDurations(datasource, labelMatcher),
	)
}

func withThanosAlertQueueGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alert Queue",
		panelgroup.PanelsPerLine(3),
		panels.AlertQueuePushed(datasource, labelMatcher),
		panels.AlertQueuePopped(datasource, labelMatcher),
		panels.DroppedRatio(datasource, labelMatcher),
	)
}

func BuildThanosRulerOverview(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("thanos-ruler-overview",
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
		withThanosRuleGroupEvaluationGroup(datasource, clusterLabelMatcher),
		withThanosAlertsSentGroup(datasource, clusterLabelMatcher),
		withThanosAlertQueueGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCUnaryGroup(datasource, clusterLabelMatcher),
		withThanosReadGRPCStreamGroup(datasource, clusterLabelMatcher),
	)
}
