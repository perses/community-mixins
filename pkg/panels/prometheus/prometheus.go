package prometheus

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	"github.com/prometheus/prometheus/model/labels"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

// PrometheusStatsTable creates a panel option for displaying Prometheus statistics.
//
// The panel uses the following Prometheus metrics:
// - prometheus_build_info: Build information about Prometheus instances
//
// The panel shows:
// - Instance count by job and version
// - Version information per instance
func PrometheusStatsTable(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Prometheus Stats",
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "job",
					Header: "Job",
				},
				{
					Name:   "instance",
					Header: "Instance",
				},
				{
					Name:   "version",
					Header: "Version",
				},
				{
					Name: "value",
					Hide: true,
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusStatsTable"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

// PrometheusTargetSync creates a panel option for monitoring Prometheus target synchronization.
//
// The panel uses the following Prometheus metrics:
// - prometheus_target_sync_length_seconds_sum: Total time taken for target synchronization
//
// The panel shows:
// - Target synchronization time per job and instance
// - Rate of synchronization over 5-minute intervals
//
// Parameters:
//   - datasourceName: The name of the data source.
//   - labelMatchers: A variadic parameter for label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusTargetSync(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Target Sync",
		panel.Description("Monitors target synchronization time for Prometheus instances"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.SecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusTargetSync"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - Metrics"),
			),
		),
	)
}

// PrometheusTargets creates a panel group option for displaying Prometheus targets.
//
// The panel uses the following Prometheus metrics:
// - prometheus_sd_discovered_targets: Number of targets discovered by service discovery
//
// The panel shows:
// - Total number of discovered targets per job
// - Breakdown by instance
//
// Parameters:
//   - datasourceName: The name of the Prometheus datasource.
//   - labelMatchers: Optional variadic parameter for PromQL label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusTargets(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Targets",
		panel.Description("Shows discovered targets across Prometheus instances"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusTargets"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - Metrics"),
			),
		),
	)
}

// PrometheusAverageScrapeIntervalDuration creates a panel option for displaying the average scrape interval duration
// for Prometheus targets. It uses the following Prometheus metrics:
// - prometheus_target_interval_length_seconds_sum: Sum of all scrape interval lengths
// - prometheus_target_interval_length_seconds_count: Count of scrape intervals
//
// The panel shows:
// - Average duration between scrapes for each target
// - Breakdown by job and instance
//
// Parameters:
//   - datasourceName: The name of the Prometheus datasource.
//   - labelMatchers: Optional PromQL label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusAverageScrapeIntervalDuration(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Scrape Interval Duration",
		panel.Description("Shows average interval between scrapes for Prometheus targets"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.SecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusAverageScrapeIntervalDuration"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - {{interval}} Configured"),
			),
		),
	)
}

// PrometheusScrapeFailures creates a panel group option for displaying Prometheus scrape failure metrics.
//
// The panel uses the following Prometheus metrics:
// - prometheus_target_scrapes_exceeded_body_size_limit_total: Number of times a scrape exceeded the body size limit
// - prometheus_target_scrapes_exceeded_sample_limit_total: Number of times a scrape exceeded the sample limit
// - prometheus_target_scrapes_sample_duplicate_timestamp_total: Number of times a scrape had duplicate timestamps
// - prometheus_target_scrapes_sample_out_of_bounds_total: Number of times a scrape had samples out of bounds
// - prometheus_target_scrapes_sample_out_of_order_total: Number of times a scrape had samples out of order
//
// The panel shows:
// - Different types of scrape failures
// - Rate of failures per type and target
func PrometheusScrapeFailures(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Scrape failures",
		panel.Description("Shows scrape failure metrics for Prometheus targets"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusScrapeFailures_exceededBodySizeLimit"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("exceeded body size limit: {{job}} - {{instance}} - Metrics"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusScrapeFailures_exceededSampleLimit"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("exceeded sample limit: {{job}} - {{instance}} - Metrics"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusScrapeFailures_duplicateTimestamp"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("duplicate timestamp: {{job}} - {{instance}} - Metrics"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusScrapeFailures_outOfBounds"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("out of bounds: {{job}} - {{instance}} - Metrics"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusScrapeFailures_outOfOrder"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("out of order: {{job}} - {{instance}} - Metrics"),
			),
		),
	)
}

// PrometheusAppendedSamples creates a panel option for displaying sample append rate.
//
// The panel uses the following Prometheus metrics:
// - prometheus_tsdb_head_samples_appended_total: Total samples appended to TSDB head
//
// The panel shows:
// - Rate of samples being appended
// - Breakdown by job and instance
func PrometheusAppendedSamples(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Appended Samples",
		panel.Description("Shows rate of samples appended to Prometheus TSDB"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusAppendedSamples"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - {{remote_name}} - {{url}}"),
			),
		),
	)
}

// PrometheusHeadSeries creates a panel option for displaying the head series metric from Prometheus.
// The panel uses the following Prometheus metrics:
// - prometheus_tsdb_head_series: Number of series in the head block
//
// The panel shows:
// - Current number of active series in TSDB head
// - Breakdown by job and instance
func PrometheusHeadSeries(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Head Series",
		panel.Description("Shows number of series in Prometheus TSDB head"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusHeadSeries"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - Head Series"),
			),
		),
	)
}

// PrometheusHeadChunks creates a panel option for displaying the "Head Chunks" metric from Prometheus.
//
// The panel uses the following Prometheus metrics:
// - prometheus_tsdb_head_chunks: Number of chunks in the head block
//
// The panel shows:
// - Current number of chunks in TSDB head
// - Breakdown by job and instance
func PrometheusHeadChunks(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Head Chunks",
		panel.Description("Shows number of chunks in Prometheus TSDB head"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusHeadChunks"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - Head Chunks"),
			),
		),
	)
}

// PrometheusQueryRate creates a panel option for displaying the query rate metrics.
// The panel uses the following Prometheus metrics:
// - prometheus_engine_query_duration_seconds_count: Number of queries executed
//
// The panel shows:
// - Query execution rate over time
// - Breakdown by job and instance
//
// Parameters:
//   - datasourceName: The name of the data source.
//   - labelMatchers: A variadic parameter for label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusQueryRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Query Rate",
		panel.Description("Shows Prometheus query rate metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusQueryRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} - {{instance}} - Query Rate"),
			),
		),
	)
}

// PrometheusQueryStateDuration creates a panel option for displaying the stage duration
// of Prometheus queries.
//
// The panel uses the following Prometheus metrics:
// - prometheus_engine_query_duration_seconds: Duration of query execution stages
//
// The panel shows:
// - Duration of different query stages
// - 90th percentile of query times
//
// Parameters:
//   - datasourceName: The name of the data source.
//   - labelMatchers: A variadic parameter for label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusQueryStateDuration(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Stage Duration",
		panel.Description("Shows duration of different Prometheus query stages"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.SecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusQueryStateDuration"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{slice}} - Duration"),
			),
		),
	)
}

// PrometheusRemoteStorageTimestampLag creates a panel option for visualizing the timestamp lag
// between the highest timestamp in Prometheus remote storage and the highest sent
// timestamp in the remote storage queue.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_highest_timestamp_in_seconds: Highest timestamp in remote storage
// - prometheus_remote_storage_queue_highest_sent_timestamp_seconds: Highest sent timestamp
//
// The panel shows:
// - Lag between storage and queue timestamps
// - Breakdown by remote storage target
func PrometheusRemoteStorageTimestampLag(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Timestamp Lag",
		panel.Description("Shows timestamp lag in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageTimestampLag"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Segment"),
			),
		),
	)
}

// PrometheusRemoteStorageRateLag creates a panel option for monitoring the rate lag of Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_highest_timestamp_in_seconds: Highest timestamp in remote storage
// - prometheus_remote_storage_queue_highest_sent_timestamp_seconds: Highest sent timestamp
//
// The panel shows:
// - Rate of lag between storage and queue timestamps
// - 5-minute rate changes per target
func PrometheusRemoteStorageRateLag(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate",
		panel.Description("Shows rate metrics over 5 minute intervals"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageRateLag"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageSampleRate creates a panel option for visualizing the rate of Prometheus remote storage samples
// over a 5-minute interval. It displays the rate of incoming samples versus succeeded or dropped samples.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_samples_in_total: Total samples received
// - prometheus_remote_storage_succeeded_samples_total: Successfully stored samples
// - prometheus_remote_storage_dropped_samples_total: Dropped samples
//
// The panel shows:
// - Rate of sample ingestion
// - Success vs drop rates
//
// Parameters:
//   - datasourceName: The name of the data source.
//   - labelMatchers: A variadic parameter for label matchers.
//
// Returns:
//   - panelgroup.Option: The configured panel option.
func PrometheusRemoteStorageSampleRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate, in vs. succeeded or dropped",
		panel.Description("Shows rate of samples in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageSampleRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageCurrentShards creates a panel option for displaying the current number of shards
// in Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_shards: Current number of shards per remote storage
//
// The panel shows:
// - Current shard count per target
// - Breakdown by instance and URL
func PrometheusRemoteStorageCurrentShards(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Current Shards",
		panel.Description("Shows current number of shards in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageCurrentShards"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageDesiredShards creates a panel option for displaying the desired shards
// of Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_shards_desired: Desired number of shards per remote storage
//
// The panel shows:
// - Target shard count per remote storage
// - Configuration vs actual shards
func PrometheusRemoteStorageDesiredShards(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Desired Shards",
		panel.Description("Shows desired number of shards in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageDesiredShards"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageMaxShards creates a panel option for displaying the maximum number of shards
// in Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_shards_max: Maximum allowed shards per remote storage
//
// The panel shows:
// - Maximum shard limit per target
// - Upper bounds for scaling
func PrometheusRemoteStorageMaxShards(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Max Shards",
		panel.Description("Shows maximum number of shards in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageMaxShards"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageMinShards creates a panel option for displaying the minimum number of shards
// in Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_shards_min: Minimum required shards per remote storage
//
// The panel shows:
// - Minimum shard requirement per target
// - Lower bounds for scaling
func PrometheusRemoteStorageMinShards(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Min Shards",
		panel.Description("Shows minimum number of shards in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageMinShards"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageShardCapacity creates a panel option for displaying the shard capacity
// in Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_shard_capacity: Current capacity of remote storage shards
//
// The panel shows:
// - Shard capacity per remote storage target
// - Breakdown by instance and URL
func PrometheusRemoteStorageShardCapacity(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Shard Capacity",
		panel.Description("Shows shard capacity in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageShardCapacity"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStoragePendingSamples creates a panel option for displaying the pending samples
// in Prometheus remote storage.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_pending_samples: Number of samples pending in remote storage
// - prometheus_remote_storage_samples_pending: Legacy metric for pending samples
//
// The panel shows:
// - Number of samples waiting to be sent
// - Breakdown by remote storage target
func PrometheusRemoteStoragePendingSamples(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Pending Samples",
		panel.Description("Shows number of pending samples in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStoragePendingSamples"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusTSDBCurrentSegment creates a panel option for displaying the current segment
// of the Prometheus TSDB WAL (Write-Ahead Log).
//
// The panel uses the following Prometheus metrics:
// - prometheus_tsdb_wal_segment_current: Current WAL segment being written to
//
// The panel shows:
// - Current WAL segment number
// - Segment progression over time
func PrometheusTSDBCurrentSegment(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("TSDB Current Segment",
		panel.Description("Shows current TSDB WAL segment"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusTSDBCurrentSegment"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Segment - Metrics"),
			),
		),
	)
}

// PrometheusRemoteWriteCurrentSegment creates a panel option for displaying the current segment
// of the Prometheus remote write WAL (Write-Ahead Log).
//
// The panel uses the following Prometheus metrics:
// - prometheus_wal_watcher_current_segment: Current segment of remote write WAL
//
// The panel shows:
// - Current remote write WAL segment
// - Segment progression over time
func PrometheusRemoteWriteCurrentSegment(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Current Segment",
		panel.Description("Shows current remote write WAL segment"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteWriteCurrentSegment"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - Segment - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageDroppedSamplesRate creates a panel option for displaying the rate of dropped samples
// in Prometheus remote storage over a 5-minute interval.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_dropped_samples_total: Total dropped samples
// - prometheus_remote_storage_samples_dropped_total: Legacy metric for dropped samples
//
// The panel shows:
// - Rate of sample drops per target
// - Drop patterns over time
func PrometheusRemoteStorageDroppedSamplesRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Dropped Samples Rate",
		panel.Description("Shows rate of dropped samples in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageDroppedSamplesRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageFailedSamplesRate creates a panel option for displaying the rate of failed samples
// in Prometheus remote storage over a 5-minute interval.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_failed_samples_total: Total failed samples
// - prometheus_remote_storage_samples_failed_total: Legacy metric for failed samples
//
// The panel shows:
// - Rate of sample failures per target
// - Failure patterns over time
func PrometheusRemoteStorageFailedSamplesRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Failed Samples",
		panel.Description("Shows rate of failed samples in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageFailedSamplesRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageRetriedSamplesRate creates a panel option for displaying the rate of retried samples
// in Prometheus remote storage over a 5-minute interval.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_retried_samples_total: Total retried samples
// - prometheus_remote_storage_samples_retried_total: Legacy metric for retried samples
//
// The panel shows:
// - Rate of sample retries per target
// - Retry patterns over time
func PrometheusRemoteStorageRetriedSamplesRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Retried Samples",
		panel.Description("Shows rate of retried samples in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageRetriedSamplesRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}

// PrometheusRemoteStorageEnqueueRetriesRate creates a panel option for displaying the rate of enqueue retries
// in Prometheus remote storage over a 5-minute interval.
//
// The panel uses the following Prometheus metrics:
// - prometheus_remote_storage_enqueue_retries_total: Total enqueue retry attempts
//
// The panel shows:
// - Rate of enqueue retries per target
// - Retry patterns over time
func PrometheusRemoteStorageEnqueueRetriesRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Enqueue Retries",
		panel.Description("Shows rate of enqueue retries in remote storage"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PrometheusCommonPanelQueries["PrometheusRemoteStorageEnqueueRetriesRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{remote_name}} - {{url}} - Metrics"),
			),
		),
	)
}
