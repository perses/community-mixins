// Copyright The Perses Authors
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	k8syaml "sigs.k8s.io/yaml"
)

const (
	JSONOutput         = "json"
	YAMLOutput         = "yaml"
	OperatorOutput     = "operator"
	OperatorJSONOutput = "operator-json"
)

type Groups struct {
	Groups []monitoringv1.RuleGroup `json:"groups,omitempty"`
}

func executeRuleBuilder(rule *monitoringv1.PrometheusRule, outputFormat string, outputDir string, errWriter io.Writer) {
	var err error
	var output []byte
	var ext string

	switch outputFormat {
	case YAMLOutput:
		output, err = k8syaml.Marshal(Groups{Groups: rule.Spec.Groups})
		ext = YAMLOutput
	case JSONOutput:
		output, err = json.MarshalIndent(Groups{Groups: rule.Spec.Groups}, "", "  ")
		ext = JSONOutput
	case OperatorOutput:
		output, err = k8syaml.Marshal(rule)
		ext = YAMLOutput
	case OperatorJSONOutput:
		output, err = json.MarshalIndent(rule, "", "  ")
		ext = JSONOutput
	default:
		err = fmt.Errorf("--output must be %q, %q, %q or %q", YAMLOutput, JSONOutput, OperatorOutput, OperatorJSONOutput)
	}

	if err != nil {
		if _, ferr := fmt.Fprint(errWriter, err); ferr != nil {
			panic(fmt.Errorf("failed to write err: %w", err))
		}
		os.Exit(-1)
	}

	// create output directory if not exists
	_, err = os.Stat(outputDir)
	if err != nil && !os.IsNotExist(err) {
		if _, ferr := fmt.Fprint(errWriter, err); ferr != nil {
			panic(fmt.Errorf("failed to write err: %w", err))
		}
		os.Exit(-1)
	}

	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(outputDir, os.ModePerm)
	}

	_ = os.WriteFile(fmt.Sprintf("%s/%s.%s", outputDir, rule.Name, ext), output, os.ModePerm)
}

func NewExec() Exec {
	output := flag.Lookup("output-rules").Value.String()
	outputDir := flag.Lookup("output-rules-dir").Value.String()

	if output == "" || outputDir == "" {
		panic("output-rules and output-rules-dir flags are required for generating rules")
	}

	return Exec{
		outputFormat: output,
		outputDir:    outputDir,
	}
}

type Exec struct {
	outputFormat string
	outputDir    string
}

// BuildDashboard is a helper to print the result of a dashboard builder in stdout and errors to stderr
func (b *Exec) BuildRule(dr RuleResult) {
	if dr.err != nil {
		fmt.Fprint(os.Stderr, dr.err)
		os.Exit(-1)
	}
	executeRuleBuilder(dr.rule, b.outputFormat, path.Join(b.outputDir, dr.component), os.Stdout)
}
