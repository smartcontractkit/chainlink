// Package deploy provides utilities for deploying LLO contracts to a blockchain
// for use in E2E testing of the streams-trigger capability.
//
// This package uses the existing gethwrappers from chainlink-evm, avoiding
// the need to import external deployment tools like gauntlet-plus-plus.
//
// References:
// - Contract source: https://github.com/smartcontractkit/chainlink-evm/tree/develop/contracts/src/v0.8/llo-feeds/v0.5.1/configuration
// - Deployment sequences: https://github.com/smartcontractkit/gauntlet-plus-plus/tree/main/packages-ethereum/operations-mercury/src/sequences/v0_5_0
package deploy

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/channel_config_store"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
)

// LLOContracts holds the deployed contract addresses and instances
type LLOContracts struct {
	ConfiguratorAddress    common.Address
	Configurator           *configurator.Configurator
	ChannelConfigStore     *channel_config_store.ChannelConfigStore
	ChannelConfigStoreAddr common.Address
}

// DeployConfig holds configuration for deploying LLO contracts
type DeployConfig struct {
	// RPC URL for the target chain
	RPCURL string
	// Private key for the deployer account (hex string without 0x prefix)
	PrivateKey string
	// Chain ID
	ChainID *big.Int
}

// OCRConfig holds the OCR configuration to set on the Configurator
type OCRConfig struct {
	DonID                 uint32
	Signers               [][]byte   // OCR signing key public keys
	Transmitters          [][32]byte // CSA public keys
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

// ChannelDefinition represents a single channel configuration
type ChannelDefinition struct {
	ReportFormat llotypes.ReportFormat
	Streams      []llotypes.Stream
	Opts         json.RawMessage
}

// Deploy deploys the Configurator and ChannelConfigStore contracts
func Deploy(ctx context.Context, cfg DeployConfig) (*LLOContracts, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, cfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	// Deploy Configurator
	configuratorAddr, tx, configuratorContract, err := configurator.DeployConfigurator(auth, client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy Configurator: %w", err)
	}
	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for Configurator deployment: %w", err)
	}
	fmt.Printf("Configurator deployed at: %s\n", configuratorAddr.Hex())

	// Deploy ChannelConfigStore
	channelConfigStoreAddr, tx, channelConfigStore, err := channel_config_store.DeployChannelConfigStore(auth, client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ChannelConfigStore: %w", err)
	}
	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for ChannelConfigStore deployment: %w", err)
	}
	fmt.Printf("ChannelConfigStore deployed at: %s\n", channelConfigStoreAddr.Hex())

	return &LLOContracts{
		ConfiguratorAddress:    configuratorAddr,
		Configurator:           configuratorContract,
		ChannelConfigStore:     channelConfigStore,
		ChannelConfigStoreAddr: channelConfigStoreAddr,
	}, nil
}

// SetProductionConfig sets the OCR production configuration on the Configurator
func (c *LLOContracts) SetProductionConfig(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	cfg OCRConfig,
) error {
	// Convert DON ID to bytes32 (configId)
	configID := DonIDToBytes32(cfg.DonID)

	// Set config on configurator
	// Signature: SetProductionConfig(configId [32]byte, signers [][]byte, offchainTransmitters [][32]byte, f uint8, onchainConfig []byte, offchainConfigVersion uint64, offchainConfig []byte)
	tx, err := c.Configurator.SetProductionConfig(
		auth,
		configID,
		cfg.Signers,
		cfg.Transmitters,
		cfg.F,
		cfg.OnchainConfig,
		cfg.OffchainConfigVersion,
		cfg.OffchainConfig,
	)
	if err != nil {
		return fmt.Errorf("failed to set production config: %w", err)
	}

	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for SetProductionConfig tx: %w", err)
	}

	fmt.Printf("Production config set for DON ID: %d\n", cfg.DonID)
	return nil
}

// SetChannelDefinitions sets channel definitions on the ChannelConfigStore
func (c *LLOContracts) SetChannelDefinitions(
	ctx context.Context,
	client *ethclient.Client,
	auth *bind.TransactOpts,
	donID uint32,
	url string, // URL where channel definitions are hosted
	sha [32]byte, // SHA256 hash of the channel definitions
) error {
	// Signature: SetChannelDefinitions(donId uint32, url string, sha [32]byte)
	tx, err := c.ChannelConfigStore.SetChannelDefinitions(
		auth,
		donID,
		url,
		sha,
	)
	if err != nil {
		return fmt.Errorf("failed to set channel definitions: %w", err)
	}

	_, err = bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for SetChannelDefinitions tx: %w", err)
	}

	fmt.Printf("Channel definitions set for DON ID: %d, URL: %s\n", donID, url)
	return nil
}

// DonIDToBytes32 converts a DON ID to a bytes32 value (left-padded)
func DonIDToBytes32(donID uint32) [32]byte {
	var result [32]byte
	big.NewInt(int64(donID)).FillBytes(result[:])
	return result
}

// GetTransactorFromPrivateKey creates a transactor from a private key
func GetTransactorFromPrivateKey(privateKeyHex string, chainID *big.Int) (*bind.TransactOpts, *ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	return auth, privateKey, nil
}

// GetSignerAddress returns the address for a private key
func GetSignerAddress(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}

// LoadExistingContracts loads existing LLO contract instances from their addresses
func LoadExistingContracts(client *ethclient.Client, configuratorAddr, channelConfigStoreAddr common.Address) (*LLOContracts, error) {
	configuratorContract, err := configurator.NewConfigurator(configuratorAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load Configurator contract: %w", err)
	}

	channelConfigStoreContract, err := channel_config_store.NewChannelConfigStore(channelConfigStoreAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to load ChannelConfigStore contract: %w", err)
	}

	return &LLOContracts{
		ConfiguratorAddress:    configuratorAddr,
		Configurator:           configuratorContract,
		ChannelConfigStore:     channelConfigStoreContract,
		ChannelConfigStoreAddr: channelConfigStoreAddr,
	}, nil
}
