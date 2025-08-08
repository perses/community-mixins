package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/istio"
	"github.com/perses/community-dashboards/pkg/promql"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
)

func withControlPlaneResources(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.MemoryUsage(datasource, labelMatcher),
		panels.MemoryAllocations(datasource, labelMatcher),
		panels.CPUUsage(datasource, labelMatcher),
		panels.Goroutines(datasource, labelMatcher),
	)
}

func withPushInformation(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Push Information",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.XDSPushes(datasource, labelMatcher),
		panels.Events(datasource, labelMatcher),
		panels.Connections(datasource, labelMatcher),
		panels.PushErrors(datasource, labelMatcher),
		panels.PushTime(datasource, labelMatcher),
		panels.PushSize(datasource, labelMatcher),
	)
}

func withDeployedVersions(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Deployed Versions",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(5),
		panels.PilotVersions(datasource, labelMatcher),
	)
}

func withWebhooks(datasource string, labelMatcher promql.LabelMatcher) dashboard.Option {
	return dashboard.AddPanelGroup("Webhooks",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.Validation(datasource, labelMatcher),
		panels.Injection(datasource, labelMatcher),
	)
}

func BuildIstioControlPlane(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := promql.LabelMatcher{}
	return dashboards.NewDashboardResult(
		dashboard.New("istio-control-plane",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Control Plane Dashboard"),
			withDeployedVersions(datasource, emptyLabelMatcher),
			withControlPlaneResources(datasource, emptyLabelMatcher),
			withPushInformation(datasource, emptyLabelMatcher),
			withWebhooks(datasource, emptyLabelMatcher),
		),
	).Component("istio")
}
