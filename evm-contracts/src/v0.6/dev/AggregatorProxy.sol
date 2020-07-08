pragma solidity 0.6.6;

import "../Owned.sol";
import "../interfaces/AggregatorInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import "../vendor/SafeMath.sol";

/**
 * @title A trusted proxy for updating where current answers are read from
 * @notice This contract provides a consistent address for the
 * CurrentAnwerInterface but delegates where it reads from to the owner, who is
 * trusted to update it.
 */
contract AggregatorProxy is AggregatorInterface, AggregatorV3Interface, Owned {
  using SafeMath for uint256;

  uint16 public epoch;
  AggregatorV3Interface public aggregator;
  AggregatorV3Interface public proposedAggregator;
  mapping(uint16 => AggregatorV3Interface) public epochAggregators;

  uint256 constant private EPOCH_OFFSET = 32;
  uint256 constant private EPOCH_BASE = 2 ** EPOCH_OFFSET;
  uint256 constant private EPOCH_MASK = 0xFFFF << EPOCH_OFFSET;

  constructor(address _aggregator) public Owned() {
    setAggregator(_aggregator);
  }

  /**
   * @notice Reads the current answer from aggregator delegated to.
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestAnswer()
    public
    view
    virtual
    override
    returns (int256 answer)
  {
    ( , answer, , , ) = latestRoundData();
  }

  /**
   * @notice Reads the last updated height from aggregator delegated to.
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestTimestamp()
    public
    view
    virtual
    override
    returns (uint256 updatedAt)
  {
    ( , , , updatedAt, ) = latestRoundData();
  }

  /**
   * @notice get past rounds answers
   * @param _roundId the answer number to retrieve the answer for
   * @dev deprecated. Use getRoundData instead.
   */
  function getAnswer(uint256 _roundId)
    public
    view
    virtual
    override
    returns (int256 answer)
  {
    ( , answer, , , ) = getRoundData(_roundId);
  }

  /**
   * @notice get block timestamp when an answer was last updated
   * @param _roundId the answer number to retrieve the updated timestamp for
   * @dev deprecated. Use getRoundData instead.
   */
  function getTimestamp(uint256 _roundId)
    public
    view
    virtual
    override
    returns (uint256 updatedAt)
  {
    ( , , , updatedAt, ) = getRoundData(_roundId);
  }

  /**
   * @notice get the latest completed round where the answer was updated
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestRound()
    public
    view
    virtual
    override
    returns (uint256 roundId)
  {
    ( roundId, , , , ) = latestRoundData();
  }

  /**
   * @notice get data about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * Note that different underlying implementations of AggregatorV3Interface
   * have slightly different semantics for some of the return values. Consumers
   * should determine what implementations they expect to receive
   * data from and validate that they can properly handle return data from all
   * of them.
   * @param _requestId the round ID to retrieve the round data for
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started.
   * (Only some AggregatorV3Interface implementations return meaningful values)
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed.
   * (Only some AggregatorV3Interface implementations return meaningful values)
   * @dev Note that answer and updatedAt may change between queries.
   */
  function getRoundData(uint256 _requestId)
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
    (reqEpoch, reqRound) = parseRequestId(_requestId);
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = epochAggregators[reqEpoch].getRoundData(reqRound);

    return (
      addEpoch(reqEpoch, roundId),
      answer,
      startedAt,
      updatedAt,
      addEpoch(reqEpoch, answeredInRound)
    );
  }

  /**
   * @notice get data about the latest round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * Note that different underlying implementations of AggregatorV3Interface
   * have slightly different semantics for some of the return values. Consumers
   * should determine what implementations they expect to receive
   * data from and validate that they can properly handle return data from all
   * of them.
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started.
   * (Only some AggregatorV3Interface implementations return meaningful values)
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed.
   * (Only some AggregatorV3Interface implementations return meaningful values)
   * @dev Note that answer and updatedAt may change between queries.
   */
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
    uint16 currentEpoch = epoch; // cache storage reads
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = aggregator.latestRoundData();

    return (
      addEpoch(currentEpoch, roundId),
      answer,
      startedAt,
      updatedAt,
      addEpoch(currentEpoch, answeredInRound)
    );
  }

  /**
   * @notice Used if an aggregator contract has been proposed.
   * @param _roundId the round ID to retrieve the round data for
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started.
   * (Only some AggregatorV3Interface implementations return meaningful values)
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed.
  */
  function proposedGetRoundData(uint256 _roundId)
    public
    view
    virtual
    hasProposal()
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return proposedAggregator.getRoundData(_roundId);
  }

  /**
   * @notice Used if an aggregator contract has been proposed.
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started.
   * (Only some AggregatorV3Interface implementations return meaningful values)
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed.
  */
  function proposedLatestRoundData()
    public
    view
    virtual
    hasProposal()
    returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    )
  {
    return proposedAggregator.latestRoundData();
  }

  /**
   * @notice represents the number of decimals the aggregator responses represent.
   */
  function decimals()
    external
    view
    override
    returns (uint8)
  {
    return aggregator.decimals();
  }

  /**
   * @notice the version number representing the type of aggregator the proxy
   * points to.
   */
  function version()
    external
    view
    override
    returns (uint256)
  {
    return aggregator.version();
  }

  /**
   * @notice returns the description of the aggregator the proxy points to.
   */
  function description()
    external
    view
    override
    returns (string memory)
  {
    return aggregator.description();
  }

  /**
   * @notice Allows the owner to propose a new address for the aggregator
   * @param _aggregator The new address for the aggregator contract
   */
  function proposeAggregator(address _aggregator)
    external
    onlyOwner()
  {
    proposedAggregator = AggregatorV3Interface(_aggregator);
  }

  /**
   * @notice Allows the owner to confirm and change the address
   * to the proposed aggregator
   * @dev Reverts if the given address doesn't match what was previously
   * proposed
   * @param _aggregator The new address for the aggregator contract
   */
  function confirmAggregator(address _aggregator)
    external
    onlyOwner()
  {
    require(_aggregator == address(proposedAggregator), "Invalid proposed aggregator");
    delete proposedAggregator;
    setAggregator(_aggregator);
  }


  /*
   * Internal
   */

  function setAggregator(address _aggregator)
    internal
  {
    epoch++;
    epochAggregators[epoch] = AggregatorV3Interface(_aggregator);
    aggregator = AggregatorV3Interface(_aggregator);
  }

  // PRIVATE

  function addEpoch(
    uint256 _epoch,
    uint256 _originalId
  )
    private
    view
    returns (uint256)
  {
    return _epoch.mul(EPOCH_BASE).add(_originalId);
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

    uint256 requestIdMask = (2**EPOCH_OFFSET).sub(1);
    uint256 roundId = requestId & requestIdMask;

    return (epochId, roundId);
  }

  // MODIFIERS

  modifier hasProposal() {
    require(address(proposedAggregator) != address(0), "No proposed aggregator present");
    _;
  }

}
