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

func withWasmVMsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Wasm VMs",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WasmVMActive(datasource, labelMatcher),
		panels.WasmVMCreated(datasource, labelMatcher),
	)
}

func withWasmModuleRemoteLoadGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Wasm Module Remote Load",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.WasmRemoteLoadCacheEntry(datasource, labelMatcher),
		panels.WasmRemoteLoadCacheVisit(datasource, labelMatcher),
		panels.WasmRemoteLoadFetch(datasource, labelMatcher),
	)
}

func withWasmProxyResourceUsageGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Proxy Resource Usage",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.WasmProxyMemory(datasource, labelMatcher),
		panels.WasmProxyVCPU(datasource, labelMatcher),
	)
}

func BuildIstioExtension(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := &labels.Matcher{}
	return dashboards.NewDashboardResult(
		dashboard.New("istio-extension-dashboard",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Wasm Extension Dashboard"),
			withWasmVMsGroup(datasource, emptyLabelMatcher),
			withWasmModuleRemoteLoadGroup(datasource, emptyLabelMatcher),
			withWasmProxyResourceUsageGroup(datasource, emptyLabelMatcher),
		),
	).Component("istio")
}
