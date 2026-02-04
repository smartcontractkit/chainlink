package capregconfig

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"sync"

	"google.golang.org/protobuf/proto"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink/v2/core/services/registrysyncer"
)

var _ OCRConfigService = (*ocrConfigService)(nil)

type ocrConfigService struct {
	services.StateMachine

	lggr            logger.Logger
	chainID         uint64
	registryAddress string
	metrics         *Metrics

	mu                sync.RWMutex
	configs           map[configKey]*cachedConfig
	registryRefreshed bool // true after at least one OnNewRegistry call
}

type configKey struct {
	CapabilityID string
	DonID        uint32
	OCRConfigKey string
}

type cachedConfig struct {
	RawConfig []byte
	// Parsed contract config ready for libocr, with a computed digest.
	ContractConfig ocrtypes.ContractConfig
}

func New(lggr logger.Logger, chainID uint64, registryAddress string) OCRConfigService {
	namedLggr := logger.Named(lggr, "OCRConfigService")

	metrics, err := InitMetrics()
	if err != nil {
		namedLggr.Warnw("failed to initialize metrics, metrics will be disabled", "error", err)
	}

	return &ocrConfigService{
		lggr:            namedLggr,
		chainID:         chainID,
		registryAddress: registryAddress,
		configs:         make(map[configKey]*cachedConfig),
		metrics:         metrics,
	}
}

func (s *ocrConfigService) Start(ctx context.Context) error {
	return s.StartOnce("OCRConfigService", func() error {
		s.lggr.Info("OCRConfigService started")
		return nil
	})
}

func (s *ocrConfigService) Close() error {
	return s.StopOnce("OCRConfigService", func() error {
		s.lggr.Info("OCRConfigService stopped")
		return nil
	})
}

// OnNewRegistry implements registrysyncer.Listener to receive registry updates with capability configurations.
func (s *ocrConfigService) OnNewRegistry(ctx context.Context, registry *registrysyncer.LocalRegistry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.registryRefreshed = true

	for donID, don := range registry.IDsToDONs {
		for capID, capConfig := range don.CapabilityConfigurations {
			if err := s.processCapabilityConfig(ctx, capID, uint32(donID), capConfig.Config); err != nil {
				s.lggr.Warnw("failed to process capability config",
					"capabilityID", capID,
					"donID", donID,
					"error", err,
				)
				s.metrics.IncrementCapabilityConfigErrors(ctx, capID, uint32(donID))
			}
		}
	}

	return nil
}

func (s *ocrConfigService) processCapabilityConfig(ctx context.Context, capabilityID string, donID uint32, configBytes []byte) error {
	if len(configBytes) == 0 {
		return nil
	}

	var capConfig capabilitiespb.CapabilityConfig
	if err := proto.Unmarshal(configBytes, &capConfig); err != nil {
		return fmt.Errorf("failed to unmarshal capability config: %w", err)
	}

	if len(capConfig.Ocr3Configs) == 0 {
		return nil
	}

	for ocrKey, ocrConfig := range capConfig.Ocr3Configs {
		key := configKey{
			CapabilityID: capabilityID,
			DonID:        donID,
			OCRConfigKey: ocrKey,
		}

		ocrConfigBytes, err := proto.Marshal(ocrConfig)
		if err != nil {
			s.metrics.IncrementParseErrors(ctx, capabilityID, donID, ocrKey)
			s.lggr.Errorw("failed to marshal OCR config for comparison",
				"capabilityID", capabilityID,
				"donID", donID,
				"ocrConfigKey", ocrKey,
				"error", err,
			)
			continue
		}

		// Check if config has changed.
		existingConfig, exists := s.configs[key]
		if exists && bytes.Equal(existingConfig.RawConfig, ocrConfigBytes) {
			s.lggr.Debugw("OCR config unchanged",
				"capabilityID", capabilityID,
				"donID", donID,
				"ocrConfigKey", ocrKey,
			)
			continue
		}

		contractConfig, err := s.parseOCR3Config(capabilityID, donID, ocrKey, ocrConfig)
		if err != nil {
			s.metrics.IncrementParseErrors(ctx, capabilityID, donID, ocrKey)
			s.lggr.Errorw("failed to parse OCR3 config",
				"capabilityID", capabilityID,
				"donID", donID,
				"ocrConfigKey", ocrKey,
				"error", err,
			)
			continue
		}

		s.configs[key] = &cachedConfig{
			RawConfig:      ocrConfigBytes,
			ContractConfig: contractConfig,
		}

		s.metrics.IncrementConfigUpdates(ctx, capabilityID, donID, ocrKey)
		s.metrics.SetConfigCount(ctx, capabilityID, donID, ocrKey, int64(ocrConfig.ConfigCount)) //#nosec G115

		s.lggr.Infow("OCR config updated from registry",
			"capabilityID", capabilityID,
			"donID", donID,
			"ocrConfigKey", ocrKey,
			"configCount", ocrConfig.ConfigCount,
			"configDigest", contractConfig.ConfigDigest.Hex(),
		)
	}

	return nil
}

// convert a protobuf OCR3Config to a libocr ContractConfig
func (s *ocrConfigService) parseOCR3Config(
	capabilityID string,
	donID uint32,
	ocrConfigKey string,
	cfg *capabilitiespb.OCR3Config,
) (ocrtypes.ContractConfig, error) {
	signers := make([]ocrtypes.OnchainPublicKey, len(cfg.Signers))
	for i, signer := range cfg.Signers {
		signers[i] = ocrtypes.OnchainPublicKey(signer)
	}

	transmitters := make([]ocrtypes.Account, len(cfg.Transmitters))
	for i, transmitter := range cfg.Transmitters {
		// Accounts are hex strings
		transmitters[i] = ocrtypes.Account(hex.EncodeToString(transmitter))
	}

	digest, err := computeConfigDigest(s.chainID, s.registryAddress, capabilityID, donID, ocrConfigKey, cfg)
	if err != nil {
		return ocrtypes.ContractConfig{}, fmt.Errorf("failed to compute config digest: %w", err)
	}

	if cfg.F > math.MaxUint8 {
		return ocrtypes.ContractConfig{}, fmt.Errorf("f value too large: %d > %d", cfg.F, math.MaxUint8)
	}

	return ocrtypes.ContractConfig{
		ConfigDigest:          digest,
		ConfigCount:           cfg.ConfigCount,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     uint8(cfg.F), //#nosec G115
		OnchainConfig:         cfg.OnchainConfig,
		OffchainConfigVersion: cfg.OffchainConfigVersion,
		OffchainConfig:        cfg.OffchainConfig,
	}, nil
}

func (s *ocrConfigService) GetConfigTracker(
	capabilityID string,
	donID uint32,
	ocrConfigKey string,
	legacyTracker ocrtypes.ContractConfigTracker,
) (ocrtypes.ContractConfigTracker, error) {
	return &dynamicConfigTracker{
		service:       s,
		capabilityID:  capabilityID,
		donID:         donID,
		ocrConfigKey:  ocrConfigKey,
		legacyTracker: legacyTracker,
		lggr:          s.lggr,
	}, nil
}

func (s *ocrConfigService) GetConfigDigester(
	capabilityID string,
	donID uint32,
	ocrConfigKey string,
	legacyDigester ocrtypes.OffchainConfigDigester,
) (ocrtypes.OffchainConfigDigester, error) {
	return &dynamicConfigDigester{
		service:        s,
		capabilityID:   capabilityID,
		donID:          donID,
		ocrConfigKey:   ocrConfigKey,
		legacyDigester: legacyDigester,
		lggr:           s.lggr,
	}, nil
}

func (s *ocrConfigService) getConfig(capabilityID string, donID uint32, ocrConfigKey string) (*cachedConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := configKey{
		CapabilityID: capabilityID,
		DonID:        donID,
		OCRConfigKey: ocrConfigKey,
	}
	cfg, exists := s.configs[key]
	if !exists {
		return nil, false
	}
	return cfg, true
}

func (s *ocrConfigService) hasRefreshedRegistry() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.registryRefreshed
}

// dynamicConfigTracker dynamically switches between legacy config and registry-based config, if available.
type dynamicConfigTracker struct {
	service       *ocrConfigService
	capabilityID  string
	donID         uint32
	ocrConfigKey  string
	legacyTracker ocrtypes.ContractConfigTracker
	lggr          logger.Logger
}

var _ ocrtypes.ContractConfigTracker = (*dynamicConfigTracker)(nil)

func (t *dynamicConfigTracker) Notify() <-chan struct{} {
	// Don't use notifications
	return nil
}

func (t *dynamicConfigTracker) LatestConfigDetails(ctx context.Context) (uint64, ocrtypes.ConfigDigest, error) {
	cfg, ok := t.service.getConfig(t.capabilityID, t.donID, t.ocrConfigKey)
	if ok {
		return cfg.ContractConfig.ConfigCount, cfg.ContractConfig.ConfigDigest, nil
	}

	// Only fall back to legacy if we've received at least one registry update
	// and confirmed there's no config for this capability/DON/key.
	if t.legacyTracker != nil && t.service.hasRefreshedRegistry() {
		return t.legacyTracker.LatestConfigDetails(ctx)
	}

	return 0, ocrtypes.ConfigDigest{}, fmt.Errorf("no config available for %s/%d/%s", t.capabilityID, t.donID, t.ocrConfigKey)
}

func (t *dynamicConfigTracker) LatestConfig(ctx context.Context, changedInBlock uint64) (ocrtypes.ContractConfig, error) {
	cfg, ok := t.service.getConfig(t.capabilityID, t.donID, t.ocrConfigKey)
	if ok {
		return cfg.ContractConfig, nil
	}

	if t.legacyTracker != nil && t.service.hasRefreshedRegistry() {
		return t.legacyTracker.LatestConfig(ctx, changedInBlock)
	}

	return ocrtypes.ContractConfig{}, fmt.Errorf("no config available for %s/%d/%s", t.capabilityID, t.donID, t.ocrConfigKey)
}

func (t *dynamicConfigTracker) LatestBlockHeight(ctx context.Context) (uint64, error) {
	// When using registry-based config, we don't have actual blocks.
	// Return config count as a placeholder. The SkipContractConfigConfirmations
	// should be set to true in LocalConfig when using registry-based config.
	cfg, ok := t.service.getConfig(t.capabilityID, t.donID, t.ocrConfigKey)
	if ok {
		return cfg.ContractConfig.ConfigCount, nil
	}

	if t.legacyTracker != nil && t.service.hasRefreshedRegistry() {
		return t.legacyTracker.LatestBlockHeight(ctx)
	}

	return 0, fmt.Errorf("no config available for %s/%d/%s", t.capabilityID, t.donID, t.ocrConfigKey)
}

// dynamicConfigDigester dynamically switches between legacy config and registry-based config, if available.
type dynamicConfigDigester struct {
	service        *ocrConfigService
	capabilityID   string
	donID          uint32
	ocrConfigKey   string
	legacyDigester ocrtypes.OffchainConfigDigester
	lggr           logger.Logger
}

var _ ocrtypes.OffchainConfigDigester = (*dynamicConfigDigester)(nil)

func (d *dynamicConfigDigester) ConfigDigest(ctx context.Context, cc ocrtypes.ContractConfig) (ocrtypes.ConfigDigest, error) {
	// If we have registry config with matching config count, use our pre-computed digest.
	cfg, ok := d.service.getConfig(d.capabilityID, d.donID, d.ocrConfigKey)
	if ok && cfg.ContractConfig.ConfigCount == cc.ConfigCount {
		return cfg.ContractConfig.ConfigDigest, nil
	}

	if d.legacyDigester != nil && d.service.hasRefreshedRegistry() {
		return d.legacyDigester.ConfigDigest(ctx, cc)
	}

	return ocrtypes.ConfigDigest{}, fmt.Errorf("no digester available for %s/%d/%s", d.capabilityID, d.donID, d.ocrConfigKey)
}

func (d *dynamicConfigDigester) ConfigDigestPrefix(ctx context.Context) (ocrtypes.ConfigDigestPrefix, error) {
	_, ok := d.service.getConfig(d.capabilityID, d.donID, d.ocrConfigKey)
	if ok {
		return ocrtypes.ConfigDigestPrefixKeystoneOCR3Capability, nil
	}

	if d.legacyDigester != nil && d.service.hasRefreshedRegistry() {
		return d.legacyDigester.ConfigDigestPrefix(ctx)
	}

	return 0, fmt.Errorf("no digester available for %s/%d/%s", d.capabilityID, d.donID, d.ocrConfigKey)
}
