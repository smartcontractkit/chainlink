// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IChainModule} from "../interfaces/v2_2/IChainModule.sol";

contract ChainModuleBase is IChainModule {
  //block number: 26,460
  //block hash: 26,917
  //get current l1 fee (0 bytes): 26,989
  // => 80366
  uint256 private constant FIXED_GAS_OVERHEAD = 20000;

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
    return (FIXED_GAS_OVERHEAD, 0);
  }
}
