package promtheusrule

import "github.com/perses/community-mixins/pkg/rules/rule-sdk/rulegroup"

func Name(name string) Option {
	return func(builder *Builder) error {
		builder.Name = name
		return nil
	}
}

func Namespace(namespace string) Option {
	return func(builder *Builder) error {
		builder.Namespace = namespace
		return nil
	}
}

func Labels(labels map[string]string) Option {
	return func(builder *Builder) error {
		builder.Labels = labels
		return nil
	}
}

func Annotations(annotations map[string]string) Option {
	return func(builder *Builder) error {
		builder.Annotations = annotations
		return nil
	}
}

func AddRuleGroup(name string, options ...rulegroup.Option) Option {
	return func(builder *Builder) error {
		ruleGroup, err := rulegroup.New(name, options...)
		if err != nil {
			return err
		}
		builder.Spec.Groups = append(builder.Spec.Groups, ruleGroup.RuleGroup)
		return nil
	}
}
