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
	panels "github.com/perses/community-mixins/pkg/panels/istio"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	markdownPanel "github.com/perses/plugins/markdown/sdk/go"
	promqlVar "github.com/perses/plugins/prometheus/sdk/go/variable/promql"
	"github.com/prometheus/prometheus/model/labels"
)

// General section
func withGeneralSection() dashboard.Option {
	return dashboard.AddPanelGroup("General",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(4),
		panelgroup.AddPanel("SERVICE",
			markdownPanel.Markdown("SERVICE Header",
				markdownPanel.Text("<div class=\"dashboard-header text-center\">\n<span>SERVICE: $service</span>\n</div>"),
			),
		),
	)
}

func withGeneralSectionII(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("General",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		// First row of stats
		panels.ClientRequestVolumeStat(datasource, labelMatcher),
		panels.ClientSuccessRateStat(datasource, labelMatcher),
		panels.ClientRequestDurationChart(datasource, labelMatcher),
		panels.TCPReceivedBytesStat(datasource, labelMatcher),
		// Second row of stats
		panels.ServerRequestVolumeStat(datasource, labelMatcher),
		panels.ServerSuccessRateStat(datasource, labelMatcher),
		panels.ServerRequestDurationChart(datasource, labelMatcher),
		panels.TCPSentBytesStat(datasource, labelMatcher),
	)
}

// Client Workloads section
func withClientWorkloadsSection(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Client Workloads",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(6),
		// Header panel
		panelgroup.AddPanel("CLIENT WORKLOADS",
			markdownPanel.Markdown("Client Workloads Header",
				markdownPanel.Text("<div class=\"dashboard-header text-center\">\n<span>CLIENT WORKLOADS</span>\n</div>"),
			),
		),
		panels.IncomingRequestsByClient(datasource, labelMatcher),
		panels.IncomingSuccessRateByClient(datasource, labelMatcher),
	)
}

// Client Workloads (II) section
func withClientWorkloadsIISection(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Client Workloads (II)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(6),
		panels.IncomingRequestDurationByClient(datasource, labelMatcher),
		panels.IncomingRequestSizeByClient(datasource, labelMatcher),
		panels.ResponseSizeByClient(datasource, labelMatcher),
	)
}

// Client Workloads (III) section
func withClientWorkloadsIIISection(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Client Workloads (III)",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(6),
		panels.BytesReceivedFromTCPClient(datasource, labelMatcher),
		panels.BytesSentToTCPClient(datasource, labelMatcher),
	)
}

// Service Workloads section
func withServiceWorkloadsSection(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Service Workloads",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(6),
		// Header panel
		panelgroup.AddPanel("SERVICE WORKLOADS",
			markdownPanel.Markdown("Service Workloads Header",
				markdownPanel.Text("<div class=\"dashboard-header text-center\">\n<span>SERVICE WORKLOADS</span>\n</div>"),
			),
		),
		panels.IncomingRequestsByService(datasource, labelMatcher),
		panels.IncomingSuccessRateByService(datasource, labelMatcher),
	)
}

// Service Workloads (II) section
func withServiceWorkloadsIISection(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Service Workloads (II)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(6),
		panels.IncomingRequestDurationByService(datasource, labelMatcher),
		panels.IncomingRequestSizeByService(datasource, labelMatcher),
		panels.ResponseSizeByService(datasource, labelMatcher),
	)
}

// Service Workloads (III) section
func withServiceWorkloadsIIISection(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Service Workloads (III)",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(6),
		panels.BytesReceivedFromTCPService(datasource, labelMatcher),
		panels.BytesSentToTCPService(datasource, labelMatcher),
	)
}

func BuildIstioService(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcherV2(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-service-dashboard",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Service Dashboard"),
			// Service variable (single-select; panels use destination_service=~"$service")
			dashboard.AddVariable("service",
				listVar.List(
					promqlVar.PrometheusPromQL(
						"sum by (destination_service) (istio_requests_total{destination_service!=\"unknown\"}) or sum by (destination_service) (istio_tcp_sent_bytes_total{destination_service!=\"unknown\"}) or sum by (destination_service) (istio_tcp_received_bytes_total{destination_service!=\"unknown\"})",
						promqlVar.Datasource(datasource),
						promqlVar.LabelName("destination_service"),
					),
					listVar.DisplayName("Service"),
					listVar.AllowAllValue(false),
					listVar.AllowMultiple(true),
				),
			),
			// Reporter variable - matches JSON exactly
			dashboard.AddVariable("qrep",
				listVar.List(
					promqlVar.PrometheusPromQL(
						"group by (reporter) (istio_requests_total) or group by (reporter) (istio_tcp_sent_bytes_total)",
						promqlVar.Datasource(datasource),
						promqlVar.LabelName("reporter"),
					),
					listVar.DisplayName("Reporter"),
					listVar.DefaultValue("destination"),
					listVar.AllowMultiple(true),
				),
			),
			// Optional client/service workload All filters are omitted: Perses/OpenShift
			// often expands All to an empty matcher and panels show No Data.
			withGeneralSection(),
			withGeneralSectionII(datasource, clusterLabelMatcher),
			withClientWorkloadsSection(datasource, clusterLabelMatcher),
			withClientWorkloadsIISection(datasource, clusterLabelMatcher),
			withClientWorkloadsIIISection(datasource, clusterLabelMatcher),
			withServiceWorkloadsSection(datasource, clusterLabelMatcher),
			withServiceWorkloadsIISection(datasource, clusterLabelMatcher),
			withServiceWorkloadsIIISection(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
