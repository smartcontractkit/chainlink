// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import {IVRFCoordinatorV2PlusMigration} from "../interfaces/IVRFCoordinatorV2PlusMigration.sol";
import {VRFConsumerBaseV2Plus} from "../VRFConsumerBaseV2Plus.sol";
import {VRFV2PlusClient} from "../libraries/VRFV2PlusClient.sol";

/// @dev this contract is only meant for testing migration
/// @dev it is a simplified example of future version (V2) of VRFCoordinatorV2Plus
// solhint-disable-next-line contract-name-camelcase
contract VRFCoordinatorV2Plus_V2Example is IVRFCoordinatorV2PlusMigration {
  error SubscriptionIDCollisionFound();

  struct Subscription {
    uint96 linkBalance;
    uint96 nativeBalance;
    uint64 reqCount;
    address owner;
    address[] consumers;
  }

  mapping(uint256 => Subscription) public s_subscriptions; /* subId */ /* subscription */
  mapping(uint256 => address) public s_requestConsumerMapping; /* RequestId */ /* consumer address */

  uint96 public s_totalLinkBalance;
  uint96 public s_totalNativeBalance;
  // request ID nonce
  uint256 public s_requestId = 0;

  // older version of coordinator, from which migration is supported
  address public s_prevCoordinator;
  address public s_link;

  constructor(address link, address prevCoordinator) {
    s_link = link;
    s_prevCoordinator = prevCoordinator;
  }

  /***************************************************************************
   * Section: Subscription
   **************************************************************************/

  /// @dev Emitted when a subscription for a given ID cannot be found
  error InvalidSubscription();

  function getSubscription(
    uint256 subId
  )
    public
    view
    returns (uint96 linkBalance, uint96 nativeBalance, uint64 reqCount, address owner, address[] memory consumers)
  {
    if (s_subscriptions[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return (
      s_subscriptions[subId].linkBalance,
      s_subscriptions[subId].nativeBalance,
      s_subscriptions[subId].reqCount,
      s_subscriptions[subId].owner,
      s_subscriptions[subId].consumers
    );
  }

  /***************************************************************************
   * Section: Migration
   **************************************************************************/

  /// @notice emitted when caller is not a previous version of VRF coordinator
  /// @param sender caller
  /// @param previousCoordinator expected coordinator address
  error MustBePreviousCoordinator(address sender, address previousCoordinator);

  /// @notice emitted when version in the request doesn't match expected version
  error InvalidVersion(uint8 requestVersion, uint8 expectedVersion);

  /// @notice emitted when transferred balance (msg.value) does not match the metadata in V1MigrationData
  error InvalidNativeBalance(uint256 transferredValue, uint96 expectedValue);

  /// @dev encapsulates data migrated over from previous coordinator
  struct V1MigrationData {
    uint8 fromVersion;
    uint256 subId;
    address subOwner;
    address[] consumers;
    uint96 linkBalance;
    uint96 nativeBalance;
  }

  /**
   * @inheritdoc IVRFCoordinatorV2PlusMigration
   */
  function onMigration(bytes calldata encodedData) external payable override {
    if (msg.sender != s_prevCoordinator) {
      revert MustBePreviousCoordinator(msg.sender, s_prevCoordinator);
    }

    V1MigrationData memory migrationData = abi.decode(encodedData, (V1MigrationData));

    if (migrationData.fromVersion != 1) {
      revert InvalidVersion(migrationData.fromVersion, 1);
    }

    if (msg.value != uint256(migrationData.nativeBalance)) {
      revert InvalidNativeBalance(msg.value, migrationData.nativeBalance);
    }

    // it should be impossible to have a subscription id collision, for two reasons:
    // 1. the subscription ID is calculated using inputs that cannot be replicated under different
    // conditions.
    // 2. once a subscription is migrated it is deleted from the previous coordinator, so it cannot
    // be migrated again.
    // however, we should have this check here in case the `migrate` function on
    // future coordinators "forgets" to delete subscription data allowing re-migration of the same
    // subscription.
    if (s_subscriptions[migrationData.subId].owner != address(0)) {
      revert SubscriptionIDCollisionFound();
    }

    s_subscriptions[migrationData.subId] = Subscription({
      nativeBalance: migrationData.nativeBalance,
      linkBalance: migrationData.linkBalance,
      reqCount: 0,
      owner: migrationData.subOwner,
      consumers: migrationData.consumers
    });
    s_totalNativeBalance += migrationData.nativeBalance;
    s_totalLinkBalance += migrationData.linkBalance;
  }

  /***************************************************************************
   * Section: Request/Response
   **************************************************************************/

  function requestRandomWords(VRFV2PlusClient.RandomWordsRequest calldata req) external returns (uint256 requestId) {
    Subscription memory sub = s_subscriptions[req.subId];
    sub.reqCount = sub.reqCount + 1;
    return _handleRequest(msg.sender);
  }

  function _handleRequest(address requester) private returns (uint256) {
    s_requestId = s_requestId + 1;
    uint256 requestId = s_requestId;
    s_requestConsumerMapping[s_requestId] = requester;
    return requestId;
  }

  function generateFakeRandomness(uint256 requestID) public pure returns (uint256[] memory) {
    uint256[] memory randomness = new uint256[](1);
    randomness[0] = uint256(keccak256(abi.encode(requestID, "not random")));
    return randomness;
  }

  function fulfillRandomWords(uint256 requestId) external {
    VRFConsumerBaseV2Plus consumer = VRFConsumerBaseV2Plus(s_requestConsumerMapping[requestId]);
    consumer.rawFulfillRandomWords(requestId, generateFakeRandomness(requestId));
  }
}
