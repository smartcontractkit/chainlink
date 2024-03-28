The code herein is the gRPC wrapper for the CCIP Reporting Plugin `Provider`s

Developing the CCIP Providers requires the follow workflow

1. Update the golang interfaces in `chainlink-common/pkg/types`, and/or  `chainlink-common/pkg/types/ccip`
2. Mirror those changes as protobuf changes in `chainlink-common/pkg/loop/internal/pb/ccip`
    note: golang `interface` <-> gRPC `service` and `struct` <-> `message`
    remember to run `make generate` at the top of the repo to generate your changes
    see docs in `internal/pb` for more details
3. If you are adding a new interface, create a file here and implement the gRPC server and client. Follow the testing pattern in the  `test` subdir to ensure roundtrip de/serialization of each and every rpc
4. If you are updating an interface, then update both the `GRPCClient` and `GRPCServer` in the relevant file.
    - The client code will not compile if you change the interface in `chainlink-common/pkg/types/ccip` and so these are obvious
    - However the server code will compile and you will receive runtime `Unimplemented` errors (this is a consequence of golang gRPC best practices that embed a `UnimplementedServer`. the curious reader can go here)
    - Add support for your new method in the static test implementation found in `test` and add a test case to the `roundTrip` test for that interface. This roundtrip test is the key to prevent runtime Unimplemented errors

Example: you add method to `OnRampReader`

- [edit the interface](https://github.com/smartcontractkit/chainlink-common/blob/d529c556ff6b98573d6e442b91d26d322c1066ac/pkg/types/ccip/onramp.go#L53)
- [edit the pb and regenerate via `make generate`](https://github.com/smartcontractkit/chainlink-common/blob/d529c556ff6b98573d6e442b91d26d322c1066ac/pkg/loop/internal/pb/ccip/onramp.proto#L1)
- [update the client by adding the func](https://github.com/smartcontractkit/chainlink-common/blob/d529c556ff6b98573d6e442b91d26d322c1066ac/pkg/loop/internal/relayer/pluginprovider/ext/ccip/onramp.go#L19)
- [update the server by the adding the func](https://github.com/smartcontractkit/chainlink-common/blob/d529c556ff6b98573d6e442b91d26d322c1066ac/pkg/loop/internal/relayer/pluginprovider/ext/ccip/onramp.go#L101)
- [update the test implementation by adding the func and test data](https://github.com/smartcontractkit/chainlink-common/blob/d529c556ff6b98573d6e442b91d26d322c1066ac/pkg/loop/internal/relayer/pluginprovider/ext/ccip/test/onramp.go#L16)
- [update the grpc test by adding a test case leveraging the test implementation](https://github.com/smartcontractkit/chainlink-common/blob/d529c556ff6b98573d6e442b91d26d322c1066ac/pkg/loop/internal/relayer/pluginprovider/ext/ccip/test/offramp_test.go#L101)

When adding interface methods, maintain alphabetical order of there declaration and implementations. You will thank yourself later! You will find yourself reading multiple files that define or implement the interfaces and this simply organizational principle makes that much easier when there are more than two or three funcs. It really pays off in the tests, where you can 

same idea in current code base: https://github.com/smartcontractkit/chainlink-common/commit/c15097958cbdd27943c474052b4bf89c845c2673