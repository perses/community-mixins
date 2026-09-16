// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istio

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	statPanel "github.com/perses/plugins/statchart/sdk/go"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func HTTPGRPCWorkloads(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP/gRPC Workloads",
		panel.Description("Request information for HTTP services"),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "destination_service",
					Header: "Service",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "destination_workload_var",
					Header: "Workload",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Requests",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.RequestsPerSecondsUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name:   "value #2",
					Header: "P50 Latency",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.MilliSecondsUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name:   "value #3",
					Header: "P90 Latency",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.MilliSecondsUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name:   "value #4",
					Header: "P99 Latency",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.MilliSecondsUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name:   "value #5",
					Header: "Success Rate",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.PercentDecimalUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name: "destination_workload_namespace",
					Hide: true,
				},
				{
					Name: "destination_workload",
					Hide: true,
				},
				{
					Name: "timestamp",
					Hide: true,
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioHTTPGRPCWorkloads"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioHTTPGRPCWorkloads50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioHTTPGRPCWorkloads90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioHTTPGRPCWorkloads99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioHTTPGRPCWorkloadsReqTotal"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func TCPServices(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("TCP Workloads",
		panel.Description("Bytes sent and recieived information for TCP services"),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "destination_service",
					Header: "Service",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "destination_workload_var",
					Header: "Workload",
					Align:  tablePanel.LeftAlign,
				},
				{
					Name:   "value #1",
					Header: "Bytes Received",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.BytesPerSecondsUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name:   "value #2",
					Header: "Bytes Sent",
					Align:  tablePanel.RightAlign,
					Format: &commonSdk.Format{
						Unit:          &dashboards.BytesPerSecondsUnit,
						DecimalPlaces: 2,
					},
				},
				{
					Name: "destination_workload_namespace",
					Hide: true,
				},
				{
					Name: "destination_workload",
					Hide: true,
				},
				{
					Name: "timestamp",
					Hide: true,
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioTCPServicesBytesRecv"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioTCPServicesBytesSent"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ destination_workload}}.{{ destination_workload_namespace }}"),
			),
		),
	)
}

func GlobalRequestVolume(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioGlobalRequestVolume"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func GlobalSuccessRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstionGlobalSuccessRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func Global4xxRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioGlobal4xxRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func Global5xxRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioGlobal5xxRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func IstioComponentVersions(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
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
				promql.SetLabelMatchersV2(
					IstioCommonPanelQueries["IstioComponentVersions"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{component}} ({{tag}})"),
			),
		),
	)
}
