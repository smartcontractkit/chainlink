pragma solidity 0.4.24;

import "../ChainlinkClient.sol";

contract ConversionRate is ChainlinkClient {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  struct Answer {
    uint256 expectedResponses;
    uint256[] responses;
  }

  uint256 public currentRate;
  bytes32[] public jobIds;
  address[] public oracles;
  uint256 private answerCounter = 1;
  uint256 private latestCompletedAnswer;
  mapping(bytes32 => uint256) private requestAnswers;
  mapping(uint256 => Answer) private answers;

  constructor(address _link, address[] _oracles, bytes32[] _jobIds) public {
    setChainlinkToken(_link);
    jobIds = _jobIds;
    oracles = _oracles;
  }

  function update()
    public
  {
    Chainlink.Request memory request;
    bytes32 requestId;

    for (uint i = 0; i < oracles.length; i++) {
      request = buildChainlinkRequest(jobIds[i], this, this.chainlinkCallback.selector);
      requestId = sendChainlinkRequestTo(oracles[i], request, ORACLE_PAYMENT);
      requestAnswers[requestId] = answerCounter;
    }
    answers[answerCounter].expectedResponses = oracles.length;
    answerCounter = answerCounter.add(1);
  }

  function chainlinkCallback(bytes32 _clRequestId, uint256 _rate)
    public
  {
    validateChainlinkCallback(_clRequestId);

    uint256 answerId = requestAnswers[_clRequestId];
    delete requestAnswers[_clRequestId];
    answers[answerId].responses.push(_rate);

    updateRecords(answerId);
  }

  function updateRecords(uint256 _answerId)
    private
    checkAllResponsesReceived(_answerId)
  {
    emit Here(1111, _answerId);
    updateRate(_answerId);
    delete answers[_answerId];
  }

  event Here(uint256 a, uint256 b);
  function updateRate(uint256 _answerId)
    private
    checkLatestAnswer(_answerId)
  {
    uint256 sumQuotients;
    uint256 sumRemainders;
    Answer memory answer = answers[_answerId];
    for (uint i = 0; i < answer.expectedResponses; i++) {
      uint256 response = answer.responses[i];
      sumQuotients = sumQuotients.add(response.div(answer.expectedResponses)); // aggregate responses and protect from overflows
      sumRemainders = sumRemainders.add(response % answer.expectedResponses); 
    }
    currentRate = sumQuotients.add(sumRemainders.div(answer.expectedResponses)); // recover lost accuracy from result 
    latestCompletedAnswer = _answerId;
  }

  modifier checkAllResponsesReceived(uint256 _answerId) {
    emit Here(answers[_answerId].responses.length, answers[_answerId].expectedResponses);
    if (answers[_answerId].responses.length == answers[_answerId].expectedResponses) {
      _;
    }
  }

  modifier checkLatestAnswer(uint256 _answerId) {
    if (latestCompletedAnswer < _answerId) {
      _;
    }
  }

}
