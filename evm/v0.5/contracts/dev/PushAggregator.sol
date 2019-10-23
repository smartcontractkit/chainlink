pragma solidity 0.5.0;

import "../Median.sol";
import "../vendor/Ownable.sol";
import "../interfaces/LinkTokenInterface.sol";

/**
 * @title The PushAggregator handles aggregating data pushed in from off-chain.
 */
contract PushAggregator is Ownable {

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

  function updateAnswer(uint256 _round, int256 _answer)
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

    uint64 min = minAnswerCount;
    uint64 max = maxAnswerCount;
    if (min > 0) {
      min = min - 1;
    }
    if (max > 0) {
      max = max - 1;
    }
    setAnswerCountRange(min, max);
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
    int256 newAnswer = Median.get(rounds[_id].answers);
    currentAnswer = newAnswer;
    emit AnswerUpdated(newAnswer, _id);
  }

  function recordAnswer(int256 _answer, uint256 _id)
    private
    ensureAcceptingAnswers(_id)
  {
    rounds[_id].answers.push(_answer);
    oracles[msg.sender].lastReportedRound = _id;
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
