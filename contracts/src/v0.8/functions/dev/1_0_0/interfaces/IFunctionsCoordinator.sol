// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsResponse} from "../libraries/FunctionsResponse.sol";

/// @title Chainlink Functions DON Coordinator interface.
interface IFunctionsCoordinator {
  /// @notice Returns the DON's threshold encryption public key used to encrypt secrets
  /// @dev All nodes on the DON have separate key shares of the threshold decryption key
  /// and nodes must participate in a threshold decryption OCR round to decrypt secrets
  /// @return thresholdPublicKey the DON's threshold encryption public key
  function getThresholdPublicKey() external view returns (bytes memory);

  /// @notice Sets the DON's threshold encryption public key used to encrypt secrets
  /// @dev Used to rotate the key
  /// @param thresholdPublicKey The new public key
  function setThresholdPublicKey(bytes calldata thresholdPublicKey) external;

  /// @notice Returns the DON's secp256k1 public key that is used to encrypt secrets
  /// @dev All nodes on the DON have the corresponding private key
  /// needed to decrypt the secrets encrypted with the public key
  /// @return publicKey the DON's public key
  function getDONPublicKey() external view returns (bytes memory);

  /// @notice Sets DON's secp256k1 public key used to encrypt secrets
  /// @dev Used to rotate the key
  /// @param donPublicKey The new public key
  function setDONPublicKey(bytes calldata donPublicKey) external;

  /// @notice Receives a request to be emitted to the DON for processing
  /// @param request The request metadata
  /// @dev see the struct for field descriptions
  /// @return commitment - The parameters of the request that must be held consistent at response time
  function startRequest(
    FunctionsResponse.RequestMeta calldata request
  ) external returns (FunctionsResponse.Commitment memory commitment);
}
