pragma solidity 0.5.0;

import "../Median.sol";
import "../vendor/Ownable.sol";
import "../vendor/SafeMath.sol";
import "./SafeMath128.sol";
import "../interfaces/LinkTokenInterface.sol";

/**
 * @title The PrepaidAggregator handles aggregating data pushed in from off-chain.
 */
contract PrepaidAggregator is Ownable {
  using SafeMath for uint256;
  using SafeMath128 for uint128;

  struct OracleStatus {
    bool enabled;
    uint128 withdrawable;
    uint256 lastReportedRound;
  }

  struct Round {
    uint64 maxAnswers;
    uint64 minAnswers;
    uint128 paymentAmount;
    int256[] answers;
  }

  int256 public currentAnswer;
  uint256 public currentRound;
  uint128 public paymentAmount;
  uint64 public oracleCount;
  uint64 public maxAnswerCount;
  uint64 public minAnswerCount;
  uint64 public roundRestartDelay;
  uint256 public updatedHeight;
  uint128 public availableFunds;
  uint128 public allocatedFunds;

  LinkTokenInterface private LINK;
  mapping(address => OracleStatus) private oracles;
  mapping(uint256 => Round) private rounds;

  event NewRound(uint256 indexed number, address indexed startedBy);
  event AnswerUpdated(int256 indexed current, uint256 indexed round);
  event AvailableFundsUpdated(uint256 indexed amount);
  event PaymentAmountUpdated(uint128 indexed amount);
  event RoundDetailsUpdated(
    uint128 indexed paymentAmount,
    uint64 indexed minAnswerCount,
    uint64 indexed maxAnswerCount,
    uint64 roundRestartDelay
  );

  constructor(address _link, uint128 _paymentAmount) public {
    LINK = LinkTokenInterface(_link);
    setPaymentAmount(_paymentAmount);
  }

  function updateAnswer(uint256 _round, int256 _answer)
    public
    onlyValidRoundId(_round)
    onlyValidOracleRound(_round)
  {
    startNewRound(_round);
    recordAnswer(_answer, _round);
    updateRoundAnswer(_round);
    payOracle(_round);
    deleteRound(_round);
  }

  function addOracle(
    address _oracle,
    uint64 _minAnswers,
    uint64 _maxAnswers,
    uint64 _roundRestartDelay
  )
    public
    onlyOwner()
    onlyUnenabledAddress(_oracle)
  {
    require(oracleCount < 42, "cannot add more than 42 oracles");
    oracles[_oracle].enabled = true;
    oracleCount += 1;
    setAnswerCountRange(_minAnswers, _maxAnswers, _roundRestartDelay);
  }

  function removeOracle(
    address _oracle,
    uint64 _minAnswers,
    uint64 _maxAnswers,
    uint64 _roundRestartDelay
  )
    public
    onlyOwner()
    onlyEnabledAddress(_oracle)
  {
    oracles[_oracle].enabled = false;
    oracleCount -= 1;
    setAnswerCountRange(_minAnswers, _maxAnswers, _roundRestartDelay);
  }

  function setPaymentAmount(uint128 _newAmount)
    public
    onlyOwner()
  {
    paymentAmount = _newAmount;
    emit PaymentAmountUpdated(_newAmount);
  }

  function setAnswerCountRange(
    uint64 _minAnswerCount,
    uint64 _maxAnswerCount,
    uint64 _roundRestartDelay
  )
    public
    onlyOwner()
    onlyValidRange(_minAnswerCount, _maxAnswerCount)
  {
    minAnswerCount = _minAnswerCount;
    maxAnswerCount = _maxAnswerCount;
    roundRestartDelay = _roundRestartDelay;

    emit RoundDetailsUpdated(paymentAmount, _minAnswerCount, _maxAnswerCount, _roundRestartDelay);
  }

  function updateAvailableFunds()
    public
  {
    uint256 available = LINK.balanceOf(address(this)).sub(allocatedFunds);
    availableFunds = uint128(available);
    emit AvailableFundsUpdated(available);
  }

  function withdrawable()
    public
    returns (uint256)
  {
    return oracles[msg.sender].withdrawable;
  }

  function withdraw(address _recipient, uint256 _amount)
    public
  {
    uint128 amount = uint128(_amount);
    uint128 available = oracles[msg.sender].withdrawable;
    require(available >= amount, "Insufficient balance");

    oracles[msg.sender].withdrawable = available.sub(amount);
    allocatedFunds = allocatedFunds.sub(amount);

    assert(LINK.transfer(_recipient, _amount));
  }

  function withdrawFunds(address _recipient, uint256 _amount)
    public
    onlyOwner()
  {
    require(availableFunds >= _amount, "Insufficient funds");
    require(LINK.transfer(_recipient, _amount), "LINK transfer failed");
    updateAvailableFunds();
  }

  function startNewRound(uint256 _id)
    private
    onlyOnNewRound(_id)
  {
    currentRound = _id;
    rounds[_id].maxAnswers = maxAnswerCount;
    rounds[_id].minAnswers = minAnswerCount;
    rounds[_id].paymentAmount = paymentAmount;
    emit NewRound(_id, msg.sender);
  }

  function updateRoundAnswer(uint256 _id)
    private
    onlyIfMinAnswersReceived(_id)
  {
    int256 newAnswer = Median.calculate(rounds[_id].answers);
    currentAnswer = newAnswer;
    updatedHeight = block.number;
    emit AnswerUpdated(newAnswer, _id);
  }

  function payOracle(uint256 _id)
    private
  {
    uint128 payment = rounds[_id].paymentAmount;
    // SafeMath128's underflow check ensures that there are enough funds to pay the oracle.
    uint128 available = availableFunds.sub(payment);

    availableFunds = available;
    allocatedFunds = allocatedFunds.add(payment);
    oracles[msg.sender].withdrawable = oracles[msg.sender].withdrawable.add(payment);

    emit AvailableFundsUpdated(available);
  }

  function recordAnswer(int256 _answer, uint256 _id)
    private
    onlyIfAcceptingAnswers(_id)
  {
    rounds[_id].answers.push(_answer);
    oracles[msg.sender].lastReportedRound = _id;
  }

  function deleteRound(uint256 _id)
    private
    onlyIfMaxAnswersReceived(_id)
  {
    delete rounds[_id];
  }

  modifier onlyValidOracleRound(uint256 _round) {
    require(oracles[msg.sender].enabled, "Only updatable by designated oracles");
    require(_round > oracles[msg.sender].lastReportedRound, "Cannot update round reports");
    _;
  }

  modifier onlyIfMinAnswersReceived(uint256 _id) {
    if (rounds[_id].answers.length == rounds[_id].minAnswers) {
      _;
    }
  }

  modifier onlyIfMaxAnswersReceived(uint256 _id) {
    if (rounds[_id].answers.length == rounds[_id].maxAnswers) {
      _;
    }
  }

  modifier onlyIfAcceptingAnswers(uint256 _id) {
    require(rounds[_id].maxAnswers != 0, "Max responses reached for round");
    _;
  }

  modifier onlyOnNewRound(uint256 _id) {
    if (_id == currentRound.add(1)) {
      _;
    }
  }

  modifier onlyValidRoundId(uint256 _id) {
    require(_id == currentRound.add(1) || _id == currentRound, "Cannot report on previous rounds");
    _;
  }

  modifier onlyValidRange(uint64 _min, uint64 _max) {
    require(oracleCount >= _max, "Cannot have the answer max higher oracle count");
    require(_max >= _min, "Cannot have the answer minimum higher the max");
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
