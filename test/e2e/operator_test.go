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

package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOperatorDashboardReconciliation(t *testing.T) {
	requireOperatorE2E(t)

	t.Run("operator deployment is ready", func(t *testing.T) {
		err := pollCondition(5*time.Minute, func() error {
			deploy, err := kubeClient.AppsV1().Deployments(testNamespace).Get(context.Background(), operatorDeploymentName, metav1.GetOptions{})
			if err != nil {
				return err
			}
			if deploy.Status.ReadyReplicas != *deploy.Spec.Replicas {
				return fmt.Errorf("expecting %d ready replicas, got %d", *deploy.Spec.Replicas, deploy.Status.ReadyReplicas)
			}
			return nil
		})
		if err != nil {
			t.Fatal(formatOperatorDebug(err))
		}
	})

	t.Run("perses instance is ready", func(t *testing.T) {
		err := pollCondition(5*time.Minute, func() error {
			return persesInstanceReady(kubeClient)
		})
		if err != nil {
			t.Fatal(formatOperatorDebug(err))
		}
	})

	t.Run("perses CR status is available and not degraded", func(t *testing.T) {
		err := pollCondition(5*time.Minute, func() error {
			conditions, err := getPersesCRConditions(kubeClient)
			if err != nil {
				return err
			}
			return assertReconciledConditions(conditions, "Perses/"+persesInstanceName)
		})
		if err != nil {
			t.Fatal(formatOperatorDebug(err))
		}
	})

	t.Run("perses service responds", func(t *testing.T) {
		err := pollCondition(5*time.Minute, func() error {
			proxyName, err := persesPodProxyName(kubeClient)
			if err != nil {
				return err
			}

			_, err = kubeClient.CoreV1().RESTClient().Get().
				Namespace(testNamespace).
				Resource("pods").
				Name(proxyName).
				SubResource("proxy").
				Suffix("/api/v1/health").
				DoRaw(context.Background())
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("dashboards reconcile", func(t *testing.T) {
		expected, err := countBuiltOperatorDashboards()
		if err != nil {
			t.Fatalf("count built operator dashboards: %v", err)
		}

		err = pollCondition(5*time.Minute, func() error {
			dashboards, err := listPersesDashboardsWithStatus(kubeClient)
			if err != nil {
				return err
			}
			if len(dashboards) != expected {
				return fmt.Errorf("expected %d PersesDashboard objects, got %d", expected, len(dashboards))
			}

			for _, dashboard := range dashboards {
				if err := assertReconciledConditions(dashboard.Conditions, "PersesDashboard/"+dashboard.Name); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(formatOperatorDebug(err))
		}
	})
}

func TestOperatorE2ESkipsWithoutEnv(t *testing.T) {
	if os.Getenv("OPERATOR_E2E") == "true" {
		t.Skip("OPERATOR_E2E=true, covered by reconciliation tests")
	}
	requireOperatorE2E(t)
}
