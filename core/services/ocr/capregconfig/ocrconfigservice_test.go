package capregconfig

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

// testPeerID returns a consistent peer ID for use in tests.
func testPeerID() ragetypes.PeerID {
	var peerID ragetypes.PeerID
	copy(peerID[:], []byte("test-peer-id-12345678901234"))
	return peerID
}

// testPeerIDProvider returns a peer ID provider function for tests.
func testPeerIDProvider() PeerIDProvider {
	return func() (ragetypes.PeerID, error) {
		return testPeerID(), nil
	}
}

func TestOCRConfigService_StartClose(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	require.NoError(t, svc.Close())
}

func TestOCRConfigService_Start_NilPeerIDProvider(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, nil, 1, "0x1234567890abcdef")

	ctx := t.Context()
	err := svc.Start(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "peerIDProvider function is required")
}

func TestOCRConfigService_OnNewRegistry(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Create a registry with OCR config.
	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           5,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	// Verify config was stored.
	_, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	assert.True(t, ok)
}

func TestOCRConfigService_GetConfigTracker(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Get tracker without any config - should work with legacy fallback.
	tracker, err := svc.GetConfigTracker("consensus@1.0.0", "__default__", nil)
	require.NoError(t, err)
	require.NotNil(t, tracker)

	// Without config or legacy, should return error.
	_, _, err = tracker.LatestConfigDetails(ctx)
	require.Error(t, err)
}

func TestOCRConfigService_GetConfigTracker_WithConfig(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Add config via registry update.
	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           5,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	// Get tracker - should return registry-based config.
	tracker, err := svc.GetConfigTracker("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey, nil)
	require.NoError(t, err)

	changedInBlock, digest, err := tracker.LatestConfigDetails(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), changedInBlock)
	assert.NotNil(t, digest)

	config, err := tracker.LatestConfig(ctx, changedInBlock)
	require.NoError(t, err)
	assert.Equal(t, uint64(5), config.ConfigCount)
	assert.Equal(t, uint8(1), config.F)
	assert.Len(t, config.Signers, 4)
	assert.Len(t, config.Transmitters, 4)
}

func TestOCRConfigService_GetConfigDigester(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Get digester without any config.
	digester, err := svc.GetConfigDigester("consensus@1.0.0", "__default__", nil)
	require.NoError(t, err)
	require.NotNil(t, digester)

	// Without config or legacy, should return error.
	_, err = digester.ConfigDigestPrefix(ctx)
	require.Error(t, err)
}

func TestOCRConfigService_GetConfigDigester_WithConfig(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Add config via registry update.
	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           5,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	// Get digester - should return registry-based prefix.
	digester, err := svc.GetConfigDigester("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey, nil)
	require.NoError(t, err)

	prefix, err := digester.ConfigDigestPrefix(ctx)
	require.NoError(t, err)
	assert.Equal(t, ocrtypes.ConfigDigestPrefixKeystoneOCR3Capability, prefix)
}

func TestOCRConfigService_ConfigChangeDetection(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Initial config.
	ocrConfig1 := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           1,
	}

	capConfig1 := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig1,
		},
	}

	configBytes1, err := proto.Marshal(capConfig1)
	require.NoError(t, err)

	don1 := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes1},
		},
	}
	don1.Members = []ragetypes.PeerID{testPeerID()}

	registry1 := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don1,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry1)
	require.NoError(t, err)

	cfg1, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)
	assert.Equal(t, uint64(1), cfg1.ContractConfig.ConfigCount)

	// Updated config with new config count.
	ocrConfig2 := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain_v2"),
		ConfigCount:           2,
	}

	capConfig2 := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig2,
		},
	}

	configBytes2, err := proto.Marshal(capConfig2)
	require.NoError(t, err)

	don2 := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes2},
		},
	}
	don2.Members = []ragetypes.PeerID{testPeerID()}

	registry2 := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don2,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry2)
	require.NoError(t, err)

	cfg2, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)
	assert.Equal(t, uint64(2), cfg2.ContractConfig.ConfigCount)
}

func TestOCRConfigService_TransmitterHexEncoding(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Use binary transmitter addresses (like Ethereum addresses).
	transmitter1 := []byte{0xde, 0xad, 0xbe, 0xef, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	transmitter2 := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}

	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{transmitter1, transmitter2, transmitter1, transmitter2},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           1,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	cfg, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)

	// Verify transmitters are hex-encoded.
	assert.Equal(t, ocrtypes.Account("0xDeadbEef00112233445566778899AABBccDDeeFf"), cfg.ContractConfig.Transmitters[0])
	assert.Equal(t, ocrtypes.Account("0x0102030405060708090a0B0c0d0e0f1011121314"), cfg.ContractConfig.Transmitters[1])
}

func TestOCRConfigService_ConfigDigestComputation(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           5,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	cfg, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)

	// Verify digest has correct prefix.
	assert.True(t, ocrtypes.ConfigDigestPrefixKeystoneOCR3Capability.IsPrefixOf(cfg.ContractConfig.ConfigDigest))

	// Verify digest is deterministic (same config produces same digest).
	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	cfg2, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)
	assert.Equal(t, cfg.ContractConfig.ConfigDigest, cfg2.ContractConfig.ConfigDigest)
}

func TestOCRConfigService_ConfigDigestUniqueness(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           1,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	// Register the same OCR config under two different capability IDs.
	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
			"consensus@2.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	require.NoError(t, svc.OnNewRegistry(ctx, registry))

	cfg1, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)

	cfg2, ok := svc.getConfig("consensus@2.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	require.True(t, ok)

	// Different capability IDs should produce different digests.
	assert.NotEqual(t, cfg1.ContractConfig.ConfigDigest, cfg2.ContractConfig.ConfigDigest)
}

func TestOCRConfigService_LegacyFallbackAfterRegistryReceived(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Create a mock legacy tracker.
	mockTracker := &mockConfigTracker{
		configCount:  10,
		configDigest: ocrtypes.ConfigDigest{0x01, 0x02},
	}

	// Get tracker before any registry update - should error (no fallback yet).
	tracker, err := svc.GetConfigTracker("consensus@1.0.0", "__default__", mockTracker)
	require.NoError(t, err)

	_, _, err = tracker.LatestConfigDetails(ctx)
	require.Error(t, err) // No registry received yet, no fallback.

	// Send an empty registry update (no config for this capability).
	registry := &registrysyncer.LocalRegistry{
		Logger:            lggr,
		IDsToDONs:         map[registrysyncer.DonID]registrysyncer.DON{},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	// Now legacy fallback should work.
	changedInBlock, digest, err := tracker.LatestConfigDetails(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(10), changedInBlock)
	assert.Equal(t, mockTracker.configDigest, digest)
}

type mockConfigTracker struct {
	configCount  uint64
	configDigest ocrtypes.ConfigDigest
}

func (m *mockConfigTracker) Notify() <-chan struct{} { return nil }
func (m *mockConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, ocrtypes.ConfigDigest, error) {
	return m.configCount, m.configDigest, nil
}
func (m *mockConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	return ocrtypes.ContractConfig{ConfigCount: m.configCount, ConfigDigest: m.configDigest}, nil
}
func (m *mockConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	return m.configCount, nil
}

func TestOCRConfigService_MultipleOCRKeys(t *testing.T) {
	lggr := logger.Test(t)
	svc := NewOCRConfigService(lggr, testPeerIDProvider(), 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	// Config with multiple OCR keys (e.g., blue/green deployment).
	blueConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain_blue"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain_blue"),
		ConfigCount:           1,
	}

	greenConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer5"), []byte("signer6"), []byte("signer7"), []byte("signer8")},
		Transmitters:          [][]byte{[]byte("tx5"), []byte("tx6"), []byte("tx7"), []byte("tx8")},
		F:                     1,
		OnchainConfig:         []byte("onchain_green"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain_green"),
		ConfigCount:           2,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			"blue":  blueConfig,
			"green": greenConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	don := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don.Members = []ragetypes.PeerID{testPeerID()}

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	// Verify both configs are stored.
	blueCfg, ok := svc.getConfig("consensus@1.0.0", "blue")
	require.True(t, ok)
	assert.Equal(t, uint64(1), blueCfg.ContractConfig.ConfigCount)

	greenCfg, ok := svc.getConfig("consensus@1.0.0", "green")
	require.True(t, ok)
	assert.Equal(t, uint64(2), greenCfg.ContractConfig.ConfigCount)
}

func TestOCRConfigService_DONMembershipFiltering(t *testing.T) {
	lggr := logger.Test(t)

	// Create a peer ID for this node.
	var myPeerID ragetypes.PeerID
	copy(myPeerID[:], []byte("my-peer-id-1234567890123456"))

	var otherPeerID ragetypes.PeerID
	copy(otherPeerID[:], []byte("other-peer-id-12345678901234"))

	peerIDProvider := func() (ragetypes.PeerID, error) { return myPeerID, nil }
	svc := NewOCRConfigService(lggr, peerIDProvider, 1, "0x1234567890abcdef")

	ctx := t.Context()
	require.NoError(t, svc.Start(ctx))
	defer svc.Close()

	ocrConfig := &capabilitiespb.OCR3Config{
		Signers:               [][]byte{[]byte("signer1"), []byte("signer2"), []byte("signer3"), []byte("signer4")},
		Transmitters:          [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")},
		F:                     1,
		OnchainConfig:         []byte("onchain"),
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte("offchain"),
		ConfigCount:           5,
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		Ocr3Configs: map[string]*capabilitiespb.OCR3Config{
			capabilitiespb.OCR3ConfigDefaultKey: ocrConfig,
		},
	}

	configBytes, err := proto.Marshal(capConfig)
	require.NoError(t, err)

	// Create capabilities.DON with Members field
	don1 := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"consensus@1.0.0": {Config: configBytes},
		},
	}
	don1.Members = []ragetypes.PeerID{myPeerID, otherPeerID}

	don2 := registrysyncer.DON{
		CapabilityConfigurations: map[string]registrysyncer.CapabilityConfiguration{
			"other_cap@1.0.0": {Config: configBytes},
		},
	}
	don2.Members = []ragetypes.PeerID{otherPeerID} // Node not a member

	registry := &registrysyncer.LocalRegistry{
		Logger: lggr,
		IDsToDONs: map[registrysyncer.DonID]registrysyncer.DON{
			1: don1,
			2: don2,
		},
		IDsToNodes:        map[ragetypes.PeerID]registrysyncer.NodeInfo{},
		IDsToCapabilities: map[string]registrysyncer.Capability{},
	}

	err = svc.OnNewRegistry(ctx, registry)
	require.NoError(t, err)

	// Config from DON 1 (node is member) should be stored.
	cfg1, ok := svc.getConfig("consensus@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	assert.True(t, ok)
	assert.Equal(t, uint64(5), cfg1.ContractConfig.ConfigCount)

	// Config from DON 2 (node is NOT member) should NOT be stored.
	_, ok = svc.getConfig("other_cap@1.0.0", capabilitiespb.OCR3ConfigDefaultKey)
	assert.False(t, ok)
}
