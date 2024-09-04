// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title Forwards a call to a target, under some conditions
interface IForwarder {
  /// @notice forward calls the `target` with `data`
  /// @param target contract address to be called
  /// @param data to send to target contract
  function forward(address target, bytes memory data) external;
}
