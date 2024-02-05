// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.16;

import {IChainModule} from "../automation/dev/interfaces/v2_2/IChainModule.sol";

contract MockChainModule is IChainModule {
  //uint256 internal blockNum;

  function blockNumber() external view returns (uint256) {
    return 1256;
  }

  function blockHash(uint256 blocknumber) external view returns (bytes32) {
    require(1000 >= blocknumber, "block too old");

    return keccak256(abi.encode(blocknumber));
  }

  function getCurrentL1Fee() external view returns (uint256) {
    return 0;
  }

  // retrieve the L1 data fee for a L2 simulation. it should return 0 for L1 chains and
  // L2 chains which don't have L1 fee component.
  function getMaxL1Fee(uint256 dataSize) external view returns (uint256) {
    return 0;
  }
}
