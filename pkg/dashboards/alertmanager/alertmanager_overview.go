// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package alertmanager

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/alertmanager"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func withAlertsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Alerts",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.Alerts(datasource, labelMatcher),
		panels.AlertsReceiveRate(datasource, labelMatcher),
	)
}

func withNotificationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Notifications",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.NotificationsSendRate(datasource, labelMatcher),
		panels.NotificationDuration(datasource, labelMatcher),
	)
}

func BuildAlertManagerOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
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
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("alertmanager_notifications_total")),
								[]*labels.Matcher{clusterLabelMatcher, {Name: "job", Type: labels.MatchEqual, Value: "$job"}},
							).Pretty(0),
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
