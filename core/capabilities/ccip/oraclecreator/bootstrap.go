package oraclecreator

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/libocr/networking"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/peergroup"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	ccipcommon "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/common"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/synchronization"
	"github.com/smartcontractkit/chainlink/v2/core/services/telemetry"
)

var _ cctypes.OracleCreator = &bootstrapOracleCreator{}

// bootstrapOracle wraps a CCIPOracle (the bootstrapper) and manages RMN-specific resources
type bootstrapOracle struct {
	baseOracle      cctypes.CCIPOracle
	peerGroupDialer *peerGroupDialer
	rmnHomeReader   ccipreaderpkg.RMNHome
	mu              sync.Mutex
}

func newBootstrapOracle(
	baseOracle cctypes.CCIPOracle,
	peerGroupDialer *peerGroupDialer,
	rmnHomeReader ccipreaderpkg.RMNHome,
) cctypes.CCIPOracle {
	return &bootstrapOracle{
		baseOracle:      baseOracle,
		peerGroupDialer: peerGroupDialer,
		rmnHomeReader:   rmnHomeReader,
	}
}

func (o *bootstrapOracle) Start() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	// Start RMNHome reader first
	if err := o.rmnHomeReader.Start(context.Background()); err != nil {
		return fmt.Errorf("failed to start RMNHome reader: %w", err)
	}

	if err := o.peerGroupDialer.Start(); err != nil {
		// Clean up RMN components if peer group dialer fails to start
		_ = o.rmnHomeReader.Close()
		return fmt.Errorf("failed to start peer group dialer: %w", err)
	}

	// Then start the base oracle (bootstrapper)
	if err := o.baseOracle.Start(); err != nil {
		// Clean up RMN components if base fails to start
		_ = o.rmnHomeReader.Close()
		_ = o.peerGroupDialer.Close()
		return fmt.Errorf("failed to start base oracle: %w", err)
	}

	return nil
}

func (o *bootstrapOracle) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	var errs []error

	if err := o.baseOracle.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close base oracle: %w", err))
	}

	if err := o.peerGroupDialer.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close peer group dialer: %w", err))
	}

	if err := o.rmnHomeReader.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close RMN home reader: %w", err))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

type bootstrapOracleCreator struct {
	peerWrapper             *ocrcommon.SingletonPeerWrapper
	bootstrapperLocators    []commontypes.BootstrapperLocator
	db                      ocr3types.Database
	monitoringEndpointGen   telemetry.MonitoringEndpointGenerator
	lggr                    logger.Logger
	homeChainContractReader types.ContractReader
	addressCodec            ccipcommon.AddressCodec
}

func NewBootstrapOracleCreator(
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	bootstrapperLocators []commontypes.BootstrapperLocator,
	db ocr3types.Database,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	lggr logger.Logger,
	homeChainContractReader types.ContractReader,
	addressCodec ccipcommon.AddressCodec,
) cctypes.OracleCreator {
	return &bootstrapOracleCreator{
		peerWrapper:             peerWrapper,
		bootstrapperLocators:    bootstrapperLocators,
		db:                      db,
		monitoringEndpointGen:   monitoringEndpointGen,
		lggr:                    lggr,
		homeChainContractReader: homeChainContractReader,
		addressCodec:            addressCodec,
	}
}

// Type implements types.OracleCreator.
func (i *bootstrapOracleCreator) Type() cctypes.OracleType {
	return cctypes.OracleTypeBootstrap
}

// Create implements types.OracleCreator.
func (i *bootstrapOracleCreator) Create(ctx context.Context, _ uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	// Assuming that the chain selector is referring to an evm chain for now.
	// TODO: add an api that returns chain family.
	// NOTE: this doesn't really matter for the bootstrap node, it doesn't do anything on-chain.
	// Its for the monitoring endpoint generation below.
	chainID, err := chainsel.ChainIdFromSelector(uint64(config.Config.ChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	destChainFamily := chaintype.EVM
	destRelayID := types.NewRelayID(string(destChainFamily), strconv.FormatUint(chainID, 10))

	oraclePeerIDs := make([]ragep2ptypes.PeerID, 0, len(config.Config.Nodes))
	for _, n := range config.Config.Nodes {
		oraclePeerIDs = append(oraclePeerIDs, n.P2pID)
	}

	rmnHomeReader, err := i.getRmnHomeReader(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get RMNHome reader: %w", err)
	}

	pgd := newPeerGroupDialer(
		logger.Named(i.lggr, "PeerGroupDialer"),
		i.peerWrapper.PeerGroupFactory,
		rmnHomeReader,
		i.bootstrapperLocators,
		oraclePeerIDs,
		config.ConfigDigest,
	)

	bootstrapperArgs := libocr3.BootstrapperArgs{
		BootstrapperFactory:   i.peerWrapper.Peer2,
		V2Bootstrappers:       i.bootstrapperLocators,
		ContractConfigTracker: ocrimpls.NewConfigTracker(config, i.addressCodec),
		Database:              i.db,
		LocalConfig:           defaultLocalConfig(),
		Logger: ocrcommon.NewOCRWrapper(
			logger.Sugared(i.lggr).
				Named("CCIPBootstrap").
				Named(destRelayID.String()).
				Named(config.Config.ChainSelector.String()).
				Named(hexutil.Encode(config.Config.OfframpAddress)),
			false, /* traceLogging */
			func(ctx context.Context, msg string) {}),
		MonitoringEndpoint: i.monitoringEndpointGen.GenMonitoringEndpoint(
			string(destChainFamily),
			destRelayID.ChainID,
			hexutil.Encode(config.Config.OfframpAddress),
			synchronization.OCR3CCIPBootstrap,
		),
		OffchainConfigDigester: ocrimpls.NewConfigDigester(config.ConfigDigest),
	}
	bootstrapper, err := libocr3.NewBootstrapper(bootstrapperArgs)
	if err != nil {
		return nil, err
	}

	bootstrapperWithCustomClose := newWrappedOracle(
		bootstrapper,
		[]io.Closer{pgd, rmnHomeReader},
	)

	return newBootstrapOracle(bootstrapperWithCustomClose, pgd, rmnHomeReader), nil
}

func (i *bootstrapOracleCreator) getRmnHomeReader(ctx context.Context, config cctypes.OCR3ConfigWithMeta) (ccipreaderpkg.RMNHome, error) {
	return ccipreaderpkg.NewRMNHomeChainReader(
		ctx,
		i.lggr,
		ccipreaderpkg.HomeChainPollingInterval,
		config.Config.ChainSelector,
		config.Config.RmnHomeAddress,
		i.homeChainContractReader,
	)
}

// peerGroupDialer keeps watching for RMNHome config changes and calls NewPeerGroup when needed.
// Required for managing RMN related peer group connections.
type peerGroupDialer struct {
	services.StateMachine

	lggr logger.Logger

	peerGroupCreator *peergroup.Creator
	rmnHomeReader    ccipreaderpkg.RMNHome

	// common oracle config
	bootstrapLocators  []commontypes.BootstrapperLocator
	oraclePeerIDs      []ragep2ptypes.PeerID
	commitConfigDigest [32]byte

	// state: accessed in mutually exclusive fashion
	// between Close() and syncLoop().
	// Close shuts down syncLoop() first and then cleans up these.
	activePeerGroups            []networking.PeerGroup
	activeEndpointConfigDigests []cciptypes.Bytes32

	syncInterval time.Duration

	stopChan chan struct{}
	wg       sync.WaitGroup
}

type syncAction struct {
	actionType           actionType
	endpointConfigDigest cciptypes.Bytes32
	rmnHomeConfigDigest  cciptypes.Bytes32
}

type actionType string

const (
	ActionCreate actionType = "create"
	ActionClose  actionType = "close"
)

func newPeerGroupDialer(
	lggr logger.Logger,
	peerGroupFactory networking.PeerGroupFactory,
	rmnHomeReader ccipreaderpkg.RMNHome,
	bootstrapLocators []commontypes.BootstrapperLocator,
	oraclePeerIDs []ragep2ptypes.PeerID,
	commitConfigDigest [32]byte,
) *peerGroupDialer {
	return &peerGroupDialer{
		lggr: lggr,

		peerGroupCreator: peergroup.NewCreator(lggr, peerGroupFactory, bootstrapLocators),
		rmnHomeReader:    rmnHomeReader,

		bootstrapLocators:  bootstrapLocators,
		oraclePeerIDs:      oraclePeerIDs,
		commitConfigDigest: commitConfigDigest,

		activePeerGroups: []networking.PeerGroup{},

		syncInterval: 12 * time.Second, // todo: make it configurable

		stopChan: make(chan struct{}),
		wg:       sync.WaitGroup{},
	}
}

func (d *peerGroupDialer) Start() error {
	return d.StateMachine.StartOnce("peerGroupDialer", func() error {
		d.lggr.Infow("Starting peer group dialer")

		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			d.syncLoop()
		}()

		return nil
	})
}

func (d *peerGroupDialer) syncLoop() {
	// eager sync on start to quickly init.
	d.sync()

	syncTicker := time.NewTicker(d.syncInterval)
	for {
		select {
		case <-syncTicker.C:
			d.sync()
		case <-d.stopChan:
			return
		}
	}
}

func (d *peerGroupDialer) Close() error {
	return d.StateMachine.StopOnce("peerGroupDialer", func() error {
		// shut down the sync goroutine.
		// the order of operations here is important:
		// * we close the stop channel and wait for the syncLoop to stop
		// * only when the sync loop is stopped do we close the peer groups
		// this avoids the race where the sync loop uses peer groups while they are being closed
		close(d.stopChan)
		d.wg.Wait()

		// close all peer groups.
		d.closeExistingPeerGroups()

		return nil
	})
}

// Pure function for calculating sync actions
func calculateSyncActions(
	commitConfigDigest cciptypes.Bytes32,
	currentEndpointConfigDigests []cciptypes.Bytes32,
	activeRmnHomeConfigDigest cciptypes.Bytes32,
	candidateRmnHomeConfigDigest cciptypes.Bytes32,
) []syncAction {
	current := mapset.NewSet[cciptypes.Bytes32]()
	for _, digest := range currentEndpointConfigDigests {
		current.Add(digest)
	}

	endpointDigestToRmnHomeDigest := make(map[cciptypes.Bytes32]cciptypes.Bytes32, 2)

	desired := mapset.NewSet[cciptypes.Bytes32]()
	if !activeRmnHomeConfigDigest.IsEmpty() {
		endpointDigest := writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo,
			sha256.Sum256(append(commitConfigDigest[:], activeRmnHomeConfigDigest[:]...)))
		desired.Add(endpointDigest)
		endpointDigestToRmnHomeDigest[endpointDigest] = activeRmnHomeConfigDigest
	}

	if !candidateRmnHomeConfigDigest.IsEmpty() {
		endpointDigest := writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo,
			sha256.Sum256(append(commitConfigDigest[:], candidateRmnHomeConfigDigest[:]...)))
		desired.Add(endpointDigest)
		endpointDigestToRmnHomeDigest[endpointDigest] = candidateRmnHomeConfigDigest
	}

	closeCount := current.Difference(desired).Cardinality()
	createCount := desired.Difference(current).Cardinality()
	actions := make([]syncAction, 0, closeCount+createCount)

	// Configs to close: in current but not in desired
	for digest := range current.Difference(desired).Iterator().C {
		actions = append(actions, syncAction{
			actionType:           ActionClose,
			rmnHomeConfigDigest:  endpointDigestToRmnHomeDigest[digest],
			endpointConfigDigest: digest,
		})
	}

	// Configs to create: in desired but not in current
	for digest := range desired.Difference(current).Iterator().C {
		actions = append(actions, syncAction{
			actionType:           ActionCreate,
			rmnHomeConfigDigest:  endpointDigestToRmnHomeDigest[digest],
			endpointConfigDigest: digest,
		})
	}

	return actions
}

func (d *peerGroupDialer) sync() {
	activeRmnHomeDigest, candidateRmnHomeDigest := d.rmnHomeReader.GetAllConfigDigests()

	lggr := logger.With(
		d.lggr,
		"method", "sync",
		"activeRmnHomeDigest", activeRmnHomeDigest.String(),
		"candidateRmnHomeDigest", candidateRmnHomeDigest.String(),
		"activeEndpointConfigDigests", fmt.Sprintf("%s", d.activeEndpointConfigDigests),
	)

	actions := calculateSyncActions(
		d.commitConfigDigest, d.activeEndpointConfigDigests, activeRmnHomeDigest, candidateRmnHomeDigest)
	if len(actions) == 0 {
		lggr.Debugw("No peer group actions needed")
		return
	}

	lggr.Infof("Syncing peer groups by applying the actions: %v", actions)

	// Handle each action
	for _, action := range actions {
		actionLggr := logger.With(lggr,
			"action", action.actionType,
			"endpointConfigDigest", action.endpointConfigDigest,
			"rmnHomeConfigDigest", action.rmnHomeConfigDigest)

		switch action.actionType {
		case ActionClose:
			d.closePeerGroup(action.endpointConfigDigest)
			actionLggr.Infow("Peer group closed successfully")
		case ActionCreate:
			if err := d.createPeerGroup(action.endpointConfigDigest, action.rmnHomeConfigDigest); err != nil {
				actionLggr.Errorw("Failed to create peer group", "err", err)
				// Consider closing all groups on error
				d.closeExistingPeerGroups()
				return
			}
		}
	}
}

// Helper function to close specific peer group
func (d *peerGroupDialer) closePeerGroup(endpointConfigDigest cciptypes.Bytes32) {
	lggr := logger.With(d.lggr, "genericEndpointConfigDigest", endpointConfigDigest.String())

	for i, digest := range d.activeEndpointConfigDigests {
		if digest == endpointConfigDigest {
			if err := d.activePeerGroups[i].Close(); err != nil {
				lggr.Warnw("Failed to close peer group", "err", err)
			} else {
				d.lggr.Infow("Closed peer group successfully")
			}

			// Remove from active groups and digests
			d.activePeerGroups = append(d.activePeerGroups[:i], d.activePeerGroups[i+1:]...)
			d.activeEndpointConfigDigests = append(
				d.activeEndpointConfigDigests[:i], d.activeEndpointConfigDigests[i+1:]...)
			return
		}
	}
}

func (d *peerGroupDialer) createPeerGroup(
	endpointConfigDigest cciptypes.Bytes32,
	rmnHomeConfigDigest cciptypes.Bytes32,
) error {
	rmnNodesInfo, err := d.rmnHomeReader.GetRMNNodesInfo(rmnHomeConfigDigest)
	if err != nil {
		return fmt.Errorf("get RMN nodes info: %w", err)
	}

	lggr := logger.With(d.lggr,
		"genericEndpointConfigDigest", endpointConfigDigest.String(),
		"rmnHomeConfigDigest", rmnHomeConfigDigest.String(),
	)

	lggr.Infow("Creating new peer group")
	peerGroup, err := d.peerGroupCreator.Create(peergroup.CreateOpts{
		CommitConfigDigest:  d.commitConfigDigest,
		RMNHomeConfigDigest: rmnHomeConfigDigest,
		OraclePeerIDs:       d.oraclePeerIDs,
		RMNNodes:            rmnNodesInfo,
	})
	if err != nil {
		return fmt.Errorf("new peer group: %w", err)
	}

	lggr.Infow("Created new peer group successfully")

	d.activePeerGroups = append(d.activePeerGroups, peerGroup.PeerGroup)
	d.activeEndpointConfigDigests = append(d.activeEndpointConfigDigests, endpointConfigDigest)

	return nil
}

// closeExistingPeerGroups closes all active peer groups
func (d *peerGroupDialer) closeExistingPeerGroups() {
	for _, pg := range d.activePeerGroups {
		if err := pg.Close(); err != nil {
			d.lggr.Warnw("failed to close peer group", "err", err)
			continue
		}
		d.lggr.Infow("Closed peer group successfully")
	}

	d.activePeerGroups = []networking.PeerGroup{}
	d.activeEndpointConfigDigests = []cciptypes.Bytes32{}
}

func writePrefix(prefix ocr2types.ConfigDigestPrefix, hash cciptypes.Bytes32) cciptypes.Bytes32 {
	var prefixBytes [2]byte
	binary.BigEndian.PutUint16(prefixBytes[:], uint16(prefix))

	hCopy := hash
	hCopy[0] = prefixBytes[0]
	hCopy[1] = prefixBytes[1]

	return hCopy
}
