package keeper

import (
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
)

const observationSource = `
    encode_check_upkeep_tx   [type=ethabiencode
                              abi="checkUpkeep(uint256 id, address from)"
                              data="{\"id\":$(jobSpec.upkeepID),\"from\":$(jobSpec.fromAddress)}"]
    check_upkeep_tx          [type=ethcall
                              failEarly=true
                              extractRevertReason=true
                              evmChainID="$(jobSpec.evmChainID)"
                              contract="$(jobSpec.contractAddress)"
                              gas="$(jobSpec.checkUpkeepGasLimit)"
                              gasPrice="$(jobSpec.gasPrice)"
                              gasTipCap="$(jobSpec.gasTipCap)"
                              gasFeeCap="$(jobSpec.gasFeeCap)"
                              data="$(encode_check_upkeep_tx)"]
    decode_check_upkeep_tx   [type=ethabidecode
                              abi="bytes memory performData, uint256 maxLinkPayment, uint256 gasLimit, uint256 adjustedGasWei, uint256 linkEth"]
    encode_perform_upkeep_tx [type=ethabiencode
                              abi="performUpkeep(uint256 id, bytes calldata performData)"
                              data="{\"id\": $(jobSpec.upkeepID),\"performData\":$(decode_check_upkeep_tx.performData)}"]
    simulate_perform_upkeep_tx  [type=ethcall
                              extractRevertReason=true
                              evmChainID="$(jobSpec.evmChainID)"
                              contract="$(jobSpec.contractAddress)"
                              from="$(jobSpec.fromAddress)"
                              gas="$(jobSpec.performUpkeepGasLimit)"
                              data="$(encode_perform_upkeep_tx)"]
    decode_check_perform_tx  [type=ethabidecode
                              abi="bool success"]
    check_success            [type=conditional
                              failEarly=true
                              data="$(decode_check_perform_tx.success)"]
    perform_upkeep_tx        [type=ethtx
                              minConfirmations=0
                              to="$(jobSpec.contractAddress)"
                              from="[$(jobSpec.fromAddress)]"
                              evmChainID="$(jobSpec.evmChainID)"
                              data="$(encode_perform_upkeep_tx)"
                              gasLimit="$(jobSpec.performUpkeepGasLimit)"
                              txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
    encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> simulate_perform_upkeep_tx -> decode_check_perform_tx -> check_success -> perform_upkeep_tx
`

// ValidatedKeeperSpec analyses the tomlString passed as parameter and
// returns a newly-created Job if there are no validation errors inside the toml.
func ValidatedKeeperSpec(tomlString string) (job.Job, error) {
	// Create a new job with a randomly generated uuid, which can be replaced with the one from tomlString.
	var j = job.Job{
		ExternalJobID: uuid.NewV4(),
	}

	tree, err := toml.Load(tomlString)
	if err != nil {
		return j, err
	}

	if err := tree.Unmarshal(&j); err != nil {
		return j, err
	}

	var spec job.KeeperSpec
	if err := tree.Unmarshal(&spec); err != nil {
		return j, err
	}

	j.KeeperSpec = &spec

	if j.Type != job.Keeper {
		return j, errors.Errorf("unsupported type %s", j.Type)
	}

	if strings.Contains(tomlString, "observationSource") ||
		strings.Contains(tomlString, "ObservationSource") {
		return j, errors.New("There should be no 'observationSource' parameter included in the toml")
	}

	return j, nil
}
