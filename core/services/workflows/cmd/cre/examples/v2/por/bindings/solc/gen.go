package solc

//go:generate find . -maxdepth 1 -name "*.sol" -exec solc --abi --bin --overwrite -o bin {} +
