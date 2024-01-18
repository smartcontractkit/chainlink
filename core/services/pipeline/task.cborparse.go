package pipeline

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/cbor"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
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
		return Result{Error: pkgerrors.Wrap(err, "task inputs")}, runInfo
	}

	var (
		data BytesParam
		mode StringParam
	)
	err = errors.Join(
		pkgerrors.Wrap(ResolveParam(&data, From(VarExpr(t.Data, vars))), "data"),
		pkgerrors.Wrap(ResolveParam(&mode, From(NonemptyString(t.Mode), "diet")), "mode"),
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
			return Result{Error: pkgerrors.Wrapf(ErrBadInput, "CBORParse: data: %v", err)}, runInfo
		}
		return Result{Value: parsed}, runInfo
	case "standard":
		parsed, err := cbor.ParseStandardCBOR(data)
		if err != nil {
			return Result{Error: pkgerrors.Wrapf(ErrBadInput, "CBORParse: data: %v", err)}, runInfo
		}
		return Result{Value: parsed}, runInfo
	default:
		return Result{Error: pkgerrors.Errorf("unrecognised mode: %s", mode)}, runInfo
	}
}
