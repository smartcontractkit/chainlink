// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IWERC20 {
  function deposit() external payable;

  function withdraw(uint256) external;
}
