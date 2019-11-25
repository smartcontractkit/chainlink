pragma solidity 0.4.24;

interface AggregatorInterface {
  function currentAnswer() external view returns (int256);
  function updatedTimestamp() external view returns (uint256);
}
