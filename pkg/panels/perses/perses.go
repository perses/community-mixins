package perses

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	timeSeriesPanel "github.com/perses/perses/go-sdk/panel/time-series"
	"github.com/perses/perses/go-sdk/prometheus/query"

	"github.com/perses/perses/go-sdk/common"
	tablePanel "github.com/perses/perses/go-sdk/panel/table"
)

// PersesStatsTable creates a panel option for displaying Perses statistics.
//
// The panel uses the following Prometheus metrics:
// - perses_build_info: Build information about Perses instances
//
// The panel shows:
// - Instance count by job and version
// - Version information per instance
func PersesStatsTable(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
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
				promql.SetLabelMatchers("count by (job, instance, version, namespace, pod) (perses_build_info{job=~'$job', instance=~'$instance'})", labelMatchers),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func HTTPRequestsLatencyPanel(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP Requests Latency",
		panel.Description("Displays the latency of HTTP requests over a 5-minute window."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.MilliSecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (handler, method) (rate(perses_http_request_duration_second_sum{job=~'$job', instance=~'$instance'}[5m])) / sum by (handler, method) (rate(perses_http_request_duration_second_count{job=~'$job', instance=~'$instance'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{handler}} {{method}}"),
			),
		),
	)
}

func HTTPRequestsRatePanel(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("HTTP Requests Rate",
		panel.Description("Displays the rate of HTTP requests over a 5-minute window."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.RightPosition,
				Mode:     timeSeriesPanel.TableMode,
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum by (handler, code) (rate(perses_http_request_total{job=~'$job', instance=~'$instance'}[5m]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{handler}} {{method}} {{code}}"),
			),
		),
	)
}

func PersesMemoryUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Memory Usage",
		panel.Description("Shows various memory usage metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.BytesUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_alloc_bytes{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Heap Allocated"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_heap_inuse_bytes{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Heap In Use"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_memstats_stack_inuse_bytes{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Stack In Use"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"process_resident_memory_bytes{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Resident Memory"),
			),
		),
	)
}

func PersesCPUUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("CPU Usage",
		panel.Description("Shows CPU usage metrics"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.PercentUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"rate(process_cpu_seconds_total{job=~'$job', instance=~'$instance'}[5m])",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}}"),
			),
		),
	)
}

func PersesGoRoutines(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Goroutines",
		panel.Description("Shows the number of goroutines currently in use"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_goroutines{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{pod}}"),
			),
		),
	)
}

func PersesGarbageCollectionPauseTime(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Garbage Collection Pause Time",
		panel.Description("Displays the pause time for garbage collection events."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.SecondsUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"go_gc_duration_seconds{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{quantile}} - {{instance}} - {{pod}}"),
			),
		),
	)
}

func PersesFileDescriptors(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("File Descriptors",
		panel.Description("Displays the number of open and max file descriptors."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &common.Format{
					Unit: string(common.DecimalUnit),
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Values:   []common.Calculation{common.LastCalculation},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"process_open_fds{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{pod}} Open FDs"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"process_max_fds{job=~'$job', instance=~'$instance'}",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{instance}} - {{pod}} - Max FDs"),
			),
		),
	)
}
