package ocr2

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	starknetrpc "github.com/NethermindEth/starknet.go/rpc"
	starknetutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/starknet"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

//go:generate mockery --name OCR2Reader --output ./mocks/
type OCR2Reader interface { //nolint:revive
	LatestConfigDetails(context.Context, *felt.Felt) (ContractConfigDetails, error)
	LatestTransmissionDetails(context.Context, *felt.Felt) (TransmissionDetails, error)
	LatestRoundData(context.Context, *felt.Felt) (RoundData, error)
	LinkAvailableForPayment(context.Context, *felt.Felt) (*big.Int, error)
	ConfigFromEventAt(context.Context, *felt.Felt, uint64) (ContractConfig, error)
	NewTransmissionsFromEventsAt(context.Context, *felt.Felt, uint64) ([]NewTransmissionEvent, error)
	BillingDetails(context.Context, *felt.Felt) (BillingDetails, error)

	BaseReader() starknet.Reader
}

var _ OCR2Reader = (*Client)(nil)

type Client struct {
	r    starknet.Reader
	lggr logger.Logger
}

func NewClient(reader starknet.Reader, lggr logger.Logger) (*Client, error) {
	return &Client{
		r:    reader,
		lggr: lggr,
	}, nil
}

func (c *Client) BaseReader() starknet.Reader {
	return c.r
}

func (c *Client) BillingDetails(ctx context.Context, address *felt.Felt) (bd BillingDetails, err error) {
	ops := starknet.CallOps{
		ContractAddress: address,
		Selector:        starknetutils.GetSelectorFromNameFelt("billing"),
	}

	res, err := c.r.CallContract(ctx, ops)
	if err != nil {
		return bd, fmt.Errorf("couldn't call the contract: %w", err)
	}

	// [0] - observation payment, [1] - transmission payment, [2] - gas base, [3] - gas per signature
	if len(res) != 4 {
		return bd, errors.New("unexpected result length")
	}

	observationPayment := res[0].BigInt(big.NewInt(0))
	transmissionPayment := res[1].BigInt(big.NewInt(0))

	bd, err = NewBillingDetails(observationPayment, transmissionPayment)
	if err != nil {
		return bd, fmt.Errorf("couldn't initialize billing details: %w", err)
	}

	return
}

func (c *Client) LatestConfigDetails(ctx context.Context, address *felt.Felt) (ccd ContractConfigDetails, err error) {
	ops := starknet.CallOps{
		ContractAddress: address,
		Selector:        starknetutils.GetSelectorFromNameFelt("latest_config_details"),
	}

	res, err := c.r.CallContract(ctx, ops)
	if err != nil {
		return ccd, fmt.Errorf("couldn't call the contract: %w", err)
	}

	// [0] - config count, [1] - block number, [2] - config digest
	if len(res) != 3 {
		return ccd, errors.New("unexpected result length")
	}

	blockNum := res[1]
	configDigest := res[2]

	ccd, err = NewContractConfigDetails(blockNum.BigInt(big.NewInt((0))), configDigest.Bytes())
	if err != nil {
		return ccd, fmt.Errorf("couldn't initialize config details: %w", err)
	}

	return
}

func (c *Client) LatestTransmissionDetails(ctx context.Context, address *felt.Felt) (td TransmissionDetails, err error) {
	ops := starknet.CallOps{
		ContractAddress: address,
		Selector:        starknetutils.GetSelectorFromNameFelt("latest_transmission_details"),
	}

	res, err := c.r.CallContract(ctx, ops)
	if err != nil {
		return td, fmt.Errorf("couldn't call the contract: %w", err)
	}

	// [0] - config digest, [1] - epoch and round, [2] - latest answer, [3] - latest timestamp
	if len(res) != 4 {
		return td, errors.New("unexpected result length")
	}

	digest := res[0]
	configDigest := types.ConfigDigest{}
	digest.BigInt(big.NewInt(0)).FillBytes(configDigest[:])

	epoch, round := parseEpochAndRound(res[1].BigInt(big.NewInt(0)))

	latestAnswer := res[2].BigInt(big.NewInt(0))
	if err != nil {
		return td, fmt.Errorf("latestAnswer invalid: %w", err)
	}

	timestampFelt := res[3]
	// TODO: Int64() can return invalid data if int is too big
	unixTime := timestampFelt.BigInt(big.NewInt(0)).Int64()
	latestTimestamp := time.Unix(unixTime, 0)

	td = TransmissionDetails{
		Digest:          configDigest,
		Epoch:           epoch,
		Round:           round,
		LatestAnswer:    latestAnswer,
		LatestTimestamp: latestTimestamp,
	}

	return td, nil
}

func (c *Client) LatestRoundData(ctx context.Context, address *felt.Felt) (round RoundData, err error) {
	ops := starknet.CallOps{
		ContractAddress: address,
		Selector:        starknetutils.GetSelectorFromNameFelt("latest_round_data"),
	}

	felts, err := c.r.CallContract(ctx, ops)
	if err != nil {
		return round, fmt.Errorf("couldn't call the contract with selector latest_round_data: %w", err)
	}

	round, err = NewRoundData(felts)
	if err != nil {
		return round, fmt.Errorf("unable to decode RoundData: %w", err)
	}
	return round, nil
}

func (c *Client) LinkAvailableForPayment(ctx context.Context, address *felt.Felt) (*big.Int, error) {
	results, err := c.r.CallContract(ctx, starknet.CallOps{
		ContractAddress: address,
		Selector:        starknetutils.GetSelectorFromNameFelt("link_available_for_payment"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call the contract with selector 'link_available_for_payment': %w", err)
	}
	if l := len(results); l != 2 {
		return nil, fmt.Errorf("insufficient data from selector 'link_available_for_payment': need 2 results but got %d", l)
	}

	isNegative := !results[0].IsZero()
	ans := results[1].BigInt(big.NewInt(0))
	if isNegative {
		ans.Neg(ans)
	}

	return ans, nil
}

func (c *Client) fetchEventsFromBlock(ctx context.Context, address *felt.Felt, eventType string, blockNum uint64) (events []starknetrpc.EmittedEvent, err error) {
	block := starknetrpc.WithBlockNumber(blockNum)

	eventKey := starknetutils.GetSelectorFromNameFelt(eventType)

	input := starknetrpc.EventsInput{
		EventFilter: starknetrpc.EventFilter{
			FromBlock: block,
			ToBlock:   block,
			Address:   address,
			Keys:      [][]*felt.Felt{{eventKey}}, // skip other event types
			// PageSize:   0,
			// PageNumber: 0,
		},
		ResultPageRequest: starknetrpc.ResultPageRequest{
			// ContinuationToken: ,
			ChunkSize: 10,
		},
	}
	chunk, err := c.r.Events(ctx, input)
	events = chunk.Events

	// TODO: check events.isLastPage, query more if needed

	if err != nil {
		return events, fmt.Errorf("couldn't fetch events for block: %w", err)
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("events not found in the block %d", blockNum)
	}
	return events, nil
}

func (c *Client) ConfigFromEventAt(ctx context.Context, address *felt.Felt, blockNum uint64) (cc ContractConfig, err error) {
	events, err := c.fetchEventsFromBlock(ctx, address, "ConfigSet", blockNum)
	if err != nil {
		return cc, fmt.Errorf("failed to fetch config_set events: %w", err)
	}
	if len(events) != 1 {
		return cc, fmt.Errorf("expected to find one config_set event in block %d for address %s but found %d", blockNum, address, len(events))
	}
	configAtEvent := events[0]
	config, err := ParseConfigSetEvent(configAtEvent)
	if err != nil {
		return cc, fmt.Errorf("couldn't parse config event: %w", err)
	}
	return ContractConfig{
		Config:      config,
		ConfigBlock: blockNum,
	}, nil
}

// NewTransmissionsFromEventsAt finds events of type new_transmission emitted by the contract address in a given block number.
func (c *Client) NewTransmissionsFromEventsAt(ctx context.Context, address *felt.Felt, blockNum uint64) (events []NewTransmissionEvent, err error) {
	rawEvents, err := c.fetchEventsFromBlock(ctx, address, "NewTransmission", blockNum)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch new_transmission events: %w", err)
	}
	if len(rawEvents) == 0 {
		return nil, fmt.Errorf("expected to find at least one new_transmission event in block %d for address %s but found %d", blockNum, address, len(rawEvents))
	}
	events = []NewTransmissionEvent{}
	for _, rawEvent := range rawEvents {
		event, err := ParseNewTransmissionEvent(rawEvent)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse new_transmission event: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}
