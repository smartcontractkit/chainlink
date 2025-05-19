// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

contract Counter {
  error AlwaysRevert();

  uint256 public count = 0;

  function increment() public returns (uint256) {
    count += 1;
    return count;
  }

  function reset() public {
    count = 0;
  }

  function alwaysRevert() public pure {
    revert AlwaysRevert();
  }

  function alwaysRevertWithString() public pure {
    revert("always revert");
  }
}
