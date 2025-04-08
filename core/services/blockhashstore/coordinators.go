package blockhashstore

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	v1 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/solidity_vrf_coordinator_interface"
	v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	v2plus "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

var (
	_ Coordinator = MultiCoordinator{}
	_ Coordinator = &V1Coordinator{}
	_ Coordinator = &V2Coordinator{}
	_ Coordinator = &V2PlusCoordinator{}
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
	c  v1.VRFCoordinatorInterface
	lp logpoller.LogPoller
}

// NewV1Coordinator creates a new V1Coordinator from the given contract.
func NewV1Coordinator(ctx context.Context, c v1.VRFCoordinatorInterface, lp logpoller.LogPoller) (*V1Coordinator, error) {
	err := lp.RegisterFilter(ctx, logpoller.Filter{
		Name: logpoller.FilterName("VRFv1CoordinatorFeeder", c.Address()),
		EventSigs: []common.Hash{
			v1.VRFCoordinatorRandomnessRequest{}.Topic(),
			v1.VRFCoordinatorRandomnessRequestFulfilled{}.Topic(),
		}, Addresses: []common.Address{c.Address()},
	})
	if err != nil {
		return nil, err
	}
	return &V1Coordinator{c, lp}, nil
}

// Requests satisfies the Coordinator interface.
func (v *V1Coordinator) Requests(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) ([]Event, error) {
	logs, err := v.lp.LogsWithSigs(
		ctx,
		int64(fromBlock),
		int64(toBlock),
		[]common.Hash{
			v1.VRFCoordinatorRandomnessRequest{}.Topic(),
		},
		v.c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "filter v1 requests")
	}

	var reqs []Event
	for _, l := range logs {
		requestLog, err := v.c.ParseLog(l.ToGethLog())
		if err != nil {
			continue // malformed log should not break flow
		}
		request, ok := requestLog.(*v1.VRFCoordinatorRandomnessRequest)
		if !ok {
			continue // malformed log should not break flow
		}
		reqs = append(reqs, Event{ID: hex.EncodeToString(request.RequestID[:]), Block: request.Raw.BlockNumber})
	}

	return reqs, nil
}

// Fulfillments satisfies the Coordinator interface.
func (v *V1Coordinator) Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error) {
	toBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching latest block")
	}

	logs, err := v.lp.LogsWithSigs(
		ctx,
		int64(fromBlock),
		toBlock.BlockNumber,
		[]common.Hash{
			v1.VRFCoordinatorRandomnessRequestFulfilled{}.Topic(),
		},
		v.c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "filter v1 fulfillments")
	}

	var fuls []Event
	for _, l := range logs {
		requestLog, err := v.c.ParseLog(l.ToGethLog())
		if err != nil {
			continue // malformed log should not break flow
		}
		request, ok := requestLog.(*v1.VRFCoordinatorRandomnessRequestFulfilled)
		if !ok {
			continue // malformed log should not break flow
		}
		fuls = append(fuls, Event{ID: hex.EncodeToString(request.RequestId[:]), Block: request.Raw.BlockNumber})
	}
	return fuls, nil
}

// V2Coordinator fetches request and fulfillment logs from a VRF V2 coordinator contract.
type V2Coordinator struct {
	c  v2.VRFCoordinatorV2Interface
	lp logpoller.LogPoller
}

// NewV2Coordinator creates a new V2Coordinator from the given contract.
func NewV2Coordinator(ctx context.Context, c v2.VRFCoordinatorV2Interface, lp logpoller.LogPoller) (*V2Coordinator, error) {
	err := lp.RegisterFilter(ctx, logpoller.Filter{
		Name: logpoller.FilterName("VRFv2CoordinatorFeeder", c.Address()),
		EventSigs: []common.Hash{
			v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(),
			v2.VRFCoordinatorV2RandomWordsFulfilled{}.Topic(),
		}, Addresses: []common.Address{c.Address()},
	})

	if err != nil {
		return nil, err
	}

	return &V2Coordinator{c, lp}, err
}

// Requests satisfies the Coordinator interface.
func (v *V2Coordinator) Requests(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) ([]Event, error) {
	logs, err := v.lp.LogsWithSigs(
		ctx,
		int64(fromBlock),
		int64(toBlock),
		[]common.Hash{
			v2.VRFCoordinatorV2RandomWordsRequested{}.Topic(),
		},
		v.c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "filter v2 requests")
	}

	var reqs []Event
	for _, l := range logs {
		requestLog, err := v.c.ParseLog(l.ToGethLog())
		if err != nil {
			continue // malformed log should not break flow
		}
		request, ok := requestLog.(*v2.VRFCoordinatorV2RandomWordsRequested)
		if !ok {
			continue // malformed log should not break flow
		}
		reqs = append(reqs, Event{ID: request.RequestId.String(), Block: request.Raw.BlockNumber})
	}

	return reqs, nil
}

// Fulfillments satisfies the Coordinator interface.
func (v *V2Coordinator) Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error) {
	toBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching latest block")
	}

	logs, err := v.lp.LogsWithSigs(
		ctx,
		int64(fromBlock),
		toBlock.BlockNumber,
		[]common.Hash{
			v2.VRFCoordinatorV2RandomWordsFulfilled{}.Topic(),
		},
		v.c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "filter v2 fulfillments")
	}

	var fuls []Event
	for _, l := range logs {
		requestLog, err := v.c.ParseLog(l.ToGethLog())
		if err != nil {
			continue // malformed log should not break flow
		}
		request, ok := requestLog.(*v2.VRFCoordinatorV2RandomWordsFulfilled)
		if !ok {
			continue // malformed log should not break flow
		}
		fuls = append(fuls, Event{ID: request.RequestId.String(), Block: request.Raw.BlockNumber})
	}
	return fuls, nil
}

// V2PlusCoordinator fetches request and fulfillment logs from a VRF V2Plus coordinator contract.
type V2PlusCoordinator struct {
	c  v2plus.IVRFCoordinatorV2PlusInternalInterface
	lp logpoller.LogPoller
}

// NewV2Coordinator creates a new V2Coordinator from the given contract.
func NewV2PlusCoordinator(ctx context.Context, c v2plus.IVRFCoordinatorV2PlusInternalInterface, lp logpoller.LogPoller) (*V2PlusCoordinator, error) {
	err := lp.RegisterFilter(ctx, logpoller.Filter{
		Name: logpoller.FilterName("VRFv2PlusCoordinatorFeeder", c.Address()),
		EventSigs: []common.Hash{
			v2plus.IVRFCoordinatorV2PlusInternalRandomWordsRequested{}.Topic(),
			v2plus.IVRFCoordinatorV2PlusInternalRandomWordsFulfilled{}.Topic(),
		}, Addresses: []common.Address{c.Address()},
	})

	if err != nil {
		return nil, err
	}

	return &V2PlusCoordinator{c, lp}, err
}

// Requests satisfies the Coordinator interface.
func (v *V2PlusCoordinator) Requests(
	ctx context.Context,
	fromBlock uint64,
	toBlock uint64,
) ([]Event, error) {
	logs, err := v.lp.LogsWithSigs(
		ctx,
		int64(fromBlock),
		int64(toBlock),
		[]common.Hash{
			v2plus.IVRFCoordinatorV2PlusInternalRandomWordsRequested{}.Topic(),
		},
		v.c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "filter v2 requests")
	}

	var reqs []Event
	for _, l := range logs {
		requestLog, err := v.c.ParseLog(l.ToGethLog())
		if err != nil {
			continue // malformed log should not break flow
		}
		request, ok := requestLog.(*v2plus.IVRFCoordinatorV2PlusInternalRandomWordsRequested)
		if !ok {
			continue // malformed log should not break flow
		}
		reqs = append(reqs, Event{ID: request.RequestId.String(), Block: request.Raw.BlockNumber})
	}

	return reqs, nil
}

// Fulfillments satisfies the Coordinator interface.
func (v *V2PlusCoordinator) Fulfillments(ctx context.Context, fromBlock uint64) ([]Event, error) {
	toBlock, err := v.lp.LatestBlock(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching latest block")
	}

	logs, err := v.lp.LogsWithSigs(
		ctx,
		int64(fromBlock),
		toBlock.BlockNumber,
		[]common.Hash{
			v2plus.IVRFCoordinatorV2PlusInternalRandomWordsFulfilled{}.Topic(),
		},
		v.c.Address())
	if err != nil {
		return nil, errors.Wrap(err, "filter v2 fulfillments")
	}

	var fuls []Event
	for _, l := range logs {
		requestLog, err := v.c.ParseLog(l.ToGethLog())
		if err != nil {
			continue // malformed log should not break flow
		}
		request, ok := requestLog.(*v2plus.IVRFCoordinatorV2PlusInternalRandomWordsFulfilled)
		if !ok {
			continue // malformed log should not break flow
		}
		fuls = append(fuls, Event{ID: request.RequestId.String(), Block: request.Raw.BlockNumber})
	}
	return fuls, nil
}
