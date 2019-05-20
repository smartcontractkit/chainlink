pragma solidity 0.4.24;

import "../ChainlinkClient.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

/**
 * @title An example Chainlink contract with aggregation
 * @notice Requesters can use this contract as a framework for creating
 * requests to multiple Chainlink nodes and running aggregation
 * as the contract receives answers.
 */
contract ConversionRate is ChainlinkClient, Ownable {
  uint256 constant private ORACLE_PAYMENT = 1 * LINK; // solium-disable-line zeppelin/no-arithmetic-operations

  struct Answer {
    uint256 expectedResponses;
    uint256[] responses;
  }

  uint256 public currentRate;
  uint256 public latestCompletedAnswer;
  bytes32[] public jobIds;
  address[] public oracles;

  uint256 private answerCounter = 1;
  mapping(bytes32 => uint256) private requestAnswers;
  mapping(uint256 => Answer) private answers;

  /**
   * @notice Deploy with the address of the LINK token and arrays of matching
   * length containing the addresses of the oracles and their corresponding
   * Job IDs.
   * @dev Sets the LinkToken address for the network, addresses of the oracles,
   * and jobIds in storage.
   * @param _link The address of the LINK token
   * @param _oracles An array of oracle addresses
   * @param _jobIds An array of Job IDs
   */
  constructor(address _link, address[] _oracles, bytes32[] _jobIds)
    public
    Ownable()
  {
    setChainlinkToken(_link);
    updateOracles(_oracles, _jobIds);
  }

  /**
   * @notice Creates a Chainlink request for each oracle in the oracles array.
   * @dev This example does not include request parameters. Reference any documentation
   * associated with the Job IDs used to determine the required parameters per-request.
   */
  function requestRateUpdate()
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

  /**
   * @notice Receives the answer from the Chainlink node.
   * @dev This function can only be called by the oracle that received the request.
   * @param _clRequestId The Chainlink request ID associated with the answer
   * @param _rate The answer provided by the Chainlink node
   */
  function chainlinkCallback(bytes32 _clRequestId, uint256 _rate)
    public
  {
    validateChainlinkCallback(_clRequestId);

    uint256 answerId = requestAnswers[_clRequestId];
    delete requestAnswers[_clRequestId];
    answers[answerId].responses.push(_rate);

    updateRecords(answerId);
  }

  /**
   * @notice Updates the arrays of oracles and jobIds with new values,
   * overwriting the old values.
   * @dev Arrays are validated to be equal length.
   * @param _oracles An array of oracle addresses
   * @param _jobIds An array of Job IDs
   */
  function updateOracles(address[] _oracles, bytes32[] _jobIds)
    public
    onlyOwner()
    checkEqualLengths(_oracles, _jobIds)
  {
    jobIds = _jobIds;
    oracles = _oracles;
  }

  /**
   * @notice Allows the owner of the contract to withdraw any LINK balance
   * available on the contract.
   * @dev The contract will need to have a LINK balance in order to create requests.
   * @param _recipient The address to receive the LINK tokens
   * @param _amount The amount of LINK to send from the contract
   */
  function transferLINK(address _recipient, uint256 _amount)
    public
    onlyOwner()
  {
    LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
    require(link.transfer(_recipient, _amount));
  }

  /**
   * @notice Called by the owner to kill the contract. This transfers all LINK
   * balance and ETH balance (if there is any) to the owner.
   */
  function destroy()
    public
    onlyOwner()
  {
    LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
    transferLINK(owner, link.balanceOf(address(this)));
    selfdestruct(owner);
  }

  /**
   * @dev Due to the ensureAllResponsesReceived modifier, this function
   * will only run if all answers have been received by the _answerId and
   * will clean up the answers mapping before running aggregation.
   * @param _answerId The answer ID associated with the group of requests
   */
  function updateRecords(uint256 _answerId)
    private
    ensureAllResponsesReceived(_answerId)
  {
    updateRate(_answerId);
    delete answers[_answerId];
  }

  /**
   * @dev Performs aggregation of the answers received from the Chainlink nodes.
   * @param _answerId The answer ID associated with the group of requests
   */
  function updateRate(uint256 _answerId)
    private
    ensureOnlyLatestAnswer(_answerId)
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

  /**
   * @dev Prevents taking an action if not all responses are received for an answer.
   * @param _answerId The the identifier of the answer that keeps track of the responses.
   */
  modifier ensureAllResponsesReceived(uint256 _answerId) {
    if (answers[_answerId].responses.length == answers[_answerId].expectedResponses) {
      _;
    }
  }

  /**
   * @dev Prevents taking an action if a newer answer has been recorded.
   * @param _answerId The current answer's identifier.
   * Answer IDs are in ascending order.
   */
  modifier ensureOnlyLatestAnswer(uint256 _answerId) {
    if (latestCompletedAnswer < _answerId) {
      _;
    }
  }

  /**
   * @dev Ensures corresponding number of oracles and jobs.
   * @param _oracles The list of oracles.
   * @param _jobIds The list of jobs.
   */
  modifier checkEqualLengths(address[] _oracles, bytes32[] _jobIds) {
    require(_oracles.length == _jobIds.length);
    _;
  }

}
