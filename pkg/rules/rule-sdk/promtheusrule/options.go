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
