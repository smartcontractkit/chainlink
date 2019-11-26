pragma solidity 0.5.0;

import "../Median.sol";
import "../vendor/Ownable.sol";
import "../vendor/SafeMath.sol";
import "./SafeMath128.sol";
import "./SafeMath32.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/WithdrawalInterface.sol";
import "./AggregatorInterface.sol";

/**
 * @title The Prepaid Aggregator contract
 * @notice Node handles aggregating data pushed in from off-chain, and unlocks
 * payment for oracles as they report. Oracles' submissions are gathered in
 * rounds, with each round aggregating the submissions for each oracle into a
 * single answer. The latest aggregated answer is exposed as well as historical
 * answers and their updated at timestamp.
 */
contract PrepaidAggregator is AggregatorInterface, Ownable, WithdrawalInterface {
  using SafeMath for uint256;
  using SafeMath128 for uint128;
  using SafeMath32 for uint32;

  struct Round {
    int256 answer;
    uint256 updatedTimestamp;
    RoundDetails details;
  }

  struct RoundDetails {
    int256[] answers;
    uint32 maxAnswers;
    uint32 minAnswers;
    uint128 paymentAmount;
  }

  struct OracleStatus {
    uint128 withdrawable;
    uint32 startingRound;
    uint32 endingRound;
    uint32 lastReportedRound;
    uint32 lastStartedRound;
    int256 latestAnswer;
  }

  uint32 private latestRoundValue;
  uint32 public currentRound;
  uint128 public allocatedFunds;
  uint128 public availableFunds;

  // Round related params
  uint128 public paymentAmount;
  uint32 public oracleCount;
  uint32 public maxAnswerCount;
  uint32 public minAnswerCount;
  uint32 public restartDelay;

  LinkTokenInterface private LINK;
  mapping(address => OracleStatus) private oracles;
  mapping(uint32 => Round) private rounds;

  event NewRound(uint32 indexed number, address indexed startedBy);
  event AnswerUpdated(int256 indexed current, uint32 indexed round);
  event AvailableFundsUpdated(uint256 indexed amount);
  event RoundDetailsUpdated(
    uint128 indexed paymentAmount,
    uint32 indexed minAnswerCount,
    uint32 indexed maxAnswerCount,
    uint32 restartDelay
  );
  event OracleAdded(address indexed oracle);
  event OracleRemoved(address indexed oracle);

  uint32 constant private ROUND_MAX = 2**32-1;

  /**
   * @notice Deploy with the address of the LINK token and initial payment amount
   * @dev Sets the LinkToken address and amount of LINK paid
   * @param _link The address of the LINK token
   * @param _paymentAmount The amount paid of LINK paid to each oracle per response
   */
  constructor(address _link, uint128 _paymentAmount) public {
    LINK = LinkTokenInterface(_link);
    paymentAmount = _paymentAmount;
  }

  /**
   * @notice called by oracles when they have witnessed a need to update
   * @param _round is the ID of the round this answer pertains to
   * @param _answer is the updated data that the oracle is submitting
   */
  function updateAnswer(uint32 _round, int256 _answer)
    external
    onlyValidRoundId(_round)
    onlyValidOracleRound(_round)
  {
    startNewRound(_round);
    recordSubmission(_answer, _round);
    updateRoundAnswer(_round);
    payOracle(_round);
    deleteRound(_round);
  }

  /**
   * @notice called by the owner to add a new Oracle and update the round
   * related parameters
   * @param _oracle is the address of the new Oracle being added
   * @param _minAnswers is the new minimum answer count for each round
   * @param _maxAnswers is the new maximum answer count for each round
   * @param _restartDelay is the number of rounds an Oracle has to wait before
   * they can initiate a round
   */
  function addOracle(
    address _oracle,
    uint32 _minAnswers,
    uint32 _maxAnswers,
    uint32 _restartDelay
  )
    external
    onlyOwner()
    onlyUnenabledAddress(_oracle)
  {
    require(oracleCount < 42, "cannot add more than 42 oracles");
    oracles[_oracle].startingRound = currentRound.add(1);
    oracles[_oracle].endingRound = ROUND_MAX;
    oracleCount += 1;

    emit OracleAdded(_oracle);

    updateFutureRounds(paymentAmount, _minAnswers, _maxAnswers, _restartDelay);
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
    oracleCount -= 1;
    oracles[_oracle].endingRound = currentRound;

    emit OracleRemoved(_oracle);

    updateFutureRounds(paymentAmount, _minAnswers, _maxAnswers, _restartDelay);
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
    uint32 _restartDelay
  )
    public
    onlyOwner()
    onlyValidRange(_minAnswers, _maxAnswers, _restartDelay)
  {
    paymentAmount = _newPaymentAmount;
    minAnswerCount = _minAnswers;
    maxAnswerCount = _maxAnswers;
    restartDelay = _restartDelay;

    emit RoundDetailsUpdated(
      paymentAmount,
      _minAnswers,
      _maxAnswers,
      _restartDelay
    );
  }

  /**
   * @notice recalculate the amount of LINK available for payouts
   */
  function updateAvailableFunds()
    public
  {
    uint256 available = LINK.balanceOf(address(this)).sub(allocatedFunds);
    availableFunds = uint128(available);
    emit AvailableFundsUpdated(available);
  }

  /**
   * @notice query the available amount of LINK for an oracle to withdraw
   */
  function withdrawable()
    external
    view
    returns (uint256)
  {
    return uint256(oracles[msg.sender].withdrawable);
  }

  /**
   * @notice get the most recently reported answer
   */
  function currentAnswer()
    external
    view
    returns (int256)
  {
    return getAnswer(latestRoundValue);
  }

  /**
   * @notice get the last updated at timestamp
   */
  function updatedTimestamp()
    external
    view
    returns (uint256)
  {
    return getUpdatedTimestamp(latestRoundValue);
  }

  /**
   * @notice get the last updated round
   */
  function latestRound()
    external
    view
    returns (uint256)
  {
    return uint256(latestRoundValue);
  }

  /**
   * @notice get past rounds answers
   * @param _id the round number to retrieve the answer for
   */
  function getAnswer(uint32 _id)
    public
    view
    returns (int256)
  {
    return rounds[_id].answer;
  }

  /**
   * @notice get timestamp when an answer was last updated
   * @param _id the round number to retrieve the updated timestamp for
   */
  function getUpdatedTimestamp(uint32 _id)
    public
    view
    returns (uint256)
  {
    return rounds[_id].updatedTimestamp;
  }

  /**
   * @notice transfers the oracle's LINK to another address
   * @param _recipient is the address to send the LINK to
   * @param _amount is the amount of LINK to send
   */
  function withdraw(address _recipient, uint256 _amount)
    external
  {
    uint128 amount = uint128(_amount);
    uint128 available = oracles[msg.sender].withdrawable;
    require(available >= amount, "Insufficient balance");

    oracles[msg.sender].withdrawable = available.sub(amount);
    allocatedFunds = allocatedFunds.sub(amount);

    assert(LINK.transfer(_recipient, _amount));
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
    returns (int256, uint32)
  {
    return (oracles[_oracle].latestAnswer, oracles[_oracle].lastReportedRound);
  }

  /**
   * @notice allows the owner to force a new round if the old round could not
   * be completed
   */
  function forceNewRound()
    external
    onlyOwner()
  {
    uint32 id = currentRound;
    rounds[id].updatedTimestamp = block.timestamp;
    startNewRound(id + 1);
  }

  /**
   * Private
   */

  function startNewRound(uint32 _id)
    private
    onlyOnNewRound(_id)
    onlyIfDelayedOrOwner(_id)
  {
    currentRound = _id;
    rounds[_id].details.maxAnswers = maxAnswerCount;
    rounds[_id].details.minAnswers = minAnswerCount;
    rounds[_id].details.paymentAmount = paymentAmount;

    recordStartedRound(_id);

    emit NewRound(_id, msg.sender);
  }

  function recordStartedRound(uint32 _id)
    private
    onlyNonOwner()
  {
    oracles[msg.sender].lastStartedRound = _id;
  }

  function updateRoundAnswer(uint32 _id)
    private
    onlyIfMinAnswersReceived(_id)
  {
    int256 newAnswer = Median.calculate(rounds[_id].details.answers);
    rounds[_id].answer = newAnswer;
    rounds[_id].updatedTimestamp = block.timestamp;
    latestRoundValue = _id;

    emit AnswerUpdated(newAnswer, _id);
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
    onlyIfAcceptingAnswers(_id)
  {
    rounds[_id].details.answers.push(_answer);
    oracles[msg.sender].lastReportedRound = _id;
    oracles[msg.sender].latestAnswer = _answer;
  }

  function deleteRound(uint32 _id)
    private
    onlyIfMaxAnswersReceived(_id)
  {
    delete rounds[_id].details;
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

  modifier onlyIfMinAnswersReceived(uint32 _id) {
    if (rounds[_id].details.answers.length >= rounds[_id].details.minAnswers) {
      _;
    }
  }

  modifier onlyIfMaxAnswersReceived(uint32 _id) {
    if (rounds[_id].details.answers.length == rounds[_id].details.maxAnswers) {
      _;
    }
  }

  modifier onlyIfAcceptingAnswers(uint32 _id) {
    require(rounds[_id].details.maxAnswers != 0, "Max responses reached for round");
    _;
  }

  modifier onlyOnNewRound(uint32 _id) {
    if (_id == currentRound.add(1)) {
      _;
    }
  }

  modifier onlyIfDelayedOrOwner(uint32 _id) {
    uint256 lastStarted = oracles[msg.sender].lastStartedRound;
    if (_id > lastStarted + restartDelay || isOwner()) {
      _;
    }
  }

  modifier onlyValidRoundId(uint32 _id) {
    require(_id == currentRound || _id == currentRound.add(1), "Must report on current round");
    require(rounds[_id.sub(1)].updatedTimestamp > 0 || _id == 1, "Cannot bump round until previous round has an answer");
    _;
  }

  modifier onlyValidRange(uint32 _min, uint32 _max, uint32 _restartDelay) {
    uint32 oracleNum = oracleCount; // Save on storage reads
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

  modifier onlyNonOwner() {
    if (!isOwner()) {
      _;
    }
  }

}
