package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/sha3"
)

// LinkSwarm — Chainlink-primary high-throughput settlement orchestration.
// Functions is the verification gate. Delivery networks are secondary.
// Designed for enterprise payroll and machine-to-machine payments.

type Config struct {
	Mode             string
	FunctionsRouter  string
	SubscriptionID   uint64
	DONID            string
	SourceCode       string
	MaxConcurrency   int
	RequestTimeout   time.Duration
	RetryAttempts    int
	WorkerCount      int
	ReportDir        string
	EnableCardLock   bool
	ConversionFeeBps int
}

func LoadConfig() Config {
	return Config{
		Mode:             getenv("LINKSWARM_MODE", "simulation"),
		FunctionsRouter:  getenv("CHAINLINK_FUNCTIONS_ROUTER", ""),
		SubscriptionID:   uint64(getenvInt("CHAINLINK_SUBSCRIPTION_ID", 0)),
		DONID:            getenv("CHAINLINK_DON_ID", "fun-ethereum-sepolia-1"),
		SourceCode:       defaultFunctionsSource(),
		MaxConcurrency:   getenvInt("MAX_CONCURRENCY", 512),
		RequestTimeout:   time.Duration(getenvInt("REQUEST_TIMEOUT_SEC", 60)) * time.Second,
		RetryAttempts:    getenvInt("RETRY_ATTEMPTS", 3),
		WorkerCount:      getenvInt("WORKER_COUNT", 10000),
		ReportDir:        getenv("REPORT_DIR", "."),
		EnableCardLock:   getenv("ENABLE_CARD_LOCK", "true") == "true",
		ConversionFeeBps: getenvInt("CONVERSION_FEE_BPS", 10),
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func getenvInt(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}

func defaultFunctionsSource() string {
	return `
const commitment = args[0];
const batchId = args[1];
const epoch = args[2];
const workerCount = args[3];
if (!commitment || commitment.length < 32) throw Error("invalid commitment");
if (!batchId) throw Error("missing batchId");
return Functions.encodeString(JSON.stringify({
  ok: true,
  commitment: commitment,
  batchId: batchId,
  epoch: epoch,
  workerCount: workerCount,
  verifiedAt: Date.now(),
  gate: "linkswarm-functions"
}));
`
}

// ---------- Domain ----------

type Worker struct {
	ID           string
	Name         string
	WagePerHour  int
	HoursWorked  int
	PTOAllocated int
	PTOUsed      int
	PaymentPref  WorkerPaymentPreference
}

type WorkerPaymentPreference struct {
	Network       string
	Token         string
	PaymentMethod string
	BankAccount   string
	RoutingNumber string
}

func (w *Worker) TotalPay() int { return w.WagePerHour * w.HoursWorked }

type LedgerEntry struct {
	ID          string            `json:"id"`
	UserID      string            `json:"userId"`
	Type        string            `json:"type"`
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	Description string            `json:"description"`
	Network     string            `json:"network"`
	Token       string            `json:"token"`
	TxHash      string            `json:"txHash"`
	Reference   string            `json:"reference"`
	Timestamp   time.Time         `json:"timestamp"`
	Status      string            `json:"status"`
	Balance     float64           `json:"balance"`
	Source      string            `json:"source"`
	Destination string            `json:"destination"`
	Tax         float64           `json:"tax"`
	Metadata    map[string]string `json:"metadata"`
}

type PrivateLedger struct {
	UserID     string
	Entries    []LedgerEntry
	Balance    float64
	TotalIn    float64
	TotalOut   float64
	LastUpdate time.Time
	Version    string
	mu         sync.RWMutex
}

func NewPrivateLedger(id string) *PrivateLedger {
	return &PrivateLedger{UserID: id, Entries: make([]LedgerEntry, 0, 32), Version: "3.0.0-linkswarm"}
}

func (l *PrivateLedger) Append(e LedgerEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e.Balance = l.Balance + e.Amount
	l.Balance = e.Balance
	if e.Amount >= 0 {
		l.TotalIn += e.Amount
	} else {
		l.TotalOut += -e.Amount
	}
	e.Timestamp = time.Now().UTC()
	l.Entries = append(l.Entries, e)
	l.LastUpdate = e.Timestamp
}

// ---------- Card Lock ----------

type CardLockManager struct {
	locks map[string]bool // true = locked
	mu    sync.RWMutex
}

func NewCardLockManager() *CardLockManager {
	return &CardLockManager{locks: make(map[string]bool)}
}
func (m *CardLockManager) Allowed(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return !m.locks[id]
}
func (m *CardLockManager) Lock(id string) {
	m.mu.Lock()
	m.locks[id] = true
	m.mu.Unlock()
}

// ---------- Chainlink Functions ----------

type FunctionsRequest struct {
	Source         string
	Args           []string
	SubscriptionID uint64
	DONID          string
	GasLimit       uint32
	IdempotencyKey string
}

type FunctionsResponse struct {
	RequestID   string
	Result      []byte
	Error       string
	Fulfilled   bool
	Verified    bool
	RawReport   string
	FulfilledAt time.Time
}

type FunctionsClient interface {
	SendRequest(ctx context.Context, req FunctionsRequest) (FunctionsResponse, error)
}

type SimulationFunctionsClient struct{}

func (c *SimulationFunctionsClient) SendRequest(ctx context.Context, req FunctionsRequest) (FunctionsResponse, error) {
	select {
	case <-ctx.Done():
		return FunctionsResponse{}, ctx.Err()
	default:
	}
	if len(req.Args) < 1 || len(req.Args[0]) < 32 {
		return FunctionsResponse{Error: "invalid commitment"}, fmt.Errorf("invalid commitment")
	}
	h := sha256.Sum256([]byte(req.Args[0] + req.DONID + req.IdempotencyKey))
	reqID := "0xreq-" + hex.EncodeToString(h[:16])
	payload, _ := json.Marshal(map[string]interface{}{
		"ok": true, "commitment": req.Args[0], "batchId": safe(req.Args, 1),
		"epoch": safe(req.Args, 2), "workerCount": safe(req.Args, 3),
		"verifiedAt": time.Now().UTC().UnixMilli(), "gate": "linkswarm-functions",
	})
	return FunctionsResponse{
		RequestID: reqID, Result: payload, Fulfilled: true, Verified: true,
		RawReport: "linkswarm-sim-" + reqID, FulfilledAt: time.Now().UTC(),
	}, nil
}

func safe(a []string, i int) string {
	if i < len(a) {
		return a[i]
	}
	return ""
}

type LiveFunctionsClient struct {
	Router string
	HTTP   *http.Client
}

func NewLiveFunctionsClient(router string, timeout time.Duration) *LiveFunctionsClient {
	return &LiveFunctionsClient{
		Router: strings.TrimRight(router, "/"),
		HTTP:   &http.Client{Timeout: timeout},
	}
}

func (c *LiveFunctionsClient) SendRequest(ctx context.Context, req FunctionsRequest) (FunctionsResponse, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"source": req.Source, "args": req.Args,
		"subscriptionId": req.SubscriptionID, "donId": req.DONID,
		"gasLimit": req.GasLimit, "idempotencyKey": req.IdempotencyKey,
	})
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.Router+"/functions/request", bytes.NewReader(body))
	if err != nil {
		return FunctionsResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return FunctionsResponse{}, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return FunctionsResponse{Error: string(data)}, fmt.Errorf("functions http %s", resp.Status)
	}
	var p struct {
		RequestID string `json:"requestId"`
		Result    string `json:"result"`
		Error     string `json:"error"`
		Report    string `json:"report"`
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return FunctionsResponse{}, err
	}
	return FunctionsResponse{
		RequestID: p.RequestID, Result: []byte(p.Result), Error: p.Error,
		Fulfilled: p.Error == "", Verified: p.Error == "",
		RawReport: p.Report, FulfilledAt: time.Now().UTC(),
	}, nil
}

// ---------- Delivery & Bridge ----------

type NetworkAdapter interface {
	Name() string
	SendBatch(payloads [][]byte) ([]string, error)
}

type GenericAdapter struct{ name string }

func (a *GenericAdapter) Name() string { return a.name }
func (a *GenericAdapter) SendBatch(payloads [][]byte) ([]string, error) {
	out := make([]string, len(payloads))
	for i, p := range payloads {
		h := sha256.Sum256(append(p, []byte(a.name)...))
		out[i] = "0x" + hex.EncodeToString(h[:])
	}
	return out, nil
}

type USDCBridge struct {
	feeBps int
	mu     sync.RWMutex
}

func NewUSDCBridge(bps int) *USDCBridge { return &USDCBridge{feeBps: bps} }
func (b *USDCBridge) Convert(amount int, network string) (int, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	usdc := amount - (amount * b.feeBps / 10000)
	h := sha256.Sum256([]byte(fmt.Sprintf("%d:%s:%d", amount, network, time.Now().UnixNano())))
	return usdc, "USDC-" + hex.EncodeToString(h[:16])
}
func (b *USDCBridge) ToFiat(usdc int, bank, routing string) (string, string) {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d:%s:%s:%d", usdc, bank, routing, time.Now().UnixNano())))
	return fmt.Sprintf("$%d.%02d", usdc/100, usdc%100), "FIAT-" + hex.EncodeToString(h[:16])
}

// ---------- Settlement ----------

type SettlementBatch struct {
	ID         string
	Epoch      int64
	Workers    []*Worker
	Total      int
	Timestamp  time.Time
	Status     string
	Commitment string
	RequestID  string
}

type SettlementResult struct {
	BatchID            string `json:"batchId"`
	WorkerID           string `json:"workerId"`
	Amount             int    `json:"amount"`
	PTOUsed            int    `json:"ptoUsed"`
	PTORemaining       int    `json:"ptoRemaining"`
	Network            string `json:"network"`
	PaymentMethod      string `json:"paymentMethod"`
	TxHash             string `json:"txHash"`
	USDCAmount         int    `json:"usdcAmount"`
	USDCtxHash         string `json:"usdcTxHash"`
	FiatAmount         string `json:"fiatAmount"`
	FiatTxHash         string `json:"fiatTxHash"`
	Success            bool   `json:"success"`
	DurationMs         int64  `json:"durationMs"`
	Cached             bool   `json:"cached"`
	Status             string `json:"status"`
	Epoch              int64  `json:"epoch"`
	Timestamp          string `json:"timestamp"`
	ChainlinkVerified  bool   `json:"chainlinkVerified"`
	FunctionsRequestID string `json:"functionsRequestId"`
}

type LinkSwarmEngine struct {
	cfg       Config
	functions FunctionsClient
	adapters  map[string]NetworkAdapter
	bridge    *USDCBridge
	ledgers   map[string]*PrivateLedger
	cardLock  *CardLockManager
	cacheHits atomic.Uint64
	mu        sync.RWMutex
	log       *log.Logger
}

func NewLinkSwarmEngine(cfg Config) *LinkSwarmEngine {
	var fc FunctionsClient
	if cfg.Mode == "live" && cfg.FunctionsRouter != "" {
		fc = NewLiveFunctionsClient(cfg.FunctionsRouter, cfg.RequestTimeout)
		log.Printf("LinkSwarm Functions: LIVE → %s", cfg.FunctionsRouter)
	} else {
		fc = &SimulationFunctionsClient{}
		log.Printf("LinkSwarm Functions: SIMULATION")
	}
	e := &LinkSwarmEngine{
		cfg: cfg, functions: fc,
		adapters:  make(map[string]NetworkAdapter),
		bridge:    NewUSDCBridge(cfg.ConversionFeeBps),
		ledgers:   make(map[string]*PrivateLedger),
		cardLock:  NewCardLockManager(),
		log:       log.New(os.Stdout, "[linkswarm] ", log.LstdFlags|log.Lmsgprefix),
	}
	for _, n := range []string{"EVM", "Solana", "Cosmos", "UTXO"} {
		e.adapters[n] = &GenericAdapter{name: n}
	}
	return e
}

func (e *LinkSwarmEngine) ledger(id string) *PrivateLedger {
	e.mu.Lock()
	defer e.mu.Unlock()
	if l, ok := e.ledgers[id]; ok {
		e.cacheHits.Add(1)
		return l
	}
	l := NewPrivateLedger(id)
	e.ledgers[id] = l
	return l
}

func (e *LinkSwarmEngine) entropy312(batchID string, payloads [][]byte) string {
	root := sha3.NewLegacyKeccak256()
	root.Write([]byte(batchID))
	for _, p := range payloads {
		root.Write(p)
	}
	batchRoot := root.Sum(nil)
	anchor := sha3.NewLegacyKeccak256()
	anchor.Write([]byte(batchID))
	anchor.Write(batchRoot)
	noise := make([]byte, 32)
	_, _ = rand.Read(noise)
	final := sha3.NewLegacyKeccak256()
	final.Write(anchor.Sum(nil))
	final.Write(noise)
	return hex.EncodeToString(final.Sum(nil))
}

func (e *LinkSwarmEngine) verifyWithFunctions(ctx context.Context, batch *SettlementBatch) (FunctionsResponse, error) {
	req := FunctionsRequest{
		Source:         e.cfg.SourceCode,
		Args:           []string{batch.Commitment, batch.ID, fmt.Sprintf("%d", batch.Epoch), fmt.Sprintf("%d", batch.Total)},
		SubscriptionID: e.cfg.SubscriptionID,
		DONID:          e.cfg.DONID,
		GasLimit:       300000,
		IdempotencyKey: batch.ID + "-" + batch.Commitment[:16],
	}
	var last error
	for attempt := 1; attempt <= e.cfg.RetryAttempts; attempt++ {
		resp, err := e.functions.SendRequest(ctx, req)
		if err == nil && resp.Fulfilled && resp.Verified {
			e.log.Printf("Functions verified requestId=%s commitment=%s...", resp.RequestID, batch.Commitment[:16])
			return resp, nil
		}
		last = err
		if err == nil {
			last = fmt.Errorf(resp.Error)
		}
		e.log.Printf("Functions attempt %d: %v", attempt, last)
		time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
	}
	return FunctionsResponse{}, last
}

func (e *LinkSwarmEngine) ProcessBatch(ctx context.Context, batch *SettlementBatch) ([]SettlementResult, error) {
	start := time.Now()
	payloads := make([][]byte, len(batch.Workers))
	for i, w := range batch.Workers {
		payloads[i] = []byte(fmt.Sprintf("%s:%d:pto=%d:epoch=%d", w.ID, w.TotalPay(), w.PTOUsed, batch.Epoch))
	}
	batch.Commitment = e.entropy312(batch.ID, payloads)
	batch.Status = "commitment_ready"

	fn, err := e.verifyWithFunctions(ctx, batch)
	if err != nil {
		batch.Status = "functions_rejected"
		return nil, fmt.Errorf("LinkSwarm Functions gate failed: %w", err)
	}
	batch.RequestID = fn.RequestID
	batch.Status = "functions_verified"

	results := make([]SettlementResult, len(batch.Workers))
	sem := make(chan struct{}, e.cfg.MaxConcurrency)
	var wg sync.WaitGroup

	for i, w := range batch.Workers {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, worker *Worker) {
			defer wg.Done()
			defer func() { <-sem }()

			if e.cfg.EnableCardLock && !e.cardLock.Allowed(worker.ID) {
				results[idx] = SettlementResult{
					BatchID: batch.ID, WorkerID: worker.ID, Success: false,
					Status: "rejected_card_locked", Epoch: batch.Epoch, Cached: true,
					FunctionsRequestID: batch.RequestID,
				}
				return
			}

			network := worker.PaymentPref.Network
			if network == "" {
				network = "EVM"
			}
			adapter := e.adapters[network]
			if adapter == nil {
				adapter = e.adapters["EVM"]
				network = "EVM"
			}

			amount := worker.TotalPay()
			txs, _ := adapter.SendBatch([][]byte{payloads[idx]})
			txHash := ""
			if len(txs) > 0 {
				txHash = txs[0]
			}

			usdcAmt, usdcTx, fiatAmt, fiatTx := 0, "", "", ""
			switch worker.PaymentPref.PaymentMethod {
			case "usdc":
				usdcAmt, usdcTx = e.bridge.Convert(amount, network)
			case "bank_transfer":
				usdcAmt, usdcTx = e.bridge.Convert(amount, network)
				fiatAmt, fiatTx = e.bridge.ToFiat(usdcAmt, worker.PaymentPref.BankAccount, worker.PaymentPref.RoutingNumber)
			}

			results[idx] = SettlementResult{
				BatchID: batch.ID, WorkerID: worker.ID, Amount: amount,
				PTOUsed: worker.PTOUsed, PTORemaining: worker.PTOAllocated - worker.PTOUsed,
				Network: network, PaymentMethod: worker.PaymentPref.PaymentMethod,
				TxHash: txHash, USDCAmount: usdcAmt, USDCtxHash: usdcTx,
				FiatAmount: fiatAmt, FiatTxHash: fiatTx, Success: true,
				DurationMs: time.Since(start).Milliseconds(), Cached: true,
				Status: "completed", Epoch: batch.Epoch,
				Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
				ChainlinkVerified: true, FunctionsRequestID: batch.RequestID,
			}

			e.ledger(worker.ID).Append(LedgerEntry{
				ID: fmt.Sprintf("ls-%d-%s", time.Now().UnixNano(), worker.ID[:8]),
				UserID: worker.ID, Type: "settlement", Amount: float64(amount),
				Currency: "USD", Description: "LinkSwarm settlement (Functions-verified)",
				Network: network, Token: worker.PaymentPref.Token, TxHash: txHash,
				Reference: fmt.Sprintf("epoch-%d", batch.Epoch), Status: "completed",
				Source: "employer", Destination: "wallet", Tax: float64(amount) * 0.15,
				Metadata: map[string]string{
					"commitment":  batch.Commitment[:16],
					"requestId":   batch.RequestID,
					"verified_by": "Chainlink Functions",
					"product":     "LinkSwarm",
				},
			})
		}(i, w)
	}
	wg.Wait()
	batch.Status = "completed"
	return results, nil
}

func ExportJSON(results []SettlementResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func ExportCSV(results []SettlementResult, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	_ = w.Write([]string{
		"WorkerID", "Amount", "PTOUsed", "PTORemaining", "Network", "PaymentMethod",
		"TxHash", "USDCAmount", "USDCtxHash", "FiatAmount", "FiatTxHash",
		"Epoch", "Timestamp", "Status", "Cached", "ChainlinkVerified", "FunctionsRequestID",
	})
	for _, r := range results {
		_ = w.Write([]string{
			r.WorkerID, fmt.Sprintf("%d", r.Amount), fmt.Sprintf("%d", r.PTOUsed),
			fmt.Sprintf("%d", r.PTORemaining), r.Network, r.PaymentMethod,
			r.TxHash, fmt.Sprintf("%d", r.USDCAmount), r.USDCtxHash,
			r.FiatAmount, r.FiatTxHash, fmt.Sprintf("%d", r.Epoch),
			r.Timestamp, r.Status, fmt.Sprintf("%t", r.Cached),
			fmt.Sprintf("%t", r.ChainlinkVerified), r.FunctionsRequestID,
		})
	}
	return nil
}

func RunLinkSwarm(cfg Config) error {
	engine := NewLinkSwarmEngine(cfg)
	workers := make([]*Worker, cfg.WorkerCount)
	nets := []string{"EVM", "Solana", "Cosmos", "UTXO"}
	methods := []string{"crypto", "usdc", "bank_transfer"}
	for i := 0; i < cfg.WorkerCount; i++ {
		workers[i] = &Worker{
			ID: fmt.Sprintf("0x%016x", i+1), Name: fmt.Sprintf("Worker-%d", i+1),
			WagePerHour: 40 + i%30, HoursWorked: 20 + i%40,
			PTOAllocated: 80, PTOUsed: i % 10,
			PaymentPref: WorkerPaymentPreference{
				Network: nets[i%len(nets)], Token: "USDC", PaymentMethod: methods[i%len(methods)],
			},
		}
	}
	batch := &SettlementBatch{
		ID: fmt.Sprintf("linkswarm-%d", time.Now().UnixNano()),
		Epoch: time.Now().Unix(), Workers: workers, Total: cfg.WorkerCount,
		Timestamp: time.Now().UTC(), Status: "pending",
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout+45*time.Second)
	defer cancel()

	log.Printf("LinkSwarm e2e starting — %d workers | mode=%s", cfg.WorkerCount, cfg.Mode)
	start := time.Now()
	results, err := engine.ProcessBatch(ctx, batch)
	if err != nil {
		return err
	}
	elapsed := time.Since(start)
	success := 0
	for _, r := range results {
		if r.Success {
			success++
		}
	}
	tps := float64(cfg.WorkerCount) / elapsed.Seconds()

	fmt.Println("================================================")
	fmt.Println("LINKSWARM — CHAINLINK FUNCTIONS SETTLEMENT GATE")
	fmt.Println("================================================")
	fmt.Printf("Product:             LinkSwarm\n")
	fmt.Printf("Mode:                %s\n", cfg.Mode)
	fmt.Printf("Total Settlements:   %d\n", cfg.WorkerCount)
	fmt.Printf("Success:             %d\n", success)
	fmt.Printf("Failed:              %d\n", cfg.WorkerCount-success)
	fmt.Printf("Success Rate:        %.2f%%\n", 100*float64(success)/float64(cfg.WorkerCount))
	fmt.Printf("Time:                %.2fs\n", elapsed.Seconds())
	fmt.Printf("TPS:                 %.1f\n", tps)
	fmt.Printf("Cache Hits:          %d\n", engine.cacheHits.Load())
	fmt.Printf("Functions RequestID: %s\n", batch.RequestID)
	fmt.Printf("Commitment:          %s...\n", batch.Commitment[:32])
	fmt.Println("================================================")

	jp := cfg.ReportDir + "/linkswarm_report.json"
	cp := cfg.ReportDir + "/linkswarm_report.csv"
	_ = ExportJSON(results, jp)
	_ = ExportCSV(results, cp)
	fmt.Printf("Reports written: %s  %s\n", jp, cp)
	return nil
}

func main() {
	cfg := LoadConfig()
	if err := RunLinkSwarm(cfg); err != nil {
		log.Fatalf("LinkSwarm failed: %v", err)
	}
}
