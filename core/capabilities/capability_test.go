package capabilities

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type capabilityUnderTest struct {
	Validatable
	CapabilityInfoProvider
}

func Test_CapabilityInfo(t *testing.T) {

	capability, err := NewCapabilityInfoProvider(fmt.Stringer(), CapabilityTypeAction, "This is a mock capability that doesn't do anything.", "test")

	require.NoError(t, err)

	capabilityUnderTest := capabilityUnderTest{
		capability,
	}

	assert.Equal(t, mockCapabilityInfo, capability.Info())
}
