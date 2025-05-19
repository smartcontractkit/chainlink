// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {AuthorizedReceiver} from "./AuthorizedReceiver.sol";
import {LinkTokenReceiver} from "./LinkTokenReceiver.sol";
import {ConfirmedOwner} from "../shared/access/ConfirmedOwner.sol";
import {LinkTokenInterface} from "../shared/interfaces/LinkTokenInterface.sol";
import {IAuthorizedReceiver} from "./interfaces/IAuthorizedReceiver.sol";
import {OperatorInterface} from "./interfaces/OperatorInterface.sol";
import {IOwnable} from "../shared/interfaces/IOwnable.sol";
import {IWithdrawal} from "./interfaces/IWithdrawal.sol";
import {OracleInterface} from "./interfaces/OracleInterface.sol";
import {SafeCast} from "../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/SafeCast.sol";

// @title The Chainlink Operator contract
// @notice Node operators can deploy this contract to fulfill requests sent to them
// solhint-disable gas-custom-errors
contract Operator is AuthorizedReceiver, ConfirmedOwner, LinkTokenReceiver, OperatorInterface, IWithdrawal {
  struct Commitment {
    bytes31 paramsHash;
    uint8 dataVersion;
  }

  uint256 public constant EXPIRYTIME = 5 minutes;
  uint256 private constant MAXIMUM_DATA_VERSION = 256;
  uint256 private constant MINIMUM_CONSUMER_GAS_LIMIT = 400000;
  uint256 private constant SELECTOR_LENGTH = 4;
  uint256 private constant EXPECTED_REQUEST_WORDS = 2;
  uint256 private constant MINIMUM_REQUEST_LENGTH = SELECTOR_LENGTH + (32 * EXPECTED_REQUEST_WORDS);
  // We initialize fields to 1 instead of 0 so that the first invocation
  // does not cost more gas.
  uint256 private constant ONE_FOR_CONSISTENT_GAS_COST = 1;
  // oracleRequest is intended for version 1, enabling single word responses
  bytes4 private constant ORACLE_REQUEST_SELECTOR = this.oracleRequest.selector;
  // operatorRequest is intended for version 2, enabling multi-word responses
  bytes4 private constant OPERATOR_REQUEST_SELECTOR = this.operatorRequest.selector;

  LinkTokenInterface internal immutable i_linkToken;
  mapping(bytes32 => Commitment) private s_commitments;
  mapping(address => bool) private s_owned;
  // Tokens sent for requests that have not been fulfilled yet
  uint256 private s_tokensInEscrow = ONE_FOR_CONSISTENT_GAS_COST;

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

  event CancelOracleRequest(bytes32 indexed requestId);

  event OracleResponse(bytes32 indexed requestId);

  event OwnableContractAccepted(address indexed acceptedContract);

  event TargetsUpdatedAuthorizedSenders(address[] targets, address[] senders, address changedBy);

  // @notice Deploy with the address of the LINK token
  // @dev Sets the LinkToken address for the imported LinkTokenInterface
  // @param link The address of the LINK token
  // @param owner The address of the owner
  constructor(address link, address owner) ConfirmedOwner(owner) {
    i_linkToken = LinkTokenInterface(link); // external but already deployed and unalterable
  }

  string public constant typeAndVersion = "Operator 1.0.0";

  // @notice Creates the Chainlink request. This is a backwards compatible API
  // with the Oracle.sol contract, but the behavior changes because
  // callbackAddress is assumed to be the same as the request sender.
  // @param callbackAddress The consumer of the request
  // @param payment The amount of payment given (specified in wei)
  // @param specId The Job Specification ID
  // @param callbackAddress The address the oracle data will be sent to
  // @param callbackFunctionId The callback function ID for the response
  // @param nonce The nonce sent by the requester
  // @param dataVersion The specified data version
  // @param data The extra request parameters
  function oracleRequest(
    address sender,
    uint256 payment,
    bytes32 specId,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion,
    bytes calldata data
  ) external override validateFromLINK {
    (bytes32 requestId, uint256 expiration) = _verifyAndProcessOracleRequest(
      sender,
      payment,
      callbackAddress,
      callbackFunctionId,
      nonce,
      dataVersion
    );
    emit OracleRequest(specId, sender, requestId, payment, sender, callbackFunctionId, expiration, dataVersion, data);
  }

  // @notice Creates the Chainlink request
  // @dev Stores the hash of the params as the on-chain commitment for the request.
  // Emits OracleRequest event for the Chainlink node to detect.
  // @param sender The sender of the request
  // @param payment The amount of payment given (specified in wei)
  // @param specId The Job Specification ID
  // @param callbackFunctionId The callback function ID for the response
  // @param nonce The nonce sent by the requester
  // @param dataVersion The specified data version
  // @param data The extra request parameters
  function operatorRequest(
    address sender,
    uint256 payment,
    bytes32 specId,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion,
    bytes calldata data
  ) external override validateFromLINK {
    (bytes32 requestId, uint256 expiration) = _verifyAndProcessOracleRequest(
      sender,
      payment,
      sender,
      callbackFunctionId,
      nonce,
      dataVersion
    );
    emit OracleRequest(specId, sender, requestId, payment, sender, callbackFunctionId, expiration, dataVersion, data);
  }

  // @notice Called by the Chainlink node to fulfill requests
  // @dev Given params must hash back to the commitment stored from `oracleRequest`.
  // Will call the callback address' callback function without bubbling up error
  // checking in a `require` so that the node can get paid.
  // @param requestId The fulfillment request ID that must match the requester's
  // @param payment The payment amount that will be released for the oracle (specified in wei)
  // @param callbackAddress The callback address to call for fulfillment
  // @param callbackFunctionId The callback function ID to use for fulfillment
  // @param expiration The expiration that the node should respond by before the requester can cancel
  // @param data The data to return to the consuming contract
  // @return Status if the external call was successful
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
    validateAuthorizedSender
    validateRequestId(requestId)
    validateCallbackAddress(callbackAddress)
    returns (bool)
  {
    _verifyOracleRequestAndProcessPayment(requestId, payment, callbackAddress, callbackFunctionId, expiration, 1);
    emit OracleResponse(requestId);
    require(gasleft() >= MINIMUM_CONSUMER_GAS_LIMIT, "Must provide consumer enough gas");
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    (bool success, ) = callbackAddress.call(abi.encodeWithSelector(callbackFunctionId, requestId, data)); // solhint-disable-line avoid-low-level-calls
    return success;
  }

  // @notice Called by the Chainlink node to fulfill requests with multi-word support
  // @dev Given params must hash back to the commitment stored from `oracleRequest`.
  // Will call the callback address' callback function without bubbling up error
  // checking in a `require` so that the node can get paid.
  // @param requestId The fulfillment request ID that must match the requester's
  // @param payment The payment amount that will be released for the oracle (specified in wei)
  // @param callbackAddress The callback address to call for fulfillment
  // @param callbackFunctionId The callback function ID to use for fulfillment
  // @param expiration The expiration that the node should respond by before the requester can cancel
  // @param data The data to return to the consuming contract
  // @return Status if the external call was successful
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
    validateAuthorizedSender
    validateRequestId(requestId)
    validateCallbackAddress(callbackAddress)
    validateMultiWordResponseId(requestId, data)
    returns (bool)
  {
    _verifyOracleRequestAndProcessPayment(requestId, payment, callbackAddress, callbackFunctionId, expiration, 2);
    emit OracleResponse(requestId);
    require(gasleft() >= MINIMUM_CONSUMER_GAS_LIMIT, "Must provide consumer enough gas");
    // All updates to the oracle's fulfillment should come before calling the
    // callback(addr+functionId) as it is untrusted.
    // See: https://solidity.readthedocs.io/en/develop/security-considerations.html#use-the-checks-effects-interactions-pattern
    (bool success, ) = callbackAddress.call(abi.encodePacked(callbackFunctionId, data)); // solhint-disable-line avoid-low-level-calls
    return success;
  }

  // @notice Transfer the ownership of ownable contracts. This is primarily
  // intended for Authorized Forwarders but could possibly be extended to work
  // with future contracts.OracleInterface
  // @param ownable list of addresses to transfer
  // @param newOwner address to transfer ownership to
  function transferOwnableContracts(address[] calldata ownable, address newOwner) external onlyOwner {
    for (uint256 i = 0; i < ownable.length; ++i) {
      s_owned[ownable[i]] = false;
      IOwnable(ownable[i]).transferOwnership(newOwner);
    }
  }

  // @notice Accept the ownership of an ownable contract. This is primarily
  // intended for Authorized Forwarders but could possibly be extended to work
  // with future contracts.
  // @dev Must be the pending owner on the contract
  // @param ownable list of addresses of Ownable contracts to accept
  function acceptOwnableContracts(address[] calldata ownable) public validateAuthorizedSenderSetter {
    for (uint256 i = 0; i < ownable.length; ++i) {
      s_owned[ownable[i]] = true;
      emit OwnableContractAccepted(ownable[i]);
      IOwnable(ownable[i]).acceptOwnership();
    }
  }

  // @notice Sets the fulfillment permission for
  // @param targets The addresses to set permissions on
  // @param senders The addresses that are allowed to send updates
  function setAuthorizedSendersOn(
    address[] calldata targets,
    address[] calldata senders
  ) public validateAuthorizedSenderSetter {
    emit TargetsUpdatedAuthorizedSenders(targets, senders, msg.sender);

    for (uint256 i = 0; i < targets.length; ++i) {
      IAuthorizedReceiver(targets[i]).setAuthorizedSenders(senders);
    }
  }

  // @notice Accepts ownership of ownable contracts and then immediately sets
  // the authorized sender list on each of the newly owned contracts. This is
  // primarily intended for Authorized Forwarders but could possibly be
  // extended to work with future contracts.
  // @param targets The addresses to set permissions on
  // @param senders The addresses that are allowed to send updates
  function acceptAuthorizedReceivers(
    address[] calldata targets,
    address[] calldata senders
  ) external validateAuthorizedSenderSetter {
    acceptOwnableContracts(targets);
    setAuthorizedSendersOn(targets, senders);
  }

  // @notice Allows the node operator to withdraw earned LINK to a given address
  // @dev The owner of the contract can be another wallet and does not have to be a Chainlink node
  // @param recipient The address to send the LINK token to
  // @param amount The amount to send (specified in wei)
  function withdraw(
    address recipient,
    uint256 amount
  ) external override(OracleInterface, IWithdrawal) onlyOwner validateAvailableFunds(amount) {
    assert(i_linkToken.transfer(recipient, amount));
  }

  // @notice Displays the amount of LINK that is available for the node operator to withdraw
  // @dev We use `ONE_FOR_CONSISTENT_GAS_COST` in place of 0 in storage
  // @return The amount of withdrawable LINK on the contract
  function withdrawable() external view override(OracleInterface, IWithdrawal) returns (uint256) {
    return _fundsAvailable();
  }

  // @notice Forward a call to another contract
  // @dev Only callable by the owner
  // @param to address
  // @param data to forward
  function ownerForward(address to, bytes calldata data) external onlyOwner validateNotToLINK(to) {
    require(to.code.length != 0, "Must forward to a contract");
    // solhint-disable-next-line avoid-low-level-calls
    (bool status, ) = to.call(data);
    require(status, "Forwarded call failed");
  }

  // @notice Interact with other LinkTokenReceiver contracts by calling transferAndCall
  // @param to The address to transfer to.
  // @param value The amount to be transferred.
  // @param data The extra data to be passed to the receiving contract.
  // @return success bool
  function ownerTransferAndCall(
    address to,
    uint256 value,
    bytes calldata data
  ) external override onlyOwner validateAvailableFunds(value) returns (bool success) {
    return i_linkToken.transferAndCall(to, value, data);
  }

  // @notice Distribute funds to multiple addresses using ETH send
  // to this payable function.
  // @dev Array length must be equal, ETH sent must equal the sum of amounts.
  // A malicious receiver could cause the distribution to revert, in which case
  // it is expected that the address is removed from the list.
  // @param receivers list of addresses
  // @param amounts list of amounts
  function distributeFunds(address payable[] calldata receivers, uint256[] calldata amounts) external payable {
    require(receivers.length > 0 && receivers.length == amounts.length, "Invalid array length(s)");
    uint256 valueRemaining = msg.value;
    for (uint256 i = 0; i < receivers.length; ++i) {
      uint256 sendAmount = amounts[i];
      valueRemaining = valueRemaining - sendAmount;
      (bool success, ) = receivers[i].call{value: sendAmount}("");
      require(success, "Address: unable to send value, recipient may have reverted");
    }
    require(valueRemaining == 0, "Too much ETH sent");
  }

  // @notice Allows recipient to cancel requests sent to this oracle contract.
  // Will transfer the LINK sent for the request back to the recipient address.
  // @dev Given params must hash to a commitment stored on the contract in order
  // for the request to be valid. Emits CancelOracleRequest event.
  // @param requestId The request ID
  // @param payment The amount of payment given (specified in wei)
  // @param callbackFunc The requester's specified callback function selector
  // @param expiration The time of the expiration for the request
  function cancelOracleRequest(
    bytes32 requestId,
    uint256 payment,
    bytes4 callbackFunc,
    uint256 expiration
  ) public override {
    bytes31 paramsHash = _buildParamsHash(payment, msg.sender, callbackFunc, expiration);
    require(s_commitments[requestId].paramsHash == paramsHash, "Params do not match request ID");
    // solhint-disable-next-line not-rely-on-time
    require(expiration <= block.timestamp, "Request is not expired");

    delete s_commitments[requestId];
    emit CancelOracleRequest(requestId);

    // Free up the escrowed funds, as we're sending them back to the requester
    s_tokensInEscrow -= payment;
    i_linkToken.transfer(msg.sender, payment);
  }

  // @notice Allows requester to cancel requests sent to this oracle contract.
  // Will transfer the LINK sent for the request back to the recipient address.
  // @dev Given params must hash to a commitment stored on the contract in order
  // for the request to be valid. Emits CancelOracleRequest event.
  // @param nonce The nonce used to generate the request ID
  // @param payment The amount of payment given (specified in wei)
  // @param callbackFunc The requester's specified callback function selector
  // @param expiration The time of the expiration for the request
  function cancelOracleRequestByRequester(
    uint256 nonce,
    uint256 payment,
    bytes4 callbackFunc,
    uint256 expiration
  ) external {
    cancelOracleRequest(keccak256(abi.encodePacked(msg.sender, nonce)), payment, callbackFunc, expiration);
  }

  // @notice Returns the address of the LINK token
  // @dev This is the public implementation for chainlinkTokenAddress, which is
  // an internal method of the ChainlinkClient contract
  function getChainlinkToken() public view override returns (address) {
    return address(i_linkToken);
  }

  // @notice Require that the token transfer action is valid
  // @dev OPERATOR_REQUEST_SELECTOR = multiword, ORACLE_REQUEST_SELECTOR = singleword
  function _validateTokenTransferAction(bytes4 funcSelector, bytes memory data) internal pure override {
    require(data.length >= MINIMUM_REQUEST_LENGTH, "Invalid request length");
    require(
      funcSelector == OPERATOR_REQUEST_SELECTOR || funcSelector == ORACLE_REQUEST_SELECTOR,
      "Must use whitelisted functions"
    );
  }

  // @notice Verify the Oracle Request and record necessary information
  // @param sender The sender of the request
  // @param payment The amount of payment given (specified in wei)
  // @param callbackAddress The callback address for the response
  // @param callbackFunctionId The callback function ID for the response
  // @param nonce The nonce sent by the requester
  function _verifyAndProcessOracleRequest(
    address sender,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 nonce,
    uint256 dataVersion
  ) private validateNotToLINK(callbackAddress) returns (bytes32 requestId, uint256 expiration) {
    requestId = keccak256(abi.encodePacked(sender, nonce));
    require(s_commitments[requestId].paramsHash == 0, "Must use a unique ID");
    // solhint-disable-next-line not-rely-on-time
    expiration = block.timestamp + EXPIRYTIME;
    bytes31 paramsHash = _buildParamsHash(payment, callbackAddress, callbackFunctionId, expiration);
    s_commitments[requestId] = Commitment(paramsHash, SafeCast.toUint8(dataVersion));
    s_tokensInEscrow = s_tokensInEscrow + payment;
    return (requestId, expiration);
  }

  // @notice Verify the Oracle request and unlock escrowed payment
  // @param requestId The fulfillment request ID that must match the requester's
  // @param payment The payment amount that will be released for the oracle (specified in wei)
  // @param callbackAddress The callback address to call for fulfillment
  // @param callbackFunctionId The callback function ID to use for fulfillment
  // @param expiration The expiration that the node should respond by before the requester can cancel
  function _verifyOracleRequestAndProcessPayment(
    bytes32 requestId,
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration,
    uint256 dataVersion
  ) internal {
    bytes31 paramsHash = _buildParamsHash(payment, callbackAddress, callbackFunctionId, expiration);
    require(s_commitments[requestId].paramsHash == paramsHash, "Params do not match request ID");
    require(s_commitments[requestId].dataVersion <= SafeCast.toUint8(dataVersion), "Data versions must match");
    s_tokensInEscrow = s_tokensInEscrow - payment;
    delete s_commitments[requestId];
  }

  // @notice Build the bytes31 hash from the payment, callback and expiration.
  // @param payment The payment amount that will be released for the oracle (specified in wei)
  // @param callbackAddress The callback address to call for fulfillment
  // @param callbackFunctionId The callback function ID to use for fulfillment
  // @param expiration The expiration that the node should respond by before the requester can cancel
  // @return hash bytes31
  function _buildParamsHash(
    uint256 payment,
    address callbackAddress,
    bytes4 callbackFunctionId,
    uint256 expiration
  ) internal pure returns (bytes31) {
    return bytes31(keccak256(abi.encodePacked(payment, callbackAddress, callbackFunctionId, expiration)));
  }

  // @notice Returns the LINK available in this contract, not locked in escrow
  // @return uint256 LINK tokens available
  function _fundsAvailable() private view returns (uint256) {
    return i_linkToken.balanceOf(address(this)) - (s_tokensInEscrow - ONE_FOR_CONSISTENT_GAS_COST);
  }

  // @notice concrete implementation of AuthorizedReceiver
  // @return bool of whether sender is authorized
  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return isAuthorizedSender(msg.sender) || owner() == msg.sender;
  }

  // MODIFIERS

  // @dev Reverts if the first 32 bytes of the bytes array is not equal to requestId
  // @param requestId bytes32
  // @param data bytes
  modifier validateMultiWordResponseId(bytes32 requestId, bytes calldata data) {
    require(data.length >= 32, "Response must be > 32 bytes");
    bytes32 firstDataWord;
    assembly {
      firstDataWord := calldataload(data.offset)
    }
    require(requestId == firstDataWord, "First word must be requestId");
    _;
  }

  // @dev Reverts if amount requested is greater than withdrawable balance
  // @param amount The given amount to compare to `s_withdrawableTokens`
  modifier validateAvailableFunds(uint256 amount) {
    require(_fundsAvailable() >= amount, "Amount requested is greater than withdrawable balance");
    _;
  }

  // @dev Reverts if request ID does not exist
  // @param requestId The given request ID to check in stored `commitments`
  modifier validateRequestId(bytes32 requestId) {
    require(s_commitments[requestId].paramsHash != 0, "Must have a valid requestId");
    _;
  }

  // @dev Reverts if the callback address is the LINK token
  // @param to The callback address
  modifier validateNotToLINK(address to) {
    require(to != address(i_linkToken), "Cannot call to LINK");
    _;
  }

  // @dev Reverts if the target address is owned by the operator
  modifier validateCallbackAddress(address callbackAddress) {
    require(!s_owned[callbackAddress], "Cannot call owned contract");
    _;
  }
}
