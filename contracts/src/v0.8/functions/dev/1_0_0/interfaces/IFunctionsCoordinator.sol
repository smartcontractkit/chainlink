// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IFunctionsSubscriptions} from "./IFunctionsSubscriptions.sol";

/**
 * @title Chainlink Functions oracle interface.
 */
interface IFunctionsCoordinator {
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
   * @param subscriptionId Identifier of the subscription that will be charged for the request
   * @param data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param caller The client contract that is sending the request
   * @param subscriptionOwner The owner of the subscription
   * @return requestId A unique request identifier (unique per DON)
   * @return estimatedCost The cost in Juels of LINK that the request is estimated to charge if market conditions were to stay the same
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    address caller,
    address subscriptionOwner
  ) external returns (bytes32, uint96, uint256, uint256);
}
