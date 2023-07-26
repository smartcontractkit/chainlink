// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouterBase} from "./IRouterBase.sol";

/**
 * @title Chainlink Functions Router interface.
 */
interface IFunctionsRouter is IRouterBase {
  /**
   * @notice The identifier of the route to retrieve the address of the access control contract
   * The access control contract controls which accounts can manage subscriptions
   * @return id - bytes32 id that can be passed to the "getContractById" of the Router
   */
  function getAllowListId() external pure returns (bytes32);

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
   * @param dataVersion Gas limit for the fulfillment callback
   * @param callbackGasLimit Gas limit for the fulfillment callback
   * @param donId An identifier used to determine which route to send the request along
   * @return requestId A unique request identifier
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    bytes32 donId
  ) external returns (bytes32);

  /**
   * @notice Fulfill the request by:
   * - calling back the data that the Oracle returned to the client contract
   * - pay the DON for processing the request
   * @dev Only callable by the Coordinator contract that is saved in the commitment
   * @param requestId The identifier for the request
   * @param response response data from DON consensus
   * @param err error from DON consensus
   * @param juelsPerGas -
   * @param costWithoutFulfillment -
   * @param transmitter -
   * @return fulfillResult -
   * @return callbackGasCostJuels -
   */
  function fulfill(
    bytes32 requestId,
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    uint96 costWithoutFulfillment,
    address transmitter
  ) external returns (uint8, uint96);

  /**
   * @notice Validate requested gas limit is below the subscription max.
   * @param subscriptionId subscription ID
   * @param callbackGasLimit desired callback gas limit
   */
  function isValidCallbackGasLimit(uint64 subscriptionId, uint32 callbackGasLimit) external view;
}
