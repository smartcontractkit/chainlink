package ocr2

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type ContractConfigDetails struct {
	Block  uint64
	Digest types.ConfigDigest
}

func NewContractConfigDetails(blockNum *big.Int, digest [32]byte) (ccd ContractConfigDetails, err error) {
	return ContractConfigDetails{
		Block:  blockNum.Uint64(),
		Digest: digest,
	}, nil
}

type ContractConfig struct {
	Config      types.ContractConfig
	ConfigBlock uint64
}

type TransmissionDetails struct {
	Digest          types.ConfigDigest
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp time.Time
}

type BillingDetails struct {
	ObservationPaymentGJuels  uint64
	TransmissionPaymentGJuels uint64
}

func NewBillingDetails(observationPaymentGJuels *big.Int, transmissionPaymentGJuels *big.Int) (bd BillingDetails, err error) {
	return BillingDetails{
		ObservationPaymentGJuels:  observationPaymentGJuels.Uint64(),
		TransmissionPaymentGJuels: transmissionPaymentGJuels.Uint64(),
	}, nil
}

type RoundData struct {
	RoundID     uint32
	Answer      *big.Int
	BlockNumber uint64
	StartedAt   time.Time
	UpdatedAt   time.Time
}

func NewRoundData(felts []*felt.Felt) (data RoundData, err error) {
	if len(felts) != 5 {
		return data, fmt.Errorf("expected number of felts to be 5 but got %d", len(felts))
	}
	roundID := felts[0].BigInt(big.NewInt(0))
	if !roundID.IsUint64() && roundID.Uint64() > math.MaxUint32 {
		return data, fmt.Errorf("aggregator round id does not fit in a uint32 '%s'", felts[0].String())
	}
	data.RoundID = uint32(roundID.Uint64())
	data.Answer = felts[1].BigInt(big.NewInt(0))
	blockNumber := felts[2].BigInt(big.NewInt(0))
	if !blockNumber.IsUint64() {
		return data, fmt.Errorf("block number '%s' does not fit into uint64", blockNumber.String())
	}
	data.BlockNumber = blockNumber.Uint64()
	startedAt := felts[3].BigInt(big.NewInt(0))
	if !startedAt.IsInt64() {
		return data, fmt.Errorf("startedAt '%s' does not fit into int64", startedAt.String())
	}
	data.StartedAt = time.Unix(startedAt.Int64(), 0)
	updatedAt := felts[4].BigInt(big.NewInt(0))
	if !updatedAt.IsInt64() {
		return data, fmt.Errorf("updatedAt '%s' does not fit into int64", startedAt.String())
	}
	data.UpdatedAt = time.Unix(updatedAt.Int64(), 0)
	return data, nil
}
