package client

import (
	"context"
	"fmt"
	"time"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-env/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Chaos is controller that manages Chaosmesh CRD instances to run experiments
type Chaos struct {
	Client         *K8sClient
	ResourceByName map[string]string
	Namespace      string
}

type ChaosState struct {
	ChaosDetails v1alpha1.ChaosStatus `json:"status"`
}

// NewChaos creates controller to run and stop chaos experiments
func NewChaos(client *K8sClient, namespace string) *Chaos {
	return &Chaos{
		Client:         client,
		ResourceByName: make(map[string]string),
		Namespace:      namespace,
	}
}

// Run runs experiment and saves its ID
func (c *Chaos) Run(app cdk8s.App, id string, resource string) (string, error) {
	log.Info().Msg("Applying chaos experiment")
	config.JSIIGlobalMu.Lock()
	manifest := *app.SynthYaml()
	config.JSIIGlobalMu.Unlock()
	log.Trace().Str("Raw", manifest).Msg("Manifest")
	c.ResourceByName[id] = resource
	if err := c.Client.Apply(manifest); err != nil {
		return id, err
	}
	err := c.waitForChaosStatus(id, v1alpha1.ConditionAllInjected)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (c *Chaos) waitForChaosStatus(id string, condition v1alpha1.ChaosConditionType) error {
	var result ChaosState
	log.Info().Msgf("waiting for chaos experiment state %s", condition)
	return wait.PollImmediate(2*time.Second, 1*time.Minute, func() (bool, error) {
		data, err := c.Client.ClientSet.
			RESTClient().
			Get().
			RequestURI(fmt.Sprintf("/apis/chaos-mesh.org/v1alpha1/namespaces/%s/%s/%s", c.Namespace, c.ResourceByName[id], id)).
			Do(context.Background()).
			Raw()
		if err == nil {
			err = json.Unmarshal(data, &result)
			if err != nil {
				return false, err
			}
			for _, c := range result.ChaosDetails.Conditions {
				if c.Type == condition && c.Status == v1.ConditionTrue {
					return true, err
				}
			}
		}
		return false, nil
	})
}

func (c *Chaos) WaitForAllRecovered(id string) error {
	return c.waitForChaosStatus(id, v1alpha1.ConditionAllRecovered)
}

// Stop removes a chaos experiment
func (c *Chaos) Stop(id string) error {
	defer delete(c.ResourceByName, id)
	return c.Client.DeleteResource(c.Namespace, c.ResourceByName[id], id)
}
