pragma solidity 0.5.0;

/**
 * @title The Owned contract
 * @notice A contract with helpers for basic contract ownership.
 */
contract Owned {

  address public owner;

  event OwnershipTransferRequested(address to, address from);

  constructor() public {
    owner = msg.sender;
  }

  function transferOwnership(address _to)
    public
    onlyOwner()
  {
    emit OwnershipTransferRequested(_to, owner);
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
