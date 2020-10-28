pragma solidity 0.7.0;

import "./LinkTokenReceiver.sol";
import "./Owned.sol";
import "../interfaces/ChainlinkRequestInterface.sol";
import "../interfaces/OracleInterface.sol";
import "../interfaces/OracleInterface2.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/WithdrawalInterface.sol";
import "../vendor/SafeMathChainlink.sol";

/**
 * @title The Chainlink Operator contract
 * @notice Node operators can deploy this contract to fulfill requests sent to them
 */
contract Operator is
  LinkTokenReceiver,
  Owned,
  ChainlinkRequestInterface,
  OracleInterface,
  OracleInterface2,
  WithdrawalInterface
{
  using SafeMathChainlink for uint256;

  uint256 constant public EXPIRY_TIME = 5 minutes;
  uint256 constant private MINIMUM_CONSUMER_GAS_LIMIT = 400000;
  // We initialize fields to 1 instead of 0 so that the first invocation
  // does not cost more gas.
  uint256 constant private ONE_FOR_CONSISTENT_GAS_COST = 1;

  LinkTokenInterface internal immutable linkToken;
  mapping(bytes32 => bytes32) private s_commitments;
  mapping(address => bool) private s_authorizedNodes;
  uint256 private s_withdrawableTokens = ONE_FOR_CONSISTENT_GAS_COST;

  event OracleRequest(
    bytes32 indexed specId,
    address requester,
    bytes32 requestId,
    uint256 payment,
    address callbackAddr,
    bytes4 callbackFunctionId,
    uint256 cancelExpiration,
    uint256 dataVersion,
    bytes data
  );

  event CancelOracleRequest(
    bytes32 indexed requestId
  );

  event OracleResponse(
    bytes32 indexed requestId,
    uint256 dataVersion
  );

  /**
   * @notice Deploy with the address of the LINK token
   * @dev Sets the LinkToken address for the imported LinkTokenInterface
   * @param link The address of the LINK token
   */
  constructor(address link)
    Owned()
  {
    linkToken = LinkTokenInterface(link); // external but already deployed and unalterable
  }

  // EXTERNAL FUNCTIONS

  /**
   * @notice Creates the Chainlink request
   * @dev Stores the hash of the params as the on-chain commitment for the request.
   * Emits OracleRequest event for the Chainlink node to detect.
   * @param sender The sender of the request
   * @param payment The amount of payment given (specified in wei)
   * @param specId The Job Specification ID
   * @param callbackAddress The callback address for the response
   * @param callbackFunctionId The callback function ID for the response
   * @param nonce The nonce sent by the requester
   * @param dataVersion The specified data version
   * @param data The CBOR payload of the request
   */
  function oracleRequest(
    address sender,
    uint256 payment,
    bytes32 specId,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion,
    bytes calldata data
  )
    external
    override
    onlyLINK()
    checkCallbackAddress(callbackAddress)
  {
    (bytes32 requestId, uint256 expiration) = verifyOracleRequest(
      sender,
      payment,
      callbackAddress,
      callbackFunctionId,
      nonce
    );
    emit OracleRequest(
      specId,
      sender,
      requestId,
      payment,
      callbackAddress,
      callbackFunctionId,
      expiration,
      dataVersion,
      data);
  }

  /**
   * @notice Called by the Chainlink node to fulfill requests
   * @dev Given params must hash back to the commitment stored from `oracleRequest`.
   * Will call the callback address' callback function without bubbling up error
   * checking in a `require` so that the node can get paid.
   * @param requestId The fulfillment request ID that must match the requester's
   * @param payment The payment amount that will be released for the oracle (specified in wei)
   * @param callbackAddress The callback address to call for fulfillment
   * @param callbackFunctionId The callback function ID to use for fulfillment
   * @param expiration The expiration that the node should respond by before the requester can cancel
   * @param data The data to return to the consuming contract
   * @return Status if the external call was successful
   */
  function fulfillOracleRequest(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    bytes32 data
  )
    external
    override
    onlyAuthorizedNode()
    isValidRequest(requestId)
    returns (bool)
  {
    verifyOracleResponse(
      requestId,
      payment,
      callbackAddress,
      callbackFunctionId,
      expiration
    );
    emit OracleResponse(requestId, 1);
    require(gasleft() >= MINIMUM_CONSUMER_GAS_LIMIT, "Must provide consumer enough gas");
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    (bool success, ) = callbackAddress.call(abi.encodeWithSelector(callbackFunctionId, requestId, data)); // solhint-disable-line avoid-low-level-calls
    return success;
  }

  /**
   * @notice Called by the Chainlink node to fulfill requests with multi-word support
   * @dev Given params must hash back to the commitment stored from `oracleRequest`.
   * Will call the callback address' callback function without bubbling up error
   * checking in a `require` so that the node can get paid.
   * @param requestId The fulfillment request ID that must match the requester's
   * @param payment The payment amount that will be released for the oracle (specified in wei)
   * @param callbackAddress The callback address to call for fulfillment
   * @param callbackFunctionId The callback function ID to use for fulfillment
   * @param expiration The expiration that the node should respond by before the requester can cancel
   * @param data The data to return to the consuming contract
   * @return Status if the external call was successful
   */
  function fulfillOracleRequest2(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    bytes memory data
  )
    external
    override
    onlyAuthorizedNode()
    isValidRequest(requestId)
    isValidMultiWord(requestId, data)
    returns (bool)
  {
    verifyOracleResponse(
      requestId,
      payment,
      callbackAddress,
      callbackFunctionId,
      expiration
    );
    emit OracleResponse(requestId, 2);
    require(gasleft() >= MINIMUM_CONSUMER_GAS_LIMIT, "Must provide consumer enough gas");
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    (bool success, ) = callbackAddress.call(abi.encodePacked(callbackFunctionId, data)); // solhint-disable-line avoid-low-level-calls
    return success;
  }

  /**
   * @notice Use this to check if a node is authorized for fulfilling requests
   * @param node The address of the Chainlink node
   * @return The authorization status of the node
   */
  function getAuthorizationStatus(address node)
    external
    view
    override
    returns (bool)
  {
    return s_authorizedNodes[node];
  }

  /**
   * @notice Sets the fulfillment permission for a given node. Use `true` to allow, `false` to disallow.
   * @param node The address of the Chainlink node
   * @param allowed Bool value to determine if the node can fulfill requests
   */
  function setFulfillmentPermission(address node, bool allowed)
    external
    override
    onlyOwner()
  {
    s_authorizedNodes[node] = allowed;
  }

  /**
   * @notice Allows the node operator to withdraw earned LINK to a given address
   * @dev The owner of the contract can be another wallet and does not have to be a Chainlink node
   * @param recipient The address to send the LINK token to
   * @param amount The amount to send (specified in wei)
   */
  function withdraw(address recipient, uint256 amount)
    external
    override(OracleInterface, WithdrawalInterface)
    onlyOwner()
    hasAvailableFunds(amount)
  {
    s_withdrawableTokens = s_withdrawableTokens.sub(amount);
    assert(linkToken.transfer(recipient, amount));
  }

  /**
   * @notice Displays the amount of LINK that is available for the node operator to withdraw
   * @dev We use `ONE_FOR_CONSISTENT_GAS_COST` in place of 0 in storage
   * @return The amount of withdrawable LINK on the contract
   */
  function withdrawable()
    external
    view
    override(OracleInterface, WithdrawalInterface)
    returns (uint256)
  {
    return s_withdrawableTokens.sub(ONE_FOR_CONSISTENT_GAS_COST);
  }

  /**
   * @notice Allows requesters to cancel requests sent to this oracle contract. Will transfer the LINK
   * sent for the request back to the requester's address.
   * @dev Given params must hash to a commitment stored on the contract in order for the request to be valid
   * Emits CancelOracleRequest event.
   * @param requestId The request ID
   * @param payment The amount of payment given (specified in wei)
   * @param callbackFunc The requester's specified callback address
   * @param expiration The time of the expiration for the request
   */
  function cancelOracleRequest(
    bytes32 requestId,
    uint256 payment,
    bytes4 callbackFunc,
    uint256 expiration
  )
    external
    override
  {
    bytes32 paramsHash = keccak256(
      abi.encodePacked(
        payment,
        msg.sender,
        callbackFunc,
        expiration)
    );
    require(paramsHash == s_commitments[requestId], "Params do not match request ID");
    // solhint-disable-next-line not-rely-on-time
    require(expiration <= block.timestamp, "Request is not expired");

    delete s_commitments[requestId];
    emit CancelOracleRequest(requestId);

    assert(linkToken.transfer(msg.sender, payment));
  }

  /**
   * @notice Returns the address of the LINK token
   * @dev This is the public implementation for chainlinkTokenAddress, which is
   * an internal method of the ChainlinkClient contract
   */
  function getChainlinkToken()
    public
    view
    override
    returns (address)
  {
    return address(linkToken);
  }

  function forward(address _to, bytes calldata _data)
    public
    onlyAuthorizedNode()
  {
    require(_to != address(linkToken), "Cannot use #forward to send messages to Link token");
    (bool status,) = _to.call(_data);
    require(status, "Forwarded call failed.");
  }

  // INTERNAL FUNCTIONS

  /**
   * @notice Verify the Oracle Request
   * @param sender The sender of the request
   * @param payment The amount of payment given (specified in wei)
   * @param callbackAddress The callback address for the response
   * @param callbackFunctionId The callback function ID for the response
   * @param nonce The nonce sent by the requester
   */
  function verifyOracleRequest(
    address sender,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce
  ) internal returns (bytes32 requestId, uint256 expiration) {
    requestId = keccak256(abi.encodePacked(sender, nonce));
    require(s_commitments[requestId] == 0, "Must use a unique ID");
    // solhint-disable-next-line not-rely-on-time
    expiration = block.timestamp.add(EXPIRY_TIME);
    s_commitments[requestId] = keccak256(
      abi.encodePacked(
        payment,
        callbackAddress,
        callbackFunctionId,
        expiration
      )
    );
    return (requestId, expiration);
  }

  /**
   * @notice Verify the Oracle Response
   * @param requestId The fulfillment request ID that must match the requester's
   * @param payment The payment amount that will be released for the oracle (specified in wei)
   * @param callbackAddress The callback address to call for fulfillment
   * @param callbackFunctionId The callback function ID to use for fulfillment
   * @param expiration The expiration that the node should respond by before the requester can cancel
   */
  function verifyOracleResponse(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration
  )
  internal
  {
    bytes32 paramsHash = keccak256(
      abi.encodePacked(
        payment,
        callbackAddress,
        callbackFunctionId,
        expiration
      )
    );
    require(s_commitments[requestId] == paramsHash, "Params do not match request ID");
    s_withdrawableTokens = s_withdrawableTokens.add(payment);
    delete s_commitments[requestId];
  }

  // MODIFIERS

  /**
   * @dev Reverts if the first 32 bytes of the bytes array is not equal to requestId
   * @param requestId bytes32
   * @param data bytes
   */
  modifier isValidMultiWord(bytes32 requestId, bytes memory data) {
    bytes32 firstWord;
    assembly{
      firstWord := mload(add(data, 0x20))
    }
    require(requestId == firstWord, "First word must be requestId");
    _;
  }


  /**
   * @dev Reverts if amount requested is greater than withdrawable balance
   * @param amount The given amount to compare to `s_withdrawableTokens`
   */
  modifier hasAvailableFunds(uint256 amount) {
    require(s_withdrawableTokens >= amount.add(ONE_FOR_CONSISTENT_GAS_COST), "Amount requested is greater than withdrawable balance");
    _;
  }

  /**
   * @dev Reverts if request ID does not exist
   * @param requestId The given request ID to check in stored `commitments`
   */
  modifier isValidRequest(bytes32 requestId) {
    require(s_commitments[requestId] != 0, "Must have a valid requestId");
    _;
  }

  /**
   * @dev Reverts if `msg.sender` is not authorized to fulfill requests
   */
  modifier onlyAuthorizedNode() {
    require(s_authorizedNodes[msg.sender], "Not an authorized node to fulfill requests");
    _;
  }

  /**
   * @dev Reverts if the callback address is the LINK token
   * @param to The callback address
   */
  modifier checkCallbackAddress(address to) {
    require(to != address(linkToken), "Cannot callback to LINK");
    _;
  }

}
