package sharedpeerabstraction

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

// ExampleUsage demonstrates how to use the SharedPeerAbstraction
func ExampleUsage(
	keyStore keystore.Master,
	keystoreP2P keystore.P2P,
	p2pConfig config.P2P,
	ocrConfig ocrcommon.PeerWrapperOCRConfig,
	ds sqlutil.DataSource,
	lggr logger.Logger,
) error {
	// Create a shared peer abstraction that can work with both OCR and non-OCR contexts
	sharedPeer := NewSharedPeerAbstraction(SharedPeerConfig{
		KeyStore:      keyStore,
		KeystoreP2P:   keystoreP2P,
		P2PConfig:     p2pConfig,
		OCRConfig:     ocrConfig,
		DataSource:    ds,
		Logger:        lggr,
		EnableOCR:     true,  // Enable OCR capabilities when available
		ForceExternal: false, // Allow automatic selection between OCR and external peer
	})

	// Start the abstraction
	ctx := context.Background()
	if err := sharedPeer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start shared peer: %w", err)
	}
	defer sharedPeer.Close()

	// Check what capabilities are available
	if sharedPeer.HasOCRCapability() {
		lggr.Info("OCR capabilities are available")

		// Get PeerGroupFactory for OCR operations
		peerGroupFactory := sharedPeer.GetPeerGroupFactory()
		if peerGroupFactory != nil {
			lggr.Info("PeerGroupFactory is available for OCR operations")
			// Use peerGroupFactory for creating OCR peer groups...
		}

		// Get the underlying SingletonPeerWrapper if needed
		wrapper := sharedPeer.GetSingletonWrapper()
		if wrapper != nil {
			lggr.Info("SingletonPeerWrapper is available")
			// Access OCR-specific functionality through the wrapper...
		}
	} else {
		lggr.Info("Using external peer mode")

		// Get the external peer for direct P2P communication
		peer := sharedPeer.GetPeer()
		if peer != nil {
			lggr.Infow("External peer is available", "peerID", peer.ID())
			// Use peer for direct P2P operations...
		}
	}

	// Use signing capability (available in both modes)
	testMessage := []byte("Hello, P2P world!")
	signature, err := sharedPeer.Sign(testMessage)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	lggr.Infow("Message signed successfully", "signatureLength", len(signature))

	return nil
}

// ExampleOCROnlyUsage shows how to force OCR mode
func ExampleOCROnlyUsage(
	keyStore keystore.Master,
	p2pConfig config.P2P,
	ocrConfig ocrcommon.PeerWrapperOCRConfig,
	ds sqlutil.DataSource,
	lggr logger.Logger,
) (SharedPeerAbstraction, error) {
	// Create abstraction that only uses OCR capabilities
	sharedPeer := NewSharedPeerAbstraction(SharedPeerConfig{
		KeyStore:      keyStore,
		KeystoreP2P:   keyStore.P2P(), // Use the P2P keystore from master
		P2PConfig:     p2pConfig,
		OCRConfig:     ocrConfig,
		DataSource:    ds,
		Logger:        lggr,
		EnableOCR:     true,
		ForceExternal: false,
	})

	ctx := context.Background()
	if err := sharedPeer.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start OCR peer: %w", err)
	}

	if !sharedPeer.HasOCRCapability() {
		sharedPeer.Close()
		return nil, fmt.Errorf("OCR capability not available")
	}

	return sharedPeer, nil
}

// ExampleExternalPeerOnlyUsage shows how to force external peer mode
func ExampleExternalPeerOnlyUsage(
	keystoreP2P keystore.P2P,
	p2pConfig config.P2P,
	ds sqlutil.DataSource,
	lggr logger.Logger,
) (SharedPeerAbstraction, error) {
	// Create abstraction that only uses external peer
	sharedPeer := NewSharedPeerAbstraction(SharedPeerConfig{
		KeyStore:      nil, // Not needed for external peer only
		KeystoreP2P:   keystoreP2P,
		P2PConfig:     p2pConfig,
		OCRConfig:     nil, // Not needed for external peer only
		DataSource:    ds,
		Logger:        lggr,
		EnableOCR:     false, // Disable OCR
		ForceExternal: false,
	})

	ctx := context.Background()
	if err := sharedPeer.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start external peer: %w", err)
	}

	peer := sharedPeer.GetPeer()
	if peer == nil {
		sharedPeer.Close()
		return nil, fmt.Errorf("external peer not available")
	}

	return sharedPeer, nil
}
