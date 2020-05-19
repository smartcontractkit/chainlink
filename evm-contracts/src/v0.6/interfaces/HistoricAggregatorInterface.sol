pragma solidity >=0.6.0;

interface HistoricAggregatorInterface {
  function latestAnswer() external returns (int256);
  function latestTimestamp() external returns (uint256);
  function latestRound() external returns (uint256);
  function getAnswer(uint256 roundId) external returns (int256);
  function getTimestamp(uint256 roundId) external returns (uint256);

  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp);
  event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt);
}
