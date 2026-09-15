package stellar

import (
	"context"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stellar/go-stellar-sdk/keypair"
	"github.com/stellar/go-stellar-sdk/xdr"

	mcmsstellar "github.com/smartcontractkit/mcms/sdk/stellar"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldfstellar "github.com/smartcontractkit/chainlink-deployments-framework/chain/stellar"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-stellar/bindings"
	"github.com/smartcontractkit/chainlink-stellar/bindings/scval"
)

// forwarderSetConfigBatchOp encodes the forwarder set_config call as an MCMS
// batch operation. The argument encoding mirrors ForwarderClient.SetConfig in
// chainlink-stellar/bindings/contracts/cre.
func forwarderSetConfigBatchOp(
	chainSelector uint64,
	forwarder string,
	donID, configVersion, f uint32,
	signers [][32]byte,
) (mcmstypes.BatchOperation, error) {
	args := []xdr.ScVal{
		scval.Uint32ToScVal(donID),
		scval.Uint32ToScVal(configVersion),
		scval.Uint32ToScVal(f),
		scval.Bytes32SliceToScVal(signers),
	}

	return mcmsstellar.NewBatchOperation(
		mcmstypes.ChainSelector(chainSelector),
		forwarder,
		"set_config",
		args,
		string(ForwarderContract),
		nil,
	)
}

// forwarderClearConfigBatchOp encodes the forwarder clear_config call as an
// MCMS batch operation. The argument encoding mirrors ForwarderClient.ClearConfig.
func forwarderClearConfigBatchOp(
	chainSelector uint64,
	forwarder string,
	donID, configVersion uint32,
) (mcmstypes.BatchOperation, error) {
	args := []xdr.ScVal{
		scval.Uint32ToScVal(donID),
		scval.Uint32ToScVal(configVersion),
	}

	return mcmsstellar.NewBatchOperation(
		mcmstypes.ChainSelector(chainSelector),
		forwarder,
		"clear_config",
		args,
		string(ForwarderContract),
		nil,
	)
}

// addForwardersBatchOp encodes one add_forwarder call per transmitter as a
// single MCMS batch operation, so the allow-list change applies atomically.
// The argument encoding mirrors ForwarderClient.AddForwarder.
func addForwardersBatchOp(
	chainSelector uint64,
	forwarder string,
	transmitters []string,
) (mcmstypes.BatchOperation, error) {
	transactions := make([]mcmstypes.Transaction, 0, len(transmitters))
	for _, transmitter := range transmitters {
		tx, err := mcmsstellar.NewTransaction(
			forwarder,
			"add_forwarder",
			[]xdr.ScVal{scval.AddressToScVal(transmitter)},
			string(ForwarderContract),
			nil,
		)
		if err != nil {
			return mcmstypes.BatchOperation{}, fmt.Errorf("build add_forwarder transaction for %s: %w", transmitter, err)
		}
		transactions = append(transactions, tx)
	}

	return mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(chainSelector),
		Transactions:  transactions,
	}, nil
}

// readOnlyInvoker builds an invoker for simulate-only reads. It signs nothing:
// the owner/pending-owner lookups go through SimulateContract, which only
// needs a source-account address, so an ephemeral keypair suffices. This keeps
// governed proposal building usable without the deployer signing key.
func readOnlyInvoker(chain cldfstellar.Chain) (bindings.Invoker, error) {
	kp, err := keypair.Random()
	if err != nil {
		return nil, fmt.Errorf("generate ephemeral stellar keypair: %w", err)
	}

	return mcmsstellar.NewInvokerWithNetworkPassphrase(chain.Client, bindings.NewStellarKeypairSigner(kp), chain.NetworkPassphrase)
}

// requireTimelockOwnership fails unless the forwarder is owned by the MCMS
// timelock the proposal input resolves to. It guards the governed path against
// building proposals the timelock cannot execute. The Stellar MCMS reader must
// be registered by the consumer (e.g. via a blank import of
// cld-changesets/mcms/stellar/readers in the CLD domain).
func requireTimelockOwnership(
	ctx context.Context,
	env cldf.Environment,
	chain cldfstellar.Chain,
	forwarder string,
	mcmsInput cldf.MCMSTimelockProposalInput,
) error {
	chainSelector := chain.Selector

	reader, ok := cldf.GetMCMSReaderRegistry().Get(chainselectors.FamilyStellar)
	if !ok {
		return fmt.Errorf("no MCMS reader registered for chain family %q", chainselectors.FamilyStellar)
	}

	timelockRef, err := reader.GetTimelockRef(env, chainSelector, mcmsInput)
	if err != nil {
		return fmt.Errorf("resolve Stellar timelock for chain %d: %w", chainSelector, err)
	}

	invoker, err := readOnlyInvoker(chain)
	if err != nil {
		return err
	}

	owner, err := mcmsstellar.NewInspectorFromInvoker(invoker).GetOwner(ctx, forwarder)
	if err != nil {
		return fmt.Errorf("read owner of stellar forwarder %s: %w", forwarder, err)
	}
	if owner == nil {
		return fmt.Errorf("stellar forwarder %s has no owner", forwarder)
	}
	if *owner != timelockRef.Address {
		return fmt.Errorf(
			"stellar forwarder %s is owned by %s, not the MCMS timelock %s; apply directly (without MCMS) while the deployer owns it, or transfer ownership to the timelock first",
			forwarder,
			*owner,
			timelockRef.Address,
		)
	}

	return nil
}
