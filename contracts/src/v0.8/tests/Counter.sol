// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract Counter {
  uint256 public count = 0;

  function increment() public returns (uint256) {
    count += 1;
    return count;
  }

  function reset() public returns (uint256) {
    count = 0;
    return count;
  }

  function alwaysRevert() public pure {
    revert("always revert");
  }
}
