package observability

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

type MethodCall struct {
	MethodName string
	Arguments  []interface{}
	Returns    []interface{}
}

// The class expected to override the observed methods.
const expectedWrapper = "core/services/ocr2/plugins/ccip/internal/observability.ObservedOnRampReader"

// TestOnRampObservedMethods tests that all methods of OnRampReader are observed by a wrapper.
// It uses the runtime to detect if the call stack contains the wrapper class.
func TestOnRampObservedMethods(t *testing.T) {
	// Methods not expected to be observed.
	// Add a method name here to exclude it from the test.
	excludedMethods := []string{
		"Address",
		"Close",
	}

	// Defines the overridden method calls to test.
	// Not defining a non-excluded method here will cause the test to fail with an explicit error.
	methodCalls := make(map[string]MethodCall)
	methodCalls["GetDynamicConfig"] = MethodCall{
		MethodName: "GetDynamicConfig",
		Arguments:  []interface{}{testutils.Context(t)},
		Returns:    []interface{}{cciptypes.OnRampDynamicConfig{}, nil},
	}
	methodCalls["GetSendRequestsBetweenSeqNums"] = MethodCall{
		MethodName: "GetSendRequestsBetweenSeqNums",
		Arguments:  []interface{}{testutils.Context(t), uint64(0), uint64(100), true},
		Returns:    []interface{}{nil, nil},
	}
	methodCalls["IsSourceChainHealthy"] = MethodCall{
		MethodName: "IsSourceChainHealthy",
		Arguments:  []interface{}{testutils.Context(t)},
		Returns:    []interface{}{false, nil},
	}
	methodCalls["IsSourceCursed"] = MethodCall{
		MethodName: "IsSourceCursed",
		Arguments:  []interface{}{testutils.Context(t)},
		Returns:    []interface{}{false, nil},
	}
	methodCalls["RouterAddress"] = MethodCall{
		MethodName: "RouterAddress",
		Arguments:  []interface{}{testutils.Context(t)},
		Returns:    []interface{}{cciptypes.Address("0x0"), nil},
	}
	methodCalls["SourcePriceRegistryAddress"] = MethodCall{
		MethodName: "SourcePriceRegistryAddress",
		Arguments:  []interface{}{testutils.Context(t)},
		Returns:    []interface{}{cciptypes.Address("0x0"), nil},
	}

	// Test each method defined in the embedded type.
	observed, reader := buildReader(t)
	observedType := reflect.TypeOf(observed)
	for i := 0; i < observedType.NumMethod(); i++ {
		method := observedType.Method(i)
		testMethod(t, method, methodCalls, excludedMethods, reader, observed)
	}
}

func testMethod(t *testing.T, method reflect.Method, methodCalls map[string]MethodCall, excludedMethods []string, reader *mocks.OnRampReader, observed ObservedOnRampReader) {
	t.Run("observability_wrapper_"+method.Name, func(t *testing.T) {
		// Skip excluded methods.
		for _, em := range excludedMethods {
			if method.Name == em {
				// Skipping ignore method (not an error).
				return
			}
		}

		// Retrieve method call from definition (fail if not present).
		mc := methodCalls[method.Name]
		if mc.MethodName == "" {
			assert.Fail(t, fmt.Sprintf("method %s not defined in methodCalls, please define it or exclude it.", method.Name))
			return
		}

		assertCallByWrapper(t, reader, mc)

		// Perform call on observed object.
		callParams := buildCallParams(mc)
		methodc := reflect.ValueOf(&observed).MethodByName(mc.MethodName)
		methodc.Call(callParams)
	})
}

// Set the mock to fail if not called by the wrapper.
func assertCallByWrapper(t *testing.T, reader *mocks.OnRampReader, mc MethodCall) {
	reader.On(mc.MethodName, mc.Arguments...).Maybe().Return(mc.Returns...).Run(func(args mock.Arguments) {
		var i = 0
		var pc uintptr
		var ok = true
		for ok {
			pc, _, _, ok = runtime.Caller(i)
			f := runtime.FuncForPC(pc)
			if strings.Contains(f.Name(), expectedWrapper) {
				// Found the expected wrapper in the call stack.
				return
			}
			i++
		}
		assert.Fail(t, fmt.Sprintf("method %s not observed by wrapper. Please implement the method or add it to the excluded list.", mc.MethodName))
	})
}

func buildCallParams(mc MethodCall) []reflect.Value {
	callParams := make([]reflect.Value, len(mc.Arguments))
	for i, arg := range mc.Arguments {
		callParams[i] = reflect.ValueOf(arg)
	}
	return callParams
}

// Build a mock reader and an observed wrapper to be used in the tests.
func buildReader(t *testing.T) (ObservedOnRampReader, *mocks.OnRampReader) {
	labels = []string{"evmChainID", "plugin", "reader", "function", "success"}
	ph := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "test_histogram",
	}, labels)
	pg := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test_gauge",
	}, labels)
	metric := metricDetails{
		interactionDuration: ph,
		resultSetSize:       pg,
		pluginName:          "test plugin",
		readerName:          "test reader",
		chainId:             1337,
	}
	reader := mocks.NewOnRampReader(t)
	observed := ObservedOnRampReader{reader, metric}
	return observed, reader
}
