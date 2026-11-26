package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrCommerceConflict = errors.New("agent commerce record conflict")
	ErrReplayClaimed    = errors.New("agent commerce replay key already claimed")
)

type CommerceRecord struct {
	Version         string
	Scope           SettlementScope
	IntentHash      string
	ExecutionID     string
	ExecutionPhase  ExecutionPhase
	OutputHash      string
	ReportHash      string
	EvidenceHash    string
	SettlementPhase SettlementPhase
	SettlementHash  string
	UpdatedAt       time.Time
	RecordHash      string
}

type CommerceAuditEntry struct {
	Sequence     uint64
	PreviousHash string
	RecordHash   string
	EntryHash    string
	RecordedAt   time.Time
}

type commerceState struct {
	SchemaVersion int                       `json:"schema_version"`
	Records       map[string]CommerceRecord `json:"records"`
	Replay        map[string]string         `json:"replay"`
	Audit         []CommerceAuditEntry      `json:"audit"`
	AuditHead     string                    `json:"audit_head"`
}

const commerceStoreSchemaVersion = 1

// FileCommerceStore is deterministic durable storage for local development and
// tests. It is not a production distributed or multi-process database.
type FileCommerceStore struct {
	mu   sync.Mutex
	path string
}

func NewFileCommerceStore(path string) (*FileCommerceStore, error) {
	if path == "" || !filepath.IsAbs(path) {
		return nil, errors.New("an absolute commerce store path is required")
	}
	return &FileCommerceStore{path: path}, nil
}

func (s *FileCommerceStore) Put(ctx context.Context, record CommerceRecord) (CommerceRecord, bool, error) {
	if ctx == nil {
		return CommerceRecord{}, false, errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return CommerceRecord{}, false, err
	}
	sealed, err := sealCommerceRecord(record)
	if err != nil {
		return CommerceRecord{}, false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return CommerceRecord{}, false, err
	}
	scopeKey, _ := settlementScopeKey(sealed.Scope)
	key := scopeKey + ":" + sealed.IntentHash + ":" + sealed.ExecutionID
	if existing, ok := state.Records[key]; ok {
		if stableCommerceIdentity(existing) == stableCommerceIdentity(sealed) {
			return existing, false, nil
		}
		if err := validateRecordUpdate(existing, sealed); err != nil {
			return CommerceRecord{}, false, fmt.Errorf("%w: %v", ErrCommerceConflict, err)
		}
	} else if sealed.ExecutionPhase != ExecutionPrepared || sealed.SettlementPhase != "" {
		return CommerceRecord{}, false, errors.New("new commerce record must begin at execution PREPARED with no settlement")
	}
	state.Records[key] = sealed
	entry := CommerceAuditEntry{Sequence: uint64(len(state.Audit) + 1), RecordHash: sealed.RecordHash, RecordedAt: sealed.UpdatedAt}
	if len(state.Audit) > 0 {
		entry.PreviousHash = state.Audit[len(state.Audit)-1].EntryHash
	}
	entry.EntryHash, err = auditEntryHash(entry)
	if err != nil {
		return CommerceRecord{}, false, err
	}
	state.Audit = append(state.Audit, entry)
	state.AuditHead = entry.EntryHash
	if err := s.save(state); err != nil {
		return CommerceRecord{}, false, err
	}
	return sealed, true, nil
}

func (s *FileCommerceStore) ClaimReplay(ctx context.Context, scope SettlementScope, key, evidenceHash string) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if key == "" {
		return errors.New("replay key is required")
	}
	scopeKey, err := settlementScopeKey(scope)
	if err != nil {
		return err
	}
	if err := ValidateSHA256Hex("evidence hash", evidenceHash); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return err
	}
	key = scopeKey + ":" + key
	if existing, ok := state.Replay[key]; ok {
		if existing == evidenceHash {
			return nil
		}
		return ErrReplayClaimed
	}
	state.Replay[key] = evidenceHash
	return s.save(state)
}

func (s *FileCommerceStore) Get(ctx context.Context, scope SettlementScope, intentHash, executionID string) (CommerceRecord, error) {
	if ctx == nil {
		return CommerceRecord{}, errors.New("context is required")
	}
	if err := ctx.Err(); err != nil {
		return CommerceRecord{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return CommerceRecord{}, err
	}
	scopeKey, err := settlementScopeKey(scope)
	if err != nil {
		return CommerceRecord{}, err
	}
	record, ok := state.Records[scopeKey+":"+intentHash+":"+executionID]
	if !ok {
		return CommerceRecord{}, os.ErrNotExist
	}
	sealed, err := sealCommerceRecord(record)
	if err != nil || sealed.RecordHash != record.RecordHash {
		return CommerceRecord{}, errors.New("commerce record integrity check failed")
	}
	return record, nil
}

func (s *FileCommerceStore) VerifyIntegrity(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.load()
	if err != nil {
		return err
	}
	previous := ""
	for i, entry := range state.Audit {
		if entry.Sequence != uint64(i+1) || entry.PreviousHash != previous {
			return errors.New("commerce audit sequence is invalid")
		}
		want, err := auditEntryHash(entry)
		if err != nil || want != entry.EntryHash {
			return errors.New("commerce audit hash is invalid")
		}
		previous = entry.EntryHash
	}
	for _, record := range state.Records {
		sealed, err := sealCommerceRecord(record)
		if err != nil || sealed.RecordHash != record.RecordHash {
			return errors.New("commerce record hash is invalid")
		}
	}
	return nil
}

func (s *FileCommerceStore) load() (commerceState, error) {
	state := commerceState{SchemaVersion: commerceStoreSchemaVersion, Records: map[string]CommerceRecord{}, Replay: map[string]string{}}
	payload, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return commerceState{}, err
	}
	if err := json.Unmarshal(payload, &state); err != nil {
		return commerceState{}, fmt.Errorf("decode commerce store: %w", err)
	}
	if state.Records == nil || state.Replay == nil {
		return commerceState{}, errors.New("invalid commerce store")
	}
	if state.SchemaVersion != commerceStoreSchemaVersion {
		return commerceState{}, fmt.Errorf("unsupported commerce store schema version %d", state.SchemaVersion)
	}
	if err := validateCommerceState(state); err != nil {
		return commerceState{}, err
	}
	return state, nil
}

func (s *FileCommerceStore) save(state commerceState) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	payload, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".agentcommerce-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(payload); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, s.path); err != nil {
		return err
	}
	if err := os.Chmod(s.path, 0o600); err != nil {
		return err
	}
	directory, err := os.Open(filepath.Dir(s.path))
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}

func sealCommerceRecord(record CommerceRecord) (CommerceRecord, error) {
	if record.Version != "1" {
		return CommerceRecord{}, errors.New("unsupported commerce record version")
	}
	if record.IntentHash == "" || record.ExecutionID == "" || record.ExecutionPhase == "" {
		return CommerceRecord{}, errors.New("incomplete commerce record")
	}
	if err := validateSettlementScope(record.Scope); err != nil {
		return CommerceRecord{}, err
	}
	if err := ValidateSHA256Hex("intent hash", record.IntentHash); err != nil {
		return CommerceRecord{}, err
	}
	wantExecutionID, err := ExecutionIDForIntent(record.IntentHash)
	if err != nil || record.ExecutionID != wantExecutionID {
		return CommerceRecord{}, errors.New("commerce record execution id mismatch")
	}
	if record.UpdatedAt.IsZero() {
		return CommerceRecord{}, errors.New("commerce record update time is required")
	}
	switch record.ExecutionPhase {
	case ExecutionPrepared, ExecutionSubmitted, ExecutionExecuting, ExecutionExecuted, ExecutionVerified, ExecutionFailed, ExecutionUnknown:
	default:
		return CommerceRecord{}, errors.New("invalid execution phase")
	}
	if record.OutputHash != "" {
		if err := ValidateSHA256Hex("output hash", record.OutputHash); err != nil {
			return CommerceRecord{}, err
		}
	}
	if record.ReportHash != "" {
		if err := ValidateSHA256Hex("report hash", record.ReportHash); err != nil {
			return CommerceRecord{}, err
		}
	}
	if record.EvidenceHash != "" {
		if err := ValidateSHA256Hex("evidence hash", record.EvidenceHash); err != nil {
			return CommerceRecord{}, err
		}
	}
	if record.SettlementHash != "" {
		if err := ValidateSHA256Hex("settlement hash", record.SettlementHash); err != nil {
			return CommerceRecord{}, err
		}
	}
	if record.SettlementPhase != "" {
		switch record.SettlementPhase {
		case SettlementPrepared, SettlementSubmitted, SettlementBroadcast, SettlementConfirmed, SettlementSettled, SettlementRefunded, SettlementFailed, SettlementUnknown:
		default:
			return CommerceRecord{}, errors.New("invalid settlement phase")
		}
	}
	if record.ExecutionPhase == ExecutionExecuted || record.ExecutionPhase == ExecutionVerified {
		if record.OutputHash == "" || record.ReportHash == "" {
			return CommerceRecord{}, errors.New("executed commerce record requires output and report hashes")
		}
	}
	if record.ExecutionPhase == ExecutionVerified && record.EvidenceHash == "" {
		return CommerceRecord{}, errors.New("verified commerce record requires evidence hash")
	}
	if record.SettlementPhase == SettlementSettled && record.SettlementHash == "" {
		return CommerceRecord{}, errors.New("settled commerce record requires settlement hash")
	}
	copy := record
	copy.RecordHash = ""
	payload, err := json.Marshal(copy)
	if err != nil {
		return CommerceRecord{}, err
	}
	sum := sha256.Sum256(payload)
	record.RecordHash = hex.EncodeToString(sum[:])
	return record, nil
}

// stableCommerceIdentity intentionally excludes UpdatedAt and RecordHash.
func stableCommerceIdentity(record CommerceRecord) string {
	record.UpdatedAt = time.Time{}
	record.RecordHash = ""
	payload, _ := json.Marshal(record)
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func validateRecordUpdate(existing, next CommerceRecord) error {
	if existing.Version != next.Version || existing.Scope != next.Scope || existing.IntentHash != next.IntentHash || existing.ExecutionID != next.ExecutionID ||
		existing.OutputHash != next.OutputHash || existing.ReportHash != next.ReportHash {
		return errors.New("immutable execution identity changed")
	}
	if existing.EvidenceHash != "" && existing.EvidenceHash != next.EvidenceHash {
		return errors.New("evidence hash changed")
	}
	if existing.SettlementHash != "" && existing.SettlementHash != next.SettlementHash {
		return errors.New("settlement hash changed")
	}
	if err := ValidateExecutionTransition(existing.ExecutionPhase, next.ExecutionPhase); err != nil {
		return err
	}
	if existing.SettlementPhase != "" || next.SettlementPhase != "" {
		if existing.SettlementPhase == "" {
			if next.SettlementPhase != SettlementPrepared {
				return errors.New("settlement must begin at PREPARED")
			}
		} else if err := ValidateSettlementTransition(existing.SettlementPhase, next.SettlementPhase); err != nil {
			return err
		}
	}
	return nil
}

func validateCommerceState(state commerceState) error {
	previous := ""
	for i, entry := range state.Audit {
		if entry.Sequence != uint64(i+1) || entry.PreviousHash != previous {
			return errors.New("commerce audit sequence is invalid")
		}
		want, err := auditEntryHash(entry)
		if err != nil || want != entry.EntryHash {
			return errors.New("commerce audit hash is invalid")
		}
		previous = entry.EntryHash
	}
	if state.AuditHead != previous {
		return errors.New("commerce audit head is invalid")
	}
	for key, record := range state.Records {
		sealed, err := sealCommerceRecord(record)
		if err != nil || sealed.RecordHash != record.RecordHash {
			return errors.New("commerce record hash is invalid")
		}
		scopeKey, _ := settlementScopeKey(record.Scope)
		if key != scopeKey+":"+record.IntentHash+":"+record.ExecutionID {
			return errors.New("commerce record storage key is invalid")
		}
	}
	for key, evidenceHash := range state.Replay {
		separator := strings.IndexByte(key, ':')
		if separator != sha256.Size*2 || separator == len(key)-1 {
			return errors.New("commerce replay key is invalid")
		}
		if err := ValidateSHA256Hex("replay scope hash", key[:separator]); err != nil {
			return err
		}
		if err := ValidateSHA256Hex("replay evidence hash", evidenceHash); err != nil {
			return err
		}
	}
	return nil
}

func auditEntryHash(entry CommerceAuditEntry) (string, error) {
	copy := entry
	copy.EntryHash = ""
	payload, err := json.Marshal(copy)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
