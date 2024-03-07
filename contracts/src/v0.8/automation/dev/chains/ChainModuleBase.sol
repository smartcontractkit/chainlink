// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IChainModule} from "../interfaces/IChainModule.sol";

contract ChainModuleBase is IChainModule {
  uint256 private constant FIXED_GAS_OVERHEAD = 300;
  uint256 private constant PER_CALLDATA_BYTE_GAS_OVERHEAD = 0;

  function blockNumber() external view virtual returns (uint256) {
    return block.number;
  }

  function blockHash(uint256 n) external view virtual returns (bytes32) {
    if (n >= block.number || block.number - n > 256) {
      return "";
    }
    return blockhash(n);
  }

  function getCurrentL1Fee() external view virtual returns (uint256) {
    return 0;
  }

  function getMaxL1Fee(uint256) external view virtual returns (uint256) {
    return 0;
  }

  function getGasOverhead()
    external
    view
    virtual
    returns (uint256 chainModuleFixedOverhead, uint256 chainModulePerByteOverhead)
  {
    return (FIXED_GAS_OVERHEAD, PER_CALLDATA_BYTE_GAS_OVERHEAD);
  }
}
