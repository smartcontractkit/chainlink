.PHONY: build_consumer
build_consumer:
	solc --abi --overwrite -o abi permissionless_feeds_consumer/PermissionlessFeedsConsumer.sol
	solc --bin --overwrite -o bin permissionless_feeds_consumer/PermissionlessFeedsConsumer.sol
	abigen --bin=bin/PermissionlessFeedsConsumer.bin --abi=abi/PermissionlessFeedsConsumer.abi --pkg=permissionless_feeds_consumer --out=permissionless_feeds_consumer/permissionless_feeds_consumer.go
