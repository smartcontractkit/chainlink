package artifacts

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ghcapabilities "github.com/smartcontractkit/chainlink/v2/core/services/gateway/handlers/capabilities"

	"github.com/smartcontractkit/chainlink-common/pkg/workflows/secrets"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/workflowkey"
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

// FetcherFunc is an abstraction for fetching the contents stored at a URL.
type FetcherFunc func(ctx context.Context, messageID string, req ghcapabilities.Request) ([]byte, error)

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

// logCustMsg emits a custom message to the external sink and logs an error if that fails.
func logCustMsg(ctx context.Context, cma custmsg.MessageEmitter, msg string, log logger.Logger) {
	err := cma.Emit(ctx, msg)
	if err != nil {
		log.Helper(1).Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
	}
}

var defaultSecretsFreshnessDuration = 24 * time.Hour

func WithMaxArtifactSize(cfg ArtifactConfig) func(*Store) {
	return func(a *Store) {
		a.limits = &cfg
	}
}

type SerialisedModuleStore interface {
	StoreModule(workflowID string, binaryID string, module []byte) error
	GetModulePath(workflowID string) (string, bool, error)
	GetBinaryID(workflowID string) (string, bool, error)
	DeleteModule(workflowID string) error
}

type decryptSecretsFn func(data []byte, owner string) (map[string]string, error)

type Store struct {
	lggr logger.Logger

	// limits sets max artifact sizes to fetch when handling events
	limits *ArtifactConfig

	orm WorkflowRegistryDS

	// fetchFn is a function that fetches the contents of a URL with a limit on the size of the response.
	fetchFn FetcherFunc

	lastFetchedAtMap         *lastFetchedAtMap
	clock                    clockwork.Clock
	secretsFreshnessDuration time.Duration

	encryptionKey workflowkey.Key

	decryptSecrets decryptSecretsFn

	emitter custmsg.MessageEmitter
}

func NewStore(lggr logger.Logger, orm WorkflowRegistryDS, fetchFn FetcherFunc, clock clockwork.Clock, encryptionKey workflowkey.Key,
	emitter custmsg.MessageEmitter, opts ...func(*Store)) *Store {
	return NewStoreWithDecryptSecretsFn(lggr, orm, fetchFn, clock, encryptionKey, emitter,
		func(data []byte, owner string) (map[string]string, error) {
			secretsPayload := secrets.EncryptedSecretsResult{}
			err := json.Unmarshal(data, &secretsPayload)
			if err != nil {
				return nil, fmt.Errorf("could not unmarshal secrets: %w", err)
			}

			return secrets.DecryptSecretsForNode(secretsPayload, encryptionKey, owner)
		},
		opts...)
}

func NewStoreWithDecryptSecretsFn(lggr logger.Logger, orm WorkflowRegistryDS, fetchFn FetcherFunc, clock clockwork.Clock, encryptionKey workflowkey.Key,
	emitter custmsg.MessageEmitter, decryptSecrets decryptSecretsFn, opts ...func(*Store)) *Store {
	limits := &ArtifactConfig{}
	limits.ApplyDefaults()

	artifactsStore := &Store{
		lggr:                     lggr,
		orm:                      orm,
		fetchFn:                  fetchFn,
		lastFetchedAtMap:         newLastFetchedAtMap(),
		clock:                    clock,
		limits:                   limits,
		secretsFreshnessDuration: defaultSecretsFreshnessDuration,
		encryptionKey:            encryptionKey,
		emitter:                  emitter,
		decryptSecrets:           decryptSecrets,
	}

	for _, o := range opts {
		o(artifactsStore)
	}

	return artifactsStore
}

// FetchWorkflowArtifacts fetches the workflow spec and config from a cache or the specified URLs if the artifacts have not
// been cached already.  Before a workflow can be started this method must be called to ensure all artifacts used by the
// workflow are available from the store.
func (h *Store) FetchWorkflowArtifacts(ctx context.Context, workflowID, binaryURL, configURL string) ([]byte, []byte, error) {
	// Check if the workflow spec is already stored in the database
	if spec, err := h.orm.GetWorkflowSpecByID(ctx, workflowID); err == nil {
		// there is no update in the BinaryURL or ConfigURL, lets decode the stored artifacts
		decodedBinary, err := hex.DecodeString(spec.Workflow)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to decode stored workflow spec: %w", err)
		}
		return decodedBinary, []byte(spec.Config), nil
	}

	// Fetch the binary and config files from the specified URLs.
	var (
		binary, decodedBinary, config []byte
		err                           error
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
		req := ghcapabilities.Request{
			URL:              configURL,
			Method:           http.MethodGet,
			MaxResponseBytes: safeUint32(h.limits.MaxConfigSize),
			WorkflowID:       workflowID,
		}
		config, err = h.fetchFn(ctx, messageID(configURL, workflowID), req)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch config from %s : %w", configURL, err)
		}
	}
	return decodedBinary, config, nil
}

func (h *Store) GetSecrets(ctx context.Context, secretsURL string, workflowID [32]byte, workflowOwner []byte) ([]byte, error) {
	wid := hex.EncodeToString(workflowID[:])
	req := ghcapabilities.Request{
		URL:              secretsURL,
		Method:           http.MethodGet,
		MaxResponseBytes: safeUint32(h.limits.MaxSecretsSize),
		WorkflowID:       wid,
	}
	fetchedSecrets, fetchErr := h.fetchFn(ctx, messageID(secretsURL, wid), req)
	if fetchErr != nil {
		return nil, fmt.Errorf("failed to fetch secrets from %s : %w", secretsURL, fetchErr)
	}

	// sanity check by decoding the secrets
	_, decryptErr := h.decryptSecrets(fetchedSecrets, hex.EncodeToString(workflowOwner))
	if decryptErr != nil {
		return nil, fmt.Errorf("failed to decrypt secrets %s: %w", secretsURL, decryptErr)
	}

	return fetchedSecrets, nil
}

func (h *Store) ForceUpdateSecrets(
	ctx context.Context,
	secretsURLHash []byte,
	owner []byte,
) (string, error) {
	// Get the URL of the secrets file from the event data
	hash := hex.EncodeToString(secretsURLHash)

	url, err := h.orm.GetSecretsURLByHash(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("failed to get URL by hash %s : %w", hash, err)
	}

	ownerHex := hex.EncodeToString(owner)
	req := ghcapabilities.Request{
		URL:              url,
		Method:           http.MethodGet,
		MaxResponseBytes: safeUint32(h.limits.MaxSecretsSize),
		// TODO -- fix, but this is used for rate limiting purposes
		WorkflowID: hex.EncodeToString(owner),
	}
	// Fetch the contents of the secrets file from the url via the fetcher
	secrets, err := h.fetchFn(ctx, messageID(url, ownerHex), req)
	if err != nil {
		return "", err
	}

	// Sanity check the payload and ensure we can decrypt it.
	// If we can't, let's return an error and we won't store the result in the DB.
	_, err = h.decryptSecrets(secrets, hex.EncodeToString(owner))
	if err != nil {
		return "", fmt.Errorf("failed to validate secrets: could not decrypt: %w", err)
	}

	h.lastFetchedAtMap.Set(hash, h.clock.Now())

	// Update the secrets in the ORM
	if _, err := h.orm.Update(ctx, hash, string(secrets)); err != nil {
		return "", fmt.Errorf("failed to update secrets: %w", err)
	}

	return string(secrets), nil
}

func (h *Store) GetSecretsURLByID(ctx context.Context, id int64) (string, error) {
	secretsURL, err := h.orm.GetSecretsURLByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to get secrets URL by ID: %w", err)
	}
	return secretsURL, nil
}

func (h *Store) GetWorkflowSpec(ctx context.Context, workflowOwner string, workflowName string) (*job.WorkflowSpec, error) {
	spec, err := h.orm.GetWorkflowSpec(ctx, workflowOwner, workflowName)
	return spec, err
}

func (h *Store) SecretsFor(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error) {
	secretsURLHash, secretsPayload, err := h.orm.GetContentsByWorkflowID(ctx, workflowID)
	if err != nil {
		// The workflow record was found, but secrets_id was empty.
		// Let's just stub out the response.
		if errors.Is(err, ErrEmptySecrets) {
			return map[string]string{}, nil
		}

		return nil, fmt.Errorf("failed to fetch secrets by workflow ID: %w", err)
	}

	lastFetchedAt, ok := h.lastFetchedAtMap.Get(secretsURLHash)
	if !ok || h.clock.Now().Sub(lastFetchedAt) > h.secretsFreshnessDuration {
		updatedSecrets, innerErr := h.refreshSecrets(ctx, workflowOwner, hexWorkflowName, workflowID, secretsURLHash)
		if innerErr != nil {
			msg := fmt.Sprintf("could not refresh secrets: proceeding with stale secrets for workflowID %s: %s", workflowID, innerErr)
			h.lggr.Error(msg)

			logCustMsg(
				ctx,
				h.emitter.With(
					platform.KeyWorkflowID, workflowID,
					platform.KeyWorkflowName, decodedWorkflowName,
					platform.KeyWorkflowOwner, workflowOwner,
				),
				msg,
				h.lggr,
			)
		} else {
			secretsPayload = updatedSecrets
		}
	}

	return h.decryptSecrets([]byte(secretsPayload), workflowOwner)
}

func (h *Store) UpsertWorkflowSpec(ctx context.Context, spec *job.WorkflowSpec) (int64, error) {
	return h.orm.UpsertWorkflowSpec(ctx, spec)
}

func (h *Store) UpsertWorkflowSpecWithSecrets(ctx context.Context, entry *job.WorkflowSpec, secretsURL, urlHash, secrets string) (int64, error) {
	return h.orm.UpsertWorkflowSpecWithSecrets(ctx, entry, secretsURL, urlHash, secrets)
}

func (h *Store) GetSecretsURLHash(workflowOwner []byte, secretsURL []byte) ([]byte, error) {
	urlHash, err := h.orm.GetSecretsURLHash(workflowOwner, secretsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets URL hash: %w", err)
	}
	return urlHash, nil
}

// DeleteWorkflowArtifacts removes the workflow spec from the database. If not found, returns nil.
func (h *Store) DeleteWorkflowArtifacts(ctx context.Context, workflowOwner string, workflowName string, workflowID string) error {
	err := h.orm.DeleteWorkflowSpec(ctx, workflowOwner, workflowName)
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
	spec, err := h.orm.GetWorkflowSpecByID(ctx, workflowID)
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

func (h *Store) refreshSecrets(ctx context.Context, workflowOwner, workflowName, workflowID, secretsURLHash string) (string, error) {
	owner, err := hex.DecodeString(workflowOwner)
	if err != nil {
		return "", err
	}

	decodedHash, err := hex.DecodeString(secretsURLHash)
	if err != nil {
		return "", err
	}

	updatedSecrets, err := h.ForceUpdateSecrets(
		ctx, decodedHash, owner)
	if err != nil {
		return "", err
	}

	return updatedSecrets, nil
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
