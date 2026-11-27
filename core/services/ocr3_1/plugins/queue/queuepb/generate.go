// Package queuepb holds the consensus-queue reporting plugin config that travels in the OCR offchain config.
package queuepb

//go:generate protoc --proto_path=../../../../../../ --go_out=../../../../../../ --go_opt=paths=source_relative core/services/ocr3_1/plugins/queue/queuepb/config.proto
