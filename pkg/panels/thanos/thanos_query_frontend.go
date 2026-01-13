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

package thanos

import (
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	"github.com/prometheus/prometheus/model/labels"

	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func QueryFrontendRequestRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Query Request Rate",
		panel.Description("Shows rate of requests on query frontend API."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.RequestsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendRequestRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{code}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendQueryRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Query Rate",
		panel.Description("Shows rate of queries passing through Query Frontend."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.RequestsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendQueryRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{code}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendErrors(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Query Error Rate",
		panel.Description("Shows the percentage of errors compared to the total number of requests on query frontend API."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.PercentUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendErrors"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{code}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendDurations(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Query Duration",
		panel.Description("Shows p50, p90 and p99 of the time taken to respond to a query via the query frontend API."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.SecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
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
					ThanosCommonPanelQueries["QueryFrontendDurations_50"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p50 {{job}} {{namespace}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendDurations_90"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p90 {{job}} {{namespace}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendDurations_99"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("p99 {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendCacheRequestRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Cache Requests",
		panel.Description("Shows rate of cache requests by cortex."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.RequestsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendCacheRequestRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tripperware}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendCacheHitRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Cache Hits",
		panel.Description("Shows rate of cache hits by cortex."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.CountsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendCacheHitRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tripperware}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendCacheMissRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Cache Misses",
		panel.Description("Shows rate of cache misses by cortex."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.CountsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendCacheMissRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tripperware}} - {{job}} {{namespace}}"),
			),
		),
	)
}

func QueryFrontendFetchedKeyRate(datasourceName string, labelMatchers ...*labels.Matcher) panelgroup.Option {
	return panelgroup.AddPanel("Fetched Keys",
		panel.Description("Shows rate of keys fetched from cache by cortex."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.DecimalUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchersV2(
					ThanosCommonPanelQueries["QueryFrontendFetchedKeyRate"],
					labelMatchers,
				).Pretty(0),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{tripperware}} - {{job}} {{namespace}}"),
			),
		),
	)
}
