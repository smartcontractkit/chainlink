package changeset

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type CsMCMSDeploy struct{}

var _ cldf.ChangeSetV2[DeployChangesetInput] = CsMCMSDeploy{}

// VerifyPreconditions checks if the input is valid before applying the changeset.
func (CsMCMSDeploy) VerifyPreconditions(env cldf.Environment, input DeployChangesetInput) error {
	if len(input.ChainSelectors) == 0 {
		return errors.New("no chain selectors provided")
	}
	if input.ConfigID == "" {
		return errors.New("configId is required")
	}
	// validate that there are no duplicates qualifier in the datastore
	//  this is a current limitation of the transfer ownership changeset which finds contracts by qualifier
	// find the duplicate qualifiers
	var err2 error
	m := make(map[string]struct{})
	for _, s := range input.ChainSelectors {
		q := input.Qualifier(s)
		if _, exists := m[q]; exists {
			err2 = errors.Join(err2, fmt.Errorf("duplicate qualifier found: %s", q))
		}
		m[q] = struct{}{}
	}
	if err2 != nil {
		return err2
	}

	// check against existing addresses in the datastore
	existingAddresses, err := env.DataStore.Addresses().Fetch()
	if err != nil {
		return err
	}
	for _, addr := range existingAddresses {
		if _, exists := m[addr.Qualifier]; exists {
			return fmt.Errorf("mcms with qualifier %s already exists in datastore, must be unique", addr.Qualifier)
		}
	}

	return nil
}

func (CsMCMSDeploy) Apply(env cldf.Environment, input DeployChangesetInput) (cldf.ChangesetOutput, error) {
	mcmsConfigPerChain := make(map[uint64]commontypes.MCMSWithTimelockConfigV2, len(input.ChainSelectors))
	// we set the qualifier per chain selector
	for _, s := range input.ChainSelectors {
		q := input.Qualifier(s)
		c := input.ContractConfiguration.Config
		c.Qualifier = &q
		mcmsConfigPerChain[s] = c
	}

	o, err := commonchangeset.DeployMCMSWithTimelockV2(env, mcmsConfigPerChain)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	// add labels
	d := o.DataStore
	a, err := d.Addresses().Fetch()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	for _, addr := range a {
		addr.Qualifier = input.Qualifier(addr.ChainSelector)
		addr.Labels.Add("mcms_config=" + input.ConfigID)
		for k, v := range input.Labels {
			addr.Labels.Add(fmt.Sprintf("%s=%s", k, v))
		}
		err := d.Addresses().Upsert(addr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to upsert address %s: %w", addr.Qualifier, err)
		}
	}
	// we don't use the AddressBook; limit scope to DataStore only
	return cldf.ChangesetOutput{Reports: o.Reports, DataStore: o.DataStore}, nil
}

// MCMSDeployChangesetInput is the input for the MCMS deployment changeset.
// It includes the environment, config ID, labels, and contract configuration.
// The contract configuration applies to all chain selectors.
type DeployChangesetInput struct {
	Environment string
	// ConfigID is the ID of the MCMS configuration in the config store
	ConfigID string `json:"configId" yaml:"configId"`
	// Descriptor is an optional descriptor to append to the computed qualifier URL as a query parameter, used to distinguish multiple MCMS deployments with the same config ID
	Descriptor *string `json:"descriptor,omitempty" yaml:"descriptor,omitempty"`
	// Labels are optional labels to add to the deployed MCMS contracts addresses in the datastore
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// MCMSContractConfiguration is the configuration for the MCMS contract to be deployed
	ContractConfiguration ContractConfiguration `json:"contractConfigurations" yaml:"contractConfigurations"`
	// ChainSelectors are the chain selectors where the MCMS contracts will be deployed
	ChainSelectors []uint64 `json:"chainSelectors" yaml:"chainSelectors"`
}

// Qualifier returns the qualifier for the given chain selector
// it is a valid URL in the form of:
// contract://<chain_selector>/mcmsv2?mcms-config=<config_id>&descriptor=<descriptor>
func (m DeployChangesetInput) Qualifier(selector uint64) string {
	// make qualifier in the form of a valid url
	u := &url.URL{
		Scheme: "contract",
		Host:   strconv.FormatUint(selector, 10),
		Path:   "mcmsv2",
	}

	// Add query parameters
	q := u.Query()
	q.Add("mcms-config", m.ConfigID)
	if m.Descriptor != nil {
		q.Add("descriptor", *m.Descriptor)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
