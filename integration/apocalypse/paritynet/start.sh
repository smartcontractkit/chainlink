#!/bin/sh

sleep 20 && parity --config /paritynet/miner.toml --db-path /paritynet/database --bootnodes=enode://8046f1ff008141321e35e27a5ca4f174e28186538d08ee6ad04ea46f909547e28f5ad48ae75528d7d5cad8029a0fb911adcdc8ea36adeb0cc978ccaa0e103f91@gethnet:30303 --network-id=1337 --node-key=21f6447581bbe082cc0e4aea0f5583f6f4d68cbfb37626addfe6c4f5df78d2a8
