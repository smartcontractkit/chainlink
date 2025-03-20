package workflows

import (
	"encoding/hex"
	"fmt"
	"sort"
	"sync"

	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

const (
	MeteringReportSchema string = "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb/capabilities.proto"
	MeteringReportDomain string = "platform"
	MeteringReportEntity string = "MeteringReport"
)

type MeteringReportStepRef string

func (s MeteringReportStepRef) String() string {
	return string(s)
}

type MeteringSpendUnit string

func (s MeteringSpendUnit) String() string {
	return string(s)
}

func (s MeteringSpendUnit) DecimalToSpendValue(value decimal.Decimal) MeteringSpendValue {
	return MeteringSpendValue{value: value, roundingPlace: 18}
}

func (s MeteringSpendUnit) IntToSpendValue(value int64) MeteringSpendValue {
	return MeteringSpendValue{value: decimal.NewFromInt(value), roundingPlace: 18}
}

func (s MeteringSpendUnit) StringToSpendValue(value string) (MeteringSpendValue, error) {
	dec, err := decimal.NewFromString(value)
	if err != nil {
		return MeteringSpendValue{}, err
	}

	return MeteringSpendValue{value: dec, roundingPlace: 18}, nil
}

type MeteringSpendValue struct {
	value         decimal.Decimal
	roundingPlace uint8
}

func (v MeteringSpendValue) Add(value MeteringSpendValue) MeteringSpendValue {
	return MeteringSpendValue{
		value:         v.value.Add(value.value),
		roundingPlace: v.roundingPlace,
	}
}

func (v MeteringSpendValue) Div(value MeteringSpendValue) MeteringSpendValue {
	return MeteringSpendValue{
		value:         v.value.Div(value.value),
		roundingPlace: v.roundingPlace,
	}
}

func (v MeteringSpendValue) GreaterThan(value MeteringSpendValue) bool {
	return v.value.GreaterThan(value.value)
}

func (v MeteringSpendValue) String() string {
	return v.value.StringFixedBank(int32(v.roundingPlace))
}

type ProtoDetail struct {
	Schema string
	Domain string
	Entity string
}

type MeteringReportStep struct {
	Peer2PeerID string
	SpendUnit   MeteringSpendUnit
	SpendValue  MeteringSpendValue
}

type MeteringReport struct {
	mu    sync.RWMutex
	steps map[MeteringReportStepRef][]MeteringReportStep
}

func NewMeteringReport() *MeteringReport {
	return &MeteringReport{
		steps: make(map[MeteringReportStepRef][]MeteringReportStep),
	}
}

func (r *MeteringReport) MedianSpend() map[MeteringSpendUnit]MeteringSpendValue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := map[MeteringSpendUnit][]MeteringSpendValue{}
	medians := map[MeteringSpendUnit]MeteringSpendValue{}

	for _, nodeVals := range r.steps {
		for _, step := range nodeVals {
			vals, ok := values[step.SpendUnit]
			if !ok {
				vals = []MeteringSpendValue{}
			}

			values[step.SpendUnit] = append(vals, step.SpendValue)
		}
	}

	for unit, set := range values {
		sort.Slice(set, func(i, j int) bool {
			return set[j].GreaterThan(set[i])
		})

		if len(set)%2 > 0 {
			medians[unit] = set[len(set)/2]

			continue
		}

		medians[unit] = set[len(set)/2-1].Add(set[len(set)/2]).Div(unit.IntToSpendValue(2))
	}

	return medians
}

func (r *MeteringReport) SetStep(ref MeteringReportStepRef, step []MeteringReportStep) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.steps[ref] = step

	return nil
}

func (r *MeteringReport) Message(lggr logger.Logger) string {
	protoReport := &pb.MeteringReport{
		Steps: map[string]*pb.MeteringReportStep{},
	}

	for key, step := range r.steps {
		nodeDetail := make([]*pb.MeteringReportNodeDetail, len(step))

		for idx, nodeVal := range step {
			nodeDetail[idx] = &pb.MeteringReportNodeDetail{
				Peer_2PeerId: nodeVal.Peer2PeerID,
				SpendUnit:    nodeVal.SpendUnit.String(),
				SpendValue:   nodeVal.SpendValue.String(),
			}
		}
		protoReport.Steps[key.String()] = &pb.MeteringReportStep{
			Nodes: nodeDetail,
		}
	}

	encoded, err := proto.Marshal(protoReport)
	if err != nil {
		lggr.Error("failed to marshal metering report")

		return ""
	}

	// using hex encoding here to normalize the data sent over otel
	return hex.EncodeToString(encoded)
}

func (r *MeteringReport) Description() string {
	return fmt.Sprintf("schema: %s; domain: %s; entity: %s", MeteringReportSchema, MeteringReportDomain, MeteringReportEntity)
}
