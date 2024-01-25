// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.16;

import "../interfaces/v2_2/IChainSpecific.sol";

contract ChainModuleBase is IChainSpecific {
  function blockNumber() external view returns (uint256) {
    return block.number;
  }

  function blockHash(uint256 blockNumber) external view returns (bytes32) {
    return blockhash(blockNumber);
  }

  function getL1Fee(bytes calldata) external pure returns (uint256) {
    return 0;
  }

  function getMaxL1Fee(uint256) external pure returns (uint256) {
    return 0;
  }
}
