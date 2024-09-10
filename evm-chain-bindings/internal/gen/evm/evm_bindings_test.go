package evm

import (
	"context"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"testing"
)

func test(t testing.T) {

	chainReaderConfig := types.ChainReaderConfig{}

	//TODO: Fix nil to actually connect to EVM.
	chainReaderService, err := evm.NewChainReaderService(context.Background(), logger.TestLogger(&t), nil, nil, nil, chainReaderConfig)

	if err != nil {
		t.Fail()
	}

	//TODO implement usage
	chainReaderService.Name()
}
