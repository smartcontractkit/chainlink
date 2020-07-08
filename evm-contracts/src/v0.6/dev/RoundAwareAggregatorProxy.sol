pragma solidity 0.6.6;


import './AggregatorProxy.sol';


contract RoundAwareAggregatorProxy is AggregatorProxy {

  uint256 constant private EPOCH_OFFSET = 32;
  uint256 constant private EPOCH_BASE = 2 ** EPOCH_OFFSET;
  uint256 constant private EPOCH_MASK = 0xFFFF << EPOCH_OFFSET;
  uint16 public epoch;
  mapping(uint16 => AggregatorV3Interface) public epochAggregators;

  constructor(address _aggregator)
    public
    AggregatorProxy(_aggregator)
  {
  }

  function latestRoundData()
    public
    view
    virtual
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = aggregator.latestRoundData();
    return (addEpoch(roundId), answer, startedAt, updatedAt, answeredInRound);
  }

  function getRoundData(uint256 requestId)
    public
    view
    virtual
    override
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    uint16 reqEpoch;
    uint256 reqRound;
    (reqEpoch, reqRound) = parseRequestId(requestId);
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = epochAggregators[reqEpoch].getRoundData(reqRound);
    return (requestId, answer, startedAt, updatedAt, answeredInRound);
  }


  // INTERNAL

  function setAggregator(address _aggregator)
    internal
    override
  {
    epoch++;
    epochAggregators[epoch] = AggregatorV3Interface(_aggregator);
    aggregator = AggregatorV3Interface(_aggregator);
  }


  // PRIVATE

  function addEpoch(
    uint256 originalId
  )
    private
    view
    returns (uint256)
  {
    return (epoch * EPOCH_BASE) + originalId;
  }

  function parseRequestId(
    uint256 requestId
  )
    private
    view
    returns (uint16, uint256)
  {
    uint256 offsetEpochId = EPOCH_MASK & requestId;
    uint16 epochId = uint16(offsetEpochId >> EPOCH_OFFSET);

    uint256 requestIdMask = (2**EPOCH_OFFSET) - 1;
    uint256 roundId = requestId & requestIdMask;

    return (epochId, roundId);
  }

}

