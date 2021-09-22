package eth

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// SendOnlyNode represents one ethereum node used as a sendonly
type SendOnlyNode interface {
	Dial(context.Context) error
	Verify(ctx context.Context, expectedChainID *big.Int) (err error)

	SendTransaction(ctx context.Context, tx *types.Transaction) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error

	String() string
}

// It only supports sending transactions
// It must a http(s) url
type sendOnlyNode struct {
	uri    url.URL
	rpc    *rpc.Client
	geth   *ethclient.Client
	log    logger.Logger
	dialed bool
	name   string
}

func NewSendOnlyNode(lggr logger.Logger, httpuri url.URL, name string) SendOnlyNode {
	s := new(sendOnlyNode)
	s.name = name
	s.log = lggr.With(
		"nodeName", name,
		"nodeTier", "sendonly",
	)
	s.uri = httpuri
	return s
}

func (s *sendOnlyNode) Dial(_ context.Context) error {
	s.log.Debugw("eth.Client#Dial(...)")
	if s.dialed {
		panic("eth.Client.Dial(...) should only be called once during the node's lifetime.")
	}

	uri := s.uri.String()
	rpc, err := rpc.DialHTTP(uri)
	if err != nil {
		return errors.Wrapf(err, "failed to dial secondary client: %v", uri)
	}
	s.dialed = true
	s.rpc = rpc
	s.geth = ethclient.NewClient(rpc)
	return nil
}

func (s sendOnlyNode) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	s.log.Debugw("eth.Client#SendTransaction(...)",
		"tx", tx,
	)
	return s.wrap(s.geth.SendTransaction(ctx, tx))
}

func (s sendOnlyNode) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	s.log.Debugw("eth.Client#BatchCall(...)",
		"nBatchElems", len(b),
	)
	return s.wrap(s.rpc.BatchCallContext(ctx, b))
}

func (s sendOnlyNode) ChainID(ctx context.Context) (chainID *big.Int, err error) {
	s.log.Debugw("eth.Client#ChainID(...)")
	chainID, err = s.geth.ChainID(ctx)
	err = s.wrap(err)
	return
}

func (s sendOnlyNode) wrap(err error) error {
	return wrap(err, fmt.Sprintf("sendonly http (%s)", s.uri.String()))
}

func (s sendOnlyNode) String() string {
	return fmt.Sprintf("(secondary)%s:%s", s.name, s.uri.String())
}

func (s sendOnlyNode) Verify(ctx context.Context, expectedChainID *big.Int) (err error) {
	if chainID, err := s.ChainID(ctx); err != nil {
		return errors.Wrap(err, "failed to verify chain ID")
	} else if chainID.Cmp(expectedChainID) != 0 {
		return errors.Errorf(
			"sendonly rpc ChainID doesn't match local chain ID: RPC ID=%s, local ID=%s, node name=%s",
			chainID.String(),
			expectedChainID.String(),
			s.name,
		)
	}
	return nil
}
