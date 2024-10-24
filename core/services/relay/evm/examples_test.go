package evm_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	clientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func ExampleChainReaderService() {
	abi := `[{"inputs": [],"stateMutability": "nonpayable","type": "constructor"},{"inputs": [{"internalType": "address","name": "token","type": "address"},{"internalType": "uint16","name": "quantity","type": "uint16"}],"name": "contractRead","outputs": [{"internalType": "uint160","name": "","type": "uint160"}],"stateMutability": "pure","type": "function"}]`
	configJSON := fmt.Sprintf(`{
		"contracts": {
			"NamedContract": {
				"contractABI": "%s", // string value containing the entire encoded ABI
				"configs": {
					"contractReadName1": {
						"chainSpecificName": "contractRead", // can be different from read name or the same
						"readType": "method", // "method" for contract methods, "event" for contract events
						"confidenceConfirmations": {
							"unconfirmed": 1, // (optional) modify confidence levels as needed
							"finalized": 21 // if the concept of 'finalized' needs to be adjusted for a chain
						}
					}
				}
			}
		}
	}`, strings.ReplaceAll(abi, `"`, `\"`))

	var config types.ChainReaderConfig
	_ = json.Unmarshal([]byte(configJSON), &config)

	client := new(clientmocks.Client)
	poller := new(mocks.LogPoller)
	reader, _ := evm.NewChainReaderService(context.Background(), logger.NewWithSync(io.Discard), poller, nil, client, config)
	_ = reader.Start(context.Background())

	defer reader.Close()

	// ContractReader usage
	type ContractReadParameters struct {
		Address  []byte // use chain agnostic address type
		Quantity uint16
	}

	parameters := ContractReadParameters{
		Address:  []byte("0x2142"),
		Quantity: 100,
	}

	const readName = "contractReadName1"

	// the following BoundContract should be provided to Bind before calling GetLatestValue
	// this only needs to be done once per address
	binding := commontypes.BoundContract{
		Address: "0x4221",
		Name:    "NamedContract",
	}

	_ = reader.Bind(context.Background(), []commontypes.BoundContract{binding})

	identifier := binding.ReadIdentifier(readName)

	var contractReadResult *big.Int // another chain agnostic type

	_ = reader.GetLatestValue(context.Background(), identifier, primitives.Finalized, parameters, contractReadResult)
}
