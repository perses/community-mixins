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

package promql

import "github.com/prometheus/prometheus/model/labels"

var NamespaceVar LabelMatcher = LabelMatcher{
	Name:  "namespace",
	Value: "$namespace",
	Type:  "=",
}

var JobVar LabelMatcher = LabelMatcher{
	Name:  "job",
	Value: "$job",
	Type:  "=~",
}

var InstanceVar LabelMatcher = LabelMatcher{
	Name:  "instance",
	Value: "$instance",
	Type:  "=~",
}

var ClusterVar LabelMatcher = LabelMatcher{
	Name:  "cluster",
	Value: "$cluster",
	Type:  "=~",
}

var NamespaceVarV2 *labels.Matcher = &labels.Matcher{
	Name:  "namespace",
	Value: "$namespace",
	Type:  labels.MatchEqual,
}

var JobVarV2 *labels.Matcher = &labels.Matcher{
	Name:  "job",
	Value: "$job",
	Type:  labels.MatchRegexp,
}

var InstanceVarV2 *labels.Matcher = &labels.Matcher{
	Name:  "instance",
	Value: "$instance",
	Type:  labels.MatchRegexp,
}

var ClusterVarV2 *labels.Matcher = &labels.Matcher{
	Name:  "cluster",
	Value: "$cluster",
	Type:  labels.MatchRegexp,
}
