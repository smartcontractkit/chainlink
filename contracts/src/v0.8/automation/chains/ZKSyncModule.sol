// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ChainModuleBase} from "./ChainModuleBase.sol";
import {SYSTEM_CONTEXT_CONTRACT} from "../interfaces/zksync/ISystemContext.sol";

contract ZKSyncModule is ChainModuleBase {
  uint256 private constant FIXED_GAS_OVERHEAD = 5_000;

  function getMaxL1Fee(uint256 gasLimit) external view override returns (uint256) {
    return gasLimit * SYSTEM_CONTEXT_CONTRACT.gasPrice();
  }

  function getGasOverhead()
    external
    pure
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (FIXED_GAS_OVERHEAD, 0);
  }
}
