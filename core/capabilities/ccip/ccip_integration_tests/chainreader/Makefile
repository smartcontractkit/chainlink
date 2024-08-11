
# IMPORTANT: If you encounter any issues try using solc 0.8.18 and abigen 1.14.5

.PHONY: build
build:
	rm -rf build/
	solc --evm-version paris --abi --bin mycontract.sol -o build
	abigen --abi build/mycontract_sol_SimpleContract.abi --bin build/mycontract_sol_SimpleContract.bin --pkg=chainreader --out=mycontract.go

.PHONY: test
test: build
	go test -v --tags "playground" ./...
