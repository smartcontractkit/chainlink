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

var SetConfigChangeset = cldf.CreateChangeSet(setConfigLogic, setConfigPrecondition)

type SetConfigConfig struct {
	ConfigsByChain map[uint64][]SetConfig
	MCMSConfig     *types.MCMSConfig
}

type SetConfig struct {
	VerifierAddress            common.Address
	ConfigDigest               [32]byte
	Signers                    []common.Address
	F                          uint8
	RecipientAddressesAndProps []verifier_v0_5_0.CommonAddressAndWeight
}

func (a SetConfig) GetContractAddress() common.Address {
	return a.VerifierAddress
}

func (cfg SetConfigConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func setConfigPrecondition(_ deployment.Environment, cc SetConfigConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid SetConfig config: %w", err)
	}
	return nil
}

func setConfigLogic(e deployment.Environment, cfg SetConfigConfig) (deployment.ChangesetOutput, error) {
	txs, err := txutil.GetTxs(
		e,
		types.Verifier.String(),
		cfg.ConfigsByChain,
		loadVerifierState,
		doSetConfig,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed building SetConfig txs: %w", err)
	}

	return mcmsutil.ExecuteOrPropose(e, txs, cfg.MCMSConfig, "SetConfig proposal")
}

func doSetConfig(v *verifier_v0_5_0.Verifier, ac SetConfig) (*goEthTypes.Transaction, error) {
	return v.SetConfig(
		deployment.SimTransactOpts(),
		ac.ConfigDigest,
		ac.Signers,
		ac.F,
		ac.RecipientAddressesAndProps,
	)
}
