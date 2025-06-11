package launcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/multierr"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragep2ptypes "github.com/smartcontractkit/libocr/ragep2p/types"

	ccipreader "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

var (
	_ job.ServiceCtx          = (*launcher)(nil)
	_ registrysyncer.Launcher = (*launcher)(nil)
)

// New creates a new instance of the CCIP launcher.
func New(
	capabilityID string,
	myP2PID ragep2ptypes.PeerID,
	lggr logger.Logger,
	homeChainReader ccipreader.HomeChain,
	tickInterval time.Duration,
	oracleCreator cctypes.OracleCreator,
) *launcher {
	return &launcher{
		myP2PID:         myP2PID,
		capabilityID:    capabilityID,
		lggr:            lggr,
		homeChainReader: homeChainReader,
		regState: registrysyncer.LocalRegistry{
			IDsToDONs:         make(map[registrysyncer.DonID]registrysyncer.DON),
			IDsToNodes:        make(map[p2ptypes.PeerID]kcr.INodeInfoProviderNodeInfo),
			IDsToCapabilities: make(map[string]registrysyncer.Capability),
		},
		tickInterval:  tickInterval,
		oracleCreator: oracleCreator,
		instances:     make(map[registrysyncer.DonID]pluginRegistry),
	}
}

// launcher manages the lifecycles of the CCIP capability on all chains.
type launcher struct {
	services.StateMachine

	// capabilityID is the fully qualified capability registry ID of the CCIP capability.
	// this is <capability_name>@<capability-semver>, e.g "ccip@1.0.0".
	capabilityID string

	// myP2PID is the peer ID of the node running this launcher.
	myP2PID         ragep2ptypes.PeerID
	lggr            logger.Logger
	homeChainReader ccipreader.HomeChain
	stopChan        services.StopChan
	// latestState is the latest capability registry state received from the syncer.
	latestState registrysyncer.LocalRegistry
	// regState is the latest capability registry state that we have successfully processed.
	regState      registrysyncer.LocalRegistry
	oracleCreator cctypes.OracleCreator
	lock          sync.RWMutex
	wg            sync.WaitGroup
	tickInterval  time.Duration

	// instances is a map of CCIP DON IDs to a map of the OCR instances that are running on them.
	// This map uses the config digest as the key, and the instance as the value.
	// We can have up to a maximum of 4 instances per CCIP DON (active/candidate) x (commit/exec)
	instances map[registrysyncer.DonID]pluginRegistry
}

// Launch implements registrysyncer.Launcher.
func (l *launcher) Launch(ctx context.Context, state *registrysyncer.LocalRegistry) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lggr.Debugw("Received new state from syncer", "dons", state.IDsToDONs)
	l.latestState = *state
	return nil
}

func (l *launcher) getLatestState() registrysyncer.LocalRegistry {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.latestState
}

func (l *launcher) runningDONIDs() []registrysyncer.DonID {
	l.lock.RLock()
	defer l.lock.RUnlock()
	var runningDONs []registrysyncer.DonID
	for id := range l.instances {
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
		for _, ceDep := range l.instances {
			err = multierr.Append(err, ceDep.CloseAll())
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

// monitor calls tick() at regular intervals to check for changes in the capability registry.
func (l *launcher) monitor() {
	defer l.wg.Done()
	ticker := time.NewTicker(l.tickInterval)

	ctx, cancel := l.stopChan.NewCtx()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := l.tick(ctx); err != nil {
				l.lggr.Errorw("Failed to tick", "err", err)
			}
		}
	}
}

// tick gets the latest registry state and processes the diff between the current and latest state.
// This may lead to starting or stopping OCR instances.
func (l *launcher) tick(ctx context.Context) error {
	// Ensure that the home chain reader is healthy.
	// For new jobs it may be possible that the home chain reader is not yet ready
	// so we won't be able to fetch configs and start any OCR instances.
	if ready := l.homeChainReader.Ready(); ready != nil {
		return fmt.Errorf("home chain reader is not ready: %w", ready)
	}

	// Fetch the latest state from the capability registry and determine if we need to
	// launch or update any OCR instances.
	latestState := l.getLatestState()

	diffRes, err := diff(l.capabilityID, l.regState, latestState)
	if err != nil {
		return fmt.Errorf("failed to diff capability registry states: %w", err)
	}

	err = l.processDiff(ctx, diffRes)
	if err != nil {
		return fmt.Errorf("failed to process diff: %w", err)
	}

	return nil
}

// processDiff processes the diff between the current and latest capability registry states.
// for any added OCR instances, it will launch them.
// for any removed OCR instances, it will shut them down.
// for any updated OCR instances, it will restart them with the new configuration.
func (l *launcher) processDiff(ctx context.Context, diff diffResult) error {
	err := l.processRemoved(diff.removed)
	err = multierr.Append(err, l.processAdded(ctx, diff.added))
	err = multierr.Append(err, l.processUpdate(ctx, diff.updated))

	return err
}

// processUpdate will manage when configurations of an existing don are updated
// If new oracles are needed, they are created and started. Old ones will be shut down
func (l *launcher) processUpdate(ctx context.Context, updated map[registrysyncer.DonID]registrysyncer.DON) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	for donID, don := range updated {
		prevPlugins, ok := l.instances[donID]
		if !ok {
			return fmt.Errorf("invariant violation: expected to find CCIP DON %d in the map of running deployments", don.ID)
		}

		latestConfigs, err := getConfigsForDon(ctx, l.homeChainReader, don)
		if err != nil {
			return err
		}

		newPlugins, err := updateDON(
			ctx,
			l.lggr,
			l.myP2PID,
			prevPlugins,
			don,
			l.oracleCreator,
			latestConfigs)
		if err != nil {
			return err
		}
		if len(newPlugins) == 0 {
			// not a member of this DON.
			continue
		}

		err = newPlugins.TransitionFrom(prevPlugins)
		if err != nil {
			return fmt.Errorf("could not transition state %w", err)
		}

		l.instances[donID] = newPlugins
		l.regState.IDsToDONs[donID] = updated[donID]
	}

	return nil
}

// processAdded is for when a new don is created. We know that all oracles
// must be created and started
func (l *launcher) processAdded(ctx context.Context, added map[registrysyncer.DonID]registrysyncer.DON) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	for donID, don := range added {
		configs, err := getConfigsForDon(ctx, l.homeChainReader, don)
		if err != nil {
			return fmt.Errorf("failed to get current configs for don %d: %w", donID, err)
		}
		newPlugins, err := createDON(
			ctx,
			l.lggr,
			l.myP2PID,
			don,
			l.oracleCreator,
			configs,
		)
		if err != nil {
			return fmt.Errorf("processAdded: call createDON %d: %w", donID, err)
		}
		if len(newPlugins) == 0 {
			// not a member of this DON.
			continue
		}

		// now that oracles are created, we need to start them. If there are issues with starting
		// we should shut them down
		if err := newPlugins.StartAll(); err != nil {
			if shutdownErr := newPlugins.CloseAll(); shutdownErr != nil {
				l.lggr.Errorw("Failed to shutdown don instances after a failed start", "donId", donID, "err", shutdownErr)
			}
			return fmt.Errorf("processAdded: start oracles for CCIP DON %d: %w", donID, err)
		}

		// update state.
		l.instances[donID] = newPlugins
		l.regState.IDsToDONs[donID] = added[donID]
	}

	return nil
}

// processRemoved handles the situation when an entire DON is removed
func (l *launcher) processRemoved(removed map[registrysyncer.DonID]registrysyncer.DON) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	for id := range removed {
		p, ok := l.instances[id]
		if !ok {
			// not running this particular DON.
			continue
		}

		if err := p.CloseAll(); err != nil {
			return fmt.Errorf("failed to shutdown oracles for CCIP DON %d: %w", id, err)
		}

		// after a successful shutdown we can safely remove the DON deployment from the map.
		delete(l.instances, id)
		delete(l.regState.IDsToDONs, id)
	}

	return nil
}

func updateDON(
	ctx context.Context,
	lggr logger.Logger,
	p2pID ragep2ptypes.PeerID,
	prevPlugins pluginRegistry,
	don registrysyncer.DON,
	oracleCreator cctypes.OracleCreator,
	latestConfigs []ccipreader.OCR3ConfigWithMeta,
) (pluginRegistry, error) {
	if !isMemberOfDON(don, p2pID) {
		lggr.Infow("Not a member of this DON, skipping", "donID", don.ID, "p2pID", p2pID.String())
		return nil, nil
	}

	newP := make(pluginRegistry)
	// If a config digest is not already in our list, we need to create an oracle
	// If a config digest is already in our list, we just need to point to the old one
	// newP.Transition will make sure we shut down the old oracles, and start the new ones
	for _, c := range latestConfigs {
		digest := c.ConfigDigest
		if _, ok := prevPlugins[digest]; !ok {
			oracle, err := oracleCreator.Create(ctx, don.ID, cctypes.OCR3ConfigWithMeta(c))
			if err != nil {
				return nil, fmt.Errorf("failed to create CCIP oracle: %w for digest %x", err, digest)
			}

			newP[digest] = oracle
		} else {
			newP[digest] = prevPlugins[digest]
		}
	}

	return newP, nil
}

// createDON is a pure function that handles the case where a new DON is added to the capability registry.
// It returns up to 4 plugins that are later started.
func createDON(
	ctx context.Context,
	lggr logger.Logger,
	p2pID ragep2ptypes.PeerID,
	don registrysyncer.DON,
	oracleCreator cctypes.OracleCreator,
	configs []ccipreader.OCR3ConfigWithMeta,
) (pluginRegistry, error) {
	if !isMemberOfDON(don, p2pID) && oracleCreator.Type() == cctypes.OracleTypePlugin {
		lggr.Infow("Not a member of this DON and not a bootstrap node either, skipping", "donID", don.ID, "p2pID", p2pID.String())
		return nil, nil
	}
	p := make(pluginRegistry)
	for _, config := range configs {
		digest, err := ocrtypes.BytesToConfigDigest(config.ConfigDigest[:])
		if err != nil {
			return nil, fmt.Errorf("digest does not match type %w", err)
		}

		oracle, err := oracleCreator.Create(ctx, don.ID, cctypes.OCR3ConfigWithMeta(config))
		if err != nil {
			return nil, fmt.Errorf("failed to create CCIP oracle: %w for digest %x", err, digest)
		}

		p[digest] = oracle
	}
	return p, nil
}

func getConfigsForDon(
	ctx context.Context,
	homeChainReader ccipreader.HomeChain,
	don registrysyncer.DON) ([]ccipreader.OCR3ConfigWithMeta, error) {
	// this should be a retryable error.
	commitOCRConfigs, err := homeChainReader.GetOCRConfigs(ctx, don.ID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OCR configs for CCIP commit plugin (don id: %d) from home chain config contract: %w",
			don.ID, err)
	}

	execOCRConfigs, err := homeChainReader.GetOCRConfigs(ctx, don.ID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OCR configs for CCIP exec plugin (don id: %d) from home chain config contract: %w",
			don.ID, err)
	}

	c := []ccipreader.OCR3ConfigWithMeta{
		commitOCRConfigs.CandidateConfig,
		commitOCRConfigs.ActiveConfig,
		execOCRConfigs.CandidateConfig,
		execOCRConfigs.ActiveConfig,
	}

	ret := make([]ccipreader.OCR3ConfigWithMeta, 0, 4)
	for _, config := range c {
		if config.ConfigDigest != [32]byte{} {
			ret = append(ret, config)
		}
	}

	return ret, nil
}
