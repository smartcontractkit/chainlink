package launcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/services/ccipcapability/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

var (
	_ job.ServiceCtx          = (*launcher)(nil)
	_ registrysyncer.Launcher = (*launcher)(nil)
)

func New(
	capabilityVersion,
	capabilityLabelledName string,
	p2pID ragep2ptypes.PeerID,
	lggr logger.Logger,
	homeChainReader cctypes.HomeChainReader,
	oracleCreator cctypes.OracleCreator,
	tickInterval time.Duration,
) *launcher {
	return &launcher{
		capabilityVersion:      capabilityVersion,
		capabilityLabelledName: capabilityLabelledName,
		p2pID:                  p2pID,
		lggr:                   lggr,
		homeChainReader:        homeChainReader,
		regState: registrysyncer.State{
			IDsToDONs:         make(map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo),
			IDsToNodes:        make(map[p2ptypes.PeerID]kcr.CapabilitiesRegistryNodeInfo),
			IDsToCapabilities: make(map[registrysyncer.HashedCapabilityID]kcr.CapabilitiesRegistryCapabilityInfo),
		},
		oracleCreator: oracleCreator,
		dons:          make(map[registrysyncer.DonID]*ccipDeployment),
		tickInterval:  tickInterval,
	}
}

// launcher manages the lifecycles of the CCIP capability on all chains.
type launcher struct {
	services.StateMachine

	capabilityVersion      string
	capabilityLabelledName string
	p2pID                  ragep2ptypes.PeerID
	lggr                   logger.Logger
	homeChainReader        cctypes.HomeChainReader
	stopChan               chan struct{}
	// latestState is the latest capability registry state received from the syncer.
	latestState registrysyncer.State
	// regState is the latest capability registry state that we have successfully processed.
	regState      registrysyncer.State
	oracleCreator cctypes.OracleCreator
	lock          sync.RWMutex
	wg            sync.WaitGroup
	tickInterval  time.Duration

	// dons is a map of CCIP DON IDs to the OCR instances that are running on them.
	// we can have up to two OCR instances per CCIP plugin, since we are running two plugins,
	// thats four OCR instances per CCIP DON maximum.
	dons map[registrysyncer.DonID]*ccipDeployment
}

// Launch implements registrysyncer.Launcher.
func (l *launcher) Launch(ctx context.Context, state registrysyncer.State) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lggr.Debugw("Received new state from syncer", "dons", state.IDsToDONs)
	l.latestState = state
	return nil
}

func (l *launcher) getLatestState() registrysyncer.State {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.latestState
}

func (l *launcher) runningDONIDs() []registrysyncer.DonID {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var runningDONs []registrysyncer.DonID
	for id := range l.dons {
		runningDONs = append(runningDONs, id)
	}
	return runningDONs
}

// Close implements job.ServiceCtx.
func (l *launcher) Close() error {
	return l.StateMachine.StopOnce("launcher", func() error {
		// shut down the monitor goroutine.
		close(l.stopChan)
		l.wg.Wait()

		// shut down all running oracles.
		var err error
		for _, ceDep := range l.dons {
			err = multierr.Append(err, ceDep.Close())
		}

		return err
	})
}

// Start implements job.ServiceCtx.
func (l *launcher) Start(context.Context) error {
	return l.StartOnce("launcher", func() error {
		l.stopChan = make(chan struct{})
		l.wg.Add(1)
		go l.monitor()
		return nil
	})
}

func (l *launcher) monitor() {
	defer l.wg.Done()
	ticker := time.NewTicker(l.tickInterval)
	for {
		select {
		case <-l.stopChan:
			return
		case <-ticker.C:
			if err := l.tick(); err != nil {
				l.lggr.Errorw("Failed to tick", "err", err)
			}
		}
	}
}

func (l *launcher) tick() error {
	// Ensure that the home chain reader is healthy.
	// For new jobs it may be possible that the home chain reader is not yet ready
	// so we won't be able to fetch configs and start any OCR instances.
	if ready := l.homeChainReader.Ready(); ready != nil {
		return fmt.Errorf("home chain reader is not ready: %w", ready)
	}

	// Fetch the latest state from the capability registry and determine if we need to
	// launch or update any OCR instances.
	latestState := l.getLatestState()

	diffRes, err := diff(l.capabilityVersion, l.capabilityLabelledName, l.regState, latestState)
	if err != nil {
		return fmt.Errorf("failed to diff capability registry states: %w", err)
	}

	err = l.processDiff(diffRes)
	if err != nil {
		return fmt.Errorf("failed to process diff: %w", err)
	}

	return nil
}

// processDiff processes the diff between the current and latest capability registry states.
// for any added OCR instances, it will launch them.
// for any removed OCR instances, it will shut them down.
// for any updated OCR instances, it will restart them with the new configuration.
func (l *launcher) processDiff(diff diffResult) error {
	err := l.processRemoved(diff.removed)
	err = multierr.Append(err, l.processAdded(diff.added))
	err = multierr.Append(err, l.processUpdate(diff.updated))

	return err
}

func (l *launcher) processUpdate(updated map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	for donID, don := range updated {
		prevDeployment, ok := l.dons[registrysyncer.DonID(don.Id)]
		if !ok {
			return fmt.Errorf("invariant violation: expected to find CCIP DON %d in the map of running deployments", don.Id)
		}

		futDeployment, err := updateDON(
			l.lggr,
			l.p2pID,
			l.homeChainReader,
			l.oracleCreator,
			*prevDeployment,
			don,
		)
		if err != nil {
			return err
		}
		if err := futDeployment.HandleBlueGreen(prevDeployment); err != nil {
			// TODO: how to handle a failed blue-green deployment?
			return fmt.Errorf("failed to handle blue-green deployment for CCIP DON %d: %w", donID, err)
		}

		// update state.
		l.dons[donID] = futDeployment
		// update the state with the latest config.
		// this way if one of the starts errors, we don't retry all of them.
		l.regState.IDsToDONs[donID] = updated[donID]
	}

	return nil
}

func (l *launcher) processAdded(added map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	for donID, don := range added {
		dep, err := createDON(
			l.lggr,
			l.p2pID,
			l.homeChainReader,
			l.oracleCreator,
			don,
		)
		if err != nil {
			return err
		}
		if dep == nil {
			// not a member of this DON.
			continue
		}

		if err := dep.StartBlue(); err != nil {
			if shutdownErr := dep.CloseBlue(); shutdownErr != nil {
				l.lggr.Errorw("Failed to shutdown blue instance after failed start", "donId", donID, "err", shutdownErr)
			}
			return fmt.Errorf("failed to start oracles for CCIP DON %d: %w", donID, err)
		}

		// update state.
		l.dons[donID] = dep
		// update the state with the latest config.
		// this way if one of the starts errors, we don't retry all of them.
		l.regState.IDsToDONs[donID] = added[donID]
	}

	return nil
}

func (l *launcher) processRemoved(removed map[registrysyncer.DonID]kcr.CapabilitiesRegistryDONInfo) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	for id := range removed {
		ceDep, ok := l.dons[id]
		if !ok {
			// not running this particular DON.
			continue
		}

		if err := ceDep.Close(); err != nil {
			return fmt.Errorf("failed to shutdown oracles for CCIP DON %d: %w", id, err)
		}

		// after a successful shutdown we can safely remove the DON deployment from the map.
		delete(l.dons, id)
		delete(l.regState.IDsToDONs, id)
	}

	return nil
}

// updateDON is a pure function that handles the case where a DON in the capability registry
// has received a new configuration.
// It returns a new ccipDeployment that can then be used to perform the blue-green deployment,
// based on the previous deployment.
func updateDON(
	lggr logger.Logger,
	p2pID ragep2ptypes.PeerID,
	homeChainReader cctypes.HomeChainReader,
	oracleCreator cctypes.OracleCreator,
	prevDeployment ccipDeployment,
	don kcr.CapabilitiesRegistryDONInfo,
) (futDeployment *ccipDeployment, err error) {
	if !isMemberOfDON(don, p2pID) {
		lggr.Infow("Not a member of this DON, skipping", "donId", don.Id, "p2pId", p2pID.String())
		return nil, nil
	}

	// this should be a retryable error.
	commitOCRConfigs, err := homeChainReader.GetOCRConfigs(context.Background(), don.Id, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OCR configs for CCIP commit plugin (don id: %d) from home chain config contract: %w",
			don.Id, err)
	}

	execOCRConfigs, err := homeChainReader.GetOCRConfigs(context.Background(), don.Id, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OCR configs for CCIP exec plugin (don id: %d) from home chain config contract: %w",
			don.Id, err)
	}

	commitBgd, err := createFutureBlueGreenDeployment(prevDeployment, commitOCRConfigs, oracleCreator, cctypes.PluginTypeCCIPCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to create future blue-green deployment for CCIP commit plugin: %w, don id: %d", err, don.Id)
	}

	execBgd, err := createFutureBlueGreenDeployment(prevDeployment, execOCRConfigs, oracleCreator, cctypes.PluginTypeCCIPExec)
	if err != nil {
		return nil, fmt.Errorf("failed to create future blue-green deployment for CCIP exec plugin: %w, don id: %d", err, don.Id)
	}

	return &ccipDeployment{
		commit: commitBgd,
		exec:   execBgd,
	}, nil
}

// valid cases:
// a) len(ocrConfigs) == 2 && !prevDeployment.HasGreenInstance(pluginType): this is a new green instance.
// b) len(ocrConfigs) == 1 && prevDeployment.HasGreenInstance(): this is a promotion of green->blue.
// All other cases are invalid. This is enforced in the ccip config contract.
func createFutureBlueGreenDeployment(
	prevDeployment ccipDeployment,
	ocrConfigs []ccipreaderpkg.OCR3ConfigWithMeta,
	oracleCreator cctypes.OracleCreator,
	pluginType cctypes.PluginType,
) (blueGreenDeployment, error) {
	var deployment blueGreenDeployment
	if isNewGreenInstance(pluginType, ocrConfigs, prevDeployment) {
		// this is a new green instance.
		greenOracle, err := oracleCreator.CreatePluginOracle(pluginType, cctypes.OCR3ConfigWithMeta(ocrConfigs[1]))
		if err != nil {
			return blueGreenDeployment{}, fmt.Errorf("failed to create CCIP commit oracle: %w", err)
		}

		deployment.blue = prevDeployment.commit.blue
		deployment.green = greenOracle
	} else if isPromotion(pluginType, ocrConfigs, prevDeployment) {
		// this is a promotion of green->blue.
		deployment.blue = prevDeployment.commit.green
	} else {
		return blueGreenDeployment{}, fmt.Errorf("invariant violation: expected 1 or 2 OCR configs for CCIP plugin (type: %d), got %d", pluginType, len(ocrConfigs))
	}

	return deployment, nil
}

// createDON is a pure function that handles the case where a new DON is added to the capability registry.
// It returns a new ccipDeployment that can then be used to start the blue instance.
func createDON(
	lggr logger.Logger,
	p2pID ragep2ptypes.PeerID,
	homeChainReader cctypes.HomeChainReader,
	oracleCreator cctypes.OracleCreator,
	don kcr.CapabilitiesRegistryDONInfo,
) (*ccipDeployment, error) {
	if !isMemberOfDON(don, p2pID) {
		lggr.Infow("Not a member of this DON, skipping", "donId", don.Id, "p2pId", p2pID.String())
		return nil, nil
	}

	// this should be a retryable error.
	commitOCRConfigs, err := homeChainReader.GetOCRConfigs(context.Background(), don.Id, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OCR configs for CCIP commit plugin (don id: %d) from home chain config contract: %w",
			don.Id, err)
	}

	execOCRConfigs, err := homeChainReader.GetOCRConfigs(context.Background(), don.Id, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OCR configs for CCIP exec plugin (don id: %d) from home chain config contract: %w",
			don.Id, err)
	}

	// upon creation we should only have one OCR config per plugin type.
	if len(commitOCRConfigs) != 1 {
		return nil, fmt.Errorf("expected exactly one OCR config for CCIP commit plugin (don id: %d), got %d", don.Id, len(commitOCRConfigs))
	}

	if len(execOCRConfigs) != 1 {
		return nil, fmt.Errorf("expected exactly one OCR config for CCIP exec plugin (don id: %d), got %d", don.Id, len(execOCRConfigs))
	}

	commitOracle, commitBootstrap, err := createOracle(p2pID, oracleCreator, cctypes.PluginTypeCCIPCommit, commitOCRConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIP commit oracle: %w", err)
	}

	execOracle, execBootstrap, err := createOracle(p2pID, oracleCreator, cctypes.PluginTypeCCIPExec, execOCRConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to create CCIP exec oracle: %w", err)
	}

	return &ccipDeployment{
		commit: blueGreenDeployment{
			blue:          commitOracle,
			bootstrapBlue: commitBootstrap,
		},
		exec: blueGreenDeployment{
			blue:          execOracle,
			bootstrapBlue: execBootstrap,
		},
	}, nil
}

func createOracle(
	p2pID ragep2ptypes.PeerID,
	oracleCreator cctypes.OracleCreator,
	pluginType cctypes.PluginType,
	ocrConfigs []ccipreaderpkg.OCR3ConfigWithMeta,
) (pluginOracle, bootstrapOracle cctypes.CCIPOracle, err error) {
	pluginOracle, err = oracleCreator.CreatePluginOracle(pluginType, cctypes.OCR3ConfigWithMeta(ocrConfigs[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CCIP plugin oracle (plugintype: %d): %w", pluginType, err)
	}

	if isMemberOfBootstrapSubcommittee(ocrConfigs[0].Config.BootstrapP2PIds, p2pID) {
		bootstrapOracle, err = oracleCreator.CreateBootstrapOracle(cctypes.OCR3ConfigWithMeta(ocrConfigs[0]))
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create CCIP bootstrap oracle (plugintype: %d): %w", pluginType, err)
		}
	}

	return pluginOracle, bootstrapOracle, nil
}
