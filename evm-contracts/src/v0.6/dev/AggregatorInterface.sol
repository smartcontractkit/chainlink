pragma solidity ^0.6.2;

interface AggregatorInterface {
  function latestAnswer() external view virtual returns (int256);
  function latestTimestamp() external view virtual returns (uint256);
  function latestRound() external view virtual returns (uint256);
  function getAnswer(uint256 roundId) external view virtual returns (int256);
  function getTimestamp(uint256 roundId) external view virtual returns (uint256);

  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp);
  event NewRound(uint256 indexed roundId, address indexed startedBy);
}
