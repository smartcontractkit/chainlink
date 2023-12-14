// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface IERC677 {
  event Transfer(address indexed from, address indexed to, uint256 value, bytes data);

  /// @notice Transfer tokens from `msg.sender` to another address and then call `onTransferReceived` on receiver
  /// @param to The address which you want to transfer to
  /// @param amount The amount of tokens to be transferred
  /// @param data bytes Additional data with no specified format, sent in call to `to`
  /// @return true unless throwing
  function transferAndCall(address to, uint256 amount, bytes memory data) external returns (bool);
}
