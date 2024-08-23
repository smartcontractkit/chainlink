// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

contract Dummy {
  uint64 internal s_counter;

  function increment() external {
    s_counter += 1;
  }
}
