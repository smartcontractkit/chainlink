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
  event OracleResponse(bytes32 requestId);

  error InvalidRequestID();
  error Unauthorized();
  error NonceMustBeUnique();

  struct Commitment {
    address client;
    uint256 subscriptionId;
  }

  uint256 private s_nonce;
  mapping(bytes32 => Commitment) private s_commitments;

  constructor(address owner) ConfirmedOwner(owner) {}

  function sendRequest(
    uint256 subscriptionId,
    bytes calldata data
  ) external override returns (bytes32) {
    s_nonce++;
    bytes32 requestId = keccak256(abi.encodePacked(msg.sender, s_nonce));
    if (s_commitments[requestId].client != address(0)) {
      revert NonceMustBeUnique();
    }
    s_commitments[requestId] = Commitment(msg.sender, subscriptionId);
    emit OracleRequest(requestId, data);
    return requestId;
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

  function _canSetAuthorizedSenders() internal view override returns (bool) {
    return isAuthorizedSender(msg.sender) || owner() == msg.sender;
  }
}
