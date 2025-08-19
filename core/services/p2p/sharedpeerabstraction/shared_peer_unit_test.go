package sharedpeerabstraction_test

import (
	"testing"

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
	require.False(t, sharedPeer.HasOCRCapability())
	require.Nil(t, sharedPeer.GetPeerGroupFactory())
	require.Nil(t, sharedPeer.GetSingletonWrapper())
	require.Nil(t, sharedPeer.GetPeer())

	// Test signing without key (should fail)
	testData := []byte("test message")
	signature, err := sharedPeer.Sign(testData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "private key not available")
	require.Empty(t, signature)

	// Test health check without starting
	err = sharedPeer.Ready()
	require.Error(t, err)
	require.Contains(t, err.Error(), "no peer implementation available")

	healthReport := sharedPeer.HealthReport()
	require.NotEmpty(t, healthReport)
	require.Contains(t, healthReport, sharedPeer.Name())
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
	require.False(t, sharedPeerOCR.HasOCRCapability())
	require.False(t, sharedPeerExternal.HasOCRCapability())
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
	require.NotEmpty(t, sharedPeer.Name())
	require.Contains(t, sharedPeer.Name(), "SharedPeerAbstraction")

	// Test that closing without starting is handled gracefully
	err := sharedPeer.Close()
	// The StateMachine implementation requires the service to be started before it can be closed
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot stop unstarted service")
}
