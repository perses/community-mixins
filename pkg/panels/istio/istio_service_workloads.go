package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

// ========== SERVICE WORKLOAD PANELS ==========

func IncomingRequestsByService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Requests By Destination Workload And Response Code",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{Min: 0}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_requests_total{connection_security_policy=\"mutual_tls\",destination_service=~\"$service\",reporter=\"destination\",destination_workload=~\"$dstwl\",destination_workload_namespace=~\"$dstns\"}[5m])) by (destination_workload, destination_workload_namespace, response_code), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} : {{ response_code }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_requests_total{connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", reporter=\"destination\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[5m])) by (destination_workload, destination_workload_namespace, response_code), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} : {{ response_code }}"),
			),
		),
	)
}

func IncomingSuccessRateByService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Success Rate (non-5xx responses) By Destination Workload",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: string(commonSdk.PercentDecimalUnit)},
				Min:    0,
				Max:    1.01,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_requests_total{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\",response_code!~\"5.*\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[5m])) by (destination_workload, destination_workload_namespace) / sum(irate(istio_requests_total{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[5m])) by (destination_workload, destination_workload_namespace)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_requests_total{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\",response_code!~\"5.*\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[5m])) by (destination_workload, destination_workload_namespace) / sum(irate(istio_requests_total{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[5m])) by (destination_workload, destination_workload_namespace)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func IncomingRequestDurationByService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Duration By Service Workload",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: string(commonSdk.SecondsUnit)},
				Min:    0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
			}),
		),
		// mTLS P50, P90, P95, P99
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.95, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.95, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P99 (🔐mTLS)"),
			),
		),
		// Non-mTLS P50, P90, P95, P99
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.95, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.95, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P99"),
			),
		),
	)
}

func IncomingRequestSizeByService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Size By Service Workload",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: string(commonSdk.BytesUnit)},
				Min:    0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
			}),
		),
		// mTLS P50, P90, P95, P99
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.50, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.90, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}  P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.95, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.99, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}  P99 (🔐mTLS)"),
			),
		),
		// Non-mTLS P50, P90, P95, P99
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.50, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.90, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.95, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.99, sum(irate(istio_request_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P99"),
			),
		),
	)
}

func ResponseSizeByService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Response Size By Service Workload",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: string(commonSdk.BytesUnit)},
				Min:    0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
			}),
		),
		// mTLS P50, P90, P95, P99
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.50, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.90, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}  P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.95, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.99, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}  P99 (🔐mTLS)"),
			),
		),
		// Non-mTLS P50, P90, P95, P99
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.50, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.90, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.95, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"histogram_quantile(0.99, sum(irate(istio_response_bytes_bucket{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[1m])) by (destination_workload, destination_workload_namespace, le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }} P99"),
			),
		),
	)
}

func BytesReceivedFromTCPService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes Received from Incoming TCP Connection",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: string(commonSdk.BytesPerSecondsUnit)},
				Min:    0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_tcp_received_bytes_total{reporter=\"destination\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace}} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_tcp_received_bytes_total{reporter=\"destination\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace}}"),
			),
		),
	)
}

func BytesSentToTCPService(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes Sent to Incoming TCP Connection",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: string(commonSdk.BytesPerSecondsUnit)},
				Min:    0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_tcp_sent_bytes_total{connection_security_policy=\"mutual_tls\", reporter=\"destination\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{destination_workload_namespace }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_tcp_sent_bytes_total{connection_security_policy!=\"mutual_tls\", reporter=\"destination\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{destination_workload_namespace }}"),
			),
		),
	)
}
