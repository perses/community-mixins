package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/istio"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	markdownPanel "github.com/perses/plugins/markdown/sdk/go"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

// General section - matches JSON layout exactly
func withGeneralSection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("General",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(4),
		// Service header panel
		panelgroup.AddPanel("SERVICE",
			markdownPanel.Markdown("SERVICE Header",
				markdownPanel.Text("<div class=\"dashboard-header text-center\">\n<span>SERVICE: $service</span>\n</div>"),
			),
		),
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
func withClientWorkloadsSection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
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
func withClientWorkloadsIISection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Client Workloads (II)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(6),
		panels.IncomingRequestDurationByClient(datasource, labelMatcher),
		panels.IncomingRequestSizeByClient(datasource, labelMatcher),
		panels.ResponseSizeByClient(datasource, labelMatcher),
	)
}

// Client Workloads (III) section
func withClientWorkloadsIIISection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Client Workloads (III)",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(6),
		panels.BytesReceivedFromTCPClient(datasource, labelMatcher),
		panels.BytesSentToTCPClient(datasource, labelMatcher),
	)
}

// Service Workloads section
func withServiceWorkloadsSection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
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
func withServiceWorkloadsIISection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Service Workloads (II)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(6),
		panels.IncomingRequestDurationByService(datasource, labelMatcher),
		panels.IncomingRequestSizeByService(datasource, labelMatcher),
		panels.ResponseSizeByService(datasource, labelMatcher),
	)
}

// Service Workloads (III) section
func withServiceWorkloadsIIISection(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Service Workloads (III)",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(6),
		panels.BytesReceivedFromTCPService(datasource, labelMatcher),
		panels.BytesSentToTCPService(datasource, labelMatcher),
	)
}

func BuildIstioService(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("istio-service-dashboard",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Service Dashboard"),
			// Service variable - matches JSON exactly
			dashboard.AddVariable("service",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_service",
						labelValuesVar.Matchers("sum(istio_requests_total{}) by (destination_service) or sum(istio_tcp_sent_bytes_total{}) by (destination_service)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service"),
					listVar.DefaultValue("details.bookinfo.svc.cluster.local"),
				),
			),
			// Reporter variable - matches JSON exactly
			dashboard.AddVariable("qrep",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("reporter",
						labelValuesVar.Matchers("up{job=~\"istio-proxy\"}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Reporter"),
					listVar.DefaultValue("destination"),
					listVar.AllowMultiple(true),
				),
			),
			// Client Cluster variable
			dashboard.AddVariable("srccluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("source_cluster",
						labelValuesVar.Matchers("sum(istio_requests_total{reporter=~\"$qrep\", destination_service=\"$service\"}) by (source_cluster) or sum(istio_tcp_sent_bytes_total{reporter=~\"$qrep\", destination_service=~\"$service\"}) by (source_cluster)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Client Cluster"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			// Client Workload Namespace variable
			dashboard.AddVariable("srcns",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("source_workload_namespace",
						labelValuesVar.Matchers("sum(istio_requests_total{reporter=~\"$qrep\", destination_service=\"$service\"}) by (source_workload_namespace) or sum(istio_tcp_sent_bytes_total{reporter=~\"$qrep\", destination_service=~\"$service\"}) by (source_workload_namespace)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Client Workload Namespace"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			// Client Workload variable
			dashboard.AddVariable("srcwl",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("source_workload",
						labelValuesVar.Matchers("sum(istio_requests_total{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload_namespace=~\"$srcns\"}) by (source_workload) or sum(istio_tcp_sent_bytes_total{reporter=~\"$qrep\", destination_service=~\"$service\", source_workload_namespace=~\"$srcns\"}) by (source_workload)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Client Workload"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			// Service Workload Cluster variable
			dashboard.AddVariable("dstcluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_cluster",
						labelValuesVar.Matchers("sum(istio_requests_total{reporter=\"destination\", destination_service=\"$service\"}) by (destination_cluster) or sum(istio_tcp_sent_bytes_total{reporter=\"destination\", destination_service=~\"$service\"}) by (destination_cluster)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service Workload Cluster"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			// Service Workload Namespace variable
			dashboard.AddVariable("dstns",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_workload_namespace",
						labelValuesVar.Matchers("sum(istio_requests_total{reporter=\"destination\", destination_service=\"$service\"}) by (destination_workload_namespace) or sum(istio_tcp_sent_bytes_total{reporter=\"destination\", destination_service=~\"$service\"}) by (destination_workload_namespace)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service Workload Namespace"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			// Service Workload variable
			dashboard.AddVariable("dstwl",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("destination_workload",
						labelValuesVar.Matchers("sum(istio_requests_total{reporter=\"destination\", destination_service=~\"$service\", destination_cluster=~\"$dstcluster\", destination_workload_namespace=~\"$dstns\"}) by (destination_workload) or sum(istio_tcp_sent_bytes_total{reporter=\"destination\", destination_service=~\"$service\", destination_cluster=~\"$dstcluster\", destination_workload_namespace=~\"$dstns\"}) by (destination_workload)"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("Service Workload"),
					listVar.AllowAllValue(true),
					listVar.AllowMultiple(true),
				),
			),
			// Add all sections that match the JSON layout
			withGeneralSection(datasource, clusterLabelMatcher),
			withClientWorkloadsSection(datasource, clusterLabelMatcher),
			withClientWorkloadsIISection(datasource, clusterLabelMatcher),
			withClientWorkloadsIIISection(datasource, clusterLabelMatcher),
			withServiceWorkloadsSection(datasource, clusterLabelMatcher),
			withServiceWorkloadsIISection(datasource, clusterLabelMatcher),
			withServiceWorkloadsIIISection(datasource, clusterLabelMatcher),
		),
	).Component("istio")
}
