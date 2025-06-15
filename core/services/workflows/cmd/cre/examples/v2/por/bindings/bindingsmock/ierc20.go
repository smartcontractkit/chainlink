package bindingsmock

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm/capabilitymock"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/cmd/cre/examples/v2/por/bindings"
)

type IERC20Mock struct {
	TotalSupply func() (*big.Int, error)
}

func NewIERC20Mock(address common.Address, clientMock *evmmock.ClientCapability) *IERC20Mock {
	erc20Mock := &IERC20Mock{}
	a := bindings.NewIERC20Abi()
	totalSupply := a.Methods["totalSupply"]
	funcMap := map[string]func([]byte) ([]byte, error){
		string(totalSupply.ID): func(payload []byte) ([]byte, error) {
			if (erc20Mock.TotalSupply) == nil {
				// TODO better if we can match the EVM's error
				return nil, errors.New("method not found on the contract")
			}

			result, err := erc20Mock.TotalSupply()
			if err != nil {
				return nil, err
			}
			return totalSupply.Outputs.Pack(result)
		},
	}
	bindings.AddInterfaceMock(address, clientMock, funcMap, nil)
	return erc20Mock
}
