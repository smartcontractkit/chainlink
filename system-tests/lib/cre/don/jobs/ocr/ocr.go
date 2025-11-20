package ocr

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
)

func BootstrapJobSpec(nodeID string, name string, ocr3CapabilityAddress string, chainID uint64) *jobv1.ProposeJobRequest {
	uuid := uuid.NewString()

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec: fmt.Sprintf(`
	type = "bootstrap"
	schemaVersion = 1
	externalJobID = "%s"
	name = "%s"
	contractID = "%s"
	contractConfigTrackerPollInterval = "1s"
	contractConfigConfirmations = 1
	relay = "evm"
	[relayConfig]
	chainID = %d
	providerType = "ocr3-capability"
`,
			uuid,
			"ocr3-bootstrap-"+name+fmt.Sprintf("-%d", chainID),
			ocr3CapabilityAddress,
			chainID),
	}
}

// ConfigMerger merges default config with overrides (either on DON or chain level)
type ConfigMerger func(flag cre.CapabilityFlag, nodeSet cre.NodeSetWithCapabilityConfigs, chainIDUint64 uint64, capabilityConfig cre.CapabilityConfig) (map[string]any, bool, error)

// JobConfigGenerator constains the logic that generates the job-specific part of the job spec
type JobConfigGenerator = func(logger zerolog.Logger, chainID uint64, nodeAddress string, mergedConfig map[string]any) (string, error)

// CapabilityEnabler determines if a capability is enabled for a given DON
type CapabilityEnabler func(don *cre.Don, flag cre.CapabilityFlag) bool

// EnabledChainsProvider provides the list of enabled chains for a given capability
type EnabledChainsProvider func(registryChainSelector uint64, nodeSet cre.NodeSetWithCapabilityConfigs, flag cre.CapabilityFlag) ([]uint64, error)

// ContractNamer is a function that returns the name of the OCR3 contract  used in the datastore
type ContractNamer func(chainID uint64) string

type DataStoreOCR3ContractKeyProvider func(contractName string, chainSelector uint64) datastore.AddressRefKey
