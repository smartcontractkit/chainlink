pragma solidity ^0.6.0;

import "../Median.sol";
import "../Owned.sol";
import "../SafeMath128.sol";
import "../SafeMath32.sol";
import "../SafeMath64.sol";
import "../interfaces/AggregatorInterface.sol";
import "../interfaces/AggregatorV3Interface.sol";
import './AnswerValidatorInterface.sol';
import "../interfaces/LinkTokenInterface.sol";
import "../vendor/SafeMath.sol";

/**
 * @title The Prepaid Aggregator contract
 * @notice Node handles aggregating data pushed in from off-chain, and unlocks
 * payment for oracles as they report. Oracles' submissions are gathered in
 * rounds, with each round aggregating the submissions for each oracle into a
 * single answer. The latest aggregated answer is exposed as well as historical
 * answers and their updated at timestamp.
 */
contract FluxAggregator is AggregatorInterface, AggregatorV3Interface, Owned {
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
    int256[] submissions;
    uint32 maxSubmissions;
    uint32 minSubmissions;
    uint32 timeout;
    uint128 paymentAmount;
  }

  struct OracleStatus {
    uint128 withdrawable;
    uint32 startingRound;
    uint32 endingRound;
    uint32 lastReportedRound;
    uint32 lastStartedRound;
    int256 latestSubmission;
    uint16 index;
    address admin;
    address pendingAdmin;
  }

  struct Requester {
    bool authorized;
    uint32 delay;
    uint32 lastStartedRound;
  }


  LinkTokenInterface public linkToken;
  AnswerValidatorInterface private answerValidator;
  uint128 public allocatedFunds;
  uint128 public availableFunds;

  // Round related params
  uint128 public paymentAmount;
  uint32 public maxSubmissionCount;
  uint32 public minSubmissionCount;
  uint32 public restartDelay;
  uint32 public timeout;
  uint8 public override decimals;
  string public override description;

  int256 immutable public minSubmissionValue;
  int256 immutable public maxSubmissionValue;

  uint256 constant public override version = 3;

  /**
   * @notice To ensure owner isn't withdrawing required funds as oracles are
   * submitting updates, we enforce that the contract maintains a minimum
   * reserve of RESERVE_ROUNDS * oracleCount() LINK earmarked for payment to
   * oracles. (Of course, this doesn't prevent the contract from running out of
   * funds without the owner's intervention.)
   */
  uint256 constant private RESERVE_ROUNDS = 2;
  uint256 constant private MAX_ORACLE_COUNT = 77;
  uint32 constant private ROUND_MAX = 2**32-1;
  uint256 private constant VALIDATOR_GAS_LIMIT = 100000;
  // An error specific to the Aggregator V3 Interface, to prevent possible
  // confusion around accidentally reading unset values as reported values.
  string constant private V3_NO_DATA_ERROR = "No data present";

  uint32 private reportingRoundId;
  uint32 internal latestRoundId;
  mapping(address => OracleStatus) private oracles;
  mapping(uint32 => Round) internal rounds;
  mapping(address => Requester) internal requesters;
  address[] private oracleAddresses;

  event AvailableFundsUpdated(
    uint256 indexed amount
  );
  event RoundDetailsUpdated(
    uint128 indexed paymentAmount,
    uint32 indexed minSubmissionCount,
    uint32 indexed maxSubmissionCount,
    uint32 restartDelay,
    uint32 timeout // measured in seconds
  );
  event OraclePermissionsUpdated(
    address indexed oracle,
    bool indexed whitelisted
  );
  event OracleAdminUpdated(
    address indexed oracle,
    address indexed newAdmin
  );
  event OracleAdminUpdateRequested(
    address indexed oracle,
    address admin,
    address newAdmin
  );
  event SubmissionReceived(
    int256 indexed submission,
    uint32 indexed round,
    address indexed oracle
  );
  event RequesterPermissionsSet(
    address indexed requester,
    bool authorized,
    uint32 delay
  );
  event AnswerValidatorUpdated(
    address indexed previous,
    address indexed current
  );

  /**
   * @notice set up the aggregator with initial configuration
   * @dev Sets the LinkToken address and amount of LINK paid
   * @param _link The address of the LINK token
   * @param _paymentAmount The amount paid of LINK paid to each oracle per submission, in wei (units of 10⁻¹⁸ LINK)
   * @param _timeout is the number of seconds after the previous round that are
   * allowed to lapse before allowing an oracle to skip an unfinished round
   * @param _answerValidator is an optional contract address for validating
   * external validation of answers
   * @param _minSubmissionValue is an immutable check for a lower bound of what
   * submission values are accepted from an oracle
   * @param _maxSubmissionValue is an immutable check for an upper bound of what
   * submission values are accepted from an oracle
   * @param _decimals represents the number of decimals to offset the answer by
   * @param _description a short description of what is being reported
   */
  constructor(
    address _link,
    uint128 _paymentAmount,
    uint32 _timeout,
    address _answerValidator,
    int256 _minSubmissionValue,
    int256 _maxSubmissionValue,
    uint8 _decimals,
    string memory _description
  ) public {
    linkToken = LinkTokenInterface(_link);
    paymentAmount = _paymentAmount;
    timeout = _timeout;
    answerValidator = AnswerValidatorInterface(_answerValidator);
    minSubmissionValue = _minSubmissionValue;
    maxSubmissionValue = _maxSubmissionValue;
    decimals = _decimals;
    description = _description;
    rounds[0].updatedAt = uint64(block.timestamp.sub(uint256(_timeout)));
  }

  /**
   * @notice called by oracles when they have witnessed a need to update
   * @param _roundId is the ID of the round this submission pertains to
   * @param _submission is the updated data that the oracle is submitting
   */
  function submit(uint256 _roundId, int256 _submission)
    external
  {
    bytes memory error = validateOracleRound(msg.sender, uint32(_roundId));
    require(_submission >= minSubmissionValue, "value below minSubmissionValue");
    require(_submission <= maxSubmissionValue, "value above maxSubmissionValue");
    require(error.length == 0, string(error));

    oracleInitializeNewRound(uint32(_roundId));
    recordSubmission(_submission, uint32(_roundId));
    updateRoundAnswer(uint32(_roundId));
    payOracle(uint32(_roundId));
    deleteRoundDetails(uint32(_roundId));
  }

  /**
   * @notice called by the owner to add new Oracles and update the round
   * related parameters
   * @param _oracles is the list of addresses of the new Oracles being added
   * @param _admins is the admin addresses of the new respective _oracles list.
   * Only this address is allowed to access the respective oracle's funds.
   * @param _minSubmissions is the new minimum submission count for each round
   * @param _maxSubmissions is the new maximum submission count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function addOracles(
    address[] calldata _oracles,
    address[] calldata _admins,
    uint32 _minSubmissions,
    uint32 _maxSubmissions,
    uint32 _restartDelay
  )
    external
    onlyOwner()
  {
    require(_oracles.length == _admins.length, "need same oracle and admin count");
    require(uint256(oracleCount()).add(_oracles.length) <= MAX_ORACLE_COUNT, "max oracles allowed");

    for (uint256 i = 0; i < _oracles.length; i++) {
      addOracle(_oracles[i], _admins[i]);
    }

    updateFutureRounds(paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, timeout);
  }

  /**
   * @notice called by the owner to remove Oracles and update the round
   * related parameters
   * @param _oracles is the address of the Oracles being removed
   * @param _minSubmissions is the new minimum submission count for each round
   * @param _maxSubmissions is the new maximum submission count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function removeOracles(
    address[] calldata _oracles,
    uint32 _minSubmissions,
    uint32 _maxSubmissions,
    uint32 _restartDelay
  )
    external
    onlyOwner()
  {
    for (uint256 i = 0; i < _oracles.length; i++) {
      removeOracle(_oracles[i]);
    }

    updateFutureRounds(paymentAmount, _minSubmissions, _maxSubmissions, _restartDelay, timeout);
  }

  /**
   * @notice update the round and payment related parameters for subsequent
   * rounds
   * @param _paymentAmount is the payment amount for subsequent rounds
   * @param _minSubmissions is the new minimum submission count for each round
   * @param _maxSubmissions is the new maximum submission count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function updateFutureRounds(
    uint128 _paymentAmount,
    uint32 _minSubmissions,
    uint32 _maxSubmissions,
    uint32 _restartDelay,
    uint32 _timeout
  )
    public
    onlyOwner()
  {
    uint32 oracleNum = oracleCount(); // Save on storage reads
    require(_maxSubmissions >= _minSubmissions, "max must equal/exceed min");
    require(oracleNum >= _maxSubmissions, "max cannot exceed total");
    require(oracleNum == 0 || oracleNum > _restartDelay, "delay cannot exceed total");
    require(availableFunds >= requiredReserve(_paymentAmount), "insufficient funds for payment");
    if (oracleCount() > 0) {
      require(_minSubmissions > 0, "min must be greater than 0");
    }

    paymentAmount = _paymentAmount;
    minSubmissionCount = _minSubmissions;
    maxSubmissionCount = _maxSubmissions;
    restartDelay = _restartDelay;
    timeout = _timeout;

    emit RoundDetailsUpdated(
      paymentAmount,
      _minSubmissions,
      _maxSubmissions,
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

    uint256 available = linkToken.balanceOf(address(this)).sub(allocatedFunds);
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
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestAnswer()
    public
    view
    virtual
    override
    returns (int256)
  {
    return rounds[latestRoundId].answer;
  }

  /**
   * @notice get the most recent updated at timestamp
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestTimestamp()
    public
    view
    virtual
    override
    returns (uint256)
  {
    return rounds[latestRoundId].updatedAt;
  }

  /**
   * @notice get the ID of the last updated round
   * @dev deprecated. Use latestRoundData instead.
   */
  function latestRound()
    public
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
   * @dev deprecated. Use getRoundData instead.
   */
  function getAnswer(uint256 _roundId)
    public
    view
    virtual
    override
    returns (int256)
  {
    return rounds[uint32(_roundId)].answer;
  }

  /**
   * @notice get timestamp when an answer was last updated
   * @param _roundId the round number to retrieve the updated timestamp for
   * @dev deprecated. Use getRoundData instead.
   */
  function getTimestamp(uint256 _roundId)
    public
    view
    virtual
    override
    returns (uint256)
  {
    return rounds[uint32(_roundId)].updatedAt;
  }

  /**
   * @notice get data about a round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values.
   * @param _roundId the round ID to retrieve the round data for
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started. This is 0
   * if the round hasn't been started yet.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed. answeredInRound may be smaller than roundId when the round
   * timed out. answerInRound is equal to roundId when the round didn't time out
   * and was completed regularly.
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet received
   * maxSubmissions) answer and updatedAt may change between queries.
   */
  function getRoundData(uint80 _roundId)
    public
    view
    virtual
    override
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    Round memory r = rounds[uint32(_roundId)];

    require(r.answeredInRound > 0, V3_NO_DATA_ERROR);

    return (
      _roundId,
      r.answer,
      r.startedAt,
      r.updatedAt,
      r.answeredInRound
    );
  }

  /**
   * @notice get data about the latest round. Consumers are encouraged to check
   * that they're receiving fresh data by inspecting the updatedAt and
   * answeredInRound return values. Consumers are encouraged to
   * use this more fully featured method over the "legacy" getAnswer/
   * latestAnswer/getTimestamp/latestTimestamp functions.
   * @return roundId is the round ID for which data was retrieved
   * @return answer is the answer for the given round
   * @return startedAt is the timestamp when the round was started. This is 0
   * if the round hasn't been started yet.
   * @return updatedAt is the timestamp when the round last was updated (i.e.
   * answer was last computed)
   * @return answeredInRound is the round ID of the round in which the answer
   * was computed. answeredInRound may be smaller than roundId when the round
   * timed out. answerInRound is equal to roundId when the round didn't time out
   * and was completed regularly.
   * @dev Note that for in-progress rounds (i.e. rounds that haven't yet received
   * maxSubmissions) answer and updatedAt may change between queries.
   */
   function latestRoundData()
    public
    view
    virtual
    override
    returns (
      uint80 roundId,
      int256 answer,
      uint256 startedAt,
      uint256 updatedAt,
      uint80 answeredInRound
    )
  {
    return getRoundData(latestRoundId);
  }


  /**
   * @notice query the available amount of LINK for an oracle to withdraw
   */
  function withdrawablePayment(address _oracle)
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
  function withdrawPayment(address _oracle, address _recipient, uint256 _amount)
    external
  {
    require(oracles[_oracle].admin == msg.sender, "only callable by admin");

    // Safe to downcast _amount because the total amount of LINK is less than 2^128.
    uint128 amount = uint128(_amount);
    uint128 available = oracles[_oracle].withdrawable;
    require(available >= amount, "insufficient withdrawable funds");

    oracles[_oracle].withdrawable = available.sub(amount);
    allocatedFunds = allocatedFunds.sub(amount);

    assert(linkToken.transfer(_recipient, uint256(amount)));
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
    require(uint256(availableFunds).sub(requiredReserve(paymentAmount)) >= _amount, "insufficient reserve funds");
    require(linkToken.transfer(_recipient, _amount), "token transfer failed");
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
    return (oracles[_oracle].latestSubmission, oracles[_oracle].lastReportedRound);
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
   * @notice transfer the admin address for an oracle
   * @param _oracle is the address of the oracle whose admin is being transfered
   * @param _newAdmin is the new admin address
   */
  function transferAdmin(address _oracle, address _newAdmin)
    external
  {
    require(oracles[_oracle].admin == msg.sender, "only callable by admin");
    oracles[_oracle].pendingAdmin = _newAdmin;

    emit OracleAdminUpdateRequested(_oracle, msg.sender, _newAdmin);
  }

  /**
   * @notice accept the admin address transfer for an oracle
   * @param _oracle is the address of the oracle whose admin is being transfered
   */
  function acceptAdmin(address _oracle)
    external
  {
    require(oracles[_oracle].pendingAdmin == msg.sender, "only callable by pending admin");
    oracles[_oracle].pendingAdmin = address(0);
    oracles[_oracle].admin = msg.sender;

    emit OracleAdminUpdated(_oracle, msg.sender);
  }

  /**
   * @notice allows non-oracles to request a new round
   */
  function requestNewRound()
    external
  {
    require(requesters[msg.sender].authorized, "not authorized requester");

    uint32 current = reportingRoundId;
    require(rounds[current].updatedAt > 0 || timedOut(current), "prev round must be supersedable");

    requesterInitializeNewRound(current.add(1));
  }

  /**
   * @notice allows the owner to specify new non-oracles to start new rounds
   * @param _requester is the address to set permissions for
   * @param _authorized is a boolean specifying whether they can start new rounds or not
   * @param _delay is the number of rounds the requester must wait before starting another round
   */
  function setRequesterPermissions(address _requester, bool _authorized, uint32 _delay)
    external
    onlyOwner()
  {
    if (requesters[_requester].authorized == _authorized) return;

    if (_authorized) {
      requesters[_requester].authorized = _authorized;
      requesters[_requester].delay = _delay;
    } else {
      delete requesters[_requester];
    }

    emit RequesterPermissionsSet(_requester, _authorized, _delay);
  }

  /**
   * @notice called through LINK's transferAndCall to update available funds
   * in the same transaction as the funds were transfered to the aggregator
   * @param _data is mostly ignored. It is checked for length, to be sure
   * nothing strange is passed in.
   */
  function onTokenTransfer(address, uint256, bytes calldata _data)
    external
  {
    require(_data.length == 0, "transfer doesn't accept calldata");
    updateAvailableFunds();
  }

  /**
   * @notice a method to provide all current info oracles need. Intended only
   * only to be callable by oracles. Not for use by contracts to read state.
   * @param _oracle the address to look up information for.
   */
  function oracleRoundState(address _oracle, uint32 _queriedRoundId)
    external
    view
    returns (
      bool _eligibleToSubmit,
      uint32 _roundId,
      int256 _latestSubmission,
      uint64 _startedAt,
      uint64 _timeout,
      uint128 _availableFunds,
      uint32 _oracleCount,
      uint128 _paymentAmount
    )
  {
    require(msg.sender == tx.origin, "off-chain reading only");

    if (_queriedRoundId > 0) {
      Round storage round = rounds[_queriedRoundId];
      return (
        eligibleForSpecificRound(_oracle, _queriedRoundId),
        _queriedRoundId,
        oracles[_oracle].latestSubmission,
        round.startedAt,
        round.details.timeout,
        availableFunds,
        oracleCount(),
        (round.startedAt > 0 ? round.details.paymentAmount : paymentAmount)
      );
    } else {
      return oracleRoundStateSuggestRound(_oracle);
    }
  }

  /**
   * @notice method to update the address which does external data validation.
   * @param _newValidator designates the address of the new validation contract.
   */
  function setAnswerValidator(address _newValidator)
    external
    onlyOwner()
  {
    address previous = address(answerValidator);

    if (previous != _newValidator) {
      answerValidator = AnswerValidatorInterface(_newValidator);

      emit AnswerValidatorUpdated(previous, _newValidator);
    }
  }


  /**
   * Private
   */

  function initializeNewRound(uint32 _roundId)
    private
  {
    updateTimedOutRoundInfo(_roundId.sub(1));

    reportingRoundId = _roundId;
    rounds[_roundId].details.maxSubmissions = maxSubmissionCount;
    rounds[_roundId].details.minSubmissions = minSubmissionCount;
    rounds[_roundId].details.paymentAmount = paymentAmount;
    rounds[_roundId].details.timeout = timeout;
    rounds[_roundId].startedAt = uint64(block.timestamp);

    emit NewRound(_roundId, msg.sender, rounds[_roundId].startedAt);
  }

  function oracleInitializeNewRound(uint32 _roundId)
    private
  {
    if (!newRound(_roundId)) return;
    uint256 lastStarted = oracles[msg.sender].lastStartedRound; // cache storage reads
    if (_roundId <= lastStarted + restartDelay && lastStarted != 0) return;

    initializeNewRound(_roundId);

    oracles[msg.sender].lastStartedRound = _roundId;
  }

  function requesterInitializeNewRound(uint32 _roundId)
    private
  {
    if (!newRound(_roundId)) return;
    uint256 lastStarted = requesters[msg.sender].lastStartedRound; // cache storage reads
    require(_roundId > lastStarted + requesters[msg.sender].delay || lastStarted == 0, "must delay requests");

    initializeNewRound(_roundId);

    requesters[msg.sender].lastStartedRound = _roundId;
  }

  function updateTimedOutRoundInfo(uint32 _roundId)
    private
  {
    if (!timedOut(_roundId)) return;

    uint32 prevId = _roundId.sub(1);
    rounds[_roundId].answer = rounds[prevId].answer;
    rounds[_roundId].answeredInRound = rounds[prevId].answeredInRound;
    rounds[_roundId].updatedAt = uint64(block.timestamp);

    delete rounds[_roundId].details;
  }

  function eligibleForSpecificRound(address _oracle, uint32 _queriedRoundId)
    private
    view
    returns (bool _eligible)
  {
    if (rounds[_queriedRoundId].startedAt > 0) {
      return acceptingSubmissions(_queriedRoundId) && validateOracleRound(_oracle, _queriedRoundId).length == 0;
    } else {
      return delayed(_oracle, _queriedRoundId) && validateOracleRound(_oracle, _queriedRoundId).length == 0;
    }
  }

  function oracleRoundStateSuggestRound(address _oracle)
    private
    view
    returns (
      bool _eligibleToSubmit,
      uint32 _roundId,
      int256 _latestSubmission,
      uint64 _startedAt,
      uint64 _timeout,
      uint128 _availableFunds,
      uint32 _oracleCount,
      uint128 _paymentAmount
    )
  {
    Round storage round = rounds[0];
    OracleStatus storage oracle = oracles[_oracle];

    bool shouldSupersede = oracle.lastReportedRound == reportingRoundId || !acceptingSubmissions(reportingRoundId);
    // Instead of nudging oracles to submit to the next round, the inclusion of
    // the shouldSupersede bool in the if condition pushes them towards
    // submitting in a currently open round.
    if (supersedable(reportingRoundId) && shouldSupersede) {
      _roundId = reportingRoundId.add(1);
      round = rounds[_roundId];

      _paymentAmount = paymentAmount;
      _eligibleToSubmit = delayed(_oracle, _roundId);
    } else {
      _roundId = reportingRoundId;
      round = rounds[_roundId];

      _paymentAmount = round.details.paymentAmount;
      _eligibleToSubmit = acceptingSubmissions(_roundId);
    }

    if (validateOracleRound(_oracle, _roundId).length != 0) {
      _eligibleToSubmit = false;
    }

    return (
      _eligibleToSubmit,
      _roundId,
      oracle.latestSubmission,
      round.startedAt,
      round.details.timeout,
      availableFunds,
      oracleCount(),
      _paymentAmount
    );
  }

  function updateRoundAnswer(uint32 _roundId)
    internal
  {
    if (rounds[_roundId].details.submissions.length < rounds[_roundId].details.minSubmissions) return;

    int256 newAnswer = Median.calculateInplace(rounds[_roundId].details.submissions);
    rounds[_roundId].answer = newAnswer;
    rounds[_roundId].updatedAt = uint64(block.timestamp);
    rounds[_roundId].answeredInRound = _roundId;
    latestRoundId = _roundId;

    validateAnswer(_roundId, newAnswer);

    emit AnswerUpdated(newAnswer, _roundId, now);
  }

  function validateAnswer(
    uint32 _roundId,
    int256 _newAnswer
  )
    private
  {
    AnswerValidatorInterface av = answerValidator; // cache storage reads
    if (address(av) == address(0)) return;

    uint32 prevRound = _roundId.sub(1);
    uint32 prevAnswerRoundId = rounds[prevRound].answeredInRound;
    int256 prevRoundAnswer = rounds[prevRound].answer;
    try av.validate.gas(VALIDATOR_GAS_LIMIT)(
      prevAnswerRoundId,
      prevRoundAnswer,
      _roundId,
      _newAnswer
    ) returns (bool _) { } catch { }
  }

  function payOracle(uint32 _roundId)
    private
  {
    uint128 payment = rounds[_roundId].details.paymentAmount;
    uint128 available = availableFunds.sub(payment);

    availableFunds = available;
    allocatedFunds = allocatedFunds.add(payment);
    oracles[msg.sender].withdrawable = oracles[msg.sender].withdrawable.add(payment);

    emit AvailableFundsUpdated(available);
  }

  function recordSubmission(int256 _submission, uint32 _roundId)
    private
  {
    require(acceptingSubmissions(_roundId), "round not accepting submissions");

    rounds[_roundId].details.submissions.push(_submission);
    oracles[msg.sender].lastReportedRound = _roundId;
    oracles[msg.sender].latestSubmission = _submission;

    emit SubmissionReceived(_submission, _roundId, msg.sender);
  }

  function deleteRoundDetails(uint32 _roundId)
    private
  {
    if (rounds[_roundId].details.submissions.length < rounds[_roundId].details.maxSubmissions) return;

    delete rounds[_roundId].details;
  }

  function timedOut(uint32 _roundId)
    private
    view
    returns (bool)
  {
    uint64 startedAt = rounds[_roundId].startedAt;
    uint32 roundTimeout = rounds[_roundId].details.timeout;
    return startedAt > 0 && roundTimeout > 0 && startedAt.add(roundTimeout) < block.timestamp;
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

  function previousAndCurrentUnanswered(uint32 _roundId, uint32 _rrId)
    private
    view
    returns (bool)
  {
    return _roundId.add(1) == _rrId && rounds[_rrId].updatedAt == 0;
  }

  function requiredReserve(uint256 payment)
    private
    view
    returns (uint256)
  {
    return payment.mul(oracleCount()).mul(RESERVE_ROUNDS);
  }

  function addOracle(
    address _oracle,
    address _admin
  )
    private
  {
    require(!oracleEnabled(_oracle), "oracle already enabled");

    require(_admin != address(0), "cannot set admin to 0");
    require(oracles[_oracle].admin == address(0) || oracles[_oracle].admin == _admin, "owner cannot overwrite admin");

    oracles[_oracle].startingRound = getStartingRound(_oracle);
    oracles[_oracle].endingRound = ROUND_MAX;
    oracles[_oracle].index = uint16(oracleAddresses.length);
    oracleAddresses.push(_oracle);
    oracles[_oracle].admin = _admin;

    emit OraclePermissionsUpdated(_oracle, true);
    emit OracleAdminUpdated(_oracle, _admin);
  }

  function removeOracle(
    address _oracle
  )
    private
  {
    require(oracleEnabled(_oracle), "oracle not enabled");

    oracles[_oracle].endingRound = reportingRoundId.add(1);
    address tail = oracleAddresses[oracleCount().sub(1)];
    uint16 index = oracles[_oracle].index;
    oracles[tail].index = index;
    delete oracles[_oracle].index;
    oracleAddresses[index] = tail;
    oracleAddresses.pop();

    emit OraclePermissionsUpdated(_oracle, false);
  }

  function validateOracleRound(address _oracle, uint32 _roundId)
    private
    view
    returns (bytes memory)
  {
    // cache storage reads
    uint32 startingRound = oracles[_oracle].startingRound;
    uint32 rrId = reportingRoundId;

    if (startingRound == 0) return "not enabled oracle";
    if (startingRound > _roundId) return "not yet enabled oracle";
    if (oracles[_oracle].endingRound < _roundId) return "no longer allowed oracle";
    if (oracles[_oracle].lastReportedRound >= _roundId) return "cannot report on previous rounds";
    if (_roundId != rrId && _roundId != rrId.add(1) && !previousAndCurrentUnanswered(_roundId, rrId)) return "invalid round to report";
    if (_roundId != 1 && !supersedable(_roundId.sub(1))) return "previous round not supersedable";
  }

  function supersedable(uint32 _roundId)
    private
    view
    returns (bool)
  {
    return rounds[_roundId].updatedAt > 0 || timedOut(_roundId);
  }

  function oracleEnabled(address _oracle)
    private
    view
    returns (bool)
  {
    return oracles[_oracle].endingRound == ROUND_MAX;
  }

  function acceptingSubmissions(uint32 _roundId)
    private
    view
    returns (bool)
  {
    return rounds[_roundId].details.maxSubmissions != 0;
  }

  function delayed(address _oracle, uint32 _roundId)
    private
    view
    returns (bool)
  {
    uint256 lastStarted = oracles[_oracle].lastStartedRound;
    return _roundId > lastStarted + restartDelay || lastStarted == 0;
  }

  function newRound(uint32 _roundId)
    private
    view
    returns (bool)
  {
    return _roundId == reportingRoundId.add(1);
  }

}
