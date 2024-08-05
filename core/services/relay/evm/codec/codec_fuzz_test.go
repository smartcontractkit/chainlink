package codec_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

func FuzzCodec(f *testing.F) {
	tester := &codecInterfaceTester{}
	interfacetests.RunCodecInterfaceFuzzTests(f, tester)
}
