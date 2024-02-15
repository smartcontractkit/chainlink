//go:generate protoc --go_out=../../../../ --go_opt=paths=source_relative --go-grpc_out=../../../../ --go-grpc_opt=paths=source_relative --proto_path=../../../../ capabilities/consensus/ocr3/types/ocr3_types.proto values/pb/values.proto
package types
