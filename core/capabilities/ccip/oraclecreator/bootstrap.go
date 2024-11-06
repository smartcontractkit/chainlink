package oraclecreator

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/libocr/networking"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-ccip/pkg/peergroup"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ocrimpls"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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

	o.peerGroupDialer.Start()

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
}

func NewBootstrapOracleCreator(
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	bootstrapperLocators []commontypes.BootstrapperLocator,
	db ocr3types.Database,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	lggr logger.Logger,
	homeChainContractReader types.ContractReader,
) cctypes.OracleCreator {
	return &bootstrapOracleCreator{
		peerWrapper:             peerWrapper,
		bootstrapperLocators:    bootstrapperLocators,
		db:                      db,
		monitoringEndpointGen:   monitoringEndpointGen,
		lggr:                    lggr,
		homeChainContractReader: homeChainContractReader,
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
	destRelayID := types.NewRelayID(string(destChainFamily), fmt.Sprintf("%d", chainID))

	oraclePeerIDs := make([]ragep2ptypes.PeerID, 0, len(config.Config.Nodes))
	for _, n := range config.Config.Nodes {
		oraclePeerIDs = append(oraclePeerIDs, n.P2pID)
	}

	rmnHomeReader, err := i.getRmnHomeReader(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get RMNHome reader: %w", err)
	}

	pgd := newPeerGroupDialer(
		i.lggr.Named("PeerGroupDialer"),
		i.peerWrapper.PeerGroupFactory,
		rmnHomeReader,
		i.bootstrapperLocators,
		oraclePeerIDs,
		config.ConfigDigest,
	)

	bootstrapperArgs := libocr3.BootstrapperArgs{
		BootstrapperFactory:   i.peerWrapper.Peer2,
		V2Bootstrappers:       i.bootstrapperLocators,
		ContractConfigTracker: ocrimpls.NewConfigTracker(config),
		Database:              i.db,
		LocalConfig:           defaultLocalConfig(),
		Logger: ocrcommon.NewOCRWrapper(
			i.lggr.
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
	rmnHomeBoundContract := types.BoundContract{
		Address: "0x" + hex.EncodeToString(config.Config.RmnHomeAddress),
		Name:    consts.ContractNameRMNHome,
	}

	if err1 := i.homeChainContractReader.Bind(ctx, []types.BoundContract{rmnHomeBoundContract}); err1 != nil {
		return nil, fmt.Errorf("failed to bind RMNHome contract: %w", err1)
	}
	rmnHomeReader := ccipreaderpkg.NewRMNHomePoller(
		i.homeChainContractReader,
		rmnHomeBoundContract,
		i.lggr,
		5*time.Second,
	)
	return rmnHomeReader, nil
}

// peerGroupDialer keeps watching for RMNHome config changes and calls NewPeerGroup when needed.
// Required for managing RMN related peer group connections.
type peerGroupDialer struct {
	lggr logger.Logger

	peerGroupCreator *peergroup.Creator
	rmnHomeReader    ccipreaderpkg.RMNHome

	// common oracle config
	bootstrapLocators  []commontypes.BootstrapperLocator
	oraclePeerIDs      []ragep2ptypes.PeerID
	commitConfigDigest [32]byte

	activePeerGroups    []networking.PeerGroup
	activeConfigDigests []cciptypes.Bytes32

	syncInterval time.Duration

	mu         *sync.Mutex
	syncCancel context.CancelFunc
}

type syncAction struct {
	actionType   actionType
	configDigest cciptypes.Bytes32
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

		mu:         &sync.Mutex{},
		syncCancel: nil,
	}
}

func (d *peerGroupDialer) Start() {
	if d.syncCancel != nil {
		d.lggr.Warnw("peer group dialer already started, should not be called twice")
		return
	}

	d.lggr.Infow("Starting peer group dialer")

	ctx, cf := context.WithCancel(context.Background())
	d.syncCancel = cf

	go func() {
		d.sync()

		syncTicker := time.NewTicker(d.syncInterval)
		for {
			select {
			case <-syncTicker.C:
				d.sync()
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (d *peerGroupDialer) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.closeExistingPeerGroups()

	if d.syncCancel != nil {
		d.syncCancel()
	}

	return nil
}

// Pure function for calculating sync actions
func calculateSyncActions(
	currentConfigDigests []cciptypes.Bytes32,
	activeConfigDigest cciptypes.Bytes32,
	candidateConfigDigest cciptypes.Bytes32,
) []syncAction {
	current := mapset.NewSet[cciptypes.Bytes32]()
	for _, digest := range currentConfigDigests {
		current.Add(digest)
	}

	desired := mapset.NewSet[cciptypes.Bytes32]()
	if !activeConfigDigest.IsEmpty() {
		desired.Add(activeConfigDigest)
	}
	if !candidateConfigDigest.IsEmpty() {
		desired.Add(candidateConfigDigest)
	}

	closeCount := current.Difference(desired).Cardinality()
	createCount := desired.Difference(current).Cardinality()
	actions := make([]syncAction, 0, closeCount+createCount)

	// Configs to close: in current but not in desired
	for digest := range current.Difference(desired).Iterator().C {
		actions = append(actions, syncAction{
			actionType:   ActionClose,
			configDigest: digest,
		})
	}

	// Configs to create: in desired but not in current
	for digest := range desired.Difference(current).Iterator().C {
		actions = append(actions, syncAction{
			actionType:   ActionCreate,
			configDigest: digest,
		})
	}

	return actions
}

func (d *peerGroupDialer) sync() {
	d.mu.Lock()
	defer d.mu.Unlock()

	activeDigest, candidateDigest := d.rmnHomeReader.GetAllConfigDigests()
	actions := calculateSyncActions(d.activeConfigDigests, activeDigest, candidateDigest)
	if len(actions) == 0 {
		d.lggr.Debugw("No peer group actions needed")
		return
	}

	d.lggr.Infow("Syncing peer groups", "actions", actions)

	// Handle each action
	for _, action := range actions {
		switch action.actionType {
		case ActionClose:
			d.closePeerGroup(action.configDigest)
		case ActionCreate:
			if err := d.createPeerGroup(action.configDigest); err != nil {
				d.lggr.Errorw("Failed to create peer group",
					"configDigest", action.configDigest,
					"err", err)
				// Consider closing all groups on error
				d.closeExistingPeerGroups()
				return
			}
		}
	}
}

// Helper function to close specific peer group
func (d *peerGroupDialer) closePeerGroup(configDigest cciptypes.Bytes32) {
	for i, digest := range d.activeConfigDigests {
		if digest == configDigest {
			if err := d.activePeerGroups[i].Close(); err != nil {
				d.lggr.Warnw("Failed to close peer group",
					"configDigest", configDigest,
					"err", err)
			} else {
				d.lggr.Infow("Closed peer group successfully",
					"configDigest", configDigest)
			}
			// Remove from active groups and digests
			d.activePeerGroups = append(d.activePeerGroups[:i], d.activePeerGroups[i+1:]...)
			d.activeConfigDigests = append(d.activeConfigDigests[:i], d.activeConfigDigests[i+1:]...)
			return
		}
	}
}

func (d *peerGroupDialer) createPeerGroup(rmnHomeConfigDigest cciptypes.Bytes32) error {
	rmnNodesInfo, err := d.rmnHomeReader.GetRMNNodesInfo(rmnHomeConfigDigest)
	if err != nil {
		return fmt.Errorf("get RMN nodes info: %w", err)
	}

	// Create generic endpoint config digest by hashing commit config digest and rmn home config digest
	h := sha256.Sum256(append(d.commitConfigDigest[:], rmnHomeConfigDigest[:]...))
	genericEndpointConfigDigest := writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo, h)

	// Combine oracle peer IDs with RMN node peer IDs
	peerIDs := make([]string, 0, len(d.oraclePeerIDs)+len(rmnNodesInfo))
	for _, p := range d.oraclePeerIDs {
		peerIDs = append(peerIDs, p.String())
	}
	for _, n := range rmnNodesInfo {
		peerIDs = append(peerIDs, n.PeerID.String())
	}

	lggr := d.lggr.With(
		"genericEndpointConfigDigest", genericEndpointConfigDigest.String(),
		"peerIDs", peerIDs,
		"bootstrappers", d.bootstrapLocators,
	)

	lggr.Infow("Creating new peer group")
	peerGroup, err := d.peerGroupCreator.Create(peergroup.CreateOpts{
		CommitConfigDigest:  cciptypes.Bytes32(d.commitConfigDigest),
		RMNHomeConfigDigest: rmnHomeConfigDigest,
		OraclePeerIDs:       d.oraclePeerIDs,
		RMNNodes:            rmnNodesInfo,
	})
	if err != nil {
		return fmt.Errorf("new peer group: %w", err)
	}
	lggr.Infow("Created new peer group successfully")

	d.activePeerGroups = append(d.activePeerGroups, peerGroup.PeerGroup)
	d.activeConfigDigests = append(d.activeConfigDigests, genericEndpointConfigDigest)

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
	d.activeConfigDigests = []cciptypes.Bytes32{}
}

func writePrefix(prefix ocr2types.ConfigDigestPrefix, hash cciptypes.Bytes32) cciptypes.Bytes32 {
	var prefixBytes [2]byte
	binary.BigEndian.PutUint16(prefixBytes[:], uint16(prefix))

	hCopy := hash
	hCopy[0] = prefixBytes[0]
	hCopy[1] = prefixBytes[1]

	return hCopy
}
