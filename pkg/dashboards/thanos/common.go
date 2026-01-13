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

package thanos

import (
	"github.com/perses/perses/go-sdk/dashboard"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/prometheus/prometheus/model/labels"

	panelsGostats "github.com/perses/community-mixins/pkg/panels/gostats"
	panels "github.com/perses/community-mixins/pkg/panels/thanos"
	"github.com/perses/community-mixins/pkg/promql"
)

func withThanosResourcesGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	labelMatchersToUse := []*labels.Matcher{
		promql.NamespaceVarV2,
		promql.JobVarV2,
	}
	labelMatchersToUse = append(labelMatchersToUse, labelMatcher)

	return dashboard.AddPanelGroup("Resources",
		panelgroup.PanelsPerLine(4),
		panelgroup.PanelHeight(10),
		panelsGostats.CPUUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.MemoryUsage(datasource, "pod", labelMatchersToUse...),
		panelsGostats.Goroutines(datasource, "pod", labelMatchersToUse...),
		panelsGostats.GarbageCollectionPauseTimeQuantiles(datasource, "pod", labelMatchersToUse...),
	)
}

func withThanosBucketOperationsGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Bucket Operations",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.BucketOperationRate(datasource, labelMatcher),
		panels.BucketOperationErrors(datasource, labelMatcher),
		panels.BucketOperationDurations(datasource, labelMatcher),
	)
}

func withThanosReadGRPCUnaryGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Read gRPC Unary (StoreAPI Info/Labels)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.ReadGRPCUnaryRate(datasource, labelMatcher),
		panels.ReadGRPCUnaryErrors(datasource, labelMatcher),
		panels.ReadGPRCUnaryDurations(datasource, labelMatcher),
	)
}

func withThanosReadGRPCStreamGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Read gRPC Stream (StoreAPI Series/Exemplars)",
		panelgroup.PanelsPerLine(3),
		panelgroup.PanelHeight(10),
		panels.ReadGRPCStreamRate(datasource, labelMatcher),
		panels.ReadGRPCStreamErrors(datasource, labelMatcher),
		panels.ReadGPRCStreamDurations(datasource, labelMatcher),
	)
}

func withThanosBucketUploadGroup(datasource string, labelMatcher *labels.Matcher) dashboard.Option {
	return dashboard.AddPanelGroup("Last Bucket Upload",
		panelgroup.PanelsPerLine(1),
		panels.BucketUploadTable(datasource, labelMatcher),
	)
}
