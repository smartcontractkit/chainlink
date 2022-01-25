package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

type gasEstimator struct {
	client        *ethclient.Client
	gasMultiplier float64
}

func (e *gasEstimator) estimate() (*big.Int, error) {
	suggested, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	return multiplyGasPrice(suggested, e.gasMultiplier), nil
}
