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
// +chaos-mesh:base

// PhysicalMachine is the Schema for the physical machine API
type PhysicalMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a physical machine
	Spec PhysicalMachineSpec `json:"spec"`
}

// PhysicalMachineSpec defines the desired state of PhysicalMachine
type PhysicalMachineSpec struct {

	// Address represents the address of the physical machine
	Address string `json:"address"`
}

// +kubebuilder:object:root=true

// PhysicalMachineList contains a list of PhysicalMachine
type PhysicalMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PhysicalMachine `json:"items"`
}
