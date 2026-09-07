// Package streamspb holds the generated types for streams.proto (the streams-adapter
// wire contract).
package streamspb

//go:generate protoc --go_out=. --go_opt=paths=source_relative streams.proto
//go:generate protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative streams.proto
