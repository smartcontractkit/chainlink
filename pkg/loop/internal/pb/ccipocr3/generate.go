//go:generate protoc --proto_path=.:.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative models.proto
//go:generate protoc --proto_path=.:.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative msghasher.proto
//go:generate protoc --proto_path=.:.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative reportcodec.proto

package ccipocr3pb
