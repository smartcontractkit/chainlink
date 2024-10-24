package oraclecreator

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	libocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"

	"github.com/smartcontractkit/chainlink-ccip/commit/merkleroot/rmn"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
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

type bootstrapOracleCreator struct {
	peerWrapper           *ocrcommon.SingletonPeerWrapper
	bootstrapperLocators  []commontypes.BootstrapperLocator
	db                    ocr3types.Database
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator
	lggr                  logger.Logger
	contractReader        types.ContractReader
}

func NewBootstrapOracleCreator(
	peerWrapper *ocrcommon.SingletonPeerWrapper,
	bootstrapperLocators []commontypes.BootstrapperLocator,
	db ocr3types.Database,
	monitoringEndpointGen telemetry.MonitoringEndpointGenerator,
	lggr logger.Logger,
	contractReader types.ContractReader,
) cctypes.OracleCreator {
	return &bootstrapOracleCreator{
		peerWrapper:           peerWrapper,
		bootstrapperLocators:  bootstrapperLocators,
		db:                    db,
		monitoringEndpointGen: monitoringEndpointGen,
		lggr:                  lggr,
		contractReader:        contractReader,
	}
}

// Type implements types.OracleCreator.
func (i *bootstrapOracleCreator) Type() cctypes.OracleType {
	return cctypes.OracleTypeBootstrap
}

// Create implements types.OracleCreator.
func (i *bootstrapOracleCreator) Create(_ uint32, config cctypes.OCR3ConfigWithMeta) (cctypes.CCIPOracle, error) {
	// Assuming that the chain selector is referring to an evm chain for now.
	// TODO: add an api that returns chain family.
	// NOTE: this doesn't really matter for the bootstrap node, it doesn't do anything on-chain.
	// Its for the monitoring endpoint generation below.
	chainID, err := chainsel.ChainIdFromSelector(uint64(config.Config.ChainSelector))
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	ctx := context.Background()
	rmnHomeReader, err := i.getRmnHomeReader(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to get RMNHome reader: %w", err)
	}
	if err = rmnHomeReader.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start RMNHome reader: %w", err)
	}

	destChainFamily := chaintype.EVM
	destRelayID := types.NewRelayID(string(destChainFamily), fmt.Sprintf("%d", chainID))

	oraclePeerIDs := make([]ragep2ptypes.PeerID, 0, len(config.Config.Nodes))
	for _, n := range config.Config.Nodes {
		oraclePeerIDs = append(oraclePeerIDs, n.P2pID)
	}

	pgd := newPeerGroupDialer(
		i.lggr.Named("PeerGroupDialer"),
		i.peerWrapper.PeerGroupFactory,
		rmnHomeReader,
		i.bootstrapperLocators,
		oraclePeerIDs,
		config.ConfigDigest,
	)
	pgd.Start()

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

	return bootstrapperWithCustomClose, nil
}

func (i *bootstrapOracleCreator) getRmnHomeReader(ctx context.Context, config cctypes.OCR3ConfigWithMeta) (ccipreaderpkg.RMNHome, error) {
	rmnHomeBoundContract := types.BoundContract{
		Address: "0x" + hex.EncodeToString(config.Config.RmnHomeAddress),
		Name:    consts.ContractNameRMNHome,
	}

	if err1 := i.contractReader.Bind(ctx, []types.BoundContract{rmnHomeBoundContract}); err1 != nil {
		return nil, fmt.Errorf("failed to bind RMNHome contract: %w", err1)
	}
	rmnHomeReader := ccipreaderpkg.NewRMNHomePoller(
		i.contractReader,
		rmnHomeBoundContract,
		i.lggr,
		5*time.Second,
	)
	return rmnHomeReader, nil
}

// peerGroupDialer keeps watching for config changes and calls NewPeerGroup when needed.
// Required for managing RMN related peer group connections.
type peerGroupDialer struct {
	lggr logger.Logger

	peerGroupFactory rmn.PeerGroupFactory
	rmnHomeReader    ccipreaderpkg.RMNHome

	// common oracle config
	bootstrapLocators  []commontypes.BootstrapperLocator
	oraclePeerIDs      []ragep2ptypes.PeerID
	commitConfigDigest [32]byte

	activePeerGroups    []rmn.PeerGroup
	activeConfigDigests []cciptypes.Bytes32

	syncInterval time.Duration

	mu        *sync.Mutex
	syncCtxCf context.CancelFunc
}

func newPeerGroupDialer(
	lggr logger.Logger,
	peerGroupFactory rmn.PeerGroupFactory,
	rmnHomeReader ccipreaderpkg.RMNHome,
	bootstrapLocators []commontypes.BootstrapperLocator,
	oraclePeerIDs []ragep2ptypes.PeerID,
	commitConfigDigest [32]byte,
) *peerGroupDialer {
	return &peerGroupDialer{
		lggr: lggr,

		peerGroupFactory: peerGroupFactory,
		rmnHomeReader:    rmnHomeReader,

		bootstrapLocators:  bootstrapLocators,
		oraclePeerIDs:      oraclePeerIDs,
		commitConfigDigest: commitConfigDigest,

		activePeerGroups: []rmn.PeerGroup{},

		syncInterval: time.Minute, // todo: make it configurable

		mu:        &sync.Mutex{},
		syncCtxCf: nil,
	}
}

func (d *peerGroupDialer) Start() {
	if d.syncCtxCf != nil {
		d.lggr.Warnw("peer group dialer already started, should not be called twice")
		return
	}

	d.lggr.Infow("Starting peer group dialer")

	ctx, cf := context.WithCancel(context.Background())
	d.syncCtxCf = cf

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

	if d.syncCtxCf != nil {
		d.syncCtxCf()
	}

	return nil
}

func (d *peerGroupDialer) sync() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.shouldSync() {
		d.lggr.Debugw("No need to sync peer groups")
		return
	}
	d.lggr.Infow("Syncing peer groups")

	d.closeExistingPeerGroups()

	if err := d.createNewPeerGroups(); err != nil {
		d.lggr.Errorw("failed to create new peer groups", "err", err)
		d.closeExistingPeerGroups() // close potentially opened peer groups
	}
}

func (d *peerGroupDialer) shouldSync() bool {
	if len(d.activePeerGroups) == 0 {
		return true
	}

	activeConfigDigest, candidateConfigDigest := d.rmnHomeReader.GetAllConfigDigests()
	var configDigests [][32]byte

	if !activeConfigDigest.IsEmpty() {
		configDigests = append(configDigests, activeConfigDigest)
	}
	if !candidateConfigDigest.IsEmpty() {
		configDigests = append(configDigests, candidateConfigDigest)
	}

	if len(configDigests) != len(d.activeConfigDigests) {
		return true
	}
	for i, rmnHomeConfigDigest := range configDigests {
		if rmnHomeConfigDigest != d.activeConfigDigests[i] {
			return true
		}
	}

	return false
}

func (d *peerGroupDialer) closeExistingPeerGroups() {
	for _, pg := range d.activePeerGroups {
		if err := pg.Close(); err != nil {
			d.lggr.Warnw("failed to close peer group", "err", err)
			continue
		}
		d.lggr.Infow("Closed peer group successfully")
	}

	d.activePeerGroups = []rmn.PeerGroup{}
	d.activeConfigDigests = []cciptypes.Bytes32{}
}

func (d *peerGroupDialer) createNewPeerGroups() error {
	activeConfigDigest, candidateConfigDigest := d.rmnHomeReader.GetAllConfigDigests()
	var configDigests [][32]byte

	if !activeConfigDigest.IsEmpty() {
		configDigests = append(configDigests, activeConfigDigest)
	}
	if !candidateConfigDigest.IsEmpty() {
		configDigests = append(configDigests, candidateConfigDigest)
	}

	d.lggr.Infow("Creating new peer groups", "configDigests", configDigests)

	for _, rmnHomeConfigDigest := range configDigests {
		rmnNodesInfo, err := d.rmnHomeReader.GetRMNNodesInfo(rmnHomeConfigDigest)
		if err != nil {
			return fmt.Errorf("get RMN nodes info: %w", err)
		}

		h := sha256.Sum256(append(d.commitConfigDigest[:], rmnHomeConfigDigest[:]...))
		genericEndpointConfigDigest := writePrefix(ocr2types.ConfigDigestPrefixCCIPMultiRoleRMNCombo, h)

		peerIDs := make([]string, 0, len(d.oraclePeerIDs))
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

		lggr.Infow("Bootstrapper is creating new peer group")
		peerGroup, err := d.peerGroupFactory.NewPeerGroup(
			[32]byte(genericEndpointConfigDigest),
			peerIDs,
			d.bootstrapLocators,
		)
		if err != nil {
			lggr.Errorw("failed to create new peer group", "err", err)
			return fmt.Errorf("new peer group: %w", err)
		}
		lggr.Infow("Created new peer group successfully")

		d.activePeerGroups = append(d.activePeerGroups, peerGroup)
		d.activeConfigDigests = append(d.activeConfigDigests, genericEndpointConfigDigest)
	}

	return nil
}

func writePrefix(prefix ocr2types.ConfigDigestPrefix, hash cciptypes.Bytes32) cciptypes.Bytes32 {
	var prefixBytes [2]byte
	binary.BigEndian.PutUint16(prefixBytes[:], uint16(prefix))

	hCopy := hash
	hCopy[0] = prefixBytes[0]
	hCopy[1] = prefixBytes[1]

	return hCopy
}
