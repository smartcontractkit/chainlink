package loadfunctions

import (
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
)

type TestMode int

const (
	ModeHTTPPayload TestMode = iota
	ModeSecretsOnlyPayload
	ModeReal
)

type SingleFunctionCallGun struct {
	ft             *FunctionsTest
	mode           TestMode
	times          uint32
	source         string
	slotID         uint8
	slotVersion    uint64
	args           []string
	subscriptionID uint64
	jobID          [32]byte
}

func NewSingleFunctionCallGun(
	ft *FunctionsTest,
	mode TestMode,
	times uint32,
	source string,
	slotID uint8,
	slotVersion uint64,
	args []string,
	subscriptionID uint64,
	jobID [32]byte,
) *SingleFunctionCallGun {
	return &SingleFunctionCallGun{
		ft:             ft,
		mode:           mode,
		times:          times,
		source:         source,
		slotID:         slotID,
		slotVersion:    slotVersion,
		args:           args,
		subscriptionID: subscriptionID,
		jobID:          jobID,
	}
}

func (m *SingleFunctionCallGun) callReal() *wasp.Response {
	err := m.ft.LoadTestClient.SendRequestWithDONHostedSecrets(
		m.times,
		m.source,
		m.slotID,
		m.slotVersion,
		m.args,
		m.subscriptionID,
		m.jobID,
	)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	return &wasp.Response{}
}

func (m *SingleFunctionCallGun) callWithSecrets() *wasp.Response {
	err := m.ft.LoadTestClient.SendRequestWithDONHostedSecrets(
		m.times,
		m.source,
		m.slotID,
		m.slotVersion,
		m.args,
		m.subscriptionID,
		m.jobID,
	)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	return &wasp.Response{}
}

func (m *SingleFunctionCallGun) callWithHTTP() *wasp.Response {
	err := m.ft.LoadTestClient.SendRequest(
		m.times,
		m.source,
		[]byte{},
		m.args,
		m.subscriptionID,
		m.jobID,
	)
	if err != nil {
		return &wasp.Response{Error: err.Error(), Failed: true}
	}
	return &wasp.Response{}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleFunctionCallGun) Call(_ *wasp.Generator) *wasp.Response {
	switch m.mode {
	case ModeSecretsOnlyPayload:
		return m.callWithSecrets()
	case ModeHTTPPayload:
		return m.callWithHTTP()
	case ModeReal:
		return m.callReal()
	default:
		panic("test mode must be ModeSecretsOnlyPayload, ModeHTTPPayload or ModeReal")
	}
}
