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

  // function isAuthorizedOracle(address oracle) external view returns (bool);

  // function getAuthorizedOracles() external view returns (address[] memory);

  // function isConsumerAllowed(address client, uint64 subscriptionId) external view;

  function sendRequest(bytes calldata data) external returns (bytes32);

  /**
   * @notice Time out all expired requests: unlocks funds and removes the ability for the request to be fulfilled
   * @param requestIdsToTimeout - A list of request IDs to time out
   */
  function timeoutRequests(bytes32[] calldata requestIdsToTimeout) external;
}
