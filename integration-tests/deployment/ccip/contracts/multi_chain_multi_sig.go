package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	ccipdeployment "github.com/smartcontractkit/chainlink/integration-tests/deployment/ccip"
	"math/big"
)

//type ManyChainMultiSig interface {
//	GetOpCount() (*big.Int, error)
//	//TODO remove common.Address
//	SetConfig(deployerKey any, signerAddresses []common.Address, signerGroups []uint8, groupQuorums [32]uint8, groupParents [32]uint8, clearRoot bool) (string, error)
//}

type Tx interface {
}

func NewManyChainMultiSig(address ccipdeployment.ContractAddress) *ManyChainMultiSig {
	return &ManyChainMultiSig{
		Address: address,
	}
}

func NewManyChainMultiSigAdapter(contract *owner_helpers.ManyChainMultiSig) *ManyChainMultiSig {
	return &ManyChainMultiSig{
		evmContract: contract,
	}
}

type ManyChainMultiSig struct {
	Address     ccipdeployment.ContractAddress
	evmContract *owner_helpers.ManyChainMultiSig
}

func (m ManyChainMultiSig) GetOpCount() (*big.Int, error) {
	if m.evmContract != nil {
		return m.evmContract.GetOpCount(nil)
	}
	panic("implement me")
}

func (m ManyChainMultiSig) SetConfig(deployerKey any, signerAddresses []common.Address, signerGroups []uint8, groupQuorums [32]uint8, groupParents [32]uint8, clearRoot bool) (string, error) {
	if m.evmContract != nil {
		m.evmContract.SetConfig(deployerKey.(*bind.TransactOpts), signerAddresses, signerGroups, groupQuorums, groupParents, clearRoot)
	}
	panic("implement me")
}
