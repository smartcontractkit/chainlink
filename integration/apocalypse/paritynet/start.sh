#!/bin/sh

parity --config /paritynet/miner.toml \
       --chain /paritynet/apocalypse-parity.json \
       --db-path /paritynet/database \
       --bootnodes=enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@172.17.0.4:30303,enode://c1cad3139b0ab583de214e3d64f7fb7793995023559f7fa1e6b01e87603145ca8e60d5d9f8e23d08df3d1c0c82294bd9515b729efec210f060b2fe3a193f9ae0@172.17.0.6:30303 \
       --network-id=1337 \
       --node-key=21f6447581bbe082cc0e4aea0f5583f6f4d68cbfb37626addfe6c4f5df78d2a8
