package ocr2

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/juno/core/felt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	starknetrpc "github.com/NethermindEth/starknet.go/rpc"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2/medianreport"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"
)

var (
	MaxObservers = 31

	// Event selectors
	NewTransmissionEventSelector = "019e22f866f4c5aead2809bf160d2b29e921e335d899979732101c6f3c38ff81"
	ConfigSetEventSelector       = "9a144bf4a6a8fd083c93211e163e59221578efcc86b93f8c97c620e7b9608a"
)

// NewTransmissionEvent represents the 'NewTransmission' event
type NewTransmissionEvent struct {
	RoundId         uint32 //nolint:revive
	LatestAnswer    *big.Int
	Transmitter     *felt.Felt
	LatestTimestamp time.Time
	Observers       []uint8
	ObservationsLen uint32
	Observations    []*big.Int
	JuelsPerFeeCoin *big.Int
	GasPrice        *big.Int
	ConfigDigest    types.ConfigDigest
	Epoch           uint32
	Round           uint8
	Reimbursement   *big.Int
}

// ParseNewTransmissionEvent is decoding binary felt data as the NewTransmissionEvent type
func ParseNewTransmissionEvent(event starknetrpc.EmittedEvent) (NewTransmissionEvent, error) {
	eventData := event.Data
	{
		const observationsLenIndex = 3
		const constNumOfElements = 9
		const constNumOfKeys = 2 + 1 // additional 1 for the automatic event ID key

		if len(eventData) < constNumOfElements {
			return NewTransmissionEvent{}, errors.New("invalid: event data")
		}

		if len(event.Keys) < constNumOfKeys {
			return NewTransmissionEvent{}, errors.New("invalid: event data")
		}

		observationsLen := eventData[observationsLenIndex].BigInt(big.NewInt(0)).Uint64()
		if len(eventData) != constNumOfElements+int(observationsLen) {
			return NewTransmissionEvent{}, errors.New("invalid: event data")
		}
	}

	// keys[0] == event_id
	// round_id
	roundID := uint32(event.Keys[1].BigInt(big.NewInt(0)).Uint64())
	// transmitter
	transmitter := event.Keys[2]

	index := 0

	// answer
	latestAnswer := eventData[index].BigInt(big.NewInt(0))

	// observation_timestamp
	index++
	unixTime := eventData[index].BigInt(big.NewInt(0)).Int64()
	latestTimestamp := time.Unix(unixTime, 0)

	// observers (raw) max 31
	index++
	observersRaw := starknet.PadBytes(eventData[index].BigInt(big.NewInt(0)).Bytes(), MaxObservers)

	// observation_len
	index++
	observationsLen := uint32(eventData[index].BigInt(big.NewInt(0)).Uint64())

	// observers (based on observationsLen)
	var observers []uint8
	for i := 0; i < int(observationsLen); i++ {
		observers = append(observers, observersRaw[i])
	}

	// observations (based on observationsLen)
	var observations []*big.Int
	for i := 0; i < int(observationsLen); i++ {
		observations = append(observations, eventData[index+i+1].BigInt(big.NewInt(0)))
	}

	// juels_per_fee_coin
	index += int(observationsLen) + 1
	juelsPerFeeCoin := eventData[index].BigInt(big.NewInt(0))

	// juels_per_fee_coin
	index++
	gasPrice := eventData[index].BigInt(big.NewInt(0))

	// config digest
	index++
	digest := eventData[index].Bytes()

	// epoch_and_round
	index++
	epoch, round := parseEpochAndRound(eventData[index].BigInt(big.NewInt(0)))

	// reimbursement
	index++
	reimbursement := eventData[index].BigInt(big.NewInt(0))

	return NewTransmissionEvent{
		RoundId:         roundID,
		LatestAnswer:    latestAnswer,
		Transmitter:     transmitter,
		LatestTimestamp: latestTimestamp,
		Observers:       observers,
		ObservationsLen: observationsLen,
		Observations:    observations,
		JuelsPerFeeCoin: juelsPerFeeCoin,
		GasPrice:        gasPrice,
		ConfigDigest:    digest,
		Epoch:           epoch,
		Round:           round,
		Reimbursement:   reimbursement,
	}, nil
}

// ParseConfigSetEvent is decoding binary felt data as the libocr ContractConfig type
func ParseConfigSetEvent(event starknetrpc.EmittedEvent) (types.ContractConfig, error) {
	eventData := event.Data
	{
		const oraclesLenIdx = 1
		if len(eventData) < oraclesLenIdx {
			return types.ContractConfig{}, errors.New("invalid: event data")
		}

		oraclesLen := eventData[oraclesLenIdx].BigInt(big.NewInt(0)).Uint64()
		onchainConfigLenIdx := oraclesLenIdx + 2*oraclesLen + 2

		if uint64(len(eventData)) < onchainConfigLenIdx {
			return types.ContractConfig{}, errors.New("invalid: event data")
		}

		onchainConfigLen := eventData[onchainConfigLenIdx].BigInt(big.NewInt(0)).Uint64()
		offchainConfigLenIdx := onchainConfigLenIdx + onchainConfigLen + 2

		if uint64(len(eventData)) < offchainConfigLenIdx {
			return types.ContractConfig{}, errors.New("invalid: event data")
		}

		offchainConfigLen := eventData[offchainConfigLenIdx].BigInt(big.NewInt(0)).Uint64()
		if uint64(len(eventData)) != offchainConfigLenIdx+offchainConfigLen+1 {
			return types.ContractConfig{}, errors.New("invalid: event data")
		}
	}

	// keys[0] == event_id
	// keys[1] == previous_config_block_number - skip

	// latest_config_digest
	digest := event.Keys[2].Bytes()

	index := 0

	// config_count
	configCount := eventData[index].BigInt(big.NewInt(0)).Uint64()

	// oracles_len
	index++
	oraclesLen := eventData[index].BigInt(big.NewInt(0)).Uint64()

	// oracles
	index++
	oracleMembers := eventData[index:(index + int(oraclesLen)*2)]
	var signers []types.OnchainPublicKey
	var transmitters []types.Account
	for i, member := range oracleMembers {
		if i%2 == 0 {
			b := member.Bytes()
			signers = append(signers, b[:]) // pad to 32 bytes
		} else {
			transmitters = append(transmitters, types.Account(member.String()))
		}
	}

	// f
	index = index + int(oraclesLen)*2
	f := eventData[index].BigInt(big.NewInt(0)).Uint64()

	// onchain_config length
	index++
	onchainConfigLen := eventData[index].BigInt(big.NewInt(0)).Uint64()

	// onchain_config (version=1, min, max)
	index++
	onchainConfigFelts := eventData[index:(index + int(onchainConfigLen))]
	onchainConfig, err := medianreport.OnchainConfigCodec{}.EncodeFromFelt(
		onchainConfigFelts[0].BigInt(big.NewInt(0)),
		onchainConfigFelts[1].BigInt(big.NewInt(0)),
		onchainConfigFelts[2].BigInt(big.NewInt(0)),
	)
	if err != nil {
		return types.ContractConfig{}, fmt.Errorf("err in encoding onchain config from felts: %w", err)
	}

	// offchain_config_version
	index += int(onchainConfigLen)
	offchainConfigVersion := eventData[index].BigInt(big.NewInt(0)).Uint64()

	// offchain_config_len
	index++
	offchainConfigLen := eventData[index].BigInt(big.NewInt(0)).Uint64()

	// offchain_config
	index++
	offchainConfigFelts := eventData[index:(index + int(offchainConfigLen))]
	offchainConfig, err := starknet.DecodeFelts(starknet.FeltsToBig(offchainConfigFelts))
	if err != nil {
		return types.ContractConfig{}, fmt.Errorf("couldn't decode offchain config: %w", err)
	}

	return types.ContractConfig{
		ConfigDigest:          digest,
		ConfigCount:           configCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     uint8(f),
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}
