package ccip

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
)

// Note if a breaking change is introduced to this struct nodes running different versions
// will not be able to unmarshal each other's observations. Do not modify unless you
// know what you are doing.
//
// IMPORTANT: Both CommitObservation and ExecutionObservation are streamed and processed by Atlas.
// Any change to that struct must be reflected in the Atlas codebase.
// Additionally, you must test if OTI telemetry ingestion works with the new struct on staging environment.
type CommitObservation struct {
	Interval                  cciptypes.CommitStoreInterval  `json:"interval"`
	TokenPricesUSD            map[cciptypes.Address]*big.Int `json:"tokensPerFeeCoin"`
	SourceGasPriceUSD         *big.Int                       `json:"sourceGasPrice"` // Deprecated
	SourceGasPriceUSDPerChain map[uint64]*big.Int            `json:"sourceGasPriceUSDPerChain"`
}

// Marshal MUST be used instead of raw json.Marshal(o) since it contains backwards compatibility related changes.
func (o CommitObservation) Marshal() ([]byte, error) {
	obsCopy := o

	// Similar to: commitObservationJSONBackComp but for commit observation marshaling.
	tokenPricesUSD := make(map[cciptypes.Address]*big.Int, len(obsCopy.TokenPricesUSD))
	for k, v := range obsCopy.TokenPricesUSD {
		tokenPricesUSD[cciptypes.Address(strings.ToLower(string(k)))] = v
	}
	obsCopy.TokenPricesUSD = tokenPricesUSD

	return json.Marshal(&obsCopy)
}

// ExecutionObservation stores messages as a map pointing from a sequence number (uint) to the message payload (MsgData)
// Having it structured this way is critical because:
// * it prevents having duplicated sequence numbers within a single ExecutionObservation (compared to the list representation)
// * prevents malicious actors from passing multiple messages with the same sequence number
// Note if a breaking change is introduced to this struct nodes running different versions
// will not be able to unmarshal each other's observations. Do not modify unless you
// know what you are doing.
//
// IMPORTANT: Both CommitObservation and ExecutionObservation are streamed and processed by Atlas.
// Any change to that struct must be reflected in the Atlas codebase.
// Additionally, you must test if OTI telemetry ingestion works with the new struct on staging environment.
type ExecutionObservation struct {
	Messages map[uint64]MsgData `json:"messages"`
}

type MsgData struct {
	TokenData [][]byte `json:"tokenData"`
}

// ObservedMessage is a transient struct used for processing convenience within the plugin. It's easier to process observed messages
// when all properties are flattened into a single structure.
// It should not be serialized and returned from types.ReportingPlugin functions, please serialize/deserialize to/from ExecutionObservation instead using NewObservedMessage
type ObservedMessage struct {
	SeqNr uint64
	MsgData
}

func NewExecutionObservation(observations []ObservedMessage) ExecutionObservation {
	denormalized := make(map[uint64]MsgData, len(observations))
	for _, o := range observations {
		denormalized[o.SeqNr] = MsgData{TokenData: o.TokenData}
	}
	return ExecutionObservation{Messages: denormalized}
}

func NewObservedMessage(seqNr uint64, tokenData [][]byte) ObservedMessage {
	return ObservedMessage{
		SeqNr:   seqNr,
		MsgData: MsgData{TokenData: tokenData},
	}
}

func (o ExecutionObservation) Marshal() ([]byte, error) {
	return json.Marshal(&o)
}

// GetParsableObservations checks the given observations for formatting and value errors.
// It returns all valid observations, potentially being an empty list. It will log
// malformed observations but never error.
//
// GetParsableObservations MUST be used instead of raw json.Unmarshal(o) since it contains backwards compatibility changes.
func GetParsableObservations[O CommitObservation | ExecutionObservation](l logger.Logger, observations []types.AttributedObservation) []O {
	var parseableObservations []O
	var observers []commontypes.OracleID
	for _, ao := range observations {
		if len(ao.Observation) == 0 {
			// Empty observation
			l.Infow("Discarded empty observation", "observer", ao.Observer)
			continue
		}
		var ob O
		var err error
		obsJSON := ao.Observation

		switch any(ob).(type) {
		case CommitObservation:
			commitObservation, err1 := commitObservationJSONBackComp(ao.Observation)
			if err1 != nil {
				l.Errorw("commit observation json backwards compatibility format failed", "err", err,
					"observation", string(ao.Observation), "observer", ao.Observer)
				continue
			}
			ob = any(commitObservation).(O)
		default:
			err = json.Unmarshal(obsJSON, &ob)
			if err != nil {
				l.Errorw("Received unmarshallable observation", "err", err, "observation", string(ao.Observation), "observer", ao.Observer)
				continue
			}
		}

		parseableObservations = append(parseableObservations, ob)
		observers = append(observers, ao.Observer)
	}
	l.Infow(
		"Parsed observations",
		"observers", observers,
		"observersLength", len(observers),
		"observationsLength", len(parseableObservations),
		"rawObservationLength", len(observations),
	)
	return parseableObservations
}

// For backwards compatibility, converts token prices to eip55.
// Prior to cciptypes.Address we were using go-ethereum common.Address type which is
// marshalled to lower-case while the string representation we used was eip55.
// Nodes that run different ccip version should generate the same observations.
func commitObservationJSONBackComp(obsJson []byte) (CommitObservation, error) {
	var obs CommitObservation
	err := json.Unmarshal(obsJson, &obs)
	if err != nil {
		return CommitObservation{}, fmt.Errorf("unmarshal observation: %w", err)
	}
	tokenPricesUSD := make(map[cciptypes.Address]*big.Int, len(obs.TokenPricesUSD))
	for k, v := range obs.TokenPricesUSD {
		tokenPricesUSD[ccipcalc.HexToAddress(string(k))] = v
	}
	obs.TokenPricesUSD = tokenPricesUSD
	return obs, nil
}
