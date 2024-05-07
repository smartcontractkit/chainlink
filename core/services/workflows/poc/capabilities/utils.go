package capabilities

import (
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

const LocalCodeActionCapability = "Local::Code::Action"
const LocalCodeConsensusCapability = "Local::Code::Consensus"

func UnwrapValue[O any](a values.Value) (O, error) {
	var o O
	if reflect.TypeOf(o).Kind() == reflect.Ptr {
		o = reflect.New(reflect.TypeOf(o).Elem()).Interface().(O)
	}

	err := a.UnwrapTo(o)
	return o, err
}
