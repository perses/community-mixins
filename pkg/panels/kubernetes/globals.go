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

	CADVISOR_MATCHER           = "job=\"" + CADVISOR_LABEL_VALUE + "\""
	KUBE_STATE_METRICS_MATCHER = "job=\"" + KUBE_STATE_METRICS_LABEL_VALUE + "\""
	KUBELET_MATCHER            = "job=\"" + KUBELET_LABEL_VALUE + "\""
	NODE_EXPORTER_MATCHER      = "job=\"" + NODE_EXPORTER_LABEL_VALUE + "\""
	CONTROLLER_MANAGER_MATCHER = "job=\"" + CONTROLLER_MANAGER_LABEL_VALUE + "\""
	KUBE_SCHEDULER_MATCHER     = "job=\"" + KUBE_SCHEDULER_LABEL_VALUE + "\""
	KUBE_PROXY_MATCHER         = "job=\"" + KUBE_PROXY_LABEL_VALUE + "\""
	API_SERVER_MATCHER         = "job=\"" + API_SERVER_LABEL_VALUE + "\""
)

// GetCAdvisorMatcher returns the matcher for the cadvisor job.
func GetCAdvisorMatcher() string {
	return CADVISOR_MATCHER
}

// GetKubeStateMetricsMatcher returns the matcher for the kube-state-metrics job.
func GetKubeStateMetricsMatcher() string {
	return KUBE_STATE_METRICS_MATCHER
}

// GetKubeletMatcher returns the matcher for the kubelet job.
func GetKubeletMatcher() string {
	return KUBELET_MATCHER
}

// GetAPIServerMatcher returns the matcher for the api server job.
func GetAPIServerMatcher() string {
	return API_SERVER_MATCHER
}

// GetNodeExporterMatcher returns the matcher for the node-exporter job.
func GetNodeExporterMatcher() string {
	return NODE_EXPORTER_MATCHER
}

// GetControllerManagerMatcher returns the matcher for the controller-manager job.
func GetControllerManagerMatcher() string {
	return CONTROLLER_MANAGER_MATCHER
}

// GetSchedulerMatcher returns the matcher for the scheduler job.
func GetSchedulerMatcher() string {
	return KUBE_SCHEDULER_MATCHER
}

// GetKubeProxyMatcher returns the matcher for the kube-proxy job.
func GetKubeProxyMatcher() string {
	return KUBE_PROXY_MATCHER
}

// SetCAdvisorMatcher sets the matcher for the cadvisor job globally.
func SetCAdvisorMatcher(matcher string) {
	CADVISOR_MATCHER = matcher
}

// SetKubeStateMetricsMatcher sets the matcher for the kube-state-metrics job globally.
func SetKubeStateMetricsMatcher(matcher string) {
	KUBE_STATE_METRICS_MATCHER = matcher
}

// SetKubeletMatcher sets the matcher for the kubelet job globally.
func SetKubeletMatcher(matcher string) {
	KUBELET_MATCHER = matcher
}

// SetNodeExporterMatcher sets the matcher for the node-exporter job globally.
func SetNodeExporterMatcher(matcher string) {
	NODE_EXPORTER_MATCHER = matcher
}

// SetControllerManagerMatcher sets the matcher for the controller-manager job globally.
func SetControllerManagerMatcher(matcher string) {
	CONTROLLER_MANAGER_MATCHER = matcher
}

// SetSchedulerMatcher sets the matcher for the scheduler job globally.
func SetSchedulerMatcher(matcher string) {
	KUBE_SCHEDULER_MATCHER = matcher
}

// SetKubeProxyMatcher sets the matcher for the kube-proxy job globally.
func SetKubeProxyMatcher(matcher string) {
	KUBE_PROXY_MATCHER = matcher
}

// SetKubeletMatcher sets the matcher for the api server job globally.
func SetAPIServeMatcher(matcher string) {
	API_SERVER_MATCHER = matcher
}

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
