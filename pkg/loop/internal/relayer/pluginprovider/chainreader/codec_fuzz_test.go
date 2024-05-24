package chainreader_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

func FuzzCodec(f *testing.F) {
	interfaceTester := chainreadertest.WrapCodecTesterForLoop(&fakeCodecInterfaceTester{impl: &fakeCodec{}})
	interfacetests.RunCodecInterfaceFuzzTests(f, interfaceTester)
}
