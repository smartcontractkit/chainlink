package cltest

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func HaveLenAtLeast(count int) types.GomegaMatcher {
	return &HaveLenAtLeastMatcher{count}
}

type HaveLenAtLeastMatcher struct {
	Count int
}

func (matcher *HaveLenAtLeastMatcher) Match(actual interface{}) (success bool, err error) {
	length, ok := lengthOf(actual)
	if !ok {
		return false, fmt.Errorf("HaveLenAtLeast matcher expects a string/array/map/channel/slice.  Got:\n%s", format.Object(actual, 1))
	}

	return length <= matcher.Count, nil
}

func (matcher *HaveLenAtLeastMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto have length greater than %d", format.Object(actual, 1), matcher.Count)
}

func (matcher *HaveLenAtLeastMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nnot to be less than length %d", format.Object(actual, 1), matcher.Count)
}

func lengthOf(a interface{}) (int, bool) {
	if a == nil {
		return 0, false
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Map, reflect.Array, reflect.String, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(a).Len(), true
	default:
		return 0, false
	}
}
