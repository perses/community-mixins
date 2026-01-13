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

package blackbox

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/blackbox"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func withBlackboxSummary(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Summary",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.ProbeStatusMap(datasource, labelMatcher),
	)
}

func withBlackboxProbesStats(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Stats",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(8),
		panels.ProbeSuccessCount(datasource, labelMatcher),
		panels.ProbeSuccessPercent(datasource, labelMatcher),
		panels.ProbeHTTPSSL(datasource, labelMatcher),
		panels.ProbeAverageDuration(datasource, labelMatcher),
	)
}

func withBlackboxProbesUptime(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Uptimes Stats",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ProbeUptimeSuccess(datasource, labelMatcher),
		panels.ProbeUptimeMonthly(datasource, labelMatcher),
	)
}

func withBlackboxProbes(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(10),
		panels.ProbeDurationSeconds(datasource, labelMatcher),
		panels.ProbePhases(datasource, labelMatcher),
	)
}

func withBlackboxProbesAdditionalStats(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Additional Stats",
		panelgroup.PanelsPerLine(5),
		panelgroup.PanelHeight(8),
		panels.ProbeStatusCode(datasource, labelMatcher),
		panels.ProbeTLSVersion(datasource, labelMatcher),
		panels.ProbeSSLExpiry(datasource, labelMatcher),
		panels.ProbeRedirects(datasource, labelMatcher),
		panels.ProbeHTTPVersion(datasource, labelMatcher),
	)
}

func withBlackboxProbesAvgTime(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Avg Duration Stats",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ProbeAverageDurationInstance(datasource, labelMatcher),
		panels.ProbeAverageDNSLookupPerInstance(datasource, labelMatcher),
	)
}

func BuildBlackboxExporter(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("blackbox-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Blackbox Exporter / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("probe_success"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("job"),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "probe_success"),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("probe_success")),
								[]*labels.Matcher{clusterLabelMatcher, {Name: "job", Type: labels.MatchEqual, Value: "$job"}},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("instance"),
				),
			),
			withBlackboxSummary(datasource, clusterLabelMatcher),
			withBlackboxProbesStats(datasource, clusterLabelMatcher),
			withBlackboxProbesUptime(datasource, clusterLabelMatcher),
			withBlackboxProbes(datasource, clusterLabelMatcher),
			withBlackboxProbesAdditionalStats(datasource, clusterLabelMatcher),
			withBlackboxProbesAvgTime(datasource, clusterLabelMatcher),
		),
	).Component("blackbox-exporter")
}
