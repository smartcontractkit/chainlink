pragma solidity 0.5.0;

interface AggregatorInterface {
  function currentAnswer() external view returns (int256);
  function updatedTimestamp() external view returns (uint256);
  function latestRound() external view returns (uint256);
  function getAnswer(uint256 id) external view returns (int256);
}
