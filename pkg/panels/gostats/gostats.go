package gostats

import (
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/perses/go-sdk/prometheus/query"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
)

func MemoryUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
		panel.Description("Shows various memory usage metrics of the component."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_alloc_bytes",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Alloc All {{pod}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_heap_alloc_bytes",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Alloc Heap {{pod}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(go_memstats_alloc_bytes_total[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Alloc Rate All {{pod}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(go_memstats_heap_alloc_bytes[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Alloc Rate Heap {{pod}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_stack_inuse_bytes",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Inuse Stack {{pod}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_heap_inuse_bytes",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Inuse Heap {{pod}}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"process_resident_memory_bytes",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Resident Memory {{pod}}"),
			),
		),
	)
}

func Goroutines(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Goroutines",
		panel.Description("Shows the number of goroutines being used by the component."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_goroutines",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}}"),
			),
		),
	)
}

func GarbageCollectionPauseTimeQuantiles(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("GC Duration",
		panel.Description("Shows the Go garbage collection pause durations for the component."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: string(commonSdk.MilliSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []commonSdk.Calculation{commonSdk.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_gc_duration_seconds",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{quantile}} - {{pod}}"),
			),
		),
	)
}
