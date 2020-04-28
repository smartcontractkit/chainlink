

# Apocalypse stress test

## Accounts

- Geth #1 (miner)
    - Persona: geth1
    - Account: 0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f
    - Enode: enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@172.17.0.3:30303
    - IP: 172.17.0.3
- Geth #2 (miner)
    - Persona: geth2
    - Account: 0x7db75251a74f40b15631109ba44d33283ed48528
    - Enode: 
- Parity (transaction relayer, does not mine): 0xde554b6c292f5e5794a68dc560a537dd89d3b03e



curl --data '{"method":"parity_addReservedPeer","params":["enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@172.17.0.3:30303"],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:28545

curl --data '{"method":"parity_addReservedPeer","params":["enode://c1cad3139b0ab583de214e3d64f7fb7793995023559f7fa1e6b01e87603145ca8e60d5d9f8e23d08df3d1c0c82294bd9515b729efec210f060b2fe3a193f9ae0@172.17.0.4:30303"],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:28545


curl --data '{"method":"parity_netPeers","params":[],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:28545

enode://9f87a48e17c61f00d1dbbfe3692d890359ccb811d10ef21517913aeb05f47947ce3c9ac6bf96ae75be527f1b1c301a7dc17c680f0bea90f33e28dd8d76d26dbe@172.17.0.3:30303


COIN=ETH \
ETHEREUM_JSONRPC_VARIANT=geth \ 
ETHEREUM_JSONRPC_HTTP_URL=http://localhost:8545 \
ETHEREUM_JSONRPC_WS_URL=ws://localhost:8546 \
make start


curl --data '{"method":"txpool_content","params":[],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:8545



for (let i = 0; i < 1000; i++) { web3.eth.sendTransaction({ from: '0xDe554B6c292f5e5794A68Dc560a537DD89d3b03E', to: '0xDe554B6c292f5e5794A68Dc560a537DD89d3b03E', value: 1, gas: 30000 }) }
