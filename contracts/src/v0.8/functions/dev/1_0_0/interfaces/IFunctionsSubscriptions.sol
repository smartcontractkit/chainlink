// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsResponse} from "../libraries/FunctionsResponse.sol";

/// @title Chainlink Functions Subscription interface.
interface IFunctionsSubscriptions {
  struct Subscription {
    uint96 balance; // ═════════╗ Common LINK balance that is controlled by the Router to be used for all consumer requests.
    address owner; // ══════════╝ The owner can fund/withdraw/cancel the subscription.
    uint96 blockedBalance; // ══╗ LINK balance that is reserved to pay for pending consumer requests.
    address proposedOwner; // ══╝ For safely transferring sub ownership.
    address[] consumers; // ════╸ Client contracts that can use the subscription
    bytes32 flags; // ══════════╸ Per-subscription flags
  }

  struct Consumer {
    bool allowed; // ══════════════╗ Owner can fund/withdraw/cancel the sub.
    uint64 initiatedRequests; //   ║ The number of requests that have been started
    uint64 completedRequests; // ══╝ The number of requests that have successfully completed or timed out
  }

  /// @notice Get details about a subscription.
  /// @param subscriptionId - the ID of the subscription
  /// @return subscription - see IFunctionsSubscriptions.Subscription for more information on the structure
  function getSubscription(uint64 subscriptionId) external view returns (Subscription memory);

  /// @notice Retrieve details about multiple subscriptions using an inclusive range
  /// @param subscriptionIdStart - the ID of the subscription to start the range at
  /// @param subscriptionIdEnd - the ID of the subscription to end the range at
  /// @return subscriptions - see IFunctionsSubscriptions.Subscription for more information on the structure
  function getSubscriptionsInRange(
    uint64 subscriptionIdStart,
    uint64 subscriptionIdEnd
  ) external view returns (Subscription[] memory);

  /// @notice Get details about a consumer of a subscription.
  /// @param client - the consumer contract address
  /// @param subscriptionId - the ID of the subscription
  /// @return consumer - see IFunctionsSubscriptions.Consumer for more information on the structure
  function getConsumer(address client, uint64 subscriptionId) external view returns (Consumer memory);

  /// @notice Get details about the total amount of LINK within the system
  /// @return totalBalance - total Juels of LINK held by the contract
  function getTotalBalance() external view returns (uint96);

  /// @notice Get details about the total number of subscription accounts
  /// @return count - total number of subscriptions in the system
  function getSubscriptionCount() external view returns (uint64);

  /// @notice Time out all expired requests: unlocks funds and removes the ability for the request to be fulfilled
  /// @param requestsToTimeoutByCommitment - A list of request commitments to time out
  /// @dev The commitment can be found on the "OracleRequest" event created when sending the request.
  function timeoutRequests(FunctionsResponse.Commitment[] calldata requestsToTimeoutByCommitment) external;

  /// @notice Oracle withdraw LINK earned through fulfilling requests
  /// @notice If amount is 0 the full balance will be withdrawn
  /// @notice Both signing and transmitting wallets will have a balance to withdraw
  /// @param recipient where to send the funds
  /// @param amount amount to withdraw
  function oracleWithdraw(address recipient, uint96 amount) external;

  /// @notice Owner cancel subscription, sends remaining link directly to the subscription owner.
  /// @dev Only callable by the Router Owner
  /// @param subscriptionId subscription id
  /// @dev notably can be called even if there are pending requests, outstanding ones may fail onchain
  function ownerCancelSubscription(uint64 subscriptionId) external;

  /// @notice Recover link sent with transfer instead of transferAndCall.
  /// @dev Only callable by the Router Owner
  /// @param to address to send link to
  function recoverFunds(address to) external;

  /// @notice Create a new subscription.
  /// @return subscriptionId - A unique subscription id.
  /// @dev You can manage the consumer set dynamically with addConsumer/removeConsumer.
  /// @dev Note to fund the subscription, use transferAndCall. For example
  /// @dev  LINKTOKEN.transferAndCall(
  /// @dev    address(ROUTER),
  /// @dev    amount,
  /// @dev    abi.encode(subscriptionId));
  function createSubscription() external returns (uint64);

  /// @notice Create a new subscription and add a consumer.
  /// @return subscriptionId - A unique subscription id.
  /// @dev You can manage the consumer set dynamically with addConsumer/removeConsumer.
  /// @dev Note to fund the subscription, use transferAndCall. For example
  /// @dev  LINKTOKEN.transferAndCall(
  /// @dev    address(ROUTER),
  /// @dev    amount,
  /// @dev    abi.encode(subscriptionId));
  function createSubscriptionWithConsumer(address consumer) external returns (uint64 subscriptionId);

  /// @notice Propose a new owner for a subscription.
  /// @dev Only callable by the Subscription's owner
  /// @param subscriptionId - ID of the subscription
  /// @param newOwner - proposed new owner of the subscription
  function proposeSubscriptionOwnerTransfer(uint64 subscriptionId, address newOwner) external;

  /// @notice Accept an ownership transfer.
  /// @param subscriptionId - ID of the subscription
  /// @dev will revert if original owner of subscriptionId has not requested that msg.sender become the new owner.
  function acceptSubscriptionOwnerTransfer(uint64 subscriptionId) external;

  /// @notice Remove a consumer from a Chainlink Functions subscription.
  /// @dev Only callable by the Subscription's owner
  /// @param subscriptionId - ID of the subscription
  /// @param consumer - Consumer to remove from the subscription
  function removeConsumer(uint64 subscriptionId, address consumer) external;

  /// @notice Add a consumer to a Chainlink Functions subscription.
  /// @dev Only callable by the Subscription's owner
  /// @param subscriptionId - ID of the subscription
  /// @param consumer - New consumer which can use the subscription
  function addConsumer(uint64 subscriptionId, address consumer) external;

  /// @notice Cancel a subscription
  /// @dev Only callable by the Subscription's owner
  /// @param subscriptionId - ID of the subscription
  /// @param to - Where to send the remaining LINK to
  function cancelSubscription(uint64 subscriptionId, address to) external;

  /// @notice Check to see if there exists a request commitment for all consumers for a given sub.
  /// @param subscriptionId - ID of the subscription
  /// @return true if there exists at least one unfulfilled request for the subscription, false otherwise.
  /// @dev Looping is bounded to MAX_CONSUMERS*(number of DONs).
  /// @dev Used to disable subscription canceling while outstanding request are present.
  function pendingRequestExists(uint64 subscriptionId) external view returns (bool);

  /// @notice Set subscription specific flags for a subscription.
  /// Each byte of the flag is used to represent a resource tier that the subscription can utilize.
  /// @param subscriptionId - ID of the subscription
  /// @param flags - desired flag values
  function setFlags(uint64 subscriptionId, bytes32 flags) external;

  /// @notice Get flags for a given subscription.
  /// @param subscriptionId - ID of the subscription
  /// @return flags - current flag values
  function getFlags(uint64 subscriptionId) external view returns (bytes32);
}
