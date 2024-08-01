// NOTE: the relative paths in the proto_path are to ensure we find common utilities, like BigInt
//go:generate protoc --proto_path=.:../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative datasource_v3.proto
//go:generate protoc --proto_path=.:../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative reportcodec_v3.proto

package mercuryv3pb
