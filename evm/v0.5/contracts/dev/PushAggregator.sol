pragma solidity 0.5.0;

import "./Quickselectable.sol";
import "../vendor/Ownable.sol";
import "../vendor/SafeMath.sol";
import "../vendor/SignedSafeMath.sol";
import "../interfaces/LinkTokenInterface.sol";

/**
 * @title The PushAggregator handles aggregating data pushed in from off-chain.
 */
contract PushAggregator is Ownable, Quickselectable {
  using SafeMath for uint256;
  using SignedSafeMath for int256;

  struct OracleStatus {
    bool enabled;
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
  uint128 public oracleCount;
  uint64 public maxAnswerCount;
  uint64 public minAnswerCount;

  LinkTokenInterface private LINK;
  mapping(address => OracleStatus) private oracles;
  mapping(uint256 => Round) private rounds;

  event NewRound(uint256 indexed number);
  event AnswerUpdated(int256 indexed current, uint256 indexed round);

  constructor(address _link, uint128 _paymentAmount)
    public
  {
    LINK = LinkTokenInterface(_link);
    setPaymentAmount(_paymentAmount);
  }

  function updateAnswer(int256 _answer, uint256 _round)
    public
    ensureValidRoundId(_round)
    validateOracleRound(_round)
  {
    startNewRound(_round);
    recordAnswer(_answer, _round);
    updateRoundAnswer(_round);
    deleteRound(_round);

    require(LINK.transfer(msg.sender, paymentAmount), "LINK transfer failed");
  }

  function addOracle(address _oracle)
    public
    onlyOwner()
  {
    require(!oracles[_oracle].enabled, "Address is already recorded as an oracle");

    oracles[_oracle].enabled = true;
    oracleCount += 1;
    setAnswerCountRange(minAnswerCount + 1, maxAnswerCount + 1);
  }

  function removeOracle(address _oracle)
    public
    onlyOwner()
  {
    require(oracles[_oracle].enabled, "Address is not an oracle");
    oracles[_oracle].enabled = false;
    oracleCount -= 1;
    setAnswerCountRange(minAnswerCount - 1, maxAnswerCount - 1);
  }

  function transferLINK(address _recipient, uint256 _amount)
    public
    onlyOwner()
  {
    require(LINK.transfer(_recipient, _amount), "LINK transfer failed");
  }

  function setPaymentAmount(uint128 _newAmount)
    public
    onlyOwner()
  {
    paymentAmount = _newAmount;
  }

  function setAnswerCountRange(uint64 _min, uint64 _max)
    public
    onlyOwner()
  {
    require(oracleCount >= _max, "Cannot have the answer max higher oracle count");
    require(_max >= _min, "Cannot have the answer minimum higher the max");
    minAnswerCount = _min;
    maxAnswerCount = _max;
  }

  function startNewRound(uint256 _id)
    private
    ensureNextRound(_id)
  {
    currentRound = _id;
    rounds[_id].maxAnswers = maxAnswerCount;
    rounds[_id].minAnswers = minAnswerCount;
    rounds[_id].paymentAmount = paymentAmount;
    emit NewRound(_id);
  }

  function updateRoundAnswer(uint256 _id)
    private
    ensureMinAnswersReceived(_id)
  {
    uint256 answerLength = rounds[_id].answers.length;
    uint256 middleIndex = answerLength.div(2);
    if (answerLength % 2 == 0) {
      int256 median1 = quickselect(rounds[_id].answers, middleIndex);
      int256 median2 = quickselect(rounds[_id].answers, middleIndex.add(1)); // quickselect is 1 indexed
      currentAnswer = median1.add(median2) / 2; // signed integers are not supported by SafeMath
    } else {
      currentAnswer = quickselect(rounds[_id].answers, middleIndex.add(1)); // quickselect is 1 indexed
    }
    emit AnswerUpdated(currentAnswer, _id);
  }

  function recordAnswer(int256 _answer, uint256 _id)
    private
    ensureAcceptingAnswers(_id)
  {
    rounds[_id].answers.push(_answer);
  }

  function deleteRound(uint256 _id)
    private
    ensureMaxAnswersReceived(_id)
  {
    delete rounds[_id];
  }

  modifier validateOracleRound(uint256 _round) {
    require(oracles[msg.sender].enabled, "Only updatable by designated oracles");
    require(_round > oracles[msg.sender].lastReportedRound, "Cannot update round reports");
    _;
  }

  modifier ensureMinAnswersReceived(uint256 _id) {
    if (rounds[_id].answers.length == rounds[_id].minAnswers) {
      _;
    }
  }

  modifier ensureMaxAnswersReceived(uint256 _id) {
    if (rounds[_id].answers.length == rounds[_id].maxAnswers) {
      _;
    }
  }

  modifier ensureAcceptingAnswers(uint256 _id) {
    require(rounds[_id].maxAnswers != 0, "Max number of responses already received for round");
    _;
  }

  modifier ensureNextRound(uint256 _id) {
    if (_id == currentRound + 1) {
      _;
    }
  }

  modifier ensureValidRoundId(uint256 _id) {
    require(_id == currentRound + 1 || _id == currentRound, "Cannot report on previous rounds");
    _;
  }
}
