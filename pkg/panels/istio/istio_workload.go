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

func IncomingRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Volume",
		panel.Description("Request volume to workload by source"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.RequestsPerSecondsUnit),
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
					"round(sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("🔒 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func IncomingSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Success Rate",
		panel.Description("Success rate of requests to workload by source"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
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
					"sum(irate(istio_requests_total{reporter=~\"$qrep\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\", response_code!~\"5.*\"}[$__rate_interval])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter=~\"$qrep\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func IncomingRequestDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Duration",
		panel.Description("Request duration percentiles for workload"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
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
					"histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P50 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P90 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("P99 {{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func OutgoingRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Request Volume",
		panel.Description("Request volume from workload to destinations"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.RequestsPerSecondsUnit),
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
					"round(sum(irate(istio_requests_total{reporter=\"source\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func OutgoingSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Success Rate",
		panel.Description("Success rate of outgoing requests from workload"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.PercentUnit),
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
					"sum(irate(istio_requests_total{reporter=\"source\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\", response_code!~\"5.*\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace) / sum(irate(istio_requests_total{reporter=\"source\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func TCPBytesReceived(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Bytes Received",
		panel.Description("TCP bytes received by workload"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
					"sum(irate(istio_tcp_received_bytes_total{reporter=~\"$qrep\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[$__rate_interval])) by (source_workload, source_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func TCPBytesSent(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Bytes Sent",
		panel.Description("TCP bytes sent from workload"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesPerSecondsUnit),
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
					"sum(irate(istio_tcp_sent_bytes_total{reporter=\"source\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\", destination_workload=~\"$dstwl\", destination_workload_namespace=~\"$dstns\"}[$__rate_interval])) by (destination_workload, destination_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload }}.{{ destination_workload_namespace }}"),
			),
		),
	)
}
