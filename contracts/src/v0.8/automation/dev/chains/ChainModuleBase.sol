// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IChainModule} from "../interfaces/v2_2/IChainModule.sol";

contract ChainModuleBase is IChainModule {
  function blockNumber() external virtual view returns (uint256) {
    return block.number;
  }

  function blockHash(uint256 n) external virtual view returns (bytes32) {
    if (n >= block.number || block.number - n > 256) {
      return "";
    }
    return blockhash(n);
  }

  function getCurrentL1Fee() external virtual view returns (uint256) {
    return 0;
  }

  function getMaxL1Fee(uint256) external virtual view returns (uint256) {
    return 0;
  }
}
