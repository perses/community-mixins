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

func withControlPlaneResources(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Resource Usage",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panels.MemoryUsage(datasource, labelMatcher),
		panels.MemoryAllocations(datasource, labelMatcher),
		panels.CPUUsage(datasource, labelMatcher),
		panels.Goroutines(datasource, labelMatcher),
	)
}

func withPushInformation(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
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

func withDeployedVersions(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Deployed Versions",
		panelgroup.PanelsPerLine(1),
		panelgroup.PanelHeight(5),
		panels.PilotVersions(datasource, labelMatcher),
	)
}

func withWebhooks(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Webhooks",
		panelgroup.PanelsPerLine(2),
		panelgroup.PanelHeight(8),
		panels.Validation(datasource, labelMatcher),
		panels.Injection(datasource, labelMatcher),
	)
}

func BuildIstioControlPlane(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := &labels.Matcher{}
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
