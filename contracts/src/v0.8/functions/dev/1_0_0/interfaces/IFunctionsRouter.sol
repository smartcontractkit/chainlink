// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IRouterBase} from "./IRouterBase.sol";
import {IFunctionsRequest} from "./IFunctionsRequest.sol";

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
   * @param response response data from DON consensus
   * @param err error from DON consensus
   * @param juelsPerGas - current rate of juels/gas
   * @param costWithoutFulfillment - The cost of processing the request (in Juels of LINK ), without fulfillment
   * @param transmitter - The Node that transmitted the OCR report
   * @param commitment - The parameters of the request that must be held consistent between request and response time
   * @return fulfillResult -
   * @return callbackGasCostJuels -
   */
  function fulfill(
    bytes memory response,
    bytes memory err,
    uint96 juelsPerGas,
    uint96 costWithoutFulfillment,
    address transmitter,
    IFunctionsRequest.Commitment memory commitment
  ) external returns (uint8, uint96);

  /**
   * @notice Validate requested gas limit is below the subscription max.
   * @param subscriptionId subscription ID
   * @param callbackGasLimit desired callback gas limit
   */
  function isValidCallbackGasLimit(uint64 subscriptionId, uint32 callbackGasLimit) external view;
}
