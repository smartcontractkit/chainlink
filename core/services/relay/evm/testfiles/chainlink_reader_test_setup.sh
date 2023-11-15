#!/bin/bash
set -e
# called from the directory above so cd to this directory
cd testfiles
solc --optimize --bin --abi chain_reader_test_contract.sol -o .
mv LatestValueHolder.bin chain_reader_test_contract_gen.bin
mv LatestValueHolder.abi chain_reader_test_contract_gen.abi
../../../../../tools/bin/abigen --bin=chain_reader_test_contract_gen.bin --abi=chain_reader_test_contract_gen.abi --pkg=testfiles --out=chain_reader_test_contract_gen.go
