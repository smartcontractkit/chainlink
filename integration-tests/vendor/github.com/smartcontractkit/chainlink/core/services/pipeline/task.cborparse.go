package pipeline

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/cbor"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// Return types:
//
//	map[string]interface{} with potential value types:
//	    float64
//	    string
//	    bool
//	    map[string]interface{}
//	    []interface{}
//	    nil
type CBORParseTask struct {
	BaseTask `mapstructure:",squash"`
	Data     string `json:"data"`
	Mode     string `json:"mode"`
}

var _ Task = (*CBORParseTask)(nil)

func (t *CBORParseTask) Type() TaskType {
	return TaskTypeCBORParse
}

func (t *CBORParseTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	_, err := CheckInputs(inputs, -1, -1, 0)
	if err != nil {
		return Result{Error: errors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		data BytesParam
		mode StringParam
	)
	err = multierr.Combine(
		errors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars))), "data"),
		errors.Wrap(ResolveParam(&mode, From(NonemptyString(t.Mode), "diet")), "mode"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	switch mode {
	case "diet":
		// NOTE: In diet mode, cbor_parse ASSUMES that the incoming CBOR is a
		// map. In the case that data is entirely missing, we assume it was the
		// empty map
		parsed, err := cbor.ParseDietCBOR(data)
		if err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "CBORParse: data: %v", err)}, runInfo
		}
		return Result{Value: parsed}, runInfo
	case "standard":
		parsed, err := cbor.ParseStandardCBOR(data)
		if err != nil {
			return Result{Error: errors.Wrapf(ErrBadInput, "CBORParse: data: %v", err)}, runInfo
		}
		return Result{Value: parsed}, runInfo
	default:
		return Result{Error: errors.Errorf("unrecognised mode: %s", mode)}, runInfo
	}
}
