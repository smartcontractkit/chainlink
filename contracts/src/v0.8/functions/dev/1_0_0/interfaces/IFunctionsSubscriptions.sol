// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title Chainlink Functions Subscription interface.
 */
interface IFunctionsSubscriptions {
  struct Subscription {
    // There are only 1e9*1e18 = 1e27 juels in existence, so the balance can fit in uint96 (2^96 ~ 7e28)
    uint96 balance; // Common LINK balance that is controlled by the Registry to be used for all consumer requests.
    address owner; // Owner can fund/withdraw/cancel the sub.
    uint96 blockedBalance; // LINK balance that is reserved to pay for pending consumer requests.
    address requestedOwner; // For safely transferring sub ownership.
    // Maintains the list of keys in s_consumers.
    // We do this for 2 reasons:
    // 1. To be able to clean up all keys from s_consumers when canceling a subscription.
    // 2. To be able to return the list of all consumers in getSubscription.
    // Note that we need the s_consumers map to be able to directly check if a
    // consumer is valid without reading all the consumers from storage.
    address[] consumers;
    bytes32 flags; // Per-subscription flags.
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
   * @notice Get the maximum number of consumers that can be added to one subscription
   * @return maxConsumers - maximum number of consumers that can be added to one subscription
   */
  function getMaxConsumers() external view returns (uint16);

  /**
   * @notice Get details about the total amount of LINK within the system
   * @return totalBalance - total Juels of LINK held by the contract
   */
  function getTotalBalance() external view returns (uint96);

  /**
   * @notice Get details about the total number of subscription accounts
   * @return count - total number of subscriptions in the system
   */
  function getSubscriptionCount() external view returns (uint64);

  /**
   * @notice Time out all expired requests: unlocks funds and removes the ability for the request to be fulfilled
   * @param requestIdsToTimeout - A list of request IDs to time out
   */
  function timeoutRequests(bytes32[] calldata requestIdsToTimeout) external;

  /**
   * @notice Oracle withdraw LINK earned through fulfilling requests
   * @dev Must be called by the Coordinator contract
   * @notice If amount is 0 the full balance will be withdrawn
   * @notice Both signing and transmitting wallets will have a balance to withdraw
   * @param recipient where to send the funds
   * @param amount amount to withdraw
   */
  function oracleWithdraw(address recipient, uint96 amount) external;

  /**
   * @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
   * @dev Only callable by the Router Owner
   * @param subscriptionId subscription id
   * @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
   */
  function ownerCancelSubscription(uint64 subscriptionId) external;

  /**
   * @notice Recover link sent with transfer instead of transferAndCall.
   * @dev Only callable by the Router Owner
   * @param to address to send link to
   */
  function recoverFunds(address to) external;

  /**
   * @notice Create a new subscription.
   * @return subscriptionId - A unique subscription id.
   * @dev You can manage the consumer set dynamically with addConsumer/removeConsumer.
   * @dev Note to fund the subscription, use transferAndCall. For example
   * @dev  LINKTOKEN.transferAndCall(
   * @dev    address(REGISTRY),
   * @dev    amount,
   * @dev    abi.encode(subscriptionId));
   */
  function createSubscription() external returns (uint64);

  /**
   * @notice Request subscription owner transfer.
   * @dev Only callable by the Subscription's owner
   * @param subscriptionId - ID of the subscription
   * @param newOwner - proposed new owner of the subscription
   */
  function requestSubscriptionOwnerTransfer(uint64 subscriptionId, address newOwner) external;

  /**
   * @notice Request subscription owner transfer.
   * @param subscriptionId - ID of the subscription
   * @dev will revert if original owner of subscriptionId has
   * not requested that msg.sender become the new owner.
   */
  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external;

  /**
   * @notice Remove a consumer from a Chainlink Functions subscription.
   * @dev Only callable by the Subscription's owner
   * @param subscriptionId - ID of the subscription
   * @param consumer - Consumer to remove from the subscription
   */
  function removeConsumer(uint64 subscriptionId, address consumer) external;

  /**
   * @notice Add a consumer to a Chainlink Functions subscription.
   * @dev Only callable by the Subscription's owner
   * @param subscriptionId - ID of the subscription
   * @param consumer - New consumer which can use the subscription
   */
  function addConsumer(uint64 subscriptionId, address consumer) external;

  /**
   * @notice Cancel a subscription
   * @dev Only callable by the Subscription's owner
   * @param subscriptionId - ID of the subscription
   * @param to - Where to send the remaining LINK to
   */
  function cancelSubscription(uint64 subscriptionId, address to) external;

  /**
   * @notice Check to see if there exists a request commitment for all consumers for a given sub.
   * @param subscriptionId - ID of the subscription
   * @return true if there exists at least one unfulfilled request for the subscription, false
   * otherwise.
   * @dev Looping is bounded to MAX_CONSUMERS*(number of DONs).
   * @dev Used to disable subscription canceling while outstanding request are present.
   */
  function pendingRequestExists(uint64 subscriptionId) external view returns (bool);

  /**
   * @notice Set flags for a given subscription.
   * @param subscriptionId - ID of the subscription
   * @param flags - desired flag values
   */
  function setFlags(uint64 subscriptionId, bytes32 flags) external;

  /**
   * @notice Get flags for a given subscription.
   * @param subscriptionId - ID of the subscription
   * @return flags - current flag values
   */
  function getFlags(uint64 subscriptionId) external view returns (bytes32);
}
