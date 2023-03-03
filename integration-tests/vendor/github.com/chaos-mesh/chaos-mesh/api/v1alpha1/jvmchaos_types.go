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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JVMChaosSpec defines the desired state of JVMChaos
type JVMChaosSpec struct {
	ContainerSelector `json:",inline"`

	// Duration represents the duration of the chaos action
	// +optional
	Duration *string `json:"duration,omitempty" webhook:"Duration"`

	// Action defines the specific jvm chaos action.
	// Supported action: latency;return;exception;stress;gc;ruleData
	// +kubebuilder:validation:Enum=latency;return;exception;stress;gc;ruleData
	Action JVMChaosAction `json:"action"`

	// JVMParameter represents the detail about jvm chaos action definition
	// +optional
	JVMParameter `json:",inline"`
}

// JVMChaosAction represents the chaos action about jvm
type JVMChaosAction string

const (
	// JVMLatencyAction represents the JVM chaos action of invoke latency
	JVMLatencyAction JVMChaosAction = "latency"

	// JVMReturnAction represents the JVM chaos action of return value
	JVMReturnAction JVMChaosAction = "return"

	// JVMExceptionAction represents the JVM chaos action of throwing custom exceptions
	JVMExceptionAction JVMChaosAction = "exception"

	// JVMStressAction represents the JVM chaos action of stress like CPU and memory
	JVMStressAction JVMChaosAction = "stress"

	// JVMGCAction represents the JVM chaos action of trigger garbage collection
	JVMGCAction JVMChaosAction = "gc"

	// JVMRuleDataAction represents inject fault with byteman's rule
	// refer to https://downloads.jboss.org/byteman/4.0.14/byteman-programmers-guide.html#the-byteman-rule-language
	JVMRuleDataAction JVMChaosAction = "ruleData"
)

// JVMParameter represents the detail about jvm chaos action definition
type JVMParameter struct {
	// +optional
	// byteman rule name, should be unique, and will use JVMChaos' name if not set
	Name string `json:"name"`

	// +optional
	// Java class
	Class string `json:"class"`

	// +optional
	// the method in Java class
	Method string `json:"method"`

	// +optional
	// the return value for action 'return'
	ReturnValue string `json:"value"`

	// +optional
	// the exception which needs to throw for action `exception`
	ThrowException string `json:"exception"`

	// +optional
	// the latency duration for action 'latency', unit ms
	LatencyDuration int `json:"latency"`

	// +optional
	// the CPU core number need to use, only set it when action is stress
	CPUCount int `json:"cpuCount"`

	// +optional
	// the memory type need to locate, only set it when action is stress, the value can be 'stack' or 'heap'
	MemoryType string `json:"memType"`

	// +optional
	// the port of agent server, default 9277
	Port int32 `json:"port"`

	// +optional
	RuleData string `json:"ruleData"`
}

// JVMChaosStatus defines the observed state of JVMChaos
type JVMChaosStatus struct {
	ChaosStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="action",type=string,JSONPath=`.spec.action`
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment

// JVMChaos is the Schema for the jvmchaos API
type JVMChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JVMChaosSpec   `json:"spec,omitempty"`
	Status JVMChaosStatus `json:"status,omitempty"`
}

var _ InnerObjectWithSelector = (*JVMChaos)(nil)
var _ InnerObject = (*JVMChaos)(nil)

func init() {
	SchemeBuilder.Register(&JVMChaos{}, &JVMChaosList{})
}

func (obj *JVMChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.ContainerSelector,
	}
}
