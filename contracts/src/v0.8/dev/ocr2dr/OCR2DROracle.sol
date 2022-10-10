// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";
import "../AuthorizedReceiver.sol";
import "../../ConfirmedOwner.sol";

/**
 * @title OCR2DR oracle contract (stub for now)
 */
contract OCR2DROracle is OCR2DROracleInterface, AuthorizedReceiver, ConfirmedOwner {
  event OracleRequest(bytes32 requestId, bytes data);
  event CancelOracleRequest(bytes32 indexed requestId);
  event OracleResponse(bytes32 indexed requestId);

  error InvalidRequestID();
  error Unauthorized();
  error NonceMustBeUnique();

  struct Commitment {
    address client;
  }

  mapping(bytes32 => Commitment) s_commitments;

  constructor(address owner) ConfirmedOwner(owner) {}

  function sendRequest(
    uint256 nonce,
    uint256, /* subscriptionId */
    bytes calldata data
  ) external override returns (bytes32) {
    bytes32 requestId = keccak256(abi.encodePacked(msg.sender, nonce));
    if (s_commitments[requestId].client != address(0)) {
      revert NonceMustBeUnique();
    }
    s_commitments[requestId].client = msg.sender;
    emit OracleRequest(requestId, data);
    return requestId;
  }

  function cancelRequest(bytes32 requestId) external override validateRequestId(requestId) {
    emit CancelOracleRequest(requestId);
    delete s_commitments[requestId];
  }

  function fulfillRequest(
    bytes32 requestId,
    bytes calldata response,
    bytes calldata err
  ) external override validateRequestId(requestId) validateAuthorizedSender {
    OCR2DRClientInterface client = OCR2DRClientInterface(s_commitments[requestId].client);
    emit OracleResponse(requestId);
    client.handleOracleFulfillment(requestId, response, err);
    delete s_commitments[requestId];
  }

  modifier validateRequestId(bytes32 requestId) {
    if (s_commitments[requestId].client == address(0)) {
      revert InvalidRequestID();
    }
    _;
  }

  /**
   * @notice concrete implementation of AuthorizedReceiver
   * @return bool of whether sender is authorized
   */
  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return isAuthorizedSender(msg.sender) || owner() == msg.sender;
  }
}
