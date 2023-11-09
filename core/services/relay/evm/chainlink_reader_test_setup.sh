#!/bin/bash
set -e
solc --optimize --bin --abi chain_reader_test_contract.sol -o .
mv LatestValueHolder.bin chain_reader_test_contract_gen.bin
mv LatestValueHolder.abi chain_reader_test_contract_gen.abi
../../../../tools/bin/abigen --bin=chain_reader_test_contract_gen.bin --abi=chain_reader_test_contract_gen.abi --pkg=evm_test --out=chain_reader_test_contract_gen_test.go
