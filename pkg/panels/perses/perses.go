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

package perses

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/prometheus/prometheus/model/labels"
)

func StatsTable(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Perses Stats",
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "job",
					Header: "Job",
				},
				{
					Name:   "instance",
					Header: "Instance",
				},
				{
					Name:   "version",
					Header: "Version",
				},
				{
					Name:   "namespace",
					Header: "Namespace",
				},
				{
					Name:   "pod",
					Header: "Pod",
				},
				{
					Name: "value",
					Hide: true,
				},
				{
					Name: "timestamp",
					Hide: true,
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesStatsTable"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func HTTPRequestsLatencyPanel(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP Requests Latency",
		panel.Description("Displays the latency of HTTP requests over a 5-minute window."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: &dashboards.SecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesHTTPLatency"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{handler}} {{method}}"),
			),
		),
	)
}

func HTTPRequestsRatePanel(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP Requests Rate",
		panel.Description("Displays the rate of HTTP requests over a 5-minute window."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesHTTPRequestRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{handler}} {{method}} {{code}}"),
			),
		),
	)
}

func HTTPErrorPercentagePanel(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel(
		"HTTP Errors Percentage",
		panel.Description("Displays the percentage of HTTP errors over total requests."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(
				timeSeriesPanel.YAxis{
					Format: &common.Format{
						Unit: &dashboards.PercentDecimalUnit,
					},
				},
			),
			timeSeriesPanel.WithLegend(
				timeSeriesPanel.Legend{
					Position: timeSeriesPanel.RightPosition,
					Mode:     timeSeriesPanel.TableMode,
				},
			),
			timeSeriesPanel.WithVisual(
				timeSeriesPanel.Visual{
					Display:      timeSeriesPanel.LineDisplay,
					ConnectNulls: false,
					LineWidth:    0.25,
					AreaOpacity:  0.5,
					Stack:        timeSeriesPanel.AllStack,
					Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
				},
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesHTTPErrorsPercentage"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{handler}} {{code}}"),
			),
		),
	)
}

func FileDescriptors(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("File Descriptors",
		panel.Description("Displays the number of open and max file descriptors."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesFileDescriptorsOpenFDS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{pod}} Open FDs"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesFileDescriptorsMaxFDS"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{pod}} - Max FDs"),
			),
		),
	)
}

func PluginSchemaLoadAttempts(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Plugin Schema Load Attempts",
		panel.Description("Displays the success and failure attempts to load plugin schemas."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.5,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					PersesCommonPanelQueries["PersesPluginSchemaLoadAttempts"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}} - {{schema}} - {{status}}"),
			),
		),
	)
}
