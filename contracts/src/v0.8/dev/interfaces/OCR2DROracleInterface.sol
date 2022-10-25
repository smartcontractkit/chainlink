// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title OCR2DR oracle interface.
 */
interface OCR2DROracleInterface {
  /**
   * @notice Returns DON secp256k1 public key used to encrypt secrets
   * @dev All Oracles nodes have the corresponding private key
   * needed to decrypt the secrets encrypted with the public key
   * @return publicKey DON's public key
   */
  function getDONPublicKey() external view returns (bytes memory);

  /**
   * @notice Sets DON secp256k1 public key used to encrypt secrets
   * @dev Used to rotate the key
   * @param donPublicKey New public key
   */
  function setDONPublicKey(bytes calldata donPublicKey) external;

  /**
   * @notice Sends a request (encoded as data) using the provided subscriptionId
   * @param subscriptionId A unique subscription ID allocated by billing system,
   * a client can make requests from different contracts referencing the same subscription
   * @param data Encoded OCR2DR request data, use OCR2DRClient API to encode a request
   * @return requestId A unique request identifier (unique per oracle)
   */
  function sendRequest(uint256 subscriptionId, bytes calldata data) external returns (bytes32);
}
