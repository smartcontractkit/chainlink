package headtracker_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	htmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/mocks"
)

func configureInMemorySaver(t *testing.T) *headtracker.EvmInMemoryHeadSaver {
	htCfg := htmocks.NewConfig(t)
	htCfg.On("EvmHeadTrackerHistoryDepth").Return(uint32(6))
	htCfg.On("EvmFinalityDepth").Return(uint32(1))
	return headtracker.NewEvmInMemoryHeadSaver()
}

func TestInMemoryHeadSaver_Save(t *testing.T)
