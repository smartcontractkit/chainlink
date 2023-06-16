// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @title Chainlink Functions Subscription interface.
 */
interface IFunctionsSubscriptions {
  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common LINK balance that is controlled by the Registry to be used for all consumer requests.
    uint96 blockedBalance; // LINK balance that is reserved to pay for pending consumer requests.
    address owner; // Owner can fund/withdraw/cancel the sub.
    address requestedOwner; // For safely transferring sub ownership.
    // Maintains the list of keys in s_consumers.
    // We do this for 2 reasons:
    // 1. To be able to clean up all keys from s_consumers when canceling a subscription.
    // 2. To be able to return the list of all consumers in getSubscription.
    // Note that we need the s_consumers map to be able to directly check if a
    // consumer is valid without reading all the consumers from storage.
    address[] consumers;
  }

  struct Consumer {
    bool allowed; // Owner can fund/withdraw/cancel the sub.
    uint64 initiatedRequests; // The number of requests that have been started
    uint64 completedRequests; // The number of requests that have successfully completed or timed out
  }

  /**
   * @notice Get details about a subscription.
   * @param subscriptionId - ID of the subscription
   * @return balance - LINK balance of the subscription in juels.
   * @return blockedBalance - amount of LINK balance of the subscription in juels that is blocked for an in flight request.
   * @return owner - owner of the subscription.
   * @return requestedOwner - proposed owner to move ownership of the subscription to.
   * @return consumers - list of consumer address which are able to use this subscription.
   */
  function getSubscription(
    uint64 subscriptionId
  )
    external
    view
    returns (uint96 balance, uint96 blockedBalance, address owner, address requestedOwner, address[] memory consumers);

  /**
   * @notice Get details about a consumer of a subscription.
   * @dev Only callable by a route
   * @param client - the consumer contract that initiated the request
   * @param subscriptionId - ID of the subscription
   * @return allowed - amount of LINK balance of the subscription in juels that is blocked for an in flight request.
   * @return initiatedRequests - owner of the subscription.
   * @return completedRequests - list of consumer address which are able to use this subscription.
   */
  function getConsumer(
    address client,
    uint64 subscriptionId
  ) external view returns (bool allowed, uint64 initiatedRequests, uint64 completedRequests);

  /**
   * @notice Time out all expired requests: unlocks funds and removes the ability for the request to be fulfilled
   * @param requestIdsToTimeout - A list of request IDs to time out
   */
  function timeoutRequests(bytes32[] calldata requestIdsToTimeout) external;
}
