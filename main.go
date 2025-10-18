package main

import (
	"flag"

	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/dashboards/alertmanager"
	"github.com/perses/community-mixins/pkg/dashboards/blackbox"
	"github.com/perses/community-mixins/pkg/dashboards/etcd"
	"github.com/perses/community-mixins/pkg/dashboards/istio"
	"github.com/perses/community-mixins/pkg/dashboards/kubernetes/apiserver"
	k8sComputeResources "github.com/perses/community-mixins/pkg/dashboards/kubernetes/compute_resources"
	"github.com/perses/community-mixins/pkg/dashboards/kubernetes/controller_manager"
	"github.com/perses/community-mixins/pkg/dashboards/kubernetes/kubelet"
	k8sNetworking "github.com/perses/community-mixins/pkg/dashboards/kubernetes/networking"
	k8sPersistentVolume "github.com/perses/community-mixins/pkg/dashboards/kubernetes/persistent_volume"
	"github.com/perses/community-mixins/pkg/dashboards/kubernetes/proxy"
	"github.com/perses/community-mixins/pkg/dashboards/kubernetes/scheduler"
	nodeexporter "github.com/perses/community-mixins/pkg/dashboards/node_exporter"
	"github.com/perses/community-mixins/pkg/dashboards/opentelemetry"
	"github.com/perses/community-mixins/pkg/dashboards/perses"
	"github.com/perses/community-mixins/pkg/dashboards/prometheus"
	"github.com/perses/community-mixins/pkg/dashboards/tempo"
	"github.com/perses/community-mixins/pkg/dashboards/thanos"
	"github.com/perses/community-mixins/pkg/rules"
	blackboxrules "github.com/perses/community-mixins/pkg/rules/blackbox"
	thanosrules "github.com/perses/community-mixins/pkg/rules/thanos"
)

var (
	project          string
	datasource       string
	clusterLabelName string
	buildRules       bool
)

func main() {
	flag.StringVar(&project, "project", "default", "The project name")
	flag.StringVar(&datasource, "datasource", "", "The datasource name")
	flag.StringVar(&clusterLabelName, "cluster-label-name", "", "The cluster label name")
	flag.BoolVar(&buildRules, "build-rules", false, "Whether to build rules")

	flag.String("output-rules", rules.YAMLOutput, "output format of the rule exec")
	flag.String("output-rules-dir", "./built/rules", "output directory of the rule exec")

	flag.String("output", dashboards.YAMLOutput, "output format of the dashboard exec")
	flag.String("output-dir", "./built", "output directory of the dashboard exec")

	flag.Parse()

	if buildRules {
		ruleWriter := rules.NewRuleWriter()
		ruleWriter.Add(
			thanosrules.BuildThanosRules(
				project,
				map[string]string{
					"app.kubernetes.io/component": "thanos",
					"app.kubernetes.io/name":      "thanos-rules",
					"app.kubernetes.io/part-of":   "thanos",
					"app.kubernetes.io/version":   "main",
				},
				map[string]string{},
				thanosrules.WithRunbookURL("https://github.com/thanos-io/thanos/blob/main/mixin/runbook.md"),
				thanosrules.WithServiceLabelValue("thanos"),
				thanosrules.WithCompactDashboardURL("https://demo.perses.dev/projects/perses/dashboards/thanoscompact"),
				thanosrules.WithQueryDashboardURL("https://demo.perses.dev/projects/perses/dashboards/thanosquery"),
				thanosrules.WithReceiveDashboardURL("https://demo.perses.dev/projects/perses/dashboards/thanosreceive"),
				thanosrules.WithStoreDashboardURL("https://demo.perses.dev/projects/perses/dashboards/thanosstore"),
				thanosrules.WithRuleDashboardURL("https://demo.perses.dev/projects/perses/dashboards/thanosrule"),
			),
		)

		ruleWriter.Add(
			blackboxrules.BuildBlackboxRules(
				project,
				map[string]string{
					"app.kubernetes.io/component": "blackbox-exporter",
					"app.kubernetes.io/name":      "blackbox-exporter-rules",
					"app.kubernetes.io/part-of":   "blackbox-exporter",
					"app.kubernetes.io/version":   "main",
				},
				map[string]string{},
				blackboxrules.WithDashboardURL("https://demo.perses.dev/projects/perses/dashboards/blackboxexporter"),
			),
		)
		ruleWriter.Write()
	} else {
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
		dashboardWriter.Add(kubelet.BuildKubeletOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(controller_manager.BuildControllerManagerOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(proxy.BuildProxyOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(scheduler.BuildSchedulerOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(k8sNetworking.BuildKubernetesClusterOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(k8sNetworking.BuildKubernetesNamespaceByPodOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(k8sNetworking.BuildKubernetesNamespaceByWorkloadOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(k8sNetworking.BuildKubernetesPodOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(k8sNetworking.BuildKubernetesWorkloadOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(k8sPersistentVolume.BuildKubernetesPersistentVolumeOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(etcd.BuildETCDOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(apiserver.BuildAPIServerOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(tempo.BuildTempoWritesOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(tempo.BuildTempoTenantOverview(project, datasource, clusterLabelName))
		dashboardWriter.Add(opentelemetry.BuildOpenTelemetryCollector(project, datasource, clusterLabelName))
		dashboardWriter.Add(istio.BuildIstioControlPlane(project, datasource, clusterLabelName))
		dashboardWriter.Add(istio.BuildIstioMesh(project, datasource, clusterLabelName))
		dashboardWriter.Add(istio.BuildIstioWorkload(project, datasource, clusterLabelName))
		dashboardWriter.Add(istio.BuildIstioService(project, datasource, clusterLabelName))
		dashboardWriter.Add(istio.BuildIstioPerformance(project, datasource, clusterLabelName))
		dashboardWriter.Add(istio.BuildIstioZtunnel(project, datasource, clusterLabelName))

		dashboardWriter.Write()
	}

}
