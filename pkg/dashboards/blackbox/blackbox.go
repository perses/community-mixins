package blackbox

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/blackbox"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

func withBlackboxSummary(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Summary",
		panelgroup.PanelsPerLine(1),
		panels.ProbeStatusMap(datasource, labelMatcher),
	)
}

func withBlackboxProbesStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Stats",
		panelgroup.PanelsPerLine(4),
		panels.ProbeSuccessCount(datasource, labelMatcher),
		panels.ProbeSuccessPercent(datasource, labelMatcher),
		panels.ProbeHTTPSSL(datasource, labelMatcher),
		panels.ProbeAverageDuration(datasource, labelMatcher),
	)
}

func withBlackboxProbesUptime(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Uptimes Stats",
		panelgroup.PanelsPerLine(2),
		panels.ProbeUptimeSuccess(datasource, labelMatcher),
		panels.ProbeUptimeMonthly(datasource, labelMatcher),
	)
}

func withBlackboxProbes(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes",
		panelgroup.PanelsPerLine(2),
		panels.ProbeDurationSeconds(datasource, labelMatcher),
		panels.ProbePhases(datasource, labelMatcher),
	)
}

func withBlackboxProbesAdditionalStats(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Additional Stats",
		panelgroup.PanelsPerLine(5),
		panels.ProbeStatusCode(datasource, labelMatcher),
		panels.ProbeTLSVersion(datasource, labelMatcher),
		panels.ProbeSSLExpiry(datasource, labelMatcher),
		panels.ProbeRedirects(datasource, labelMatcher),
		panels.ProbeHTTPVersion(datasource, labelMatcher),
	)
}

func withBlackboxProbesAvgTime(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes Avg Duration Stats",
		panelgroup.PanelsPerLine(2),
		panels.ProbeAverageDurationInstance(datasource, labelMatcher),
		panels.ProbeAverageDNSLookupPerInstance(datasource, labelMatcher),
	)
}

func BuildBlackboxExporter(project string, datasource string, clusterLabelName string) (dashboard.Builder, error) {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboard.New("blackbox-overview",
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
						promql.SetLabelMatchers(
							"probe_success",
							[]promql.LabelMatcher{clusterLabelMatcher, {Name: "job", Type: "=", Value: "$job"}},
						),
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
	)
}
