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
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPersesFormatDashboardsApplied(t *testing.T) {
	requirePersesE2E(t)

	t.Run("perses instance is ready", func(t *testing.T) {
		err := pollCondition(5*time.Minute, func() error {
			return persesInstanceReady(kubeClient)
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("perses service responds", func(t *testing.T) {
		err := pollCondition(5*time.Minute, func() error {
			_, err := persesAPIGet(kubeClient, "/api/v1/health")
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("perses-format dashboards are listed via API", func(t *testing.T) {
		expected, err := countBuiltPersesDashboards()
		if err != nil {
			t.Fatalf("count built Perses dashboards: %v", err)
		}

		err = pollCondition(5*time.Minute, func() error {
			names, err := listPersesAPIDashboards(kubeClient, persesFormatProject)
			if err != nil {
				return err
			}
			if len(names) != expected {
				return fmt.Errorf("expected %d Perses API dashboards in project %q, got %d", expected, persesFormatProject, len(names))
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestPersesE2ESkipsWithoutEnv(t *testing.T) {
	if os.Getenv("PERSES_E2E") == "true" {
		t.Skip("PERSES_E2E=true, covered by apply tests")
	}
	requirePersesE2E(t)
}
