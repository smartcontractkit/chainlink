package discoverer

import (
	"testing"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestNewFactory(t *testing.T) {
	type args struct {
		lggr logger.Logger
		opts []Opt
	}
	tests := []struct {
		name string
		args args
		want *factory
	}{
		{
			"no opts",
			args{
				lggr: logger.TestLogger(t),
			},
			&factory{
				evmDeps: make(map[models.NetworkSelector]evmDep),
				lggr:    logger.TestLogger(t),
			},
		},
		{
			"with opts",
			args{
				lggr: logger.TestLogger(t),
				opts: []Opt{
					WithEvmDep(models.NetworkSelector(1), nil),
					WithEvmDep(models.NetworkSelector(2), nil),
					WithEvmDep(models.NetworkSelector(3), nil),
				},
			},
			&factory{
				evmDeps: map[models.NetworkSelector]evmDep{
					models.NetworkSelector(1): {},
					models.NetworkSelector(2): {},
					models.NetworkSelector(3): {},
				},
				lggr: logger.TestLogger(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFactory(tt.args.lggr, tt.args.opts...)
			require.Equal(t, tt.want.evmDeps, got.evmDeps)
		})
	}
}

func Test_factory_NewDiscoverer(t *testing.T) {
	t.Run("uncached", func(t *testing.T) {
		network := models.NetworkSelector(chainsel.TEST_90000001.Selector)
		f := NewFactory(logger.TestLogger(t), WithEvmDep(network, nil))
		got, err := f.NewDiscoverer(network, models.Address{})
		require.NoError(t, err)
		require.NotNil(t, got)
	})

	t.Run("cached", func(t *testing.T) {
		network := models.NetworkSelector(chainsel.TEST_90000001.Selector)
		f := NewFactory(logger.TestLogger(t), WithEvmDep(network, nil))

		want, err := f.NewDiscoverer(network, models.Address{})
		require.NoError(t, err)

		got, err := f.NewDiscoverer(network, models.Address{})
		require.NoError(t, err)

		require.Equal(t, want, got)
		_, ok := f.cachedDiscoverers.Load(f.cacheKey(network, models.Address{}))
		require.True(t, ok)
	})

	t.Run("network doesn't exist", func(t *testing.T) {
		network := models.NetworkSelector(chainsel.TEST_90000001.Selector)
		f := NewFactory(logger.TestLogger(t), WithEvmDep(network, nil))

		otherNetwork := models.NetworkSelector(chainsel.TEST_90000002.Selector)
		_, err := f.NewDiscoverer(otherNetwork, models.Address{})
		require.Error(t, err)
	})
}
