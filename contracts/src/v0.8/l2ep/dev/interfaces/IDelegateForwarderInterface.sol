// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

/// @title IDelegateForwarderInterface - forwards a delegatecall to a target, under some conditions
interface IDelegateForwarderInterface {
  /**
   * @notice forward delegatecalls the `target` with `data`
   * @param target contract address to be delegatecalled
   * @param data to send to target contract
   */
  function forwardDelegate(address target, bytes memory data) external;
}
