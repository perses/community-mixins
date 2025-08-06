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
					"histogram_quantile(0.50, sum(rate(istio_request_duration_milliseconds_bucket{reporter=~\"source|waypoint\"}[$__rate_interval])) by (destination_workload,destination_workload_namespace,le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.90, sum(rate(istio_request_duration_milliseconds_bucket{reporter=~\"source|waypoint\"}[$__rate_interval])) by (destination_workload,destination_workload_namespace,le))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"histogram_quantile(0.99, sum(rate(istio_request_duration_milliseconds_bucket{reporter=~\"source|waypoint\"}[$__rate_interval])) by (destination_workload,destination_workload_namespace,le))",
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
	return panelgroup.AddPanel("TCP Services",
		panel.Description("TCP traffic information for services"),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "Value #sent",
					Header: "Bytes Sent",
				},
				{
					Name:   "Value #received",
					Header: "Bytes Received",
				},
				{
					Name:   "destination_service_name",
					Header: "Service",
				},
				{
					Name:   "destination_service_namespace",
					Header: "Namespace",
				},
				{
					Name: "timestamp",
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (destination_service_name,destination_service_namespace)(rate(istio_tcp_sent_bytes_total{reporter=~\"source|waypoint\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service_name}}.{{ destination_service_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (destination_service_name,destination_service_namespace)(rate(istio_tcp_received_bytes_total{reporter=~\"source|waypoint\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_service_name}}.{{ destination_service_namespace }}"),
			),
		),
	)
}

func GlobalRequestVolume(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Global Request Volume",
		panel.Description("Total request volume across the mesh"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: string(commonSdk.RequestsPerSecondsUnit),
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"round(sum(rate(istio_requests_total{reporter=~\"source|waypoint\"}[$__rate_interval])), 0.01)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func GlobalSuccessRate(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Global Success Rate",
		panel.Description("Overall success rate across the mesh"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: string(commonSdk.PercentUnit),
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
	return panelgroup.AddPanel("Global 4xx Rate",
		panel.Description("4xx error rate across the mesh"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: string(commonSdk.RequestsPerSecondsUnit),
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
	return panelgroup.AddPanel("Global 5xx Rate",
		panel.Description("5xx error rate across the mesh"),
		statPanel.Chart(
			statPanel.Calculation(commonSdk.LastCalculation),
			statPanel.Format(commonSdk.Format{
				Unit: string(commonSdk.RequestsPerSecondsUnit),
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
