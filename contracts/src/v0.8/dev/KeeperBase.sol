// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract KeeperBase {

  error CannotExecute();

  /**
   * @notice method that allows it to be simulated via eth_call by checking that
   * the sender is the zero address.
   */
  function preventExecution() internal view {
    if(tx.origin != address(0)) {
      revert CannotExecute();
    }
  }

  /**
   * @notice modifier that allows it to be simulated via eth_call by checking
   * that the sender is the zero address.
   */
  modifier cannotExecute() {
    preventExecution();
    _;
  }
}
