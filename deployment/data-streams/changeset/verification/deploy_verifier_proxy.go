package verification

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/data-streams/view/v0_5"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"

	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/metadata"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/mcmsutil"
)

// DeployVerifierProxyChangeset deploys VerifierProxy to the chains specified in the config.
var DeployVerifierProxyChangeset = cldf.CreateChangeSet(verifierProxyDeployLogic, verifierProxyDeployPrecondition)

type DeployVerifierProxyConfig struct {
	// ChainsToDeploy is a list of chain selectors to deploy the contract to.
	ChainsToDeploy map[uint64]DeployVerifierProxy
	Ownership      types.OwnershipSettings
	Version        semver.Version
}

func (cfg DeployVerifierProxyConfig) GetOwnershipConfig() types.OwnershipSettings {
	return cfg.Ownership
}

type DeployVerifierProxy struct {
	AccessControllerAddress common.Address
}

func (cfg DeployVerifierProxyConfig) Validate() error {
	switch cfg.Version {
	case deployment.Version0_5_0:
		// no-op
	default:
		return fmt.Errorf("unsupported contract version %s", cfg.Version)
	}
	if len(cfg.ChainsToDeploy) == 0 {
		return errors.New("ChainsToDeploy is empty")
	}
	for chain := range cfg.ChainsToDeploy {
		if err := cldf.IsValidChainSelector(chain); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", chain, err)
		}
	}
	return nil
}

func verifierProxyDeployLogic(e cldf.Environment, cc DeployVerifierProxyConfig) (cldf.ChangesetOutput, error) {
	dataStore := ds.NewMemoryDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata]()
	err := deploy(e, dataStore, cc)
	if err != nil {
		e.Logger.Errorw("Failed to deploy VerifierProxy", "err", err)
		return cldf.ChangesetOutput{}, cldf.MaybeDataErr(err)
	}

	records, err := dataStore.Addresses().Fetch()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	proposals, err := mcmsutil.GetTransferOwnershipProposals(e, cc, records)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
	}

	sealedDS, err := ds.ToDefault(dataStore.Seal())
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to convert data store to default format: %w", err)
	}

	return cldf.ChangesetOutput{
		DataStore:             sealedDS,
		MCMSTimelockProposals: proposals,
	}, nil
}

type HasOwnershipConfig interface {
	GetOwnershipConfig() types.OwnershipSettings
}

func verifierProxyDeployPrecondition(_ cldf.Environment, cc DeployVerifierProxyConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}
	return nil
}

func deploy(e cldf.Environment,
	dataStore ds.MutableDataStore[metadata.SerializedContractMetadata, ds.DefaultMetadata],
	cfg DeployVerifierProxyConfig) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid DeployVerifierProxyConfig: %w", err)
	}

	for chainSel, chainCfg := range cfg.ChainsToDeploy {
		chain, ok := e.Chains[chainSel]
		if !ok {
			return fmt.Errorf("chain not found for chain selector %d", chainSel)
		}

		res, err := changeset.DeployContract(e, dataStore, chain, verifyProxyDeployFn(chainCfg), nil)
		if err != nil {
			return fmt.Errorf("failed to deploy verifier proxy: %w", err)
		}
		contractMetadata := metadata.GenericContractMetadata[v0_5.VerifierProxyView]{
			Metadata: metadata.CommonContractMetadata{
				DeployBlock: res.Block,
			},
		}

		serialized, err := metadata.NewSerializedContractMetadata(contractMetadata)
		if err != nil {
			return fmt.Errorf("failed to serialize contract metadata: %w", err)
		}

		if err = dataStore.ContractMetadata().Upsert(
			ds.ContractMetadata[metadata.SerializedContractMetadata]{
				ChainSelector: chain.Selector,
				Address:       res.Address.String(),
				Metadata:      *serialized,
			},
		); err != nil {
			return fmt.Errorf("failed to upser contract metadata: %w", err)
		}
	}

	return nil
}

// verifyProxyDeployFn returns a function that deploys a VerifyProxy contract.
func verifyProxyDeployFn(cfg DeployVerifierProxy) changeset.ContractDeployFn[*verifier_proxy_v0_5_0.VerifierProxy] {
	return func(chain cldf.Chain) *changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy] {
		addr, tx, contract, err := verifier_proxy_v0_5_0.DeployVerifierProxy(
			chain.DeployerKey,
			chain.Client,
			cfg.AccessControllerAddress,
		)
		if err != nil {
			return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
				Err: err,
			}
		}
		bn, err := chain.Confirm(tx)
		if err != nil {
			return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
				Err: err,
			}
		}
		return &changeset.ContractDeployment[*verifier_proxy_v0_5_0.VerifierProxy]{
			Address:  addr,
			Block:    bn,
			Contract: contract,
			Tx:       tx,
			Tv:       cldf.NewTypeAndVersion(types.VerifierProxy, deployment.Version0_5_0),
			Err:      nil,
		}
	}
}
