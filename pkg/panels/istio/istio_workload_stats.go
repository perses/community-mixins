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

// Stat panels for General section
func IncomingRequestVolumeStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Volume",
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
				promql.SetLabelMatchers(
					"round(sum(irate(istio_requests_total{reporter=~\"$qrep\",destination_workload_namespace=~\"$namespace\",destination_workload=~\"$workload\"}[5m])), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func IncomingSuccessRateStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Success Rate (non-5xx responses)",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.MeanCalculation),
			statPanel.Format(commonSdk.Format{Unit: string(commonSdk.PercentDecimalUnit)}),
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
				promql.SetLabelMatchers(
					"sum(irate(istio_requests_total{reporter=~\"$qrep\",destination_workload_namespace=~\"$namespace\",destination_workload=~\"$workload\",response_code!~\"5.*\"}[5m])) / sum(irate(istio_requests_total{reporter=~\"$qrep\",destination_workload_namespace=~\"$namespace\",destination_workload=~\"$workload\"}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func RequestDurationChart(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Request Duration",
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.SecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0.1,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\",destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\"}[1m])) by (le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\",destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\"}[1m])) by (le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\",destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\"}[1m])) by (le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\",destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\"}[1m])) by (le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\",destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\"}[1m])) by (le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\",destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\"}[1m])) by (le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P99"),
			),
		),
	)
}

func TCPServerTrafficStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Server Traffic",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.MeanCalculation),
			statPanel.Format(commonSdk.Format{Unit: string(commonSdk.BytesPerSecondsUnit)}),
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
				promql.SetLabelMatchers(
					"sum(irate(istio_tcp_sent_bytes_total{reporter=~\"$qrep\", destination_workload_namespace=~\"$namespace\", destination_workload=~\"$workload\"}[1m])) + sum(irate(istio_tcp_received_bytes_total{reporter=~\"$qrep\", destination_workload_namespace=~\"$namespace\", destination_workload=~\"$workload\"}[1m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TCPClientTrafficStat(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Client Traffic",
		statPanel.Chart(
			statPanel.Calculation(commonSdk.MeanCalculation),
			statPanel.Format(commonSdk.Format{Unit: string(commonSdk.BytesPerSecondsUnit)}),
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
				promql.SetLabelMatchers(
					"sum(irate(istio_tcp_sent_bytes_total{reporter=\"source\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\"}[1m])) + sum(irate(istio_tcp_received_bytes_total{reporter=\"source\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\"}[1m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
