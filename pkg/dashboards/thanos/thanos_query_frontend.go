package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/prometheus/prometheus/model/labels"

	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/thanos"
)

func withThanosQueryFrontendRequestsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Frontend API",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.QueryFrontendRequestRate(datasource, labelMatcher),
		panels.QueryFrontendQueryRate(datasource, labelMatcher),
		panels.QueryFrontendErrors(datasource, labelMatcher),
		panels.QueryFrontendDurations(datasource, labelMatcher),
	)
}

func withThanosQueryFrontendCacheGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Query Frontend Cache Operations",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.QueryFrontendCacheRequestRate(datasource, labelMatcher),
		panels.QueryFrontendCacheHitRate(datasource, labelMatcher),
		panels.QueryFrontendCacheMissRate(datasource, labelMatcher),
		panels.QueryFrontendFetchedKeyRate(datasource, labelMatcher),
	)
}

func BuildThanosQueryFrontendOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcherV2 := dashboards.GetClusterLabelMatcherV2(clusterLabelName)

	return dashboards.NewDashboardResult(
		dashboard.New("thanos-query-frontend-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Thanos / Query Frontend / Overview"),
			dashboard.AddVariable("job",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("job",
						labelValuesVar.Matchers("thanos_build_info{container=\"thanos-query-frontend\"}"),
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
			withThanosQueryFrontendRequestsGroup(datasource, clusterLabelMatcherV2),
			withThanosQueryFrontendCacheGroup(datasource, clusterLabelMatcherV2),
			withThanosResourcesGroup(datasource, clusterLabelMatcherV2),
		),
	).Component("thanos")
}
