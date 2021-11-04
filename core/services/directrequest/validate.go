package directrequest

import (
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type DirectRequestToml struct {
	ContractAddress    ethkey.EIP55Address      `toml:"contractAddress"`
	Requesters         models.AddressCollection `toml:"requesters"`
	MinContractPayment *assets.Link             `toml:"minContractPaymentLinkJuels"`
	EVMChainID         *utils.Big               `toml:"evmChainID"`
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
		ContractAddress:    spec.ContractAddress,
		Requesters:         spec.Requesters,
		MinContractPayment: spec.MinContractPayment,
		EVMChainID:         spec.EVMChainID,
	}

	if jb.Type != job.DirectRequest {
		return jb, errors.Errorf("unsupported type %s", jb.Type)
	}
	return jb, nil
}
