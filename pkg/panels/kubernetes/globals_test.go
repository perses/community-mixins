// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import "testing"

func TestSettersUpdateMatchers(t *testing.T) {
	tests := []struct {
		name         string
		setter       func(string)
		getter       func() string
		defaultValue string
		customValue  string
		globalVar    *string
	}{
		{"APIServer", SetAPIServerLabelValue, GetAPIServerMatcher, "kube-apiserver", "my-apiserver", &API_SERVER_LABEL_VALUE},
		{"Kubelet", SetKubeletLabelValue, GetKubeletMatcher, "kubelet", "my-kubelet", &KUBELET_LABEL_VALUE},
		{"NodeExporter", SetNodeExporterLabelValue, GetNodeExporterMatcher, "node-exporter", "my-node-exporter", &NODE_EXPORTER_LABEL_VALUE},
		{"ControllerManager", SetControllerManagerLabelValue, GetControllerManagerMatcher, "kube-controller-manager", "my-cm", &CONTROLLER_MANAGER_LABEL_VALUE},
		{"Scheduler", SetSchedulerLabelValue, GetSchedulerMatcher, "kube-scheduler", "my-scheduler", &KUBE_SCHEDULER_LABEL_VALUE},
		{"KubeProxy", SetKubeProxyLabelValue, GetKubeProxyMatcher, "kube-proxy", "my-proxy", &KUBE_PROXY_LABEL_VALUE},
		{"KubeStateMetrics", SetKubeStateMetricsLabelValue, GetKubeStateMetricsMatcher, "kube-state-metrics", "my-ksm", &KUBE_STATE_METRICS_LABEL_VALUE},
		{"CAdvisor", SetCAdvisorLabelValue, GetCAdvisorMatcher, "cadvisor", "my-cadvisor", &CADVISOR_LABEL_VALUE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := *tt.globalVar
			defer func() { *tt.globalVar = original }()

			wantDefault := `job="` + tt.defaultValue + `"`
			if got := tt.getter(); got != wantDefault {
				t.Errorf("default matcher = %q, want %q", got, wantDefault)
			}

			tt.setter(tt.customValue)

			wantCustom := `job="` + tt.customValue + `"`
			if got := tt.getter(); got != wantCustom {
				t.Errorf("after setter(%q), matcher = %q, want %q", tt.customValue, got, wantCustom)
			}
		})
	}
}
