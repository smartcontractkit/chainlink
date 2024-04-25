// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract MockReceiver {
  uint256 public value;

  function receiveData(uint256 _value) public {
    value = _value;
  }

  function revertMessage() public pure {
    revert("test revert message");
  }
}
