pragma solidity ^0.8.0;

import "../ConfirmedOwner.sol";

contract Greeter is ConfirmedOwner {
  string public greeting;

  constructor(address owner) ConfirmedOwner(owner) {}

  function setGreeting(string calldata _greeting) external onlyOwner {
    require(bytes(_greeting).length > 0, "Invalid greeting length");
    greeting = _greeting;
  }

  function triggerRevert() external pure {
    require(false, "Greeter: revert triggered");
  }
}
