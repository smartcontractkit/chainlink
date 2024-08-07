// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ChainModuleBase} from "./ChainModuleBase.sol";

contract ZKSyncModule is ChainModuleBase {
  function getGasOverhead()
    external
    pure
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (5_000, 0);
  }
}
