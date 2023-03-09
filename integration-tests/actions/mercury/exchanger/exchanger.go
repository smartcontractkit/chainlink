package exchanger

import (
	"math/big"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Order struct {
	FeedID       [32]byte
	CurrencySrc  [32]byte
	CurrencyDst  [32]byte
	AmountSrc    *big.Int
	MinAmountDst *big.Int
	Sender       common.Address
	Receiver     common.Address
}

func CreateEncodedCommitment(order Order) ([]byte, error) {
	// bytes32 feedID, bytes32 currencySrc, bytes32 currencyDst, uint256 amountSrc, uint256 minAmountDst, address sender, address receiver
	orderType, _ := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "feedID", Type: "bytes32"},
		{Name: "currencySrc", Type: "bytes32"},
		{Name: "currencyDst", Type: "bytes32"},
		{Name: "amountSrc", Type: "uint256"},
		{Name: "minAmountDst", Type: "uint256"},
		{Name: "sender", Type: "address"},
		{Name: "receiver", Type: "address"},
	})
	var args abi.Arguments = []abi.Argument{{Type: orderType}}
	return args.Pack(order)
}

func CreateCommitmentHash(order Order) common.Hash {
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "", nil)
	addressTy, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{
			Type: bytes32Ty,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: uint256Ty,
		},
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
	}

	bytes, _ := arguments.Pack(
		order.FeedID,
		order.CurrencySrc,
		order.CurrencyDst,
		order.AmountSrc,
		order.MinAmountDst,
		order.Sender,
		order.Receiver,
	)

	return crypto.Keccak256Hash(bytes)
}
