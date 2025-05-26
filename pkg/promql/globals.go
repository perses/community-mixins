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
