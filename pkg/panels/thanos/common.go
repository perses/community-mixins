package thanos

import (
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func BucketOperationRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bucket Operations",
		panel.Description("Shows rate of executions of operations against object storage bucket."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.RequestsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, operation) (rate(thanos_objstore_bucket_operations_total{namespace='$namespace', job=~'$job'}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{operation}} {{namespace}}"),
			),
		),
	)
}

func BucketOperationErrors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bucket Operation Errors",
		panel.Description("Shows percentage of errors of operations against object storage bucket."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(sum by (namespace, job, operation) (rate(thanos_objstore_bucket_operation_failures_total{namespace='$namespace', job=~'$job'}[$__rate_interval])) / sum by (namespace, job, operation) (rate(thanos_objstore_bucket_operations_total{namespace='$namespace', job=~'$job'}[$__rate_interval]))) * 100",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{operation}} {{namespace}}"),
			),
		),
	)
}

func BucketOperationDurations(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bucket Operation Latency",
		panel.Description("Shows latency of operations against object storage bucket."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, operation, le) (rate(thanos_objstore_bucket_operation_duration_seconds_bucket{namespace='$namespace', job=~'$job'}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{job}} {{operation}} {{namespace}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, operation, le) (rate(thanos_objstore_bucket_operation_duration_seconds_bucket{namespace='$namespace', job=~'$job'}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{job}} {{operation}} {{namespace}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, operation, le) (rate(thanos_objstore_bucket_operation_duration_seconds_bucket{namespace='$namespace', job=~'$job'}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{job}} {{operation}} {{namespace}}"),
			),
		),
	)
}

func ReadGRPCUnaryRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unary gRPC Read request rate",
		panel.Description("Shows rate of handled Unary gRPC Read requests (StoreAPI)."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.RequestsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_method, grpc_code) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGRPCUnaryErrors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unary gRPC Read error rate",
		panel.Description("Shows percentage of errors of Unary gRPC Read requests (StoreAPI)."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_code) (rate(grpc_server_handled_total{grpc_code=~\"Unknown|ResourceExhausted|Internal|Unavailable|DataLoss\",namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[$__rate_interval])) / ignoring (grpc_code) group_left() sum by (namespace, job) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGPRCUnaryDurations(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unary gRPC Read duration",
		panel.Description("Shows duration percentiles of handled Unary gRPC Read requests (StoreAPI)."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{namespace}} {{job}}"),
			),
		),
	)
}

func ReadGRPCStreamRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Stream gRPC Read request rate",
		panel.Description("Shows rate of handled Stream gRPC Read requests (StoreAPI Series/Exemplar calls)."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.RequestsPerSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_method, grpc_code) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGRPCStreamErrors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Stream gRPC Read error rate",
		panel.Description("Shows percentage of errors of Stream gRPC Read requests (StoreAPI Series/Exemplar calls)."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_code) (rate(grpc_server_handled_total{grpc_code=~\"Unknown|ResourceExhausted|Internal|Unavailable|DataLoss\",namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[$__rate_interval])) / ignoring (grpc_code) group_left() sum by (namespace, job) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGPRCStreamDurations(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Stream gRPC Read duration",
		panel.Description("Shows duration percentiles of handled Stream gRPC Read requests (StoreAPI Series/Exemplar calls)."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[$__rate_interval])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{namespace}} {{job}}"),
			),
		),
	)
}
