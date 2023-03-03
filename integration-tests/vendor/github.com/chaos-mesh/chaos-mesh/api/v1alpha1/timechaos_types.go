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

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment

// TimeChaos is the Schema for the timechaos API
type TimeChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a time chaos experiment
	Spec TimeChaosSpec `json:"spec"`

	// +optional
	// Most recently observed status of the time chaos experiment
	Status TimeChaosStatus `json:"status"`
}

var _ InnerObjectWithSelector = (*TimeChaos)(nil)
var _ InnerObject = (*TimeChaos)(nil)

// TimeChaosSpec defines the desired state of TimeChaos
type TimeChaosSpec struct {
	ContainerSelector `json:",inline"`

	// TimeOffset defines the delta time of injected program. It's a possibly signed sequence of decimal numbers, such as
	// "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	TimeOffset string `json:"timeOffset" webhook:"TimeOffset"`

	// ClockIds defines all affected clock id
	// All available options are ["CLOCK_REALTIME","CLOCK_MONOTONIC","CLOCK_PROCESS_CPUTIME_ID","CLOCK_THREAD_CPUTIME_ID",
	// "CLOCK_MONOTONIC_RAW","CLOCK_REALTIME_COARSE","CLOCK_MONOTONIC_COARSE","CLOCK_BOOTTIME","CLOCK_REALTIME_ALARM",
	// "CLOCK_BOOTTIME_ALARM"]
	// Default value is ["CLOCK_REALTIME"]
	ClockIds []string `json:"clockIds,omitempty" webhook:"ClockIds,nilable"`

	// Duration represents the duration of the chaos action
	Duration *string `json:"duration,omitempty"`
}

// TimeChaosStatus defines the observed state of TimeChaos
type TimeChaosStatus struct {
	ChaosStatus `json:",inline"`
}

func (in *TimeChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &in.Spec.ContainerSelector,
	}
}
