//go:generate protoc --proto_path=.:..:. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ocr3_reporting.proto
//go:generate protoc --proto_path=.:..:. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative contract_transmitter.proto ocr3_reporting.proto
package ocr3pb
