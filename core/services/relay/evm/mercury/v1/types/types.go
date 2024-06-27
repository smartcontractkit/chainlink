package reporttypes

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var schema = GetSchema()

func GetSchema() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "feedId", Type: mustNewType("bytes32")},
		{Name: "observationsTimestamp", Type: mustNewType("uint32")},
		{Name: "benchmarkPrice", Type: mustNewType("int192")},
		{Name: "bid", Type: mustNewType("int192")},
		{Name: "ask", Type: mustNewType("int192")},
		{Name: "currentBlockNum", Type: mustNewType("uint64")},
		{Name: "currentBlockHash", Type: mustNewType("bytes32")},
		{Name: "validFromBlockNum", Type: mustNewType("uint64")},
		{Name: "currentBlockTimestamp", Type: mustNewType("uint64")},
	})
}

type Report struct {
	FeedId                [32]byte
	ObservationsTimestamp uint32
	BenchmarkPrice        *big.Int
	Bid                   *big.Int
	Ask                   *big.Int
	CurrentBlockNum       uint64
	CurrentBlockHash      [32]byte
	ValidFromBlockNum     uint64
	CurrentBlockTimestamp uint64
}

// Decode is made available to external users (i.e. mercury server)
func Decode(report []byte) (*Report, error) {
	values, err := schema.Unpack(report)
	if err != nil {
		return nil, fmt.Errorf("failed to decode report: %w", err)
	}
	decoded := new(Report)
	if err = schema.Copy(decoded, values); err != nil {
		return nil, fmt.Errorf("failed to copy report values to struct: %w", err)
	}
	return decoded, nil
}
