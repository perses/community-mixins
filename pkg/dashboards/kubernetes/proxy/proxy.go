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
	"github.com/prometheus/prometheus/model/labels"
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
			Value: panels.KUBE_PROXY_LABEL_VALUE,
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

func withProxyResources(datasource string, clusterLabelMatcher *labels.Matcher) dashboard.Option {
	// TODO(saswatamcode): Add a way to configure these.
	labelMatchersToUse := []*labels.Matcher{
		promql.ClusterVarV2,
		promql.InstanceVarV2,
		{
			Name:  "job",
			Value: panels.KUBE_PROXY_LABEL_VALUE,
			Type:  labels.MatchEqual,
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

func BuildProxyOverview(project string, datasource string, clusterLabelName string, variableOverrides ...dashboard.Option) dashboards.DashboardResult {
	defaultVars := []dashboard.Option{
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
	}

	clusterLabelMatcher := dashboards.GetClusterLabelMatcher(clusterLabelName)
	clusterLabelMatcherV2 := dashboards.GetClusterLabelMatcherV2(clusterLabelName)

	vars := defaultVars
	if len(variableOverrides) > 0 {
		vars = variableOverrides
	}
	options := append([]dashboard.Option{
		dashboard.ProjectName(project),
		dashboard.Name("Kubernetes / Proxy"),
	}, vars...)
	options = append(options,
		withProxyStatsGroup(datasource, clusterLabelMatcher),
		withProxyRulesSyncRateGroup(datasource, clusterLabelMatcher),
		withProxyNetworkProgrammingRateGroup(datasource, clusterLabelMatcher),
		withProxyKubeAPIRequestsGroup(datasource, clusterLabelMatcher),
		withProxyResources(datasource, clusterLabelMatcherV2),
	)
	return dashboards.NewDashboardResult(
		dashboard.New("proxy-overview", options...),
	).Component("kubernetes")
}
