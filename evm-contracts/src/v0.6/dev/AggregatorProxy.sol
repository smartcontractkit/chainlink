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
    AggregatorV3Interface aggregator;
  }
  Phase private currentPhase;
  AggregatorV3Interface public proposedAggregator;
  mapping(uint16 => AggregatorV3Interface) public phaseAggregators;

  uint256 constant private PHASE_OFFSET = 64;
  uint256 constant private PHASE_BASE = 2 ** PHASE_OFFSET;
  uint256 constant private PHASE_MASK = 0xFFFF << PHASE_OFFSET;
  uint256 constant private REQUEST_ID_MASK = ~PHASE_MASK;
  bytes32 constant private EXPECTED_V3_ERROR = keccak256(bytes("No data present"));

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
    ( , answer, , , ) = tryLatestRoundData();
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
    ( , , , updatedAt, ) = tryLatestRoundData();
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
    ( , answer, , , ) = tryGetRoundData(_roundId);
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
    ( , , , updatedAt, ) = tryGetRoundData(_roundId);
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
    ( roundId, , , , ) = tryLatestRoundData();
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
    uint16 requestPhaseId;
    uint256 requestRoundId;
    (requestPhaseId, requestRoundId) = parseRequestId(_requestId);
    (
      roundId,
      answer,
      startedAt,
      updatedAt,
      answeredInRound
    ) = phaseAggregators[requestPhaseId].getRoundData(requestRoundId);
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
    ) = current.aggregator.latestRoundData();
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
   * @notice returns the current phase's aggregator address.
   */
  function aggregator()
    external
    view
    returns (address)
  {
    return address(currentPhase.aggregator);
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
    return currentPhase.aggregator.decimals();
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
    return currentPhase.aggregator.version();
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
    return currentPhase.aggregator.description();
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
    currentPhase.id++;
    phaseAggregators[currentPhase.id] = AggregatorV3Interface(_aggregator);
    currentPhase.aggregator = AggregatorV3Interface(_aggregator);
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

  function parseRequestId(
    uint256 _requestId
  )
    internal
    view
    returns (uint16, uint256)
  {
    uint16 phaseId = uint16((PHASE_MASK & _requestId) >> PHASE_OFFSET);
    uint256 roundId = _requestId & REQUEST_ID_MASK;

    return (phaseId, roundId);
  }
  function tryLatestRoundData()
    internal
    view
    returns (uint256, int256, uint256, uint256, uint256)
  {
    try this.latestRoundData() returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    ) {
      return (roundId, answer, startedAt, updatedAt, answeredInRound);
    } catch Error(string memory reason) {
      return handleRoundDataError(reason);
    }
  }

  function tryGetRoundData(
    uint256 _requestId
  )
    internal
    view
    returns (uint256, int256, uint256, uint256, uint256)
  {
    try this.getRoundData(_requestId) returns (
      uint256 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint256 answeredInRound
    ) {
      return (roundId, answer, startedAt, updatedAt, answeredInRound);
    } catch Error(string memory reason) {
      return handleRoundDataError(reason);
    }
  }

  function handleRoundDataError(
    string memory reason
  )
    internal
    view
    returns (uint256, int256, uint256, uint256, uint256)
  {
    require(keccak256(bytes(reason)) == EXPECTED_V3_ERROR, reason);
    return (0, 0, 0, 0, 0);
  }


  /*
   * Modifiers
   */

  modifier hasProposal() {
    require(address(proposedAggregator) != address(0), "No proposed aggregator present");
    _;
  }

}
