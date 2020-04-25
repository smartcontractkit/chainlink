// package triggerfns contains logic for triggering a fluxmonitor report
// according to arbitrary rules specified in this package.
//
// Rules should be added to triggerFnFactories using
// RegisterTriggerFunctionFactory. A natural place to put such factories is the
// fluxmonitor triggerfns package.
//
// The factoryFunction should take a parameter object, and use that to construct
// the trigger function it returns. The parameter object is parsed from the json
// value associated to "name" in the jobspec params object.
//
// Triggering is the main method on a TriggerFn. It is called to check whether
// the new answer merits a fresh report. The Parameters and Factory methods are
// used during serialization to JSON of TriggerFns. Each entry will appear as if
// in a fluxmonitor initiator params valueTriggers object as a key-value pair:
// {tfn.Factory(): tfn.Parameters}.n
package triggerfns

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// triggerFnFactories maps the names of the trigger functions used in a JSON job
// spec with a fluxmonitor initiator to the corresponding factory functions. New
// threshold functions should be added here.
var triggerFnFactories = map[string]triggerFnFactory{}

// RegisterTriggerFunctionFactory adds the given trigger function factory to the
// register, under the given name. The name is the key to be used to refer to
// this factory in the "valueTriggers" object of a fluxmonitor jobspec's
// initiator's params object.
func RegisterTriggerFunctionFactory(name string, factory triggerFnFactory) {
	triggerFnFactories[name] = factory
}

// A triggerFnFactory returns a trigger function based on its name and the other
// params. The factory is responsible for returning an error if the params are
// invalid in some way.
type triggerFnFactory func(name string, params interface{}) (TriggerFn, error)

// TriggerFn represents a trigger function used by fluxmonitor initiator to
// determine whether to report a recently observed feed value onchain.
type TriggerFn interface {
	// Triggering returns true if the deviation between the onchain and recently
	// observed values implies that the new value should be reported to the
	// fluxAggregator contract.
	Triggering(onchain, recent decimal.Decimal, extraData ...interface{}) (bool, error)
	// Parameters returns the parameters passed to the factory to create this
	// trigger, in a form which can be marshaled to JSON.
	Parameters() interface{}
	// Factory returns the name of the factory function which created this trigger
	Factory() string
}

// TriggerFns is a collection of ThresholdFn's with convenient serialization.
// Despite being represented as alist, the ordering of the functions should not
// be relied upon, as it could change upon (de)serialization.
type TriggerFns []TriggerFn

var ( // interface assertions
	_ driver.Valuer    = TriggerFns{}
	_ sql.Scanner      = TriggerFns{}
	_ json.Unmarshaler = &TriggerFns{}
	_ json.Marshaler   = TriggerFns{}
)

// AllTriggered returns true if all trigger functions in t trigger given the change
// between the onchain value and the recently observed value from the feed
func (t TriggerFns) AllTriggered(onchain, recent decimal.Decimal) (bool, error) {
	if len(t) == 0 {
		return false, errors.Errorf("fluxmonitor must have at least one trigger function")
	}
	trigger := true
	for _, tfn := range t {
		doesTrigger, err := tfn.Triggering(onchain, recent)
		fmt.Printf("triggering predicate %s %+v, %s -> %s %s\n",
			tfn.Factory(), tfn.Parameters(), onchain, recent, doesTrigger)
		if err != nil {
			return false, errors.Wrapf(err, "could not determine whether feed change "+
				"%s -> %s merits onchain report according to %s trigger function with "+
				"parameters %+v", onchain, recent, tfn.Factory(), tfn.Parameters())
		}
		trigger = trigger && doesTrigger
	}
	return trigger, nil
}

// getParameters extracts the trigger function parameters from a valueTriggers
// object, as a map from factory function names to trigger parameters
func getParameters(b []byte) (map[string]interface{}, error) {
	rawParameters := make(map[string]interface{})
	if err := json.Unmarshal(b, &rawParameters); err != nil {
		return nil, errors.Wrapf(err, `while parsing "%s" as trigger function`,
			string(b))
	}
	return rawParameters, nil
}

// makeTriggerFn returns the trigger function from the given factory, with the
// given trigger parameters, or an error
func makeTriggerFn(triggerFunctionName string, params interface{}) (TriggerFn, error) {
	triggerFnFactory, ok := triggerFnFactories[triggerFunctionName]
	if !ok {
		return nil, errors.Errorf(`trigger function "%s" uknown`,
			triggerFunctionName)
	}
	triggerFn, err := triggerFnFactory(triggerFunctionName, params)
	if err != nil {
		return nil, errors.Wrapf(err,
			`while deserializing trigger function "%s" from parameters %+v`,
			triggerFunctionName, params)
	}
	return triggerFn, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (f *TriggerFns) UnmarshalJSON(b []byte) error {
	// Length could be zero, since not all initiators require valueTriggers param
	rawParameters, err := getParameters(b)
	if err != nil {
		return err
	}
	for triggerFunctionName, params := range rawParameters {
		triggerFn, err := makeTriggerFn(triggerFunctionName, params)
		if err != nil {
			return err
		}
		*f = append(*f, triggerFn)
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface
func (f TriggerFns) MarshalJSON() ([]byte, error) {
	// The length of f could be zero, here, since the initiator might not be a fluxmonitor
	rv := []string{"{"}
	for _, tfn := range f {
		params, err := json.Marshal(tfn.Parameters())
		if err != nil {
			return nil, errors.Wrapf(err,
				"could not marshal parameters %+v for trigger function %s",
				tfn.Parameters(), tfn.Factory())
		}
		rv = append(rv, fmt.Sprintf(`"%s":%s`, tfn.Factory(), string(params)), ",")
	}
	rv = append(rv, "}")
	return []byte(strings.Join(rv, "")), nil

}

// Scan deserializes the trigger functions from the JSON in value, and appends
// them to f, which must be empty when Scan is called
func (f TriggerFns) Scan(value interface{}) error {
	var raw []byte
	switch v := value.(type) {
	case string:
		raw = []byte(v)
	case []byte:
		raw = v
	}
	if err := f.UnmarshalJSON(raw); err != nil {
		return err
	}
	return nil
}

// Value implements the driver.Valuer interface
func (f TriggerFns) Value() (driver.Value, error) {
	return f.MarshalJSON()
}
