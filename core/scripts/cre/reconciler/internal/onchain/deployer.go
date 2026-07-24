// Package onchain owns the CRE on-chain deployment + configuration engine and the
// post-environment job-sync engine: building the CLDF/CRE environment, topology,
// deploying registry contracts, configuring CapabilitiesRegistry/WorkflowRegistry,
// preparing JD node chain configs, and creating/recreating jobs via JD.
package onchain

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	cldfjd "github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"

	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/domain"
	griddleinfra "github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/infra"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/reconciler/internal/jobs"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

// K8sAPI is the narrow Kubernetes surface the on-chain engine needs — a subset of
// the root package's K8sAPI (declared here so onchain doesn't import root).
type K8sAPI interface {
	GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*griddleinfra.NodeAPIInfo, error)
	GetNodeSecretsToml(ctx context.Context, nodeName, namespace string) (string, error)
}

// ConfirmFunc is the per-step confirm gate, injected from the CLI layer. It
// returns nil to proceed, or an error (e.g. ErrStepDeclined-equivalent) to abort.
type ConfirmFunc func(title, details string) error

// Deployer builds the CLDF/CRE environment, deploys and configures the on-chain
// registries, and syncs jobs via JD. It operates on the caller-owned
// *domain.StateFile passed into Apply/SyncJobs: it reads prior results for
// idempotency (HasAddress/DONIDs) and writes new results back into it.
type Deployer struct {
	k8s         K8sAPI
	deployerKey string
	log         zerolog.Logger
	confirm     ConfirmFunc

	// transient caches populated by runPreEnvStartup and consumed by configureCapReg,
	// reset implicitly at the start of each Apply (a fresh Deployer per run).
	donsCapabilities                map[uint64][]keystone_changeset.DONCapabilityWithConfig
	capabilityToOCR3Config          map[string]*ocr3.OracleConfig
	capabilityToExtraSignerFamilies map[string][]string

	mu sync.Mutex // guards state.JDNodeIDs writes during parallel JD chain-config prep
}

// NewDeployer creates a Deployer. confirm gates each on-chain/job step; it should
// return nil to proceed or a non-nil error (propagated to the caller) to abort.
func NewDeployer(k8s K8sAPI, deployerKey string, log zerolog.Logger, confirm ConfirmFunc) *Deployer {
	return &Deployer{
		k8s:         k8s,
		deployerKey: deployerKey,
		log:         log,
		confirm:     confirm,
	}
}

// skipUnlessConfirmed gates a step behind the injected ConfirmFunc.
func (d *Deployer) skipUnlessConfirmed(title, details string) error {
	return d.confirm(title, details)
}

// Apply implements P1-P8 on-chain deployment and configuration. It must only be
// called when the caller has determined on-chain configuration is not yet
// complete — it does not re-check that itself. persist is called after each
// state-mutating milestone so a crash mid-flow can resume from the last
// successfully persisted step.
func (d *Deployer) Apply(
	ctx context.Context,
	desired *domain.DesiredState,
	cv *domain.ChartValues,
	state *domain.StateFile,
	persist func(),
) error {
	d.log.Info().Msg("=== P1-P6: On-chain configuration ===")

	// The JD access token is validated once, up front, in Run (requireJDAccessToken) —
	// before discovery even starts — so no redundant check is needed here.

	if err := d.skipUnlessConfirmed("P1: Build CLDF environment", d.onChainEnvSummary(desired, cv)); err != nil {
		return err
	}
	env, chainSelector, err := d.buildCldfEnv(ctx, desired)
	if err != nil {
		return errors.Wrap(err, "failed to build cldf environment")
	}
	d.log.Info().
		Uint64("chainSelector", chainSelector).
		Str("envName", env.Name).
		Bool("hasJD", env.Offchain != nil).
		Msg("Built cldf environment")

	if err = d.skipUnlessConfirmed("P2: Build topology from desired state + K8s secrets", d.topologySummary(desired, cv)); err != nil {
		return err
	}
	topology, err := d.buildTopology(ctx, desired, cv, state)
	if err != nil {
		return errors.Wrap(err, "failed to build topology")
	}
	d.storeGatewayConnectors(desired, cv, state, topology)
	for _, donMeta := range topology.DonsMetadata.List() {
		d.log.Info().
			Str("don", donMeta.Name).
			Uint64("id", donMeta.ID).
			Strs("flags", donMeta.Flags).
			Int("nodes", len(donMeta.NodesMetadata)).
			Msg("Topology DON")
	}

	if err = d.skipUnlessConfirmed("P3: Build CRE environment (blockchains + contract versions)", fmt.Sprintf("Registry chain selector: %d\nDeclared chains: %d", chainSelector, len(desired.Chains))); err != nil {
		return err
	}
	creEnv, err := d.buildCreEnvironment(ctx, desired, env, chainSelector)
	if err != nil {
		return errors.Wrap(err, "failed to build CRE environment")
	}
	d.log.Info().
		Uint64("registryChainSelector", creEnv.RegistryChainSelector).
		Int("blockchains", len(creEnv.Blockchains)).
		Int("contractVersions", len(creEnv.ContractVersions)).
		Str("provider", creEnv.Provider.Type).
		Msg("Built CRE environment")

	deploySummary := "Deploy CapabilitiesRegistry + WorkflowRegistry v2 to Anvil"
	if state.HasAddress(keystone_changeset.CapabilitiesRegistry.String()) {
		deploySummary = "Contracts already deployed — will hydrate addresses from state"
	}
	if err = d.skipUnlessConfirmed("P4: Deploy registry contracts", deploySummary); err != nil {
		return err
	}
	if err = d.deployContracts(env, chainSelector, state); err != nil {
		return errors.Wrap(err, "failed to deploy contracts")
	}
	if err = d.syncAddressBook(env, state); err != nil {
		return errors.Wrap(err, "failed to sync address book after P4")
	}
	persist()

	if err = d.skipUnlessConfirmed("P5: Run Features.PreEnvStartup", "Deploy forwarders and collect capability configs per DON"); err != nil {
		return err
	}
	if err = d.runPreEnvStartup(ctx, topology, creEnv); err != nil {
		return errors.Wrap(err, "failed to run PreEnvStartup")
	}
	if err = d.syncAddressBook(env, state); err != nil {
		return errors.Wrap(err, "failed to sync address book after P5")
	}
	d.storeGatewayServiceConfigs(state, topology)
	persist()

	if err = d.skipUnlessConfirmed("P5b: Prepare JD node chain configs for CapReg", d.jdChainConfigSummary(desired, topology, state)); err != nil {
		return err
	}
	if err = d.prepareJDChainConfigsForCapReg(ctx, topology, creEnv, chainSelector, desired, cv, state); err != nil {
		return errors.Wrap(err, "failed to prepare JD chain configs")
	}

	if err = d.skipUnlessConfirmed("P6: Configure Capabilities Registry on-chain", "CapReg: "+state.GetAddress(keystone_changeset.CapabilitiesRegistry.String())); err != nil {
		return err
	}
	capReg, err := d.configureCapReg(topology, creEnv, chainSelector, state)
	if err != nil {
		return errors.Wrap(err, "failed to configure capabilities registry")
	}

	if err := d.skipUnlessConfirmed("P7: Resolve DON IDs from CapReg contract", fmt.Sprintf("DONs: %d", len(desired.DONs))); err != nil {
		return err
	}
	if err := d.resolveDONIDs(capReg, desired, cv, state); err != nil {
		return errors.Wrap(err, "failed to resolve DON IDs")
	}
	if err := d.syncAddressBook(env, state); err != nil {
		return errors.Wrap(err, "failed to sync address book after P6/P7")
	}
	persist()

	if err := d.skipUnlessConfirmed("P8: Configure Workflow Registry on-chain", "WfReg: "+state.GetAddress(keystone_changeset.WorkflowRegistry.String())); err != nil {
		return err
	}
	if err := d.configureWorkflowReg(ctx, env, chainSelector, desired, state); err != nil {
		return errors.Wrap(err, "failed to configure workflow registry")
	}
	if err := d.syncAddressBook(env, state); err != nil {
		return errors.Wrap(err, "failed to sync address book after P8")
	}

	state.Phase = domain.PhaseOnChain
	persist()
	d.log.Info().Msg("On-chain configuration complete")
	return nil
}

// SyncJobs builds env/topology/dons, restores persisted gateway service configs onto
// the fresh topology, deletes existing jobs, and runs each feature's PostEnvStartup +
// gateway job creation. Must only be called once on-chain configuration is complete.
func (d *Deployer) SyncJobs(ctx context.Context, desired *domain.DesiredState, cv *domain.ChartValues, state *domain.StateFile) error {
	d.log.Info().Msg("=== J1: Job creation ===")

	if !onChainComplete(state) {
		return errors.New("on-chain configuration must complete before job creation")
	}

	env, chainSelector, err := d.buildCldfEnv(ctx, desired)
	if err != nil {
		return errors.Wrap(err, "failed to build cldf environment for jobs")
	}
	memDs, err := hydrateMemoryDataStoreFromState(state, chainSelector)
	if err != nil {
		return errors.Wrap(err, "failed to hydrate contract addresses for jobs")
	}
	env.DataStore = memDs.Seal()

	topology, err := d.buildTopology(ctx, desired, cv, state)
	if err != nil {
		return errors.Wrap(err, "failed to build topology for jobs")
	}
	// PreEnvStartup does not run in the jobs phase, so the fresh topology has only
	// the default web-api-capabilities handler. Restore the handler set (e.g.
	// http-capabilities) that Features populated and we persisted during the
	// on-chain phase; otherwise the generated gateway job omits those handlers.
	d.applyStoredGatewayServiceConfigs(state, topology)

	creEnv, err := d.buildCreEnvironment(ctx, desired, env, chainSelector)
	if err != nil {
		return errors.Wrap(err, "failed to build CRE environment for jobs")
	}
	creEnv.FreshExternalJobIDs = true

	if creEnv.CldfEnvironment.Offchain == nil {
		return errors.New("JD client is required for job creation")
	}
	jd, ok := creEnv.CldfEnvironment.Offchain.(*cldfjd.JobDistributor)
	if !ok {
		return fmt.Errorf("offchain client must be a Job Distributor, got %T", creEnv.CldfEnvironment.Offchain)
	}

	if err = d.skipUnlessConfirmed("J1: Build DONs for PostEnvStartup", d.jobCreationSummary(desired, topology)); err != nil {
		return err
	}
	dons, err := d.buildDonsForJobs(ctx, topology, jd, desired, cv, state)
	if err != nil {
		return errors.Wrap(err, "failed to build DONs for jobs")
	}

	if err := d.skipUnlessConfirmed("J1a: Delete all existing jobs before recreating", d.jobCreationSummary(desired, topology)); err != nil {
		return err
	}
	if err := jobs.DeleteAllForDons(ctx, d.log, jobs.NewJobDistributorAdapter(jd), jobs.GQLProposalCanceller{}, dons); err != nil {
		return errors.Wrap(err, "failed to delete existing jobs")
	}

	if err := d.skipUnlessConfirmed("J1b: PostEnvStartup (gateway jobs + capability jobs)", d.jobCreationSummary(desired, topology)); err != nil {
		return err
	}
	if err := d.runPostEnvStartup(ctx, desired, creEnv, topology, dons); err != nil {
		return errors.Wrap(err, "failed to run PostEnvStartup")
	}

	state.Phase = domain.PhaseJobs
	d.log.Info().Msg("Job creation complete")
	return nil
}

// onChainComplete mirrors the root package's on-chain-complete check: P1-P8
// actually finished (contracts deployed, DON IDs resolved, workflow registry
// configured). SyncJobs uses it as its own precondition guard since it cannot
// call back into the root package.
func onChainComplete(state *domain.StateFile) bool {
	return state.WorkflowReg != nil &&
		state.HasAddress(keystone_changeset.CapabilitiesRegistry.String()) &&
		len(state.DONIDs) > 0
}

func (d *Deployer) onChainEnvSummary(desired *domain.DesiredState, cv *domain.ChartValues) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Namespace: %s\n", cv.Namespace)
	fmt.Fprintf(&b, "Declared chains: %d\n", len(desired.Chains))
	if registryChain, ok := desired.RegistryChain(); ok {
		fmt.Fprintf(&b, "Registry chain ID: %d\n", registryChain.ChainID)
	}
	fmt.Fprintf(&b, "JD gRPC: %s\n", desired.JD.GRPC)
	fmt.Fprintf(&b, "JD configured: %t\n", griddleinfra.JDAccessToken() != "")
	return b.String()
}

func (d *Deployer) topologySummary(desired *domain.DesiredState, cv *domain.ChartValues) string {
	var b strings.Builder
	fmt.Fprintf(&b, "DONs in desired state: %d\n", len(desired.DONs))
	for _, don := range desired.DONs {
		if don.IsBootstrapOnly(cv) || don.IsGatewayDon() {
			continue
		}
		fmt.Fprintf(&b, "  • %s: %d workers, caps=%v\n", don.Name, len(don.WorkerNodes(cv)), don.Capabilities)
	}
	if desired.NeedsGateway() {
		fmt.Fprintf(&b, "Gateway nodes: %d\n", len(cv.FindGatewayNodes()))
	}
	return b.String()
}
