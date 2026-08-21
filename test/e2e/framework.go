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
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
)

const (
	testNamespace              = "perses-dev"
	persesInstanceName         = "perses"
	operatorDeploymentName     = "perses-operator-controller-manager"
	operatorPodLabelSelector   = "control-plane=controller-manager"
	operatorManagerContainer   = "manager"
	persesAvailableCondition   = "Available"
	persesDegradedCondition    = "Degraded"
)

var kubeClient kubernetes.Interface

func builtOperatorDashboardDir() string {
	return filepath.Join("..", "..", "built", "dashboards", "operator")
}

func countBuiltOperatorDashboards() (int, error) {
	dir := builtOperatorDashboardDir()
	var count int
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			count++
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk %s: %w", dir, err)
	}
	if count == 0 {
		return 0, fmt.Errorf("no operator dashboard YAML found under %s", dir)
	}
	return count, nil
}

func pollCondition(timeout time.Duration, conditionFunc func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var conditionErr error
	if err := wait.PollUntilContextCancel(ctx, 5*time.Second, true, func(context.Context) (bool, error) {
		conditionErr = conditionFunc()
		return conditionErr == nil, nil
	}); err != nil {
		return fmt.Errorf("%w: %w", err, conditionErr)
	}

	return nil
}

func requireOperatorE2E(t *testing.T) {
	t.Helper()
	if os.Getenv("OPERATOR_E2E") != "true" {
		t.Skip("OPERATOR_E2E != true, skipping operator e2e tests")
	}
	if kubeClient == nil {
		t.Fatal("kubeClient is not initialized; set KUBECONFIG when OPERATOR_E2E=true")
	}
}

func persesInstanceReady(kClient kubernetes.Interface) error {
	statefulSet, err := kClient.AppsV1().StatefulSets(testNamespace).Get(context.Background(), persesInstanceName, metav1.GetOptions{})
	if err == nil {
		if statefulSet.Status.ReadyReplicas < 1 {
			return fmt.Errorf("expecting at least 1 ready perses replica, got %d", statefulSet.Status.ReadyReplicas)
		}
		return nil
	}

	deploy, err := kClient.AppsV1().Deployments(testNamespace).Get(context.Background(), persesInstanceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("perses workload not found: %w", err)
	}
	if deploy.Status.ReadyReplicas != *deploy.Spec.Replicas {
		return fmt.Errorf("expecting %d ready perses replicas, got %d", *deploy.Spec.Replicas, deploy.Status.ReadyReplicas)
	}
	return nil
}

func persesServiceEndpointsReady(kClient kubernetes.Interface) error {
	endpoints, err := kClient.CoreV1().Endpoints(testNamespace).Get(context.Background(), persesInstanceName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for _, subset := range endpoints.Subsets {
		if len(subset.Addresses) > 0 {
			return nil
		}
	}

	return fmt.Errorf("perses service has no ready endpoints")
}

type persesDashboardItem struct {
	Name       string
	Conditions []metav1.Condition
}

func listPersesDashboardsWithStatus(kClient kubernetes.Interface) ([]persesDashboardItem, error) {
	raw, err := kClient.CoreV1().RESTClient().Get().
		AbsPath("/apis/perses.dev/v1alpha2/namespaces", testNamespace, "persesdashboards").
		DoRaw(context.Background())
	if err != nil {
		return nil, err
	}

	var list struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
			Status struct {
				Conditions []metav1.Condition `json:"conditions"`
			} `json:"status"`
		} `json:"items"`
	}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}

	dashboards := make([]persesDashboardItem, 0, len(list.Items))
	for _, item := range list.Items {
		dashboards = append(dashboards, persesDashboardItem{
			Name:       item.Metadata.Name,
			Conditions: item.Status.Conditions,
		})
	}

	return dashboards, nil
}

func getPersesCRConditions(kClient kubernetes.Interface) ([]metav1.Condition, error) {
	raw, err := kClient.CoreV1().RESTClient().Get().
		AbsPath("/apis/perses.dev/v1alpha2/namespaces", testNamespace, "perses", persesInstanceName).
		DoRaw(context.Background())
	if err != nil {
		return nil, err
	}

	var perses struct {
		Status struct {
			Conditions []metav1.Condition `json:"conditions"`
		} `json:"status"`
	}
	if err := json.Unmarshal(raw, &perses); err != nil {
		return nil, err
	}

	return perses.Status.Conditions, nil
}

func getCondition(conditions []metav1.Condition, conditionType string) (metav1.Condition, bool) {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition, true
		}
	}
	return metav1.Condition{}, false
}

func assertReconciledConditions(conditions []metav1.Condition, resourceName string) error {
	available, ok := getCondition(conditions, persesAvailableCondition)
	if !ok || available.Status != metav1.ConditionTrue {
		detail := "condition not set"
		if ok {
			detail = fmt.Sprintf("reason=%s message=%s", available.Reason, available.Message)
		}
		return fmt.Errorf("%s: Available!=True (%s)", resourceName, detail)
	}

	degraded, ok := getCondition(conditions, persesDegradedCondition)
	if !ok {
		return fmt.Errorf("%s: Degraded condition not set", resourceName)
	}
	if degraded.Status == metav1.ConditionTrue {
		return fmt.Errorf("%s: Degraded=True reason=%s message=%s", resourceName, degraded.Reason, degraded.Message)
	}
	if degraded.Status != metav1.ConditionFalse {
		return fmt.Errorf("%s: Degraded=%s reason=%s message=%s", resourceName, degraded.Status, degraded.Reason, degraded.Message)
	}

	return nil
}

func getOperatorManagerLogs(kClient kubernetes.Interface, tailLines int64) (string, error) {
	pods, err := kClient.CoreV1().Pods(testNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: operatorPodLabelSelector,
	})
	if err != nil {
		return "", err
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("perses-operator pod not found")
	}

	tail := tailLines
	req := kClient.CoreV1().Pods(testNamespace).GetLogs(pods.Items[0].Name, &corev1.PodLogOptions{
		Container: operatorManagerContainer,
		TailLines: &tail,
	})
	logBytes, err := req.DoRaw(context.Background())
	if err != nil {
		return "", err
	}

	return string(logBytes), nil
}

func formatOperatorDebug(err error) error {
	if kubeClient == nil {
		return err
	}

	logs, logErr := getOperatorManagerLogs(kubeClient, 200)
	if logErr != nil {
		return fmt.Errorf("%w\n\nfailed to fetch perses-operator logs: %v", err, logErr)
	}

	return fmt.Errorf("%w\n\nrecent perses-operator manager logs:\n%s", err, logs)
}
