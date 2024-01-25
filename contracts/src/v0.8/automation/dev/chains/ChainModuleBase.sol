// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import {IChainModule} from "../interfaces/v2_2/IChainModule.sol";

contract ChainModuleBase is IChainModule {
  function blockNumber() external view returns (uint256) {
    return block.number;
  }

  function blockHash(uint256 blocknumber) external view returns (bytes32) {
    return blockhash(blocknumber);
  }

  function getL1Fee(bytes calldata) external pure returns (uint256) {
    return 0;
  }

  function getMaxL1Fee(uint256) external pure returns (uint256) {
    return 0;
  }
}
