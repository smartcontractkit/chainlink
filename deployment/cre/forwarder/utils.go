package forwarder

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kf "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"

	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
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

	MCMSOperation *mcmstypes.BatchOperation // if using MCMS, the proposed operation for the config change
}

// configureForwarder sets the config for the forwarder contract on the chain for all Dons that accept workflows
// dons that don't accept workflows are not registered with the forwarder
func configureForwarder(
	lggr logger.Logger,
	chain cldf_evm.Chain,
	fwdr *kf.KeystoneForwarder,
	cfg Config,
	useMCMS bool,
	strategy strategies.TransactionStrategy,
) (*configureFowarderResponse, error) {
	if fwdr == nil {
		return nil, errors.New("nil forwarder contract")
	}

	ver := cfg.ConfigVersion // note config count on the don info is the version on the forwarder
	signers := cfg.Signers

	operation, _, err := strategy.Apply(func(txOpts *bind.TransactOpts) (*types.Transaction, error) {
		return fwdr.SetConfig(txOpts, cfg.DonID, ver, cfg.F, signers)
	})
	if err != nil {
		err = cldf.DecodeErr(kf.KeystoneForwarderABI, err)
		return nil, fmt.Errorf("failed to call SetConfig for forwarder %s on chain %d: %w", fwdr.Address().String(), chain.Selector, err)
	}

	if useMCMS {
		lggr.Infow("Created MCMS proposal for forwarder", "address", fwdr.Address().String(), "donId", cfg.DonID, "version", ver, "f", cfg.F, "signers", signers)

		return &configureFowarderResponse{
			ChainSelector: chain.Selector,
			DonID:         cfg.DonID,
			Forwarder:     fwdr.Address(),
			MCMSOperation: operation,
		}, nil
	}

	lggr.Infow("Successfully configured forwarder", "address", fwdr.Address().String(), "donId", cfg.DonID, "version", ver, "f", cfg.F, "signers", signers)

	return &configureFowarderResponse{
		ChainSelector: chain.Selector,
		DonID:         cfg.DonID,
		Forwarder:     fwdr.Address(),
	}, nil
}
