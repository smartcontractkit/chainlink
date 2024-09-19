package llo

import (
	"context"
	"fmt"
	sync "sync"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
)

// RetirementReportCacheReader is used by the plugin-scoped
// RetirementReportCache
type RetirementReportCacheReader interface {
	AttestedRetirementReport(cd ocr2types.ConfigDigest) ([]byte, bool)
	Config(cd ocr2types.ConfigDigest) (Config, bool)
}

// RetirementReportCache is intended to be a global singleton that is wrapped
// by a PluginScopedRetirementReportCache for a given plugin
type RetirementReportCache interface {
	services.Service
	StoreAttestedRetirementReport(ctx context.Context, cd ocrtypes.ConfigDigest, seqNr uint64, retirementReport []byte, sigs []types.AttributedOnchainSignature) error
	StoreConfig(ctx context.Context, cd ocr2types.ConfigDigest, signers [][]byte, f uint8) error
	RetirementReportCacheReader
}

type retirementReportCache struct {
	services.Service
	eng *services.Engine

	mu      sync.RWMutex
	arrs    map[ocr2types.ConfigDigest][]byte
	configs map[ocr2types.ConfigDigest]Config

	orm RetirementReportCacheORM
}

func NewRetirementReportCache(lggr logger.Logger, ds sqlutil.DataSource) RetirementReportCache {
	orm := &retirementReportCacheORM{ds: ds}
	return newRetirementReportCache(lggr, orm)
}

func newRetirementReportCache(lggr logger.Logger, orm RetirementReportCacheORM) *retirementReportCache {
	r := &retirementReportCache{
		arrs:    make(map[ocr2types.ConfigDigest][]byte),
		configs: make(map[ocr2types.ConfigDigest]Config),
		orm:     orm,
	}
	r.Service, r.eng = services.Config{
		Name:  "RetirementReportCache",
		Start: r.start,
	}.NewServiceEngine(lggr)
	return r
}

// NOTE: Could do this lazily instead if we wanted to avoid a performance hit
// or potential tables missing etc on application startup (since
// RetirementReportCache is global)
func (r *retirementReportCache) start(ctx context.Context) (err error) {
	// Load all attested retirement reports from the ORM
	// and store them in the cache
	r.arrs, err = r.orm.LoadAttestedRetirementReports(ctx)
	if err != nil {
		return fmt.Errorf("failed to load attested retirement reports: %w", err)
	}
	configs, err := r.orm.LoadConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configs: %w", err)
	}
	for _, c := range configs {
		r.configs[c.Digest] = c
	}
	return nil
}

func (r *retirementReportCache) StoreAttestedRetirementReport(ctx context.Context, cd ocr2types.ConfigDigest, seqNr uint64, retirementReport []byte, sigs []types.AttributedOnchainSignature) error {
	r.mu.RLock()
	if _, ok := r.arrs[cd]; ok {
		r.mu.RUnlock()
		return nil
	}
	r.mu.RUnlock()

	var pbSigs []*AttributedOnchainSignature
	for _, s := range sigs {
		pbSigs = append(pbSigs, &AttributedOnchainSignature{
			Signer:    uint32(s.Signer),
			Signature: s.Signature,
		})
	}
	attestedRetirementReport := AttestedRetirementReport{
		RetirementReport: retirementReport,
		SeqNr:            seqNr,
		Sigs:             pbSigs,
	}

	serialized, err := proto.Marshal(&attestedRetirementReport)
	if err != nil {
		return fmt.Errorf("StoreAttestedRetirementReport failed; failed to marshal protobuf: %w", err)
	}

	if err := r.orm.StoreAttestedRetirementReport(ctx, cd, serialized); err != nil {
		return fmt.Errorf("StoreAttestedRetirementReport failed; failed to persist to ORM: %w", err)
	}

	r.mu.Lock()
	r.arrs[cd] = serialized
	r.mu.Unlock()

	return nil
}

func (r *retirementReportCache) StoreConfig(ctx context.Context, cd ocr2types.ConfigDigest, signers [][]byte, f uint8) error {
	r.mu.RLock()
	if _, ok := r.configs[cd]; ok {
		r.mu.RUnlock()
		return nil
	}
	r.mu.RUnlock()

	r.mu.Lock()
	r.configs[cd] = Config{
		Digest:  cd,
		Signers: signers,
		F:       f,
	}
	r.mu.Unlock()

	return r.orm.StoreConfig(ctx, cd, signers, f)
}

func (r *retirementReportCache) AttestedRetirementReport(predecessorConfigDigest ocr2types.ConfigDigest) ([]byte, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	arr, exists := r.arrs[predecessorConfigDigest]
	return arr, exists
}

func (r *retirementReportCache) Config(cd ocr2types.ConfigDigest) (Config, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.configs[cd]
	return c, exists
}
