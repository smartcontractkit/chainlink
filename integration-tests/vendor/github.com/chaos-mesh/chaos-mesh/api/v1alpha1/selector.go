// Copyright 2021 Chaos Mesh Authors.
//
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
//

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// LabelSelectorRequirements is list of LabelSelectorRequirement
type LabelSelectorRequirements []metav1.LabelSelectorRequirement

// SelectorMode represents the mode to run chaos action.
type SelectorMode string

const (
	// OneMode represents that the system will do the chaos action on one object selected randomly.
	OneMode SelectorMode = "one"
	// AllMode represents that the system will do the chaos action on all objects
	// regardless of status (not ready or not running pods includes).
	// Use this label carefully.
	AllMode SelectorMode = "all"
	// FixedMode represents that the system will do the chaos action on a specific number of running objects.
	FixedMode SelectorMode = "fixed"
	// FixedPercentMode to specify a fixed % that can be inject chaos action.
	FixedPercentMode SelectorMode = "fixed-percent"
	// RandomMaxPercentMode to specify a maximum % that can be inject chaos action.
	RandomMaxPercentMode SelectorMode = "random-max-percent"
)

// GenericSelectorSpec defines some selectors to select objects.
type GenericSelectorSpec struct {
	// Namespaces is a set of namespace to which objects belong.
	// +optional
	Namespaces []string `json:"namespaces,omitempty"`

	// Map of string keys and values that can be used to select objects.
	// A selector based on fields.
	// +optional
	FieldSelectors map[string]string `json:"fieldSelectors,omitempty"`

	// Map of string keys and values that can be used to select objects.
	// A selector based on labels.
	// +optional
	LabelSelectors map[string]string `json:"labelSelectors,omitempty"`

	// a slice of label selector expressions that can be used to select objects.
	// A list of selectors based on set-based label expressions.
	// +optional
	ExpressionSelectors LabelSelectorRequirements `json:"expressionSelectors,omitempty" swaggerignore:"true"`

	// Map of string keys and values that can be used to select objects.
	// A selector based on annotations.
	// +optional
	AnnotationSelectors map[string]string `json:"annotationSelectors,omitempty"`
}

// PodSelectorSpec defines the some selectors to select objects.
// If the all selectors are empty, all objects will be used in chaos experiment.
type PodSelectorSpec struct {
	GenericSelectorSpec `json:",inline"`

	// Nodes is a set of node name and objects must belong to these nodes.
	// +optional
	Nodes []string `json:"nodes,omitempty"`

	// Pods is a map of string keys and a set values that used to select pods.
	// The key defines the namespace which pods belong,
	// and the each values is a set of pod names.
	// +optional
	Pods map[string][]string `json:"pods,omitempty"`

	// Map of string keys and values that can be used to select nodes.
	// Selector which must match a node's labels,
	// and objects must belong to these selected nodes.
	// +optional
	NodeSelectors map[string]string `json:"nodeSelectors,omitempty"`

	// PodPhaseSelectors is a set of condition of a pod at the current time.
	// supported value: Pending / Running / Succeeded / Failed / Unknown
	// +optional
	PodPhaseSelectors []string `json:"podPhaseSelectors,omitempty"`
}

func (in *PodSelectorSpec) DefaultNamespace(namespace string) {
	if len(in.Namespaces) == 0 {
		in.Namespaces = []string{namespace}
	}
}

type PodSelector struct {
	// Selector is used to select pods that are used to inject chaos action.
	Selector PodSelectorSpec `json:"selector"`

	// Mode defines the mode to run chaos action.
	// Supported mode: one / all / fixed / fixed-percent / random-max-percent
	// +kubebuilder:validation:Enum=one;all;fixed;fixed-percent;random-max-percent
	Mode SelectorMode `json:"mode"`

	// Value is required when the mode is set to `FixedMode` / `FixedPercentMode` / `RandomMaxPercentMode`.
	// If `FixedMode`, provide an integer of pods to do chaos action.
	// If `FixedPercentMode`, provide a number from 0-100 to specify the percent of pods the server can do chaos action.
	// IF `RandomMaxPercentMode`,  provide a number from 0-100 to specify the max percent of pods to do chaos action
	// +optional
	Value string `json:"value,omitempty"`
}

type ContainerSelector struct {
	PodSelector `json:",inline"`

	// ContainerNames indicates list of the name of affected container.
	// If not set, the first container will be injected
	// +optional
	ContainerNames []string `json:"containerNames,omitempty"`
}

// ClusterScoped returns true if the selector selects Pods in the cluster
func (in PodSelectorSpec) ClusterScoped() bool {
	// in fact, this will never happened, will add namespace if it is empty, so len(s.Namespaces) can not be 0,
	// but still add judgentment here for safe
	// https://github.com/chaos-mesh/chaos-mesh/blob/478d00d01bb0f9fb08a1085428a7da8c8f9df4e8/api/v1alpha1/common_webhook.go#L22
	if len(in.Namespaces) == 0 && len(in.Pods) == 0 {
		return true
	}

	return false
}
