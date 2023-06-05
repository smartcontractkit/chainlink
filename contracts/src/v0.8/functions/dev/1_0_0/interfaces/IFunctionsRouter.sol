// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IRouterBase} from "./IRouterBase.sol";

/**
 * @title Chainlink Functions Router interface.
 */
interface IFunctionsRouter is IRouterBase {
  /**
   * @notice The fee that will be paid to the Router owner for operating the network
   * @return fee Cost in Juels (1e18) of LINK
   */
  function getAdminFee() external view returns (uint96);

  /**
   * @notice Sends a request (encoded as data) using the provided subscriptionId
   * @param subscriptionId A unique subscription ID allocated by billing system,
   * a client can make requests from different contracts referencing the same subscription
   * @param data Encoded Chainlink Functions request data, use FunctionsClient API to encode a request
   * @param gasLimit Gas limit for the fulfillment callback
   * @return requestId A unique request identifier (unique per DON)
   */
  function sendRequest(uint64 subscriptionId, bytes calldata data, uint32 gasLimit) external returns (bytes32);
}
