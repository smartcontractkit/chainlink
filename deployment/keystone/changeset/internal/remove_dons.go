package internal

import (
	"errors"
	"fmt"
	"math/big"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// RemoveDONsRequest holds the parameters for the RemoveDONs operation.
type RemoveDONsRequest struct {
	Chain                cldf_evm.Chain
	CapabilitiesRegistry *kcr.CapabilitiesRegistry
	DONs                 []uint32
	UseMCMS              bool
	// If UseMCMS is true and Ops is not nil then the RemoveDONs contract operation
	// will be added to the Ops.Batch.
	Ops *mcmstypes.BatchOperation
}

// RemoveDONsResponse represents the response from calling RemoveDONs.
type RemoveDONsResponse struct {
	// TxHash is the hash of the transaction if not using MCMS.
	TxHash string
	// Ops contains the MCMS operation (if any).
	Ops *mcmstypes.BatchOperation
}

func (r *RemoveDONsRequest) Validate() error {
	if len(r.DONs) == 0 {
		return errors.New("DONs list is empty")
	}

	if r.CapabilitiesRegistry == nil {
		return errors.New("registry is required")
	}
	return nil
}

// RemoveDONs calls the RemoveDONs method on the capabilities registry contract.
// It takes a list of DON IDs to remove.
func RemoveDONs(lggr logger.Logger, req *RemoveDONsRequest) (*RemoveDONsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	txOpts := req.Chain.DeployerKey
	for _, don := range req.DONs {
		if _, err := req.CapabilitiesRegistry.GetDON(nil, don); err != nil {
			return nil, fmt.Errorf("DON ID %d not found in registry: %w", don, err)
		}
	}

	if req.UseMCMS {
		txOpts = cldf.SimTransactOpts()
	}

	tx, err := req.CapabilitiesRegistry.RemoveDONs(txOpts, req.DONs)
	if err != nil {
		err = cldf.DecodeErr(kcr.CapabilitiesRegistryABI, err)
		return nil, fmt.Errorf("failed to call RemoveDONs: %w", err)
	}

	var ops mcmstypes.BatchOperation
	if !req.UseMCMS {
		_, err = req.Chain.Confirm(tx)
		if err != nil {
			return nil, fmt.Errorf("failed to confirm RemoveDONs transaction %s: %w", tx.Hash().String(), err)
		}
	} else {
		ops, err = proposalutils.BatchOperationForChain(req.Chain.Selector, req.CapabilitiesRegistry.Address().Hex(), tx.Data(), big.NewInt(0), string(CapabilitiesRegistry), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create batch operation: %w", err)
		}
	}

	return &RemoveDONsResponse{
		TxHash: tx.Hash().String(),
		Ops:    &ops,
	}, nil
}
