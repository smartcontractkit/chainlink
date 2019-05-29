pragma solidity 0.4.24;

interface AggregatorInterface {
  function currentAnswer() external returns (int256);
  function updatedHeight() external returns (uint256);
}
