package directrequest

import (
	"crypto/sha256"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

func ValidatedDirectRequestSpec(tomlString string) (job.SpecDB, error) {
	var specDB = job.SpecDB{
		Pipeline: *pipeline.NewTaskDAG(),
	}
	var spec job.DirectRequestSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return specDB, err
	}
	err = tree.Unmarshal(&specDB)
	if err != nil {
		return specDB, err
	}
	err = tree.Unmarshal(&spec)
	if err != nil {
		return specDB, err
	}
	spec.OnChainJobSpecID = sha256.Sum256([]byte(tomlString))
	specDB.DirectRequestSpec = &spec

	if specDB.Type != job.DirectRequest {
		return specDB, errors.Errorf("unsupported type %s", specDB.Type)
	}
	if specDB.SchemaVersion != uint32(1) {
		return specDB, errors.Errorf("the only supported schema version is currently 1, got %v", specDB.SchemaVersion)
	}
	return specDB, nil
}
