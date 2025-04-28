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
	builder dashboard.Builder
	err     error
}

func NewDashboardWriter() *DashboardWriter {
	return &DashboardWriter{
		executor: NewExec(),
	}
}

// Add adds a dashboard to the writer.
func (w *DashboardWriter) Add(builder dashboard.Builder, err error) {
	w.dashboardResults = append(w.dashboardResults, DashboardResult{
		builder: builder,
		err:     err,
	})
}

// Write writes the dashboards to the output directory.
func (w *DashboardWriter) Write() {
	for _, result := range w.dashboardResults {
		w.executor.BuildDashboard(result.builder, result.err)
	}
}

// OperatorResources returns the operator resources of the dashboards added to the writer.
func (w *DashboardWriter) OperatorResources() []runtime.Object {
	operatorResources := []runtime.Object{}
	for _, result := range w.dashboardResults {
		operatorResources = append(operatorResources, w.executor.BuildDashboardOperatorResource(result.builder))
	}
	return operatorResources
}
