package host

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bytecodealliance/wasmtime-go/v19"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/values/pb"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/poc/workflow"
)

const entryPoint = "_start"

func NewWasmGuestRunner(data []byte) (GuestRunner, error) {
	engine := wasmtime.NewEngine()
	module, err := wasmtime.NewModule(engine, data)
	if err != nil {
		return nil, err
	}

	entryFound := false
	memoryFound := false
	for _, export := range module.Exports() {
		if export.Name() == entryPoint {
			entryFound = true
			if export.Type().FuncType() == nil {
				return nil, fmt.Errorf("entry point %s is not a function", entryPoint)
			}

			if memoryFound {
				break
			}
		}

		if export.Name() == "memory" {
			if export.Type().MemoryType() == nil {
				return nil, fmt.Errorf("memory export is not a memory")
			}
			memoryFound = true
			if entryFound {
				break
			}
		}
	}

	if !entryFound {
		return nil, fmt.Errorf("entry point %s not entryFound in", entryPoint)
	}

	if !memoryFound {
		return nil, fmt.Errorf("memory not found in")
	}

	return &wasmGuestRunner{module: module, engine: engine, inputs: make(chan computeRequest, 100), exitCh: make(chan runResult, 100)}, nil
}

type wasmGuestRunner struct {
	module *wasmtime.Module
	engine *wasmtime.Engine
	store  *wasmtime.Store
	inputs chan computeRequest
	exitCh chan runResult
	retCh  chan runResult
}

func (w *wasmGuestRunner) Run() (*workflow.Spec, error) {
	w.store = wasmtime.NewStore(w.engine)
	linker := wasmtime.NewLinker(w.engine)
	wasiConfig := wasmtime.NewWasiConfig()
	// TODO remove, it's useful for debugging
	wasiConfig.InheritStdout()
	wasiConfig.InheritStderr()
	w.store.SetWasi(wasiConfig)
	if err := linker.DefineWasi(); err != nil {
		return nil, err
	}
	if err := linker.DefineFunc(w.store, "env", "step", w.step); err != nil {
		return nil, err
	}

	instance, err := linker.Instantiate(w.store, w.module)
	if err != nil {
		return nil, err
	}

	w.retCh = make(chan runResult)
	start := instance.GetFunc(w.store, "_start")
	w.exitCh = make(chan runResult, 10)
	go func() {
		r, err := realResultOrError(start.Call(w.store))
		if err != nil {
			w.exitCh <- runResult{err: err}
		}
		rv, err := values.Wrap(int64(r.(int32)))
		w.exitCh <- runResult{retVal: rv, err: err}
	}()

	select {
	case rawSpec := <-w.retCh:
		spec := &workflow.Spec{}
		if rawSpec.err != nil {
			return nil, rawSpec.err
		}
		if err = rawSpec.retVal.UnwrapTo(&spec.SpecBase); err != nil {
			return nil, err
		}

		spec.LocalExecutions = map[string]workflow.LocalCapability{}
		w.addWasmLocalCapabilities(spec, spec.Actions, commoncap.CapabilityTypeAction)
		w.addWasmLocalCapabilities(spec, spec.Consensus, commoncap.CapabilityTypeConsensus)

		return spec, nil
	case r := <-w.exitCh:
		return nil, fmt.Errorf("program exited before creating a spec, return code %v, error %w", r.retVal, r.err)
	}
}

func (w *wasmGuestRunner) step(caller *wasmtime.Caller, valueBuf, valueSize, retVal, retValMaxSize int64) int64 {
	valueRaw, err := safeMem(caller, valueBuf, valueSize, w.store)
	// TODO something better
	if err != nil {
		return -1
	}

	// TODO make a proto for this, don't just override map...
	wrappedValue := &pb.Value{}
	if err = proto.Unmarshal(valueRaw, wrappedValue); err != nil {
		return -2
	}
	value := wrappedValue.GetMapValue()
	fields := value.Fields
	result := runResult{retVal: values.FromProto(fields["result"]), cont: fields["continue"].GetBoolValue()}
	eVal := value.Fields["error"]
	if eVal != nil {
		result.err = errors.New(eVal.GetStringValue())
	}

	w.retCh <- result
	nextInput := w.getNextInput(retValMaxSize)
	return copyBuffer(caller, nextInput, retVal, retValMaxSize, w.store)
}

func safeMem(caller *wasmtime.Caller, ptr int64, size int64, store *wasmtime.Store) ([]byte, error) {
	mem := caller.GetExport("memory").Memory()
	data := mem.UnsafeData(store)
	if ptr+size > int64(len(data)) {
		return nil, errors.New("out of bounds memory access")
	}
	return data[ptr : ptr+size], nil
}

func (w *wasmGuestRunner) getNextInput(maxSize int64) []byte {
	for {
		req := <-w.inputs
		// TODO make a real proto for this, don't just re-use the map as a hack
		wrapped, err := values.NewMap(map[string]any{
			"request": req.input,
			"stepRef": req.stepRef,
		})
		if err != nil {
			req.retCh <- runResult{err: err}
		}

		reqBytes, err := proto.Marshal(values.Proto(wrapped))
		if err != nil {
			req.retCh <- runResult{err: err}
		}

		if len(reqBytes) > int(maxSize) {
			req.retCh <- runResult{err: fmt.Errorf("input too large: %d > %d", len(reqBytes), maxSize)}
		}

		w.retCh = req.retCh
		return reqBytes
	}
}

type computeRequest struct {
	stepRef string
	input   values.Value
	retCh   chan runResult
}

type runResult struct {
	retVal values.Value
	cont   bool
	err    error
}

func realResultOrError(rawResult any, err error) (any, error) {
	if err == nil {
		return rawResult, nil
	}
	var werr *wasmtime.Error
	ok := errors.As(err, &werr)
	if ok {
		if i, ok := werr.ExitStatus(); ok {
			return i, nil
		} else {
			return nil, err
		}
	}
	return nil, err
}

func copyBuffer(caller *wasmtime.Caller, src []byte, ptr int64, size int64, store *wasmtime.Store) int64 {
	mem := caller.GetExport("memory").Memory()
	rawData := mem.UnsafeData(store)
	if len(rawData) < int(ptr+size) {
		return -1
	}
	buffer := rawData[ptr : ptr+size]
	dataLen := int64(len(src))
	copy(buffer, src)
	return dataLen
}

func (w *wasmGuestRunner) addWasmLocalCapabilities(spec *workflow.Spec, steps []workflow.StepDefinition, ctype commoncap.CapabilityType) {
	for _, step := range steps {
		if strings.HasPrefix(step.TypeRef, capabilities.LocalCapabilityPrefix) {
			spec.LocalExecutions[step.Ref] = &wasmCapability{runner: w, capabilityType: ctype}
		}
	}
}
