package customendpoint

import (
	"testing"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
)

// Test util to wait till in-progress transmissions are done.
func WaitForTransmitters(t *testing.T, transmitter types.ContractTransmitter) {
	tracker, ok := transmitter.(*contractTracker)
	if !ok {
		require.True(t, ok, "Unsuccessful cast to contractTracker")
	}
	tracker.transmittersWg.Wait()
}

func CreateConfigDigester(EndpointName string, EndpointTarget string, PayloadType string) types.OffchainConfigDigester {
	return offchainConfigDigester{
		EndpointName:   EndpointName,
		EndpointTarget: EndpointTarget,
		PayloadType:    PayloadType,
	}
}
