pragma solidity >=0.6.0;

interface AggregatorV3Interface {

  function decimals() external view returns (uint8);
  function description() external view returns (string memory);
  function version() external view returns (uint256);

  // getRoundData and latestRoundData should both raise "No data present"
  // if they do not have data to report, instead of returning unset values
  // which could be misinterpreted as actual reported values.
  function getRoundData(uint80 _roundId)
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );
  function latestRoundData()
    external
    view
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    );

  // latestAnswer returns the latest answer, but does not throw if there is not
  // yet any reported answer, so it is impossible to tell a zero answer from a
  // lack of anwer. For this reason, users should take caution and only point to
  // Aggregators that have already come to consensus on at least one answer.
  function latestAnswer()
    external
    view
    returns(
      int256 answer
    );

  event NewRound(
    uint256 indexed roundId,
    address indexed startedBy,
    uint256 startedAt
  );
  event AnswerUpdated(
    int256 indexed current,
    uint256 indexed roundId,
    uint256 updatedAt
  );

}
