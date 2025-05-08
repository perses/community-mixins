package promql

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
