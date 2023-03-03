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

// PodHttpChaosSpec defines the desired state of PodHttpChaos.
type PodHttpChaosSpec struct {
	// Rules are a list of injection rule for http request.
	// +optional
	Rules []PodHttpChaosRule `json:"rules,omitempty"`
}

// PodHttpChaosStatus defines the actual state of PodHttpChaos.
type PodHttpChaosStatus struct {
	// Pid represents a running tproxy process id.
	// +optional
	Pid int64 `json:"pid,omitempty"`

	// StartTime represents the start time of a tproxy process.
	// +optional
	StartTime int64 `json:"startTime,omitempty"`

	// +optional
	FailedMessage string `json:"failedMessage,omitempty"`

	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// PodHttpChaosRule defines the injection rule for http.
type PodHttpChaosRule struct {
	PodHttpChaosBaseRule `json:",inline"`

	// Source represents the source of current rules
	Source string `json:"source,omitempty"`

	// Port represents the target port to be proxy of.
	Port int32 `json:"port"`
}

// PodHttpChaosBaseRule defines the injection rule without source and port.
type PodHttpChaosBaseRule struct {
	// Target is the object to be selected and injected, <Request|Response>.
	Target PodHttpChaosTarget `json:"target"`

	// Selector contains the rules to select target.
	Selector PodHttpChaosSelector `json:"selector"`

	// Actions contains rules to inject target.
	Actions PodHttpChaosActions `json:"actions"`
}

type PodHttpChaosSelector struct {
	// Port is a rule to select server listening on specific port.
	// +optional
	Port *int32 `json:"port,omitempty"`

	// Path is a rule to select target by uri path in http request.
	// +optional
	Path *string `json:"path,omitempty"`

	// Method is a rule to select target by http method in request.
	// +optional
	Method *string `json:"method,omitempty"`

	// Code is a rule to select target by http status code in response.
	// +optional
	Code *int32 `json:"code,omitempty"`

	// RequestHeaders is a rule to select target by http headers in request.
	// The key-value pairs represent header name and header value pairs.
	// +optional
	RequestHeaders map[string]string `json:"request_headers,omitempty"`

	// ResponseHeaders is a rule to select target by http headers in response.
	// The key-value pairs represent header name and header value pairs.
	// +optional
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
}

// PodHttpChaosActions defines possible actions of HttpChaos.
type PodHttpChaosActions struct {
	// Abort is a rule to abort a http session.
	// +optional
	Abort *bool `json:"abort,omitempty"`

	// Delay represents the delay of the target request/response.
	// A duration string is a possibly unsigned sequence of
	// decimal numbers, each with optional fraction and a unit suffix,
	// such as "300ms", "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	// +optional
	Delay *string `json:"delay,omitempty"`

	// Replace is a rule to replace some contents in target.
	// +optional
	Replace *PodHttpChaosReplaceActions `json:"replace,omitempty"`

	// Patch is a rule to patch some contents in target.
	// +optional
	Patch *PodHttpChaosPatchActions `json:"patch,omitempty"`
}

// PodHttpChaosPatchActions defines possible patch-actions of HttpChaos.
type PodHttpChaosPatchActions struct {
	// Body is a rule to patch message body of target.
	// +optional
	Body *PodHttpChaosPatchBodyAction `json:"body,omitempty"`

	// Queries is a rule to append uri queries of target(Request only).
	// For example: `[["foo", "bar"], ["foo", "unknown"]]`.
	// +optional
	Queries [][]string `json:"queries,omitempty"`

	// Headers is a rule to append http headers of target.
	// For example: `[["Set-Cookie", "<one cookie>"], ["Set-Cookie", "<another cookie>"]]`.
	// +optional
	Headers [][]string `json:"headers,omitempty"`
}

// PodHttpChaosPatchBodyAction defines patch body action of HttpChaos.
type PodHttpChaosPatchBodyAction struct {
	// Type represents the patch type, only support `JSON` as [merge patch json](https://tools.ietf.org/html/rfc7396) currently.
	Type string `json:"type"`

	// Value is the patch contents.
	Value string `json:"value"`
}

// PodHttpChaosReplaceActions defines possible replace-actions of HttpChaos.
type PodHttpChaosReplaceActions struct {
	// Path is rule to to replace uri path in http request.
	// +optional
	Path *string `json:"path,omitempty"`

	// Method is a rule to replace http method in request.
	// +optional
	Method *string `json:"method,omitempty"`

	// Code is a rule to replace http status code in response.
	// +optional
	Code *int32 `json:"code,omitempty"`

	// Body is a rule to replace http message body in target.
	// +optional
	Body []byte `json:"body,omitempty"`

	// Queries is a rule to replace uri queries in http request.
	// For example, with value `{ "foo": "unknown" }`, the `/?foo=bar` will be altered to `/?foo=unknown`,
	// +optional
	Queries map[string]string `json:"queries,omitempty"`

	// Headers is a rule to replace http headers of target.
	// The key-value pairs represent header name and header value pairs.
	// +optional
	Headers map[string]string `json:"headers,omitempty"`
}

// PodHttpChaosTarget represents the type of an HttpChaos Action
type PodHttpChaosTarget string

const (
	// PodHttpRequest represents injecting chaos for http request
	PodHttpRequest PodHttpChaosTarget = "Request"

	// PodHttpResponse represents injecting chaos for http response
	PodHttpResponse PodHttpChaosTarget = "Response"
)

// +kubebuilder:object:root=true

// +chaos-mesh:base
// +chaos-mesh:webhook:enableUpdate
// +kubebuilder:subresource:status
// PodHttpChaos is the Schema for the podhttpchaos API
type PodHttpChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodHttpChaosSpec   `json:"spec,omitempty"`
	Status PodHttpChaosStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodHttpChaosList contains a list of PodHttpChaos
type PodHttpChaosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodHttpChaos `json:"items"`
}
