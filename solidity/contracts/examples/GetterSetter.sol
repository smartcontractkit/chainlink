pragma solidity ^0.4.18;

contract GetterSetter {
  bytes32 public value;
  uint256 public requestId;

  function setValue(uint256 _requestId, bytes32 _value) public {
    requestId = _requestId;
    value = _value;
  }
}
