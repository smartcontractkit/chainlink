// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../../ChainSpecificUtil.sol";

/// @dev A helper contract that exposes ChainSpecificUtil methods for testing
contract ChainSpecificUtilHelper {
  function getBlockhash(uint64 blockNumber) external view returns (bytes32) {
    return ChainSpecificUtil.getBlockhash(blockNumber);
  }

  function getBlockNumber() external view returns (uint256) {
    return ChainSpecificUtil.getBlockNumber();
  }

  function getCurrentTxL1GasFees(string memory txCallData) external view returns (uint256) {
    return ChainSpecificUtil.getCurrentTxL1GasFees(bytes(txCallData));
  }

  function getL1CalldataGasCost(uint256 calldataSize) external view returns (uint256) {
    return ChainSpecificUtil.getL1CalldataGasCost(calldataSize);
  }
}
