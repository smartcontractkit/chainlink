pragma solidity ^0.6.0;

interface AggregatorInterface {
  /*
   * Historic Aggregator (a.k.a. Solidity v0.4 Aggregator, Aggregator version 2)
   * functions
   */
  function latestAnswer() external view returns (int256);
  function latestTimestamp() external view returns (uint256);
  function latestRound() external view returns (uint256);
  function getAnswer(uint256 roundId) external view returns (int256);
  function getTimestamp(uint256 roundId) external view returns (uint256);

  /*
   * FluxAggregator functions
   */
  function getRound(uint256 _roundId)
    external
    view
    returns (
      uint256 roundId,
      int256 answer,
      uint64 startedAt,
      uint64 updatedAt,
      uint256 answeredInRound
    );
  function decimals() external view returns (uint8);

  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 timestamp);
  event NewRound(uint256 indexed roundId, address indexed startedBy, uint256 startedAt);
}
