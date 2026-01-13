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

func withProcessGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Process",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.ZtunnelVersions(datasource, labelMatcher),
		panels.ZtunnelMemoryUsage(datasource, labelMatcher),
		panels.ZtunnelCPUUsage(datasource, labelMatcher),
	)
}

func withNetworkGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Network",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.ZtunnelConnections(datasource, labelMatcher),
		panels.ZtunnelBytesTransmitted(datasource, labelMatcher),
		panels.ZtunnelDNSRequest(datasource, labelMatcher),
	)
}

func withOperationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Operations",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(8),
		panels.ZtunnelXDSConnections(datasource, labelMatcher),
		panels.ZtunnelXDSPushes(datasource, labelMatcher),
		panels.ZtunnelWorkloadManager(datasource, labelMatcher),
	)
}

func BuildIstioZtunnel(project string, datasource string, clusterLabelName string) dashboards.DashboardResult {
	emptyLabelMatcher := &labels.Matcher{}
	return dashboards.NewDashboardResult(
		dashboard.New("istio-ztunnel-dashboard",
			dashboard.ProjectName(project),
			dashboard.Name("Istio Ztunnel Dashboard"),
			withProcessGroup(datasource, emptyLabelMatcher),
			withNetworkGroup(datasource, emptyLabelMatcher),
			withOperationsGroup(datasource, emptyLabelMatcher),
		),
	).Component("istio")
}
