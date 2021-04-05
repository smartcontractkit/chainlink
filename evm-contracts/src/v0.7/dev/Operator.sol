// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "./LinkTokenReceiver.sol";
import "./ConfirmedOwner.sol";
import "./OperatorForwarder.sol";
import "../interfaces/ChainlinkRequestInterface.sol";
import "../interfaces/OperatorInterface.sol";
import "../interfaces/LinkTokenInterface.sol";
import "../interfaces/WithdrawalInterface.sol";
import "../vendor/SafeMathChainlink.sol";

/**
 * @title The Chainlink Operator contract
 * @notice Node operators can deploy this contract to fulfill requests sent to them
 */
contract Operator is
  LinkTokenReceiver,
  ConfirmedOwner,
  ChainlinkRequestInterface,
  OperatorInterface,
  WithdrawalInterface
{
  using SafeMathChainlink for uint256;

  struct Commitment {
    bytes31 paramsHash;
    uint8 dataVersion;
  }

  uint256 constant public EXPIRY_TIME = 5 minutes;
  uint256 constant private MAXIMUM_DATA_VERSION = 256;
  uint256 constant private MINIMUM_CONSUMER_GAS_LIMIT = 400000;
  // We initialize fields to 1 instead of 0 so that the first invocation
  // does not cost more gas.
  uint256 constant private ONE_FOR_CONSISTENT_GAS_COST = 1;

  LinkTokenInterface internal immutable linkToken;
  mapping(bytes32 => Commitment) private s_commitments;
  mapping(address => bool) private s_authorizedSenders;
  address[] private s_authorizedSenderList;
  // Tokens sent for requests that have not been fulfilled yet
  uint256 private s_tokensInEscrow = ONE_FOR_CONSISTENT_GAS_COST;
  // Forwarders
  address[] private s_forwardersList;

  event AuthorizedSendersChanged(
    address[] senders
  );

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
    bytes32 indexed requestId
  );

  event ForwarderCreated(
    address indexed addr
  );

  /**
   * @notice Deploy with the address of the LINK token
   * @dev Sets the LinkToken address for the imported LinkTokenInterface
   * @param link The address of the LINK token
   * @param owner The address of the owner
   */
  constructor(
    address link,
    address owner
  )
    ConfirmedOwner(owner)
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
    (bytes32 requestId, uint256 expiration) = _verifyOracleRequest(
      sender,
      payment,
      callbackAddress,
      callbackFunctionId,
      nonce,
      dataVersion
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
    onlyAuthorizedSender()
    isValidRequest(requestId)
    returns (
      bool
    )
  {
    _verifyOracleResponse(
      requestId,
      payment,
      callbackAddress,
      callbackFunctionId,
      expiration,
      1
    );
    emit OracleResponse(requestId);
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
    bytes calldata data
  )
    external
    override
    onlyAuthorizedSender()
    isValidRequest(requestId)
    isValidMultiWord(requestId, data)
    returns (
      bool
    )
  {
    _verifyOracleResponse(
      requestId,
      payment,
      callbackAddress,
      callbackFunctionId,
      expiration,
      2
    );
    emit OracleResponse(requestId);
    require(gasleft() >= MINIMUM_CONSUMER_GAS_LIMIT, "Must provide consumer enough gas");
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    (bool success, ) = callbackAddress.call(abi.encodePacked(callbackFunctionId, data)); // solhint-disable-line avoid-low-level-calls
    return success;
  }

  /**
   * @notice Use this to check if a node is authorized for fulfilling requests
   * @param sender The address of the Chainlink node
   * @return The authorization status of the node
   */
  function isAuthorizedSender(
    address sender
  )
    external
    view
    override
    returns (bool)
  {
    return s_authorizedSenders[sender];
  }

  /**
   * @notice Sets the fulfillment permission for a given node. Use `true` to allow, `false` to disallow.
   * @param senders The addresses of the authorized Chainlink node
   */
  function setAuthorizedSenders(
    address[] calldata senders
  )
    external
    override
    onlyOwner()
  {
    require(senders.length > 0, "Must have at least 1 authorized sender");
    // Set previous authorized senders to false
    uint256 authorizedSendersLength = s_authorizedSenderList.length;
    for (uint256 i = 0; i < authorizedSendersLength; i++) {
      s_authorizedSenders[s_authorizedSenderList[i]] = false;
    }
    // Set new to true
    for (uint256 i = 0; i < senders.length; i++) {
      s_authorizedSenders[senders[i]] = true;
    }
    // Replace list
    s_authorizedSenderList = senders;
    emit AuthorizedSendersChanged(senders);
  }

  /**
   * @notice Retrieve a list of authorized senders
   * @return array of addresses
   */
  function getAuthorizedSenders()
    external
    view
    override
    returns (
      address[] memory
    )
  {
    return s_authorizedSenderList;
  }

  /**
   * @notice Retrive a list of forwarders
   * @return array of addresses
   */
  function getForwarders()
    external
    view
    override
    returns (
      address[] memory
    )
  {
    return s_forwardersList;
  }

  /**
   * @notice Allows the node operator to withdraw earned LINK to a given address
   * @dev The owner of the contract can be another wallet and does not have to be a Chainlink node
   * @param recipient The address to send the LINK token to
   * @param amount The amount to send (specified in wei)
   */
  function withdraw(
    address recipient,
    uint256 amount
  )
    external
    override(OracleInterface, WithdrawalInterface)
    onlyOwner()
    hasAvailableFunds(amount)
  {
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
    return _fundsAvailable();
  }

  /**
   * @notice Interact with other LinkTokenReceiver contracts by calling transferAndCall
   * @param to The address to transfer to.
   * @param value The amount to be transferred.
   * @param data The extra data to be passed to the receiving contract.
   * @return success bool
   */
  function operatorTransferAndCall(
    address to,
    uint256 value,
    bytes calldata data
  )
    external
    override
    onlyOwner()
    hasAvailableFunds(value)
    returns (
      bool success
    )
  {
    return linkToken.transferAndCall(to, value, data);
  }

  /**
   * @notice Distribute funds to multiple addresses using ETH send
   * to this payable function.
   * @dev Array length must be equal, ETH sent must equal the sum of amounts.
   * @param receivers list of addresses
   * @param amounts list of amounts
   */
  function distributeFunds(
    address payable[] calldata receivers,
    uint[] calldata amounts
  )
    external
    override
    payable
  {
    require(receivers.length > 0 && receivers.length == amounts.length, "Invalid array length(s)");
    uint256 valueRemaining = msg.value;
    for (uint256 i = 0; i < receivers.length; i++) {
      uint256 sendAmount = amounts[i];
      valueRemaining = valueRemaining.sub(sendAmount);
      receivers[i].transfer(sendAmount);
    }
    require(valueRemaining == 0, "Too much ETH sent");
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
    bytes31 paramsHash = _buildFunctionHash(payment, msg.sender, callbackFunc, expiration);
    require(s_commitments[requestId].paramsHash == paramsHash, "Params do not match request ID");
    // solhint-disable-next-line not-rely-on-time
    require(expiration <= block.timestamp, "Request is not expired");

    delete s_commitments[requestId];
    emit CancelOracleRequest(requestId);

    assert(linkToken.transfer(msg.sender, payment));
  }

  // PUBLIC FUNCTIONS

  /**
   * @notice Create a new forwarder contract
   * @dev This function uses create2 to deploy at a predetermined address.
   * @return addr deployed address
   */
  function createForwarder()
    public
    onlyAuthorizedSender()
    returns (
      address addr
    )
  {
    bytes32 salt = bytes32(s_forwardersList.length);
    OperatorForwarder forwarder = new OperatorForwarder{salt: salt}(address(linkToken));
    addr = address(forwarder);
    require(addr != address(0), "Create2: Failed deployment");
    s_forwardersList.push(addr);
    emit ForwarderCreated(addr);
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
    returns (
      address
    )
  {
    return address(linkToken);
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
  function _verifyOracleRequest(
    address sender,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion
  ) 
    internal
    returns (
      bytes32 requestId,
      uint256 expiration
    )
  {
    requestId = keccak256(abi.encodePacked(sender, nonce));
    require(s_commitments[requestId].paramsHash == 0, "Must use a unique ID");
    // solhint-disable-next-line not-rely-on-time
    expiration = block.timestamp.add(EXPIRY_TIME);
    bytes31 paramsHash = _buildFunctionHash(payment, callbackAddress, callbackFunctionId, expiration);
    s_commitments[requestId] = Commitment(paramsHash, _safeCastToUint8(dataVersion));
    s_tokensInEscrow = s_tokensInEscrow.add(payment);
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
  function _verifyOracleResponse(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    uint256 dataVersion
  )
    internal
  {
    bytes31 paramsHash = _buildFunctionHash(payment, callbackAddress, callbackFunctionId, expiration);
    require(s_commitments[requestId].paramsHash == paramsHash, "Params do not match request ID");
    require(s_commitments[requestId].dataVersion <= _safeCastToUint8(dataVersion), "Data versions must match");
    s_tokensInEscrow = s_tokensInEscrow.sub(payment);
    delete s_commitments[requestId];
  }

  /**
   * @notice Build the bytes31 function hash from the payment, callback and expiration.
   * @param payment The payment amount that will be released for the oracle (specified in wei)
   * @param callbackAddress The callback address to call for fulfillment
   * @param callbackFunctionId The callback function ID to use for fulfillment
   * @param expiration The expiration that the node should respond by before the requester can cancel
   * @return hash bytes31
   */
  function _buildFunctionHash(
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration
  )
    internal
    pure
    returns (
      bytes31
    )
  {
    return bytes31(keccak256(
      abi.encodePacked(
        payment,
        callbackAddress,
        callbackFunctionId,
        expiration
      )
    ));
  }

  /**
   * @notice Safely cast uint256 to uint8
   * @param number uint256
   * @return uint8 number
   */
  function _safeCastToUint8(
    uint256 number
  )
    internal
    pure
    returns (
      uint8
    )
  {
    require(number < MAXIMUM_DATA_VERSION, "number too big to cast");
    return uint8(number);
  }

  // PRIVATE FUNCTIONS

  /**
   * @notice Returns the LINK available in this contract, not locked in escrow
   * @return uint256 LINK tokens available
   */
  function _fundsAvailable()
    private
    view
    returns (
      uint256
    )
  {
    uint256 inEscrow = s_tokensInEscrow.sub(ONE_FOR_CONSISTENT_GAS_COST);
    return linkToken.balanceOf(address(this)).sub(inEscrow);
  }

  // MODIFIERS

  /**
   * @dev Reverts if the first 32 bytes of the bytes array is not equal to requestId
   * @param requestId bytes32
   * @param data bytes
   */
  modifier isValidMultiWord(
    bytes32 requestId,
    bytes memory data
  ) {
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
  modifier hasAvailableFunds(
    uint256 amount
  ) {
    require(_fundsAvailable() >= amount, "Amount requested is greater than withdrawable balance");
    _;
  }

  /**
   * @dev Reverts if request ID does not exist
   * @param requestId The given request ID to check in stored `commitments`
   */
  modifier isValidRequest(
    bytes32 requestId
  ) {
    require(s_commitments[requestId].paramsHash != 0, "Must have a valid requestId");
    _;
  }

  /**
   * @dev Reverts if `msg.sender` is not authorized to fulfill requests
   */
  modifier onlyAuthorizedSender() {
    require(s_authorizedSenders[msg.sender], "Not an authorized node to fulfill requests");
    _;
  }

  /**
   * @dev Reverts if the callback address is the LINK token
   * @param to The callback address
   */
  modifier checkCallbackAddress(
    address to
  ) {
    require(to != address(linkToken), "Cannot callback to LINK");
    _;
  }

}
