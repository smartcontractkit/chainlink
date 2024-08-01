package chainreadertest

import (
	"testing"

	"google.golang.org/grpc"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

// WrapCodecTesterForLoop allows you to test a [types.Codec] implementation behind a LOOP server
func WrapCodecTesterForLoop(wrapped interfacetests.CodecInterfaceTester) interfacetests.CodecInterfaceTester {
	return &codecReaderLoopTester{CodecInterfaceTester: wrapped}
}

type codecReaderLoopTester struct {
	interfacetests.CodecInterfaceTester
	lst loopServerTester
}

func (c *codecReaderLoopTester) Setup(t *testing.T) {
	c.CodecInterfaceTester.Setup(t)
	codec := c.CodecInterfaceTester.GetCodec(t)
	c.lst.registerHook = func(server *grpc.Server) {
		if codec != nil {
			impl := chainreader.NewCodecServer(codec)
			pb.RegisterCodecServer(server, impl)
		}
	}
	c.lst.Setup(t)
}

func (c *codecReaderLoopTester) Name() string {
	return c.CodecInterfaceTester.Name() + " on loop"
}

func (c *codecReaderLoopTester) GetCodec(t *testing.T) types.Codec {
	return chainreader.NewCodecClient(nil, c.lst.GetConn(t))
}

func (c *codecReaderLoopTester) IncludeArrayEncodingSizeEnforcement() bool {
	return false
}
