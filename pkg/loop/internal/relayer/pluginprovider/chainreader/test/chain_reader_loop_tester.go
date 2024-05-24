package chainreadertest

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

// WrapChainReaderTesterForLoop allows you to test a [types.ContractReader] implementation behind a LOOP server
func WrapChainReaderTesterForLoop(wrapped interfacetests.ChainReaderInterfaceTester) interfacetests.ChainReaderInterfaceTester {
	return &chainReaderLoopTester{ChainReaderInterfaceTester: wrapped}
}

type chainReaderLoopTester struct {
	interfacetests.ChainReaderInterfaceTester
	lst loopServerTester
}

func (c *chainReaderLoopTester) Setup(t *testing.T) {
	c.ChainReaderInterfaceTester.Setup(t)
	chainReader := c.ChainReaderInterfaceTester.GetChainReader(t)
	c.lst.registerHook = func(server *grpc.Server) {
		if chainReader != nil {
			impl := chainreader.NewServer(chainReader)
			pb.RegisterChainReaderServer(server, impl)
		}
	}
	c.lst.Setup(t)
}

func (c *chainReaderLoopTester) GetChainReader(t *testing.T) types.ContractReader {
	return chainreader.NewClient(nil, c.lst.GetConn(t))
}

func (c *chainReaderLoopTester) Name() string {
	return c.ChainReaderInterfaceTester.Name() + " on loop"
}
