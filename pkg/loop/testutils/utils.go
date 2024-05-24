package testutils

import (
	test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader/test"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
)

// This file exposes functions from pkg/loop/internal/test without exposing internal details.
// the duplication of the function is required so that the test of the LOOP servers themselves
// can dog food the same testers without creating a circular dependency.

// WrapChainReaderTesterForLoop allows you to test a [types.ContractReader] implementation behind a LOOP server
func WrapChainReaderTesterForLoop(wrapped interfacetests.ChainReaderInterfaceTester) interfacetests.ChainReaderInterfaceTester {
	return test.WrapChainReaderTesterForLoop(wrapped)
}

func WrapCodecTesterForLoop(wrapped interfacetests.CodecInterfaceTester) interfacetests.CodecInterfaceTester {
	return test.WrapCodecTesterForLoop(wrapped)
}
