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

func withBlackboxProbes(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Probes",
		panelgroup.PanelsPerLine(1),
		panels.ProbeDurationSeconds(datasource, labelMatcher),
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
		withBlackboxProbes(datasource, clusterLabelMatcher),
	)
}
