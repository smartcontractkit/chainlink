package keeper

import (
	"fmt"
	"reflect"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
)

const (
	// ExpectedObservationSource this is the expected observation source of the keeper job.
	ExpectedObservationSource = `
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
perform_upkeep_tx        [type=ethtx
                          minConfirmations=0
                          to="$(jobSpec.contractAddress)"
                          from="[$(jobSpec.fromAddress)]"
                          evmChainID="$(jobSpec.evmChainID)"
                          data="$(encode_perform_upkeep_tx)"
                          gasLimit="$(jobSpec.performUpkeepGasLimit)"
                          txMeta="{\"jobID\":$(jobSpec.jobID),\"upkeepID\":$(jobSpec.prettyID)}"]
encode_check_upkeep_tx -> check_upkeep_tx -> decode_check_upkeep_tx -> encode_perform_upkeep_tx -> perform_upkeep_tx`
)

// expectedPipeline it is basically parsed ExpectedObservationSource value
var expectedPipeline pipeline.Pipeline

func init() {
	pp, err := pipeline.Parse(ExpectedObservationSource)
	if err != nil {
		panic(fmt.Sprintf("failed to parse default observation source: %v", err))
	}
	expectedPipeline = *pp
}

func ValidatedKeeperSpec(tomlString string) (job.Job, error) {
	var j = job.Job{
		ExternalJobID: uuid.NewV4(), // Default to generating a uuid, can be overwritten by the specified one in tomlString.
	}

	var spec job.KeeperSpec
	tree, err := toml.Load(tomlString)
	if err != nil {
		return j, err
	}

	if err := tree.Unmarshal(&j); err != nil {
		return j, err
	}

	if err := tree.Unmarshal(&spec); err != nil {
		return j, err
	}
	j.KeeperSpec = &spec

	if j.Type != job.Keeper {
		return j, errors.Errorf("unsupported type %s", j.Type)
	}

	if !reflect.DeepEqual(j.Pipeline.Tasks, expectedPipeline.Tasks) {
		return j, errors.New("invalid observation source provided")
	}

	return j, nil
}
