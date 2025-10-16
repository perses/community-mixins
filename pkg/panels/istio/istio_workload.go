package istio

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func IncomingRequestVolume(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Requests By Source And Response Code",
		panel.Description("Request volume to workload by source"),
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestVolume"],
					labelMatchers,
				).Pretty(0),
				// promql.SetLabelMatchers(
				// 	"round(sum(irate(istio_requests_total{reporter=~\"$qrep\", connection_security_policy=\"mutual_tls\", destination_workload=~\"$workload\", destination_workload_namespace=~\"$namespace\", source_workload=~\"$srcwl\", source_workload_namespace=~\"$srcns\"}[1m])) by (source_workload, source_workload_namespace, response_code), 0.001)",
				// 	labelMatchers,
				// ),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} : {{ response_code }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestVolumeNonmTLS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} : {{ response_code }}"),
			),
		),
	)
}

func IncomingSuccessRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Success Rate (non-5xx responses) By Source",
		panel.Description("Success rate of requests to workload by source"),
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingSuccessRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingSuccessRateNonmTLS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ source_workload }}.{{ source_workload_namespace }}"),
			),
		),
	)
}

func IncomingRequestDuration(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Incoming Request Duration By Source",
		panel.Description("Request duration percentiles for workload"),
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDuration50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDuration90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDuration95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDuration99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDurationNonmTLS50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDurationNonmTLS90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDurationNonmTLS95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IncomingRequestDurationNonmTLS99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{source_workload}}.{{source_workload_namespace}} P99"),
			),
		),
	)
}

func OutgoingRequestVolume(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestVolume"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} : {{ response_code }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestVolumeNonmTLS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} : {{ response_code }}"),
			),
		),
	)
}

func OutgoingSuccessRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Outgoing Success Rate (non-5xx responses) By Destination",
		panel.Description("Success rate of outgoing requests from workload"),
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
				LineWidth:    1,
				AreaOpacity:  0,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingSuccessRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingSuccessRateNonmTLS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}"),
			),
		),
	)
}

func OutgoingRequestDuration(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDuration50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDuration90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDuration95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDuration99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDurationNonmTLS50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDurationNonmTLS90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDurationNonmTLS95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestDurationNonmTLS99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99"),
			),
		),
	)
}

func OutgoingRequestSize(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSize50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSize90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSize95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSize99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSizeNonmTLS50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSizeNonmTLS90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSizeNonmTLS95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingRequestSizeNonmLTS99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99"),
			),
		),
	)
}
func OutgoingResponseSize(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSize50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSize90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSize95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSize99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}  P99 (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSizeNonmTLS50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P50"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSizeNonmTLS90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P90"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSizeNonmMTLS95"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P95"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["OutgoingResponseSizeNonmTLS99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} P99"),
			),
		),
	)
}

func TCPBytesReceived(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes Received from Outgoing TCP Connection",
		panel.Description("TCP bytes received by workload"),
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["TCPBytesReceived"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["TCPBytesReceivedNonmTLS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}"),
			),
		),
	)
}

func TCPBytesSent(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes Sent on Outgoing TCP Connection",
		panel.Description("TCP bytes sent from workload"),
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["TCPBytesSent"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }} (🔐mTLS)"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["TCPBytesSentNonmTLS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service }}"),
			),
		),
	)
}
