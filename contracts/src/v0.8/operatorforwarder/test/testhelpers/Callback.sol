// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Callback {
  address private s_operator;
  uint256 private s_callbacksReceived = 0;

  constructor(address _operator) {
    s_operator = _operator;
  }

  // Callback function for oracle request fulfillment
  function callback(bytes32) public {
    // solhint-disable-next-line gas-custom-errors
    require(msg.sender == s_operator, "Only Operator can call this function");
    s_callbacksReceived += 1;
  }

  function getCallbacksReceived() public view returns (uint256) {
    return s_callbacksReceived;
  }
}
