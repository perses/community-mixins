package kubernetes

// Globals to emulate customization behaviors https://github.com/kubernetes-monitoring/kubernetes-mixin?tab=readme-ov-file#customising-the-mixin.
var (
	CADVISOR_MATCHER           = "job=\"cadvisor\""
	KUBE_STATE_METRICS_MATCHER = "job=\"kube-state-metrics\""
	KUBELET_MATCHER            = "job=\"kubelet\""
	NODE_EXPORTER_MATCHER      = "job=\"node-exporter\""
	CONTROLLER_MANAGER_MATCHER = "job=\"kube-controller-manager\""
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

// GetNodeExporterMatcher returns the matcher for the node-exporter job.
func GetNodeExporterMatcher() string {
	return NODE_EXPORTER_MATCHER
}

// GetControllerManagerMatcher returns the matcher for the controller-manager job.
func GetControllerManagerMatcher() string {
	return CONTROLLER_MANAGER_MATCHER
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
