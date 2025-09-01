package v2

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	storage_service "github.com/smartcontractkit/chainlink-protos/storage-service/go"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type lastFetchedAtMap struct {
	m map[string]time.Time
	sync.RWMutex
}

func (l *lastFetchedAtMap) Set(url string, at time.Time) {
	l.Lock()
	defer l.Unlock()
	l.m[url] = at
}

func (l *lastFetchedAtMap) Get(url string) (time.Time, bool) {
	l.RLock()
	defer l.RUnlock()
	got, ok := l.m[url]
	return got, ok
}

func newLastFetchedAtMap() *lastFetchedAtMap {
	return &lastFetchedAtMap{
		m: map[string]time.Time{},
	}
}

func safeUint32(n uint64) uint32 {
	if n > math.MaxUint32 {
		return math.MaxUint32
	}
	return uint32(n)
}

type ArtifactConfig struct {
	MaxConfigSize  uint64
	MaxSecretsSize uint64
	MaxBinarySize  uint64
}

// By default, if type is unknown, the largest artifact size is 26.4KB.  Configure the artifact size
// via the ArtifactConfig to override this default.
const defaultMaxArtifactSizeBytes = 26.4 * utils.KB

func (cfg *ArtifactConfig) ApplyDefaults() {
	if cfg.MaxConfigSize == 0 {
		cfg.MaxConfigSize = defaultMaxArtifactSizeBytes
	}
	if cfg.MaxSecretsSize == 0 {
		cfg.MaxSecretsSize = defaultMaxArtifactSizeBytes
	}
	if cfg.MaxBinarySize == 0 {
		cfg.MaxBinarySize = defaultMaxArtifactSizeBytes
	}
}

var defaultSecretsFreshnessDuration = 24 * time.Hour

func WithMaxArtifactSize(cfg ArtifactConfig) func(*Store) {
	return func(a *Store) {
		a.limits = &cfg
	}
}

type StoreConfig struct {
	ArtifactStorageHost string
}

func WithConfig(cfg StoreConfig) func(*Store) {
	return func(a *Store) {
		a.config = &cfg
	}
}

type SerialisedModuleStore interface {
	StoreModule(workflowID string, binaryID string, module []byte) error
	GetModulePath(workflowID string) (string, bool, error)
	GetBinaryID(workflowID string) (string, bool, error)
	DeleteModule(workflowID string) error
}

type Store struct {
	lggr logger.Logger

	// limits sets max artifact sizes to fetch when handling events
	limits *ArtifactConfig
	config *StoreConfig

	orm WorkflowRegistryDS

	// retrieveFunc is a function that retrieves a URL to download an artifact.
	retrieveFunc types.LocationRetrieverFunc
	// fetchFn is a function that fetches the contents of a URL with a limit on the size of the response.
	fetchFn types.FetcherFunc

	lastFetchedAtMap         *lastFetchedAtMap
	clock                    clockwork.Clock
	secretsFreshnessDuration time.Duration

	encryptionKey workflowkey.Key

	emitter custmsg.MessageEmitter
}

func NewStore(lggr logger.Logger, orm WorkflowRegistryDS, fetchFn types.FetcherFunc, retrieveFunc types.LocationRetrieverFunc, clock clockwork.Clock, encryptionKey workflowkey.Key,
	emitter custmsg.MessageEmitter, opts ...func(*Store)) (*Store, error) {
	limits := &ArtifactConfig{}
	limits.ApplyDefaults()

	artifactsStore := &Store{
		lggr:                     lggr,
		orm:                      orm,
		retrieveFunc:             retrieveFunc,
		fetchFn:                  fetchFn,
		lastFetchedAtMap:         newLastFetchedAtMap(),
		clock:                    clock,
		limits:                   limits,
		config:                   &StoreConfig{},
		secretsFreshnessDuration: defaultSecretsFreshnessDuration,
		encryptionKey:            encryptionKey,
		emitter:                  emitter,
	}

	for _, o := range opts {
		o(artifactsStore)
	}

	if retrieveFunc != nil && artifactsStore.config.ArtifactStorageHost == "" {
		return nil, errors.New("storage service URL prefix must be set in the store config")
	}

	return artifactsStore, nil
}

// FetchWorkflowArtifacts fetches the workflow spec and config from a cache or the specified URLs if the artifacts have not
// been cached already.  Before a workflow can be started this method must be called to ensure all artifacts used by the
// workflow are available from the store.
func (h *Store) FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryURL, configURL string) ([]byte, []byte, error) {
	// Check if the workflow spec is already stored in the database
	if spec, err := h.orm.GetWorkflowSpec(ctx, workflowID); err == nil {
		// there is no update in the BinaryURL or ConfigURL, lets decode the stored artifacts
		decodedBinary, err := hex.DecodeString(spec.Workflow)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode stored workflow spec: %w", err)
		}
		return decodedBinary, []byte(spec.Config), nil
	}

	// Determine which URL to retrieve workflow binary artifacts from
	parsedBinaryURL, err := url.Parse(binaryURL)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid binary URL: %w", err)
	}

	// If the binary URL points to the artifact storage host, use the retrieve function to get the signed URL.
	// NOTE: retrieveFunc may be nil if the fetcherFunc was overridden.
	// TODO CRE-632: retrieverFunc should enforced made to always be set, once local CRE can support it.
	if h.retrieveFunc != nil && parsedBinaryURL.Host == h.config.ArtifactStorageHost {
		signedBinaryURL, err2 := h.retrieveFunc(ctx, &storage_service.DownloadArtifactRequest{
			Id:   workflowID,
			Type: storage_service.ArtifactType_ARTIFACT_TYPE_BINARY,
		})
		if err2 != nil {
			return nil, nil, fmt.Errorf("failed to get binary artifact URL: %w", err2)
		}
		binaryURL = signedBinaryURL
	}

	// Fetch the binary files from the specified URLs.
	var (
		binary, decodedBinary, config []byte
	)

	req := ghcapabilities.Request{
		URL:              binaryURL,
		Method:           http.MethodGet,
		MaxResponseBytes: safeUint32(h.limits.MaxBinarySize),
		WorkflowID:       workflowID,
	}
	binary, err = h.fetchFn(ctx, messageID(binaryURL, workflowID), req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch binary from %s : %w", binaryURL, err)
	}

	if decodedBinary, err = base64.StdEncoding.DecodeString(string(binary)); err != nil {
		return nil, nil, fmt.Errorf("failed to decode binary: %w", err)
	}

	if configURL != "" {
		// Determine which URL to retrieve config binary artifacts from
		parsedConfigURL, err2 := url.Parse(configURL)
		if err2 != nil {
			return nil, nil, fmt.Errorf("invalid config URL: %w", err2)
		}

		// If the config URL points to the artifact storage host, use the retrieve function to get the signed URL.
		// NOTE: retrieveFunc may be nil if the fetcherFunc was overridden.
		// TODO CRE-632: retrieverFunc should enforced made to always be set, once local CRE can support it.
		if h.retrieveFunc != nil && parsedConfigURL.Host == h.config.ArtifactStorageHost {
			signedConfigURL, configErr := h.retrieveFunc(ctx, &storage_service.DownloadArtifactRequest{
				Id:   workflowID,
				Type: storage_service.ArtifactType_ARTIFACT_TYPE_CONFIG,
			})
			if configErr != nil {
				return nil, nil, fmt.Errorf("failed to get config artifact URL: %w", configErr)
			}
			configURL = signedConfigURL
		}

		// Fetch the config files from the specified URLs.
		req := ghcapabilities.Request{
			URL:              configURL,
			Method:           http.MethodGet,
			MaxResponseBytes: safeUint32(h.limits.MaxConfigSize),
			WorkflowID:       workflowID,
		}

		config, err2 = h.fetchFn(ctx, messageID(configURL, workflowID), req)
		if err2 != nil {
			return nil, nil, fmt.Errorf("failed to fetch config from %s : %w", configURL, err2)
		}
	}
	return decodedBinary, config, nil
}

func (h *Store) GetWorkflowSpec(ctx context.Context, workflowID string) (*job.WorkflowSpec, error) {
	spec, err := h.orm.GetWorkflowSpec(ctx, workflowID)
	return spec, err
}

func (h *Store) UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	return h.orm.UpsertWorkflowSpec(ctx, spec)
}

// DeleteWorkflowArtifacts removes the workflow spec from the database. If not found, returns nil.
func (h *Store) DeleteWorkflowArtifacts(ctx context.Context, workflowID string) error {
	err := h.orm.DeleteWorkflowSpec(ctx, workflowID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.lggr.Warnw("failed to delete workflow spec: not found", "workflowID", workflowID)
			return nil
		}
		return fmt.Errorf("failed to delete workflow spec: %w", err)
	}

	return nil
}

func (h *Store) GetWasmBinary(ctx context.Context, workflowID string) ([]byte, error) {
	spec, err := h.orm.GetWorkflowSpec(ctx, workflowID)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow spec by workflow ID: %w", err)
	}

	// there is no update in the BinaryURL or ConfigURL, lets decode the stored artifacts
	decodedBinary, err := hex.DecodeString(spec.Workflow)
	if err != nil {
		return nil, fmt.Errorf("failed to decode stored workflow string: %w", err)
	}

	return decodedBinary, nil
}

func messageID(url string, parts ...string) string {
	h := sha256.New()
	h.Write([]byte(url))
	for _, p := range parts {
		h.Write([]byte(p))
	}
	hash := hex.EncodeToString(h.Sum(nil))
	p := []string{ghcapabilities.MethodWorkflowSyncer, hash}
	return strings.Join(p, "/")
}
