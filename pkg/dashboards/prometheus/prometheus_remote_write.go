package prometheus

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	panels "github.com/perses/community-mixins/pkg/panels/prometheus"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/model/labels"
)

func withPrometheusRwTimestamps(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Timestamps",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusRemoteStorageTimestampLag(datasource, labelMatcher),
		panels.PrometheusRemoteStorageRateLag(datasource, labelMatcher),
	)
}

func withPrometheusRwSamples(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Samples",
		panelgroup.PanelsPerLine(1),
		panels.PrometheusRemoteStorageSampleRate(datasource, labelMatcher),
	)
}

func withPrometheusRwShard(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Shards",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusRemoteStorageCurrentShards(datasource, labelMatcher),
		panels.PrometheusRemoteStorageDesiredShards(datasource, labelMatcher),
		panels.PrometheusRemoteStorageMaxShards(datasource, labelMatcher),
		panels.PrometheusRemoteStorageMinShards(datasource, labelMatcher),
	)
}

func withPrometheusRwShardDetails(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Shard Details",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusRemoteStorageShardCapacity(datasource, labelMatcher),
		panels.PrometheusRemoteStoragePendingSamples(datasource, labelMatcher),
	)
}

func withPrometheusRwSegments(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Segments",
		panelgroup.PanelsPerLine(2),
		panels.PrometheusTSDBCurrentSegment(datasource, labelMatcher),
		panels.PrometheusRemoteWriteCurrentSegment(datasource, labelMatcher),
	)
}

func withPrometheusRwMiscRates(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Misc. Rates",
		panelgroup.PanelsPerLine(4),
		panels.PrometheusRemoteStorageDroppedSamplesRate(datasource, labelMatcher),
		panels.PrometheusRemoteStorageFailedSamplesRate(datasource, labelMatcher),
		panels.PrometheusRemoteStorageRetriedSamplesRate(datasource, labelMatcher),
		panels.PrometheusRemoteStorageEnqueueRetriesRate(datasource, labelMatcher),
	)
}

func BuildPrometheusRemoteWrite(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("prometheus-remote-write",
			dashboard.Name("Prometheus / Remote Write"),
			dashboard.ProjectName(project),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "prometheus_remote_storage_shards"),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(vector.WithMetricName("prometheus_remote_storage_shards")),
								[]*labels.Matcher{clusterLabelMatcher},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("instance"),
				),
			),
			dashboard.AddVariable("url",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("url",
						labelValuesVar.Matchers(
							promql.SetLabelMatchersV2(
								vector.New(
									vector.WithMetricName("prometheus_remote_storage_shards"),
									vector.WithLabelMatchers(
										label.New("instance").Equal("$instance"),
									),
								),
								[]*labels.Matcher{clusterLabelMatcher},
							).Pretty(0),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("url"),
				),
			),
			withPrometheusRwTimestamps(datasource, clusterLabelMatcher),
			withPrometheusRwSamples(datasource, clusterLabelMatcher),
			withPrometheusRwShard(datasource, clusterLabelMatcher),
			withPrometheusRwShardDetails(datasource, clusterLabelMatcher),
			withPrometheusRwSegments(datasource, clusterLabelMatcher),
			withPrometheusRwMiscRates(datasource, clusterLabelMatcher),
		),
	).Component("prometheus")
}
