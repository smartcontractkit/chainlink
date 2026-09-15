package stellar

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	stellarfwd "github.com/smartcontractkit/chainlink/deployment/cre/forwarder/stellar"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
)

func forwarderAddress(ds datastore.DataStore, chainSelector uint64) (string, bool) {
	key := datastore.NewAddressRefKey(
		chainSelector,
		stellarfwd.ForwarderContract,
		semver.MustParse(stellarfwd.DefaultForwarderVersion),
		stellarfwd.DefaultForwarderQualifier,
	)
	ref, err := ds.Addresses().Get(key)
	if err != nil {
		return "", false
	}
	return ref.Address, true
}

func mustForwarderAddress(ds datastore.DataStore, chainSelector uint64) string {
	addr, ok := forwarderAddress(ds, chainSelector)
	if !ok {
		panic(fmt.Sprintf("missing Stellar forwarder address for chain selector %d", chainSelector))
	}
	return addr
}

// forwarderVersion returns the Stellar forwarder contract version pinned by the
// environment, falling back to DefaultForwarderVersion when none is set.
func forwarderVersion(creEnv *cre.Environment) *semver.Version {
	if v := creEnv.ContractVersions[stellarfwd.ForwarderContract.String()]; v != nil {
		return v
	}
	return semver.MustParse(stellarfwd.DefaultForwarderVersion)
}

// deployStellarForwarders deploys the Soroban CRE forwarder via the CLDF
// changeset and records its C-address in the environment datastore. Mirrors
// features/evm/v2 deployEVMForwarders.
func deployStellarForwarders(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	chain *stellchain.Blockchain,
) error {
	if addr, ok := forwarderAddress(creEnv.CldfEnvironment.DataStore, chain.ChainSelector()); ok {
		testLogger.Info().
			Uint64("chainSelector", chain.ChainSelector()).
			Str("forwarder", addr).
			Msg("Stellar CRE forwarder already present in datastore")
		return nil
	}

	// Fund the CLDF deployer keypair before the changeset uploads WASM.
	stellarChain, ok := creEnv.CldfEnvironment.BlockChains.StellarChains()[chain.ChainSelector()]
	if !ok {
		return fmt.Errorf("stellar chain selector %d not registered in CLDF BlockChains", chain.ChainSelector())
	}
	if stellarChain.Signer != nil {
		if fundErr := chain.Fund(ctx, stellarChain.Signer.Address(), 0); fundErr != nil {
			return fmt.Errorf("failed to fund stellar deployer %s via friendbot: %w", stellarChain.Signer.Address(), fundErr)
		}
	}

	version := forwarderVersion(creEnv)

	out, err := (stellarfwd.DeployForwarder{}).Apply(*creEnv.CldfEnvironment, &stellarfwd.DeployForwarderRequest{
		ChainSel:  chain.ChainSelector(),
		Qualifier: stellarfwd.DefaultForwarderQualifier,
		Version:   version.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to deploy stellar forwarder via CLDF: %w", err)
	}

	memoryDatastore, err := contracts.NewDataStoreFromExisting(creEnv.CldfEnvironment.DataStore)
	if err != nil {
		return fmt.Errorf("failed to create memory datastore: %w", err)
	}
	if out.DataStore != nil {
		if mergeErr := memoryDatastore.Merge(out.DataStore.Seal()); mergeErr != nil {
			return fmt.Errorf("failed to merge stellar forwarder datastore: %w", mergeErr)
		}
	}
	creEnv.CldfEnvironment.DataStore = memoryDatastore.Seal()

	addr := mustForwarderAddress(creEnv.CldfEnvironment.DataStore, chain.ChainSelector())
	testLogger.Info().
		Uint64("chainSelector", chain.ChainSelector()).
		Str("forwarder", addr).
		Msg("Deployed Stellar CRE forwarder to CLDF datastore")

	return nil
}

// configureStellarForwarders configures the deployed forwarder with the
// consensus DON's ed25519 signer set. Mirrors features/evm/v2 configureEVMForwarders.
func configureStellarForwarders(
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	consensusDon *cre.Don,
	chain *stellchain.Blockchain,
) error {
	if _, ok := forwarderAddress(creEnv.CldfEnvironment.DataStore, chain.ChainSelector()); !ok {
		return fmt.Errorf("missing Stellar forwarder in datastore for chain selector %d", chain.ChainSelector())
	}

	version := forwarderVersion(creEnv)

	_, err := (stellarfwd.ConfigureForwarders{}).Apply(*creEnv.CldfEnvironment, &stellarfwd.ConfigureForwarderRequest{
		DON: forwarder.DonConfiguration{
			Name:    consensusDon.Name,
			ID:      libc.MustSafeUint32FromUint64(consensusDon.ID),
			F:       consensusDon.F,
			Version: 1,
			NodeIDs: consensusDon.KeystoneDONConfig().NodeIDs,
		},
		Chains:    map[uint64]struct{}{chain.ChainSelector(): {}},
		Qualifier: stellarfwd.DefaultForwarderQualifier,
		Version:   version.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to configure stellar forwarder for DON %q: %w", consensusDon.Name, err)
	}

	testLogger.Info().
		Uint64("chainSelector", chain.ChainSelector()).
		Str("don", consensusDon.Name).
		Msg("Configured Stellar CRE forwarder")
	return nil
}

// authorizeStellarForwarderTransmitters registers each capability worker's Stellar
// signing account as an authorized transmitter on the forwarder. report() calls
// require_valid_forwarder(transmitter); the worker submits report() from its own
// relayer signing account (n.Keys.Stellar[chainID].Account), so each worker account
// must be authorized or report() reverts with UnauthorizedForwarder. Uses the
// capability DON (the one hosting the Stellar worker / submitting reports).
func authorizeStellarForwarderTransmitters(
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	don *cre.Don,
	chain *stellchain.Blockchain,
) error {
	workers, err := don.Workers()
	if err != nil {
		return fmt.Errorf("failed to list worker nodes for DON %q: %w", don.Name, err)
	}
	chainID := chain.StellarChainID()
	transmitters := make([]string, 0, len(workers))
	for _, worker := range workers {
		if worker.Keys == nil || worker.Keys.Stellar == nil {
			return fmt.Errorf("missing Stellar keys for worker %q in DON %q", worker.Name, don.Name)
		}
		key, ok := worker.Keys.Stellar[chainID]
		if !ok || key == nil || key.Account == "" {
			return fmt.Errorf("missing Stellar account for worker %q (chain %s) in DON %q", worker.Name, chainID, don.Name)
		}
		transmitters = append(transmitters, key.Account)
	}

	version := forwarderVersion(creEnv)

	if _, err := (stellarfwd.AddForwarders{}).Apply(*creEnv.CldfEnvironment, &stellarfwd.AddForwardersRequest{
		Transmitters: transmitters,
		Chains:       map[uint64]struct{}{chain.ChainSelector(): {}},
		Qualifier:    stellarfwd.DefaultForwarderQualifier,
		Version:      version.String(),
	}); err != nil {
		return fmt.Errorf("failed to authorize Stellar forwarder transmitters for DON %q: %w", don.Name, err)
	}

	testLogger.Info().
		Uint64("chainSelector", chain.ChainSelector()).
		Str("don", don.Name).
		Int("transmitters", len(transmitters)).
		Msg("Authorized Stellar forwarder transmitters")
	return nil
}
