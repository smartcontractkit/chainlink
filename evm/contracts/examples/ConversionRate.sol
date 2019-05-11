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
  uint256 private requestCounter = 1;
  uint256 private latestAnswer;
  mapping(bytes32 => uint256) private requestMap;
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
      requestMap[requestId] = requestCounter;
    }
    answers[requestCounter].expectedResponses = oracles.length;
    requestCounter = requestCounter.add(1);
  }

  function chainlinkCallback(bytes32 _clRequestId, uint256 _rate)
    public
  {
    validateChainlinkCallback(_clRequestId);

    uint256 requestId = requestMap[_clRequestId];
    delete requestMap[_clRequestId];
    answers[requestId].responses.push(_rate);

    updateRecords(requestId);
  }

  function updateRecords(uint256 _requestId)
    private
    checkAllResponsesReceived(_requestId)
  {
    updateRate(_requestId);
    delete answers[_requestId];
  }

  function updateRate(uint256 _requestId)
    private
    checkLatestAnswer(_requestId)
  {
    uint256 sumQuotients;
    uint256 sumRemainders;
    Answer memory answer = answers[_requestId];
    for (uint i = 0; i < answer.expectedResponses; i++) {
      uint256 response = answer.responses[i];
      sumQuotients = sumQuotients.add(response.div(answer.expectedResponses)); // aggregate responses and protect from overflows
      sumRemainders = sumRemainders.add(response % answer.expectedResponses); 
    }
    currentRate = sumQuotients.add(sumRemainders.div(answer.expectedResponses)); // recover lost accuracy from result 
    latestAnswer = _requestId;
  }

  modifier checkAllResponsesReceived(uint256 _requestId) {
    if (answers[_requestId].responses.length == answers[_requestId].expectedResponses) {
      _;
    }
  }

  modifier checkLatestAnswer(uint256 _requestId) {
    if (latestAnswer < _requestId) {
      _;
    }
  }

}
