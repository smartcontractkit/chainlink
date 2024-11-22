package actions

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/client"
)

// ReorgSuite is a test suite that generates reorgs on source/dest chains
type ReorgSuite struct {
	t         *testing.T
	Cfg       *ReorgConfig
	Logger    *zerolog.Logger
	SrcClient *client.RPCClient
	DstClient *client.RPCClient
}

// ReorgConfig is a configuration for reorg tests
type ReorgConfig struct {
	// SrcGethHTTPURL source chain Geth HTTP URL
	SrcGethHTTPURL string
	// DstGethHTTPURL dest chain Geth HTTP URL
	DstGethHTTPURL string
	// SrcFinalityDepth source chain finality depth
	SrcFinalityDepth int
	// DstGethHTTPURL dest chain finality depth
	DstFinalityDepth int
	// FinalityDelta blocks to rewind below or above finality
	FinalityDelta int
}

// Validate validates ReorgConfig params
func (rc *ReorgConfig) Validate() error {
	if rc.FinalityDelta >= rc.SrcFinalityDepth || rc.FinalityDelta >= rc.DstFinalityDepth {
		return fmt.Errorf(
			"finality delta can't be higher than source or dest chain finality, delta: %d, src: %d, dst: %d",
			rc.FinalityDelta, rc.SrcFinalityDepth, rc.DstFinalityDepth,
		)
	}
	return nil
}

// NewReorgSuite creates new reorg suite with source/dest RPC clients, works only with Geth
func NewReorgSuite(t *testing.T, lggr *zerolog.Logger, cfg *ReorgConfig) (*ReorgSuite, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &ReorgSuite{
		t:         t,
		Cfg:       cfg,
		Logger:    lggr,
		SrcClient: client.NewRPCClient(cfg.SrcGethHTTPURL, nil),
		DstClient: client.NewRPCClient(cfg.DstGethHTTPURL, nil),
	}, nil
}

// RunReorg rollbacks given chain, for N blocks back
func (r *ReorgSuite) RunReorg(client *client.RPCClient, blocksBack int, network string, startDelay time.Duration) {
	go func() {
		time.Sleep(startDelay)
		r.Logger.Info().
			Str("Network", network).
			Str("URL", client.URL).
			Int("BlocksBack", blocksBack).
			Msg(fmt.Sprintf("Rewinding blocks on %s chain", network))
		blockNumber, err := client.BlockNumber()
		assert.NoError(r.t, err, "error getting block number")
		r.Logger.Info().
			Int64("Number", blockNumber).
			Str("Network", network).
			Msg("Block number before rewinding")
		err = client.GethSetHead(blocksBack)
		assert.NoError(r.t, err, "error setting block head")
		blockNumber, err = client.BlockNumber()
		assert.NoError(r.t, err, "error getting block number")
		r.Logger.Info().
			Int64("Number", blockNumber).
			Str("Network", network).
			Msg("Block number after rewinding")
	}()
}
