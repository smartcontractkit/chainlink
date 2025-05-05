package verification

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
)

var LatestConfigDetailsChangeset = cldf.CreateChangeSet(LatestConfigDetailsLogic, LatestConfigDetailsPrecondition)

type LatestConfigDetailsConfig struct {
	ConfigsByChain map[uint64][]LatestConfigDetails
}

type LatestConfigDetails struct {
	VerifierAddress common.Address
	ConfigDigest    [32]byte
}

func (cfg LatestConfigDetailsConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func LatestConfigDetailsPrecondition(_ deployment.Environment, cc LatestConfigDetailsConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid LatestConfigDetails config: %w", err)
	}
	return nil
}

func LatestConfigDetailsLogic(e deployment.Environment, cfg LatestConfigDetailsConfig) (deployment.ChangesetOutput, error) {
	for chainSelector, configs := range cfg.ConfigsByChain {
		for _, config := range configs {
			confState, err := loadVerifierState(e, chainSelector, config.VerifierAddress.Hex())
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			block, err := confState.LatestConfigDetails(&bind.CallOpts{Context: e.GetContext()}, config.ConfigDigest)
			fmt.Println("AddressBookV2 -> The latest config block number is: ", block)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
		}
	}

	return deployment.ChangesetOutput{}, nil
}
