// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {IRouterBase} from "./IRouterBase.sol";

/**
 * @title Chainlink Functions Router interface.
 */
interface IFunctionsRouter is IRouterBase {
  enum FulfillResult {
    USER_SUCCESS,
    USER_ERROR,
    INVALID_REQUEST_ID,
    INSUFFICIENT_GAS,
    INSUFFICIENT_SUBSCRIPTION_BALANCE,
    INTERNAL_ERROR
  }

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
   * @param jobId A jobId that identifies which route to send the request to
   * @return requestId A unique request identifier
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit,
    bytes32 jobId
  ) external returns (bytes32);

  /**
   * @notice The fee that will be paid to the Router owner for operating the network
   * @dev Only callable by the Coordinator contract that is saved in the commitment
   * @param requestId The identifier for the request
   * @param response response data from DON consensus
   * @param err error from DON consensus
   * @param juelsPerGas -
   * @param transmitter -
   * @param to -
   * @param amount -
   */
  function fulfill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    address transmitter,
    address[] memory to,
    uint96[] memory amount
  ) external returns (FulfillResult);
}
