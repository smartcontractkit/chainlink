package handler

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	link "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

// baseHandler is the common handler with a common logic
type baseHandler struct {
	cfg *config.Config

	client        *ethclient.Client
	privateKey    *ecdsa.PrivateKey
	linkToken     *link.LinkToken
	fromAddr      common.Address
	approveAmount *big.Int
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

	approveAmount := big.NewInt(0)
	approveAmount.SetString(cfg.ApproveAmount, 10)

	return &baseHandler{
		cfg:           cfg,
		client:        nodeClient,
		privateKey:    privateKey,
		linkToken:     linkToken,
		fromAddr:      fromAddr,
		approveAmount: approveAmount,
	}
}

func (h *baseHandler) buildTxOpts(ctx context.Context) *bind.TransactOpts {
	nonce, err := h.client.PendingNonceAt(ctx, h.fromAddr)
	if err != nil {
		log.Fatal("PendingNonceAt failed: ", err)
	}

	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("SuggestGasPrice failed: ", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(h.privateKey, big.NewInt(h.cfg.ChainID))
	if err != nil {
		log.Fatal("NewKeyedTransactorWithChainID failed: ", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = h.cfg.GasLimit // in units
	auth.GasPrice = gasPrice

	return auth
}

func (h *baseHandler) waitDeployment(ctx context.Context, tx *types.Transaction) {
	if _, err := bind.WaitDeployed(ctx, h.client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}
}

func (h *baseHandler) waitTx(ctx context.Context, tx *types.Transaction) {
	if _, err := bind.WaitMined(ctx, h.client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}
}
