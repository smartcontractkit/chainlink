package multienv

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	chainsel "github.com/smartcontractkit/chain-selectors"
)

// Environment variables used to configure
// the environment for the rebalancer
const (
	// OwnerKey is the private key used to deploy contracts and send funds to the rebalancer nodes
	OwnerKey = "OWNER_KEY"
	// RPCPrefix is the prefix for the environment variable that contains the RPC URL for a chain
	RPCPrefix = "RPC_"
	// WebSocketPerfix is the prefix for the environment variable that contains the WebSocket URL for a chain
	WebSocketPerfix = "WS_"
)

type Env struct {
	Transactors map[uint64]*bind.TransactOpts
	Clients     map[uint64]*ethclient.Client
	JRPCs       map[uint64]*rpc.Client
	HTTPURLs    map[uint64]string
	WSURLs      map[uint64]string
}

// nolint
func New(websocket bool, overrideNonce bool) Env {
	env := Env{
		Transactors: make(map[uint64]*bind.TransactOpts),
		Clients:     make(map[uint64]*ethclient.Client),
		JRPCs:       make(map[uint64]*rpc.Client),
		HTTPURLs:    make(map[uint64]string),
		WSURLs:      make(map[uint64]string),
	}
	for _, chainID := range []uint64{
		chainsel.ETHEREUM_TESTNET_SEPOLIA.EvmChainID,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.EvmChainID,
		chainsel.ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1.EvmChainID,
	} {
		client, rpcClient, err := GetClient(chainID, websocket)
		if err != nil {
			log.Println("error getting client for chain, assuming not specified, chain id:", chainID, ", err:", err)
		} else {
			env.Clients[chainID] = client
			env.JRPCs[chainID] = rpcClient
			env.Transactors[chainID] = GetTransactor(big.NewInt(int64(chainID)))
			env.HTTPURLs[chainID], _ = GetRPC(chainID)
			if websocket {
				env.WSURLs[chainID], _ = GetWS(chainID)
			}
		}
	}
	return env
}

func GetRPC(chainID uint64) (string, error) {
	envVariable := RPCPrefix + strconv.FormatUint(chainID, 10)
	rpc := os.Getenv(envVariable)
	if rpc != "" {
		return rpc, nil
	}
	return "", fmt.Errorf("RPC not found. Please set the environment variable for chain %d e.g. RPC_420=https://rpc.420.com", chainID)
}

func GetWS(chainID uint64) (string, error) {
	envVariable := WebSocketPerfix + strconv.FormatUint(chainID, 10)
	ws := os.Getenv(envVariable)
	if ws != "" {
		return ws, nil
	}
	return "", fmt.Errorf("WS not found. Please set the environment variable for chain %d e.g. WS_420=wss://ws.420.com", chainID)
}

func GetTransactor(chainID *big.Int) *bind.TransactOpts {
	ownerKey := os.Getenv(OwnerKey)
	if ownerKey != "" {
		b, err := hex.DecodeString(ownerKey)
		if err != nil {
			panic(err)
		}
		d := new(big.Int).SetBytes(b)

		pkX, pkY := crypto.S256().ScalarBaseMult(d.Bytes())
		privateKey := ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: crypto.S256(),
				X:     pkX,
				Y:     pkY,
			},
			D: d,
		}
		owner, err := bind.NewKeyedTransactorWithChainID(&privateKey, chainID)
		if err != nil {
			panic(err)
		}
		return owner
	}
	panic("OWNER_KEY not found. Please set the environment variable OWNER_KEY with the private key of the owner")
}

func GetClient(chainID uint64, websocket bool) (*ethclient.Client, *rpc.Client, error) {
	rpcURL, err := GetRPC(chainID)
	if err != nil {
		return nil, nil, err
	}
	_, err = GetWS(chainID)
	if err != nil {
		return nil, nil, err
	}
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		// failing to dial is blocking, so panic
		panic(err)
	}
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		// failing to dial is blocking, so panic
		panic(err)
	}
	return client, rpcClient, nil
}
