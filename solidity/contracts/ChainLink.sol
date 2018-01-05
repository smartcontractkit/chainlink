pragma solidity ^0.4.17;

contract ChainLink {

  bytes32 public value;

  function ChainLink() public {
    value = "Hello World!";
  }

  function setValue(bytes32 newValue) public {
    value = newValue;
  }
}
