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
