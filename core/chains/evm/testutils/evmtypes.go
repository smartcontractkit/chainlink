package testutils

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	mrand "math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// FixtureChainID matches the chain always added by fixtures.sql
// It is set to 0 since no real chain ever has this ID and allows a virtual
// "test" chain ID to be used without clashes
var FixtureChainID = big.NewInt(0)

// SimulatedChainID is the chain ID for the go-ethereum simulated backend
var SimulatedChainID = big.NewInt(1337)

// NewRandomEVMChainID returns a suitable random chain ID that will not conflict
// with fixtures
func NewRandomEVMChainID() *big.Int {
	id := mrand.Int63n(math.MaxInt32) + 10000
	return big.NewInt(id)
}

// NewAddress return a random new address
func NewAddress() common.Address {
	return common.BytesToAddress(randomBytes(20))
}

// NewHash return random Keccak256
func NewHash() common.Hash {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(b)
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// Head given the value convert it into an Head
func Head(val interface{}) *evmtypes.Head {
	var h evmtypes.Head
	time := uint64(0)
	switch t := val.(type) {
	case int:
		h = evmtypes.NewHead(big.NewInt(int64(t)), evmutils.NewHash(), evmutils.NewHash(), time, ubig.New(FixtureChainID))
	case uint64:
		h = evmtypes.NewHead(big.NewInt(int64(t)), evmutils.NewHash(), evmutils.NewHash(), time, ubig.New(FixtureChainID))
	case int64:
		h = evmtypes.NewHead(big.NewInt(t), evmutils.NewHash(), evmutils.NewHash(), time, ubig.New(FixtureChainID))
	case *big.Int:
		h = evmtypes.NewHead(t, evmutils.NewHash(), evmutils.NewHash(), time, ubig.New(FixtureChainID))
	default:
		panic(fmt.Sprintf("Could not convert %v of type %T to Head", val, val))
	}
	return &h
}

func NewLegacyTransaction(nonce uint64, to common.Address, value *big.Int, gasLimit uint32, gasPrice *big.Int, data []byte) *types.Transaction {
	tx := types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      uint64(gasLimit),
		GasPrice: gasPrice,
		Data:     data,
	}
	return types.NewTx(&tx)
}

func NewAddressPtr() *common.Address {
	a := common.BytesToAddress(randomBytes(20))
	return &a
}
