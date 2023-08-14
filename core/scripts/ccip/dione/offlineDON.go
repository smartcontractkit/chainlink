package dione

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/core/scripts/common"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type OfflineDON struct {
	Config NodesConfig
	env    Environment
	lggr   logger.Logger
}

func NewOfflineDON(env Environment, lggr logger.Logger) OfflineDON {
	config := MustReadNodeConfig(env)

	return OfflineDON{
		Config: config,
		env:    env,
		lggr:   lggr,
	}
}

func (don *OfflineDON) GenerateOracleIdentities(chain uint64) []confighelper2.OracleIdentityExtra {
	var oracles []confighelper2.OracleIdentityExtra

	for _, node := range don.Config.Nodes {
		evmKeys := GetOCRkeysForChainType(node.OCRKeys, "evm")

		oracles = append(oracles,
			confighelper2.OracleIdentityExtra{
				OracleIdentity: confighelper2.OracleIdentity{
					TransmitAccount:   ocr2types.Account(node.EthKeys[strconv.FormatUint(chain, 10)]),
					OnchainPublicKey:  gethcommon.HexToAddress(strings.TrimPrefix(evmKeys.Attributes.OnChainPublicKey, "ocr2on_evm_")).Bytes(),
					OffchainPublicKey: common.ToOffchainPublicKey("0x" + strings.TrimPrefix(evmKeys.Attributes.OffChainPublicKey, "ocr2off_evm_")),
					PeerID:            node.PeerID,
				},
				ConfigEncryptionPublicKey: common.StringTo32Bytes("0x" + strings.TrimPrefix(evmKeys.Attributes.ConfigPublicKey, "ocr2cfg_evm_")),
			})
	}
	return oracles
}

func (don *OfflineDON) GetSendingKeys(chain uint64) (keys []gethcommon.Address) {
	for _, node := range don.Config.Nodes {
		keys = append(keys, gethcommon.HexToAddress(node.EthKeys[strconv.FormatUint(chain, 10)]))
	}
	return
}

func (don *OfflineDON) FundNodeKeys(chainConfig *rhea.EvmDeploymentConfig, ownerPrivKey string, amount *big.Int, fundingThreshold *big.Int) {
	currentNonce, err := chainConfig.Client.PendingNonceAt(context.Background(), chainConfig.Owner.From)
	helpers.PanicErr(err)
	var gasTipCap *big.Int
	if chainConfig.ChainConfig.GasSettings.EIP1559 {
		gasTipCap, err = chainConfig.Client.SuggestGasTipCap(context.Background())
		helpers.PanicErr(err)
	}
	gasPrice, err := chainConfig.Client.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)

	ownerKey, err := crypto.HexToECDSA(ownerPrivKey)
	helpers.PanicErr(err)

	don.lggr.Infof("Chain id %d", chainConfig.ChainConfig.EvmChainId)

	nonceIncrement := 0
	for i, node := range don.Config.Nodes {
		eoa := gethcommon.HexToAddress(node.EthKeys[strconv.FormatUint(chainConfig.ChainConfig.EvmChainId, 10)])
		if eoa == gethcommon.HexToAddress("0x") {
			don.lggr.Warnf("Node %2d has no sending key configured. Skipping funding", i)
			continue
		}
		balanceAt, err := chainConfig.Client.BalanceAt(context.Background(), eoa, nil)
		helpers.PanicErr(err)

		if balanceAt.Cmp(fundingThreshold) == -1 {
			don.lggr.Infof("❌ Node %2d has a balance of %s eth, which is lower than the set minimum. Funding...", i, EthBalanceToString(balanceAt))

			if chainConfig.ChainConfig.GasSettings.EIP1559 {
				sendEthEIP1559(eoa, *chainConfig, currentNonce+uint64(nonceIncrement), gasTipCap, ownerKey, amount)
			} else {
				sendEth(eoa, *chainConfig, currentNonce+uint64(nonceIncrement), gasPrice, ownerKey, amount)
			}
			nonceIncrement++
			don.lggr.Infof("Sent %s eth to %s", EthBalanceToString(amount), eoa.Hex())
		} else {
			don.lggr.Infof("✅ Node %2d has a balance of %s eth ", i, EthBalanceToString(balanceAt))
		}
	}
}

type NodeWallet struct {
	ChainID uint64
	Address string
}

func (don *OfflineDON) GetAllNodesWallets(chainId uint64) []NodeWallet {
	fmt.Printf("ChainId %d\n", chainId)
	var wallets []NodeWallet
	for _, node := range don.Config.Nodes {
		eoa := gethcommon.HexToAddress(node.EthKeys[strconv.FormatUint(chainId, 10)]).String()
		if eoa == (gethcommon.Address{}).String() {
			don.lggr.Warnf("Node %s has no sending key configured. Skipping.\n", eoa)
			continue
		}
		fmt.Printf("%s\n", eoa)
		wallets = append(wallets, NodeWallet{chainId, eoa})
	}
	return wallets
}

func EthBalanceToString(balance *big.Int) string {
	return new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18)).String()
}

func (don *OfflineDON) WriteToFile() error {
	path := getFileLocation(don.env, NODES_FOLDER)
	file, err := json.MarshalIndent(don.Config, "", "  ")
	if err != nil {
		return err
	}
	return WriteJSON(path, file)
}

func (don *OfflineDON) PrintConfig() {
	file, err := json.MarshalIndent(don.Config, "", "  ")
	common.PanicErr(err)

	don.lggr.Infof(string(file))
}

func sendEth(to gethcommon.Address, chainConfig rhea.EvmDeploymentConfig, nonce uint64, gasPrice *big.Int, ownerKey *ecdsa.PrivateKey, amount *big.Int) {
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      21_000,
			To:       &to,
			Value:    amount,
			Data:     []byte{},
		},
	)

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(0).SetUint64(chainConfig.ChainConfig.EvmChainId)), ownerKey)
	helpers.PanicErr(err)
	err = chainConfig.Client.SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
}

func sendEthEIP1559(to gethcommon.Address, chainConfig rhea.EvmDeploymentConfig, nonce uint64, gasTipCap *big.Int, ownerKey *ecdsa.PrivateKey, amount *big.Int) {
	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:    big.NewInt(0).SetUint64(chainConfig.ChainConfig.EvmChainId),
			Nonce:      nonce,
			GasTipCap:  gasTipCap,
			GasFeeCap:  big.NewInt(2e9),
			Gas:        uint64(21_000),
			To:         &to,
			Value:      amount,
			Data:       []byte{},
			AccessList: types.AccessList{},
		},
	)

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(big.NewInt(0).SetUint64(chainConfig.ChainConfig.EvmChainId)), ownerKey)
	helpers.PanicErr(err)
	err = chainConfig.Client.SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
}
