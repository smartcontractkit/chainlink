package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SubmissionReceivedEvent emitted from solidity contract as:
// event SubmissionReceived(
//   int256 indexed answer,
//   uint32 indexed round,
//   address indexed oracle
// );
type SubmissionReceivedEvent struct {
	Answer  string
	RoundID *big.Int
	Oracle  common.Address
}
