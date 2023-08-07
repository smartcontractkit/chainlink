1. go to `core/scripts/vrfv2/testnet/docker` folder and start containers - `docker compose up`
2. execute from `core/scripts/vrfv2/testnet/setup-envs` folder
```
go run . \
--vrf-primary-node-url=http://localhost:6610 \
--vrf-backup-node-url=http://localhost:6611 \
--bhs-node-url=http://localhost:6612 \
--bhf-node-url=http://localhost:6614 \
--creds-file ../docker/secrets/apicredentials 
```

Optional parameters - will not be deployed if specified (NOT WORKING YET)
```
   --link-address 0x8606681e2295B2C4fD44D08E4E8a4D3180071559 \
   --link-eth-feed 0xD46fbB21875EBA71aF1b669c25EEe3692e5B9F13 \
   --bhs-address 0xbD84DbFaC527bd150384B75BA6cD75286F48da28 \
   --batch-bhs-address 0x126AeF1346A81003B8f1B6FBE8742536B2D22F71 \
   --coordinator-address 0x24d42DcD17C92d99100dce9D6133A7496e988A79 \
   --batch-coordinator-address 0xA276bCB4e67d787ce97F8787f2F454A4253Ce0Da 
```