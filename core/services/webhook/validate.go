package webhook

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type TOMLWebhookSpecExternalInitiator struct {
	Name string      `toml:"name"`
	Spec models.JSON `toml:"spec"`
}

type TOMLWebhookSpec struct {
	ExternalInitiators []TOMLWebhookSpecExternalInitiator `toml:"externalInitiators"`
}

func ValidatedWebhookSpec(tomlString string, externalInitiatorManager ExternalInitiatorManager) (jb job.Job, err error) {
	var tree *toml.Tree
	tree, err = toml.Load(tomlString)
	if err != nil {
		return
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return
	}
	if jb.Type != job.Webhook {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}

	var tomlSpec TOMLWebhookSpec
	err = tree.Unmarshal(&tomlSpec)
	if err != nil {
		return jb, err
	}

	var externalInitiatorWebhookSpecs []job.ExternalInitiatorWebhookSpec
	for _, eiSpec := range tomlSpec.ExternalInitiators {
		ei, findErr := externalInitiatorManager.FindExternalInitiatorByName(eiSpec.Name)
		if findErr != nil {
			err = multierr.Combine(err, errors.Wrapf(findErr, "unable to find external initiator named %s", eiSpec.Name))
			continue
		}
		eiWS := job.ExternalInitiatorWebhookSpec{
			ExternalInitiatorID: ei.ID,
			WebhookSpecID:       0, // It will be populated later, on save
			Spec:                eiSpec.Spec,
		}
		externalInitiatorWebhookSpecs = append(externalInitiatorWebhookSpecs, eiWS)
	}

	if err != nil {
		return jb, err
	}

	jb.WebhookSpec = &job.WebhookSpec{
		ExternalInitiatorWebhookSpecs: externalInitiatorWebhookSpecs,
	}

	return jb, nil
}
