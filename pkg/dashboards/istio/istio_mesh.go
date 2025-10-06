package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/istio"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/prometheus/prometheus/model/labels"
)

func withMeshOverview(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Global Traffic",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(6),
		panels.GlobalRequestVolume(datasource, labelMatcher),
		panels.GlobalSuccessRate(datasource, labelMatcher),
		panels.Global4xxRate(datasource, labelMatcher),
		panels.Global5xxRate(datasource, labelMatcher),
	)
}

func withMeshWorkloads(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Global Traffic",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(16),
		panels.HTTPGRPCWorkloads(datasource, labelMatcher),
		panels.TCPServices(datasource, labelMatcher),
	)
}

func withIstioComponentVersions(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Istio Component Versions",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.IstioComponentVersions(datasource, labelMatcher),
	)
}

func BuildIstioMesh(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := &labels.Matcher{}
	return dashboards.NewDashboardResult(
		dashboard.New("istio-mesh",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Mesh Dashboard"),
			withMeshOverview(datasource, emptyLabelMatcher),
			withMeshWorkloads(datasource, emptyLabelMatcher),
			withIstioComponentVersions(datasource, emptyLabelMatcher),
		),
	).Component("istio")
}
