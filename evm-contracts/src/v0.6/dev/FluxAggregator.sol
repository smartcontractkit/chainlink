pragma solidity 0.6.2;

import "../Median.sol";
import "../vendor/SafeMath.sol";
import "./SafeMath128.sol";
import "./SafeMath64.sol";
import "./SafeMath32.sol";
import "../interfaces/LinkTokenInterface.sol";
import "./AggregatorInterface.sol";
import "../Owned.sol";

/**
 * @title The Prepaid Aggregator contract
 * @notice Node handles aggregating data pushed in from off-chain, and unlocks
 * payment for oracles as they report. Oracles' submissions are gathered in
 * rounds, with each round aggregating the submissions for each oracle into a
 * single answer. The latest aggregated answer is exposed as well as historical
 * answers and their updated at timestamp.
 */
contract FluxAggregator is AggregatorInterface, Owned {
  using SafeMath for uint256;
  using SafeMath128 for uint128;
  using SafeMath64 for uint64;
  using SafeMath32 for uint32;

  struct Round {
    int256 answer;
    uint64 startedAt;
    uint64 updatedAt;
    uint32 answeredInRound;
    RoundDetails details;
  }

  struct RoundDetails {
    int256[] answers;
    uint32 maxAnswers;
    uint32 minAnswers;
    uint32 timeout;
    uint128 paymentAmount;
  }

  struct OracleStatus {
    uint128 withdrawable;
    uint32 startingRound;
    uint32 endingRound;
    uint32 lastReportedRound;
    uint32 lastStartedRound;
    int256 latestAnswer;
    uint16 index;
    address admin;
  }

  uint256 constant public VERSION = 2;

  uint128 public allocatedFunds;
  uint128 public availableFunds;

  // Round related params
  uint128 public paymentAmount;
  uint32 public maxAnswerCount;
  uint32 public minAnswerCount;
  uint32 public restartDelay;
  uint32 public timeout;
  uint8 public decimals;
  bytes32 public description;

  uint32 private reportingRoundId;
  uint32 internal latestRoundId;
  LinkTokenInterface private LINK;
  mapping(address => OracleStatus) private oracles;
  mapping(uint32 => Round) internal rounds;
  mapping(address => bool) internal authorizedRequesters;
  address[] private oracleAddresses;

  event AvailableFundsUpdated(uint256 indexed amount);
  event RoundDetailsUpdated(
    uint128 indexed paymentAmount,
    uint32 indexed minAnswerCount,
    uint32 indexed maxAnswerCount,
    uint32 restartDelay,
    uint32 timeout // measured in seconds
  );
  event OracleAdded(address indexed oracle);
  event OracleRemoved(address indexed oracle);
  event OracleAdminUpdated(address indexed oracle, address indexed newAdmin);
  event SubmissionReceived(
    int256 indexed answer,
    uint32 indexed round,
    address indexed oracle
  );
  event RequesterAuthorizationSet(address indexed requester, bool allowed);

  uint32 constant private ROUND_MAX = 2**32-1;

  /**
   * @notice Deploy with the address of the LINK token and initial payment amount
   * @dev Sets the LinkToken address and amount of LINK paid
   * @param _link The address of the LINK token
   * @param _paymentAmount The amount paid of LINK paid to each oracle per response
   * @param _timeout is the number of seconds after the previous round that are
   * allowed to lapse before allowing an oracle to skip an unfinished round
   */
  constructor(
    address _link,
    uint128 _paymentAmount,
    uint32 _timeout,
    uint8 _decimals,
    bytes32 _description
  ) public {
    LINK = LinkTokenInterface(_link);
    paymentAmount = _paymentAmount;
    timeout = _timeout;
    decimals = _decimals;
    description = _description;
    rounds[0].updatedAt = uint64(block.timestamp.sub(uint256(_timeout)));
  }

  /**
   * @notice called by oracles when they have witnessed a need to update
   * @param _round is the ID of the round this answer pertains to
   * @param _answer is the updated data that the oracle is submitting
   */
  function updateAnswer(uint256 _round, int256 _answer)
    external
    onlyValidRoundId(uint32(_round))
    onlyValidOracleRound(uint32(_round))
  {
    initializeNewRound(uint32(_round));
    recordSubmission(_answer, uint32(_round));
    updateRoundAnswer(uint32(_round));
    payOracle(uint32(_round));
    deleteRoundDetails(uint32(_round));
  }

  /**
   * @notice called by the owner to add a new Oracle and update the round
   * related parameters
   * @param _oracle is the address of the new Oracle being added
   * @param _admin is the admin address of the new oracle. Only this address
   * is allowed to access the oracle's funds.
   * @param _minAnswers is the new minimum answer count for each round
   * @param _maxAnswers is the new maximum answer count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function addOracle(
    address _oracle,
    address _admin,
    uint32 _minAnswers,
    uint32 _maxAnswers,
    uint32 _restartDelay
  )
    external
    onlyOwner()
    onlyUnenabledAddress(_oracle)
  {
    require(oracleCount() < 42, "cannot add more than 42 oracles");
    require(_admin != address(0), "admin address must not be 0x0");
    require(
      oracles[_oracle].admin == address(0) || oracles[_oracle].admin == _admin,
      "cannot modify previously-set admin address"
    );
    oracles[_oracle].startingRound = getStartingRound(_oracle);
    oracles[_oracle].endingRound = ROUND_MAX;
    oracleAddresses.push(_oracle);
    oracles[_oracle].index = uint16(oracleAddresses.length.sub(1));
    oracles[_oracle].admin = _admin;

    emit OracleAdded(_oracle);
    emit OracleAdminUpdated(_oracle, _admin);

    updateFutureRounds(paymentAmount, _minAnswers, _maxAnswers, _restartDelay, timeout);
  }

  /**
   * @notice called by the owner to remove an Oracle and update the round
   * related parameters
   * @param _oracle is the address of the Oracle being removed
   * @param _minAnswers is the new minimum answer count for each round
   * @param _maxAnswers is the new maximum answer count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function removeOracle(
    address _oracle,
    uint32 _minAnswers,
    uint32 _maxAnswers,
    uint32 _restartDelay
  )
    external
    onlyOwner()
    onlyEnabledAddress(_oracle)
  {
    oracles[_oracle].endingRound = reportingRoundId;
    address tail = oracleAddresses[oracleCount().sub(1)];
    uint16 index = oracles[_oracle].index;
    oracles[tail].index = index;
    delete oracles[_oracle].index;
    oracleAddresses[index] = tail;
    oracleAddresses.pop();

    emit OracleRemoved(_oracle);

    updateFutureRounds(paymentAmount, _minAnswers, _maxAnswers, _restartDelay, timeout);
  }

  /**
   * @notice update the round and payment related parameters for subsequent
   * rounds
   * @param _newPaymentAmount is the payment amount for subsequent rounds
   * @param _minAnswers is the new minimum answer count for each round
   * @param _maxAnswers is the new maximum answer count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function updateFutureRounds(
    uint128 _newPaymentAmount,
    uint32 _minAnswers,
    uint32 _maxAnswers,
    uint32 _restartDelay,
    uint32 _timeout
  )
    public
    onlyOwner()
    onlyValidRange(_minAnswers, _maxAnswers, _restartDelay)
  {
    paymentAmount = _newPaymentAmount;
    minAnswerCount = _minAnswers;
    maxAnswerCount = _maxAnswers;
    restartDelay = _restartDelay;
    timeout = _timeout;

    emit RoundDetailsUpdated(
      paymentAmount,
      _minAnswers,
      _maxAnswers,
      _restartDelay,
      _timeout
    );
  }

  /**
   * @notice recalculate the amount of LINK available for payouts
   */
  function updateAvailableFunds()
    public
  {
    uint128 pastAvailableFunds = availableFunds;

    uint256 available = LINK.balanceOf(address(this)).sub(allocatedFunds);
    availableFunds = uint128(available);

    if (pastAvailableFunds != available) {
      emit AvailableFundsUpdated(available);
    }
  }

  /**
   * @notice returns the number of oracles
   */
  function oracleCount() public view returns (uint32) {
    return uint32(oracleAddresses.length);
  }

  /**
   * @notice returns an array of addresses containing the oracles on contract
   */
  function getOracles() external view returns (address[] memory) {
    return oracleAddresses;
  }

  /**
   * @notice get the most recently reported answer
   */
  function latestAnswer()
    external
    view
    virtual
    override
    returns (int256)
  {
    return _latestAnswer();
  }

  /**
   * @notice get the most recent updated at timestamp
   */
  function latestTimestamp()
    external
    view
    virtual
    override
    returns (uint256)
  {
    return _latestTimestamp();
  }

  /**
   * @notice get the ID of the last updated round
   */
  function latestRound()
    external
    view
    override
    returns (uint256)
  {
    return latestRoundId;
  }

  /**
   * @notice get the ID of the round most recently reported on
   */
  function reportingRound()
    external
    view
    returns (uint256)
  {
    return reportingRoundId;
  }

  /**
   * @notice get past rounds answers
   * @param _roundId the round number to retrieve the answer for
   */
  function getAnswer(uint256 _roundId)
    external
    view
    virtual
    override
    returns (int256)
  {
    return _getAnswer(_roundId);
  }

  /**
   * @notice get timestamp when an answer was last updated
   * @param _roundId the round number to retrieve the updated timestamp for
   */
  function getTimestamp(uint256 _roundId)
    external
    view
    virtual
    override
    returns (uint256)
  {
    return _getTimestamp(_roundId);
  }

  /**
   * @notice get the timed out status of a given round
   * @param _roundId the round number to retrieve the timed out status for
   */
  function getTimedOutStatus(uint256 _roundId)
    external
    view
    returns (bool)
  {
    uint32 roundId = uint32(_roundId);
    uint32 answeredIn = rounds[roundId].answeredInRound;
    return answeredIn > 0 && answeredIn != roundId;
  }

  /**
   * @notice get the start time of the current reporting round
   */
  function reportingRoundStartedAt()
    external
    view
    returns (uint256)
  {
    return rounds[reportingRoundId].startedAt;
  }

  /**
   * @notice get the start time of a round
   * @param _roundId the round number to retrieve the startedAt time for
   */
  function getRoundStartedAt(uint256 _roundId)
    external
    view
    returns (uint256)
  {
    return rounds[uint32(_roundId)].startedAt;
  }

  /**
   * @notice get the round ID that an answer was originally reported in
   * @param _roundId the round number to retrieve the answer for
   */
  function getOriginatingRoundOfAnswer(uint256 _roundId)
    external
    view
    returns (uint256)
  {
    return rounds[uint32(_roundId)].answeredInRound;
  }

  /**
   * @notice query the available amount of LINK for an oracle to withdraw
   */
  function withdrawable(address _oracle)
    external
    view
    returns (uint256)
  {
    return oracles[_oracle].withdrawable;
  }

  /**
   * @notice transfers the oracle's LINK to another address. Can only be called
   * by the oracle's admin.
   * @param _oracle is the oracle whose LINK is transferred
   * @param _recipient is the address to send the LINK to
   * @param _amount is the amount of LINK to send
   */
  function withdraw(address _oracle, address _recipient, uint256 _amount)
    external
  {
    require(oracles[_oracle].admin == msg.sender, "Only admin can withdraw");

    uint128 amount = uint128(_amount);
    uint128 available = oracles[_oracle].withdrawable;
    require(available >= amount, "Insufficient balance");

    oracles[_oracle].withdrawable = available.sub(amount);
    allocatedFunds = allocatedFunds.sub(amount);

    assert(LINK.transfer(_recipient, uint256(amount)));
  }

  /**
   * @notice transfers the owner's LINK to another address
   * @param _recipient is the address to send the LINK to
   * @param _amount is the amount of LINK to send
   */
  function withdrawFunds(address _recipient, uint256 _amount)
    external
    onlyOwner()
  {
    require(availableFunds >= _amount, "Insufficient funds");
    require(LINK.transfer(_recipient, _amount), "LINK transfer failed");
    updateAvailableFunds();
  }

  /**
   * @notice get the latest submission for any oracle
   * @param _oracle is the address to lookup the latest submission for
   */
  function latestSubmission(address _oracle)
    external
    view
    returns (int256, uint256)
  {
    return (oracles[_oracle].latestAnswer, oracles[_oracle].lastReportedRound);
  }

  /**
   * @notice get the admin address of an oracle
   * @param _oracle is the address of the oracle whose admin is being queried
   */
  function getAdmin(address _oracle)
    external
    view
    returns (address)
  {
    return oracles[_oracle].admin;
  }

  /**
   * @notice update the admin address for an oracle
   * @param _oracle is the address of the oracle whose admin is being updated
   * @param _newAdmin is the new admin address
   */
  function updateAdmin(address _oracle, address _newAdmin)
    external
  {
    require(oracles[_oracle].admin == msg.sender, "Only admin can update admin");
    oracles[_oracle].admin = _newAdmin;

    emit OracleAdminUpdated(_oracle, _newAdmin);
  }

  /**
   * @notice allows non-oracles to request a new round
   */
  function startNewRound()
    external
    onlyAuthorizedRequesters()
  {
    uint32 current = reportingRoundId;

    require(rounds[current].updatedAt > 0 || timedOut(current), 'Cannot start a round mid-round');

    initializeNewRound(current.add(1));
  }

  /**
   * @notice allows the owner to specify new non-oracles to start new rounds
   * @param _requester is the address to set permissions for
   * @param _allowed is a boolean specifying whether they can start new rounds or not
   */
  function setAuthorization(address _requester, bool _allowed)
    external
    onlyOwner()
  {
    authorizedRequesters[_requester] = _allowed;

    emit RequesterAuthorizationSet(_requester, _allowed);
  }

  /**
   * @notice called through LINK's transferAndCall to update available funds
   * in the same transaction as the funds were transfered to the aggregator
   */
  function onTokenTransfer(address, uint256, bytes memory) public {
    updateAvailableFunds();
  }

  /**
   * Internal
   */

  /**
   * @dev Internal implementation for latestAnswer
   */
  function _latestAnswer()
    internal
    view
    returns (int256)
  {
    return rounds[latestRoundId].answer;
  }

  /**
   * @dev Internal implementation of latestTimestamp
   */
  function _latestTimestamp()
    internal
    view
    returns (uint256)
  {
    return rounds[latestRoundId].updatedAt;
  }

  /**
   * @dev Internal implementation of getAnswer
   */
  function _getAnswer(uint256 _roundId)
    internal
    view
    returns (int256)
  {
    return rounds[uint32(_roundId)].answer;
  }

  /**
   * @dev Internal implementation of getTimestamp
   */
  function _getTimestamp(uint256 _roundId)
    internal
    view
    returns (uint256)
  {
    return rounds[uint32(_roundId)].updatedAt;
  }

  /**
   * Private
   */

  function initializeNewRound(uint32 _id)
    private
    ifNewRound(_id)
    ifDelayed(_id)
  {
    updateTimedOutRoundInfo(_id.sub(1));

    reportingRoundId = _id;
    rounds[_id].details.maxAnswers = maxAnswerCount;
    rounds[_id].details.minAnswers = minAnswerCount;
    rounds[_id].details.paymentAmount = paymentAmount;
    rounds[_id].details.timeout = timeout;
    rounds[_id].startedAt = uint64(block.timestamp);

    oracles[msg.sender].lastStartedRound = _id;

    emit NewRound(_id, msg.sender, rounds[_id].startedAt);
  }

  function updateTimedOutRoundInfo(uint32 _id)
    private
    ifTimedOut(_id)
    onlyWithPreviousAnswer(_id)
  {
    uint32 prevId = _id.sub(1);
    rounds[_id].answer = rounds[prevId].answer;
    rounds[_id].answeredInRound = rounds[prevId].answeredInRound;
    rounds[_id].updatedAt = uint64(block.timestamp);

    delete rounds[_id].details;
  }

  function updateRoundAnswer(uint32 _id)
    private
    ifMinAnswersReceived(_id)
  {
    int256 newAnswer = Median.calculateInplace(rounds[_id].details.answers);
    rounds[_id].answer = newAnswer;
    rounds[_id].updatedAt = uint64(block.timestamp);
    rounds[_id].answeredInRound = _id;
    latestRoundId = _id;

    emit AnswerUpdated(newAnswer, _id, now);
  }

  function payOracle(uint32 _id)
    private
  {
    uint128 payment = rounds[_id].details.paymentAmount;
    uint128 available = availableFunds.sub(payment);

    availableFunds = available;
    allocatedFunds = allocatedFunds.add(payment);
    oracles[msg.sender].withdrawable = oracles[msg.sender].withdrawable.add(payment);

    emit AvailableFundsUpdated(available);
  }

  function recordSubmission(int256 _answer, uint32 _id)
    private
    onlyWhenAcceptingAnswers(_id)
  {
    rounds[_id].details.answers.push(_answer);
    oracles[msg.sender].lastReportedRound = _id;
    oracles[msg.sender].latestAnswer = _answer;

    emit SubmissionReceived(_answer, _id, msg.sender);
  }

  function deleteRoundDetails(uint32 _id)
    private
    ifMaxAnswersReceived(_id)
  {
    delete rounds[_id].details;
  }

  function timedOut(uint32 _id)
    private
    view
    returns (bool)
  {
    uint64 startedAt = rounds[_id].startedAt;
    uint32 roundTimeout = rounds[_id].details.timeout;
    return startedAt > 0 && roundTimeout > 0 && startedAt.add(roundTimeout) < block.timestamp;
  }

  function finished(uint32 _id)
    private
    view
    returns (bool)
  {
    return rounds[_id].updatedAt > 0;
  }

  function getStartingRound(address _oracle)
    private
    view
    returns (uint32)
  {
    uint32 currentRound = reportingRoundId;
    if (currentRound != 0 && currentRound == oracles[_oracle].endingRound) {
      return currentRound;
    }
    return currentRound.add(1);
  }

  function roundState()
    external
    view
    returns (uint32 _reportableRoundId, bool _eligibleToSubmit, int256 _latestAnswer)
  {
    if (finished(reportingRoundId) || timedOut(reportingRoundId)) {
      return (reportingRoundId.add(1), eligibleToSubmit(), rounds[latestRoundId].answer);
    }
    return (reportingRoundId, eligibleToSubmit(), rounds[latestRoundId].answer);
  }

  function eligibleToSubmit()
    private
    view
    returns (bool)
  {
    bool finishedOrTimedOut = finished(reportingRoundId) || timedOut(reportingRoundId);
    uint32 startingRound = oracles[msg.sender].startingRound;
    if (startingRound == 0) {
      return false;
    }
    uint32 reportableRoundId = finishedOrTimedOut ? reportingRoundId.add(1) : reportingRoundId;
    if (startingRound > reportableRoundId) {
      return false;
    } else if (oracles[msg.sender].endingRound < reportableRoundId) {
      return false;
    } else if (oracles[msg.sender].lastReportedRound >= reportableRoundId) {
      return false;
    } else if (rounds[reportableRoundId].details.maxAnswers == 0) {
      return false;
    } else if (finishedOrTimedOut) {
      uint32 lastStartedRound = oracles[msg.sender].lastStartedRound;
      if (reportableRoundId <= lastStartedRound + restartDelay && lastStartedRound > 0) {
        return false;
      }
    }

    return true;
  }

  /**
   * Modifiers
   */

  modifier onlyValidOracleRound(uint32 _id) {
    uint32 startingRound = oracles[msg.sender].startingRound;
    require(startingRound != 0, "Only updatable by whitelisted oracles");
    require(startingRound <= _id, "New oracles cannot participate in in-progress rounds");
    require(oracles[msg.sender].endingRound >= _id, "Oracle has been removed from whitelist");
    require(oracles[msg.sender].lastReportedRound < _id, "Cannot update round reports");
    _;
  }

  modifier ifMinAnswersReceived(uint32 _id) {
    if (rounds[_id].details.answers.length >= rounds[_id].details.minAnswers) {
      _;
    }
  }

  modifier ifMaxAnswersReceived(uint32 _id) {
    if (rounds[_id].details.answers.length == rounds[_id].details.maxAnswers) {
      _;
    }
  }

  modifier onlyWhenAcceptingAnswers(uint32 _id) {
    require(rounds[_id].details.maxAnswers != 0, "Round not currently eligible for reporting");
    _;
  }

  modifier ifNewRound(uint32 _id) {
    if (_id == reportingRoundId.add(1)) {
      _;
    }
  }

  modifier ifDelayed(uint32 _id) {
    uint256 lastStarted = oracles[msg.sender].lastStartedRound;
    if (_id > lastStarted + restartDelay || lastStarted == 0) {
      _;
    }
  }

  modifier onlyValidRoundId(uint32 _id) {
    require(_id == reportingRoundId || _id == reportingRoundId.add(1), "Must report on current round");
    require(_id == 1 || finished(_id.sub(1)) || timedOut(_id.sub(1)), "Not eligible to bump round");
    _;
  }

  modifier onlyValidRange(uint32 _min, uint32 _max, uint32 _restartDelay) {
    uint32 oracleNum = oracleCount(); // Save on storage reads
    require(oracleNum >= _max, "Cannot have the answer max higher oracle count");
    require(_max >= _min, "Cannot have the answer minimum higher the max");
    require(oracleNum == 0 || oracleNum > _restartDelay, "Restart delay must be less than oracle count");
    _;
  }

  modifier onlyUnenabledAddress(address _oracle) {
    require(oracles[_oracle].endingRound != ROUND_MAX, "Address is already recorded as an oracle");
    _;
  }

  modifier onlyEnabledAddress(address _oracle) {
    require(oracles[_oracle].endingRound == ROUND_MAX, "Address is not a whitelisted oracle");
    _;
  }

  modifier ifTimedOut(uint32 _id) {
    if (timedOut(_id)) {
      _;
    }
  }

  modifier onlyWithPreviousAnswer(uint32 _id) {
    require(rounds[_id.sub(1)].updatedAt != 0, "Must have a previous answer to pull from");
    _;
  }

  modifier onlyAuthorizedRequesters() {
    require(authorizedRequesters[msg.sender], "Only authorized requesters can call");
    _;
  }

}
