package sdk

import (
	"fmt"
	"unsafe"

	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

func NewRunner() workflow.Runner {
	return &wasmRunner{}
}

type wasmRunner struct {
}

//go:wasmimport env step
func step(valueBuf, valueSize, retVal, retValMaxSize int64) int64

func (wasmRunner) Run(spec *workflow.Spec) error {
	result, err := values.NewMap(map[string]any{"continue": true, "result": spec.SpecBase})
	if err != nil {
		return err
	}

	for {
		inputBuf, err := proto.Marshal(values.Proto(result))
		if err != nil {
			return err
		}
		inputBufRaw := int64(uintptr(unsafe.Pointer(&inputBuf[0])))

		// TODO determine the right max size
		requestBuf := make([]byte, 1024)
		requestBufRaw := int64(uintptr(unsafe.Pointer(&requestBuf[0])))
		reqSize := step(inputBufRaw, int64(len(inputBuf)), requestBufRaw, int64(len(requestBuf)))
		if reqSize < 0 {
			fmt.Println("Too small buffer")
			return fmt.Errorf("buffer of size %d too small for request")
		}
		requestBuf = requestBuf[:reqSize]
		wrappedRequest := &pb.Value{}
		if err := proto.Unmarshal(requestBuf, wrappedRequest); err != nil {
			fmt.Printf("Error unmarshalling request: %v", err)
			return err
		}
		fields := wrappedRequest.GetMapValue().Fields
		fmt.Printf("fields: %v\n", fields)
		stepRef := fields["stepRef"].GetStringValue()
		request := fields["request"]

		callback, ok := spec.LocalExecutions[stepRef]
		if !ok {
			fmt.Printf("stepRef %s not found in localExecutions", stepRef)
			return fmt.Errorf("stepRef %s not found in localExecutions", stepRef)
		}
		rawResult, cont, err := callback.Run(stepRef, values.FromProto(request))
		// TODO use a real proto, don't just use values.Value
		wrappedResult := map[string]any{"continue": cont, "result": rawResult}
		if err != nil {
			wrappedResult["error"] = err.Error()
		}
		result, err = values.NewMap(wrappedResult)
		if err != nil {
			fmt.Printf("Error creating new map: %v", err)
			return err
		}
	}
}
