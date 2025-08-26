package kubernetes

// Globals to emulate customization behaviors https://github.com/kubernetes-monitoring/kubernetes-mixin?tab=readme-ov-file#customising-the-mixin.
var (
	API_SERVER_LABEL_VALUE         = "kube-apiserver"
	KUBELET_LABEL_VALUE            = "kubelet"
	NODE_EXPORTER_LABEL_VALUE      = "node-exporter"
	CONTROLLER_MANAGER_LABEL_VALUE = "kube-controller-manager"
	KUBE_SCHEDULER_LABEL_VALUE     = "kube-scheduler"
	KUBE_PROXY_LABEL_VALUE         = "kube-proxy"
	KUBE_STATE_METRICS_LABEL_VALUE = "kube-state-metrics"
	CADVISOR_LABEL_VALUE           = "cadvisor"

	// Matchers are now computed dynamically via getter functions
)

// GetCAdvisorMatcher returns the matcher for the cadvisor job.
func GetCAdvisorMatcher() string {
	return "job=\"" + CADVISOR_LABEL_VALUE + "\""
}

// GetKubeStateMetricsMatcher returns the matcher for the kube-state-metrics job.
func GetKubeStateMetricsMatcher() string {
	return "job=\"" + KUBE_STATE_METRICS_LABEL_VALUE + "\""
}

// GetKubeletMatcher returns the matcher for the kubelet job.
func GetKubeletMatcher() string {
	return "job=\"" + KUBELET_LABEL_VALUE + "\""
}

// GetAPIServerMatcher returns the matcher for the api server job.
func GetAPIServerMatcher() string {
	return "job=\"" + API_SERVER_LABEL_VALUE + "\""
}

// GetNodeExporterMatcher returns the matcher for the node-exporter job.
func GetNodeExporterMatcher() string {
	return "job=\"" + NODE_EXPORTER_LABEL_VALUE + "\""
}

// GetControllerManagerMatcher returns the matcher for the controller-manager job.
func GetControllerManagerMatcher() string {
	return "job=\"" + CONTROLLER_MANAGER_LABEL_VALUE + "\""
}

// GetSchedulerMatcher returns the matcher for the scheduler job.
func GetSchedulerMatcher() string {
	return "job=\"" + KUBE_SCHEDULER_LABEL_VALUE + "\""
}

// GetKubeProxyMatcher returns the matcher for the kube-proxy job.
func GetKubeProxyMatcher() string {
	return "job=\"" + KUBE_PROXY_LABEL_VALUE + "\""
}

// NOTE: Matcher setter functions have been removed since matchers are now computed dynamically.
// Use the SetXXXLabelValue functions instead to change the underlying label values.

// SetAPIServerLabelValue sets the label value for the api server job globally.
// WARNING: Ensure you only set this to value that represents a Kube API Server specifically.
func SetAPIServerLabelValue(labelValue string) {
	API_SERVER_LABEL_VALUE = labelValue
}

// SetKubeletLabelValue sets the label value for the kubelet job globally.
// WARNING: Ensure you only set this to value that represents a Kubelet specifically.
func SetKubeletLabelValue(labelValue string) {
	KUBELET_LABEL_VALUE = labelValue
}

// SetNodeExporterLabelValue sets the label value for the node-exporter job globally.
// WARNING: Ensure you only set this to value that represents a Node Exporter specifically.
func SetNodeExporterLabelValue(labelValue string) {
	NODE_EXPORTER_LABEL_VALUE = labelValue
}

// SetControllerManagerLabelValue sets the label value for the controller-manager job globally.
// WARNING: Ensure you only set this to value that represents a Controller Manager specifically.
func SetControllerManagerLabelValue(labelValue string) {
	CONTROLLER_MANAGER_LABEL_VALUE = labelValue
}

// SetCadvisorLabelValue sets the label value for the cadvisor job globally.
// WARNING: Ensure you only set this to value that represents a Cadvisor specifically.
func SetCadvisorLabelValue(labelValue string) {
	CADVISOR_LABEL_VALUE = labelValue
}

// SetSchedulerLabelValue sets the label value for the scheduler job globally.
// WARNING: Ensure you only set this to value that represents a Scheduler specifically.
func SetSchedulerLabelValue(labelValue string) {
	KUBE_SCHEDULER_LABEL_VALUE = labelValue
}

// SetKubeProxyLabelValue sets the label value for the kube-proxy job globally.
// WARNING: Ensure you only set this to value that represents a Kube Proxy specifically.
func SetKubeProxyLabelValue(labelValue string) {
	KUBE_PROXY_LABEL_VALUE = labelValue
}

// SetKubeStateMetricsLabelValue sets the label value for the kube-state-metrics job globally.
// WARNING: Ensure you only set this to value that represents a Kube State Metrics specifically.
func SetKubeStateMetricsLabelValue(labelValue string) {
	KUBE_STATE_METRICS_LABEL_VALUE = labelValue
}

// Metrics deprecation considerations: https://github.com/kubernetes-monitoring/kubernetes-mixin?tab=readme-ov-file#metrics-deprecation
var (
	NODE_NS_CPU_SECONDS_RECORDING_RULE = "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_rate5m"
)

// GetNodeNSCPUSecondsRecordingRule returns the recording rule for the node namespace pod container cpu usage.
func GetNodeNSCPUSecondsRecordingRule() string {
	return NODE_NS_CPU_SECONDS_RECORDING_RULE
}

// SetNodeNSCPUSecondsRecordingRule sets the recording rule for the node namespace pod container cpu usage to use old deprecated rule.
// Metrics deprecation considerations: https://github.com/kubernetes-monitoring/kubernetes-mixin?tab=readme-ov-file#metrics-deprecation
func SetNodeNSCPUSecondsRecordingRuleToDeprecated() {
	NODE_NS_CPU_SECONDS_RECORDING_RULE = "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate"
}
