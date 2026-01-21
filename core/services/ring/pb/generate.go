//go:generate protoc --go_out=. --go_opt=paths=source_relative shared.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative arbiter.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative consensus.proto

package pb
