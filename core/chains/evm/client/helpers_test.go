package client

import (
	"fmt"
	"math/big"
	"net/url"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TestNodeConfig struct {
	NoNewHeadsThreshold  time.Duration
	PollFailureThreshold uint32
	PollInterval         time.Duration
	SelectionMode        string
	SyncThreshold        uint32
}

func (tc TestNodeConfig) NodeNoNewHeadsThreshold() time.Duration { return tc.NoNewHeadsThreshold }
func (tc TestNodeConfig) NodePollFailureThreshold() uint32       { return tc.PollFailureThreshold }
func (tc TestNodeConfig) NodePollInterval() time.Duration        { return tc.PollInterval }
func (tc TestNodeConfig) NodeSelectionMode() string              { return tc.SelectionMode }
func (tc TestNodeConfig) NodeSyncThreshold() uint32              { return tc.SyncThreshold }

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
	n.(*node).setLatestReceived(0, utils.NewBigI(0))
	primaries := []Node{n}

	var sendonlys []SendOnlyNode
	for i, url := range sendonlyRPCURLs {
		if url.Scheme != "http" && url.Scheme != "https" {
			return nil, errors.Errorf("sendonly ethereum rpc url scheme must be http(s): %s", url.String())
		}
		s := NewSendOnlyNode(lggr, url, fmt.Sprintf("eth-sendonly-%d", i), chainID)
		sendonlys = append(sendonlys, s)
	}

	pool := NewPool(lggr, cfg, primaries, sendonlys, chainID, "")
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

const HeadResult = `{"difficulty":"0xf3a00","extraData":"0xd883010503846765746887676f312e372e318664617277696e","gasLimit":"0xffc001","gasUsed":"0x0","hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a","logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","miner":"0xd1aeb42885a43b72b518182ef893125814811048","mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e","nonce":"0x0ece08ea8c49dfd9","number":"0x1","parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d","receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347","size":"0x218","stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b","timestamp":"0x58318da2","totalDifficulty":"0x1f3a00","transactions":[],"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","uncles":[]}`

func IsDialed(s SendOnlyNode) bool {
	return s.(*sendOnlyNode).dialed
}
