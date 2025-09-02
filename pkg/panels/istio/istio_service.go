package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	statPanel "github.com/perses/plugins/statchart/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func ClientRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Client Request Volume",
		panel.Description("Request volume from client workloads to service"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.RequestsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
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
					"round(sum(irate(istio_requests_total{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func ClientSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Client Success Rate",
		panel.Description("Success rate of requests from client workloads to service"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PercentUnit,
				},
				Min: 0,
				Max: 1,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
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
					"sum(irate(istio_requests_total{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\", response_code!~\"5.*\"}[$__rate_interval])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func ClientRequestDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Client Request Duration",
		panel.Description("Request duration percentiles from client workloads to service"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.MilliSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
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
					"histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P50 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P90 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P99 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func ServerRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Server Request Volume",
		panel.Description("Request volume to service workloads"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.RequestsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
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
					"round(sum(irate(istio_requests_total{reporter=\"destination\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func ServiceTCPBytesReceived(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Bytes Received",
		panel.Description("TCP bytes received by service workloads"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
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
					"sum(irate(istio_tcp_received_bytes_total{reporter=\"destination\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func ServiceTCPBytesSent(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Bytes Sent",
		panel.Description("TCP bytes sent from service workloads"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
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
					"sum(irate(istio_tcp_sent_bytes_total{reporter=\"destination\", destination_service=~\"$service\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

// ========== STAT PANELS (for General section) ==========

func ClientRequestVolumeStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Client Request Volume",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.WithSparkline(statPanel.Sparkline{
				Width: 1,
			}),
			statPanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{Color: "green", Value: 0},
					{Color: "red", Value: 80},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_requests_total{reporter=~\"$qrep\",destination_service=~\"$service\"}[5m])), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClientSuccessRateStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Client Success Rate (non-5xx responses)",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{Unit: &dashboards.PercentDecimalUnit}),
			statPanel.WithSparkline(statPanel.Sparkline{
				Width: 1,
			}),
			statPanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{Color: "dark-red", Value: 0},
					{Color: "dark-yellow", Value: 0.95},
					{Color: "dark-green", Value: 0.99},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_requests_total{reporter=~\"$qrep\",destination_service=~\"$service\",response_code!~\"5.*\"}[5m])) / (sum(irate(istio_requests_total{reporter=~\"$qrep\",destination_service=~\"$service\"}[5m])) or on () vector(1))",
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ClientRequestDurationChart(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Client Request Duration",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: &dashboards.SecondsUnit},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\",destination_service=~\"$service\"}[1m])) by (le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\",destination_service=~\"$service\"}[1m])) by (le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\",destination_service=~\"$service\"}[1m])) by (le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\",destination_service=~\"$service\"}[1m])) by (le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\",destination_service=~\"$service\"}[1m])) by (le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\",destination_service=~\"$service\"}[1m])) by (le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P99"),
			),
		),
	)
}

func TCPReceivedBytesStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Received Bytes",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.MeanCalculation),
			statPanel.Format(commonSdk.Format{Unit: &dashboards.BytesPerSecondsUnit}),
			statPanel.WithSparkline(statPanel.Sparkline{
				Width: 1,
			}),
			statPanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{Color: "green", Value: 0},
					{Color: "red", Value: 80},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_tcp_received_bytes_total{reporter=~\"$qrep\", destination_service=~\"$service\"}[1m]))",
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ServerRequestVolumeStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Server Request Volume",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.WithSparkline(statPanel.Sparkline{
				Width: 1,
			}),
			statPanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{Color: "green", Value: 0},
					{Color: "red", Value: 80},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_requests_total{reporter=\"destination\",destination_service=~\"$service\"}[5m])), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ServerSuccessRateStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Server Success Rate (non-5xx responses)",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{Unit: &dashboards.PercentDecimalUnit}),
			statPanel.WithSparkline(statPanel.Sparkline{
				Width: 1,
			}),
			statPanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{Color: "rgba(50, 172, 45, 0.97)", Value: 0},
					{Color: "rgba(237, 129, 40, 0.89)", Value: 95},
					{Color: "rgba(245, 54, 54, 0.9)", Value: 99},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_requests_total{reporter=\"destination\",destination_service=~\"$service\",response_code!~\"5.*\"}[5m])) / (sum(irate(istio_requests_total{reporter=\"destination\",destination_service=~\"$service\"}[5m])) or on () vector(1))",
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func ServerRequestDurationChart(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Server Request Duration",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: &dashboards.SecondsUnit},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\",destination_service=~\"$service\"}[1m])) by (le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\",destination_service=~\"$service\"}[1m])) by (le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\",destination_service=~\"$service\"}[1m])) by (le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\",destination_service=~\"$service\"}[1m])) by (le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"destination\",destination_service=~\"$service\"}[1m])) by (le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=\"destination\",destination_service=~\"$service\"}[1m])) by (le))",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P99"),
			),
		),
	)
}

func TCPSentBytesStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Sent Bytes",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.MeanCalculation),
			statPanel.Format(commonSdk.Format{Unit: &dashboards.BytesPerSecondsUnit}),
			statPanel.WithSparkline(statPanel.Sparkline{
				Width: 1,
			}),
			statPanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{Color: "green", Value: 0},
					{Color: "red", Value: 80},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_tcp_sent_bytes_total{reporter=~\"$qrep\", destination_service=~\"$service\"}[1m]))",
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

// ========== CLIENT WORKLOAD PANELS ==========

func IncomingRequestsByClient(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Requests By Source And Response Code",
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
				"round(sum(irate(istio_requests_total{connection_security_policy=\"mutual_tls\",destination_service=~\"$service\",reporter=~\"$qrep\",source_workload=~\"$srcwl\",source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace, response_code), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} : {{ response_code }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"round(sum(irate(istio_requests_total{connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", reporter=~\"$qrep\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace, response_code), 0.001)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} : {{ response_code }}"),
			),
		),
	)
}

func IncomingSuccessRateByClient(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Success Rate (non-5xx responses) By Source",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{Unit: &dashboards.PercentDecimalUnit},
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
				"sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\",response_code!~\"5.*\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				"sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\",response_code!~\"5.*\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_service=~\"$service\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace)",
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}
