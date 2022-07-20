package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
)

//
// Return types:
//    map[string]interface{}{
//        "results": []interface{} containing any other type other pipeline tasks can return
//        "occurrences": (int64)
//    }
//
type ModeTask struct {
	BaseTask      `mapstructure:",squash"`
	Values        string `json:"values"`
	AllowedFaults string `json:"allowedFaults"`
}

var _ Task = (*ModeTask)(nil)

func (t *ModeTask) Type() TaskType {
	return TaskTypeMode
}

func (t *ModeTask) Run(_ context.Context, _ logger.Logger, vars Vars, inputs []Result) (result Result, runInfo RunInfo) {
	var (
		maybeAllowedFaults MaybeUint64Param
		valuesAndErrs      SliceParam
		allowedFaults      int
		faults             int
	)
	err := multierr.Combine(
		errors.Wrap(ResolveParam(&maybeAllowedFaults, From(t.AllowedFaults)), "allowedFaults"),
		errors.Wrap(ResolveParam(&valuesAndErrs, From(VarExpr(t.Values, vars), JSONWithVarExprs(t.Values, vars, true), Inputs(inputs))), "values"),
	)
	if err != nil {
		return Result{Error: err}, runInfo
	}

	if allowed, isSet := maybeAllowedFaults.Uint64(); isSet {
		allowedFaults = int(allowed)
	} else {
		allowedFaults = len(valuesAndErrs) - 1
	}

	values, faults := valuesAndErrs.FilterErrors()
	if faults > allowedFaults {
		return Result{Error: errors.Wrapf(ErrTooManyErrors, "Number of faulty inputs %v to mode task > number allowed faults %v", faults, allowedFaults)}, runInfo
	} else if len(values) == 0 {
		return Result{Error: errors.Wrap(ErrWrongInputCardinality, "values")}, runInfo
	}

	type entry struct {
		count    uint64
		original interface{}
	}

	var (
		m     = make(map[string]entry, len(values))
		max   uint64
		modes []interface{}
	)
	for _, val := range values {
		var comparable string
		switch v := val.(type) {
		case []byte:
			comparable = string(v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float64, float32,
			string, bool:
			comparable = fmt.Sprintf("%v", v)
		case *big.Int:
			comparable = v.String()
		case big.Int:
			comparable = v.String()
		case *decimal.Decimal:
			comparable = v.String()
		case decimal.Decimal:
			comparable = v.String()
		default:
			bs, err := json.Marshal(v)
			if err != nil {
				return Result{Error: errors.Wrapf(ErrBadInput, "could not json stringify value: %v", err)}, runInfo
			}
			comparable = string(bs)
		}

		m[comparable] = entry{
			count:    m[comparable].count + 1,
			original: val,
		}

		if m[comparable].count > max {
			modes = []interface{}{val}
			max = m[comparable].count
		} else if m[comparable].count == max {
			modes = append(modes, val)
		}
	}
	return Result{Value: map[string]interface{}{
		"results":     modes,
		"occurrences": max,
	}}, runInfo
}
