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

  struct Phase {
    uint16 id;
    address aggregator;
  }
  Phase private currentPhase;
  address public proposedAggregator;
  mapping(uint16 => address) public phaseAggregators;

  uint256 constant private PHASE_OFFSET = 64;
  uint256 constant private PHASE_BASE = 2 ** PHASE_OFFSET;
  uint256 constant private PHASE_MASK = 0xFFFF << PHASE_OFFSET;
  uint256 constant private REQUEST_ID_MASK = ~PHASE_MASK;

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
    return AggregatorInterface(currentPhase.aggregator).latestAnswer();
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
    return AggregatorInterface(currentPhase.aggregator).latestTimestamp();
  }

  /**
   * @notice get past rounds answers
   * @param _requestId the answer number to retrieve the answer for
   * @dev deprecated. Use getRoundData instead.
   */
  function getAnswer(uint256 _requestId)
    public
    view
    virtual
    override
    returns (int256 answer)
  {
    (
      uint256 requestRoundId, ,
      address aggregator
    ) = getRoundIdPhaseIdAndAggregator(_requestId);
    return AggregatorInterface(aggregator).getAnswer(requestRoundId);
  }

  /**
   * @notice get block timestamp when an answer was last updated
   * @param _requestId the answer number to retrieve the updated timestamp for
   * @dev deprecated. Use getRoundData instead.
   */
  function getTimestamp(uint256 _requestId)
    public
    view
    virtual
    override
    returns (uint256 updatedAt)
  {
    (
      uint256 requestRoundId, ,
      address aggregator
    ) = getRoundIdPhaseIdAndAggregator(_requestId);
    return AggregatorInterface(aggregator).getTimestamp(requestRoundId);
  }

  /**
   * @notice get the latest completed round where the answer was updated. This
   * ID includes the proxy's phase, to make sure round IDs increase even when
   * switching to a newly deployed aggregator.
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestRound()
    public
    view
    virtual
    override
    returns (uint256 roundId)
  {
    Phase memory current = currentPhase; // cache storage reads
    uint256 roundId = AggregatorInterface(current.aggregator).latestRound();
    return addPhase(current.id, roundId);
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
   * @return roundId is the round ID from the aggregator for which the data was
   * retrieved combined with an phase to ensure that round IDs get larger as
   * time moves forward.
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
    (
      uint256 requestRoundId,
      uint16 requestPhaseId,
      address aggregator
    ) = getRoundIdPhaseIdAndAggregator(_requestId);
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = AggregatorV3Interface(aggregator).getRoundData(requestRoundId);
    roundId = addPhase(requestPhaseId, roundId);
    answeredInRound = addPhase(requestPhaseId, answeredInRound);
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
   * @return roundId is the round ID from the aggregator for which the data was
   * retrieved combined with an phase to ensure that round IDs get larger as
   * time moves forward.
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
    Phase memory current = currentPhase; // cache storage reads
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = AggregatorV3Interface(current.aggregator).latestRoundData();
    roundId = addPhase(current.id, roundId);
    answeredInRound = addPhase(current.id, answeredInRound);
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
    return AggregatorV3Interface(proposedAggregator).getRoundData(_roundId);
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
    return AggregatorV3Interface(proposedAggregator).latestRoundData();
  }

  /**
   * @notice returns the current phase's aggregator address.
   */
  function aggregator()
    external
    view
    returns (address)
  {
    return currentPhase.aggregator;
  }

  /**
   * @notice returns the current phase's ID.
   */
  function phaseId()
    external
    view
    returns (uint16)
  {
    return currentPhase.id;
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
    return AggregatorV3Interface(currentPhase.aggregator).decimals();
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
    return AggregatorV3Interface(currentPhase.aggregator).version();
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
    return AggregatorV3Interface(currentPhase.aggregator).description();
  }

  /**
   * @notice Allows the owner to propose a new address for the aggregator
   * @param _aggregator The new address for the aggregator contract
   */
  function proposeAggregator(address _aggregator)
    external
    onlyOwner()
  {
    proposedAggregator = _aggregator;
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
    currentPhase.id++;
    phaseAggregators[currentPhase.id] = _aggregator;
    currentPhase.aggregator = _aggregator;
  }

  function addPhase(
    uint256 _phase,
    uint256 _originalId
  )
    internal
    view
    returns (uint256)
  {
    return (_originalId & REQUEST_ID_MASK) | _phase.mul(PHASE_BASE);
  }

  function getRoundIdPhaseIdAndAggregator(
    uint256 _requestId
  )
    internal
    view
    returns (uint256, uint16, address)
  {
    uint16 requestPhaseId = uint16((PHASE_MASK & _requestId) >> PHASE_OFFSET);
    uint256 requestRoundId = _requestId & REQUEST_ID_MASK;

    return (
      requestRoundId,
      requestPhaseId,
      phaseAggregators[requestPhaseId]
    );
  }


  /*
   * Modifiers
   */

  modifier hasProposal() {
    require(address(proposedAggregator) != address(0), "No proposed aggregator present");
    _;
  }

}
