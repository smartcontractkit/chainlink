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
  struct Answer {
    uint256 minimumResponses;
    uint256 maxResponses;
    uint256[] responses;
  }

  uint256 public currentRate;
  uint256 public latestCompletedAnswer;
  uint256 public paymentAmount;
  uint256 public minimumResponses;
  bytes32[] public jobIds;
  address[] public oracles;

  uint256 private answerCounter = 1;
  mapping(address => bool) public authorizedRequesters;
  mapping(bytes32 => uint256) private requestAnswers;
  mapping(uint256 => Answer) private answers;

  /**
   * @notice Deploy with the address of the LINK token and arrays of matching
   * length containing the addresses of the oracles and their corresponding
   * Job IDs.
   * @dev Sets the LinkToken address for the network, addresses of the oracles,
   * and jobIds in storage.
   * @param _link The address of the LINK token
   * @param _paymentAmount the amount of LINK to be sent to each oracle for each request
   * @param _minimumResponses the minimum number of responses
   * before an answer will be calculated
   * @param _oracles An array of oracle addresses
   * @param _jobIds An array of Job IDs
   */
  constructor(
    address _link,
    uint256 _paymentAmount,
    uint256 _minimumResponses,
    address[] _oracles,
    bytes32[] _jobIds
  )
    public
    Ownable()
  {
    setChainlinkToken(_link);
    updateRequestDetails(_paymentAmount, _minimumResponses, _oracles, _jobIds);
  }

  /**
   * @notice Creates a Chainlink request for each oracle in the oracles array.
   * @dev This example does not include request parameters. Reference any documentation
   * associated with the Job IDs used to determine the required parameters per-request.
   */
  function requestRateUpdate()
    public
    ensureAuthorizedRequester()
  {
    Chainlink.Request memory request;
    bytes32 requestId;
    uint256 oraclePayment = paymentAmount;

    for (uint i = 0; i < oracles.length; i++) {
      request = buildChainlinkRequest(jobIds[i], this, this.chainlinkCallback.selector);
      requestId = sendChainlinkRequestTo(oracles[i], request, oraclePayment);
      requestAnswers[requestId] = answerCounter;
    }
    answers[answerCounter].minimumResponses = minimumResponses;
    answers[answerCounter].maxResponses = oracles.length;
    answerCounter = answerCounter.add(1);
  }

  /**
   * @notice Receives the answer from the Chainlink node.
   * @dev This function can only be called by the oracle that received the request.
   * @param _clRequestId The Chainlink request ID associated with the answer
   * @param _response The answer provided by the Chainlink node
   */
  function chainlinkCallback(bytes32 _clRequestId, uint256 _response)
    public
  {
    validateChainlinkCallback(_clRequestId);

    uint256 answerId = requestAnswers[_clRequestId];
    delete requestAnswers[_clRequestId];

    answers[answerId].responses.push(_response);
    updateLatestAnswer(answerId);
    deleteAnswer(answerId);
  }

  /**
   * @notice Updates the arrays of oracles and jobIds with new values,
   * overwriting the old values.
   * @dev Arrays are validated to be equal length.
   * @param _paymentAmount the amount of LINK to be sent to each oracle for each request
   * @param _minimumResponses the minimum number of responses
   * before an answer will be calculated
   * @param _oracles An array of oracle addresses
   * @param _jobIds An array of Job IDs
   */
  function updateRequestDetails(
    uint256 _paymentAmount,
    uint256 _minimumResponses,
    address[] _oracles,
    bytes32[] _jobIds
  )
    public
    onlyOwner()
    validateAnswerRequirements(_minimumResponses, _oracles, _jobIds)
  {
    paymentAmount = _paymentAmount;
    minimumResponses = _minimumResponses;
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
    require(link.transfer(_recipient, _amount), "LINK transfer failed");
  }

  /**
   * @notice Called by the owner to permission other addresses to generate new
   * requests to oracles.
   * @param _requester the address whose permissions are being set
   * @param _allowed boolean that determines whether the requester is
   * permissioned or not
   */
  function setAuthorization(address _requester, bool _allowed)
    public
    onlyOwner
  {
    authorizedRequesters[_requester] = _allowed;
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
   * @dev Returns the kth value of the ordered array
   * See: http://www.cs.yale.edu/homes/aspnes/pinewiki/QuickSelect.html
   * @param _list The list of elements to pull from
   * @param _k The index, 1 based, of the elements you want to pull from when ordered
   */
  function quickselect(uint256[] memory _list, uint256 _k)
    private
    returns (uint256)
  {
    uint256 pivot = _list[_list.length / 2];
    uint256[] memory lt = new uint256[](_list.length);
    uint256[] memory gt = new uint256[](_list.length);
    uint256 ltLen = 0;
    uint256 gtLen = 0;
    for (uint256 i = 0; i < _list.length; i++) {
      if (_list[i] < pivot) {
        lt[ltLen] = _list[i];
        ltLen++;
      } else if (_list[i] > pivot) {
        gt[gtLen] = _list[i];
        gtLen++;
      }
    }
    if (_k <= ltLen) {
      return quickselect(trim(lt, ltLen), _k);
    } else if (_k > (_list.length - gtLen)) {
      return quickselect(trim(gt, gtLen), (_k - (_list.length - gtLen)));
    } else {
      return pivot;
    }
  }

  /**
   * @dev Returns a copy of the array, truncated to the specified length.
   * @param _list The list that will be copied and truncated
   * @param _length The length to turncate the list to
   */
  function trim(uint256[] memory _list, uint256 _length)
    private
    returns (uint256[] memory)
  {
      uint256[] memory trimmed = new uint256[](_length);
      for (uint256 i = 0; i < _length; i++) {
        trimmed[i] = _list[i];
      }
      return trimmed;
  }

  /**
   * @dev Performs aggregation of the answers received from the Chainlink nodes.
   * Assumes that at least half the oracles are honest and so can't contol the
   * middle of the ordered responses.
   * @param _answerId The answer ID associated with the group of requests
   */
  function updateLatestAnswer(uint256 _answerId)
    private
    ensureMinResponsesReceived(_answerId)
    ensureOnlyLatestAnswer(_answerId)
  {
    Answer memory answer = answers[_answerId];
    uint256 responseLength = answer.responses.length;
    uint256 middleIndex = responseLength.div(2);
    if (responseLength % 2 == 0) {
      uint256 median1 = quickselect(answers[_answerId].responses, middleIndex);
      uint256 median2 = quickselect(answers[_answerId].responses, middleIndex.add(1));
      currentRate = median1.add(median2).div(2);
    } else {
      currentRate = quickselect(answers[_answerId].responses, middleIndex.add(1));
    }

    latestCompletedAnswer = _answerId;
  }

  /**
   * @dev Cleans up the answer record if all responses have been received.
   * @param _answerId The identifier of the answer to be deleted
   */
  function deleteAnswer(uint256 _answerId)
    private
    ensureAllResponsesReceived(_answerId)
  {
    delete answers[_answerId];
  }


  /**
   * @dev Prevents taking an action if the minimum number of responses has not
   * been received for an answer.
   * @param _answerId The the identifier of the answer that keeps track of the responses.
   */
  modifier ensureMinResponsesReceived(uint256 _answerId) {
    if (answers[_answerId].responses.length >= answers[_answerId].minimumResponses) {
      _;
    }
  }

  /**
   * @dev Prevents taking an action if not all responses are received for an answer.
   * @param _answerId The the identifier of the answer that keeps track of the responses.
   */
  modifier ensureAllResponsesReceived(uint256 _answerId) {
    if (answers[_answerId].responses.length == answers[_answerId].maxResponses) {
      _;
    }
  }

  /**
   * @dev Prevents taking an action if a newer answer has been recorded.
   * @param _answerId The current answer's identifier.
   * Answer IDs are in ascending order.
   */
  modifier ensureOnlyLatestAnswer(uint256 _answerId) {
    if (latestCompletedAnswer <= _answerId) {
      _;
    }
  }

  /**
   * @dev Ensures corresponding number of oracles and jobs.
   * @param _oracles The list of oracles.
   * @param _jobIds The list of jobs.
   */
  modifier validateAnswerRequirements(
    uint256 _minimumResponses,
    address[] _oracles,
    bytes32[] _jobIds
  ) {
    require(_oracles.length >= _minimumResponses, "must have at least as many oracles as responses");
    require(_oracles.length == _jobIds.length, "must have exactly as many oracles as job IDs");
    _;
  }

  /**
   * @dev Reverts if `msg.sender` is not authorized to make requests.
   */
  modifier ensureAuthorizedRequester() {
    require(authorizedRequesters[msg.sender] || msg.sender == owner, "Not an authorized address for creating requests");
    _;
  }

}
