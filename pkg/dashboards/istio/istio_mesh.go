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
