package kubernetes

var (
	CADVISOR_MATCHER           = "job=\"cadvisor\""
	KUBE_STATE_METRICS_MATCHER = "job=\"kube-state-metrics\""
	KUBELET_MATCHER            = "job=\"kubelet\""
)

func GetCAdvisorMatcher() string {
	return CADVISOR_MATCHER
}

func GetKubeStateMetricsMatcher() string {
	return KUBE_STATE_METRICS_MATCHER
}

func GetKubeletMatcher() string {
	return KUBELET_MATCHER
}

func SetCAdvisorMatcher(matcher string) {
	CADVISOR_MATCHER = matcher
}

func SetKubeStateMetricsMatcher(matcher string) {
	KUBE_STATE_METRICS_MATCHER = matcher
}

func SetKubeletMatcher(matcher string) {
	KUBELET_MATCHER = matcher
}
