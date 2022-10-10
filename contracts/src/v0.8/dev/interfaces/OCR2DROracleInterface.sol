// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title OCR2DR oracle interface.
 */
interface OCR2DROracleInterface {
  function sendRequest(
    address sender,
    uint256 nonce,
    uint256 subscriptionId,
    bytes calldata data
  ) external returns (bytes32);

  function cancelRequest(bytes32 requestId) external;

  function fulfillRequest(
    bytes32 requestId,
    address callbackAddress,
    bytes calldata response,
    bytes calldata err
  ) external;
}
