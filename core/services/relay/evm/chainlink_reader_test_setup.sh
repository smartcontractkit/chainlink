#!/bin/bash
solc --optimize --bin --abi chain_reader_test_contract.sol -o . && mv LatestValueHolder.bin chain_reader_test_contract.bin && mv LatestValueHolder.abi chain_reader_test_contract.abi
