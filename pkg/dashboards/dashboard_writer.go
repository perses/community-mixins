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

package dashboards

import (
	"github.com/perses/perses/go-sdk/dashboard"
	"k8s.io/apimachinery/pkg/runtime"
)

type DashboardWriter struct {
	dashboardResults []DashboardResult
	executor         Exec
}

type DashboardResult struct {
	builder   dashboard.Builder
	component string
	err       error
}

func NewDashboardResult(builder dashboard.Builder, err error) DashboardResult {
	return DashboardResult{
		builder: builder,
		err:     err,
	}
}

// Components sets the component field of the DashboardResult.
// This component field is used by DashboardWriter, as the subdirectory name for the dashboard.
func (d DashboardResult) Component(component string) DashboardResult {
	d.component = component
	return d
}

func NewDashboardWriter() *DashboardWriter {
	return &DashboardWriter{
		executor: NewExec(),
	}
}

// Add adds a dashboard to the writer.
func (w *DashboardWriter) Add(dr DashboardResult) {
	w.dashboardResults = append(w.dashboardResults, dr)
}

// Write writes the dashboards to the output directory.
func (w *DashboardWriter) Write() {
	for _, result := range w.dashboardResults {
		w.executor.BuildDashboard(result)
	}
}

// OperatorResources returns the operator resources of the dashboards added to the writer.
func (w *DashboardWriter) OperatorResources() []runtime.Object {
	operatorResources := []runtime.Object{}
	for _, result := range w.dashboardResults {
		operatorResources = append(operatorResources, w.executor.BuildDashboardOperatorResource(result))
	}
	return operatorResources
}
