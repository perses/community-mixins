package thanos

import (
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"

	commonSdk "github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func RemoteWriteRequestRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Rate",
		panel.Description("Shows rate of incoming Remote Write v1 requests."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers("sum by (namespace, job, handler, code) (rate(http_requests_total{namespace='$namespace', job=~'$job', handler=\"receive\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{namespace}} {{handler}} {{code}}"),
			),
		),
	)
}

func RemoteWriteRequestErrors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Errors",
		panel.Description("Shows percentage of errors for incoming Remote Write v1 requests."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(sum by (namespace, job, code) (rate(http_requests_total{namespace='$namespace', job=~'$job', handler=\"receive\",code=~\"5..\"}[5m])) / ignoring (code) group_left() sum by (namespace, job) (rate(http_requests_total{namespace='$namespace', job=~'$job', handler=\"receive\"}[5m]))) * 100",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}} {{namespace}} {{code}} {{handler}}"),
			),
		),
	)
}

func RemoteWriteRequestDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Duration",
		panel.Description("Duration percentiles of successful Remote Write requests."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(http_request_duration_seconds_bucket{namespace='$namespace', job=~'$job', handler=\"receive\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{job}} - {{namespace}} duration"),
			),
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(http_request_duration_seconds_bucket{namespace='$namespace', job=~'$job', handler=\"receive\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{job}} - {{namespace}} duration"),
			),
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(http_request_duration_seconds_bucket{namespace='$namespace', job=~'$job', handler=\"receive\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{job}} {{namespace}} duration"),
			),
		),
	)
}

func TenantedRemoteWriteRequestRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Rate by tenant",
		panel.Description("Shows rate of incoming Remote Write v1 requests split by tenant."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (tenant, job, handler, code) (rate(http_requests_total{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}} {{code}} {{job}} {{namespace}} {{handler}}"),
			),
		),
	)
}

func TenantedRemoteWriteRequestErrors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Remote Write Errors by tenant",
		panel.Description("Shows percentage of errors for incoming Remote Write v1 requests split by tenant."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(sum by (tenant, namespace, job, code) (rate(http_requests_total{namespace='$namespace', job=~'$job', handler=\"receive\", code!~\"2..\", tenant=~'$tenant'}[5m])) / ignoring (code) group_left() sum by (tenant, namespace, job) (rate(http_requests_total{namespace='$namespace', job=~'$job', handler=\"receive\", tenant=~'$tenant'}[5m]))) * 100",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}} {{code}} {{job}} {{namespace}}"),
			),
		),
	)
}

func TenantedRemoteWriteRequestDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Remote Write Duration by tenant",
		panel.Description("Average duration of Remote Write requests by tenants."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, tenant) (rate(http_request_duration_seconds_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\"}[5m])) / sum by (namespace, job, tenant) (http_request_duration_seconds_count{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}} {{namespace}} {{job}}"),
			),
		),
	)
}

func AvgRemoteWriteRequestSize(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Successful Remote Write request size",
		panel.Description("Shows average size of successful remote write request."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, tenant) (rate(http_request_size_bytes_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\", code=~\"2..\"}[5m])) / sum by (namespace, job, tenant) (rate(http_request_size_bytes_count{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\", code=~\"2..\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}}"),
			),
		),
	)
}

func AvgFailedRemoteWriteRequestSize(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Average Failed Remote Write request size",
		panel.Description("Shows average size of failed remote write request."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, tenant) (rate(http_request_size_bytes_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\", code!~\"2..\"}[5m])) / sum by (namespace, job, tenant) (rate(http_request_size_bytes_count{namespace='$namespace', job=~'$job', tenant=~'$tenant', handler=\"receive\", code!~\"2..\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}}"),
			),
		),
	)
}

func InflightRemoteWriteRequests(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Inflight remote write requests",
		panel.Description("Shows inflight remote write HTTP requests."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, tenant, method) (http_inflight_requests{namespace='$namespace', job=~'$job', tenant=~'tenant', handler=\"receive\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{method}} {{tenant}}"),
			),
		),
	)
}

const (
	seriesPerSecond  commonSdk.ThroughputUnit = "series/sec"
	samplesPerSecond commonSdk.ThroughputUnit = "samples/sec"
)

func RemoteWriteSeriesRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Series ingested",
		panel.Description("Shows rate of timeseries ingested by Receive, split by tenant."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(seriesPerSecond),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(thanos_receive_write_timeseries_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', code=~\"2..\"}[5m])) by (namespace, job, tenant)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}}"),
			),
		),
	)
}

func RemoteWriteSeriesNotWrittenRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Series not ingested",
		panel.Description("Shows rate of timeseries not ingested by Receive, split by tenant."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(seriesPerSecond),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(thanos_receive_write_timeseries_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', code!~\"2..\"}[5m])) by (namespace, job, tenant)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}}"),
			),
		),
	)
}

func RemoteWriteSamplesRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Samples ingested",
		panel.Description("Shows rate of samples ingested by Receive, split by tenant."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(seriesPerSecond),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(thanos_receive_write_samples_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', code=~\"2..\"}[5m])) by (namespace, job, tenant) ",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}}"),
			),
		),
	)
}

func RemoteWriteSamplesNotWrittenRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of Samples not ingested",
		panel.Description("Shows rate of samples not ingested by Receive, split by tenant."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(seriesPerSecond),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(thanos_receive_write_samples_sum{namespace='$namespace', job=~'$job', tenant=~'$tenant', code!~\"2..\"}[5m])) by (namespace, job, tenant) ",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tenant}}"),
			),
		),
	)
}

func RemoteWriteReplicationRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of replication requests",
		panel.Description("Shows rate of replication requests between Receives."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job) (rate(thanos_receive_replications_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}"),
			),
		),
	)
}

func RemoteWriteReplicationErrorRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of replication errors",
		panel.Description("Shows rate of replication errors between Receives."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job) (rate(thanos_receive_replications_total{namespace='$namespace', job=~'$job', result=\"error\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}"),
			),
		),
	)
}

func RemoteWriteForwardRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of forward requests",
		panel.Description("Shows rate of forward requests between Receives."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job) (rate(thanos_receive_forward_requests_total{namespace='$namespace', job=~'$job'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}"),
			),
		),
	)
}

func RemoteWriteForwardErrRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Rate of forward errors",
		panel.Description("Shows rate of forward errors between Receives."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job) (rate(thanos_receive_forward_requests_total{namespace='$namespace', job=~'$job', result=\"error\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{job}}"),
			),
		),
	)
}

func WriteGRPCUnaryRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unary gRPC Write request rate",
		panel.Description("Shows rate of handled Unary gRPC Write requests (WritableStore)."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_method, grpc_code) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method=\"RemoteWrite\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func WriteGRPCUnaryErrRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unary gRPC Write error rate",
		panel.Description("Shows percentage of errors of Unary gRPC Write requests (WritableStore)."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_code) (rate(grpc_server_handled_total{grpc_code=~\"Unknown|ResourceExhausted|Internal|Unavailable|DataLoss\",namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method=\"RemoteWrite\"}[5m])) / ignoring (grpc_code) group_left() sum by (namespace, job) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method=\"RemoteWrite\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func WriteGPRCUnaryDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Unary gRPC Write duration",
		panel.Description("Shows duration percentiles of handled Unary gRPC Write requests (WritableStore)."),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method=\"RemoteWrite\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method=\"RemoteWrite\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method=\"RemoteWrite\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{namespace}} {{job}}"),
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_method, grpc_code) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGRPCUnaryErrRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_code) (rate(grpc_server_handled_total{grpc_code=~\"Unknown|ResourceExhausted|Internal|Unavailable|DataLoss\",namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[5m])) / ignoring (grpc_code) group_left() sum by (namespace, job) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGPRCUnaryDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"unary\", grpc_method!=\"RemoteWrite\"}[5m])))",
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_method, grpc_code) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGRPCStreamErrRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (namespace, job, grpc_code) (rate(grpc_server_handled_total{grpc_code=~\"Unknown|ResourceExhausted|Internal|Unavailable|DataLoss\",namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[5m])) / ignoring (grpc_code) group_left() sum by (namespace, job) (rate(grpc_server_handled_total{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{namespace}} {{job}} {{grpc_method}} {{grpc_code}}"),
			),
		),
	)
}

func ReadGPRCStreamDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
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
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{namespace}} {{job}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum by (namespace, job, le) (rate(grpc_server_handling_seconds_bucket{namespace='$namespace', job=~'$job', grpc_type=\"server_stream\"}[5m])))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{namespace}} {{job}}"),
			),
		),
	)
}

func BucketUploadTable(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Last Uploaded Block",
		panel.Description("Shows the last uploaded block time for Receive."),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "job",
					Header: "Job",
				},
				{
					Name:   "bucket",
					Header: "Bucket",
				},
				{
					Name:   "namespace",
					Header: "Namespace",
				},
				{
					Name:   "value",
					Header: "Uploaded Ago",
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"time() - max by (namespace, job, bucket) (thanos_objstore_bucket_last_successful_upload_time{namespace='$namespace', job=~'$job'})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
