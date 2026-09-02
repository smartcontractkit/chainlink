package pipeline

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// Sample is one measurement from one source at one instant.
// It bundles the value with the metadata (source, weight, unit, timestamp) that
// downstream tasks consume. A source with multiple fields (e.g. mid, bid, ask)
// becomes multiple self-contained samples, not parallel arrays of metadata
// that can drift out of alignment.
type Sample struct {
	Source string          `json:"source"`         // grouping + weighting key
	Value  decimal.Decimal `json:"value"`          // the measured value
	TsMs   int64           `json:"ts_ms"`          // the source's own clock, milliseconds
	Weight decimal.Decimal `json:"weight"`         // nominal weight, scaled in place by staleness
	Unit   string          `json:"unit,omitempty"` // "" = already the common unit
	Tag    string          `json:"tag,omitempty"`  // free label, diagnostics only
}

// SampleParam parses a single Sample from various representations.
type SampleParam Sample

func (s *SampleParam) UnmarshalPipelineParam(val any) error {
	switch v := val.(type) {
	case nil:
		return errors.Wrap(ErrBadInput, "nil sample")
	case Sample:
		*s = SampleParam(v)
		return nil
	case *Sample:
		if v == nil {
			return errors.Wrap(ErrBadInput, "nil sample pointer")
		}
		*s = SampleParam(*v)
		return nil
	case map[string]any:
		return sampleFromMap(v, (*Sample)(s))
	case []byte:
		var m map[string]any
		if err := json.Unmarshal(v, &m); err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		return sampleFromMap(m, (*Sample)(s))
	case string:
		return s.UnmarshalPipelineParam([]byte(v))
	default:
		return errors.Wrapf(ErrBadInput, "expected sample, got %T", val)
	}
}

func sampleFromMap(m map[string]any, s *Sample) error {
	var err error
	if src, ok := m["source"].(string); ok {
		s.Source = src
	}

	val, ok := m["value"]
	if !ok {
		return errors.Wrap(ErrBadInput, "sample missing value")
	}
	val, err = normalizeJSONNumber(val)
	if err != nil {
		return err
	}
	s.Value, err = utils.ToDecimal(val)
	if err != nil {
		return errors.Wrapf(ErrBadInput, "sample value: %v", err)
	}

	ts, ok := m["ts_ms"]
	if !ok {
		return errors.Wrap(ErrBadInput, "sample missing ts_ms")
	}
	ts, err = normalizeJSONNumber(ts)
	if err != nil {
		return err
	}
	s.TsMs, err = toInt64(ts)
	if err != nil {
		return errors.Wrapf(ErrBadInput, "sample ts_ms: %v", err)
	}

	if w, ok := m["weight"]; ok {
		w, err = normalizeJSONNumber(w)
		if err != nil {
			return err
		}
		s.Weight, err = utils.ToDecimal(w)
		if err != nil {
			return errors.Wrapf(ErrBadInput, "sample weight: %v", err)
		}
	} else {
		s.Weight = decimal.NewFromInt(1)
	}

	if u, ok := m["unit"].(string); ok {
		s.Unit = u
	}
	if tag, ok := m["tag"].(string); ok {
		s.Tag = tag
	}
	return nil
}

func normalizeJSONNumber(v any) (any, error) {
	if n, ok := v.(json.Number); ok {
		return n.String(), nil
	}
	return v, nil
}

func toInt64(v any) (int64, error) {
	switch tv := v.(type) {
	case int64:
		return tv, nil
	case int:
		return int64(tv), nil
	case int32:
		return int64(tv), nil
	case uint64:
		return int64(tv), nil
	case uint:
		return int64(tv), nil
	case float64:
		return int64(tv), nil
	case float32:
		return int64(tv), nil
	case string:
		return strconv.ParseInt(tv, 10, 64)
	case decimal.Decimal:
		if !tv.IsInteger() {
			return 0, errors.Errorf("decimal is not an integer: %v", tv)
		}
		return tv.IntPart(), nil
	default:
		return 0, errors.Errorf("cannot convert %T to int64", v)
	}
}

// SampleSliceParam parses a slice of Samples, flattening one level of nested slices.
type SampleSliceParam []Sample

func (s *SampleSliceParam) UnmarshalPipelineParam(val any) error {
	switch v := val.(type) {
	case nil:
		*s = nil
		return nil
	case []Sample:
		*s = v
		return nil
	case Sample:
		*s = []Sample{v}
		return nil
	case *Sample:
		if v == nil {
			return errors.Wrap(ErrBadInput, "nil sample in slice")
		}
		*s = []Sample{*v}
		return nil
	case []any:
		var samples []Sample
		for _, x := range v {
			switch inner := x.(type) {
			case Sample:
				samples = append(samples, inner)
			case *Sample:
				if inner == nil {
					return errors.Wrap(ErrBadInput, "nil sample in slice")
				}
				samples = append(samples, *inner)
			case []Sample:
				samples = append(samples, inner...)
			case []any:
				var innerSlice SampleSliceParam
				if err := innerSlice.UnmarshalPipelineParam(inner); err != nil {
					return err
				}
				samples = append(samples, innerSlice...)
			case map[string]any:
				var sp SampleParam
				if err := sp.UnmarshalPipelineParam(inner); err != nil {
					return err
				}
				samples = append(samples, Sample(sp))
			default:
				return errors.Wrapf(ErrBadInput, "expected sample in slice, got %T", x)
			}
		}
		*s = samples
		return nil
	case map[string]any:
		var sp SampleParam
		if err := sp.UnmarshalPipelineParam(v); err != nil {
			return err
		}
		*s = []Sample{Sample(sp)}
		return nil
	case SliceParam:
		return s.UnmarshalPipelineParam([]any(v))
	case []byte:
		var arr []any
		if err := json.Unmarshal(v, &arr); err != nil {
			return errors.Wrap(ErrBadInput, err.Error())
		}
		return s.UnmarshalPipelineParam(arr)
	case string:
		return s.UnmarshalPipelineParam([]byte(v))
	default:
		return errors.Wrapf(ErrBadInput, "expected sample slice, got %T", val)
	}
}

// AnyParam accepts any value without conversion. It is used by SampleTask to
// receive raw upstream output (string, map, bytes) before JSON path extraction.
type AnyParam struct{ v any }

func (a *AnyParam) UnmarshalPipelineParam(val any) error {
	a.v = val
	return nil
}

// Experimental: SampleTask extracts a single value and timestamp from an upstream bridge/HTTP
// response and bundles them into a Sample with source, weight, and unit metadata.
//
// Input:  input (any) — raw upstream result (JSON string, map, or bytes)
// Output: Sample — one labelled measurement ready for normalize/staleness/aggregation
// Fails:  if valuePath or tsPath is missing, or the upstream input errored
// Fails:  if the JSON path does not resolve to a valid number
type SampleTask struct {
	BaseTask  `mapstructure:",squash"`
	Input     string `json:"input"`
	Source    string `json:"source"`
	Weight    string `json:"weight"`
	Unit      string `json:"unit"`
	Tag       string `json:"tag"`
	ValuePath string `json:"valuePath"`
	TsPath    string `json:"tsPath"`
	TsUnit    string `json:"tsUnit"`
	Decimals  string `json:"decimals"`
}

var _ Task = (*SampleTask)(nil)

func (t *SampleTask) Type() TaskType {
	return TaskTypeSample
}

func (t *SampleTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		input     AnyParam
		source    StringParam
		weight    DecimalParam
		unit      StringParam
		tag       StringParam
		valuePath JSONPathParam
		tsPath    JSONPathParam
		tsUnit    StringParam
		decimals  MaybeInt32Param
	)
	// Default separator for JSONPathParam is comma.
	valuePath = NewJSONPathParam("")
	tsPath = NewJSONPathParam("")

	err := stderrors.Join(
		errors.Wrap(ResolveParam(&input, From(VarExpr(t.Input, vars), Input(inputs, 0))), "input"),
		errors.Wrap(ResolveParam(&source, From(NonemptyString(t.Source))), "source"),
		errors.Wrap(ResolveParam(&weight, From(NonemptyString(t.Weight), "1")), "weight"),
		errors.Wrap(ResolveParam(&unit, From(t.Unit, "")), "unit"),
		errors.Wrap(ResolveParam(&tag, From(t.Tag, "")), "tag"),
		errors.Wrap(ResolveParam(&valuePath, From(NonemptyString(t.ValuePath))), "valuePath"),
		errors.Wrap(ResolveParam(&tsPath, From(NonemptyString(t.TsPath))), "tsPath"),
		errors.Wrap(ResolveParam(&tsUnit, From(NonemptyString(t.TsUnit), "ms")), "tsUnit"),
		errors.Wrap(ResolveParam(&decimals, From(VarExpr(t.Decimals, vars), t.Decimals)), "decimals"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	decoded, err := sampleInputToMap(input.v)
	if err != nil {
		return Result{Error: errors.Wrap(err, "input")}, runInfo
	}

	valueRaw, err := traverseJSONPath(decoded, valuePath, false)
	if err != nil {
		return Result{Error: errors.Wrap(err, "valuePath")}, runInfo
	}
	valueRaw, err = normalizeJSONNumber(valueRaw)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	value, err := utils.ToDecimal(valueRaw)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "value: %v", err)}, runInfo
	}

	if dec, isSet := decimals.Int32(); isSet && dec > 0 {
		value = value.Div(decimal.New(1, dec))
	}

	tsRaw, err := traverseJSONPath(decoded, tsPath, false)
	if err != nil {
		return Result{Error: errors.Wrap(err, "tsPath")}, runInfo
	}
	tsRaw, err = normalizeJSONNumber(tsRaw)
	if err != nil {
		return Result{Error: err}, runInfo
	}
	tsMs, err := toInt64(tsRaw)
	if err != nil {
		return Result{Error: errors.Wrapf(ErrBadInput, "timestamp: %v", err)}, runInfo
	}
	if strings.EqualFold(string(tsUnit), "s") {
		tsMs *= 1000
	}

	return Result{Value: Sample{
		Source: string(source),
		Value:  value,
		TsMs:   tsMs,
		Weight: weight.Decimal(),
		Unit:   string(unit),
		Tag:    string(tag),
	}}, runInfo
}

func sampleInputToMap(v any) (map[string]any, error) {
	switch x := v.(type) {
	case map[string]any:
		return x, nil
	case string:
		return sampleInputToMap([]byte(x))
	case []byte:
		var m map[string]any
		if err := json.Unmarshal(x, &m); err != nil {
			return nil, err
		}
		return m, nil
	default:
		return nil, errors.Errorf("expected JSON string or map, got %T", v)
	}
}
