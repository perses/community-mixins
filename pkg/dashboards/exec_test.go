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

package dashboards

import (
	"strings"
	"testing"

	operatorv2 "github.com/perses/perses-operator/api/v1alpha2"
	"github.com/perses/perses/go-sdk/dashboard"
	k8syaml "sigs.k8s.io/yaml"
)

func TestBuilderToOperatorResource_v1alpha2(t *testing.T) {
	builder, err := dashboard.New("test-dashboard",
		dashboard.ProjectName("perses-dev"),
		dashboard.Name("Test Dashboard"),
	)
	if err != nil {
		t.Fatalf("dashboard.New() returned error: %v", err)
	}

	obj := builderToOperatorResource(builder)
	cr, ok := obj.(*operatorv2.PersesDashboard)
	if !ok {
		t.Fatalf("expected *operatorv2.PersesDashboard, got %T", obj)
	}

	if cr.APIVersion != "perses.dev/v1alpha2" {
		t.Errorf("APIVersion = %q, want perses.dev/v1alpha2", cr.APIVersion)
	}
	if cr.Kind != "PersesDashboard" {
		t.Errorf("Kind = %q, want PersesDashboard", cr.Kind)
	}
	if cr.Spec.Config.Display == nil || cr.Spec.Config.Display.Name != "Test Dashboard" {
		t.Errorf("spec.config.display.name = %q, want Test Dashboard", cr.Spec.Config.Display.Name)
	}
	if cr.Labels["app.kubernetes.io/part-of"] != "perses-operator" {
		t.Errorf("labels[app.kubernetes.io/part-of] = %q, want perses-operator", cr.Labels["app.kubernetes.io/part-of"])
	}

	yamlOutput, err := k8syaml.Marshal(cr)
	if err != nil {
		t.Fatalf("failed to marshal operator resource: %v", err)
	}
	output := string(yamlOutput)
	if !strings.Contains(output, "apiVersion: perses.dev/v1alpha2") {
		t.Errorf("yaml missing v1alpha2 apiVersion:\n%s", output)
	}
	if !strings.Contains(output, "config:") {
		t.Errorf("yaml missing spec.config wrapper:\n%s", output)
	}
	if strings.Contains(output, "perses.dev/v1alpha1") {
		t.Errorf("yaml should not contain v1alpha1:\n%s", output)
	}
}
