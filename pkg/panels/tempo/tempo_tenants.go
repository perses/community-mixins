package tempo

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/promql"
	commonSdk "github.com/perses/perses/go-sdk/common"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	tablePanel "github.com/perses/plugins/table/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func TenantInfo(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Tenant info",
		panel.Description("Displays the effective Tempo limits (limit_name)"),
		tablePanel.Table(
			tablePanel.WithColumnSettings([]tablePanel.ColumnSettings{
				{
					Name:   "limit_name",
					Header: "limit name",
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(\n  max by (cluster, namespace, limit_name) (tempo_limits_overrides{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",user=\"$tenant\"})\n  or max by (cluster, namespace, limit_name) (tempo_limits_defaults{cluster=~\"$cluster\", job=~\"($namespace)/compactor\"})\n) by (limit_name)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func TenantDistributorBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Distributor bytes/s",
		panel.Description("Data ingestion rate (in bytes/sec) for the Tempo distributor, filtered by tenant, cluster, and namespace."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_bytes_received_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\",tenant=\"$tenant\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("received"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(\n  max by (cluster, namespace, limit_name) (tempo_limits_overrides{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",user=\"$tenant\",limit_name=\"ingestion_rate_limit_bytes\"})\n  or max by (cluster, namespace, limit_name) (tempo_limits_defaults{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",limit_name=\"ingestion_rate_limit_bytes\"})\n) by (ingestion_rate_limit_bytes)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("limit"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(\n  max by (cluster, namespace, limit_name) (tempo_limits_overrides{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",user=\"$tenant\",limit_name=\"ingestion_burst_size_bytes\"})\n  or max by (cluster, namespace, limit_name) (tempo_limits_defaults{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",limit_name=\"ingestion_burst_size_bytes\"})\n) by (ingestion_burst_size_bytes)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("burst limit"),
			),
		),
	)
}

func TenantDistributorSpan(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Distributor span/s",
		panel.Description("Rate of Spans per Second for Tempo Distributor"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_distributor_spans_received_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\",tenant=\"$tenant\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("accepted"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_discarded_spans_total{cluster=~\"$cluster\", job=~\"($namespace)/distributor\",tenant=\"$tenant\"}[$__rate_interval])) by (reason)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("refused {{ reason }}"),
			),
		),
	)
}

func TenantLiveTraces(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Live traces",
		panel.Description("Tempo ingester traces in memory"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(tempo_ingester_live_traces{cluster=~\"$cluster\", job=~\"($namespace)/ingester\",tenant=\"$tenant\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("live traces"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(\n  max by (cluster, namespace, limit_name) (tempo_limits_overrides{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",user=\"$tenant\",limit_name=\"max_global_traces_per_user\"})\n  or max by (cluster, namespace, limit_name) (tempo_limits_defaults{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",limit_name=\"max_global_traces_per_user\"})\n) by (max_global_traces_per_user)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("global limit"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(\n  max by (cluster, namespace, limit_name) (tempo_limits_overrides{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",user=\"$tenant\",limit_name=\"max_local_traces_per_user\"})\n  or max by (cluster, namespace, limit_name) (tempo_limits_defaults{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",limit_name=\"max_local_traces_per_user\"})\n) by (max_local_traces_per_user)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("local limit"),
			),
		),
	)
}

func TenantQueriesID(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Queries/s (ID lookup)",
		panel.Description("Rate per second of traces queries handled by Tempo's query-frontend"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_query_frontend_queries_total{cluster=~\"$cluster\", job=~\"($namespace)/query-frontend\",tenant=\"$tenant\",op=\"traces\"}[$__rate_interval])) by (status)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ status }}"),
			),
		),
	)
}

func TenantQueriesSearch(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Queries/s (search)",
		panel.Description("Rate per second of Tempo search queries processed by the query-frontend"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_query_frontend_queries_total{cluster=~\"$cluster\", job=~\"($namespace)/query-frontend\",tenant=\"$tenant\",op=\"search\"}[$__rate_interval])) by (status)",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ status }}"),
			),
		),
	)
}

func TenantBlockslistLength(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Blockslist length",
		panel.Description("Average number of blocks currently listed in the blocklist of the Tempo compactor"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"avg(tempodb_blocklist_length{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",tenant=\"$tenant\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("length"),
			),
		),
	)
}

func TenantOutstandingCompactions(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Outstanding compactions",
		panel.Description("Average number of Tempo blocks awaiting compaction per compactor instance"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(tempodb_compaction_outstanding_blocks{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",tenant=\"$tenant\"})\n/\ncount(tempo_build_info{cluster=~\"$cluster\", job=~\"($namespace)/compactor\"})\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("blocks"),
			),
		),
	)
}

func TenantMetricGeneratorBytes(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Bytes/s",
		panel.Description("Rate of trace data (in bytes/sec) ingested by the Tempo metrics generato"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(rate(tempo_metrics_generator_bytes_received_total{cluster=~\"$cluster\", job=~\"($namespace)/metrics-generator\",tenant=\"$tenant\"}[$__rate_interval]))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("rate"),
			),
		),
	)
}

func TenantMetricGeneratorActiveSeries(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Active series",
		panel.Description("Number of active metric series registered in the Tempo Metrics Generator"),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit:          &dashboards.DecimalUnit,
					DecimalPlaces: 0,
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
				Palette:      &timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum(tempo_metrics_generator_registry_active_series{cluster=~\"$cluster\", job=~\"($namespace)/metrics-generator\",tenant=\"$tenant\"})",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("{{ tenant }}"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max(\n  max by (cluster, namespace, limit_name) (tempo_limits_overrides{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",user=\"$tenant\",limit_name=\"metrics_generator_max_active_series\"})\n  or max by (cluster, namespace, limit_name) (tempo_limits_defaults{cluster=~\"$cluster\", job=~\"($namespace)/compactor\",limit_name=\"metrics_generator_max_active_series\"})\n) by (metrics_generator_max_active_series)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("limit"),
			),
		),
	)
}
