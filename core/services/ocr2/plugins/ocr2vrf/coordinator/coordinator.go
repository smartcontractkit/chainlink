package coordinator

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

	vrf_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/ocr2vrf/generated/vrf_beacon_coordinator"

	dkg_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/ocr2vrf/generated/dkg"

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

// block is used to key into a set that tracks beacon blocks.
type block struct {
	blockNumber uint64
	confDelay   uint32
}

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

	// set of blocks that have been scheduled for transmission.
	toBeTransmittedBlocks map[block]struct{}
	// set of request id's that have been scheduled for transmission.
	toBeTransmittedCallbacks map[uint64]struct{}
	// transmittedMu protects the toBeTransmittedBlocks and toBeTransmittedCallbacks
	transmittedMu sync.Mutex
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
	// Call MergeFilter once for each event signature, otherwise the log poller won't
	// index the logs we want.
	for _, sig := range []common.Hash{
		t.randomnessRequestedTopic,
		t.randomnessFulfillmentRequestedTopic,
		t.randomWordsFulfilledTopic,
		t.configSetTopic,
		t.newTransmissionTopic,
	} {
		logPoller.MergeFilter([]common.Hash{sig}, coordinatorAddress)
	}

	// We need ConfigSet events from the DKG contract as well.
	logPoller.MergeFilter([]common.Hash{
		t.configSetTopic,
	}, dkgAddress)

	return &coordinator{
		coordinatorContract:      coordinatorContract,
		coordinatorAddress:       coordinatorAddress,
		dkgAddress:               dkgAddress,
		lp:                       logPoller,
		topics:                   t,
		lookbackBlocks:           lookbackBlocks,
		evmClient:                client,
		orm:                      orm,
		lggr:                     lggr.Named("OCR2VRFCoordinator"),
		toBeTransmittedBlocks:    make(map[block]struct{}),
		toBeTransmittedCallbacks: make(map[uint64]struct{}),
		transmittedMu:            sync.Mutex{},
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

	c.lggr.Infow("current chain height", "currentHeight", currentHeight)

	logs, err := c.lp.LogsWithSigs(
		currentHeight-c.lookbackBlocks,
		currentHeight,
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

	c.lggr.Info(fmt.Sprintf("vrf LogsWithSigs: %+v", logs))

	randomnessRequestedLogs,
		randomnessFulfillmentRequestedLogs,
		randomWordsFulfilledLogs,
		newTransmissionLogs,
		err := c.unmarshalLogs(logs)
	if err != nil {
		err = errors.Wrap(err, "unmarshal logs")
		return
	}

	c.lggr.Info(fmt.Sprintf("finished unmarshalLogs: RandomnessRequested: %+v , RandomnessFulfillmentRequested: %+v , RandomWordsFulfilled: %+v , NewTransmission: %+v",
		randomnessRequestedLogs, randomWordsFulfilledLogs, newTransmissionLogs, randomnessFulfillmentRequestedLogs))

	blocksRequested := make(map[block]struct{})
	unfulfilled := c.filterEligibleRandomnessRequests(randomnessRequestedLogs, confirmationDelays, currentHeight)
	for _, uf := range unfulfilled {
		blocksRequested[uf] = struct{}{}
	}

	c.lggr.Info(fmt.Sprintf("filtered eligible randomness requests: %+v", unfulfilled))

	callbacksRequested, unfulfilled := c.filterEligibleCallbacks(randomnessFulfillmentRequestedLogs, confirmationDelays, currentHeight)
	for _, uf := range unfulfilled {
		blocksRequested[uf] = struct{}{}
	}

	c.lggr.Info(fmt.Sprintf("filtered eligible callbacks: %+v, unfulfilled: %+v", callbacksRequested, unfulfilled))

	// Remove blocks that have already received responses so that we don't
	// respond to them again.
	fulfilledBlocks := c.getFulfilledBlocks(newTransmissionLogs)
	for _, f := range fulfilledBlocks {
		delete(blocksRequested, f)
	}

	c.lggr.Info(fmt.Sprintf("got fulfilled blocks: %+v", fulfilledBlocks))

	// Construct the slice of blocks to return. At this point
	// we only need to fetch the blockhashes of the blocks that
	// need a VRF output.
	blocks, err = c.getBlocks(ctx, blocksRequested)
	if err != nil {
		return
	}

	c.lggr.Info(fmt.Sprintf("got blocks: %+v", blocks))

	// Find unfulfilled callback requests by filtering out already fulfilled callbacks.
	fulfilledRequestIDs := c.getFulfilledRequestIDs(randomWordsFulfilledLogs)
	callbacks = c.filterUnfulfilledCallbacks(callbacksRequested, fulfilledRequestIDs, confirmationDelays, currentHeight)

	c.lggr.Info(fmt.Sprintf("filtered unfulfilled callbacks: %+v, fulfilled: %+v", callbacks, fulfilledRequestIDs))

	return
}

func (c *coordinator) getFulfilledBlocks(newTransmissionLogs []vrf_wrapper.VRFBeaconCoordinatorNewTransmission) (fulfilled []block) {
	for _, r := range newTransmissionLogs {
		for _, o := range r.OutputsServed {
			fulfilled = append(fulfilled, block{
				blockNumber: o.Height,
				confDelay:   uint32(o.ConfirmationDelay.Uint64()),
			})
		}
	}
	return
}

// getBlocks returns the blocks that require a VRF output.
func (c *coordinator) getBlocks(
	ctx context.Context,
	blocksRequested map[block]struct{},
) (blocks []ocr2vrftypes.Block, err error) {
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

	headSet := make(map[uint64]*evmtypes.Head)
	for _, h := range heads {
		headSet[uint64(h.Number)] = h
	}

	for k, _ := range blocksRequested {
		if head, ok := headSet[k.blockNumber]; ok {
			blocks = append(blocks, ocr2vrftypes.Block{
				Hash:              head.Hash,
				Height:            k.blockNumber,
				ConfirmationDelay: k.confDelay,
			})
		}
	}
	return
}

// getFulfilledRequestIDs returns the request IDs referenced by the given RandomWordsFulfilled logs slice.
func (c *coordinator) getFulfilledRequestIDs(randomWordsFulfilledLogs []vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled) map[uint64]struct{} {
	fulfilledRequestIDs := make(map[uint64]struct{})
	for _, r := range randomWordsFulfilledLogs {
		for i, requestID := range r.RequestIDs {
			if r.SuccessfulFulfillment[i] == 1 {
				fulfilledRequestIDs[requestID.Uint64()] = struct{}{}
			}
		}
	}
	return fulfilledRequestIDs
}

// filterUnfulfilledCallbacks returns unfulfilled callback requests given the
// callback request logs and the already fulfilled callback request IDs.
func (c *coordinator) filterUnfulfilledCallbacks(
	callbacksRequested []vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested,
	fulfilledRequestIDs map[uint64]struct{},
	confirmationDelays map[uint32]struct{},
	currentHeight int64,
) (callbacks []ocr2vrftypes.AbstractCostedCallbackRequest) {
	// TODO: check if subscription has enough funds (eventually)
	for _, r := range callbacksRequested {
		requestID := r.Callback.RequestID
		if _, ok := fulfilledRequestIDs[requestID.Uint64()]; !ok {
			// The on-chain machinery will revert requests that specify an unsupported
			// confirmation delay, so this is more of a sanity check than anything else.
			if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
				// if we can't find the conf delay in the map then just ignore this request
				c.lggr.Errorw("ignoring bad request, unsupported conf delay",
					"confDelay", r.ConfDelay.String(),
					"supportedConfDelays", confirmationDelays)
				continue
			}

			// NOTE: we already check if the callback has been fulfilled in filterEligibleCallbacks,
			// so we don't need to do that again here.
			if isBlockEligible(r.NextBeaconOutputHeight, r.ConfDelay, currentHeight) {
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
	return callbacks
}

// filterEligibleCallbacks extracts valid callback requests from the given logs,
// based on their readiness to be fulfilled. It also returns any unfulfilled blocks
// associated with those callbacks.
func (c *coordinator) filterEligibleCallbacks(
	randomnessFulfillmentRequestedLogs []vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested,
	confirmationDelays map[uint32]struct{},
	currentHeight int64,
) (callbacks []vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested, unfulfilled []block) {
	c.transmittedMu.Lock()
	defer c.transmittedMu.Unlock()

	for _, r := range randomnessFulfillmentRequestedLogs {
		// The on-chain machinery will revert requests that specify an unsupported
		// confirmation delay, so this is more of a sanity check than anything else.
		if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
			// if we can't find the conf delay in the map then just ignore this request
			c.lggr.Errorw("ignoring bad request, unsupported conf delay",
				"confDelay", r.ConfDelay.String(),
				"supportedConfDelays", confirmationDelays)
			continue
		}

		// check that the callback hasn't been scheduled for transmission
		// so we don't fulfill the callback twice.
		_, transmitted := c.toBeTransmittedCallbacks[r.Callback.RequestID.Uint64()]
		if isBlockEligible(r.NextBeaconOutputHeight, r.ConfDelay, currentHeight) && !transmitted {
			callbacks = append(callbacks, r)

			// We could have a callback request that was made in a different block than what we
			// have possibly already received from regular requests.
			unfulfilled = append(unfulfilled, block{
				blockNumber: r.NextBeaconOutputHeight,
				confDelay:   uint32(r.ConfDelay.Uint64()),
			})
		}
	}
	return
}

// filterEligibleRandomnessRequests extracts valid randomness requests from the given logs,
// based on their readiness to be fulfilled.
func (c *coordinator) filterEligibleRandomnessRequests(
	randomnessRequestedLogs []vrf_wrapper.VRFBeaconCoordinatorRandomnessRequested,
	confirmationDelays map[uint32]struct{},
	currentHeight int64,
) (unfulfilled []block) {
	c.transmittedMu.Lock()
	defer c.transmittedMu.Unlock()

	for _, r := range randomnessRequestedLogs {
		// The on-chain machinery will revert requests that specify an unsupported
		// confirmation delay, so this is more of a sanity check than anything else.
		if _, ok := confirmationDelays[uint32(r.ConfDelay.Uint64())]; !ok {
			// if we can't find the conf delay in the map then just ignore this request
			c.lggr.Errorw("ignoring bad request, unsupported conf delay",
				"confDelay", r.ConfDelay.String(),
				"supportedConfDelays", confirmationDelays)
			continue
		}

		// check if the block has been scheduled for transmission so that we don't
		// retransmit for the same block.
		_, blockTransmitted := c.toBeTransmittedBlocks[block{
			blockNumber: r.NextBeaconOutputHeight,
			confDelay:   uint32(r.ConfDelay.Uint64()),
		}]
		if isBlockEligible(r.NextBeaconOutputHeight, r.ConfDelay, currentHeight) && !blockTransmitted {
			unfulfilled = append(unfulfilled, block{
				blockNumber: r.NextBeaconOutputHeight,
				confDelay:   uint32(r.ConfDelay.Uint64()),
			})
		}
	}
	return
}

func (c *coordinator) unmarshalLogs(
	logs []logpoller.Log,
) (
	randomnessRequestedLogs []vrf_wrapper.VRFBeaconCoordinatorRandomnessRequested,
	randomnessFulfillmentRequestedLogs []vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested,
	randomWordsFulfilledLogs []vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled,
	newTransmissionLogs []vrf_wrapper.VRFBeaconCoordinatorNewTransmission,
	err error,
) {
	for _, lg := range logs {
		switch {
		case bytes.Equal(lg.EventSig, c.randomnessRequestedTopic[:]):
			unpacked, err2 := unmarshalRandomnessRequested(lg)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomnessRequested log")
				return
			}
			randomnessRequestedLogs = append(randomnessRequestedLogs, unpacked)
		case bytes.Equal(lg.EventSig, c.randomnessFulfillmentRequestedTopic[:]):
			unpacked, err2 := unmarshalRandomnessFulfillmentRequested(lg)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomnessFulfillmentRequested log")
				return
			}
			randomnessFulfillmentRequestedLogs = append(randomnessFulfillmentRequestedLogs, unpacked)
		case bytes.Equal(lg.EventSig, c.randomWordsFulfilledTopic[:]):
			unpacked, err2 := unmarshalRandomWordsFulfilled(lg)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal RandomWordsFulfilled log")
				return
			}
			randomWordsFulfilledLogs = append(randomWordsFulfilledLogs, unpacked)
		case bytes.Equal(lg.EventSig, c.newTransmissionTopic[:]):
			unpacked, err2 := unmarshalNewTransmission(lg)
			if err2 != nil {
				// should never happen
				err = errors.Wrap(err2, "unmarshal NewTransmission log")
				return
			}
			newTransmissionLogs = append(newTransmissionLogs, unpacked)
		default:
			c.lggr.Error(fmt.Sprintf("Unexpected event sig: %s", hexutil.Encode(lg.EventSig)))
			c.lggr.Error(fmt.Sprintf("expected one of: %s %s %s %s",
				hexutil.Encode(c.randomnessRequestedTopic[:]), hexutil.Encode(c.randomnessFulfillmentRequestedTopic[:]),
				hexutil.Encode(c.randomWordsFulfilledTopic[:]), hexutil.Encode(c.newTransmissionTopic[:])))
		}
	}
	return
}

// ReportWillBeTransmitted registers to the CoordinatorInterface that the
// local node has accepted the AbstractReport for transmission, so that its
// blocks and callbacks can be tracked for possible later retransmission
func (c *coordinator) ReportWillBeTransmitted(ctx context.Context, report ocr2vrftypes.AbstractReport) error {
	c.transmittedMu.Lock()
	defer c.transmittedMu.Unlock()
	for _, output := range report.Outputs {
		c.toBeTransmittedBlocks[block{
			blockNumber: output.BlockHeight,
			confDelay:   output.ConfirmationDelay,
		}] = struct{}{}
		for _, cb := range output.Callbacks {
			c.toBeTransmittedCallbacks[cb.RequestID] = struct{}{}
		}
	}
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

	// len(signers) == len(transmitters), this is guaranteed by libocr.
	for i := range vrfConfigSetLog.Signers {
		vrfCommittee.Signers = append(vrfCommittee.Signers, vrfConfigSetLog.Signers[i])
		vrfCommittee.Transmitters = append(vrfCommittee.Transmitters, vrfConfigSetLog.Transmitters[i])
	}

	for i := range dkgConfigSetLog.Signers {
		dkgCommittee.Signers = append(dkgCommittee.Signers, dkgConfigSetLog.Signers[i])
		dkgCommittee.Transmitters = append(dkgCommittee.Transmitters, dkgConfigSetLog.Transmitters[i])
	}

	return
}

// ProvingKeyHash returns the VRF current proving block, in view of the local
// node. On ethereum this can be retrieved from the VRF contract's attribute
// s_provingKeyHash
func (c *coordinator) ProvingKeyHash(ctx context.Context) (common.Hash, error) {
	h, err := c.coordinatorContract.SProvingKeyHash(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "get proving block hash")
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

// isBlockEligible returns true if and only if the nextBeaconOutputHeight plus
// the confDelay is less than the current blockchain height, meaning that the beacon
// output height has enough confirmations.
//
// NextBeaconOutputHeight is always greater than the request block, therefore
// a number of confirmations on the beacon block is always enough confirmations
// for the request block.
func isBlockEligible(nextBeaconOutputHeight uint64, confDelay *big.Int, currentHeight int64) bool {
	cond := confDelay.Int64() < currentHeight // Edge case: for simulated chains with low block numbers
	cond = cond && (nextBeaconOutputHeight+confDelay.Uint64()) < uint64(currentHeight)
	return cond
}
