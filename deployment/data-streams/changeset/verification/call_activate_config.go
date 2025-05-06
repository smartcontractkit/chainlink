package verification

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

var ActivateConfigChangeset = cldf.CreateChangeSet(activateConfigLogic, activateConfigPrecondition)

type ActivateConfigConfig struct {
	ConfigsByChain map[uint64][]ActivateConfig
	MCMSConfig     *types.MCMSConfig
}

type ActivateConfig struct {
	VerifierAddress common.Address
	ConfigDigest    [32]byte
}

func (a ActivateConfig) GetContractAddress() common.Address {
	return a.VerifierAddress
}

func (cfg ActivateConfigConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func activateConfigPrecondition(_ deployment.Environment, cc ActivateConfigConfig) error {
	return cc.Validate()
}

func activateConfigLogic(e deployment.Environment, cfg ActivateConfigConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.VerifierProxy.String(),
		cfg.ConfigsByChain,
		loadVerifierState,
		doActivateConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building ActivateConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "ActivateConfig proposal")
}

func doActivateConfig(
	v *verifier_v0_5_0.Verifier,
	ac ActivateConfig,
) (*goEthTypes.Transaction, error) {
	return v.ActivateConfig(
		deployment.SimTransactOpts(),
		ac.ConfigDigest,
	)
}
