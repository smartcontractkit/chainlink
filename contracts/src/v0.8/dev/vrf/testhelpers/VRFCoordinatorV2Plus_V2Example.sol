// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "../../../shared/interfaces/LinkTokenInterface.sol";
import "../../interfaces/IVRFCoordinatorV2PlusMigration.sol";
import "../../interfaces/IVRFMigratableCoordinatorV2Plus.sol";
import "../VRFConsumerBaseV2Plus.sol";

/// @dev this contract is only meant for testing migration
/// @dev it is a simplified example of future version (V2) of VRFCoordinatorV2Plus
contract VRFCoordinatorV2Plus_V2Example is IVRFCoordinatorV2PlusMigration, IVRFMigratableCoordinatorV2Plus {
  struct Subscription {
    address owner;
    address[] consumers;
    uint96 linkBalance;
    uint96 nativeBalance;
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
  ) public view returns (address owner, address[] memory consumers, uint96 linkBalance, uint96 nativeBalance) {
    if (s_subscriptions[subId].owner == address(0)) {
      revert InvalidSubscription();
    }
    return (
      s_subscriptions[subId].owner,
      s_subscriptions[subId].consumers,
      s_subscriptions[subId].linkBalance,
      s_subscriptions[subId].nativeBalance
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
    uint96 ethBalance;
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

    if (msg.value != uint256(migrationData.ethBalance)) {
      revert InvalidNativeBalance(msg.value, migrationData.ethBalance);
    }

    s_subscriptions[migrationData.subId] = Subscription({
      owner: migrationData.subOwner,
      consumers: migrationData.consumers,
      nativeBalance: migrationData.ethBalance,
      linkBalance: migrationData.linkBalance
    });
    s_totalNativeBalance += migrationData.ethBalance;
    s_totalLinkBalance += migrationData.linkBalance;
  }

  /***************************************************************************
   * Section: Request/Response
   **************************************************************************/

  /**
   * @inheritdoc IVRFMigratableCoordinatorV2Plus
   */
  function requestRandomWords(
    VRFV2PlusClient.RandomWordsRequest calldata /* req */
  ) external override returns (uint256 requestId) {
    return handleRequest(msg.sender);
  }

  function handleRequest(address requester) private returns (uint256) {
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
