package networking

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"

	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withClusterCurrentRateBytesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Current Rate of Bytes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesCurrentRateOfBytesReceived("cluster-networking", datasource, labelMatcher),
		panels.KubernetesCurrentRateOfBytesTransmitted("cluster-networking", datasource, labelMatcher),
	)
}

func withClusterNetworkingCurrentStatusGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Current Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(10),
		panels.ClusterNetworkingCurrentStatus(datasource, labelMatcher),
	)
}

func withClusterAvgRateOfBytesGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Average Rate of Bytes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesAverageRateOfBytesReceived("cluster-networking", datasource, labelMatcher),
		panels.KubernetesAverageRateOfBytesTransmitted("cluster-networking", datasource, labelMatcher),
	)
}

func withClusterBandwidthGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bandwidth",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceiveBandwidth("cluster-networking", datasource, labelMatcher),
		panels.KubernetesTransmitBandwidth("cluster-networking", datasource, labelMatcher),
	)
}

func withClusterRateOfPacketsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPackets("cluster-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPackets("cluster-networking", datasource, labelMatcher),
	)
}

func withClusterRateOfPacketsDroppedGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Rate of Packets Dropped",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.KubernetesReceivedPacketsDropped("cluster-networking", datasource, labelMatcher),
		panels.KubernetesTransmittedPacketsDropped("cluster-networking", datasource, labelMatcher),
	)
}

func withClusterTCPRetransmitRateGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("TCP Retransmit Rate",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.ClusterTCPRetransmitRate(datasource, labelMatcher),
		panels.ClusterTCPSYNRetransmitRate(datasource, labelMatcher),
	)
}

func BuildKubernetesClusterOverview(project string, datasource string, clusterLabelName string, variableOverrides []dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
		dashboard.AddVariable("cluster",
			listVar.List(
				labelValuesVar.PrometheusLabelValues("cluster",
					labelValuesVar.Matchers("up{"+panels.GetKubeletMatcher()+", metrics_path=\"/metrics/cadvisor\"}"),
					dashboards.AddVariableDatasource(datasource),
				),
				listVar.DisplayName("cluster"),
			),
		),
	}
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)

	vars := defaultVars
	if len(variableOverrides) > 0 {
		vars = variableOverrides
	}
	options := append([]dashboard.Option{
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Networking / Cluster"),
	}, vars...)
	options = append(options,
		withClusterCurrentRateBytesGroup(datasource, clusterLabelMatcher),
		withClusterNetworkingCurrentStatusGroup(datasource, clusterLabelMatcher),
		withClusterAvgRateOfBytesGroup(datasource, clusterLabelMatcher),
		withClusterBandwidthGroup(datasource, clusterLabelMatcher),
		withClusterRateOfPacketsGroup(datasource, clusterLabelMatcher),
		withClusterRateOfPacketsDroppedGroup(datasource, clusterLabelMatcher),
		withClusterTCPRetransmitRateGroup(datasource, clusterLabelMatcher),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("kubernetes-cluster-networking-overview", options...),
	).Component("kubernetes")
}
