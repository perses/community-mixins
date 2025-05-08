package main

import (
	"flag"

	dashboards "github.com/perses/community-dashboards/pkg/dashboards"
	"github.com/perses/community-dashboards/pkg/dashboards/alertmanager"
	"github.com/perses/community-dashboards/pkg/dashboards/blackbox"
	k8sComputeResources "github.com/perses/community-dashboards/pkg/dashboards/kubernetes/compute_resources"
	nodeexporter "github.com/perses/community-dashboards/pkg/dashboards/node_exporter"
	"github.com/perses/community-dashboards/pkg/dashboards/perses"
	"github.com/perses/community-dashboards/pkg/dashboards/prometheus"
	"github.com/perses/community-dashboards/pkg/dashboards/thanos"
)

var (
	project          string
	datasource       string
	clusterLabelName string
)

func main() {

	flag.StringVar(&project, "project", "default", "The project name")
	flag.StringVar(&datasource, "datasource", "", "The datasource name")
	flag.StringVar(&clusterLabelName, "cluster-label-name", "", "The cluster label name")
	flag.Parse()

	dashboardWriter := dashboards.NewDashboardWriter()

	dashboardWriter.Add(perses.BuildPersesOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(prometheus.BuildPrometheusOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(prometheus.BuildPrometheusRemoteWrite(project, datasource, clusterLabelName))
	dashboardWriter.Add(nodeexporter.BuildNodeExporterNodes(project, datasource, clusterLabelName))
	dashboardWriter.Add(nodeexporter.BuildNodeExporterClusterUseMethod(project, datasource, clusterLabelName))
	dashboardWriter.Add(alertmanager.BuildAlertManagerOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(thanos.BuildThanosReceiveOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(thanos.BuildThanosQueryOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(thanos.BuildThanosStoreOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(thanos.BuildThanosRulerOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(thanos.BuildThanosQueryFrontendOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(thanos.BuildThanosCompactOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(blackbox.BuildBlackboxExporter(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesNodeResourcesOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesClusterOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesNamespaceOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesPodOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesWorkloadOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesWorkloadNamespaceOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesMultiClusterOverview(project, datasource, clusterLabelName))

	dashboardWriter.Write()
}
