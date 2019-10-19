pragma solidity 0.5.0;

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
    uint128 minimumResponses;
    uint128 paymentAmount;
    int256[] answers;
  }

  int256 public currentAnswer;
  uint256 public answerRound;
  uint128 public paymentAmount;
  uint128 public oracleCount;

  LinkTokenInterface private LINK;
  mapping(address => OracleStatus) private oracles;
  mapping(uint256 => Round) private rounds;

  event NewRound(uint256 indexed number);

  constructor(address _link, uint128 _paymentAmount)
    public
  {
    LINK = LinkTokenInterface(_link);
    updatePaymentAmount(_paymentAmount);
  }

  function updateAnswer(int256 _answer, uint256 _round)
    public
  {
    OracleStatus memory oracle = oracles[msg.sender];
    require(oracle.enabled, "Only updatable by designated oracles");
    require(_round > oracle.lastReportedRound, "Cannot update round reports");
    require(_round == answerRound + 1, "Cannot report on previous rounds");

    if (_round == answerRound + 1) {
      startNewRound(_round);
    }
    rounds[_round].answers.push(_answer);
    currentAnswer = _answer;
    require(LINK.transfer(msg.sender, paymentAmount), "LINK transfer failed");
  }

  function startNewRound(uint256 _number)
    internal
  {
    answerRound = _number;
    rounds[_number].minimumResponses = oracleCount;
    rounds[_number].paymentAmount = paymentAmount;
    emit NewRound(_number);
  }

  function addOracle(address _oracle)
    public
    onlyOwner()
  {
    require(!oracles[_oracle].enabled, "Address is already recorded as an oracle");

    oracles[_oracle].enabled = true;
    oracleCount += 1;
  }

  function removeOracle(address _oracle)
    public
    onlyOwner()
  {
    require(oracles[_oracle].enabled, "Address is not an oracle");
    oracles[_oracle].enabled = false;
    oracleCount -= 1;
  }

  function transferLINK(address _recipient, uint256 _amount)
    public
    onlyOwner()
  {
    require(LINK.transfer(_recipient, _amount), "LINK transfer failed");
  }

  function updatePaymentAmount(uint128 _newAmount)
    public
    onlyOwner()
  {
    paymentAmount = _newAmount;
  }

}
