// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

interface ICommonAggregator {
  function description() external view returns (string memory);

  function version() external view returns (uint256);
}
