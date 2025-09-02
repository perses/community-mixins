package kubernetes

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"

	commonSdk "github.com/perses/perses/go-sdk/common"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
)

func KubernetesIOPS(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		panelName = "IOPS(Reads+Writes)"
		description = "Shows IOPS(Reads+Writes) by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"ceil(sum by(namespace) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]) + rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval])))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		panelName = "IOPS(Reads+Writes)"
		description = "Shows IOPS(Reads+Writes) by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"ceil(sum by(pod) (rate(container_fs_reads_total{container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace=\"$namespace\"}[$__rate_interval]) + rate(container_fs_writes_total{container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace=\"$namespace\"}[$__rate_interval])))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		panelName = "IOPS(Pod)"
		description = "Shows IOPS of a pod, split by read and write."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"ceil(sum by(pod) (rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\",namespace=\"$namespace\", pod=~\"$pod\"}[$__rate_interval])))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("Writes"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"ceil(sum by(pod) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[$__rate_interval])))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("Reads"),
				),
			),
		}
	case "pod-container":
		panelName = "IOPS(Reads+Writes Containers) "
		description = "Shows IOPS(Reads+Writes) by containers in a pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"ceil(sum by(container) (rate(container_fs_reads_total{"+GetCAdvisorMatcher()+", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\"}[$__rate_interval]) + rate(container_fs_writes_total{"+GetCAdvisorMatcher()+", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\"}[$__rate_interval])))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.OpsPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}

func KubernetesThroughput(granularity, datasourceName string, labelMatchers ...promql.LabelMatcher) panelgroup.Option {
	var panelName, description string
	var queries []panel.Option

	switch granularity {
	case "cluster":
		panelName = "ThroughPut(Read+Write)"
		description = "Shows Throughput(Read+Write) by namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by(namespace) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]) + rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", container!=\"\", cluster=\"$cluster\", namespace!=\"\"}[$__rate_interval]))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{namespace}}"),
				),
			),
		}
	case "namespace-pod":
		panelName = "ThroughPut(Read+Write)"
		description = "Shows Throughput(Read+Write) by pods in a namespace."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by(pod) (rate(container_fs_reads_bytes_total{container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace=\"$namespace\"}[$__rate_interval]) + rate(container_fs_writes_bytes_total{container!=\"\", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", cluster=\"$cluster\", namespace=\"$namespace\"}[$__rate_interval]))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{pod}}"),
				),
			),
		}
	case "pod":
		panelName = "ThroughPut(Pod)"
		description = "Shows Throughput of a pod, split by read and write."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by(pod) (rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[$__rate_interval]))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("Writes"),
				),
			),
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by(pod) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", device=~\"(/dev.+)|mmcblk.p.+|nvme.+|rbd.+|sd.+|vd.+|xvd.+|dm-.+|dasd.+\", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=~\"$pod\"}[$__rate_interval]))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("Reads"),
				),
			),
		}
	case "pod-container":
		panelName = "ThroughPut(Reads+Writes Containers) "
		description = "Shows Throughput(Reads+Writes) by containers in a pod."
		queries = []panel.Option{
			panel.AddQuery(
				query.PromQL(
					promql.SetLabelMatchers(
						"sum by(container) (rate(container_fs_reads_bytes_total{"+GetCAdvisorMatcher()+", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\"}[$__rate_interval]) + rate(container_fs_writes_bytes_total{"+GetCAdvisorMatcher()+", container!=\"\", cluster=\"$cluster\", namespace=\"$namespace\", pod=\"$pod\"}[$__rate_interval]))",
						labelMatchers,
					),
					dashboards.AddQueryDataSource(datasourceName),
					query.SeriesNameFormat("{{container}}"),
				),
			),
		}
	}

	panelOpts := []panel.Option{
		panel.Description(description),
		timeSeriesPanel.Chart(
			timeSeriesPanel.WithYAxis(timeSeriesPanel.YAxis{
				Format: &commonSdk.Format{
					Unit: &dashboards.BytesPerSecondsUnit,
				},
			}),
			timeSeriesPanel.WithLegend(timeSeriesPanel.Legend{
				Position: timeSeriesPanel.BottomPosition,
				Mode:     timeSeriesPanel.ListMode,
				Size:     timeSeriesPanel.SmallSize,
			}),
			timeSeriesPanel.WithVisual(timeSeriesPanel.Visual{
				Display:      timeSeriesPanel.LineDisplay,
				ConnectNulls: false,
				LineWidth:    0.25,
				AreaOpacity:  0.75,
				Palette:      timeSeriesPanel.Palette{Mode: timeSeriesPanel.AutoMode},
			}),
		),
	}
	panelOpts = append(panelOpts, queries...)

	return panelgroup.AddPanel(panelName, panelOpts...)
}
