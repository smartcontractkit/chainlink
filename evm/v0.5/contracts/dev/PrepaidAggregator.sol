pragma solidity 0.5.0;

import "../Median.sol";
import "../vendor/Ownable.sol";
import "../vendor/SafeMath.sol";
import "./SafeMath128.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/WithdrawalInterface.sol";

/**
 * @title The Prepaid Aggregator contract
 * @notice Node handles aggregating data pushed in from off-chain, and unlocks
 * payment for oracles as they report. Oracles' submissions are gathered in
 * rounds, with each round aggregating the submissions for each oracle into a
 * single answer. The latest aggregated answer is exposed as well as historical
 * answers and their updated at timestamp.
 */
contract PrepaidAggregator is Ownable, WithdrawalInterface {
  using SafeMath for uint256;
  using SafeMath128 for uint128;

  struct Round {
    int256 answer;
    uint256 updatedTimestamp;
    RoundDetails details;
  }

  struct RoundDetails {
    int256[] answers;
    uint64 maxAnswers;
    uint64 minAnswers;
    uint128 paymentAmount;
  }

  struct OracleStatus {
    bool enabled;
    uint128 withdrawable;
    uint128 startingRound;
    uint128 lastReportedRound;
    uint128 lastStartedRound;
    int256 latestAnswer;
  }

  uint128 public latestRound;
  uint128 public currentRound;
  uint128 public allocatedFunds;
  uint128 public availableFunds;

  // Round related params
  uint128 public paymentAmount;
  uint64 public oracleCount;
  uint64 public maxAnswerCount;
  uint64 public minAnswerCount;
  uint64 public restartDelay;

  LinkTokenInterface private LINK;
  mapping(address => OracleStatus) private oracles;
  mapping(uint128 => Round) private rounds;

  event NewRound(uint128 indexed number, address indexed startedBy);
  event AnswerUpdated(int256 indexed current, uint128 indexed round);
  event AvailableFundsUpdated(uint256 indexed amount);
  event RoundDetailsUpdated(
    uint128 indexed paymentAmount,
    uint64 indexed minAnswerCount,
    uint64 indexed maxAnswerCount,
    uint64 restartDelay
  );
  event OracleAdded(address indexed oracle);
  event OracleRemoved(address indexed oracle);

  /**
   * @notice Deploy with the address of the LINK token and initial payment amount
   * @dev Sets the LinkToken address and amount of LINK paid
   * @param _link The address of the LINK token
   * @param _paymentAmount The amount paid of LINK paid to each oracle per response
   */
  constructor(address _link, uint128 _paymentAmount) public {
    LINK = LinkTokenInterface(_link);
    updateFutureRounds(_paymentAmount, 0, 0, 0);
  }

  /**
   * @notice called by oracles when they have witnessed a need to update
   * @param _round is the ID of the round this answer pertains to
   * @param _answer is the updated data that the oracle is submitting
   */
  function updateAnswer(uint128 _round, int256 _answer)
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
    uint64 _minAnswers,
    uint64 _maxAnswers,
    uint64 _restartDelay
  )
    external
    onlyOwner()
    onlyUnenabledAddress(_oracle)
  {
    require(oracleCount < 42, "cannot add more than 42 oracles");
    oracles[_oracle].startingRound = currentRound.add(1);
    oracles[_oracle].enabled = true;
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
    uint64 _minAnswers,
    uint64 _maxAnswers,
    uint64 _restartDelay
  )
    external
    onlyOwner()
    onlyEnabledAddress(_oracle)
  {
    oracles[_oracle].enabled = false;
    oracleCount -= 1;

    emit OracleAdded(_oracle);

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
    uint64 _minAnswers,
    uint64 _maxAnswers,
    uint64 _restartDelay
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
    return getAnswer(latestRound);
  }

  /**
   * @notice get the last updated at timestamp
   */
  function updatedTimestamp()
    external
    view
    returns (uint256)
  {
    return getUpdatedTimestamp(latestRound);
  }

  /**
   * @notice get past rounds answers
   * @param _id the round number to retrieve the answer for
   */
  function getAnswer(uint128 _id)
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
  function getUpdatedTimestamp(uint128 _id)
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
    returns (int256, uint256)
  {
    return (oracles[_oracle].latestAnswer, oracles[_oracle].lastReportedRound);
  }

  /**
   * Private
   */

  function startNewRound(uint128 _id)
    private
    onlyOnNewRound(_id)
    onlyIfDelayed(_id)
  {
    currentRound = _id;
    rounds[_id].details.maxAnswers = maxAnswerCount;
    rounds[_id].details.minAnswers = minAnswerCount;
    rounds[_id].details.paymentAmount = paymentAmount;

    oracles[msg.sender].lastStartedRound = _id;

    emit NewRound(_id, msg.sender);
  }

  function updateRoundAnswer(uint128 _id)
    private
    onlyIfMinAnswersReceived(_id)
  {
    int256 newAnswer = Median.calculate(rounds[_id].details.answers);
    rounds[_id].answer = newAnswer;
    rounds[_id].updatedTimestamp = block.number;
    latestRound = _id;

    emit AnswerUpdated(newAnswer, _id);
  }

  function payOracle(uint128 _id)
    private
  {
    uint128 payment = rounds[_id].details.paymentAmount;
    uint128 available = availableFunds.sub(payment);

    availableFunds = available;
    allocatedFunds = allocatedFunds.add(payment);
    oracles[msg.sender].withdrawable = oracles[msg.sender].withdrawable.add(payment);

    emit AvailableFundsUpdated(available);
  }

  function recordSubmission(int256 _answer, uint128 _id)
    private
    onlyIfAcceptingAnswers(_id)
  {
    rounds[_id].details.answers.push(_answer);
    oracles[msg.sender].lastReportedRound = _id;
    oracles[msg.sender].latestAnswer = _answer;
  }

  function deleteRound(uint128 _id)
    private
    onlyIfMaxAnswersReceived(_id)
  {
    delete rounds[_id].details;
  }

  /**
   * Modifiers
   */

  modifier onlyValidOracleRound(uint128 _id) {
    require(oracles[msg.sender].enabled, "Only updatable by designated oracles");
    require(oracles[msg.sender].startingRound <= _id, "New oracles cannot participate in in-progress rounds");
    require(_id > oracles[msg.sender].lastReportedRound, "Cannot update round reports");
    _;
  }

  modifier onlyIfMinAnswersReceived(uint128 _id) {
    if (rounds[_id].details.answers.length >= rounds[_id].details.minAnswers) {
      _;
    }
  }

  modifier onlyIfMaxAnswersReceived(uint128 _id) {
    if (rounds[_id].details.answers.length == rounds[_id].details.maxAnswers) {
      _;
    }
  }

  modifier onlyIfAcceptingAnswers(uint128 _id) {
    require(rounds[_id].details.maxAnswers != 0, "Max responses reached for round");
    _;
  }

  modifier onlyOnNewRound(uint128 _id) {
    if (_id == currentRound.add(1)) {
      _;
    }
  }

  modifier onlyIfDelayed(uint128 _id) {
    uint256 lastStarted = oracles[msg.sender].lastStartedRound;
    if (_id > lastStarted + restartDelay) {
      _;
    }
  }

  modifier onlyValidRoundId(uint128 _id) {
    require(_id == currentRound || _id == currentRound.add(1), "Must report on current round");
    if (_id > 1) {
      require(rounds[_id.sub(1)].updatedTimestamp > 0, "Cannot bump round until previous round has an answer");
    }
    _;
  }

  modifier onlyValidRange(uint64 _min, uint64 _max, uint64 _restartDelay) {
    uint64 oracleNum = oracleCount; // Save on storage reads
    require(oracleNum >= _max, "Cannot have the answer max higher oracle count");
    require(_max >= _min, "Cannot have the answer minimum higher the max");
    if (oracleNum > 0) {
      require(oracleNum > _restartDelay, "Restart delay must be less than oracle count");
    }
    _;
  }

  modifier onlyUnenabledAddress(address _oracle) {
    require(!oracles[_oracle].enabled, "Address is already recorded as an oracle");
    _;
  }

  modifier onlyEnabledAddress(address _oracle) {
    require(oracles[_oracle].enabled, "Address is not an oracle");
    _;
  }

}
