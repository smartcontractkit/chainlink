// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "./OCR2DRBillingInterface.sol";

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
   * @param data Encoded OCR2DR request data, use OCR2DRClient API to encode a request
   * @param billing Configuration for payment of the request
   * @return requestId A unique request identifier (unique per oracle)
   */
  function sendRequest(bytes calldata data, OCR2DRBillingInterface.RequestBilling calldata billing)
    external
    returns (bytes32);
}
