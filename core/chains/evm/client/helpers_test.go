package client

import (
	"fmt"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/logger"
)

type TestNodeConfig struct {
	NoNewHeadsThreshold  time.Duration
	PollFailureThreshold uint32
	PollInterval         time.Duration
	SelectionMode        string
}

func (tc TestNodeConfig) NodeNoNewHeadsThreshold() time.Duration { return tc.NoNewHeadsThreshold }
func (tc TestNodeConfig) NodePollFailureThreshold() uint32       { return tc.PollFailureThreshold }
func (tc TestNodeConfig) NodePollInterval() time.Duration        { return tc.PollInterval }
func (tc TestNodeConfig) NodeSelectionMode() string              { return tc.SelectionMode }

func NewClientWithTestNode(t *testing.T, cfg NodeConfig, rpcUrl string, rpcHTTPURL *url.URL, sendonlyRPCURLs []url.URL, id int32, chainID *big.Int) (*client, error) {
	parsed, err := url.ParseRequestURI(rpcUrl)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, errors.Errorf("ethereum url scheme must be websocket: %s", parsed.String())
	}

	lggr := logger.TestLogger(t)
	n := NewNode(cfg, lggr, *parsed, rpcHTTPURL, "eth-primary-0", id, chainID)
	n.(*node).setLatestReceivedBlockNumber(0)
	primaries := []Node{n}

	var sendonlys []SendOnlyNode
	for i, url := range sendonlyRPCURLs {
		if url.Scheme != "http" && url.Scheme != "https" {
			return nil, errors.Errorf("sendonly ethereum rpc url scheme must be http(s): %s", url.String())
		}
		s := NewSendOnlyNode(lggr, url, fmt.Sprintf("eth-sendonly-%d", i), chainID)
		sendonlys = append(sendonlys, s)
	}

	pool := NewPool(lggr, cfg, primaries, sendonlys, chainID)
	c := &client{logger: lggr, pool: pool}
	t.Cleanup(c.Close)
	return c, nil
}

func Wrap(err error, s string) error {
	return wrap(err, s)
}

type TestableSendOnlyNode interface {
	SendOnlyNode
	SetEthClient(newBatchSender BatchSender, newSender TxSender)
}
