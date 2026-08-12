package cre

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	crescriptenv "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/environment"
	crelib "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	creenv "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
	feature_sets "github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/sets"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

// This file provides an in-process CRE environment start for tests that must
// supply capabilities and features computed at runtime.
//
// The standard helper (t_helpers.SetupTestEnvironmentWithConfig) starts the
// environment by shelling out to `cre env start`, so a test cannot hand it Go
// values. The confidential workflows test needs to, twice over: the enclave list
// (only known once the enclaves are running) goes into the capability's on-chain
// registry config, and the trusted measurements go into the relay's capability
// config.
//
// startConfidentialCreEnvironment therefore calls the exported
// crescriptenv.StartCLIEnvironment directly and writes the local CRE state file.
// Callers then invoke the standard helper, which finds that state file and skips
// starting the environment, building the TestEnvironment from saved state.

const (
	// confidentialCleanupWait matches the CLI's own cleanup grace period.
	confidentialCleanupWait = 15 * time.Second

	// Resolved against the CRE environment directory rather than used as-is: the
	// CLI defaults assume a working directory of core/scripts/cre/environment,
	// but this test runs from smoke/cre.
	confidentialCapabilityDefaultsConfig = "configs/capability_defaults.toml"
)

// startConfidentialCreEnvironment starts a local CRE environment with extra
// capabilities and features, then persists the state file so the standard test
// environment helper can build from it.
func startConfidentialCreEnvironment(
	ctx context.Context,
	relativePathToRepoRoot string,
	environmentDirPath string,
	extraCapabilities []crelib.InstallableCapability,
	extraAllowedPorts []int,
) error {
	in, err := confidentialPreConfigure(relativePathToRepoRoot, environmentDirPath)
	if err != nil {
		return err
	}

	// Extra capabilities are registered dynamically rather than declared in the
	// topology config. Job spec delivery and on-chain registration both check
	// don.HasFlag(name), which reads NodeSet.Capabilities, so the flags have to be
	// present there too.
	for _, c := range extraCapabilities {
		flag := c.Flag()
		for i, ns := range in.NodeSets {
			if slices.Contains(ns.DONTypes, "workflow") && !slices.Contains(ns.Capabilities, flag) {
				in.NodeSets[i].Capabilities = append(in.NodeSets[i].Capabilities, flag)
			}
		}
	}

	// Config.Validate rejects capability flags the provider doesn't know, so every
	// flag this run declares has to be present. don-time is absent from the
	// extensible provider's built-in globals, unlike the default provider's.
	extraFlags := []string{string(crelib.DONTimeCapability)}
	for _, c := range extraCapabilities {
		extraFlags = append(extraFlags, c.Flag())
	}

	envDependencies := crelib.NewEnvironmentDependencies(
		flags.NewExtensibleCapabilityFlagsProvider(extraFlags),
		crelib.NewContractVersionsProvider(envconfig.DefaultContractSet()),
	)
	if err := in.Validate(envDependencies); err != nil {
		return errors.Wrap(err, "failed to validate environment configuration")
	}

	// The standard feature set carries the confidential relay; the topology's
	// capability config supplies its settings.
	features := feature_sets.New()

	// Same allowlist `cre env start` grants, plus this test's enclave ports.
	gatewayWhitelistConfig := crescriptenv.DefaultGatewayWhitelistConfig(in, extraAllowedPorts)

	output, startErr := crescriptenv.StartCLIEnvironment(
		ctx,
		relativePathToRepoRoot,
		in,
		extraCapabilities,
		features,
		nil, // no extra job spec functions
		envDependencies,
		gatewayWhitelistConfig,
	)
	if startErr != nil {
		if stopErr := stopConfidentialCreEnvironment(relativePathToRepoRoot); stopErr != nil {
			return errors.Wrapf(startErr, "failed to start environment, and cleanup also failed: %s", stopErr)
		}
		return errors.Wrap(startErr, "failed to start environment")
	}

	addresses, aErr := output.CreEnvironment.CldfEnvironment.DataStore.Addresses().Fetch()
	if aErr != nil {
		return errors.Wrap(aErr, "failed to fetch addresses from datastore")
	}
	if err := in.SetAddresses(addresses); err != nil {
		return errors.Wrap(err, "failed to set addresses on config")
	}
	if storeErr := in.Store(envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)); storeErr != nil {
		return errors.Wrap(storeErr, "failed to store local CRE state")
	}

	return nil
}

// confidentialPreConfigure clears any prior environment state and loads the
// topology config, so a stale state file from a different topology is not
// merged into this run.
//
// This deliberately skips crescriptenv.RunSetup, matching `cre env start`,
// which only runs setup under --auto-setup (off by default, and unset in CI).
// Setup installs host tooling such as Bun, which is unavailable on CI runners.
func confidentialPreConfigure(relativePathToRepoRoot, environmentDirPath string) (*envconfig.Config, error) {
	_ = stopConfidentialCreEnvironment(relativePathToRepoRoot)

	if err := framework.RemoveTestContainers(); err != nil {
		return nil, errors.Wrap(err, "failed to remove test containers")
	}
	defer func() {
		crescriptenv.StartCmdRecoverHandlerFunc(nil, nil, true, confidentialCleanupWait)
	}()

	if cleanUpErr := envconfig.RemoveAllEnvironmentStateDir(relativePathToRepoRoot); cleanUpErr != nil {
		return nil, errors.Wrap(cleanUpErr, "failed to clean up environment state files")
	}

	// Re-prepend the capability defaults to whatever CTF_CONFIGS the caller set.
	// Stripping the prefix first keeps this idempotent across repeated calls.
	defaultsConfig := filepath.Join(environmentDirPath, confidentialCapabilityDefaultsConfig)
	userConfigs := strings.TrimPrefix(os.Getenv("CTF_CONFIGS"), defaultsConfig+",")
	ctfConfigs := defaultsConfig
	if userConfigs != "" && userConfigs != defaultsConfig {
		ctfConfigs = defaultsConfig + "," + userConfigs
	}
	if err := os.Setenv("CTF_CONFIGS", ctfConfigs); err != nil {
		return nil, fmt.Errorf("failed to set CTF_CONFIGS: %w", err)
	}

	if pkErr := creenv.SetDefaultPrivateKeyIfEmpty(blockchain.DefaultAnvilPrivateKey); pkErr != nil {
		return nil, errors.Wrap(pkErr, "failed to set default private key")
	}

	// Keep Ryuk from reaping the containers when this process exits; the test
	// tears them down itself.
	if setErr := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); setErr != nil {
		return nil, fmt.Errorf("failed to set TESTCONTAINERS_RYUK_DISABLED: %w", setErr)
	}

	in := &envconfig.Config{}
	if err := in.Load(os.Getenv("CTF_CONFIGS")); err != nil {
		return nil, errors.Wrap(err, "failed to load environment configuration")
	}

	return in, nil
}

// stopConfidentialCreEnvironment removes the environment containers and the local
// CRE state file.
func stopConfidentialCreEnvironment(relativePathToRepoRoot string) error {
	if removeErr := framework.RemoveTestContainers(); removeErr != nil {
		return errors.Wrap(removeErr, "failed to remove environment containers")
	}

	creStateFile := envconfig.MustLocalCREStateFileAbsPath(relativePathToRepoRoot)
	if cErr := os.Remove(creStateFile); cErr != nil && !os.IsNotExist(cErr) {
		framework.L.Warn().Msgf("failed to remove local CRE state file: %s", cErr)
	}

	return nil
}
