package mercury_common_test

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
	mercury_types "github.com/smartcontractkit/chainlink-common/pkg/types/mercury"
)

var ChainReader = staticMercuryChainReader{
	staticMercuryChainReaderConfig: staticMercuryChainReaderConfig{
		latestHeads: []mercury_types.Head{
			{
				Hash:   []byte{1},
				Number: 1,
			},
			{
				Hash:   []byte{2},
				Number: 2,
			},
		},
	},
}

type MercuryChainReaderEvaluator interface {
	mercury_types.ChainReader
	testtypes.Evaluator[mercury_types.ChainReader]
}

type staticMercuryChainReaderConfig struct {
	latestHeads []mercury_types.Head
}

type staticMercuryChainReader struct {
	staticMercuryChainReaderConfig
}

var _ mercury_types.ChainReader = staticMercuryChainReader{}

func (s staticMercuryChainReader) LatestHeads(ctx context.Context, n int) ([]mercury_types.Head, error) {
	return s.latestHeads, nil
}

func (s staticMercuryChainReader) Evaluate(ctx context.Context, other mercury_types.ChainReader) error {
	heads, err := other.LatestHeads(ctx, 2)
	if err != nil {
		return err
	}
	if !assert.ObjectsAreEqual(s.latestHeads, heads) {
		return fmt.Errorf("expected lastestHeads %v but got %v", s.latestHeads, heads)
	}

	return nil
}
