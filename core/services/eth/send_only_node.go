package eth

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// SendOnlyNode represents one ethereum node used as a sendonly
// It only supports sending transactions
// It must a http(s) url
type SendOnlyNode struct {
	uri    url.URL
	rpc    *rpc.Client
	geth   *ethclient.Client
	log    *logger.Logger
	dialed bool
	name   string
}

func NewSendOnlyNode(lggr *logger.Logger, httpuri url.URL, name string) (s *SendOnlyNode) {
	s = new(SendOnlyNode)
	s.name = name
	s.log = lggr.With(
		"nodeName", name,
		"nodeTier", "sendonly",
	)
	s.uri = httpuri
	return
}

func (s *SendOnlyNode) Dial() error {
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

func (s SendOnlyNode) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	s.log.Debugw("eth.Client#SendTransaction(...)",
		"tx", tx,
	)
	return s.wrap(s.geth.SendTransaction(ctx, tx))
}

func (s SendOnlyNode) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	s.log.Debugw("eth.Client#BatchCall(...)",
		"nBatchElems", len(b),
	)
	return s.wrap(s.rpc.BatchCallContext(ctx, b))
}

func (s SendOnlyNode) ChainID(ctx context.Context) (chainID *big.Int, err error) {
	s.log.Debugw("eth.Client#ChainID(...)")
	chainID, err = s.geth.ChainID(ctx)
	err = s.wrap(err)
	return
}

func (s SendOnlyNode) wrap(err error) error {
	return wrap(err, fmt.Sprintf("sendonly http (%s)", s.uri.String()))
}
