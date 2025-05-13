package verification

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	goEthTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/txutil"
)

var UpdateConfigChangeset = cldf.CreateChangeSet(updateConfigLogic, updateConfigPrecondition)

type UpdateConfigConfig struct {
	ConfigsByChain map[uint64][]UpdateConfig
	MCMSConfig     *types.MCMSConfig
}

type UpdateConfig struct {
	VerifierAddress common.Address
	ConfigDigest    [32]byte
	PrevSigners     []common.Address
	NewSigners      []common.Address
	F               uint8
}

func (a UpdateConfig) GetContractAddress() common.Address {
	return a.VerifierAddress
}

func (cfg UpdateConfigConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func updateConfigPrecondition(_ cldf.Environment, cc UpdateConfigConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid UpdateConfig config: %w", err)
	}
	return nil
}

func updateConfigLogic(e cldf.Environment, cfg UpdateConfigConfig) (cldf.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.VerifierProxy.String(),
		cfg.ConfigsByChain,
		loadVerifierState,
		doUpdateConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed building UpdateConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "UpdateConfig proposal")
}

func doUpdateConfig(v *verifier_v0_5_0.Verifier, ac UpdateConfig) (*goEthTypes.Transaction, error) {
	return v.UpdateConfig(
		cldf.SimTransactOpts(),
		ac.ConfigDigest,
		ac.PrevSigners,
		ac.NewSigners,
		ac.F,
	)
}
