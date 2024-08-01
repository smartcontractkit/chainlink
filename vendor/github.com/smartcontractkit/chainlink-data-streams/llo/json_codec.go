package llo

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

var _ ReportCodec = JSONReportCodec{}

type JSONStreamValue struct {
	Type  LLOStreamValue_Type
	Value string
}

func UnmarshalJSONStreamValue(enc *JSONStreamValue) (StreamValue, error) {
	if enc == nil {
		return nil, ErrNilStreamValue
	}
	switch enc.Type {
	case LLOStreamValue_Decimal:
		sv := new(Decimal)
		if err := (sv).UnmarshalText([]byte(enc.Value)); err != nil {
			return nil, err
		}
		return sv, nil
	case LLOStreamValue_Quote:
		sv := new(Quote)
		if err := (sv).UnmarshalText([]byte(enc.Value)); err != nil {
			return nil, err
		}
		return sv, nil
	default:
		return nil, fmt.Errorf("unknown StreamValueType %d", enc.Type)
	}
}

// JSONReportCodec is a chain-agnostic reference implementation

type JSONReportCodec struct{}

func (cdc JSONReportCodec) Encode(r Report, _ llotypes.ChannelDefinition) ([]byte, error) {
	type encode struct {
		ConfigDigest                types.ConfigDigest
		SeqNr                       uint64
		ChannelID                   llotypes.ChannelID
		ValidAfterSeconds           uint32
		ObservationTimestampSeconds uint32
		Values                      []JSONStreamValue
		Specimen                    bool
	}
	values := make([]JSONStreamValue, len(r.Values))
	for i, sv := range r.Values {
		b, err := sv.MarshalText()
		if err != nil {
			return nil, fmt.Errorf("failed to encode StreamValue: %w", err)
		}
		values[i] = JSONStreamValue{
			Type:  sv.Type(),
			Value: string(b),
		}
	}
	e := encode{
		ConfigDigest:                r.ConfigDigest,
		SeqNr:                       r.SeqNr,
		ChannelID:                   r.ChannelID,
		ValidAfterSeconds:           r.ValidAfterSeconds,
		ObservationTimestampSeconds: r.ObservationTimestampSeconds,
		Values:                      values,
		Specimen:                    r.Specimen,
	}
	return json.Marshal(e)
}

func (cdc JSONReportCodec) Decode(b []byte) (r Report, err error) {
	type decode struct {
		ConfigDigest                string
		SeqNr                       uint64
		ChannelID                   llotypes.ChannelID
		ValidAfterSeconds           uint32
		ObservationTimestampSeconds uint32
		Values                      []JSONStreamValue
		Specimen                    bool
	}
	d := decode{}
	err = json.Unmarshal(b, &d)
	if err != nil {
		return r, fmt.Errorf("failed to decode report: expected JSON (got: %s); %w", b, err)
	}
	cdBytes, err := hex.DecodeString(d.ConfigDigest)
	if err != nil {
		return r, fmt.Errorf("invalid ConfigDigest; %w", err)
	}
	cd, err := types.BytesToConfigDigest(cdBytes)
	if err != nil {
		return r, fmt.Errorf("invalid ConfigDigest; %w", err)
	}
	values := make([]StreamValue, len(d.Values))
	for i := range d.Values {
		values[i], err = UnmarshalJSONStreamValue(&d.Values[i])
		if err != nil {
			return r, fmt.Errorf("failed to decode StreamValue: %w", err)
		}
	}

	return Report{
		ConfigDigest:                cd,
		SeqNr:                       d.SeqNr,
		ChannelID:                   d.ChannelID,
		ValidAfterSeconds:           d.ValidAfterSeconds,
		ObservationTimestampSeconds: d.ObservationTimestampSeconds,
		Values:                      values,
		Specimen:                    d.Specimen,
	}, err
}
