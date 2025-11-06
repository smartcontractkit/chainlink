package forwarder

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3/v2/changeset/operations/contracts"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	kf "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// Config is the configuration to set on a Keystone Forwarder contract
type Config struct {
	DonID         uint32           // the DON id as registered in the capabilities registry. Is an id corresponding to a DON that run consensus capability
	F             uint8            // the F value for the DON as registered in the capabilities registry
	ConfigVersion uint32           // the config version for the DON as registered in the capabilities registry
	Signers       []common.Address // the onchain public keys of the nodes in the DON corresponding to DonID
}

type configureFowarderResponse struct {
	ChainSelector uint64
	DonID         uint32
	Forwarder     common.Address

	BatchOperation *mcmstypes.BatchOperation // if using MCMS, the batch operation to propose the change
}

// configureForwarder sets the config for the forwarder contract on the chain for all Dons that accept workflows
// dons that don't accept workflows are not registered with the forwarder
func configureForwarder(lggr logger.Logger, chain cldf_evm.Chain, fwdr *kf.KeystoneForwarder, cfg Config, useMCMS bool) (*configureFowarderResponse, error) {
	if fwdr == nil {
		return nil, errors.New("nil forwarder contract")
	}

	ver := cfg.ConfigVersion // note config count on the don info is the version on the forwarder
	signers := cfg.Signers
	txOpts := chain.DeployerKey
	if useMCMS {
		txOpts = cldf.SimTransactOpts()
	}
	tx, err := fwdr.SetConfig(txOpts, cfg.DonID, ver, cfg.F, signers)
	if err != nil {
		err = cldf.DecodeErr(kf.KeystoneForwarderABI, err)
		return nil, fmt.Errorf("failed to call SetConfig for forwarder %s on chain %d: %w", fwdr.Address().String(), chain.Selector, err)
	}
	var op *mcmstypes.BatchOperation
	if !useMCMS {
		_, err = chain.Confirm(tx)
		if err != nil {
			err = cldf.DecodeErr(kf.KeystoneForwarderABI, err)
			return nil, fmt.Errorf("failed to confirm SetConfig for forwarder %s: %w", fwdr.Address().String(), err)
		}
	} else {
		// create the mcms proposals
		op2, err := proposalutils.BatchOperationForChain(chain.Selector, fwdr.Address().Hex(), tx.Data(), big.NewInt(0), string(contracts.KeystoneForwarder), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create proposal batch operation for chain %d: %w", chain.Selector, err)
		}
		op = &op2
	}
	lggr.Debugw("configured forwarder", "forwarder", fwdr.Address().String(), "donId", cfg.DonID, "version", ver, "f", cfg.F, "signers", signers)

	return &configureFowarderResponse{
		ChainSelector:  chain.Selector,
		DonID:          cfg.DonID,
		Forwarder:      fwdr.Address(),
		BatchOperation: op,
	}, nil
}
