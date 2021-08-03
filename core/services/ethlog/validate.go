package ethlog

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
)

type EthLogToml struct {
	ContractAddress ethkey.EIP55Address `toml:"contractAddress"`
}

func ValidatedEthLogSpec(tomlString string) (job.Job, error) {
	var jb = job.Job{}
	tree, err := toml.Load(tomlString)
	if err != nil {
		return jb, err
	}
	err = tree.Unmarshal(&jb)
	if err != nil {
		return jb, err
	}
	var spec EthLogToml
	err = tree.Unmarshal(&spec)
	if err != nil {
		return jb, err
	}
	jb.EthLogSpec = &job.EthLogSpec{ContractAddress: spec.ContractAddress}

	if jb.Type != job.EthLog {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	return jb, nil
}
