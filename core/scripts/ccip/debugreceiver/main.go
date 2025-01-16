package main

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ccipAny struct {
	SourceChainID *big.Int
	Sender        []byte
	Data          []byte
	Tokens        []common.Address
	Amounts       []*big.Int
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ABIEncode is the equivalent of abi.encode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIEncode(abiStr string, values ...interface{}) ([]byte, error) {
	// Create a dummy method with arguments
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "inputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}
	return res[4:], nil
}

// ABIDecode is the equivalent of abi.decode.
// See a full set of examples https://github.com/ethereum/go-ethereum/blob/420b78659bef661a83c5c442121b13f13288c09f/accounts/abi/packing_test.go#L31
func ABIDecode(abiStr string, data []byte) ([]interface{}, error) {
	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "outputs": %s}]`, abiStr)
	inAbi, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}
	return inAbi.Unpack("method", data)
}

func main() {
	// User inputs
	source, err := ethclient.Dial("TODO goerli URL")
	panicErr(err)
	dest, err := ethclient.Dial("TODO optimism URL")
	panicErr(err)
	var (
		requestBlock = int64(7916078)
		receiveBlock = int64(2518035)
		ccipReceiver = common.HexToAddress("0x9fE056F44510F970d724adA16903ba5D75CC4742")
	)

	log, err := source.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(requestBlock),
		ToBlock:   big.NewInt(requestBlock),
		Topics:    [][]common.Hash{{common.HexToHash("0x73dfb9df8214728e699dbaaf6ba97aa125afaaba83a5d0de7903062e7c5b3139")}}, // CCIPSendRequested
	})
	panicErr(err)
	encodedMsg, err := ABIDecode(`[{"components":[
{"internalType":"uint256","name":"sourceChainId","type":"uint256"},
{"internalType":"uint64","name":"sequenceNumber","type":"uint64"},
{"internalType":"address","name":"sender","type":"address"}, 
{"internalType":"address","name":"receiver","type":"address"},
{"internalType":"uint64","name":"nonce","type":"uint64"},
{"internalType":"bytes","name":"data","type":"bytes"},
{"internalType":"address[]","name":"tokens","type":"address[]"},
{"internalType":"uint256[]","name":"amounts","type":"uint256[]"},
{"internalType":"uint256","name":"gasLimit","type":"uint256"}],
"internalType":"structCCIP.EVM2EVMSubscriptionMessage","name":"message","type":"tuple"}]`,
		log[0].Data)
	panicErr(err)
	send := encodedMsg[0].(struct {
		SourceChainID  *big.Int         `json:"sourceChainId"`
		SequenceNumber uint64           `json:"sequenceNumber"`
		Sender         common.Address   `json:"sender"`
		Receiver       common.Address   `json:"receiver"`
		Nonce          uint64           `json:"nonce"`
		Data           []uint8          `json:"data"`
		Tokens         []common.Address `json:"tokens"`
		Amounts        []*big.Int       `json:"amounts"`
		GasLimit       *big.Int         `json:"gasLimit"`
	})
	sender, err := ABIEncode(`[{"type":"bytes", "name":"sender"}]`, send.Sender.Bytes())
	panicErr(err)
	any2evm, err := ABIEncode(`[{"components":[
{"internalType":"uint256","name":"sourceChainId","type":"uint256"}, 
{"internalType":"bytes","name":"sender","type":"bytes"}, 
{"internalType":"bytes","name":"data","type":"bytes"},
{"internalType":"address[]","name":"tokens","type":"address[]"},
{"internalType":"uint256[]","name":"amounts","type":"uint256[]"}],
"internalType":"structCCIP.Any2EVMMessage","name":"message","type":"tuple"}]`,
		ccipAny{send.SourceChainID, sender, send.Data, send.Tokens, send.Amounts})
	panicErr(err)
	a, err := dest.CallContract(context.Background(), ethereum.CallMsg{
		From:       common.HexToAddress("0x2b7ab40413da5077e168546ea376920591aee8e7"), // offramp router
		To:         &ccipReceiver,
		Gas:        send.GasLimit.Uint64(),
		GasPrice:   big.NewInt(0),
		Data:       bytes.Join([][]byte{hexutil.MustDecode("0xa0c6df15"), any2evm}, []byte{}), // ccipReceive selector
		AccessList: nil,
	}, big.NewInt(receiveBlock))
	fmt.Println(a, err)
}
