package prometheus

import (
	"github.com/perses/community-dashboards/pkg/promql"
	promqlbuilder "github.com/perses/promql-builder"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/matrix"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var PrometheusCommonPanelQueries = map[string]parser.Expr{
	"PrometheusStatsTable": promqlbuilder.Count(
		vector.New(
			vector.WithMetricName("prometheus_build_info"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp("$job"),
				label.New("instance").EqualRegexp("$instance"),
			),
		),
	).By("job", "instance", "version"),
	"PrometheusTargetSync": promql.SumByRate(
		"prometheus_target_sync_length_seconds_sum",
		[]string{"job", "scrape_job", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PrometheusTargets": promqlbuilder.Sum(
		vector.New(
			vector.WithMetricName("prometheus_sd_discovered_targets"),
			vector.WithLabelMatchers(
				label.New("job").EqualRegexp("$job"),
				label.New("instance").EqualRegexp("$instance"),
			),
		),
	).By("job", "instance"),
	"PrometheusAverageScrapeIntervalDuration": promqlbuilder.Div(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_target_interval_length_seconds_sum"),
					vector.WithLabelMatchers(
						label.New("job").EqualRegexp("$job"),
						label.New("instance").EqualRegexp("$instance"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_target_interval_length_seconds_count"),
					vector.WithLabelMatchers(
						label.New("job").EqualRegexp("$job"),
						label.New("instance").EqualRegexp("$instance"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	),
	"PrometheusScrapeFailures_exceededBodySizeLimit": promql.SumByRate(
		"prometheus_target_scrapes_exceeded_body_size_limit_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PrometheusScrapeFailures_exceededSampleLimit": promql.SumByRate(
		"prometheus_target_scrapes_exceeded_sample_limit_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PrometheusScrapeFailures_duplicateTimestamp": promql.SumByRate(
		"prometheus_target_scrapes_sample_duplicate_timestamp_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PrometheusScrapeFailures_outOfBounds": promql.SumByRate(
		"prometheus_target_scrapes_sample_out_of_bounds_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PrometheusScrapeFailures_outOfOrder": promql.SumByRate(
		"prometheus_target_scrapes_sample_out_of_order_total",
		[]string{"job", "instance"},
		label.New("job").EqualRegexp("$job"),
		label.New("instance").EqualRegexp("$instance"),
	),
	"PrometheusAppendedSamples": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("prometheus_tsdb_head_samples_appended_total"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp("$job"),
					label.New("instance").EqualRegexp("$instance"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"PrometheusHeadSeries": vector.New(
		vector.WithMetricName("prometheus_tsdb_head_series"),
		vector.WithLabelMatchers(
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PrometheusHeadChunks": vector.New(
		vector.WithMetricName("prometheus_tsdb_head_chunks"),
		vector.WithLabelMatchers(
			label.New("job").EqualRegexp("$job"),
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PrometheusQueryRate": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("prometheus_engine_query_duration_seconds_count"),
				vector.WithLabelMatchers(
					label.New("job").EqualRegexp("$job"),
					label.New("instance").EqualRegexp("$instance"),
					label.New("slice").Equal("inner_eval"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
	"PrometheusQueryStateDuration": promqlbuilder.Max(
		vector.New(
			vector.WithMetricName("prometheus_engine_query_duration_seconds"),
			vector.WithLabelMatchers(
				label.New("quantile").Equal("0.9"),
				label.New("job").EqualRegexp("$job"),
				label.New("instance").EqualRegexp("$instance"),
			),
		),
	).By("slice"),
	"PrometheusRemoteStorageTimestampLag": promqlbuilder.Sub(
		vector.New(
			vector.WithMetricName("prometheus_remote_storage_highest_timestamp_in_seconds"),
			vector.WithLabelMatchers(
				label.New("instance").EqualRegexp("$instance"),
			),
		),
		promqlbuilder.Neq(
			vector.New(
				vector.WithMetricName("prometheus_remote_storage_queue_highest_sent_timestamp_seconds"),
				vector.WithLabelMatchers(
					label.New("instance").EqualRegexp("$instance"),
					label.New("url").Equal("$url"),
				),
			),
			&parser.NumberLiteral{
				Val: 0,
			},
		),
	).Ignoring("remote_name", "url").GroupRight("instance"),
	"PrometheusRemoteStorageRateLag": promqlbuilder.ClampMin(
		promqlbuilder.Sub(
			promqlbuilder.Rate(
				matrix.New(
					vector.New(
						vector.WithMetricName("prometheus_remote_storage_highest_timestamp_in_seconds"),
						vector.WithLabelMatchers(
							label.New("instance").EqualRegexp("$instance"),
						),
					),
					matrix.WithRangeAsVariable("$__rate_interval"),
				),
			),
			promqlbuilder.Rate(
				matrix.New(
					vector.New(
						vector.WithMetricName("prometheus_remote_storage_queue_highest_sent_timestamp_seconds"),
						vector.WithLabelMatchers(
							label.New("instance").EqualRegexp("$instance"),
							label.New("url").Equal("$url"),
						),
					),
					matrix.WithRangeAsVariable("$__rate_interval"),
				),
			),
		).Ignoring("remote_name", "url").GroupRight("instance"),
		0,
	),
	"PrometheusRemoteStorageSampleRate": promqlbuilder.Sub(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_samples_in_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
		promqlbuilder.Sub(
			promqlbuilder.Or(
				promqlbuilder.Rate(
					matrix.New(
						vector.New(
							vector.WithMetricName("prometheus_remote_storage_succeeded_samples_total"),
							vector.WithLabelMatchers(
								label.New("instance").EqualRegexp("$instance"),
								label.New("url").Equal("$url"),
							),
						),
						matrix.WithRangeAsVariable("$__rate_interval"),
					),
				),
				promqlbuilder.Rate(
					matrix.New(
						vector.New(
							vector.WithMetricName("prometheus_remote_storage_samples_total"),
							vector.WithLabelMatchers(
								label.New("instance").EqualRegexp("$instance"),
								label.New("url").Equal("$url"),
							),
						),
						matrix.WithRangeAsVariable("$__rate_interval"),
					),
				),
			),
			promqlbuilder.Or(
				promqlbuilder.Rate(
					matrix.New(
						vector.New(
							vector.WithMetricName("prometheus_remote_storage_dropped_samples_total"),
							vector.WithLabelMatchers(
								label.New("instance").EqualRegexp("$instance"),
								label.New("url").Equal("$url"),
							),
						),
						matrix.WithRangeAsVariable("$__rate_interval"),
					),
				),
				promqlbuilder.Rate(
					matrix.New(
						vector.New(
							vector.WithMetricName("prometheus_remote_storage_samples_dropped_total"),
							vector.WithLabelMatchers(
								label.New("instance").EqualRegexp("$instance"),
								label.New("url").Equal("$url"),
							),
						),
						matrix.WithRangeAsVariable("$__rate_interval"),
					),
				),
			),
		),
	).Ignoring("remote_name", "url").GroupRight("instance"),
	"PrometheusRemoteStorageCurrentShards": vector.New(
		vector.WithMetricName("prometheus_remote_storage_shards"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
			label.New("url").Equal("$url"),
		),
	),
	"PrometheusRemoteStorageDesiredShards": vector.New(
		vector.WithMetricName("prometheus_remote_storage_shards_desired"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
			label.New("url").Equal("$url"),
		),
	),
	"PrometheusRemoteStorageMaxShards": vector.New(
		vector.WithMetricName("prometheus_remote_storage_shards_max"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
			label.New("url").Equal("$url"),
		),
	),
	"PrometheusRemoteStorageMinShards": vector.New(
		vector.WithMetricName("prometheus_remote_storage_shards_min"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
			label.New("url").Equal("$url"),
		),
	),
	"PrometheusRemoteStorageShardCapacity": vector.New(
		vector.WithMetricName("prometheus_remote_storage_shard_capacity"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
			label.New("url").Equal("$url"),
		),
	),
	"PrometheusRemoteStoragePendingSamples": promqlbuilder.Or(
		vector.New(
			vector.WithMetricName("prometheus_remote_storage_pending_samples"),
			vector.WithLabelMatchers(
				label.New("instance").EqualRegexp("$instance"),
				label.New("url").Equal("$url"),
			),
		),
		vector.New(
			vector.WithMetricName("prometheus_remote_storage_samples_pending"),
			vector.WithLabelMatchers(
				label.New("instance").EqualRegexp("$instance"),
				label.New("url").Equal("$url"),
			),
		),
	),
	"PrometheusTSDBCurrentSegment": vector.New(
		vector.WithMetricName("prometheus_tsdb_wal_segment_current"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PrometheusRemoteWriteCurrentSegment": vector.New(
		vector.WithMetricName("prometheus_wal_watcher_current_segment"),
		vector.WithLabelMatchers(
			label.New("instance").EqualRegexp("$instance"),
		),
	),
	"PrometheusRemoteStorageDroppedSamplesRate": promqlbuilder.Or(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_dropped_samples_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
						label.New("url").Equal("$url"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_samples_dropped_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
						label.New("url").Equal("$url"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	),
	"PrometheusRemoteStorageFailedSamplesRate": promqlbuilder.Or(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_failed_samples_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
						label.New("url").Equal("$url"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_samples_failed_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
						label.New("url").Equal("$url"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	),
	"PrometheusRemoteStorageRetriedSamplesRate": promqlbuilder.Or(
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_retried_samples_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
						label.New("url").Equal("$url"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
		promqlbuilder.Rate(
			matrix.New(
				vector.New(
					vector.WithMetricName("prometheus_remote_storage_samples_retried_total"),
					vector.WithLabelMatchers(
						label.New("instance").EqualRegexp("$instance"),
						label.New("url").EqualRegexp("$url"),
					),
				),
				matrix.WithRangeAsVariable("$__rate_interval"),
			),
		),
	),
	"PrometheusRemoteStorageEnqueueRetriesRate": promqlbuilder.Rate(
		matrix.New(
			vector.New(
				vector.WithMetricName("prometheus_remote_storage_enqueue_retries_total"),
				vector.WithLabelMatchers(
					label.New("instance").EqualRegexp("$instance"),
					label.New("url").EqualRegexp("$url"),
				),
			),
			matrix.WithRangeAsVariable("$__rate_interval"),
		),
	),
}

// OverridePrometheusPanelQueries overrides the PrometheusCommonPanelQueries global.
// Refer to panel queries in the map, that you'd like to override.
// The convention of naming followed, is to use Panel function name (with _suffix, in case panel has multiple queries)
func OverridePrometheusPanelQueries(queries map[string]parser.Expr) {
	for k, v := range queries {
		PrometheusCommonPanelQueries[k] = v
	}
}
