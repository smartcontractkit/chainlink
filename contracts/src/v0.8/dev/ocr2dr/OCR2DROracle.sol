// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";

/**
 * @title OCR2DR oracle contract (stub for now)
 */
contract OCR2DROracle is OCR2DROracleInterface {
  event OracleRequest(address sender, bytes32 requestId, bytes data);
  event CancelOracleRequest(bytes32 indexed requestId);
  event OracleResponse(bytes32 indexed requestId);

  function sendRequest(
    address sender,
    uint256 nonce,
    uint256, /* subscriptionId */
    bytes calldata data
  ) external override returns (bytes32) {
    bytes32 requestId = keccak256(abi.encodePacked(msg.sender, nonce));
    emit OracleRequest(sender, requestId, data);
    return requestId;
  }

  function cancelRequest(bytes32 requestId) external override {
    emit CancelOracleRequest(requestId);
  }

  function fulfillRequest(
    bytes32 requestId,
    address callbackAddress,
    bytes calldata response,
    bytes calldata err
  ) external override {
    OCR2DRClientInterface client = OCR2DRClientInterface(callbackAddress);
    emit OracleResponse(requestId);
    client.handleOracleFulfillment(requestId, response, err);
  }
}
