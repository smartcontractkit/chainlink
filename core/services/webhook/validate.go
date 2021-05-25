package webhook

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

type WebhookToml struct {
	OnChainJobSpecID uuid.UUID `toml:"jobID"`
}

var ErrMissingJobID = errors.New("missing job ID")

func ValidateWebhookSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	var spec WebhookToml
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}
	if spec.OnChainJobSpecID == (uuid.UUID{}) {
		return jb, ErrMissingJobID
	}
	jb.WebhookSpec = &job.WebhookSpec{}
	copy(jb.WebhookSpec.OnChainJobSpecID[:], spec.OnChainJobSpecID.Bytes())

	if jb.Type != job.Webhook {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}
	return jb, nil
}
