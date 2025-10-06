package istio

import (
	"github.com/perses/community-dashboards/pkg/dashboards"
	panels "github.com/perses/community-dashboards/pkg/panels/istio"
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/prometheus/prometheus/model/labels"
)

func withPerformanceNotes(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Performance Dashboard Notes",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(6),
		panels.PerformanceDashboardReadme(datasource, labelMatcher),
	)
}

func withVCPUUsage(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("vCPU Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.VCPUPer1kRPS(datasource, labelMatcher),
		panels.VCPU(datasource, labelMatcher),
	)
}

func withMemoryAndDataRates(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Memory and Data Rates",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.PerformanceMemoryUsage(datasource, labelMatcher),
		panels.BytesTransferred(datasource, labelMatcher),
	)
}

func withIstioComponentVersionsPerf(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Istio Component Versions",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(8),
		panels.IstioComponentsByVersion(datasource, labelMatcher),
	)
}

func withProxyResourceUsage(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Proxy Resource Usage",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(7),
		panels.ProxyMemory(datasource, labelMatcher),
		panels.ProxyVCPU(datasource, labelMatcher),
		panels.ProxyDisk(datasource, labelMatcher),
	)
}

func withIstiodResourceUsage(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Istiod Resource Usage",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(7),
		panels.IstiodMemory(datasource, labelMatcher),
		panels.IstiodVCPU(datasource, labelMatcher),
		panels.IstiodDisk(datasource, labelMatcher),
		panels.IstiodGoroutines(datasource, labelMatcher),
	)
}

func BuildIstioPerformance(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := &labels.Matcher{}
	return dashboards.NewDashboardResult(
		dashboard.New("istio-performance",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Performance Dashboard"),
			withPerformanceNotes(datasource, emptyLabelMatcher),
			withVCPUUsage(datasource, emptyLabelMatcher),
			withMemoryAndDataRates(datasource, emptyLabelMatcher),
			withIstioComponentVersionsPerf(datasource, emptyLabelMatcher),
			withProxyResourceUsage(datasource, emptyLabelMatcher),
			withIstiodResourceUsage(datasource, emptyLabelMatcher),
		),
	).Component("istio")
}
