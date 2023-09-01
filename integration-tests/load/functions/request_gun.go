package loadfunctions

import (
	"github.com/smartcontractkit/wasp"
)

/* SingleFunctionCallGun is a gun that constantly requests randomness for one feed  */

type SingleFunctionCallGun struct {
	ft                         *FunctionsTest
	times                      uint32
	source                     string
	encryptedSecretsReferences []byte
	args                       []string
	subscriptionId             uint64
	jobId                      [32]byte
}

func NewSingleFunctionCallGun(ft *FunctionsTest, times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) *SingleFunctionCallGun {
	return &SingleFunctionCallGun{
		ft:                         ft,
		times:                      times,
		source:                     source,
		encryptedSecretsReferences: encryptedSecretsReferences,
		args:                       args,
		subscriptionId:             subscriptionId,
		jobId:                      jobId,
	}
}

// Call implements example gun call, assertions on response bodies should be done here
func (m *SingleFunctionCallGun) Call(l *wasp.Generator) *wasp.CallResult {
	err := m.ft.LoadTestClient.SendRequest(
		m.times,
		m.source,
		m.encryptedSecretsReferences,
		m.args,
		m.subscriptionId,
		m.jobId,
	)
	if err != nil {
		return &wasp.CallResult{Error: err.Error(), Failed: true}
	}
	return &wasp.CallResult{}
}
