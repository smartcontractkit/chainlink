// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title Chainlink Functions Subscription interface.
 */
interface IFunctionsSubscriptions {
  /**
   * @notice Get details about a subscription.
   * @param subscriptionId - ID of the subscription
   * @return balance - LINK balance of the subscription in juels.
   * @return blockedBalance - amount of LINK balance of the subscription in juels that is blocked for an in flight request.
   * @return owner - owner of the subscription.
   * @return consumers - list of consumer address which are able to use this subscription.
   */
  function getSubscription(uint64 subscriptionId)
    external
    view
    returns (
      uint96 balance,
      uint96 blockedBalance,
      address owner,
      address[] memory consumers
    );

  /**
   * @notice Get details about a consumer of a subscription.
   * @dev Only callable by a route
   * @param client - the consumer contract that initiated the request
   * @param subscriptionId - ID of the subscription
   * @return allowed - amount of LINK balance of the subscription in juels that is blocked for an in flight request.
   * @return initiatedRequests - owner of the subscription.
   * @return completedRequests - list of consumer address which are able to use this subscription.
   */
  function getConsumer(address client, uint64 subscriptionId)
    external
    view
    returns (
      bool allowed,
      uint64 initiatedRequests,
      uint64 completedRequests
    );

  /**
   * @notice Blocks funds on a subscription account to be used for an in flight request
   * @dev Only callable by a route
   * @param client - the consumer contract that initiated the request
   * @param subscriptionId - The subscription ID to block funds for
   * @param amount - The amount to transfer
   */
  function blockBalance(
    address client,
    uint64 subscriptionId,
    uint96 amount
  ) external;

  /**
   * @notice Unblocks funds on a subscription account once a request has completed
   * @dev Only callable by a route
   * @param client - the consumer contract that initiated the request
   * @param subscriptionId - The subscription ID to block funds for
   * @param amount - The amount to transfer
   */
    function unblockBalance(
    address client,
    uint64 subscriptionId,
    uint96 amount
  ) external;

  /**
   * @notice Moves funds from one subscription account to another.
   * @dev Only callable by a route
   * @param from - The subscription ID to remove funds from
   * @param to - The address to pay funds to, allowing them to withdraw
   * @param amount - The amount to transfer
   */
  function pay(
    uint64 from,
    address to,
    uint96 amount
  ) external;
}
