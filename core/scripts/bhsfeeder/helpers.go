package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func getConfig() (*config, error) {
	ethURL, set := os.LookupEnv("ETH_URL")
	if !set {
		return &config{}, errors.New("ETH_URL environment variable not set. Please set it to the appropriate value")
	}

	chainIDEnv, set := os.LookupEnv("ETH_CHAIN_ID")
	if !set {
		return &config{}, errors.New("ETH_CHAIN_ID environment variable not set. Please set it to the appropriate value")
	}

	chainID, err := strconv.ParseInt(chainIDEnv, 10, 64)
	if err != nil {
		return &config{}, errors.Wrap(err, "chain ID not an integer")
	}

	privateKey, set := os.LookupEnv("PRIVATE_KEY")
	if !set {
		return &config{}, errors.New("PRIVATE_KEY environment variable not set. Please set it to the appropriate value")
	}

	account, err := common.GetAccount(privateKey, big.NewInt(chainID))
	if err != nil {
		return &config{}, errors.Wrap(err, "account from private key")
	}

	return &config{
		ethURL:  ethURL,
		chainID: chainID,
		account: account,
	}, nil
}

func multiplyGasPrice(gasPrice *big.Int, multiplier float64) *big.Int {
	gpFloat := new(big.Float).SetInt(gasPrice)
	multiplied := gpFloat.Mul(gpFloat, new(big.Float).SetFloat64(multiplier))
	multipliedInt, _ := multiplied.Int(nil)
	return multipliedInt
}

func serializedBlockHeader(ethClient *ethclient.Client, blockNumber *big.Int, chainID *big.Int) ([]byte, error) {
	block, err := ethClient.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, err
	}

	encodedHeader, err := rlp.EncodeToBytes(block.Header())
	if err != nil {
		return nil, err
	}

	keccak, err := utils.Keccak256(encodedHeader)
	if err != nil {
		return nil, err
	}

	if gethcommon.BytesToHash(keccak) != block.Hash() {
		return nil, fmt.Errorf("blockhash does not match. Received invalid block header? %s vs %s", block.Hash().String(), gethcommon.BytesToHash(keccak).String())
	}

	return encodedHeader, nil
}
