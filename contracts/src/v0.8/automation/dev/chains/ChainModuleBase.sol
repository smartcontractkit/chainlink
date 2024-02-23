// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IChainModule} from "../interfaces/v2_2/IChainModule.sol";

contract ChainModuleBase is IChainModule {
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
    return (0, 0);
  }
}
