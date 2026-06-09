package solana

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

type DeploySolanaMCMS struct{}

var _ cldf.ChangeSetV2[DeploySolanaMCMSConfig] = DeploySolanaMCMS{}

// DeploySolanaMCMSConfig mirrors the CCIP Solana deploy input shape: optional BuildConfig plus
// MCMSWithTimelockConfig for deploy/init, with CRE-specific metadata for qualifiers and labels.
type DeploySolanaMCMSConfig struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`

	// BuildConfig is used to clone/build or download MCMS program artifacts before deploy.
	BuildConfig *helpers.BuildSolanaConfig `json:"buildConfig,omitempty" yaml:"buildConfig,omitempty"`

	// MCMSWithTimelockConfig is required. Programs are deployed and roles initialized with this policy.
	MCMSWithTimelockConfig *cldfproposalutils.MCMSWithTimelockConfig `json:"mcmsWithTimelockConfig,omitempty" yaml:"mcmsWithTimelockConfig,omitempty"`

	Environment string            `json:"environment,omitempty" yaml:"environment,omitempty"`
	ConfigID    string            `json:"configId" yaml:"configId"`
	Descriptor  *string           `json:"descriptor,omitempty" yaml:"descriptor,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
}

func (c DeploySolanaMCMSConfig) Qualifier() string {
	u := &url.URL{
		Scheme: "contract",
		Host:   strconv.FormatUint(c.ChainSelector, 10),
		Path:   "mcmsv2",
	}

	q := u.Query()
	q.Add("mcms-config", c.ConfigID)
	if c.Descriptor != nil {
		q.Add("descriptor", *c.Descriptor)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (DeploySolanaMCMS) VerifyPreconditions(env cldf.Environment, cfg DeploySolanaMCMSConfig) error {
	if cfg.MCMSWithTimelockConfig == nil {
		return errors.New("mcmsWithTimelockConfig is required")
	}

	if cfg.ConfigID == "" {
		return errors.New("configId is required")
	}

	if err := verifySelector(env, cfg.ChainSelector); err != nil {
		return err
	}

	qualifier := cfg.Qualifier()
	if env.DataStore != nil {
		existingAddresses, err := env.DataStore.Addresses().Fetch()
		if err != nil {
			return err
		}
		for _, addr := range existingAddresses {
			if addr.Qualifier == qualifier {
				return fmt.Errorf("mcms with qualifier %s already exists in datastore, must be unique", qualifier)
			}
		}
	}

	return nil
}

func (DeploySolanaMCMS) Apply(env cldf.Environment, cfg DeploySolanaMCMSConfig) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	if cfg.BuildConfig != nil {
		env.Logger.Debugw("Building solana MCMS artifacts", "buildConfig", cfg.BuildConfig)
		if err := BuildMCMSPrograms(env, *cfg.BuildConfig); err != nil {
			return out, fmt.Errorf("failed to build solana MCMS artifacts: %w", err)
		}
	} else {
		env.Logger.Debugw("Skipping solana MCMS build; no buildConfig provided")
	}

	chain, ok := env.BlockChains.SolanaChains()[cfg.ChainSelector]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", cfg.ChainSelector)
	}

	mcmsCfg := *cfg.MCMSWithTimelockConfig
	qualifier := cfg.Qualifier()
	mcmsCfg.Qualifier = &qualifier

	outDS := datastore.NewMemoryDataStore()
	if _, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(env, outDS, chain, mcmsCfg); err != nil {
		return out, fmt.Errorf("failed to deploy MCMS with timelock on solana: %w", err)
	}

	if err := tagMCMSAddresses(outDS, cfg); err != nil {
		return out, err
	}

	out.DataStore = outDS
	return out, nil
}

func tagMCMSAddresses(ds datastore.MutableDataStore, cfg DeploySolanaMCMSConfig) error {
	qualifier := cfg.Qualifier()
	refs, err := ds.Addresses().Fetch()
	if err != nil {
		return fmt.Errorf("failed to fetch deployed MCMS addresses: %w", err)
	}
	for _, addr := range refs {
		if addr.ChainSelector != cfg.ChainSelector {
			continue
		}
		addr.Qualifier = qualifier
		addr.Labels.Add("mcms_config=" + cfg.ConfigID)
		for k, v := range cfg.Labels {
			addr.Labels.Add(fmt.Sprintf("%s=%s", k, v))
		}
		if err := ds.Addresses().Upsert(addr); err != nil {
			return fmt.Errorf("failed to upsert address %s: %w", addr.Qualifier, err)
		}
	}
	return nil
}

func verifySelector(env cldf.Environment, selector uint64) error {
	if selector == 0 {
		return errors.New("chainSelector is required")
	}
	family, err := chain_selectors.GetSelectorFamily(selector)
	if err != nil {
		return err
	}
	if family != chain_selectors.FamilySolana {
		return fmt.Errorf("chain selector %d is not a solana chain", selector)
	}
	if _, ok := env.BlockChains.SolanaChains()[selector]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", selector)
	}

	return nil
}
