// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Callback {
  address private operator;
  uint256 private callbacksReceived = 0;

  constructor(address _operator) {
    operator = _operator;
  }

  // Callback function for oracle request fulfillment
  function callback(bytes32) public {
    require(msg.sender == operator, "Only Operator can call this function");
    callbacksReceived += 1;
  }

  function getCallbacksReceived() public view returns (uint256) {
    return callbacksReceived;
  }
}
