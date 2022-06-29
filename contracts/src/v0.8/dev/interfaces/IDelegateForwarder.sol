// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title IDelegateForwarder - forwards a delegatecall to a target, under some conditions
interface IDelegateForwarder {
  /**
   * @notice forward delegatecalls the `target` with `data`
   * @param target contract address to be delegatecalled
   * @param data to send to target contract
   */
  function forwardDelegate(address target, bytes memory data) external;
}
