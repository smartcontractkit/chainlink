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

var DeactivateConfigChangeset = cldf.CreateChangeSet(deactivateConfigLogic, deactivateConfigPrecondition)

type DeactivateConfigConfig struct {
	ConfigsByChain map[uint64][]DeactivateConfig
	MCMSConfig     *types.MCMSConfig
}

type DeactivateConfig struct {
	VerifierAddress common.Address
	ConfigDigest    [32]byte
}

func (a DeactivateConfig) GetContractAddress() common.Address {
	return a.VerifierAddress
}

func (cfg DeactivateConfigConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func deactivateConfigPrecondition(_ deployment.Environment, cc DeactivateConfigConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid ActivateConfig config: %w", err)
	}
	return nil
}

func deactivateConfigLogic(e deployment.Environment, cfg DeactivateConfigConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.VerifierProxy.String(),
		cfg.ConfigsByChain,
		loadVerifierState,
		doDeactivateConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building DeactivateConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "ActivateConfig proposal")
}

func doDeactivateConfig(v *verifier_v0_5_0.Verifier, ac DeactivateConfig) (*goEthTypes.Transaction, error) {
	return v.DeactivateConfig(
		deployment.SimTransactOpts(),
		ac.ConfigDigest,
	)
}
