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
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="action",type=string,JSONPath=`.spec.action`
// +kubebuilder:printcolumn:name="duration",type=string,JSONPath=`.spec.duration`
// +chaos-mesh:experiment
// +chaos-mesh:oneshot=in.Spec.Action==Ec2Restart

// AWSChaos is the Schema for the awschaos API
type AWSChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSChaosSpec   `json:"spec"`
	Status AWSChaosStatus `json:"status,omitempty"`
}

var _ InnerObjectWithSelector = (*AWSChaos)(nil)
var _ InnerObject = (*AWSChaos)(nil)

// AWSChaosAction represents the chaos action about aws.
type AWSChaosAction string

const (
	// Ec2Stop represents the chaos action of stopping ec2.
	Ec2Stop AWSChaosAction = "ec2-stop"
	// Ec2Restart represents the chaos action of restarting ec2.
	Ec2Restart AWSChaosAction = "ec2-restart"
	// DetachVolume represents the chaos action of detaching the volume of ec2.
	DetachVolume AWSChaosAction = "detach-volume"
)

// AWSChaosSpec is the content of the specification for an AWSChaos
type AWSChaosSpec struct {
	// Action defines the specific aws chaos action.
	// Supported action: ec2-stop / ec2-restart / detach-volume
	// Default action: ec2-stop
	// +kubebuilder:validation:Enum=ec2-stop;ec2-restart;detach-volume
	Action AWSChaosAction `json:"action"`

	// Duration represents the duration of the chaos action.
	// +optional
	Duration *string `json:"duration,omitempty" webhook:"Duration"`

	// SecretName defines the name of kubernetes secret.
	// +optional
	SecretName *string `json:"secretName,omitempty" webhook:",nilable"`

	AWSSelector `json:",inline"`
}

// AWSChaosStatus represents the status of an AWSChaos
type AWSChaosStatus struct {
	ChaosStatus `json:",inline"`
}

type AWSSelector struct {
	// TODO: it would be better to split them into multiple different selector and implementation
	// but to keep the minimal modification on current implementation, it hasn't been splited.

	// Endpoint indicates the endpoint of the aws server. Just used it in test now.
	// +ui:form:ignore
	// +optional
	Endpoint *string `json:"endpoint,omitempty"`

	// AWSRegion defines the region of aws.
	AWSRegion string `json:"awsRegion"`

	// Ec2Instance indicates the ID of the ec2 instance.
	Ec2Instance string `json:"ec2Instance"`

	// EbsVolume indicates the ID of the EBS volume.
	// Needed in detach-volume.
	// +ui:form:when=action=='detach-volume'
	// +optional
	EbsVolume *string `json:"volumeID,omitempty" webhook:"EbsVolume,nilable"`

	// DeviceName indicates the name of the device.
	// Needed in detach-volume.
	// +ui:form:when=action=='detach-volume'
	// +optional
	DeviceName *string `json:"deviceName,omitempty" webhook:"AWSDeviceName,nilable"`
}

func (obj *AWSChaos) GetSelectorSpecs() map[string]interface{} {
	return map[string]interface{}{
		".": &obj.Spec.AWSSelector,
	}
}

func (selector *AWSSelector) Id() string {
	// TODO: handle the error here
	// or ignore it is enough ?
	json, _ := json.Marshal(selector)

	return string(json)
}
