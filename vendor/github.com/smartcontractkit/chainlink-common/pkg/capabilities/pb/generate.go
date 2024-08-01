//go:generate protoc --go_out=../../ --go_opt=paths=source_relative --go-grpc_out=../../ --go-grpc_opt=paths=source_relative --proto_path=../../ capabilities/pb/capabilities.proto values/pb/values.proto
//go:generate protoc --go_out=../../ --go_opt=paths=source_relative --go-grpc_out=../../ --go-grpc_opt=paths=source_relative --proto_path=../../ capabilities/pb/registry.proto values/pb/values.proto
package pb
