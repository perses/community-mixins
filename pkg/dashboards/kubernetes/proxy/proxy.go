package proxy

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panelsGostats "github.com/perses/community-dashboards/pkg/panels/gostats"
	panels "github.com/perses/community-dashboards/pkg/panels/kubernetes"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	listVar "github.com/perses/perses/go-sdk/variable/list-variable"
	labelValuesVar "github.com/perses/plugins/prometheus/sdk/go/variable/label-values"
)

func withProxyStatsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Proxy Status",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.ProxyUpStatus(datasource, labelMatcher),
	)
}

func withProxyRulesSyncRateGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Proxy Rules Sync Rate",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.RulesSyncRate(datasource, labelMatcher),
		panels.RulesSyncLatency(datasource, labelMatcher),
	)
}

func withProxyNetworkProgrammingRateGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Proxy Network Programming Rate",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.NetworkProgrammingRate(datasource, labelMatcher),
		panels.NetworkProgrammingLatency(datasource, labelMatcher),
	)
}

func withProxyKubeAPIRequestsGroup(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []promql.LabelMatcher{
		promql.ClusterVar,
		promql.InstanceVar,
		{
			Name:  "job",
			Value: "kube-proxy",
			Type:  "=",
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, labelMatcher)

	return dashboard.AddPanelGroup("Kube API Requests",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.KubeAPIRequestRate(datasource, labelMatchersToUse...),
		panels.PostRequestLatency(datasource, labelMatchersToUse...),
		panels.GetRequestLatency(datasource, labelMatchersToUse...),
	)
}

func withProxyResources(datasource string, clusterLabelMatcher promql.LabelMatcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []promql.LabelMatcher{
		promql.ClusterVar,
		promql.InstanceVar,
		{
			Name:  "job",
			Value: "kube-proxy",
			Type:  "=",
		},
	}

	labelMatchersToUse = append(labelMatchersToUse, clusterLabelMatcher)

	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panelsGostats.MemoryUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.CPUUsage(datasource, "instance", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "instance", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "instance", labelMatchersToUse...),
	)
}

func BuildProxyOverview(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	return dashboards.NewDashboardResult(
		dashboard.New("proxy-overview",
			dashboard.ProjectName(project),
			dashboard.Name("Kubernetes / Proxy"),
			dashboard.AddVariable("cluster",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("cluster",
						labelValuesVar.Matchers("up{"+panels.GetKubeProxyMatcher()+"}"),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("cluster"),
				),
			),
			dashboard.AddVariable("instance",
				listVar.List(
					labelValuesVar.PrometheusLabelValues("instance",
						labelValuesVar.Matchers(
							promql.SetLabelMatchers(
								"up{"+panels.GetKubeProxyMatcher()+"}",
								[]promql.LabelMatcher{{Name: "cluster", Type: "=", Value: "$cluster"}},
							),
						),
						dashboards.AddVariableDatasource(datasource),
					),
					listVar.DisplayName("instance"),
				),
			),
			withProxyStatsGroup(datasource, clusterLabelMatcher),
			withProxyRulesSyncRateGroup(datasource, clusterLabelMatcher),
			withProxyNetworkProgrammingRateGroup(datasource, clusterLabelMatcher),
			withProxyKubeAPIRequestsGroup(datasource, clusterLabelMatcher),
			withProxyResources(datasource, clusterLabelMatcher),
		),
	).Component("kubernetes")
}
