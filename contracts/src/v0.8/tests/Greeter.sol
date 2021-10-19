pragma solidity ^0.8.0;

import "../ConfirmedOwner.sol";

contract Greeter is ConfirmedOwner(msg.sender) {
  string public greeting;

  function setGreeting(string calldata _greeting) external onlyOwner {
    greeting = _greeting;
  }
}
