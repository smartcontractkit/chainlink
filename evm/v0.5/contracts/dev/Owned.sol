pragma solidity 0.5.0;

/**
 * @title The Owned contract
 * @notice A contract with helpers for basic contract ownership.
 */
contract Owned {

  address public owner;

  constructor() public {
    owner = msg.sender;
  }

  modifier onlyOwner() {
    require(msg.sender == owner, "Only callable by owner");
    _;
  }

  modifier ifOwner() {
    if (msg.sender == owner) {
      _;
    }
  }

}
