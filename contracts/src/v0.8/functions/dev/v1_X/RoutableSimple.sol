// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

/// @title This abstract should be inherited by contracts that will be used
/// as the destinations to a route (id=>contract) on the Router.
/// It provides a Router getter and modifiers.
contract Routable {
  function add(uint256 a, uint256 b) public pure returns (uint256) {
    return a + b;
  }
}
