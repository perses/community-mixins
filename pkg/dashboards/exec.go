package dashboards

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	persesv1 "github.com/perses/perses-operator/api/v1alpha1"
	"github.com/perses/perses/go-sdk/dashboard"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "sigs.k8s.io/yaml"
)

const (
	JSONOutput     = "json"
	YAMLOutput     = "yaml"
	OperatorOutput = "operator"
)

func init() {
	flag.String("output", YAMLOutput, "output format of the exec")
	flag.String("output-dir", "./dist", "output directory of the exec")
}

func executeDashboardBuilder(builder dashboard.Builder, outputFormat string, outputDir string, errWriter io.Writer) {
	var err error
	var output []byte
	var ext string

	switch outputFormat {
	case YAMLOutput:
		output, err = yaml.Marshal(builder.Dashboard)
		ext = YAMLOutput
	case JSONOutput:
		output, err = json.Marshal(builder.Dashboard)
		ext = JSONOutput
	case OperatorOutput:
		op := persesv1.PersesDashboard{
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

		output, err = k8syaml.Marshal(op)
		ext = YAMLOutput
	default:
		err = fmt.Errorf("--output must be %q, %q or %q", JSONOutput, YAMLOutput, OperatorOutput)
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

func NewExec() Exec {
	output := flag.Lookup("output").Value.String()
	outputDir := flag.Lookup("output-dir").Value.String()

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
func (b *Exec) BuildDashboard(builder dashboard.Builder, err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(-1)
	}
	executeDashboardBuilder(builder, b.outputFormat, b.outputDir, os.Stdout)
}
