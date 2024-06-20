// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {ChainModuleBase} from "./ChainModuleBase.sol";

ISystemContext constant SYSTEM_CONTEXT_CONTRACT = ISystemContext(address(0x800b));

interface ISystemContext {
  function gasPrice() external view returns (uint256);
}

contract ZKSyncModule is ChainModuleBase {
  uint256 private constant FIXED_GAS_OVERHEAD = 5_000;

  function getMaxL1Fee(uint256 maxCalldataSize) external view override returns (uint256) {
    return maxCalldataSize * SYSTEM_CONTEXT_CONTRACT.gasPrice();
  }

  function getGasOverhead()
    external
    view
    override
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (FIXED_GAS_OVERHEAD, 0);
  }
}
