package agentcommerce

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func storeRecordFixture(t *testing.T) CommerceRecord {
	t.Helper()
	intent := sha256.Sum256([]byte("store-intent"))
	intentHash := hex.EncodeToString(intent[:])
	executionID, _ := ExecutionIDForIntent(intentHash)
	output := sha256.Sum256([]byte("output"))
	report := sha256.Sum256([]byte("report"))
	return CommerceRecord{Version: "1", Scope: SettlementScope{TenantID: "tenant-a", UserID: "user-a", ScopeID: "scope-a"}, IntentHash: intentHash, ExecutionID: executionID, ExecutionPhase: ExecutionPrepared, OutputHash: hex.EncodeToString(output[:]), ReportHash: hex.EncodeToString(report[:]), UpdatedAt: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
}

func TestFileCommerceStoreScopeIdempotencyAndTerminalConflict(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "store.json")
	store, _ := NewFileCommerceStore(path)
	record := storeRecordFixture(t)
	if _, created, err := store.Put(ctx, record); err != nil || !created {
		t.Fatalf("create=%v err=%v", created, err)
	}
	retry := record
	retry.UpdatedAt = retry.UpdatedAt.Add(time.Hour)
	if _, created, err := store.Put(ctx, retry); err != nil || created {
		t.Fatalf("retry create=%v err=%v", created, err)
	}
	wrongScope := record.Scope
	wrongScope.TenantID = "tenant-b"
	if _, err := store.Get(ctx, wrongScope, record.IntentHash, record.ExecutionID); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("cross scope read=%v", err)
	}
	for _, phase := range []ExecutionPhase{ExecutionSubmitted, ExecutionExecuting, ExecutionExecuted, ExecutionVerified} {
		if phase == ExecutionVerified {
			evidence := sha256.Sum256([]byte("verified-evidence"))
			record.EvidenceHash = hex.EncodeToString(evidence[:])
		}
		record.ExecutionPhase = phase
		record.UpdatedAt = record.UpdatedAt.Add(time.Second)
		if _, _, err := store.Put(ctx, record); err != nil {
			t.Fatal(err)
		}
	}
	record.ExecutionPhase = ExecutionFailed
	record.UpdatedAt = record.UpdatedAt.Add(time.Second)
	if _, _, err := store.Put(ctx, record); !errors.Is(err, ErrCommerceConflict) {
		t.Fatalf("terminal conflict=%v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("store permissions=%o", info.Mode().Perm())
	}
}

func TestFileCommerceStoreRejectsImpossibleInitialStates(t *testing.T) {
	for _, phase := range []ExecutionPhase{ExecutionSubmitted, ExecutionExecuting, ExecutionExecuted, ExecutionVerified, ExecutionFailed, ExecutionUnknown} {
		t.Run(string(phase), func(t *testing.T) {
			store, _ := NewFileCommerceStore(filepath.Join(t.TempDir(), "store.json"))
			record := storeRecordFixture(t)
			record.ExecutionPhase = phase
			if _, _, err := store.Put(context.Background(), record); err == nil {
				t.Fatal("impossible initial phase accepted")
			}
		})
	}
	store, _ := NewFileCommerceStore(filepath.Join(t.TempDir(), "store.json"))
	record := storeRecordFixture(t)
	record.SettlementPhase = SettlementPrepared
	if _, _, err := store.Put(context.Background(), record); err == nil {
		t.Fatal("initial settlement accepted before execution record")
	}
}

func TestFileCommerceStoreConcurrentReplayClaims(t *testing.T) {
	store, _ := NewFileCommerceStore(filepath.Join(t.TempDir(), "store.json"))
	scope := SettlementScope{TenantID: "tenant", UserID: "user", ScopeID: "scope"}
	a := sha256.Sum256([]byte("a"))
	b := sha256.Sum256([]byte("b"))
	hashes := []string{hex.EncodeToString(a[:]), hex.EncodeToString(b[:])}
	var wg sync.WaitGroup
	var mu sync.Mutex
	success := map[string]int{}
	conflicts := 0
	for i := 0; i < 40; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			hash := hashes[i%2]
			err := store.ClaimReplay(context.Background(), scope, "shared", hash)
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				success[hash]++
			} else if errors.Is(err, ErrReplayClaimed) {
				conflicts++
			} else {
				t.Errorf("claim=%v", err)
			}
		}(i)
	}
	wg.Wait()
	if len(success) != 1 || conflicts == 0 {
		t.Fatalf("success=%v conflicts=%d", success, conflicts)
	}
}

func TestFileCommerceStoreLoadCorruptionMatrix(t *testing.T) {
	mutations := map[string]func(map[string]any){
		"schema":     func(doc map[string]any) { doc["schema_version"] = float64(99) },
		"audit head": func(doc map[string]any) { doc["audit_head"] = "bad" },
		"record": func(doc map[string]any) {
			for _, raw := range doc["records"].(map[string]any) {
				raw.(map[string]any)["RecordHash"] = "bad"
				break
			}
		},
		"replay": func(doc map[string]any) {
			for key := range doc["replay"].(map[string]any) {
				delete(doc["replay"].(map[string]any), key)
				doc["replay"].(map[string]any)["corrupt"] = "bad"
				break
			}
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			path := filepath.Join(t.TempDir(), "store.json")
			store, _ := NewFileCommerceStore(path)
			record := storeRecordFixture(t)
			if _, _, err := store.Put(ctx, record); err != nil {
				t.Fatal(err)
			}
			evidence := sha256.Sum256([]byte("evidence"))
			if err := store.ClaimReplay(ctx, record.Scope, "replay", hex.EncodeToString(evidence[:])); err != nil {
				t.Fatal(err)
			}
			payload, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var doc map[string]any
			if err := json.Unmarshal(payload, &doc); err != nil {
				t.Fatal(err)
			}
			mutate(doc)
			payload, _ = json.Marshal(doc)
			if err := os.WriteFile(path, payload, 0o600); err != nil {
				t.Fatal(err)
			}
			reopened, _ := NewFileCommerceStore(path)
			if err := reopened.VerifyIntegrity(ctx); err == nil {
				t.Fatal("corruption accepted")
			}
		})
	}
}
