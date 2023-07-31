// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title Chainlink Functions oracle interface.
 */
interface IFunctionsCoordinator {
  struct Request {
    address requestingContract; // The client contract that is sending the request
    address subscriptionOwner; // The owner of the subscription
    bytes data; // Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
    uint64 subscriptionId; // Identifier of the subscription that will be charged for the request
    uint16 dataVersion; // The version of the structure of the encoded data
    bytes32 flags; // Per-subscription flags
    uint32 callbackGasLimit; // The amount of gas that the callback to the consuming contract can utilize
  }

  /**
   * @notice Returns the DON's threshold encryption public key used to encrypt secrets
   * @dev All nodes on the DON have separate key shares of the threshold decryption key
   * and nodes must participate in a threshold decryption OCR round to decrypt secrets
   * @return thresholdPublicKey the DON's threshold encryption public key
   */
  function getThresholdPublicKey() external view returns (bytes memory);

  /**
   * @notice Sets the DON's threshold encryption public key used to encrypt secrets
   * @dev Used to rotate the key
   * @param thresholdPublicKey The new public key
   */
  function setThresholdPublicKey(bytes calldata thresholdPublicKey) external;

  /**
   * @notice Returns the DON's secp256k1 public key that is used to encrypt secrets
   * @dev All nodes on the DON have the corresponding private key
   * needed to decrypt the secrets encrypted with the public key
   * @return publicKey the DON's public key
   */
  function getDONPublicKey() external view returns (bytes memory);

  /**
   * @notice Sets DON's secp256k1 public key used to encrypt secrets
   * @dev Used to rotate the key
   * @param donPublicKey The new public key
   */
  function setDONPublicKey(bytes calldata donPublicKey) external;

  /**
   * @notice Sets a per-node secp256k1 public key used to encrypt secrets for that node
   * @dev Callable only by contract owner and DON members
   * @param node node's address
   * @param publicKey node's public key
   */
  function setNodePublicKey(address node, bytes calldata publicKey) external;

  /**
   * @notice Deletes node's public key
   * @dev Callable only by contract owner or the node itself
   * @param node node's address
   */
  function deleteNodePublicKey(address node) external;

  /**
   * @notice Return two arrays of equal size containing DON members' addresses and their corresponding
   * public keys (or empty byte arrays if per-node key is not defined)
   */
  function getAllNodePublicKeys() external view returns (address[] memory, bytes[] memory);

  /**
   * @notice Sends a request (encoded as data) using the provided subscriptionId
   * @dev Callable only by the Router
   * @param request The request information, @dev see the struct for field descriptions
   * @return requestId A unique request identifier (unique per DON)
   * @return estimatedCost The cost in Juels of LINK that the request is estimated to charge if market conditions were to stay the same
   * @return gasAfterPaymentCalculation The amount of gas overhead that will be used after balances have already been changed
   * @return requestTimeoutSeconds The amount of time in seconds before this request is considered stale
   */
  function sendRequest(Request calldata request) external returns (bytes32, uint96, uint256, uint256);
}
