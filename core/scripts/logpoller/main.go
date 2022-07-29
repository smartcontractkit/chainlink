package main

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"time"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

type nodeConfig struct {
}

func (n nodeConfig) NodeNoNewHeadsThreshold() time.Duration {
	return time.Minute
}

func (n nodeConfig) NodePollFailureThreshold() uint32 {
	return 10
}

func (n nodeConfig) NodePollInterval() time.Duration {
	return time.Minute
}

type cfg struct {
}

func (c cfg) LogSQL() bool {
	return true
}

func main() {
	lggr, done := logger.NewLogger()
	defer done()
	err := os.Setenv("DATABASE_URL", "TODO")
	panicErr(err)
	db, err := pg.OpenUnlockedDB(config.NewGeneralConfig(lggr), lggr)
	panicErr(err)
	defer db.Close()
	chainID := big.NewInt(137) // E.g. polygon mainnet.
	_, err = db.Exec(fmt.Sprintf("INSERT INTO evm_chains (id, created_at, updated_at) VALUES (%d, NOW(), NOW()) ON CONFLICT DO NOTHING", chainID.Int64()))
	panicErr(err)
	// Can try different RPCs.
	ws, _ := url.Parse("TODO")
	client, err := evmclient.NewClientWithNodes(lggr, []evmclient.Node{evmclient.NewNode(nodeConfig{}, lggr, *ws, nil, "poly", 1, chainID)}, nil, chainID)
	panicErr(err)
	err = client.Dial(context.Background())
	panicErr(err)
	s := time.Now()
	b, err := client.BlockByNumber(context.Background(), nil)
	panicErr(err)
	fmt.Println(b.Number().Int64(), time.Since(s))
	lp := logpoller.NewLogPoller(logpoller.NewORM(chainID, db, lggr, cfg{}),
		client, lggr, 1*time.Second, 1000, 3)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// TODO: can add filters to test log inserts.
	err = lp.Start(ctx)
	<-ctx.Done()
	// Inspect DB to check the log poller can keep up with the chain, has the logs expected etc.
}
