package directrequest

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

type DirectRequestToml struct {
	ContractAddress          ethkey.EIP55Address `toml:"contractAddress"`
	MinIncomingConfirmations clnull.Uint32       `toml:"minIncomingConfirmations"`
}

func ValidatedDirectRequestSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{}
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}
	var spec DirectRequestToml
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}
	jb.DirectRequestSpec = &job.DirectRequestSpec{
		ContractAddress:          spec.ContractAddress,
		MinIncomingConfirmations: spec.MinIncomingConfirmations,
	}

	if jb.Type != job.DirectRequest {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	return jb, nil
}
