// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {ArbSys} from "./vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol";
import {ChainSpecificUtil} from "./ChainSpecificUtil.sol";

contract TestArbitrum {
    function getBlockhash(uint64 lookback) external returns (bytes32) {
        return ChainSpecificUtil.getBlockhash(uint64(ChainSpecificUtil.getBlockNumber()) - lookback);
    }
}