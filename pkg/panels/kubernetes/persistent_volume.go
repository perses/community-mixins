package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	gaugePanel "github.com/perses/plugins/gaugechart/sdk/go"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func VolumeSpaceUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Volume Space Usage",
		panel.Description("Shows the space usage of persistent volume in a namespace by a PV claim."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(\n  sum without(instance, node) (topk(1, (kubelet_volume_stats_capacity_bytes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})))\n  -\n  sum without(instance, node) (topk(1, (kubelet_volume_stats_available_bytes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})))\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Used Space"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum without(instance, node) (topk(1, (kubelet_volume_stats_available_bytes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})))\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Free Space"),
			),
		),
	)
}

func VolumeSpaceUsageGauge(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Volume Space Usage",
		panel.Description("Shows the space usage of persistent volume in a namespace by a PV claim."),
		gaugePanel.Chart(
			gaugePanel.Calculation(commonSdk.LastCalculation),
			gaugePanel.Format(commonSdk.Format{
				Unit: &dashboards.PercentUnit,
			}),
			gaugePanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "green",
						Value: 0,
					},
					{
						Color: "orange",
						Value: 80,
					},
					{
						Color: "red",
						Value: 90,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max without(instance,node) (\n(\n  topk(1, kubelet_volume_stats_capacity_bytes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})\n  -\n  topk(1, kubelet_volume_stats_available_bytes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})\n)\n/\ntopk(1, kubelet_volume_stats_capacity_bytes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})\n* 100)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}

func VolumeInodesUsage(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Volume inodes Usage",
		panel.Description("Shows the inodes usage of persistent volume in a namespace by a PV claim."),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.TableMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  1,
				Stack:        timeSeriesPanel.AllStack,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"sum without(instance, node) (topk(1, (kubelet_volume_stats_inodes_used{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})))",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Used inodes"),
			),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"(\n  sum without(instance, node) (topk(1, (kubelet_volume_stats_inodes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})))\n  -\n  sum without(instance, node) (topk(1, (kubelet_volume_stats_inodes_used{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})))\n)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
				query.SeriesNameFormat("Free inodes"),
			),
		),
	)
}

func VolumeInodesUsageGauge(datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	return panelgroup.AddPanel("Volume inodes Usage",
		panel.Description("Shows the inodes usage of persistent volume in a namespace by a PV claim."),
		gaugePanel.Chart(
			gaugePanel.Calculation(commonSdk.LastCalculation),
			gaugePanel.Format(commonSdk.Format{
				Unit: &dashboards.PercentUnit,
			}),
			gaugePanel.Thresholds(commonSdk.Thresholds{
				Mode:         commonSdk.AbsoluteMode,
				DefaultColor: "green",
				Steps: []commonSdk.StepOption{
					{
						Color: "green",
						Value: 0,
					},
					{
						Color: "orange",
						Value: 80,
					},
					{
						Color: "red",
						Value: 90,
					},
				},
			}),
		),
		panel.AddQuery(
			query.PromQL(
				promql.SetLabelMatchers(
					"max without(instance,node) (\ntopk(1, kubelet_volume_stats_inodes_used{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})\n/\ntopk(1, kubelet_volume_stats_inodes{cluster=\"$cluster\", "+GetKubeletMatcher()+", namespace=\"$namespace\", persistentvolumeclaim=\"$volume\"})\n* 100)\n",
					labelMatchers,
				),
				dashboards.AddQueryDataSource(datasourceName),
			),
		),
	)
}
