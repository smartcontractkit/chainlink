package withdrawprover

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l2_to_l1_message_passer"
)

var (
	l2ToL1MessagePasserABI *abi.ABI
)

func init() {
	abi, err := optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMetaData.GetAbi()
	if err != nil {
		panic(err)
	}

	// check that we have the message passed event
	_, ok := abi.Events["MessagePassed"]
	if !ok {
		panic("OptimismL2ToL1MessagePasser gethwrapper ABI does not contain MessagePassed event")
	}

	l2ToL1MessagePasserABI = abi
}

func hashMessageHash(h [32]byte) ([32]byte, error) {
	var zeroHash [32]byte
	encoded, err := utils.ABIEncode(`[{"type": "bytes32"}, {"type": "bytes32"}]`, h, zeroHash)
	if err != nil {
		return [32]byte{}, err
	}
	return crypto.Keccak256Hash(encoded), nil
}

func GetMessagePassedLog(logs []*gethtypes.Log) *gethtypes.Log {
	for _, lg := range logs {
		if lg.Topics[0] == MessagePassedTopic {
			return lg
		}
	}
	return nil
}

func ParseMessagePassedLog(log *gethtypes.Log) (*optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed, error) {
	// address doesn't matter for parsing, and neither does the contract backend.
	cdm, err := optimism_l2_to_l1_message_passer.NewOptimismL2ToL1MessagePasser(common.HexToAddress("0x0"), nil)
	if err != nil {
		return nil, err
	}

	messagePassed, err := cdm.ParseMessagePassed(*log)
	if err != nil {
		return nil, err
	}

	return messagePassed, nil
}

func hashLowLevelMessage(llm *optimism_l2_to_l1_message_passer.OptimismL2ToL1MessagePasserMessagePassed) ([32]byte, error) {
	event, ok := l2ToL1MessagePasserABI.Events["MessagePassed"]
	// should be impossible, checked in init
	if !ok {
		return [32]byte{}, fmt.Errorf("event MessagePassed not found in L2ToL1MessagePasser ABI")
	}

	// last event argument is withdrawalHash, but we don't want to encode that.
	var arguments abi.Arguments
	for i, arg := range event.Inputs {
		if i == len(event.Inputs)-1 {
			break
		}
		arguments = append(arguments, arg)
	}

	encoded, err := arguments.Pack(llm.Nonce, llm.Sender, llm.Target, llm.Value, llm.GasLimit, llm.Data)
	if err != nil {
		return [32]byte{}, err
	}

	return crypto.Keccak256Hash(encoded), nil
}

func toProofBytes(proof []hexutil.Bytes) [][]byte {
	proofBytes := make([][]byte, len(proof))
	for i, p := range proof {
		proofBytes[i] = p
	}
	return proofBytes
}
