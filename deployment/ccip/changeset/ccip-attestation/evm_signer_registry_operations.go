package ccip_attestation

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/evm/operations/contract"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
)

var addSignersOperation = contract.NewWrite(contract.WriteParams[[]Signer, *signer_registry.SignerRegistry]{
	Name:            "ccip-attestation:signer-registry:add-signers",
	Version:         &deployment.Version1_0_0,
	Description:     "Adds signers to an EVM SignerRegistry",
	ContractType:    shared.EVMSignerRegistry,
	ContractABI:     signer_registry.SignerRegistryABI,
	NewContract:     signer_registry.NewSignerRegistry,
	IsAllowedCaller: contract.OnlyOwner[*signer_registry.SignerRegistry, []Signer],
	Validate:        validateSigners,
	CallContract: func(registry *signer_registry.SignerRegistry, opts *bind.TransactOpts, signers []Signer) (*types.Transaction, error) {
		return registry.AddSigners(opts, toContractSigners(signers))
	},
})

func validateSigners(signers []Signer) error {
	if len(signers) == 0 {
		return errors.New("no signers provided")
	}

	seen := make(map[common.Address]struct{}, len(signers)*2)
	for i, signer := range signers {
		if signer.EVMAddress == (common.Address{}) {
			return fmt.Errorf("signer %d has a zero EVM address", i)
		}
		if signer.EVMAddress == signer.NewEVMAddress {
			return fmt.Errorf("signer %d has identical active and pending EVM addresses", i)
		}
		if _, exists := seen[signer.EVMAddress]; exists {
			return fmt.Errorf("signer %d reuses EVM address %s", i, signer.EVMAddress.Hex())
		}
		seen[signer.EVMAddress] = struct{}{}

		if signer.NewEVMAddress == (common.Address{}) {
			continue
		}
		if _, exists := seen[signer.NewEVMAddress]; exists {
			return fmt.Errorf("signer %d reuses pending EVM address %s", i, signer.NewEVMAddress.Hex())
		}
		seen[signer.NewEVMAddress] = struct{}{}
	}

	return nil
}

func toContractSigners(signers []Signer) []signer_registry.ISignerRegistrySigner {
	contractSigners := make([]signer_registry.ISignerRegistrySigner, len(signers))
	for i, signer := range signers {
		contractSigners[i] = signer_registry.ISignerRegistrySigner{
			EvmAddress:    signer.EVMAddress,
			NewEVMAddress: signer.NewEVMAddress,
		}
	}
	return contractSigners
}
