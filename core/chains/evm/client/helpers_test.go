package client

import (
	"fmt"
	"math/big"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
)

func init() {
	dialRetryInterval = 100 * time.Millisecond
}

func NewClient(lggr logger.Logger, rpcUrl string, rpcHTTPURL *url.URL, sendonlyRPCURLs []url.URL, chainID *big.Int) (*client, error) {
	parsed, err := url.ParseRequestURI(rpcUrl)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, errors.Errorf("ethereum url scheme must be websocket: %s", parsed.String())
	}

	primaries := []Node{NewNode(lggr, *parsed, rpcHTTPURL, "eth-primary-0")}

	var sendonlys []SendOnlyNode
	for i, url := range sendonlyRPCURLs {
		if url.Scheme != "http" && url.Scheme != "https" {
			return nil, errors.Errorf("sendonly ethereum rpc url scheme must be http(s): %s", url.String())
		}
		s := NewSendOnlyNode(lggr, url, fmt.Sprintf("eth-sendonly-%d", i))
		sendonlys = append(sendonlys, s)
	}

	pool := NewPool(lggr, primaries, sendonlys, chainID)
	return &client{logger: lggr, pool: pool}, nil
}

func Wrap(err error, s string) error {
	return wrap(err, s)
}
