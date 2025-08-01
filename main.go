package main

import (
	"flag"
	"github.com/perses/community-dashboards/pkg/dashboards/alertmanager"
	"github.com/perses/community-dashboards/pkg/dashboards/blackbox"
	"github.com/perses/community-dashboards/pkg/dashboards/etcd"
	"github.com/perses/community-dashboards/pkg/dashboards/kubernetes/apiserver"
	"github.com/perses/community-dashboards/pkg/dashboards/kubernetes/controller_manager"
	"github.com/perses/community-dashboards/pkg/dashboards/kubernetes/kubelet"
	k8sNetworking "github.com/perses/community-dashboards/pkg/dashboards/kubernetes/networking"
	k8sPersistentVolume "github.com/perses/community-dashboards/pkg/dashboards/kubernetes/persistent_volume"
	"github.com/perses/community-dashboards/pkg/dashboards/kubernetes/proxy"
	"github.com/perses/community-dashboards/pkg/dashboards/kubernetes/scheduler"
	nodeexporter "github.com/perses/community-dashboards/pkg/dashboards/node_exporter"
	"github.com/perses/community-dashboards/pkg/dashboards/opentelemetry"
	"github.com/perses/community-dashboards/pkg/dashboards/perses"
	"github.com/perses/community-dashboards/pkg/dashboards/prometheus"
	"github.com/perses/community-dashboards/pkg/dashboards/tempo"
	"github.com/perses/community-dashboards/pkg/dashboards/thanos"

	"github.com/perses/community-dashboards/pkg/dashboards"
	k8sComputeResources "github.com/perses/community-dashboards/pkg/dashboards/kubernetes/compute_resources"
	"github.com/perses/perses/go-sdk/dashboard"
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

	customVars := []dashboard.Option{}

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
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesNodeResourcesOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesClusterOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesNamespaceOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesPodOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesWorkloadOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesWorkloadNamespaceOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sComputeResources.BuildKubernetesMultiClusterOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(kubelet.BuildKubeletOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(controller_manager.BuildControllerManagerOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(proxy.BuildProxyOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(scheduler.BuildSchedulerOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sNetworking.BuildKubernetesClusterOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sNetworking.BuildKubernetesNamespaceByPodOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sNetworking.BuildKubernetesNamespaceByWorkloadOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sNetworking.BuildKubernetesPodOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sNetworking.BuildKubernetesWorkloadOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(k8sPersistentVolume.BuildKubernetesPersistentVolumeOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(etcd.BuildETCDOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(apiserver.BuildAPIServerOverview(project, datasource, clusterLabelName, customVars))
	dashboardWriter.Add(tempo.BuildTempoWritesOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(tempo.BuildTempoTenantOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(opentelemetry.BuildOpenTelemetryCollector(project, datasource, clusterLabelName))

	dashboardWriter.Write()
}
