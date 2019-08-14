pragma solidity 0.5.0;

import "./ChainlinkClient.sol";
import { Ownable as Ownable_Chainlink } from "./vendor/Ownable.sol";
import { SafeMath as SafeMath_Chainlink } from "./vendor/SafeMath.sol";
import { SignedSafeMath as SignedSafeMath_Chainlink } from "./vendor/SignedSafeMath.sol";

/**
 * @title PreCoordinator is a contract that builds on-chain service agreements
 * using the current architecture of 1 request to 1 oracle contract.
 * @dev This contract accepts requests as service agreement IDs and loops over
 * the corresponding list of oracles to create distinct requests to each one.
 */
contract PreCoordinator is ChainlinkClient, Ownable_Chainlink, ChainlinkRequestInterface {
  using SafeMath_Chainlink for uint256;
  using SignedSafeMath_Chainlink for int256;

  uint256 constant private SELECTOR_LENGTH = 4;
  uint256 constant private EXPECTED_REQUEST_WORDS = 2;
  // solium-disable-next-line zeppelin/no-arithmetic-operations
  uint256 constant private MINIMUM_REQUEST_LENGTH = SELECTOR_LENGTH + (32 * EXPECTED_REQUEST_WORDS);
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
  event ServiceAgreementRequest(bytes32 saId, uint256 payment);

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
   * @param _totalPayment The sum of the _payments array. Compute this off-chain.
   * @param _minResponses The minimum number of responses before the requesting
   * contract is called with the response data.
   * @param _oracles The list of oracle contract addresses.
   * @param _jobIds The corresponding list of Job IDs.
   * @param _payments The corresponding list of payment amounts.
   */
  function createServiceAgreement(
    uint256 _totalPayment,
    uint256 _minResponses,
    address[] calldata _oracles,
    bytes32[] calldata _jobIds,
    uint256[] calldata _payments
  )
    external onlyOwner returns (bytes32 saId)
  {
    require(_oracles.length == _jobIds.length && _oracles.length == _payments.length, "Unmet length");
    require(_oracles.length <= MAX_ORACLE_COUNT, "Cannot have more than 45 oracles");
    require(_oracles.length >= _minResponses, "Invalid min responses");
    saId = keccak256(abi.encodePacked(globalNonce, now));
    globalNonce++; // yes, let it overflow
    // Manually calculate total payment off-chain
    serviceAgreements[saId] = ServiceAgreement(_totalPayment, _minResponses, _oracles, _jobIds, _payments);

    emit NewServiceAgreement(saId, _totalPayment, _minResponses);
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
   * @notice Deletes a service agreement from storage
   * @dev Use this with caution since there may be responses waiting to come
   * back for a service agreement. This can be monitored off-chain by looking
   * for the ServiceAgreementRequest event.
   * @param _saId The service agreement ID
   */
  function deleteServiceAgreement(bytes32 _saId) external onlyOwner {
    delete serviceAgreements[_saId];
  }

  /**
   * @notice Returns the address of the LINK token
   * @dev This is the public implementation for chainlinkTokenAddress, which is
   * an internal method of the ChainlinkClient contract
   */
  function getChainlinkToken() external view returns (address) {
    return chainlinkTokenAddress();
  }

  /**
   * @notice Called when LINK is sent to the contract via `transferAndCall`
   * @dev The data payload's first 2 words will be overwritten by the `_sender` and `_amount`
   * values to ensure correctness. Calls oracleRequest.
   * @param _sender Address of the sender
   * @param _amount Amount of LINK sent (specified in wei)
   * @param _data Payload of the transaction
   */
  function onTokenTransfer(
    address _sender,
    uint256 _amount,
    bytes memory _data
  )
    public
    onlyLINK
    validRequestLength(_data)
    permittedFunctionsForLINK(_data)
  {
    assembly {
      // solhint-disable-next-line avoid-low-level-calls
      mstore(add(_data, 36), _sender) // ensure correct sender is passed
      // solhint-disable-next-line avoid-low-level-calls
      mstore(add(_data, 68), _amount)    // ensure correct amount is passed
    }
    // solhint-disable-next-line avoid-low-level-calls
    (bool success, ) = address(this).delegatecall(_data); // calls oracleRequest
    require(success, "Unable to create request");
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
    checkCallbackAddress(_callbackAddress)
  {
    ServiceAgreement memory sa = serviceAgreements[_saId];
    // this revert message does not bubble up
    require(_payment >= sa.totalPayment, "Insufficient payment");
    bytes32 callbackRequestId = keccak256(abi.encodePacked(_sender, _nonce));
    requesters[callbackRequestId].callbackFunctionId = _callbackFunctionId;
    requesters[callbackRequestId].callbackAddress = _sender;
    createRequests(_saId, callbackRequestId, _data);
    if (_payment > sa.totalPayment) {
      uint256 overage = _payment.sub(sa.totalPayment);
      LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
      assert(link.transfer(_sender, overage));
    }
  }

  /**
   * @dev Creates Chainlink requests to each oracle in the service agreement with the
   * same data payload supplied by the requester
   * @param _saId The service agreement ID
   * @param _requestId The requester-supplied request ID
   * @param _data The data payload (request parameters) to send to each oracle
   */
  function createRequests(bytes32 _saId, bytes32 _requestId, bytes memory _data) private {
    ServiceAgreement memory sa = serviceAgreements[_saId];
    Chainlink.Request memory request;
    bytes32 requestId;
    emit ServiceAgreementRequest(_saId, sa.minResponses);
    for (uint i = 0; i < sa.oracles.length; i++) {
      request = buildChainlinkRequest(sa.jobIds[i], address(this), this.chainlinkCallback.selector);
      request.setBuffer(_data);
      requestId = sendChainlinkRequestTo(sa.oracles[i], request, sa.payments[i]);
      requests[requestId] = _requestId;
      serviceAgreementRequests[requestId] = _saId;
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
    delete requests[_requestId];
    delete serviceAgreementRequests[_requestId];
    if (requesters[cbRequestId].responses.push(_data) == sa.minResponses) {
      Requester memory req = requesters[cbRequestId];
      delete requesters[cbRequestId];
      int256 result = getMedian(req.responses);
      // solhint-disable-next-line avoid-low-level-calls
      (bool success, ) = req.callbackAddress.call(abi.encodeWithSelector(req.callbackFunctionId, cbRequestId, result));
      return success;
    }
    return true;
  }

  /**
   * @dev Performs aggregation of the answers received from the Chainlink nodes.
   * Assumes that at least half the oracles are honest and so can't contol the
   * middle of the ordered responses.
   * @param _responses The answer ID associated with the group of requests
   */
  function getMedian(int256[] memory _responses)
    private pure returns (int256 result)
  {
    uint256 responseLength = _responses.length;
    uint256 middleIndex = responseLength.div(2);
    if (responseLength % 2 == 0) {
      int256 median1 = quickselect(_responses, middleIndex);
      int256 median2 = quickselect(_responses, middleIndex.add(1)); // quickselect is 1 indexed
      result = median1.add(median2) / 2; // signed integers are not supported by SafeMath
    } else {
      result = quickselect(_responses, middleIndex.add(1)); // quickselect is 1 indexed
    }
  }

  /**
   * @dev Returns the kth value of the ordered array
   * See: http://www.cs.yale.edu/homes/aspnes/pinewiki/QuickSelect.html
   * @param _a The list of elements to pull from
   * @param _k The index, 1 based, of the elements you want to pull from when ordered
   */
  function quickselect(int256[] memory _a, uint256 _k)
    private
    pure
    returns (int256)
  {
    int256[] memory a = _a;
    uint256 k = _k;
    uint256 aLen = a.length;
    int256[] memory a1 = new int256[](aLen);
    int256[] memory a2 = new int256[](aLen);
    uint256 a1Len;
    uint256 a2Len;
    int256 pivot;
    uint256 i;

    while (true) {
      pivot = a[aLen.div(2)];
      a1Len = 0;
      a2Len = 0;
      for (i = 0; i < aLen; i++) {
        if (a[i] < pivot) {
          a1[a1Len] = a[i];
          a1Len++;
        } else if (a[i] > pivot) {
          a2[a2Len] = a[i];
          a2Len++;
        }
      }
      if (k <= a1Len) {
        aLen = a1Len;
        (a, a1) = swap(a, a1);
      } else if (k > (aLen.sub(a2Len))) {
        k = k.sub(aLen.sub(a2Len));
        aLen = a2Len;
        (a, a2) = swap(a, a2);
      } else {
        return pivot;
      }
    }
  }

  /**
   * @dev Swaps the pointers to two uint256 arrays in memory
   * @param _a The pointer to the first in memory array
   * @param _b The pointer to the second in memory array
   */
  function swap(int256[] memory _a, int256[] memory _b)
    private
    pure
    returns(int256[] memory, int256[] memory)
  {
    return (_b, _a);
  }

  /**
   * @notice Allows the owner to withdraw any LINK balance on the contract
   * @dev The only valid case for there to be remaining LINK on this contract
   * is if a user accidentally sent LINK directly to this contract's address.
   */
  function withdrawLink() external onlyOwner {
    LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
    require(link.transfer(msg.sender, link.balanceOf(address(this))), "Unable to transfer");
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
  {
    bytes32 cbRequestId = requests[_requestId];
    delete requests[_requestId];
    delete serviceAgreementRequests[_requestId];
    Requester memory req = requesters[cbRequestId];
    delete requesters[cbRequestId];
    cancelChainlinkRequest(_requestId, _payment, _callbackFunctionId, _expiration);
    LinkTokenInterface link = LinkTokenInterface(chainlinkTokenAddress());
    require(link.transfer(req.callbackAddress, _payment), "Unable to transfer");
  }

  /**
   * @dev Reverts if not sent from the LINK token
   */
  modifier onlyLINK() {
    require(msg.sender == chainlinkTokenAddress(), "Must use LINK token");
    _;
  }

  /**
   * @dev Reverts if the given data does not begin with the `oracleRequest` function selector
   * @param _data The data payload of the request
   */
  modifier permittedFunctionsForLINK(bytes memory _data) {
    bytes4 funcSelector;
    assembly {
      // solhint-disable-next-line avoid-low-level-calls
      funcSelector := mload(add(_data, 32))
    }
    require(funcSelector == this.oracleRequest.selector, "Must use whitelisted functions");
    _;
  }

  /**
   * @dev Reverts if the callback address is the LINK token
   * @param _to The callback address
   */
  modifier checkCallbackAddress(address _to) {
    require(_to != chainlinkTokenAddress(), "Cannot callback to LINK");
    _;
  }

  /**
   * @dev Reverts if the given payload is less than needed to create a request
   * @param _data The request payload
   */
  modifier validRequestLength(bytes memory _data) {
    require(_data.length >= MINIMUM_REQUEST_LENGTH, "Invalid request length");
    _;
  }
}