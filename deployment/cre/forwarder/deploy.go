package forwarder

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
)

type DeployOpDeps struct {
	Env *cldf.Environment
}

type DeployOpInput struct {
	ChainSelector uint64
	Qualifier     string
	Labels        []string // optional labels to add to the deployed contract
}

type DeployOpOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // Keeping the address book for backward compatibility, as not everything has been migrated to datastore

	AddressRef datastore.AddressRef // The address ref of the deployed Keystone Forwarder
}

// DeployOp is an operation that deploys the Keystone Forwarder contract.
var DeployOp = operations.NewOperation[DeployOpInput, DeployOpOutput, DeployOpDeps](
	"deploy-keystone-forwarder-op",
	semver.MustParse("1.0.0"),
	"Deploy KeystoneForwarder Contract",
	func(b operations.Bundle, deps DeployOpDeps, input DeployOpInput) (DeployOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployOpOutput{}, fmt.Errorf("deploy-keystone-forwarder-op failed: chain selector %d not found in environment", input.ChainSelector)
		}
		addr, tv, err := deploy(b.GetContext(), chain.DeployerKey, chain)
		if err != nil {
			return DeployOpOutput{}, fmt.Errorf("deploy-keystone-forwarder-op failed: %w", err)
		}
		labels := tv.Labels.List()
		labels = append(labels, input.Labels...)
		r := datastore.AddressRef{
			ChainSelector: input.ChainSelector,
			Address:       addr.String(),
			Type:          datastore.ContractType(tv.Type),
			Version:       &tv.Version,
			Qualifier:     input.Qualifier,
			Labels:        datastore.NewLabelSet(labels...),
		}
		ds := datastore.NewMemoryDataStore()
		if err := ds.AddressRefStore.Add(r); err != nil {
			return DeployOpOutput{}, fmt.Errorf("deploy-keystone-forwarder-op failed: failed to add address ref to datastore: %w", err)
		}
		addressBook := cldf.NewMemoryAddressBook()
		if err := addressBook.Save(input.ChainSelector, addr.String(), *tv); err != nil {
			return DeployOpOutput{}, fmt.Errorf("deploy-keystone-forwarder-op failed: failed to add address ref to address book: %w", err)
		}

		return DeployOpOutput{
			Addresses:   ds.Addresses(),
			AddressBook: addressBook,
			AddressRef:  r,
		}, nil
	},
)

const (
	DeploymentBlockLabel = "deployment-block"
	DeploymentHashLabel  = "deployment-hash"
)

func deploy(ctx context.Context, auth *bind.TransactOpts, chain evm.Chain) (*common.Address, *cldf.TypeAndVersion, error) {
	forwarderAddr, tx, forwarder, err := forwarder.DeployKeystoneForwarder(
		auth,
		chain.Client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}

	_, err = chain.Confirm(tx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to confirm and save KeystoneForwarder: %w", err)
	}
	tvStr, err := forwarder.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get type and version: %w", err)
	}
	tv, err := cldf.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	txHash := tx.Hash()
	txReceipt, err := chain.Client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}
	hashLabel := fmt.Sprintf("%s: %s", DeploymentHashLabel, txHash.Hex())
	blockLabel := fmt.Sprintf("%s: %s", DeploymentBlockLabel, txReceipt.BlockNumber.String())
	tv.Labels.Add(blockLabel)
	tv.Labels.Add(hashLabel)

	return &forwarderAddr, &tv, nil
}

type DeploySequenceDeps struct {
	Env *cldf.Environment // The environment in which the Keystone Forwarders will be deployed
}

type DeploySequenceInput struct {
	Targets   []uint64 // The target chains for the Keystone Forwarders
	Qualifier string   // The qualifier for the forwarder deployment
}

type DeploySequenceOutput struct {
	Addresses   datastore.AddressRefStore
	AddressBook cldf.AddressBook // The address book containing the deployed Keystone Forwarders
	Datastore   datastore.DataStore
}

var DeploySequence = operations.NewSequence[DeploySequenceInput, DeploySequenceOutput, DeploySequenceDeps](
	"deploy-keystone-forwarders-seq",
	semver.MustParse("1.0.0"),
	"Deploy Keystone Forwarders",
	func(b operations.Bundle, deps DeploySequenceDeps, input DeploySequenceInput) (DeploySequenceOutput, error) {
		ab := cldf.NewMemoryAddressBook()
		as := datastore.NewMemoryDataStore()
		contractErrGroup := &errgroup.Group{}
		for _, target := range input.Targets {
			contractErrGroup.Go(func() error {
				r, err := operations.ExecuteOperation(b, DeployOp, DeployOpDeps(deps), DeployOpInput{
					ChainSelector: target,
					Qualifier:     input.Qualifier,
				})
				if err != nil {
					return err
				}
				err = ab.Merge(r.Output.AddressBook)
				if err != nil {
					return fmt.Errorf("failed to save Keystone Forwarder address on address book for target %d: %w", target, err)
				}
				addrs, err := r.Output.Addresses.Fetch()
				if err != nil {
					return fmt.Errorf("failed to fetch Keystone Forwarder addresses for target %d: %w", target, err)
				}
				for _, addr := range addrs {
					if addrRefErr := as.AddressRefStore.Add(addr); addrRefErr != nil {
						return fmt.Errorf("failed to save Keystone Forwarder address on datastore for target %d: %w", target, addrRefErr)
					}
				}

				return nil
			})
		}
		if err := contractErrGroup.Wait(); err != nil {
			return DeploySequenceOutput{AddressBook: ab, Addresses: as.Addresses()}, fmt.Errorf("failed to deploy Keystone contracts: %w", err)
		}
		return DeploySequenceOutput{AddressBook: ab, Addresses: as.Addresses(), Datastore: as.Seal()}, nil
	},
)
