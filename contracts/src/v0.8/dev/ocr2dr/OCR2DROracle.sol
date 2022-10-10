// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";

/**
 * @title OCR2DR oracle contract (stub for now)
 */
contract OCR2DROracle is OCR2DROracleInterface {
  event OracleRequest(bytes32 requestId, bytes data);
  event CancelOracleRequest(bytes32 indexed requestId);
  event OracleResponse(bytes32 indexed requestId);

  error UnknownRequest();
  error NonceMustBeUnique();

  struct Request {
    address client;
  }

  mapping(bytes32 => Request) s_requests;

  function sendRequest(
    uint256 nonce,
    uint256, /* subscriptionId */
    bytes calldata data
  ) external override returns (bytes32) {
    bytes32 requestId = keccak256(abi.encodePacked(msg.sender, nonce));
    if (s_requests[requestId].client != address(0)) {
      revert NonceMustBeUnique();
    }
    s_requests[requestId].client = msg.sender;
    emit OracleRequest(requestId, data);
    return requestId;
  }

  function cancelRequest(bytes32 requestId) external override {
    if (s_requests[requestId].client == address(0)) {
      revert UnknownRequest();
    }
    emit CancelOracleRequest(requestId);
    delete s_requests[requestId];
  }

  function fulfillRequest(
    bytes32 requestId,
    bytes calldata response,
    bytes calldata err
  ) external override {
    if (s_requests[requestId].client == address(0)) {
      revert UnknownRequest();
    }
    OCR2DRClientInterface client = OCR2DRClientInterface(s_requests[requestId].client);
    emit OracleResponse(requestId);
    client.handleOracleFulfillment(requestId, response, err);
    delete s_requests[requestId];
  }
}
