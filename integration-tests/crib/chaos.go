package crib

import (
	"time"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/smartcontractkit/havoc/k8schaos"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func rebootCLNamespace(delay time.Duration, namespace string) (*k8schaos.Chaos, error) {
	k8sClient, err := k8schaos.NewChaosMeshClient()
	if err != nil {
		return nil, err
	}
	return k8schaos.NewChaos(k8schaos.ChaosOpts{
		Description: "Reboot CRIB",
		DelayCreate: delay,
		Object: &v1alpha1.PodChaos{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodChaos",
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "reboot-crib",
				Namespace: namespace,
			},
			Spec: v1alpha1.PodChaosSpec{
				ContainerSelector: v1alpha1.ContainerSelector{
					PodSelector: v1alpha1.PodSelector{
						Mode: v1alpha1.AllMode,
						Selector: v1alpha1.PodSelectorSpec{
							GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
								Namespaces: []string{namespace},
							},
						},
					},
				},
				Action: v1alpha1.PodKillAction,
			},
		},
		Client: k8sClient,
		Logger: &k8schaos.Logger,
	})
}
