//go:generate protoc --proto_path=.:..:./v1:./v2:./v3:./v4 --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative mercury_loop.proto
//go:generate protoc --proto_path=.:..:./v1:./v2:./v3:./v4 --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative mercury_plugin.proto
package mercurypb
