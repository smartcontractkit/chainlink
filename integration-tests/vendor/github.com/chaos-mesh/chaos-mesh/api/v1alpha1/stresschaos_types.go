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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Stress chaos is a chaos to generate plenty of stresses over a collection of pods.

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment

// StressChaos is the Schema for the stresschaos API
type StressChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a time chaos experiment
	Spec StressChaosSpec `json:"spec"`

	// +optional
	// Most recently observed status of the time chaos experiment
	Status StressChaosStatus `json:"status"`
}

var _ InnerObjectWithCustomStatus = (*StressChaos)(nil)
var _ InnerObjectWithSelector = (*StressChaos)(nil)
var _ InnerObject = (*StressChaos)(nil)

// StressChaosSpec defines the desired state of StressChaos
type StressChaosSpec struct {
	ContainerSelector `json:",inline"`

	// Stressors defines plenty of stressors supported to stress system components out.
	// You can use one or more of them to make up various kinds of stresses. At least
	// one of the stressors should be specified.
	// +optional
	Stressors *Stressors `json:"stressors,omitempty"`

	// StressngStressors defines plenty of stressors just like `Stressors` except that it's an experimental
	// feature and more powerful. You can define stressors in `stress-ng` (see also `man stress-ng`) dialect,
	// however not all of the supported stressors are well tested. It maybe retired in later releases. You
	// should always use `Stressors` to define the stressors and use this only when you want more stressors
	// unsupported by `Stressors`. When both `StressngStressors` and `Stressors` are defined, `StressngStressors`
	// wins.
	// +optional
	StressngStressors string `json:"stressngStressors,omitempty"`

	// Duration represents the duration of the chaos action
	// +optional
	Duration *string `json:"duration,omitempty" webhook:"Duration"`
}

// StressChaosStatus defines the observed state of StressChaos
type StressChaosStatus struct {
	ChaosStatus `json:",inline"`
	// Instances always specifies stressing instances
	// +optional
	Instances map[string]StressInstance `json:"instances,omitempty"`
}

// StressInstance is an instance generates stresses
type StressInstance struct {
	// UID is the stress-ng identifier
	// +optional
	UID string `json:"uid"`
	// MemoryUID is the memStress identifier
	// +optional
	MemoryUID string `json:"memoryUid"`
	// StartTime specifies when the stress-ng starts
	// +optional
	StartTime *metav1.Time `json:"startTime"`
	// MemoryStartTime specifies when the memStress starts
	// +optional
	MemoryStartTime *metav1.Time `json:"memoryStartTime"`
}

// Stressors defines plenty of stressors supported to stress system components out.
// You can use one or more of them to make up various kinds of stresses
type Stressors struct {
	// MemoryStressor stresses virtual memory out
	// +optional
	MemoryStressor *MemoryStressor `json:"memory,omitempty"`
	// CPUStressor stresses CPU out
	// +optional
	CPUStressor *CPUStressor `json:"cpu,omitempty"`
}

// Normalize the stressors to comply with stress-ng
func (in *Stressors) Normalize() (string, string, error) {
	cpuStressors := ""
	memoryStressors := ""
	if in.MemoryStressor != nil && in.MemoryStressor.Workers != 0 {
		memoryStressors += fmt.Sprintf(" --workers %d", in.MemoryStressor.Workers)
		if len(in.MemoryStressor.Size) != 0 {
			memoryStressors += fmt.Sprintf(" --size %s", in.MemoryStressor.Size)
		}

		if in.MemoryStressor.Options != nil {
			for _, v := range in.MemoryStressor.Options {
				memoryStressors += fmt.Sprintf(" %v ", v)
			}
		}
	}
	if in.CPUStressor != nil && in.CPUStressor.Workers != 0 {
		cpuStressors += fmt.Sprintf(" --cpu %d", in.CPUStressor.Workers)
		if in.CPUStressor.Load != nil {
			cpuStressors += fmt.Sprintf(" --cpu-load %d",
				*in.CPUStressor.Load)
		}

		if in.CPUStressor.Options != nil {
			for _, v := range in.CPUStressor.Options {
				cpuStressors += fmt.Sprintf(" %v ", v)
			}
		}
	}
	return cpuStressors, memoryStressors, nil
}

// Stressor defines common configurations of a stressor
type Stressor struct {
	// Workers specifies N workers to apply the stressor.
	// Maximum 8192 workers can run by stress-ng
	// +kubebuilder:validation:Maximum=8192
	Workers int `json:"workers"`
}

// MemoryStressor defines how to stress memory out
type MemoryStressor struct {
	Stressor `json:",inline"`

	// Size specifies N bytes consumed per vm worker, default is the total available memory.
	// One can specify the size as % of total available memory or in units of B, KB/KiB,
	// MB/MiB, GB/GiB, TB/TiB.
	// +optional
	Size string `json:"size,omitempty" webhook:"Bytes"`

	// extend stress-ng options
	// +optional
	Options []string `json:"options,omitempty"`
}

// CPUStressor defines how to stress CPU out
type CPUStressor struct {
	Stressor `json:",inline"`
	// Load specifies P percent loading per CPU worker. 0 is effectively a sleep (no load) and 100
	// is full loading.
	// +optional
	Load *int `json:"load,omitempty"`

	// extend stress-ng options
	// +optional
	Options []string `json:"options,omitempty"`
}

func (obj *StressChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.ContainerSelector,
	}
}

func (obj *StressChaos) GetCustomStatus() interface{} {
	return &obj.Status.Instances
}
