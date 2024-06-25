package reader

import (
	"context"
	"fmt"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	libocrtypes "github.com/smartcontractkit/libocr/ragep2p/types"

	"github.com/smartcontractkit/ccipocr3/internal/mocks"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/libocr/commontypes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	chainA       = cciptypes.ChainSelector(1)
	chainB       = cciptypes.ChainSelector(2)
	chainC       = cciptypes.ChainSelector(3)
	oracleAId    = commontypes.OracleID(1)
	p2pOracleAId = libocrtypes.PeerID{byte(oracleAId)}
	oracleBId    = commontypes.OracleID(2)
	p2pOracleBId = libocrtypes.PeerID{byte(oracleBId)}
	oracleCId    = commontypes.OracleID(3)
	p2pOracleCId = libocrtypes.PeerID{byte(oracleCId)}
)

func TestHomeChainConfigPoller_HealthReport(t *testing.T) {
	homeChainReader := mocks.NewContractReaderMock()
	homeChainReader.On(
		"GetLatestValue",
		mock.Anything,
		"CCIPCapabilityConfiguration",
		"getAllChainConfigs",
		mock.Anything,
		mock.Anything).Return(fmt.Errorf("error"))

	configPoller := NewHomeChainConfigPoller(
		homeChainReader,
		logger.Test(t),
		50*time.Millisecond,
	)
	_ = configPoller.Start(context.Background())

	// Initially it's healthy
	healthy := configPoller.HealthReport()
	assert.Equal(t, map[string]error{configPoller.Name(): error(nil)}, healthy)

	// After one second it will try polling 10 times and fail
	time.Sleep(1 * time.Second)

	errors := configPoller.HealthReport()

	err := configPoller.Close()
	time.Sleep(100 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(errors))
	assert.Errorf(t, errors[configPoller.Name()], "polling failed %d times in a row", MaxFailedPolls)
}

func Test_PollingWorking(t *testing.T) {
	onChainConfigs := []ChainConfigInfo{
		{
			ChainSelector: chainA,
			ChainConfig: HomeChainConfigMapper{
				FChain: 1,
				Readers: []libocrtypes.PeerID{
					p2pOracleAId,
					p2pOracleBId,
					p2pOracleCId,
				},
				Config: []byte{0},
			},
		},
		{
			ChainSelector: chainB,
			ChainConfig: HomeChainConfigMapper{
				FChain: 2,
				Readers: []libocrtypes.PeerID{
					p2pOracleAId,
					p2pOracleBId,
				},
				Config: []byte{0},
			},
		},
		{
			ChainSelector: chainC,
			ChainConfig: HomeChainConfigMapper{
				FChain: 3,
				Readers: []libocrtypes.PeerID{
					p2pOracleCId,
				},
				Config: []byte{0},
			},
		},
	}
	homeChainConfig := map[cciptypes.ChainSelector]ChainConfig{
		chainA: {
			FChain:         1,
			SupportedNodes: mapset.NewSet(p2pOracleAId, p2pOracleBId, p2pOracleCId),
		},
		chainB: {
			FChain:         2,
			SupportedNodes: mapset.NewSet(p2pOracleAId, p2pOracleBId),
		},
		chainC: {
			FChain:         3,
			SupportedNodes: mapset.NewSet(p2pOracleCId),
		},
	}

	homeChainReader := mocks.NewContractReaderMock()
	homeChainReader.On(
		"GetLatestValue", mock.Anything, "CCIPCapabilityConfiguration", "getAllChainConfigs", mock.Anything, mock.Anything).Run(
		func(args mock.Arguments) {
			arg := args.Get(4).(*[]ChainConfigInfo)
			*arg = onChainConfigs
		}).Return(nil)

	configPoller := NewHomeChainConfigPoller(
		homeChainReader,
		logger.Test(t),
		400*time.Millisecond,
	)

	ctx := context.Background()
	err := configPoller.Start(ctx)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	err = configPoller.Close()
	assert.NoError(t, err)

	// called 3 times, once when it's started, and 2 times when it's polling
	homeChainReader.AssertNumberOfCalls(t, "GetLatestValue", 3)

	configs, err := configPoller.GetAllChainConfigs()
	assert.NoError(t, err)
	assert.Equal(t, homeChainConfig, configs)
}
