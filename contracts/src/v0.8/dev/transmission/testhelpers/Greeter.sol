// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

/// @dev Ownerless greeter contract.
contract Greeter {
  string private s_greeting;

  function setGreeting(string memory greeting) external {
    s_greeting = greeting;
  }

  function getGreeting() external view returns (string memory) {
    return s_greeting;
  }
}
