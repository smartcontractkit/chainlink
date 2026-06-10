package changeset

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// acceptOwnershipABI is the minimal ABI needed to call acceptOwnership on any ownable contract.
var acceptOwnershipABI, _ = abi.JSON(strings.NewReader(`[{"name":"acceptOwnership","type":"function","stateMutability":"nonpayable","inputs":[],"outputs":[]}]`))

// AcceptOwnershipInput specifies a contract address on which acceptOwnership should be called.
type AcceptOwnershipInput struct {
	ChainSelector   uint64 `json:"chainSelector" yaml:"chainSelector"`
	ContractAddress string `json:"contractAddress" yaml:"contractAddress"`
}

// AcceptOwnershipEOA calls acceptOwnership on an arbitrary contract using the deployer key.
// This is for non-MCMS use only — it completes an EOA-to-EOA ownership transfer.
// Do not use when the pending owner is a timelock/MCMS contract.
type AcceptOwnershipEOA struct{}

var _ cldf.ChangeSetV2[AcceptOwnershipInput] = AcceptOwnershipEOA{}

func (AcceptOwnershipEOA) VerifyPreconditions(e cldf.Environment, input AcceptOwnershipInput) error {
	if _, ok := e.BlockChains.EVMChains()[input.ChainSelector]; !ok {
		return fmt.Errorf("chain selector %d not found in environment", input.ChainSelector)
	}
	if input.ContractAddress == "" {
		return errors.New("contractAddress is required")
	}
	return nil
}

func (AcceptOwnershipEOA) Apply(e cldf.Environment, input AcceptOwnershipInput) (cldf.ChangesetOutput, error) {
	chain := e.BlockChains.EVMChains()[input.ChainSelector]
	addr := common.HexToAddress(input.ContractAddress)

	bc := bind.NewBoundContract(addr, acceptOwnershipABI, chain.Client, chain.Client, chain.Client)
	tx, err := bc.Transact(chain.DeployerKey, "acceptOwnership")
	if _, confErr := cldf.ConfirmIfNoError(chain, tx, err); confErr != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("AcceptOwnership failed for %s on chain %d: %w", input.ContractAddress, input.ChainSelector, confErr)
	}

	e.Logger.Infow("Accepted ownership", "address", input.ContractAddress, "chainSelector", input.ChainSelector)
	return cldf.ChangesetOutput{}, nil
}
