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
	"context"
	"encoding/json"
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var podhttpchaoslog = logf.Log.WithName("rawpodhttp-resource")

// +kubebuilder:object:generate=false

// PodHttpChaosHandler represents the implementation of podhttpchaos
type PodHttpChaosHandler interface {
	Apply(context.Context, *PodHttpChaos) (int32, error)
}

var podHttpChaosHandler PodHttpChaosHandler

// RegisterPodHttpHandler registers handler into webhook
func RegisterPodHttpHandler(newHandler PodHttpChaosHandler) {
	podHttpChaosHandler = newHandler
}

// SetupWebhookWithManager setup PodHttpChaos's webhook with manager
func (in *PodHttpChaos) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mgr.GetWebhookServer().
		Register("/mutate-chaos-mesh-org-v1alpha1-podhttpchaos", &webhook.Admission{Handler: &PodHttpChaosWebhookRunner{}})
	return nil
}

// +kubebuilder:webhook:path=/mutate-chaos-mesh-org-v1alpha1-podhttpchaos,mutating=true,failurePolicy=fail,groups=chaos-mesh.org,resources=podhttpchaos,verbs=create;update,versions=v1alpha1,name=mpodhttpchaos.kb.io

// +kubebuilder:object:generate=false

// PodHttpChaosWebhookRunner runs webhook for podhttpchaos
type PodHttpChaosWebhookRunner struct {
	decoder *admission.Decoder
}

// Handle will run podhttpchaoshandler for this resource
func (r *PodHttpChaosWebhookRunner) Handle(ctx context.Context, req admission.Request) admission.Response {
	chaos := &PodHttpChaos{}
	err := r.decoder.Decode(req, chaos)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if podHttpChaosHandler != nil {
		statusCode, err := podHttpChaosHandler.Apply(ctx, chaos)
		if err != nil {
			return admission.Errored(statusCode, err)
		}
	}

	// mutate the fields in pod
	marshaledPodHttpChaos, err := json.Marshal(chaos)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPodHttpChaos)
}

// InjectDecoder injects decoder into webhook runner
func (r *PodHttpChaosWebhookRunner) InjectDecoder(d *admission.Decoder) error {
	r.decoder = d
	return nil
}
