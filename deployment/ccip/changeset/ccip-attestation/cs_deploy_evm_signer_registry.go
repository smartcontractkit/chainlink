package ccip_attestation

import (
	"fmt"
	"math/big"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
)

var (
	EVMSignerRegistryDeploymentChangeset = cldf.CreateChangeSet(signerRegistryDeploymentLogic, signerRegistryDeploymentPrecondition)
)

const (
	MaxSigners          = 20
	BaseMainnetSelector = 15971525489660198786
	BaseSepoliaSelector = 10344971235874465080
)

type SignerRegistryChangesetConfig struct {
	// MaxSigners is the maximum number of signers that can be registered.
	MaxSigners uint32
	// Signers is the initial set of signers to register.
	Signers []signer_registry.ISignerRegistrySigner
}

func signerRegistryDeploymentPrecondition(env cldf.Environment, config SignerRegistryChangesetConfig) error {
	if config.MaxSigners != MaxSigners {
		return fmt.Errorf("max signers must be %d", MaxSigners)
	}

	signers := config.Signers
	if len(signers) > int(MaxSigners) {
		return fmt.Errorf("too many signers: %d > %d", len(signers), MaxSigners)
	}
	// ensure no duplicates among all EVM addresses and non-zero new EVM addresses
	seen := make(map[string]struct{})
	for _, signer := range signers {
		if signer.EvmAddress == utils.ZeroAddress {
			return fmt.Errorf("signer %s has zero evm address", signer.EvmAddress)
		}
		if signer.EvmAddress == signer.NewEVMAddress {
			return fmt.Errorf("signer %s has the same evm address and new evm address", signer.EvmAddress)
		}
		// Check duplicates for EvmAddress
		evmAddrHex := signer.EvmAddress.Hex()
		if _, ok := seen[evmAddrHex]; ok {
			return fmt.Errorf("duplicate signer evm address: %s", signer.EvmAddress)
		}
		seen[evmAddrHex] = struct{}{}
		// Check duplicates for non-zero NewEVMAddress
		if signer.NewEVMAddress != utils.ZeroAddress {
			newEVMAddrHex := signer.NewEVMAddress.Hex()
			if _, ok := seen[newEVMAddrHex]; ok {
				return fmt.Errorf("duplicate signer new EVM address: %s", signer.NewEVMAddress)
			}
			seen[newEVMAddrHex] = struct{}{}
		}
	}
	return nil
}

func signerRegistryDeploymentLogic(e cldf.Environment, config SignerRegistryChangesetConfig) (cldf.ChangesetOutput, error) {
	addressBook := cldf.NewMemoryAddressBook()

	for _, chain := range e.BlockChains.EVMChains() {
		// Only deploy on Base Mainnet and Base Sepolia
		if chain.ChainSelector() != BaseMainnetSelector && chain.ChainSelector() != BaseSepoliaSelector {
			e.Logger.Infof("Skipping deployment on chain %s (selector: %d) - only deploying on Base chains", chain.String(), chain.ChainSelector())
			continue
		}
		signerRegistry, err := cldf.DeployContract(e.Logger, chain, addressBook,
			func(chain cldf_evm.Chain) cldf.ContractDeploy[*signer_registry.SignerRegistry] {
				address, tx, signerRegistry, err := signer_registry.DeploySignerRegistry(
					chain.DeployerKey,
					chain.Client,
					big.NewInt(int64(config.MaxSigners)),
					config.Signers,
				)

				return cldf.ContractDeploy[*signer_registry.SignerRegistry]{
					Address:  address,
					Contract: signerRegistry,
					Tx:       tx,
					Tv:       cldf.NewTypeAndVersion(shared.EVMSignerRegistry, deployment.Version1_0_0),
					Err:      err,
				}
			},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy signer registry: %w", err)
		}

		e.Logger.Infof("Successfully deployed signer registry %s on %s", signerRegistry.Address.String(), chain.String())
	}

	return cldf.ChangesetOutput{AddressBook: addressBook}, nil
}
