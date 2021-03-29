package directrequest

import (
	"github.com/gofrs/uuid"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func ValidatedDirectRequestSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	var spec job.DirectRequestSpec
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
	onChainJobSpecID, err := uuid.NewV4()
	if err != nil {
		return jb, err
	}
	copy(spec.OnChainJobSpecID[:], onChainJobSpecID[:])
	jb.DirectRequestSpec = &spec

	if jb.Type != job.DirectRequest {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	if jb.SchemaVersion != uint32(1) {
		return jb, errors.Errorf("the only supported schema version is currently 1, got %v", jb.SchemaVersion)
	}
	return jb, nil
}
