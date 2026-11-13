package agentcommerce

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	minHashLength = 16
	maxHashLength = 64
)

// Agent models the end-to-end Agent Commerce Protocol lifecycle.
type Agent interface {
	Discover(context.Context, ServiceQuery) ([]AgentProfile, error)
	Negotiate(context.Context, AgentProfile, ServiceRequest) (IntentTerms, error)
	Execute(context.Context, SignedIntent) (Proof, error)
	Verify(context.Context, SignedIntent, Proof) (VerificationResult, error)
	Settle(context.Context, SignedIntent, VerificationResult) (SettlementReceipt, error)
}

// Wallet signs and verifies protocol payloads.
type Wallet interface {
	Address() string
	Sign([]byte) ([]byte, error)
	Verify(address string, payload, signature []byte) bool
}

type Escrow interface {
	Lock(context.Context, SignedIntent) (EscrowReceipt, error)
	Release(context.Context, string) (SettlementReceipt, error)
	Refund(context.Context, string) (SettlementReceipt, error)
}

type Settlement interface {
	Send(context.Context, SettlementInstruction) (SettlementReceipt, error)
	Receive(context.Context, SettlementReceipt) error
}

type Reputation interface {
	Update(context.Context, ReputationEvent) error
	Query(context.Context, string) (ReputationScore, error)
}

type AgentProfile struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Capabilities []string          `json:"capabilities"`
	Pricing      map[string]Amount `json:"pricing"`
	Reputation   ReputationScore   `json:"reputation"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

type Amount struct {
	Value    int64  `json:"value"`
	Currency string `json:"currency"`
	Rail     string `json:"rail"`
}

type ServiceQuery struct {
	Capability       string
	MaxPrice         Amount
	MinReputation    int64
	SettlementRail   string
	RequiredMetadata map[string]string
}

type ServiceRequest struct {
	Service         string            `json:"service"`
	SLA             string            `json:"sla"`
	Deliverables    []string          `json:"deliverables"`
	SettlementChain string            `json:"settlement_chain"`
	EscrowTerms     EscrowTerms       `json:"escrow_terms"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

type EscrowTerms struct {
	ReleaseConditions []string      `json:"release_conditions"`
	Expiration        time.Duration `json:"expiration"`
}

type IntentTerms struct {
	ServiceRequest
	Price     Amount    `json:"price"`
	Buyer     string    `json:"buyer"`
	Seller    string    `json:"seller"`
	Timestamp time.Time `json:"timestamp"`
}

type SignedIntent struct {
	Terms           IntentTerms `json:"terms"`
	Hash            string      `json:"hash"`
	BuyerSignature  []byte      `json:"buyer_signature"`
	SellerSignature []byte      `json:"seller_signature"`
}

type Proof struct {
	OutputHash  string            `json:"output_hash"`
	Timestamp   time.Time         `json:"timestamp"`
	ExecutionID string            `json:"execution_id"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type VerificationResult struct {
	Verified bool      `json:"verified"`
	Method   string    `json:"method"`
	Reason   string    `json:"reason,omitempty"`
	Time     time.Time `json:"time"`
}

type EscrowReceipt struct {
	EscrowID   string    `json:"escrow_id"`
	IntentHash string    `json:"intent_hash"`
	Buyer      string    `json:"buyer"`
	Seller     string    `json:"seller"`
	Amount     Amount    `json:"amount"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type SettlementInstruction struct {
	EscrowID string
	To       string
	Amount   Amount
	Rail     string
}

type SettlementReceipt struct {
	SettlementID string    `json:"settlement_id"`
	EscrowID     string    `json:"escrow_id"`
	To           string    `json:"to"`
	Amount       Amount    `json:"amount"`
	Rail         string    `json:"rail"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

type ReputationEvent struct {
	AgentID    string
	IntentHash string
	Delta      int64
	Reason     string
}

type ReputationScore struct {
	AgentID    string `json:"agent_id"`
	Successful int64  `json:"successful"`
	Disputed   int64  `json:"disputed"`
	Score      int64  `json:"score"`
}

type AuditEntry struct {
	Kind         string    `json:"kind"`
	Ref          string    `json:"ref"`
	Payload      []byte    `json:"payload"`
	PayloadHash  string    `json:"payload_hash"`
	PreviousHash string    `json:"previous_hash,omitempty"`
	EntryHash    string    `json:"entry_hash"`
	Timestamp    time.Time `json:"timestamp"`
}

type Policy struct {
	MaxSpend              int64
	AllowedCurrencies     []string
	AllowedRails          []string
	MinSellerReputation   int64
	RequireWalletApproval bool
}

func IntentHash(terms IntentTerms) (string, error) {
	if err := ValidateIntentTerms(terms); err != nil {
		return "", err
	}
	payload, err := canonicalJSON(terms)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func ValidateIntentTerms(terms IntentTerms) error {
	if terms.Service == "" {
		return errors.New("service is required")
	}
	if terms.Price.Value <= 0 {
		return errors.New("price must be positive")
	}
	if terms.Price.Currency == "" {
		return errors.New("currency is required")
	}
	if terms.Price.Rail == "" {
		return errors.New("settlement rail is required")
	}
	if terms.Buyer == "" {
		return errors.New("buyer is required")
	}
	if terms.Seller == "" {
		return errors.New("seller is required")
	}
	if terms.Buyer == terms.Seller {
		return errors.New("buyer and seller must differ")
	}
	if terms.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}
	if terms.EscrowTerms.Expiration <= 0 {
		return errors.New("escrow expiration must be positive")
	}
	if len(terms.EscrowTerms.ReleaseConditions) == 0 {
		return errors.New("at least one release condition is required")
	}
	return nil
}

func SignIntent(terms IntentTerms, buyer, seller Wallet) (SignedIntent, error) {
	if buyer == nil {
		return SignedIntent{}, errors.New("buyer wallet is required")
	}
	if seller == nil {
		return SignedIntent{}, errors.New("seller wallet is required")
	}
	if err := ValidateIntentTerms(terms); err != nil {
		return SignedIntent{}, err
	}
	hash, err := IntentHash(terms)
	if err != nil {
		return SignedIntent{}, err
	}
	if hash == "" {
		return SignedIntent{}, errors.New("generated empty hash")
	}
	payload := []byte(hash)
	bs, err := buyer.Sign(payload)
	if err != nil {
		return SignedIntent{}, err
	}
	ss, err := seller.Sign(payload)
	if err != nil {
		return SignedIntent{}, err
	}
	if len(bs) == 0 {
		return SignedIntent{}, errors.New("buyer signature is empty")
	}
	if len(ss) == 0 {
		return SignedIntent{}, errors.New("seller signature is empty")
	}
	return SignedIntent{
		Terms:           terms,
		Hash:            hash,
		BuyerSignature:  bs,
		SellerSignature: ss,
	}, nil
}

func VerifySignedIntent(intent SignedIntent, wallet Wallet) error {
	if wallet == nil {
		return errors.New("wallet is required")
	}
	hash, err := IntentHash(intent.Terms)
	if err != nil {
		return err
	}
	if hash != intent.Hash {
		return errors.New("intent hash mismatch")
	}
	payload := []byte(intent.Hash)
	if !wallet.Verify(intent.Terms.Buyer, payload, intent.BuyerSignature) {
		return errors.New("invalid buyer signature")
	}
	if !wallet.Verify(intent.Terms.Seller, payload, intent.SellerSignature) {
		return errors.New("invalid seller signature")
	}
	return nil
}

func EvaluatePolicy(policy Policy, seller AgentProfile, terms IntentTerms) error {
	if err := ValidateIntentTerms(terms); err != nil {
		return err
	}
	if policy.MaxSpend > 0 && terms.Price.Value > policy.MaxSpend {
		return fmt.Errorf("price %d exceeds max spend %d", terms.Price.Value, policy.MaxSpend)
	}
	if len(policy.AllowedCurrencies) > 0 && !contains(policy.AllowedCurrencies, terms.Price.Currency) {
		return fmt.Errorf("currency %s not allowed", terms.Price.Currency)
	}
	if len(policy.AllowedRails) > 0 && !contains(policy.AllowedRails, terms.Price.Rail) {
		return fmt.Errorf("rail %s not allowed", terms.Price.Rail)
	}
	if seller.Reputation.Score < policy.MinSellerReputation {
		return fmt.Errorf("seller reputation %d below threshold %d", seller.Reputation.Score, policy.MinSellerReputation)
	}
	for key, want := range terms.Metadata {
		if got, ok := seller.Metadata[key]; ok && want != "" && got != want {
			return fmt.Errorf("seller metadata %s=%s does not match required value %s", key, got, want)
		}
	}
	return nil
}

type Ed25519Wallet struct {
	address string
	priv    ed25519.PrivateKey
	pubkeys *sync.Map
}

func NewEd25519Wallet(registry *sync.Map) (*Ed25519Wallet, error) {
	if registry == nil {
		return nil, errors.New("wallet registry is required")
	}
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	sum := sha256.Sum256(pub)
	address := hex.EncodeToString(sum[:20])
	if address == "" {
		return nil, errors.New("generated empty address")
	}
	registry.Store(address, pub)
	return &Ed25519Wallet{
		address: address,
		priv:    priv,
		pubkeys: registry,
	}, nil
}

func (w *Ed25519Wallet) Address() string { return w.address }

func (w *Ed25519Wallet) Sign(payload []byte) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errors.New("payload is empty")
	}
	return ed25519.Sign(w.priv, payload), nil
}

func (w *Ed25519Wallet) Verify(address string, payload, signature []byte) bool {
	if address == "" || len(payload) == 0 || len(signature) == 0 {
		return false
	}
	v, ok := w.pubkeys.Load(address)
	if !ok {
		return false
	}
	pubKey, ok := v.(ed25519.PublicKey)
	if !ok || pubKey == nil {
		return false
	}
	if len(pubKey) != ed25519.PublicKeySize {
		return false
	}
	return ed25519.Verify(pubKey, payload, signature)
}

// SafeWallet wraps a Wallet with concurrency protection.
type SafeWallet struct {
	wallet Wallet
	mu     sync.RWMutex
}

func NewSafeWallet(w Wallet) *SafeWallet {
	return &SafeWallet{wallet: w}
}

func (s *SafeWallet) Address() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.wallet == nil {
		return ""
	}
	return s.wallet.Address()
}

func (s *SafeWallet) Sign(payload []byte) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.wallet == nil {
		return nil, errors.New("wallet is nil")
	}
	if len(payload) == 0 {
		return nil, errors.New("payload is empty")
	}
	return s.wallet.Sign(payload)
}

func (s *SafeWallet) Verify(address string, payload, signature []byte) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.wallet == nil {
		return false
	}
	return s.wallet.Verify(address, payload, signature)
}

type InMemoryDirectory struct {
	mu       sync.RWMutex
	profiles map[string]AgentProfile
}

func NewInMemoryDirectory() *InMemoryDirectory {
	return &InMemoryDirectory{profiles: map[string]AgentProfile{}}
}

func (d *InMemoryDirectory) Register(profile AgentProfile) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.profiles[profile.ID] = profile
}

func (d *InMemoryDirectory) Discover(q ServiceQuery) []AgentProfile {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var out []AgentProfile
	for _, p := range d.profiles {
		price, ok := p.Pricing[q.Capability]
		if !ok || !contains(p.Capabilities, q.Capability) {
			continue
		}
		if q.MaxPrice.Value > 0 && price.Value > q.MaxPrice.Value {
			continue
		}
		if q.SettlementRail != "" && price.Rail != q.SettlementRail {
			continue
		}
		if p.Reputation.Score < q.MinReputation {
			continue
		}
		if !metadataMatches(p.Metadata, q.RequiredMetadata) {
			continue
		}
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Reputation.Score > out[j].Reputation.Score })
	return out
}

type InMemoryEscrow struct {
	mu       sync.Mutex
	locked   map[string]EscrowReceipt
	released map[string]bool
}

func NewInMemoryEscrow() *InMemoryEscrow {
	return &InMemoryEscrow{locked: map[string]EscrowReceipt{}, released: map[string]bool{}}
}

func (e *InMemoryEscrow) Lock(_ context.Context, intent SignedIntent) (EscrowReceipt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := ValidateIntentTerms(intent.Terms); err != nil {
		return EscrowReceipt{}, err
	}
	if len(intent.Hash) < minHashLength {
		return EscrowReceipt{}, errors.New("intent hash too short")
	}
	if _, ok := e.locked[intent.Hash]; ok {
		return EscrowReceipt{}, errors.New("escrow already locked")
	}
	escrowID := "escrow-" + intent.Hash[:minHashLength]
	r := EscrowReceipt{
		EscrowID:   escrowID,
		IntentHash: intent.Hash,
		Buyer:      intent.Terms.Buyer,
		Seller:     intent.Terms.Seller,
		Amount:     intent.Terms.Price,
		ExpiresAt:  intent.Terms.Timestamp.Add(intent.Terms.EscrowTerms.Expiration),
	}
	e.locked[intent.Hash] = r
	return r, nil
}

func (e *InMemoryEscrow) Release(_ context.Context, escrowID string) (SettlementReceipt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	r, err := e.find(escrowID)
	if err != nil {
		return SettlementReceipt{}, err
	}
	if e.released[escrowID] {
		return SettlementReceipt{}, errors.New("escrow already settled")
	}
	e.released[escrowID] = true
	return SettlementReceipt{
		SettlementID: "settle-" + r.IntentHash[:minHashLength],
		EscrowID:     escrowID,
		To:           r.Seller,
		Amount:       r.Amount,
		Rail:         r.Amount.Rail,
		Status:       "released",
		Timestamp:    time.Now().UTC(),
	}, nil
}

func (e *InMemoryEscrow) Refund(_ context.Context, escrowID string) (SettlementReceipt, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	r, err := e.find(escrowID)
	if err != nil {
		return SettlementReceipt{}, err
	}
	if e.released[escrowID] {
		return SettlementReceipt{}, errors.New("escrow already settled")
	}
	e.released[escrowID] = true
	return SettlementReceipt{
		SettlementID: "refund-" + r.IntentHash[:minHashLength],
		EscrowID:     escrowID,
		To:           r.Buyer,
		Amount:       r.Amount,
		Rail:         r.Amount.Rail,
		Status:       "refunded",
		Timestamp:    time.Now().UTC(),
	}, nil
}

func (e *InMemoryEscrow) find(id string) (EscrowReceipt, error) {
	for _, r := range e.locked {
		if r.EscrowID == id {
			return r, nil
		}
	}
	return EscrowReceipt{}, errors.New("escrow not found")
}

type InMemoryReputation struct {
	mu     sync.Mutex
	scores map[string]ReputationScore
}

func NewInMemoryReputation() *InMemoryReputation {
	return &InMemoryReputation{scores: map[string]ReputationScore{}}
}

func (r *InMemoryReputation) Update(_ context.Context, ev ReputationEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	s := r.scores[ev.AgentID]
	s.AgentID = ev.AgentID
	if ev.Delta >= 0 {
		s.Successful += ev.Delta
	} else {
		s.Disputed -= ev.Delta
	}
	s.Score += ev.Delta
	r.scores[ev.AgentID] = s
	return nil
}

func (r *InMemoryReputation) Query(_ context.Context, id string) (ReputationScore, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.scores[id], nil
}

type AuditLog struct {
	mu      sync.Mutex
	entries []AuditEntry
}

func (l *AuditLog) Store(kind, ref string, payload any) error {
	b, err := canonicalJSON(payload)
	if err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	payloadSum := sha256.Sum256(b)
	previousHash := ""
	if len(l.entries) > 0 {
		previousHash = l.entries[len(l.entries)-1].EntryHash
	}
	timestamp := time.Now().UTC()
	entryPayload, err := canonicalJSON(struct {
		Kind         string
		Ref          string
		PayloadHash  string
		PreviousHash string
		Timestamp    time.Time
	}{Kind: kind, Ref: ref, PayloadHash: hex.EncodeToString(payloadSum[:]), PreviousHash: previousHash, Timestamp: timestamp})
	if err != nil {
		return err
	}
	entrySum := sha256.Sum256(entryPayload)
	l.entries = append(l.entries, AuditEntry{
		Kind:         kind,
		Ref:          ref,
		Payload:      b,
		PayloadHash:  hex.EncodeToString(payloadSum[:]),
		PreviousHash: previousHash,
		EntryHash:    hex.EncodeToString(entrySum[:]),
		Timestamp:    timestamp,
	})
	return nil
}

func (l *AuditLog) Entries() []AuditEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]AuditEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

func canonicalJSON(v any) ([]byte, error) { return json.Marshal(v) }

func contains(values []string, want string) bool {
	for _, v := range values {
		if v == want {
			return true
		}
	}
	return false
}

func metadataMatches(have, want map[string]string) bool {
	for k, v := range want {
		if have[k] != v {
			return false
		}
	}
	return true
}
