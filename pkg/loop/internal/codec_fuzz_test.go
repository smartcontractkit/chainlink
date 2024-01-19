package internal_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

func FuzzCodec(f *testing.F) {
	interfaceTester := test.WrapCodecTesterForLoop(&fakeCodecInterfaceTester{impl: &fakeCodec{}})
	interfacetests.RunCodecInterfaceFuzzTests(f, interfaceTester)
}
