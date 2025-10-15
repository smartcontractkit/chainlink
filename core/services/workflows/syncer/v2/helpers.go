package v2

import (
	"context"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
)

var rlConfig = ratelimiter.Config{
	GlobalRPS:      1000.0,
	GlobalBurst:    1000,
	PerSenderRPS:   30.0,
	PerSenderBurst: 30,
}

var wlConfig = syncerlimiter.Config{
	Global:   200,
	PerOwner: 200,
}

var ErrCouldNotDecode = errors.New("failed to decode revert data")

type testEvtHandler struct {
	events []Event
	mux    sync.Mutex
	errFn  func() error
}

func (m *testEvtHandler) Close() error { return nil }

func (m *testEvtHandler) Handle(ctx context.Context, event Event) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.events = append(m.events, event)
	if m.errFn != nil {
		return m.errFn()
	}
	return nil
}

func (m *testEvtHandler) ClearEvents() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.events = make([]Event, 0)
}

func (m *testEvtHandler) GetEvents() []Event {
	m.mux.Lock()
	defer m.mux.Unlock()

	eventsCopy := make([]Event, len(m.events))
	copy(eventsCopy, m.events)

	return eventsCopy
}

func newTestEvtHandler(errFn func() error) *testEvtHandler {
	return &testEvtHandler{
		errFn:  errFn,
		events: make([]Event, 0),
	}
}

type testDonNotifier struct {
	don capabilities.DON
	err error
}

func (t *testDonNotifier) WaitForDon(ctx context.Context) (capabilities.DON, error) {
	return t.don, t.err
}

type mockService struct{}

func (m *mockService) Start(context.Context) error { return nil }

func (m *mockService) Close() error { return nil }

func (m *mockService) HealthReport() map[string]error { return map[string]error{"svc": nil} }

func (m *mockService) Ready() error { return nil }

func (m *mockService) Name() string { return "svc" }

func HandleRevertData(err error) (any, error) {
	var ec rpc.Error
	var ed rpc.DataError
	if errors.As(err, &ec) && errors.As(err, &ed) && ec.ErrorCode() == 3 {
		if eds, ok := ed.ErrorData().(string); ok {
			revertData, err := hexutil.Decode(eds)
			if err == nil {
				ss, err := workflow_registry_wrapper_v2.WorkflowRegistryMetaData.GetAbi()
				if err != nil {
					return nil, err
				}
				var firstFour [4]byte
				if len(revertData) >= 4 {
					copy(firstFour[:], revertData[0:4])
				}
				er, err := ss.ErrorByID(firstFour)
				if err != nil {
					return nil, err
				}
				interf, err := er.Unpack(revertData)
				if err != nil {
					return nil, err
				}
				return interf, nil
			}
		}
	}
	return nil, ErrCouldNotDecode
}
