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
	return panelgroup.AddPanel("Incoming Requests By Source And Response Code",
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
					"round(sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, response_code), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} : {{ response_code }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, response_code), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} : {{ response_code }}"),
			),
		),
	)
}

func IncomingSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Success Rate (non-5xx responses) By Source",
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
					"sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload_namespace=~\"$namespace\", destination_workload=~\"$workload\",response_code!~\"5.*\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload_namespace=~\"$namespace\", destination_workload=~\"$workload\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload_namespace=~\"$namespace\", destination_workload=~\"$workload\",response_code!~\"5.*\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace) / sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload_namespace=~\"$namespace\", destination_workload=~\"$workload\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[5m])) by (source_workload, source_workload_namespace)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func IncomingRequestDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Duration By Source",
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
					"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.95, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.95, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.95, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.95, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=~\"$qrep\", connection_security_policy!=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P99"),
			),
		),
	)
}

func OutgoingRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Requests By Destination And Response Code",
		panel.Description("Request volume from workload to destinations"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Min: 0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(irate(istio_requests_total{destination_principal=~\"spiffe.*\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", reporter=\"source\", destination_service=~\"$dstsvc\"}[5m])) by (destination_service, response_code), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} : {{ response_code }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(irate(istio_requests_total{destination_principal!~\"spiffe.*\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", reporter=\"source\", destination_service=~\"$dstsvc\"}[5m])) by (destination_service, response_code), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} : {{ response_code }}"),
			),
		),
	)
}

func OutgoingSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Success Rate (non-5xx responses) By Destination",
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
				LineWidth:    1,
				AreaOpacity:  0,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(istio_requests_total{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\",response_code!~\"5.*\", destination_service=~\"$dstsvc\"}[5m])) by (destination_service) / sum(irate(istio_requests_total{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", destination_service=~\"$dstsvc\"}[5m])) by (destination_service)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(irate(istio_requests_total{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\",response_code!~\"5.*\", destination_service=~\"$dstsvc\"}[5m])) by (destination_service) / sum(irate(istio_requests_total{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", destination_service=~\"$dstsvc\"}[5m])) by (destination_service)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}"),
			),
		),
	)
}

func OutgoingRequestDuration(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Request Duration By Destination",
		panel.Description("Request volume from workload to destinations"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Min: 0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.95, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.95, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.50, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.50, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.90, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.90, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.95, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.95, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(histogram_quantile(0.99, sum(irate(istio_request_duration_milliseconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le)) / 1000) or histogram_quantile(0.99, sum(irate(istio_request_duration_seconds_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99"),
			),
		),
	)
}

func OutgoingRequestSize(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Request Size By Destination",
		panel.Description("Request volume from workload to destinations"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Min: 0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.95, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.95, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(irate(istio_request_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99"),
			),
		),
	)
}
func OutgoingResponseSize(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Response Size By Destination",
		panel.Description("Request volume from workload to destinations"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Min: 0,
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    1,
				AreaOpacity:  0,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.95, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}  P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.50, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.95, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(irate(istio_response_bytes_bucket{reporter=\"source\", connection_security_policy!=\"mutual_tls\", source_workload=~\"$workload\", source_workload_namespace=~\"$namespace\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service, le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99"),
			),
		),
	)
}

func TCPBytesReceived(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes Received from Outgoing TCP Connection",
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
					"round(sum(irate(istio_tcp_received_bytes_total{connection_security_policy=\"mutual_tls\", reporter=\"source\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", destination_service=~\"$dstsvc\"}[$__rate_interval])) by (destination_service), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(irate(istio_tcp_received_bytes_total{connection_security_policy!=\"mutual_tls\", reporter=\"source\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", destination_service=~\"$dstsvc\"}[$__rate_interval])) by (destination_service), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}"),
			),
		),
	)
}

func TCPBytesSent(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes Sent on Outgoing TCP Connection",
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
					"round(sum(irate(istio_tcp_received_bytes_total{connection_security_policy=\"mutual_tls\", reporter=\"source\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(irate(istio_tcp_received_bytes_total{connection_security_policy!=\"mutual_tls\", reporter=\"source\", source_workload_namespace=~\"$namespace\", source_workload=~\"$workload\", destination_service=~\"$dstsvc\"}[1m])) by (destination_service), 0.001)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}"),
			),
		),
	)
}
