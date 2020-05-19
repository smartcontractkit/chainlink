pragma solidity ^0.6.0;

import "../interfaces/AggregatorInterface.sol";

/**
 * @title The MockAggregator contract
 * @notice Use this contract when you need to test
 * other contract's ability to read data from an
 * aggregator contract, but how the aggregator got
 * its answer is unimportant
 */
contract MockAggregator is AggregatorInterface {
  uint8 public override decimals;
  int256 public override latestAnswer;
  uint256 public override latestTimestamp;
  uint256 public override latestRound;

  mapping(uint256 => int256) public override getAnswer;
  mapping(uint256 => uint256) public override getTimestamp;
  mapping(uint256 => uint256) private getStartedAt;

  constructor(
    uint8 _decimals,
    int256 _initialAnswer
  ) public {
    decimals = _decimals;
    updateAnswer(_initialAnswer);
  }

  function updateAnswer(
    int256 _answer
  ) public {
    latestAnswer = _answer;
    latestTimestamp = block.timestamp;
    latestRound++;
    getAnswer[latestRound] = _answer;
    getTimestamp[latestRound] = block.timestamp;
  }

  function updateRoundData(
    uint256 _roundId,
    int256 _answer,
    uint256 _timestamp,
    uint256 _startedAt
  ) public {
    latestRound = _roundId;
    latestAnswer = _answer;
    latestTimestamp = _timestamp;
    getAnswer[latestRound] = _answer;
    getTimestamp[latestRound] = _timestamp;
    getStartedAt[latestRound] = _startedAt;
  }

  function getRoundData(uint256 _roundId)
    external
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return (
      _roundId,
      getAnswer[_roundId],
      getStartedAt[_roundId],
      getTimestamp[_roundId],
      _roundId
    );
  }

  function latestRoundData()
    external
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return (
      latestRound,
      getAnswer[latestRound],
      getStartedAt[latestRound],
      getTimestamp[latestRound],
      latestRound
    );
  }
}
