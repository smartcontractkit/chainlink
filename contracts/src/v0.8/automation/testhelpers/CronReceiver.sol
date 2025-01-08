// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

contract CronReceiver {
  event Received1();
  event Received2();

  function handler1() external {
    emit Received1();
  }

  function handler2() external {
    emit Received2();
  }

  function revertHandler() external {
    revert("revert!");
  }
}
