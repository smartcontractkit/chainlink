pragma solidity ^0.4.18;

contract GetterSetter {
  bytes32 public value;
  uint256 public nonce;

  function setValue(uint256 _nonce, bytes32 _value) public {
    nonce = _nonce;
    value = _value;
  }
}
