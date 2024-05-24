// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title DelegateForwarderInterface - forwards a delegatecall to a target, under some conditions
// solhint-disable-next-line interface-starts-with-i
interface DelegateForwarderInterface {
  /**
   * @notice forward delegatecalls the `target` with `data`
   * @param target contract address to be delegatecalled
   * @param data to send to target contract
   */
  function forwardDelegate(address target, bytes memory data) external;
}
