package sharedpeerabstraction_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p/sharedpeerabstraction"
)

type mockOCRConfigSimple struct{}

func (m *mockOCRConfigSimple) TraceLogging() bool {
	return false
}

func TestSharedPeerAbstraction_Creation(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Test creating shared peer abstraction with minimal config
	sharedPeer := sharedpeerabstraction.NewSharedPeerAbstraction(sharedpeerabstraction.SharedPeerConfig{
		KeyStore:      nil,
		KeystoreP2P:   nil,
		P2PConfig:     nil,
		OCRConfig:     &mockOCRConfigSimple{},
		DataSource:    nil,
		Logger:        lggr,
		EnableOCR:     false,
		ForceExternal: false,
	})

	require.NotNil(t, sharedPeer)

	// Test initial state
	assert.False(t, sharedPeer.HasOCRCapability())
	assert.Nil(t, sharedPeer.GetPeerGroupFactory())
	assert.Nil(t, sharedPeer.GetSingletonWrapper())
	assert.Nil(t, sharedPeer.GetPeer())

	// Test signing without key (should fail)
	testData := []byte("test message")
	signature, err := sharedPeer.Sign(testData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private key not available")
	assert.Empty(t, signature)

	// Test health check without starting
	err = sharedPeer.Ready()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no peer implementation available")

	healthReport := sharedPeer.HealthReport()
	assert.NotEmpty(t, healthReport)
	assert.Contains(t, healthReport, sharedPeer.Name())
}

func TestSharedPeerAbstraction_Configuration(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Test OCR enabled configuration
	sharedPeerOCR := sharedpeerabstraction.NewSharedPeerAbstraction(sharedpeerabstraction.SharedPeerConfig{
		KeyStore:      nil,
		KeystoreP2P:   nil,
		P2PConfig:     nil,
		OCRConfig:     &mockOCRConfigSimple{},
		DataSource:    nil,
		Logger:        lggr,
		EnableOCR:     true,
		ForceExternal: false,
	})

	require.NotNil(t, sharedPeerOCR)

	// Test force external configuration
	sharedPeerExternal := sharedpeerabstraction.NewSharedPeerAbstraction(sharedpeerabstraction.SharedPeerConfig{
		KeyStore:      nil,
		KeystoreP2P:   nil,
		P2PConfig:     nil,
		OCRConfig:     &mockOCRConfigSimple{},
		DataSource:    nil,
		Logger:        lggr,
		EnableOCR:     true, // OCR enabled but...
		ForceExternal: true, // ...forced to use external
	})

	require.NotNil(t, sharedPeerExternal)

	// Both should start in same state without being started
	assert.False(t, sharedPeerOCR.HasOCRCapability())
	assert.False(t, sharedPeerExternal.HasOCRCapability())
}

func TestSharedPeerAbstraction_ServiceLifecycle(t *testing.T) {
	lggr := logger.TestLogger(t)

	sharedPeer := sharedpeerabstraction.NewSharedPeerAbstraction(sharedpeerabstraction.SharedPeerConfig{
		KeyStore:      nil,
		KeystoreP2P:   nil,
		P2PConfig:     nil,
		OCRConfig:     &mockOCRConfigSimple{},
		DataSource:    nil,
		Logger:        lggr,
		EnableOCR:     false,
		ForceExternal: false,
	})

	require.NotNil(t, sharedPeer)

	// Test name
	assert.NotEmpty(t, sharedPeer.Name())
	assert.Contains(t, sharedPeer.Name(), "SharedPeerAbstraction")

	// Test that closing without starting is handled gracefully
	err := sharedPeer.Close()
	// The StateMachine implementation requires the service to be started before it can be closed
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot stop unstarted service")
}
