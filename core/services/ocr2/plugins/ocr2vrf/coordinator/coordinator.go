package coordinator

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	dkg_wrapper "github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	vrf_wrapper "github.com/smartcontractkit/ocr2vrf/gethwrappers/vrfbeaconcoordinator"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var _ ocr2vrftypes.CoordinatorInterface = &coordinator{}

var (
	dkgABI = evmtypes.MustGetABI(dkg_wrapper.DKGMetaData.ABI)
	vrfABI = evmtypes.MustGetABI(vrf_wrapper.VRFBeaconCoordinatorMetaData.ABI)
)

const (
	// VRF-only events.
	randomnessRequestedEvent            string = "RandomnessRequested"
	randomnessFulfillmentRequestedEvent        = "RandomnessFulfillmentRequested"
	randomWordsFulfilledEvent                  = "RandomWordsFulfilled"
	newTransmissionEvent                       = "NewTransmission"

	// Both VRF and DKG contracts emit this, it's an OCR event.
	configSetEvent = "ConfigSet"
)

type coordinator struct {
	lggr logger.Logger

	lp logpoller.LogPoller
	topics
	lookbackBlocks int64

	coordinatorContract VRFBeaconCoordinator
	coordinatorAddress  common.Address

	// We need to keep track of DKG ConfigSet events as well.
	dkgAddress common.Address

	evmClient evmclient.Client
	orm       ORM
}

// New creates a new CoordinatorInterface implementor.
func New(
	lggr logger.Logger,
	coordinatorAddress common.Address,
	dkgAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
	logPoller logpoller.LogPoller,
	orm ORM,
) (ocr2vrftypes.CoordinatorInterface, error) {
	coordinatorContract, err := vrf_wrapper.NewVRFBeaconCoordinator(coordinatorAddress, client)
	if err != nil {
		return nil, errors.Wrap(err, "coordinator wrapper creation")
	}

	t, err := newTopics()

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	logPoller.MergeFilter([]common.Hash{
		t.randomnessRequestedTopic,
		t.randomnessFulfillmentRequestedTopic,
		t.randomWordsFulfilledTopic,
		t.configSetTopic,
		t.newTransmissionTopic,
	}, coordinatorAddress)

	// We need ConfigSet events from the DKG contract as well.
	logPoller.MergeFilter([]common.Hash{
		t.configSetTopic,
	}, dkgAddress)

	return &coordinator{
		coordinatorContract: coordinatorContract,
		coordinatorAddress:  coordinatorAddress,
		dkgAddress:          dkgAddress,
		lp:                  logPoller,
		topics:              t,
		lookbackBlocks:      lookbackBlocks,
		evmClient:           client,
		orm:                 orm,
		lggr:                lggr.Named("OCR2VRFCoordinator"),
	}, nil
}

// ReportBlocks returns the heights and hashes of the blocks which require VRF
// proofs in the current report, and the callback requests which should be
// served as part of processing that report. Everything returned by this
// should concern blocks older than the corresponding confirmationDelay.
// Blocks and callbacks it has returned previously may be returned again, as
// long as retransmissionDelay blocks have passed since they were last
// returned. The callbacks returned do not have to correspond to the blocks.
//
// The implementor is responsible for only returning well-funded callback
// requests, and blocks for which clients have actually requested random output
//
// This can be implemented on ethereum using the RandomnessRequested and
// RandomnessFulfillmentRequested events, to identify which blocks and
// callbacks need to be served, along with the NewTransmission and
// RandomWordsFulfilled events, to identify which have already been served.
func (c *coordinator) ReportBlocks(
	ctx context.Context,
	slotInterval uint16, // TODO: unused for now
	confirmationDelays map[uint32]struct{},
	retransmissionDelay time.Duration, // TODO: unused for now
	maxBlocks, // TODO: unused for now
	maxCallbacks int, // TODO: unused for now
) (blocks []ocr2vrftypes.Block, callbacks []ocr2vrftypes.AbstractCostedCallbackRequest, err error) {
	// TODO: use head broadcaster instead?
	currentHead, err := c.evmClient.HeadByNumber(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "header by number")
		return
	}
	currentHeight := currentHead.Number

	logs, err := c.lp.LogsWithSigs(
		currentHeight-c.lookbackBlocks,
		currentHeight-1,
		[]common.Hash{
			c.randomnessRequestedTopic,
			c.randomnessFulfillmentRequestedTopic,
			c.randomWordsFulfilledTopic,
			c.newTransmissionTopic,
		},
		c.coordinatorAddress,
		pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "logs with topics")
		return
	}

	var (
		randomnessRequestedLogs            []vrf_wrapper.VRFBeaconCoordinatorRandomnessRequested
		randomnessFulfillmentRequestedLogs []vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested
		randomWordsFulfilledLogs           []vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled
		newTransmissionLogs                []vrf_wrapper.VRFBeaconCoordinatorNewTransmission
	)
	for _, lg := range logs {
		switch {
		case bytes.Equal(lg.EventSig, c.randomnessRequestedTopic[:]):
			unpacked, err2 := unmarshalRandomnessRequested(lg) // FIXME
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomnessRequested log")
			}
			randomnessRequestedLogs = append(randomnessRequestedLogs, unpacked)
		case bytes.Equal(lg.EventSig, c.randomnessFulfillmentRequestedTopic[:]):
			unpacked, err2 := unmarshalRandomnessFulfillmentRequested(lg) // FIXME
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomnessFulfillmentRequested log")
			}
			randomnessFulfillmentRequestedLogs = append(randomnessFulfillmentRequestedLogs, unpacked)
		case bytes.Equal(lg.EventSig, c.randomWordsFulfilledTopic[:]):
			unpacked, err2 := unmarshalRandomWordsFulfilled(lg) // FIXME
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomWordsFulfilled log")
			}
			randomWordsFulfilledLogs = append(randomWordsFulfilledLogs, unpacked)
		case bytes.Equal(lg.EventSig, c.newTransmissionTopic[:]):
			unpacked, err2 := unmarshalNewTransmission(lg) // FIXME
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal NewTransmission log")
			}
			newTransmissionLogs = append(newTransmissionLogs, unpacked)
		}
	}

	// Scan for blocks where an output is required
	// blocksRequested maps block number to the block object.
	type key struct {
		blockNumber uint64
		confDelay   uint32
	}
	blocksRequested := make(map[key]struct{})
	for _, r := range randomnessRequestedLogs {
		if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
			// if we can't find the conf delay in the map then just ignore this request
			c.lggr.Infow("ignoring bad request, unsupported conf delay",
				"confDelay", r.ConfDelay.String(),
				"supportedConfDelays", confirmationDelays)
			continue
		}
		// If the next beacon output height is less than currentHeight - conf delay
		// AND the log has enough confirmations, then we can schedule it to be fulfilled.
		cond := r.ConfDelay.Int64() < currentHeight // TODO: is this necessary? Won't this always be true?
		cond = cond && r.NextBeaconOutputHeight < uint64(currentHeight-r.ConfDelay.Int64())
		cond = cond && currentHeight >= int64(r.Raw.BlockNumber+r.ConfDelay.Uint64()) // TODO: is this redundant?
		if cond {
			blocksRequested[key{
				blockNumber: r.NextBeaconOutputHeight,
				confDelay:   uint32(r.ConfDelay.Uint64()),
			}] = struct{}{}
		}
	}

	// Scan for blocks where a callback is requested
	var callbacksRequested []vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested
	for _, r := range randomnessFulfillmentRequestedLogs {
		if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
			// if we can't find the conf delay in the map then just ignore this request
			c.lggr.Infow("ignoring bad request, unsupported conf delay",
				"confDelay", r.ConfDelay.String(),
				"supportedConfDelays", confirmationDelays)
			continue
		}

		// If the next beacon output height is less than currentHeight - conf delay
		// AND the log has enough confirmations, then we can schedule it to be fulfilled.
		cond := r.ConfDelay.Int64() < currentHeight // TODO: is this necessary? Won't this always be true?
		cond = cond && r.NextBeaconOutputHeight < uint64(currentHeight-r.ConfDelay.Int64())
		cond = cond && currentHeight >= int64(r.Raw.BlockNumber+r.ConfDelay.Uint64()) // TODO: is this redundant?
		if cond {
			callbacksRequested = append(callbacksRequested, r)

			// We could have a callback request that was made in a different block than what we
			// have possibly already received from regular requests.
			blocksRequested[key{
				blockNumber: r.NextBeaconOutputHeight,
				confDelay:   uint32(r.ConfDelay.Uint64()),
			}] = struct{}{}
		}
	}

	// Prune blocks that have already received responses so that we don't
	// respond to them again.
	for _, r := range newTransmissionLogs {
		for _, o := range r.OutputsServed {
			k := key{
				blockNumber: o.Height,
				confDelay:   uint32(o.ConfirmationDelay.Uint64()),
			}
			if _, ok := blocksRequested[k]; ok {
				delete(blocksRequested, k)
			}
		}
	}

	// Get all the block hashes for the blocks that we need to service from the head saver.
	// Note that we do this to avoid making an RPC call for each block height separately.
	// Alternatively, we could do a batch RPC call.
	var blockHeights []uint64
	for k, _ := range blocksRequested {
		blockHeights = append(blockHeights, k.blockNumber)
	}

	// TODO: is it possible that the head saver doesn't have some of these heights?
	heads, err := c.orm.HeadsByNumbers(ctx, blockHeights)
	if err != nil {
		err = errors.Wrap(err, "heads by numbers")
		return
	}
	if len(heads) != len(blockHeights) {
		err = fmt.Errorf("could not find all heads in db: want %d got %d", len(blockHeights), len(heads))
		return
	}
	for k, _ := range blocksRequested {
		for _, head := range heads {
			if k.blockNumber == uint64(head.Number) {
				blocks = append(blocks, ocr2vrftypes.Block{
					Hash:              head.Hash,
					Height:            k.blockNumber,
					ConfirmationDelay: k.confDelay,
				})
			}
		}
	}

	// Find the callback requests that have been fulfilled on-chain
	fulfilledRequestIDs := make(map[uint64]struct{})
	for _, r := range randomWordsFulfilledLogs {
		for i, requestID := range r.RequestIDs {
			if r.SuccessfulFulfillment[i] == 1 {
				fulfilledRequestIDs[requestID.Uint64()] = struct{}{}
			}
		}
	}

	// Find unfulfilled callback requests
	// TODO: check if subscription has enough funds (eventually)
	for _, r := range callbacksRequested {
		requestID := r.Callback.RequestID
		if _, ok := fulfilledRequestIDs[requestID.Uint64()]; !ok {
			if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
				// if we can't find the conf delay in the map then just ignore this request
				c.lggr.Infow("ignoring bad request, unsupported conf delay",
					"confDelay", r.ConfDelay.String(),
					"supportedConfDelays", confirmationDelays)
				continue
			}
			cond := r.ConfDelay.Int64() < currentHeight // TODO: is this necessary? Won't this always be true?
			cond = cond && r.NextBeaconOutputHeight < uint64(currentHeight-r.ConfDelay.Int64())
			cond = cond && currentHeight >= int64(r.Raw.BlockNumber+r.ConfDelay.Uint64()) // TODO: is this redundant?
			if cond {
				callbacks = append(callbacks, ocr2vrftypes.AbstractCostedCallbackRequest{
					BeaconHeight:      r.NextBeaconOutputHeight,
					ConfirmationDelay: uint32(r.ConfDelay.Uint64()),
					SubscriptionID:    r.SubID,
					Price:             big.NewInt(0), // TODO: no price tracking
					RequestID:         requestID.Uint64(),
					NumWords:          r.Callback.NumWords,
					Requester:         r.Callback.Requester,
					Arguments:         r.Callback.Arguments,
					GasAllowance:      r.Callback.GasAllowance,
					RequestHeight:     r.Raw.BlockNumber,
					RequestBlockHash:  r.Raw.BlockHash,
				})
			}
		}
	}

	return
}

// ReportWillBeTransmitted registers to the CoordinatorInterface that the
// local node has accepted the AbstractReport for transmission, so that its
// blocks and callbacks can be tracked for possible later retransmission
func (c *coordinator) ReportWillBeTransmitted(ctx context.Context, report ocr2vrftypes.AbstractReport) error {
	// TODO: implement me
	// Improve upon the implementation in the future in a more optimized version.
	return nil
}

// DKGVRFCommittees returns the addresses of the signers and transmitters
// for the DKG and VRF OCR committees. On ethereum, these can be retrieved
// from the most recent ConfigSet events for each contract.
func (c *coordinator) DKGVRFCommittees(ctx context.Context) (dkgCommittee, vrfCommittee ocr2vrftypes.OCRCommittee, err error) {
	latestVRF, err := c.lp.LatestLogByEventSigWithConfs(
		c.configSetTopic,
		c.coordinatorAddress,
		1,
	)
	if err != nil {
		err = errors.Wrap(err, "latest vrf ConfigSet by sig with confs")
		return
	}

	latestDKG, err := c.lp.LatestLogByEventSigWithConfs(
		c.configSetTopic,
		c.dkgAddress,
		1,
	)
	if err != nil {
		err = errors.Wrap(err, "latest dkg ConfigSet by sig with confs")
		return
	}

	var vrfConfigSetLog vrf_wrapper.VRFBeaconCoordinatorConfigSet
	err = vrfABI.UnpackIntoInterface(&vrfConfigSetLog, configSetEvent, latestVRF.Data)
	if err != nil {
		err = errors.Wrap(err, "unpack vrf ConfigSet into interface")
		return
	}

	var dkgConfigSetLog dkg_wrapper.DKGConfigSet
	err = dkgABI.UnpackIntoInterface(&dkgConfigSetLog, configSetEvent, latestDKG.Data)
	if err != nil {
		err = errors.Wrap(err, "unpack dkg ConfigSet into interface")
		return
	}

	// NOTE: is it guaranteed that len(signers) == len(transmitters)?
	// in that case, we can simplify the below to a single loop.
	for _, signer := range vrfConfigSetLog.Signers {
		vrfCommittee.Signers = append(vrfCommittee.Signers, signer)
	}
	for _, transmitter := range vrfConfigSetLog.Transmitters {
		vrfCommittee.Transmitters = append(vrfCommittee.Transmitters, transmitter)
	}

	for _, signer := range dkgConfigSetLog.Signers {
		dkgCommittee.Signers = append(dkgCommittee.Signers, signer)
	}
	for _, transmitter := range dkgConfigSetLog.Transmitters {
		dkgCommittee.Transmitters = append(dkgCommittee.Transmitters, transmitter)
	}

	return
}

// ProvingKeyHash returns the VRF current proving key, in view of the local
// node. On ethereum this can be retrieved from the VRF contract's attribute
// s_provingKeyHash
func (c *coordinator) ProvingKeyHash(ctx context.Context) (common.Hash, error) {
	h, err := c.coordinatorContract.SProvingKeyHash(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "get proving key hash")
	}

	return h, nil
}

// BeaconPeriod returns the period used in the coordinator's contract
func (c *coordinator) BeaconPeriod(ctx context.Context) (uint16, error) {
	beaconPeriodBlocks, err := c.coordinatorContract.IBeaconPeriodBlocks(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return 0, errors.Wrap(err, "get beacon period blocks")
	}

	return uint16(beaconPeriodBlocks.Int64()), nil
}
