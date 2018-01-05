pragma solidity ^0.4.18;

contract GetterSetter {
  bytes32 public value;

  function setValue(bytes32 _value) public {
    value = _value;
  }
}
