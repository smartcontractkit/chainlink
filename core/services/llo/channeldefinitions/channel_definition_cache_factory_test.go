package channeldefinitions

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	lloconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo/config"
)

func Test_ChannelDefinitionCacheFactory(t *testing.T) {
	lggr := logger.TestLogger(t)
	cdcFactory := NewChannelDefinitionCacheFactory(lggr, nil, nil, nil)

	t.Run("NewCache", func(t *testing.T) {
		t.Run("when ChannelDefinitions is present, returns static cache", func(t *testing.T) {
			_, err := cdcFactory.NewCache(lloconfig.PluginConfig{ChannelDefinitions: "..."})
			require.EqualError(t, err, "failed to unmarshal static channel definitions: invalid character '.' looking for beginning of value")

			cdc, err := cdcFactory.NewCache(lloconfig.PluginConfig{ChannelDefinitions: "{}"})
			require.NoError(t, err)
			require.IsType(t, &staticCDC{}, cdc)
		})
		t.Run("when ChannelDefinitions is not present, returns dynamic cache", func(t *testing.T) {
			cdc, err := cdcFactory.NewCache(lloconfig.PluginConfig{
				ChannelDefinitionsContractAddress: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				DonID:                             1,
			})
			require.NoError(t, err)
			require.IsType(t, &channelDefinitionCache{}, cdc)

			// creates another one if you try to do it again with the same addr/donID
			cdc, err = cdcFactory.NewCache(lloconfig.PluginConfig{
				ChannelDefinitionsContractAddress: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				DonID:                             1,
			})
			require.NoError(t, err)
			require.IsType(t, &channelDefinitionCache{}, cdc)

			// is fine if you do it again with different addr
			cdc, err = cdcFactory.NewCache(lloconfig.PluginConfig{
				ChannelDefinitionsContractAddress: common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
				DonID:                             1,
			})
			require.NoError(t, err)
			require.IsType(t, &channelDefinitionCache{}, cdc)

			// is fine if you do it again with different don ID
			cdc, err = cdcFactory.NewCache(lloconfig.PluginConfig{
				ChannelDefinitionsContractAddress: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				DonID:                             2,
			})
			require.NoError(t, err)
			require.IsType(t, &channelDefinitionCache{}, cdc)
		})
	})
}
