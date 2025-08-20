package evm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/ocr"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	chainsel "github.com/smartcontractkit/chain-selectors"
)

var EVMJobSpecFactoryFn = func(logger zerolog.Logger, chainID uint64, config map[string]any,
	capabilitiesAwareNodeSets []*cre.CapabilitiesAwareNodeSet,
	infraInput infra.Input,
	evmBinaryPath string) cre.JobSpecFactoryFn {
	return func(input *cre.JobSpecFactoryInput) (cre.DonsToJobSpecs, error) {
		configGen := func(nodeAddress string) (string, error) {
			pollIntervalStr, ok := config["logTriggerPollInterval"]
			if !ok || pollIntervalStr == nil {
				pollIntervalStr = "1s"
			}
			pollInterval, err := time.ParseDuration(pollIntervalStr.(string))
			if err != nil {
				return "", errors.Wrap(err, fmt.Sprintf("failed to parse poll interval '%s'", pollIntervalStr))
			}

			receiverGasMinimumStr, ok := config["receiverGasMinimum"]
			if !ok || receiverGasMinimumStr == nil {
				receiverGasMinimumStr = "1"
			}
			receiverGasMinimum, err := strconv.ParseUint(receiverGasMinimumStr.(string), 10, 64)
			if err != nil {
				return "", errors.Wrap(err, fmt.Sprintf("invalid receiverGasMinimum '%s'", receiverGasMinimumStr))
			}

			cs, ok := chainsel.EvmChainIdToChainSelector()[chainID]
			if !ok {
				return "", fmt.Errorf("chain selector not found for chainID: %d", chainID)
			}

			creForwarderKey := datastore.NewAddressRefKey(
				cs,
				datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
				semver.MustParse("1.0.0"),
				"",
			)
			creForwarderAddress, err := input.CldEnvironment.DataStore.Addresses().Get(creForwarderKey)
			if err != nil {
				return "", errors.Wrap(err, "failed to get CRE Forwarder address")
			}

			logger.Debug().Msgf("Found CRE Forwarder contract on chain %d at %s", chainID, creForwarderAddress.Address)

			return fmt.Sprintf(
				`'{"chainId":%d,"network":"evm","logTriggerPollInterval":%d, "creForwarderAddress":"%s","receiverGasMinimum":%d,"nodeAddress":"%s"}'`,
				chainID,
				pollInterval.Nanoseconds(),
				creForwarderAddress.Address,
				receiverGasMinimum,
				nodeAddress,
			), nil
		}

		return ocr.GenerateJobSpecsForStandardCapabilityWithOCR(
			logger,
			input.DonTopology,
			input.CldEnvironment.DataStore,
			chainID,
			capabilitiesAwareNodeSets,
			infraInput,
			evmBinaryPath,
			contracts.GetCapabilityContractIdentifier(chainID),
			cre.EVMCapability,
			fmt.Sprintf("evm-capability-%d", chainID),
			configGen,
		)
	}
}
