package blackbox

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"

	labelValuesVar "github.com/perses/perses/go-sdk/prometheus/variable/label-values"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
)

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
	)
}
