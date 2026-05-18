package hd

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type SovereignHD struct {
	Mnemonic string
	Root     *bip32.Key
}

func New(mnemonic string) (*SovereignHD, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, "")
	root, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	return &SovereignHD{Mnemonic: mnemonic, Root: root}, nil
}

func (h *SovereignHD) DeriveReceiveOnly() (string, error) {
	k := h.Root
	k, _ = k.NewChildKey(402)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	privKey, _ := crypto.ToECDSA(k.Key)
	return crypto.PubkeyToAddress(privKey.PublicKey).Hex(), nil
}

func (h *SovereignHD) DeriveEphemeral(index uint32) (string, string, error) {
	k := h.Root
	k, _ = k.NewChildKey(402)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(index)
	privHex := hex.EncodeToString(k.Key)
	privKey, _ := crypto.ToECDSA(k.Key)
	addr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	return addr, privHex, nil
}
