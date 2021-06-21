package eth

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// secondarynode represents one ethereum node used as a secondary
// It only supports sending transactions
// It must a http(s) url
type secondarynode struct {
	uri    url.URL
	rpc    *rpc.Client
	geth   *ethclient.Client
	log    *logger.Logger
	dialed bool
}

func newSecondaryNode(httpuri url.URL, name string) (s *secondarynode) {
	s = new(secondarynode)
	s.log = logger.CreateLogger(logger.Default.With(
		"nodeName", name,
		"nodeTier", "secondary",
	))
	s.uri = httpuri
	return
}

func (s *secondarynode) Dial() error {
	s.log.Debugw("eth.Client#Dial(...)")
	if s.dialed {
		panic("eth.Client.Dial(...) should only be called once during the node's lifetime.")
	}

	uri := s.uri.String()
	rpc, err := rpc.DialHTTP(uri)
	if err != nil {
		return err
	}
	s.dialed = true
	s.rpc = rpc
	s.geth = ethclient.NewClient(rpc)
	return nil
}

func (s secondarynode) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	s.log.Debugw("eth.Client#SendTransaction(...)",
		"tx", tx,
	)
	return s.wrap(s.geth.SendTransaction(ctx, tx))
}

func (s secondarynode) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	s.log.Debugw("eth.Client#BatchCall(...)",
		"nBatchElems", len(b),
	)
	return s.wrap(s.rpc.BatchCallContext(ctx, b))
}

func (s secondarynode) wrap(err error) error {
	return wrap(err, fmt.Sprintf("secondary http (%s)", s.uri.String()))
}
