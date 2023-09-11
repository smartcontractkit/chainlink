package loadfunctions

import (
	"github.com/smartcontractkit/wasp"
)

type TestMode int

const (
	ModeHTTPPayload TestMode = iota
	ModeSecretsOnlyPayload
	ModeReal
)

type SingleFunctionCallGun struct {
	ft               *FunctionsTest
	mode             TestMode
	times            uint32
	source           string
	slotID           uint8
	slotVersion      uint64
	encryptedSecrets []byte
	args             []string
	subscriptionId   uint64
	jobId            [32]byte
}

func NewSingleFunctionCallGun(
	ft *FunctionsTest,
	mode TestMode,
	times uint32,
	source string,
	slotID uint8,
	slotVersion uint64,
	args []string,
	subscriptionId uint64,
	jobId [32]byte,
) *SingleFunctionCallGun {
	return &SingleFunctionCallGun{
		ft:             ft,
		mode:           mode,
		times:          times,
		source:         source,
		slotID:         slotID,
		slotVersion:    slotVersion,
		args:           args,
		subscriptionId: subscriptionId,
		jobId:          jobId,
	}
}

func (m *SingleFunctionCallGun) callReal() *wasp.CallResult {
	err := m.ft.LoadTestClient.SendRequestWithDONHostedSecrets(
		m.times,
		m.source,
		m.slotID,
		m.slotVersion,
		m.args,
		m.subscriptionId,
		m.jobId,
	)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	return &wasp.CallResult{}
}

func (m *SingleFunctionCallGun) callWithSecrets() *wasp.CallResult {
	err := m.ft.LoadTestClient.SendRequestWithDONHostedSecrets(
		m.times,
		m.source,
		m.slotID,
		m.slotVersion,
		m.args,
		m.subscriptionId,
		m.jobId,
	)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	return &wasp.CallResult{}
}

func (m *SingleFunctionCallGun) callWithHttp() *wasp.CallResult {
	err := m.ft.LoadTestClient.SendRequest(
		m.times,
		m.source,
		[]byte{},
		m.args,
		m.subscriptionId,
		m.jobId,
	)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	return &wasp.CallResult{}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleFunctionCallGun) Call(_ *wasp.Generator) *wasp.CallResult {
	switch m.mode {
	case ModeSecretsOnlyPayload:
		return m.callWithSecrets()
	case ModeHTTPPayload:
		return m.callWithHttp()
	case ModeReal:
		return m.callReal()
	default:
		panic("test mode must be ModeSecretsOnlyPayload, ModeHTTPPayload or ModeReal")
	}
}
