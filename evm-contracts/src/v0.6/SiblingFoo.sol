pragma solidity 0.6.6;

import "./FluxAggregator.sol";

contract SiblingFoo is FluxAggregator {
  uint public bar;

  constructor()
  FluxAggregator(address(0), 0, 0, address(0), 0, 0, 1, "desc")
  public {
    bar = 5;
  }
}
