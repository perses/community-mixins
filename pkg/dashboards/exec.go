package dashboards

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	persesv1 "github.com/perses/perses-operator/api/v1alpha1"
	"github.com/perses/perses/go-sdk/dashboard"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8syaml "sigs.k8s.io/yaml"
)

const (
	JSONOutput         = "json"
	YAMLOutput         = "yaml"
	OperatorOutput     = "operator"
	OperatorJSONOutput = "operator-json"
)

func executeDashboardBuilder(builder dashboard.Builder, outputFormat string, outputDir string, errWriter io.Writer) {
	var err error
	var output []byte
	var ext string

	switch outputFormat {
	case YAMLOutput:
		output, err = yaml.Marshal(builder.Dashboard)
		ext = YAMLOutput
	case JSONOutput:
		output, err = json.MarshalIndent(builder.Dashboard, "", "  ")
		ext = JSONOutput
	case OperatorOutput:
		output, err = k8syaml.Marshal(builderToOperatorResource(builder))
		ext = YAMLOutput
	case OperatorJSONOutput:
		output, err = json.MarshalIndent(builderToOperatorResource(builder), "", "  ")
		ext = JSONOutput
	default:
		err = fmt.Errorf("--output must be %q, %q, %q or %q", JSONOutput, YAMLOutput, OperatorOutput, OperatorJSONOutput)
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

	_ = os.WriteFile(fmt.Sprintf("%s/%s.%s", outputDir, builder.Dashboard.Metadata.Name, ext), output, os.ModePerm)
}

func builderToOperatorResource(builder dashboard.Builder) runtime.Object {
	return &persesv1.PersesDashboard{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersesDashboard",
			APIVersion: "perses.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      builder.Dashboard.Metadata.Name,
			Namespace: builder.Dashboard.Metadata.Project,
			Labels: map[string]string{
				"app.kubernetes.io/name":      "perses-dashboard",
				"app.kubernetes.io/instance":  builder.Dashboard.Metadata.Name,
				"app.kubernetes.io/part-of":   "perses-operator",
				"app.kubernetes.io/component": "dashboard",
			},
		},
		Spec: persesv1.Dashboard{
			DashboardSpec: builder.Dashboard.Spec,
		},
	}
}

func NewExec() Exec {
	output := flag.Lookup("output").Value.String()
	outputDir := flag.Lookup("output-dir").Value.String()

	if output == "" || outputDir == "" {
		panic("output and output-dir flags are required for generating dashboards")
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
func (b *Exec) BuildDashboard(dr DashboardResult) {
	if dr.err != nil {
		fmt.Fprint(os.Stderr, dr.err)
		os.Exit(-1)
	}
	executeDashboardBuilder(dr.builder, b.outputFormat, path.Join(b.outputDir, dr.component), os.Stdout)
}

// BuildDashboardOperatorResource is a helper to return the operator resource of a dashboard builder as a runtime.Object.
func (b *Exec) BuildDashboardOperatorResource(dr DashboardResult) runtime.Object {
	if dr.err != nil {
		fmt.Fprint(os.Stderr, dr.err)
		os.Exit(-1)
	}
	return builderToOperatorResource(dr.builder)
}
