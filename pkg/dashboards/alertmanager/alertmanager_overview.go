package alertmanager

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/alertmanager"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withAlertsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alerts",
		panelgroup.PanelsPerLine(2),
		panels.Alerts(datasource, labelMatcher),
		panels.AlertsReceiveRate(datasource, labelMatcher),
	)
}

func withNotificationsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Notifications",
		panelgroup.PanelsPerLine(2),
		panels.NotificationsSendRate(datasource, labelMatcher),
		panels.NotificationDuration(datasource, labelMatcher),
	)
}

func BuildAlertManagerOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("alertmanager-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Alertmanager / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("alertmanager_alerts"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("job"),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "alertmanager_alerts"),
			dashboard.AddVariable("integration",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("integration",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"alertmanager_notifications_total",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
					listVar.DisplayName("integration"),
				),
			),
			withAlertsGroup(datasource, clusterLabelMatcher),
			withNotificationsGroup(datasource, clusterLabelMatcher),
		),
	).Component("alertmanager")
}
