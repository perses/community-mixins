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

package opentelemetry

import (
	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/panels/opentelemetry"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	promqlVar "github.com/perses/plugins/prometheus/sdk/go/variable/promql"
	"github.com/prometheus/prometheus/model/labels"
)

func withReceiversGroup(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Receivers",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		opentelemetry.SpanRate(datasource, clusterLabelMatcher),
		opentelemetry.MeticPointsRate(datasource, clusterLabelMatcher),
		opentelemetry.LogRecordsRate(datasource, clusterLabelMatcher),
	)
}

func withProcessorsGroup(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Processors",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		opentelemetry.SpanProcessorRate(datasource, clusterLabelMatcher),
		opentelemetry.MetricProcessorRate(datasource, clusterLabelMatcher),
		opentelemetry.LogProcessorRate(datasource, clusterLabelMatcher),
		opentelemetry.BatchProcessorBatchSendSize(datasource, clusterLabelMatcher),
		opentelemetry.BatchProcessorBatchSendSizeCount(datasource, clusterLabelMatcher),
		opentelemetry.BatchProcessorBatchSizeTriggerSend(datasource, clusterLabelMatcher),
	)
}

func withExportersGroup(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Exporters",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		opentelemetry.SpanExporterRate(datasource, clusterLabelMatcher),
		opentelemetry.MetricExporterRate(datasource, clusterLabelMatcher),
		opentelemetry.LogExporterRate(datasource, clusterLabelMatcher),
		opentelemetry.QueueSizeExporterRate(datasource, clusterLabelMatcher),
		opentelemetry.QueueCapacityExporterRate(datasource, clusterLabelMatcher),
		opentelemetry.QueueUtilizationExporterRate(datasource, clusterLabelMatcher),
	)
}

func BuildOpenTelemetryCollector(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("opentelemetry-collector",
			dashboard.ProjectName(project),
			dashboard.Name("OpenTelemetry Collector"),
			dashboard.AddVariable("job",
				listVar.List(
					promqlVar.PrometheusPromQL(
						"group by (job) ({__name__=~'otelcol_process_uptime.*'})",
						promqlVar.Datasource(datasource),
						promqlVar.LabelName("job"),
					),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "up{job='$job'}"),
			dashboard.AddVariable(
				"receiver",
				listVar.List(
					promqlVar.PrometheusPromQL(
						"group by (receiver) ({__name__=~'otelcol_receiver_.+', job='$job'})",
						promqlVar.Datasource(datasource),
						promqlVar.LabelName("receiver"),
					),
					listVar.AllowAllValue(true),
				),
			),
			dashboard.AddVariable(
				"processor",
				listVar.List(
					promqlVar.PrometheusPromQL(
						"group by (processor) ({__name__=~'otelcol_processor_.+', job='$job'})",
						promqlVar.Datasource(datasource),
						promqlVar.LabelName("processor"),
					),
					listVar.AllowAllValue(true),
				),
			),
			dashboard.AddVariable(
				"exporter",
				listVar.List(
					promqlVar.PrometheusPromQL(
						"group by (exporter) ({__name__=~'otelcol_exporter_.+', job='$job'})",
						promqlVar.Datasource(datasource),
						promqlVar.LabelName("exporter"),
					),
					listVar.AllowAllValue(true),
				),
			),
			withReceiversGroup(datasource, clusterLabelMatcher),
			withProcessorsGroup(datasource, clusterLabelMatcher),
			withExportersGroup(datasource, clusterLabelMatcher),
		),
	).Component("opentelemetry-collector")
}
