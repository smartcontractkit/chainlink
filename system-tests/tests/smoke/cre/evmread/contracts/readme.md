## How to build the contract and generate Go bindings
### 1 Compile the contract
```bash
solc --abi --bin --overwrite -o ./ ./MessageEmitter.sol
```
### 2 Generate Go bindings
```bash
abigen --bin=./MessageEmitter.bin --abi=./MessageEmitter.abi --pkg=contracts --out=./MessageEmitter.go --type MessageEmitter
```

