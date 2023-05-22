package s4

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Envelope represents a JSON object that is signed for address verification.
// All []byte values are encoded as base64 (default JSON behavior).
// Hex is not used to avoid confusion due to case-sensivity and 0x prefix.
// A signer is responsible for generating a JSON that has no whitespace and
// the keys appear in this exact order:
// {"address":base64,"slotid":int,"payload":base64,"version":int,"expiration":int}
type Envelope struct {
	Address    []byte `json:"address"`
	SlotID     uint   `json:"slotid"`
	Payload    []byte `json:"payload"`
	Version    uint64 `json:"version"`
	Expiration int64  `json:"expiration"`
}

func NewEnvelopeFromRecord(key *Key, record *Record) *Envelope {
	return &Envelope{
		Address:    key.Address.Bytes(),
		SlotID:     key.SlotId,
		Payload:    record.Payload,
		Version:    key.Version,
		Expiration: record.Expiration,
	}
}

// Sign calculates signature for the serialized envelope data.
func (e Envelope) Sign(privateKey *ecdsa.PrivateKey) (signature []byte, err error) {
	js, err := e.ToJson()
	if err != nil {
		return nil, err
	}
	hash := crypto.Keccak256Hash(js)
	return crypto.Sign(hash[:], privateKey)
}

// GetSignerAddress verifies the signature and returns the signing address.
func (e Envelope) GetSignerAddress(signature []byte) (address common.Address, err error) {
	js, err := e.ToJson()
	if err != nil {
		return common.Address{}, err
	}
	hash := crypto.Keccak256Hash(js)
	sigPublicKey, err := crypto.SigToPub(hash[:], signature)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*sigPublicKey), nil
}

func (e Envelope) ToJson() ([]byte, error) {
	address, err := json.Marshal(e.Address)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return nil, err
	}
	js := fmt.Sprintf(`{"address":%s,"slotid":%d,"payload":%s,"version":%d,"expiration":%d}`, address, e.SlotID, payload, e.Version, e.Expiration)
	return []byte(js), nil
}
