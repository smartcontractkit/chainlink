package handler

import (
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	link "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

// baseHandler is the common handler with a common logic
type baseHandler struct {
	cfg *config.Config

	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	linkToken  *link.LinkToken
	fromAddr   common.Address
}

// newBaseHandler is the constructor of baseHandler
func newBaseHandler(cfg *config.Config) *baseHandler {
	// Created a client by the given node address
	nodeClient, err := ethclient.Dial(cfg.NodeURL)
	if err != nil {
		log.Fatal("failed to deal with ETH node", err)
	}

	// Parse private key
	d := new(big.Int).SetBytes(common.FromHex(cfg.PrivateKey))
	pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pkX,
			Y:     pkY,
		},
		D: d,
	}

	// Init from address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddr := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Create link token wrapper
	linkToken, err := link.NewLinkToken(common.HexToAddress(cfg.LinkTokenAddr), nodeClient)
	if err != nil {
		log.Fatal(err)
	}

	return &baseHandler{
		cfg:        cfg,
		client:     nodeClient,
		privateKey: privateKey,
		linkToken:  linkToken,
		fromAddr:   fromAddr,
	}
}
