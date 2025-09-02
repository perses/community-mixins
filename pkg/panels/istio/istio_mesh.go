package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	statPanel "github.com/perses/plugins/statchart/sdk/go"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func HTTPGRPCWorkloads(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP/gRPC Workloads",
		panel.Description("Request information for HTTP services"),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "Value #requests",
					Header: "Requests",
				},
				{
					Name:   "Value #p50",
					Header: "P50 Latency",
				},
				{
					Name:   "Value #p90",
					Header: "P90 Latency",
				},
				{
					Name:   "Value #p99",
					Header: "P99 Latency",
				},
				{
					Name:   "Value #success",
					Header: "Success Rate",
				},
				{
					Name:   "destination_workload_var",
					Header: "Workload",
				},
				{
					Name:   "destination_service",
					Header: "Service",
				},
				{
					Name: "destination_workload_namespace",
				},
				{
					Name: "destination_workload",
				},
				{
					Name: "timestamp",
				},
			}),

			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(sum by (destination_workload,destination_workload_namespace,destination_service)(rate(istio_requests_total{reporter=~\"source|waypoint\"}[$__rate_interval])), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(histogram_quantile(0.5, sum by (le,destination_workload,destination_workload_namespace) (rate(istio_request_duration_milliseconds_bucket{reporter=~\"source|waypoint\"}[$__rate_interval]))), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(histogram_quantile(0.9, sum by (le,destination_workload,destination_workload_namespace) (rate(istio_request_duration_milliseconds_bucket{reporter=~\"source|waypoint\"}[$__rate_interval]))), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(histogram_quantile(0.99, sum by (le,destination_workload,destination_workload_namespace) (rate(istio_request_duration_milliseconds_bucket{reporter=~\"source|waypoint\"}[$__rate_interval]))), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(sum by (destination_workload,destination_workload_namespace)(rate(istio_requests_total{reporter=~\"source|waypoint\",response_code!~\"5..\"}[$__rate_interval]))/sum by (destination_workload,destination_workload_namespace)(rate(istio_requests_total{reporter=~\"source|waypoint\"}[$__rate_interval])), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func TCPServices(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Workloads",
		panel.Description("Bytes sent and recieived information for TCP services"),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "Value #recv",
					Header: "Bytes Received",
				},
				{
					Name:   "Value #sent",
					Header: "Bytes Sent",
				},
				{
					Name:   "destination_workload_var",
					Header: "Workload",
				},
				{
					Name:   "destination_service",
					Header: "Service",
				},
				{
					Name: "destination_workload_namespace",
				},
				{
					Name: "destination_workload",
				},
				{
					Name: "timestamp",
				},
			}),

			tablePanel.Transform([]commonSdk.Transform{
				{
					Kind: commonSdk.MergeSeriesKind,
					Spec: commonSdk.MergeSeriesSpec{
						Disabled: false,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(sum by (destination_workload,destination_workload_namespace,destination_service) (rate(istio_tcp_received_bytes_total{reporter=~\"source|waypoint\"}[$__rate_interval])), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"label_join(sum by (destination_workload,destination_workload_namespace,destination_service) (rate(istio_tcp_sent_bytes_total{reporter=~\"source|waypoint\"}[$__rate_interval])), \"destination_workload_var\", \".\", \"destination_workload\", \"destination_workload_namespace\")",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func GlobalRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Traffic Volume",
		panel.Description("Total requests in the cluster"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: &dashboards.RequestsPerSecondsUnit,
			}),
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
					"round(sum (rate(istio_requests_total{reporter=~\"source|waypoint\"}[$__rate_interval])), 0.01)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func GlobalSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Success Rate",
		panel.Description("Total success rate of requests in the cluster"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: &dashboards.PercentDecimalUnit,
			}),
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
					"sum(rate(istio_requests_total{reporter=~\"source|waypoint\",response_code!~\"5..\"}[$__rate_interval])) / sum(rate(istio_requests_total{reporter=~\"source|waypoint\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func Global4xxRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("4xxs",
		panel.Description("Total 4xx requests in in the cluster"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: &dashboards.RequestsPerSecondsUnit,
			}),
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
					"round(sum(rate(istio_requests_total{reporter=~\"source|waypoint\",response_code=~\"4..\"}[$__rate_interval])), 0.01)or vector(0)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func Global5xxRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("5xxs",
		panel.Description("Total 5xx requests in in the cluster"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: &dashboards.RequestsPerSecondsUnit,
			}),
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
					"round(sum(rate(istio_requests_total{reporter=~\"source|waypoint\",response_code=~\"5..\"}[$__rate_interval])), 0.01)or vector(0)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func IstioComponentVersions(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Istio Component Versions",
		panel.Description("Version number of each running instance"),
		timeSeriesPanel.Chart(
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
				promql.SetLabelMatchers(
					"sum by (component,tag) (istio_build)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{component}} ({{tag}})"),
			),
		),
	)
}
