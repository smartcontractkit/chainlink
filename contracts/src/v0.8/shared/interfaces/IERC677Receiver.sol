// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

interface IERC677Receiver {
  function onTokenTransfer(address sender, uint256 amount, bytes calldata data) external;
}
