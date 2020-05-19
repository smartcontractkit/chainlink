pragma solidity >=0.5.0;

interface AggregatorInterface {
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp);
  event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt);

  function latestAnswer() external returns (int256);
  function latestTimestamp() external returns (uint256);
  function latestRound() external returns (uint256);
  function getAnswer(uint256 roundId) external returns (int256);
  function getTimestamp(uint256 roundId) external returns (uint256);

  // post-Historic

  function decimals() external returns (uint8);
  function getRoundData(uint256 _roundId)
    external
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    );
  function latestRoundData()
    external
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    );
}
