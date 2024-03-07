package ccip

import (
	"math/big"

	json2 "github.com/goccy/go-json"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
)

// Note if a breaking change is introduced to this struct nodes running different versions
// will not be able to unmarshal each other's observations. Do not modify unless you
// know what you are doing.
type CommitObservation struct {
	Interval          cciptypes.CommitStoreInterval  `json:"interval"`
	TokenPricesUSD    map[cciptypes.Address]*big.Int `json:"tokensPerFeeCoin"`
	SourceGasPriceUSD *big.Int                       `json:"sourceGasPrice"`
}

func (o CommitObservation) Marshal() ([]byte, error) {
	return json2.Marshal(&o)
}

// ExecutionObservation stores messages as a map pointing from a sequence number (uint) to the message payload (MsgData)
// Having it structured this way is critical because:
// * it prevents having duplicated sequence numbers within a single ExecutionObservation (compared to the list representation)
// * prevents malicious actors from passing multiple messages with the same sequence number
// Note if a breaking change is introduced to this struct nodes running different versions
// will not be able to unmarshal each other's observations. Do not modify unless you
// know what you are doing.
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
	return json2.Marshal(&o)
}

// GetParsableObservations checks the given observations for formatting and value errors.
// It returns all valid observations, potentially being an empty list. It will log
// malformed observations but never error.
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
		err := json2.Unmarshal(ao.Observation, &ob)
		if err != nil {
			l.Errorw("Received unmarshallable observation", "err", err, "observation", string(ao.Observation), "observer", ao.Observer)
			continue
		}
		parseableObservations = append(parseableObservations, ob)
		observers = append(observers, ao.Observer)
	}
	l.Infow(
		"Parsed observations",
		"observers", observers,
		"observersLength", len(observers),
		"observations", parseableObservations,
		"observationsLength", len(parseableObservations),
		"rawObservationLength", len(observations),
	)
	return parseableObservations
}
