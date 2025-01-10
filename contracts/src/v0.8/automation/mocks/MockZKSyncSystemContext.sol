// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

contract MockZKSyncSystemContext {
  function gasPrice() external pure returns (uint256) {
    return 250000000; // 0.25 gwei
  }

  function gasPerPubdataByte() external pure returns (uint256) {
    return 500;
  }

  function getCurrentPubdataSpent() external pure returns (uint256 currentPubdataSpent) {
    return 1000;
  }
}
