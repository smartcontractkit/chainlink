package main

import (
	"github.com/ethereum/go-ethereum/common"
)

type remote struct {
	host     string
	login    string
	password string
}

func (r remote) IsTerminal() bool {
	return false
}

func (r remote) PasswordPrompt(p string) string {
	return r.password
}

func (r remote) Prompt(p string) string {
	return r.login
}

type ocr2Bundle struct {
	ID                string `json:"id"`
	ChainType         string `json:"chainType"`
	OnchainPublicKey  string `json:"onchainPublicKey"`
	OffchainPublicKey string `json:"offchainPublicKey"`
	ConfigPublicKey   string `json:"configPublicKey"`
}

type Node struct {
	Host                string   `json:"host"`
	ETHKeys             []string `json:"eth_key"`
	P2PPeerIDS          []string `json:"p2p_ids"`
	OCR2KeyIDs          []string `json:"ocr2_ids"`
	OCR2ConfigPubKeys   []string `json:"ocr2_pub_keys"`
	OCR2OffchainPubKeys []string `json:"ocr2_offchain_pub_keys"`
	OCR2OnchainPubKeys  []string `json:"ocr2_onchain_pub_keys"`
}

type config struct {
	ChainID            int64          `yaml:"chain-id"`
	P2PPort            int64          `yaml:"p2p-port"`
	DONContractAddress common.Address `yaml:"don-contract-address"`
}
