// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

import "./ChainlinkClient.sol";
import "./LinkTokenReceiver.sol";
import "./Median.sol";
import "./vendor/Ownable.sol";
import "./vendor/SafeMathChainlink.sol";

/**
 * @title PreCoordinator is a contract that builds on-chain service agreements
 * using the current architecture of 1 request to 1 oracle contract.
 * @dev This contract accepts requests as service agreement IDs and loops over
 * the corresponding list of oracles to create distinct requests to each one.
 */
contract PreCoordinator is ChainlinkClient, Ownable, ChainlinkRequestInterface, LinkTokenReceiver {
  using SafeMathChainlink for uint256;

  uint256 constant private MAX_ORACLE_COUNT = 45;

  uint256 private globalNonce;

  struct ServiceAgreement {
    uint256 totalPayment;
    uint256 minResponses;
    address[] oracles;
    bytes32[] jobIds;
    uint256[] payments;
  }

  struct Requester {
    bytes4 callbackFunctionId;
    address sender;
    address callbackAddress;
    int256[] responses;
  }

  // Service Agreement ID => ServiceAgreement
  mapping(bytes32 => ServiceAgreement) internal serviceAgreements;
  // Local Request ID => Service Agreement ID
  mapping(bytes32 => bytes32) internal serviceAgreementRequests;
  // Requester's Request ID => Requester
  mapping(bytes32 => Requester) internal requesters;
  // Local Request ID => Requester's Request ID
  mapping(bytes32 => bytes32) internal requests;

  event NewServiceAgreement(bytes32 indexed saId, uint256 payment, uint256 minresponses);
  event ServiceAgreementRequested(bytes32 indexed saId, bytes32 indexed requestId, uint256 payment);
  event ServiceAgreementResponseReceived(bytes32 indexed saId, bytes32 indexed requestId, address indexed oracle, int256 answer);
  event ServiceAgreementAnswerUpdated(bytes32 indexed saId, bytes32 indexed requestId, int256 answer);
  event ServiceAgreementDeleted(bytes32 indexed saId);

  /**
   * @notice Deploy the contract with a specified address for the LINK
   * and Oracle contract addresses
   * @dev Sets the storage for the specified addresses
   * @param _link The address of the LINK token contract
   */
  constructor(address _link) public {
    if(_link == address(0)) {
      setPublicChainlinkToken();
    } else {
      setChainlinkToken(_link);
    }
  }

  /**
   * @notice Allows the owner of the contract to create new service agreements
   * with multiple oracles. Each oracle will have their own Job ID and can have
   * their own payment amount.
   * @dev The globalNonce keeps service agreement IDs unique. Assume one cannot
   * create the max uint256 number of service agreements in the same block.
   * @param _minResponses The minimum number of responses before the requesting
   * contract is called with the response data.
   * @param _oracles The list of oracle contract addresses.
   * @param _jobIds The corresponding list of Job IDs.
   * @param _payments The corresponding list of payment amounts.
   */
  function createServiceAgreement(
    uint256 _minResponses,
    address[] calldata _oracles,
    bytes32[] calldata _jobIds,
    uint256[] calldata _payments
  )
    external returns (bytes32 saId)
  {
    require(_minResponses > 0, "Min responses must be > 0");
    require(_oracles.length == _jobIds.length && _oracles.length == _payments.length, "Unmet length");
    require(_oracles.length <= MAX_ORACLE_COUNT, "Cannot have more than 45 oracles");
    require(_oracles.length >= _minResponses, "Invalid min responses");
    uint256 totalPayment;
    for (uint i = 0; i < _payments.length; i++) {
      totalPayment = totalPayment.add(_payments[i]);
    }
    saId = keccak256(abi.encodePacked(globalNonce, now));
    globalNonce++; // yes, let it overflow
    serviceAgreements[saId] = ServiceAgreement(totalPayment, _minResponses, _oracles, _jobIds, _payments);

    emit NewServiceAgreement(saId, totalPayment, _minResponses);
  }

  /**
   * @notice This is a helper function to retrieve the details of a service agreement
   * by its given service agreement ID.
   * @dev This function is used instead of the public mapping to return the values
   * of the arrays: oracles, jobIds, and payments.
   */
  function getServiceAgreement(bytes32 _saId)
    external view returns
  (
    uint256 totalPayment,
    uint256 minResponses,
    address[] memory oracles,
    bytes32[] memory jobIds,
    uint256[] memory payments
  )
  {
    return
    (
      serviceAgreements[_saId].totalPayment,
      serviceAgreements[_saId].minResponses,
      serviceAgreements[_saId].oracles,
      serviceAgreements[_saId].jobIds,
      serviceAgreements[_saId].payments
    );
  }

  /**
   * @notice Returns the address of the LINK token
   * @dev This is the public implementation for chainlinkTokenAddress, which is
   * an internal method of the ChainlinkClient contract
   */
  function getChainlinkToken() public view override returns (address) {
    return chainlinkTokenAddress();
  }

  /**
   * @notice Creates the Chainlink request
   * @dev Stores the hash of the params as the on-chain commitment for the request.
   * Emits OracleRequest event for the Chainlink node to detect.
   * @param _sender The sender of the request
   * @param _payment The amount of payment given (specified in wei)
   * @param _saId The Job Specification ID
   * @param _callbackAddress The callback address for the response
   * @param _callbackFunctionId The callback function ID for the response
   * @param _nonce The nonce sent by the requester
   * @param _data The CBOR payload of the request
   */
  function oracleRequest(
    address _sender,
    uint256 _payment,
    bytes32 _saId,
    address _callbackAddress,
    bytes4 _callbackFunctionId,
    uint256 _nonce,
    uint256,
    bytes calldata _data
  )
    external
    onlyLINK
    override
    checkCallbackAddress(_callbackAddress)
  {
    uint256 totalPayment = serviceAgreements[_saId].totalPayment;
    // this revert message does not bubble up
    require(_payment >= totalPayment, "Insufficient payment");
    bytes32 callbackRequestId = keccak256(abi.encodePacked(_sender, _nonce));
    require(requesters[callbackRequestId].sender == address(0), "Nonce already in-use");
    requesters[callbackRequestId].callbackFunctionId = _callbackFunctionId;
    requesters[callbackRequestId].callbackAddress = _callbackAddress;
    requesters[callbackRequestId].sender = _sender;
    createRequests(_saId, callbackRequestId, _data);
    if (_payment > totalPayment) {
      uint256 overage = _payment.sub(totalPayment);
      LinkTokenInterface _link = LinkTokenInterface(chainlinkTokenAddress());
      assert(_link.transfer(_sender, overage));
    }
  }

  /**
   * @dev Creates Chainlink requests to each oracle in the service agreement with the
   * same data payload supplied by the requester
   * @param _saId The service agreement ID
   * @param _incomingRequestId The requester-supplied request ID
   * @param _data The data payload (request parameters) to send to each oracle
   */
  function createRequests(bytes32 _saId, bytes32 _incomingRequestId, bytes memory _data) private {
    ServiceAgreement memory sa = serviceAgreements[_saId];
    require(sa.minResponses > 0, "Invalid service agreement");
    Chainlink.Request memory request;
    bytes32 outgoingRequestId;
    emit ServiceAgreementRequested(_saId, _incomingRequestId, sa.totalPayment);
    for (uint i = 0; i < sa.oracles.length; i++) {
      request = buildChainlinkRequest(sa.jobIds[i], address(this), this.chainlinkCallback.selector);
      request.setBuffer(_data);
      outgoingRequestId = sendChainlinkRequestTo(sa.oracles[i], request, sa.payments[i]);
      requests[outgoingRequestId] = _incomingRequestId;
      serviceAgreementRequests[outgoingRequestId] = _saId;
    }
  }

  /**
   * @notice The fulfill method from requests created by this contract
   * @dev The recordChainlinkFulfillment protects this function from being called
   * by anyone other than the oracle address that the request was sent to
   * @param _requestId The ID that was generated for the request
   * @param _data The answer provided by the oracle
   */
  function chainlinkCallback(bytes32 _requestId, int256 _data)
    external
    recordChainlinkFulfillment(_requestId)
    returns (bool)
  {
    ServiceAgreement memory sa = serviceAgreements[serviceAgreementRequests[_requestId]];
    bytes32 cbRequestId = requests[_requestId];
    bytes32 saId = serviceAgreementRequests[_requestId];
    delete requests[_requestId];
    delete serviceAgreementRequests[_requestId];
    emit ServiceAgreementResponseReceived(saId, cbRequestId, msg.sender, _data);
    requesters[cbRequestId].responses.push(_data);
    Requester memory req = requesters[cbRequestId];
    if (req.responses.length == sa.oracles.length) delete requesters[cbRequestId];
    bool success = true;
    if (req.responses.length == sa.minResponses) {
      int256 result = Median.calculate(req.responses);
      emit ServiceAgreementAnswerUpdated(saId, cbRequestId, result);
      // solhint-disable-next-line avoid-low-level-calls
      (success, ) = req.callbackAddress.call(abi.encodeWithSelector(req.callbackFunctionId, cbRequestId, result));
    }
    return success;
  }

  /**
   * @notice Allows the owner to withdraw any LINK balance on the contract
   * @dev The only valid case for there to be remaining LINK on this contract
   * is if a user accidentally sent LINK directly to this contract's address.
   */
  function withdrawLink() external onlyOwner {
    LinkTokenInterface _link = LinkTokenInterface(chainlinkTokenAddress());
    require(_link.transfer(msg.sender, _link.balanceOf(address(this))), "Unable to transfer");
  }

  /**
   * @notice Call this method if no response is received within 5 minutes
   * @param _requestId The ID that was generated for the request to cancel
   * @param _payment The payment specified for the request to cancel
   * @param _callbackFunctionId The bytes4 callback function ID specified for
   * the request to cancel
   * @param _expiration The expiration generated for the request to cancel
   */
  function cancelOracleRequest(
    bytes32 _requestId,
    uint256 _payment,
    bytes4 _callbackFunctionId,
    uint256 _expiration
  )
    external
    override
  {
    bytes32 cbRequestId = requests[_requestId];
    delete requests[_requestId];
    delete serviceAgreementRequests[_requestId];
    Requester memory req = requesters[cbRequestId];
    require(req.sender == msg.sender, "Only requester can cancel");
    delete requesters[cbRequestId];
    cancelChainlinkRequest(_requestId, _payment, _callbackFunctionId, _expiration);
    LinkTokenInterface _link = LinkTokenInterface(chainlinkTokenAddress());
    require(_link.transfer(req.sender, _payment), "Unable to transfer");
  }

  /**
   * @dev Reverts if the callback address is the LINK token
   * @param _to The callback address
   */
  modifier checkCallbackAddress(address _to) {
    require(_to != chainlinkTokenAddress(), "Cannot callback to LINK");
    _;
  }
}
