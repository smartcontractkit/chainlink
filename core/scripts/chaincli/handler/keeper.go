package handler

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	keeper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/keeper_registry_wrapper"
	link "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	upkeep "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/upkeep_perform_counter_restrictive_wrapper"
	"github.com/smartcontractkit/chainlink/core/scripts/chaincli/config"
)

// Keeper is the keepers commands handler
type Keeper struct {
	cfg *config.Config

	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	linkToken  *link.LinkToken
	fromAddr   common.Address

	approveAmount  *big.Int
	addFundsAmount *big.Int
}

// NewKeeper is the constructor of Keeper
func NewKeeper(cfg *config.Config) *Keeper {
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

	addFundsAmount := big.NewInt(0)
	addFundsAmount.SetString(cfg.AddFundsAmount, 10)

	return &Keeper{
		cfg:            cfg,
		client:         nodeClient,
		privateKey:     privateKey,
		linkToken:      linkToken,
		fromAddr:       fromAddr,
		approveAmount:  approveAmount,
		addFundsAmount: addFundsAmount,
	}
}

// DeployKeepers contains a logic to deploy keepers
func (h *Keeper) DeployKeepers(ctx context.Context) {
	// Deploy keeper registry
	log.Println("Deploying keeper registry...")
	registryAddr, deployKeeperRegistryTx, registryInstance, err := keeper.DeployKeeperRegistry(h.buildTxOpts(ctx), h.client,
		common.HexToAddress(h.cfg.LinkTokenAddr),
		common.HexToAddress(h.cfg.LinkETHFeedAddr),
		common.HexToAddress(h.cfg.FastGasFeedAddr),
		h.cfg.PaymentPremiumPBB,
		h.cfg.FlatFeeMicroLink,
		big.NewInt(h.cfg.BlockCountPerTurn),
		h.cfg.CheckGasLimit,
		big.NewInt(h.cfg.StalenessSeconds),
		h.cfg.GasCeilingMultiplier,
		big.NewInt(h.cfg.FallbackGasPrice),
		big.NewInt(h.cfg.FallbackLinkPrice),
	)
	if err != nil {
		log.Fatal("DeployAbi failed: ", err)
	}
	log.Println("Waiting for keeper registry contract deployment confirmation...", deployKeeperRegistryTx.Hash().Hex())
	h.waitDeployment(ctx, deployKeeperRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry deployed - ", deployKeeperRegistryTx.Hash().Hex())

	// Approve keeper registry
	approveRegistryTx, err := h.linkToken.Approve(h.buildTxOpts(ctx), registryAddr, h.approveAmount)
	if err != nil {
		log.Fatal(registryAddr.Hex(), ": Approve failed - ", err)
	}
	h.waitTx(ctx, approveRegistryTx)
	log.Println(registryAddr.Hex(), ": KeeperRegistry approved - ", approveRegistryTx.Hash().Hex())

	// Deploy Upkeeps
	h.deployUpkeeps(ctx, registryInstance)

	// Set Keepers
	log.Println("Set keepers...")
	keepers, owners := h.keepers()
	setKeepersTx, err := registryInstance.SetKeepers(h.buildTxOpts(ctx), keepers, owners)
	if err != nil {
		log.Fatal("SetKeepers failed: ", err)
	}
	h.waitTx(ctx, setKeepersTx)
	log.Println("Keepers registered:", setKeepersTx.Hash().Hex())
}

func (h *Keeper) keepers() ([]common.Address, []common.Address) {
	var addrs []common.Address
	var fromAddrs []common.Address
	for _, addr := range h.cfg.Keepers {
		addrs = append(addrs, common.HexToAddress(addr))
		fromAddrs = append(fromAddrs, h.fromAddr)
	}
	return addrs, fromAddrs
}

func (h *Keeper) deployUpkeeps(ctx context.Context, registryInstance *keeper.KeeperRegistry) {
	fmt.Println()
	log.Println("Deploying upkeeps...")
	for i := int64(0); i < h.cfg.UpkeepCount; i++ {
		fmt.Println()
		// Deploy
		upkeepAddr, deployUpkeepTx, _, err := upkeep.DeployUpkeepPerformCounterRestrictive(h.buildTxOpts(ctx), h.client,
			big.NewInt(h.cfg.UpkeepTestRange), big.NewInt(h.cfg.UpkeepAverageEligibilityCadence),
		)
		if err != nil {
			log.Fatal(i, ": DeployAbi failed - ", err)
		}
		h.waitDeployment(ctx, deployUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep deployed - ", deployUpkeepTx.Hash().Hex())

		// Approve
		approveUpkeepTx, err := h.linkToken.Approve(h.buildTxOpts(ctx), upkeepAddr, h.approveAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": Approve failed - ", err)
		}
		h.waitTx(ctx, approveUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep approved - ", approveUpkeepTx.Hash().Hex())

		// Register
		registerUpkeepTx, err := registryInstance.RegisterUpkeep(h.buildTxOpts(ctx),
			upkeepAddr, h.cfg.UpkeepGasLimit, h.fromAddr, []byte(h.cfg.UpkeepCheckData),
		)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": RegisterUpkeep failed - ", err)
		}
		h.waitTx(ctx, registerUpkeepTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep registered - ", registerUpkeepTx.Hash().Hex())

		// Fund
		addFundsTx, err := registryInstance.AddFunds(h.buildTxOpts(ctx), big.NewInt(int64(i)), h.addFundsAmount)
		if err != nil {
			log.Fatal(i, upkeepAddr.Hex(), ": AddFunds failed - ", err)
		}
		h.waitTx(ctx, addFundsTx)
		log.Println(i, upkeepAddr.Hex(), ": Upkeep funded - ", addFundsTx.Hash().Hex())
	}
	fmt.Println()
}

func (h *Keeper) buildTxOpts(ctx context.Context) *bind.TransactOpts {
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

func (h *Keeper) waitDeployment(ctx context.Context, tx *types.Transaction) {
	if _, err := bind.WaitDeployed(ctx, h.client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}
}

func (h *Keeper) waitTx(ctx context.Context, tx *types.Transaction) {
	if _, err := bind.WaitMined(ctx, h.client, tx); err != nil {
		log.Fatal("WaitDeployed failed: ", err)
	}
}
