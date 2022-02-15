package blockhashstore

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"

	v1 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	v2 "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
)

var (
	_ Coordinator = MultiCoordinator{}
	_ Coordinator = &V1Coordinator{}
	_ Coordinator = &V2Coordinator{}
)

// MultiCoordinator combines the data from multiple coordinators.
type MultiCoordinator []Coordinator

// NewMultiCoordinator creates a new Coordinator that combines the results of the given
// coordinators.
func NewMultiCoordinator(coordinators ...Coordinator) Coordinator {
	if len(coordinators) == 1 {
		return coordinators[0]
	}
	return MultiCoordinator(coordinators)
}

// Requests satisfies the Coordinator interface.
func (m MultiCoordinator) Requests(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) ([]Event, error) {
	var reqs []Event
	for _, c := range m {
		r, err := c.Requests(ctx, fromBlock, toBlock)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		reqs = append(reqs, r...)
	}
	return reqs, nil
}

// Fulfillments satisfies the Coordinator interface.
func (m MultiCoordinator) Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error) {
	var fuls []Event
	for _, c := range m {
		f, err := c.Fulfillments(ctx, fromBlock)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		fuls = append(fuls, f...)
	}
	return fuls, nil
}

// V1Coordinator fetches request and fulfillment logs from a VRF V1 coordinator contract.
type V1Coordinator struct {
	c v1.VRFCoordinatorInterface
}

// NewV1Coordinator creates a new V1Coordinator from the given contract.
func NewV1Coordinator(c v1.VRFCoordinatorInterface) *V1Coordinator {
	return &V1Coordinator{c}
}

// Requests satisfies the Coordinator interface.
func (v *V1Coordinator) Requests(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) ([]Event, error) {
	iter, err := v.c.FilterRandomnessRequest(&bind.FilterOpts{
		Start:   fromBlock,
		End:     &toBlock,
		Context: ctx,
	}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "filter v1 requests")
	}
	defer iter.Close()
	var reqs []Event
	for iter.Next() {
		reqs = append(reqs, Event{
			ID:    hex.EncodeToString(iter.Event.RequestID[:]),
			Block: iter.Event.Raw.BlockNumber,
		})
	}
	return reqs, nil
}

// Fulfillments satisfies the Coordinator interface.
func (v *V1Coordinator) Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error) {
	iter, err := v.c.FilterRandomnessRequestFulfilled(&bind.FilterOpts{
		Start:   fromBlock,
		Context: ctx,
	})
	if err != nil {
		return nil, errors.Wrap(err, "filter v1 fulfillments")
	}
	defer iter.Close()
	var fuls []Event
	for iter.Next() {
		fuls = append(fuls, Event{
			ID:    hex.EncodeToString(iter.Event.RequestId[:]),
			Block: iter.Event.Raw.BlockNumber,
		})
	}
	return fuls, nil
}

// V2Coordinator fetches request and fulfillment logs from a VRF V2 coordinator contract.
type V2Coordinator struct {
	c v2.VRFCoordinatorV2Interface
}

// NewV2Coordinator creates a new V2Coordinator from the given contract.
func NewV2Coordinator(c v2.VRFCoordinatorV2Interface) *V2Coordinator {
	return &V2Coordinator{c}
}

// Requests satisfies the Coordinator interface.
func (v *V2Coordinator) Requests(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) ([]Event, error) {
	iter, err := v.c.FilterRandomWordsRequested(&bind.FilterOpts{
		Start:   fromBlock,
		End:     &toBlock,
		Context: ctx,
	}, nil, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "filter v2 requests")
	}
	defer iter.Close()
	var reqs []Event
	for iter.Next() {
		reqs = append(reqs, Event{
			ID:    iter.Event.RequestId.String(),
			Block: iter.Event.Raw.BlockNumber,
		})
	}
	return reqs, nil
}

// Fulfillments satisfies the Coordinator interface.
func (v *V2Coordinator) Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error) {
	iter, err := v.c.FilterRandomWordsFulfilled(&bind.FilterOpts{
		Start:   fromBlock,
		Context: ctx,
	}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "filter v2 fulfillments")
	}
	defer iter.Close()
	var fuls []Event
	for iter.Next() {
		fuls = append(fuls, Event{
			ID:    iter.Event.RequestId.String(),
			Block: iter.Event.Raw.BlockNumber,
		})
	}
	return fuls, nil
}
