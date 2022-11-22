#!/bin/zsh

./contracts/scripts/native_solc8_6_compile dev/ocr2dr/OCR2DRClientExample.sol
./contracts/scripts/native_solc8_6_compile dev/ocr2dr/OCR2DROracle.sol
./contracts/scripts/native_solc8_6_compile dev/ocr2dr/OCR2DRClient.sol
./contracts/scripts/native_solc8_6_compile dev/ocr2dr/OCR2DRRegistry.sol

cd core/gethwrappers
go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/OCR2DRClientExample.abi ../../contracts/solc/v0.8.6/OCR2DRClientExample.bin OCR2DRClientExample ocr2dr_client_example
go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/OCR2DRClient.abi ../../contracts/solc/v0.8.6/OCR2DRClient.bin OCR2DRClient ocr2dr_client
go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/OCR2DROracle.abi ../../contracts/solc/v0.8.6/OCR2DROracle.bin OCR2DROracle ocr2dr_oracle
go run ./generation/generate/wrap.go ../../contracts/solc/v0.8.6/OCR2DRRegistry.abi ../../contracts/solc/v0.8.6/OCR2DRRegistry.bin OCR2DRRegistry ocr2dr_registry
