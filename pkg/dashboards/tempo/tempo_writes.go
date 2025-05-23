package tempo

import (
	panels "github.com/perses/community-dashboards/pkg/panels/tempo"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"

	"github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/promql"
)

func withWritesGateway(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Gateway",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesGatewayQPS(datasource, labelMatcher),
		panels.WritesGatewayLatency(datasource, labelMatcher),
	)
}

func withWritesEnvoyProxy(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Envoy Proxy",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesEnvoyProxyQPS(datasource, labelMatcher),
		panels.WritesEnvoygRPCStatusCodes(datasource, labelMatcher),
	)
}

func withWritesDistributor(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Distributor",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.WritesDistributorSpansSecond(datasource, labelMatcher),
		panels.WritesDistributorBytesPerSecond(datasource, labelMatcher),
		panels.WritesDistributorLatency(datasource, labelMatcher),
	)
}

func withWritesKafkaProducedRecords(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Kafka produced records",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesDistributorKafkaAppendRecords(datasource, labelMatcher),
		panels.WritesDistributorKafkaAppendFail(datasource, labelMatcher),
	)
}

func withWritesKafkaWrites(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Kafka Writes",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesDistributorKafkaWrite(datasource, labelMatcher),
		panels.WritesDistributorKafkaWriteLatency(datasource, labelMatcher),
	)
}

func withWritesIngester(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Ingester",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesIngesterQPS(datasource, labelMatcher),
		panels.WritesIngesterLatency(datasource, labelMatcher),
	)
}

func withWritesMemcachedIngester(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memcached - Ingester",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesMemcachedIngesterQPS(datasource, labelMatcher),
		panels.WritesMemcachedIngesterLatency(datasource, labelMatcher),
	)
}

func withWritesBackendIngester(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Backend - Ingester",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesBackendIngesterQPS(datasource, labelMatcher),
		panels.WritesBackendIngesterLatency(datasource, labelMatcher),
	)
}

func withWritesMemcachedCompactor(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memcached - Compactor",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesMemcachedCompactorQPS(datasource, labelMatcher),
		panels.WritesMemcachedCompactorLatency(datasource, labelMatcher),
	)
}

func withWritesBackendCompactor(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Backend - Compactor",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WritesBackendCompactorQPS(datasource, labelMatcher),
		panels.WritesBackendCompactorLatency(datasource, labelMatcher),
	)
}

func BuildTempoWritesOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("tempo-writes-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Tempo / Writes"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("tempo_build_info"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
					listVar.AllowMultiple(true),
				),
			),
			dashboards.AddClusterVariable(datasource, clusterLabelName, "tempo_build_info"),
			dashboard.AddVariable("namespace",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("namespace",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"tempo_build_info",
								[]promql.LabelMatcher{clusterLabelMatcher, {Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("namespace"),
				),
			),
			withWritesGateway(datasource, clusterLabelMatcher),
			withWritesEnvoyProxy(datasource, clusterLabelMatcher),
			withWritesDistributor(datasource, clusterLabelMatcher),
			withWritesKafkaProducedRecords(datasource, clusterLabelMatcher),
			withWritesKafkaWrites(datasource, clusterLabelMatcher),
			withWritesIngester(datasource, clusterLabelMatcher),
			withWritesMemcachedIngester(datasource, clusterLabelMatcher),
			withWritesBackendIngester(datasource, clusterLabelMatcher),
			withWritesMemcachedCompactor(datasource, clusterLabelMatcher),
			withWritesBackendCompactor(datasource, clusterLabelMatcher),
		),
	).Component("tempo")
}
