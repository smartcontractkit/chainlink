package environment

import (
	"errors"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

type Config struct {
	Blockchains        []*cre.WrappedBlockchainInput `toml:"blockchains" validate:"required"`
	NodeSets           []*ns.Input                   `toml:"nodesets" validate:"required"`
	JD                 *jd.Input                     `toml:"jd" validate:"required"`
	Infra              *infra.Input                  `toml:"infra" validate:"required"`
	Fake               *fake.Input                   `toml:"fake" validate:"required"`
	ExtraCapabilities  ExtraCapabilitiesConfig       `toml:"extra_capabilities"`
	CapabilitiesConfig CapabilitiesConfig            `toml:"capabilities_configs"`
	S3ProviderInput    *s3provider.Input             `toml:"s3provider"`
}

type CapabilitiesConfig struct {
	EVM map[string]map[string]any `toml:"evm"`
}

func (c Config) Validate() error {
	if c.JD.CSAEncryptionKey == "" {
		return errors.New("jd.csa_encryption_key must be provided")
	}
	return nil
}

type ExtraCapabilitiesConfig struct {
	CronCapabilityBinaryPath      string `toml:"cron_capability_binary_path"`
	EVMCapabilityBinaryPath       string `toml:"evm_capability_binary_path"`
	LogEventTriggerBinaryPath     string `toml:"log_event_trigger_binary_path"`
	ReadContractBinaryPath        string `toml:"read_contract_capability_binary_path"`
	HTTPTriggerBinaryPath         string `toml:"http_trigger_capability_binary_path"`
	HTTPActionBinaryPath          string `toml:"http_action_capability_binary_path"`
	ConsensusCapabilityBinaryPath string `toml:"consensus_capability_binary_path"`
}
