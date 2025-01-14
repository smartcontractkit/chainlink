pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {VRFCoordinatorV2Plus_V2Example} from "../dev/testhelpers/VRFCoordinatorV2Plus_V2Example.sol";
import {ExposedVRFCoordinatorV2_5} from "../dev/testhelpers/ExposedVRFCoordinatorV2_5.sol";
import {VRFCoordinatorV2_5} from "../dev/VRFCoordinatorV2_5.sol";
import {SubscriptionAPI} from "../dev/SubscriptionAPI.sol";
import {VRFV2PlusConsumerExample} from "../dev/testhelpers/VRFV2PlusConsumerExample.sol";
import {MockLinkToken} from "../../functions/tests/v1_X/testhelpers/MockLinkToken.sol";
import {MockV3Aggregator} from "../../shared/mocks/MockV3Aggregator.sol";
import {VRFV2PlusMaliciousMigrator} from "../dev/testhelpers/VRFV2PlusMaliciousMigrator.sol";

contract VRFCoordinatorV2Plus_Migration is BaseTest {
  uint256 internal constant DEFAULT_LINK_FUNDING = 10 ether; // 10 LINK
  uint256 internal constant DEFAULT_NATIVE_FUNDING = 50 ether; // 50 ETH
  uint32 internal constant DEFAULT_CALLBACK_GAS_LIMIT = 50_000;
  uint16 internal constant DEFAULT_REQUEST_CONFIRMATIONS = 3;
  uint32 internal constant DEFAULT_NUM_WORDS = 1;
  // VRF KeyV2 generated from a node; not sensitive information.
  // The secret key used to generate this key is: 10.
  bytes internal constant UNCOMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c7893aba425419bc27a3b6c7e693a24c696f794c2ed877a1593cbee53b037368d7";
  bytes internal constant COMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c701";
  bytes32 internal constant KEY_HASH = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";
  uint64 internal constant GAS_LANE_MAX_GAS = 5000 gwei;

  ExposedVRFCoordinatorV2_5 v1Coordinator;
  VRFCoordinatorV2Plus_V2Example v2Coordinator;
  ExposedVRFCoordinatorV2_5 v1Coordinator_noLink;
  VRFCoordinatorV2Plus_V2Example v2Coordinator_noLink;
  uint256 subId;
  uint256 subId_noLink;
  VRFV2PlusConsumerExample testConsumer;
  VRFV2PlusConsumerExample testConsumer_noLink;
  MockLinkToken linkToken;
  address linkTokenAddr;
  MockV3Aggregator linkNativeFeed;
  address v1CoordinatorAddr;
  address v2CoordinatorAddr;
  address v1CoordinatorAddr_noLink;
  address v2CoordinatorAddr_noLink;

  event CoordinatorRegistered(address coordinatorAddress);
  event CoordinatorDeregistered(address coordinatorAddress);
  event MigrationCompleted(address newCoordinator, uint256 subId);

  function setUp() public override {
    BaseTest.setUp();
    vm.deal(OWNER, 100 ether);
    address bhs = makeAddr("bhs");
    v1Coordinator = new ExposedVRFCoordinatorV2_5(bhs);
    v1Coordinator_noLink = new ExposedVRFCoordinatorV2_5(bhs);
    subId = v1Coordinator.createSubscription();
    subId_noLink = v1Coordinator_noLink.createSubscription();
    linkToken = new MockLinkToken();
    linkNativeFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)
    v1Coordinator.setLINKAndLINKNativeFeed(address(linkToken), address(linkNativeFeed));
    linkTokenAddr = address(linkToken);
    v2Coordinator = new VRFCoordinatorV2Plus_V2Example(address(linkToken), address(v1Coordinator));
    v2Coordinator_noLink = new VRFCoordinatorV2Plus_V2Example(address(0), address(v1Coordinator_noLink));
    v1CoordinatorAddr = address(v1Coordinator);
    v2CoordinatorAddr = address(v2Coordinator);
    v1CoordinatorAddr_noLink = address(v1Coordinator_noLink);
    v2CoordinatorAddr_noLink = address(v2Coordinator_noLink);

    vm.expectEmit(
      false, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    emit CoordinatorRegistered(v2CoordinatorAddr);
    v1Coordinator.registerMigratableCoordinator(v2CoordinatorAddr);
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));

    vm.expectEmit(
      false, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    emit CoordinatorRegistered(v2CoordinatorAddr_noLink);
    v1Coordinator_noLink.registerMigratableCoordinator(v2CoordinatorAddr_noLink);
    assertTrue(v1Coordinator_noLink.isTargetRegisteredExternal(v2CoordinatorAddr_noLink));

    testConsumer = new VRFV2PlusConsumerExample(address(v1Coordinator), address(linkToken));
    testConsumer_noLink = new VRFV2PlusConsumerExample(address(v1Coordinator_noLink), address(0));
    v1Coordinator.setConfig(
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      600,
      10_000,
      20_000,
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );
    v1Coordinator_noLink.setConfig(
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      600,
      10_000,
      20_000,
      500_000, // fulfillmentFlatFeeNativePPM
      100_000, // fulfillmentFlatFeeLinkDiscountPPM
      15, // nativePremiumPercentage
      10 // linkPremiumPercentage
    );
    registerProvingKey();
    testConsumer.setCoordinator(v1CoordinatorAddr);
    testConsumer_noLink.setCoordinator(v1CoordinatorAddr_noLink);
  }

  function testDeregister() public {
    vm.expectEmit(
      false, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    emit CoordinatorDeregistered(v2CoordinatorAddr);
    v1Coordinator.deregisterMigratableCoordinator(v2CoordinatorAddr);
    assertFalse(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));

    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2_5.CoordinatorNotRegistered.selector, v2CoordinatorAddr));
    v1Coordinator.migrate(subId, v2CoordinatorAddr);

    // test register/deregister multiple coordinators
    address v3CoordinatorAddr = makeAddr("v3Coordinator");
    v1Coordinator.registerMigratableCoordinator(v2CoordinatorAddr);
    v1Coordinator.registerMigratableCoordinator(v3CoordinatorAddr);
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v3CoordinatorAddr));

    v1Coordinator.deregisterMigratableCoordinator(v3CoordinatorAddr);
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));
    assertFalse(v1Coordinator.isTargetRegisteredExternal(v3CoordinatorAddr));

    v1Coordinator.registerMigratableCoordinator(v3CoordinatorAddr);
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v3CoordinatorAddr));

    v1Coordinator.deregisterMigratableCoordinator(v2CoordinatorAddr);
    assertFalse(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));
    assertTrue(v1Coordinator.isTargetRegisteredExternal(v3CoordinatorAddr));

    v1Coordinator.deregisterMigratableCoordinator(v3CoordinatorAddr);
    assertFalse(v1Coordinator.isTargetRegisteredExternal(v2CoordinatorAddr));
    assertFalse(v1Coordinator.isTargetRegisteredExternal(v3CoordinatorAddr));
  }

  function testMigration() public {
    linkToken.transferAndCall(v1CoordinatorAddr, DEFAULT_LINK_FUNDING, abi.encode(subId));
    v1Coordinator.fundSubscriptionWithNative{value: DEFAULT_NATIVE_FUNDING}(subId);
    v1Coordinator.addConsumer(subId, address(testConsumer));

    // subscription exists in V1 coordinator before migration
    (uint96 balance, uint96 nativeBalance, uint64 reqCount, address owner, address[] memory consumers) = v1Coordinator
      .getSubscription(subId);
    assertEq(balance, DEFAULT_LINK_FUNDING);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(owner, address(OWNER));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(testConsumer));

    assertEq(v1Coordinator.s_totalBalance(), DEFAULT_LINK_FUNDING);
    assertEq(v1Coordinator.s_totalNativeBalance(), DEFAULT_NATIVE_FUNDING);

    // Update consumer to point to the new coordinator
    vm.expectEmit(
      false, // no first indexed field
      false, // no second indexed field
      false, // no third indexed field
      true // check data fields
    );
    emit MigrationCompleted(v2CoordinatorAddr, subId);
    v1Coordinator.migrate(subId, v2CoordinatorAddr);

    // subscription no longer exists in v1 coordinator after migration
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    v1Coordinator.getSubscription(subId);
    assertEq(v1Coordinator.s_totalBalance(), 0);
    assertEq(v1Coordinator.s_totalNativeBalance(), 0);
    assertEq(linkToken.balanceOf(v1CoordinatorAddr), 0);
    assertEq(v1CoordinatorAddr.balance, 0);

    // subscription exists in v2 coordinator
    (balance, nativeBalance, reqCount, owner, consumers) = v2Coordinator.getSubscription(subId);
    assertEq(owner, address(OWNER));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(testConsumer));
    assertEq(reqCount, 0);
    assertEq(balance, DEFAULT_LINK_FUNDING);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(v2Coordinator.s_totalLinkBalance(), DEFAULT_LINK_FUNDING);
    assertEq(v2Coordinator.s_totalNativeBalance(), DEFAULT_NATIVE_FUNDING);
    assertEq(linkToken.balanceOf(v2CoordinatorAddr), DEFAULT_LINK_FUNDING);
    assertEq(v2CoordinatorAddr.balance, DEFAULT_NATIVE_FUNDING);

    // calling migrate again on V1 coordinator should fail
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    v1Coordinator.migrate(subId, v2CoordinatorAddr);

    // test request still works after migration
    testConsumer.requestRandomWords(
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_NUM_WORDS,
      KEY_HASH,
      false
    );
    assertEq(testConsumer.s_recentRequestId(), 1);

    v2Coordinator.fulfillRandomWords(testConsumer.s_recentRequestId());
    assertEq(
      testConsumer.getRandomness(testConsumer.s_recentRequestId(), 0),
      v2Coordinator.generateFakeRandomness(testConsumer.s_recentRequestId())[0]
    );
  }

  function testMigrationNoLink() public {
    v1Coordinator_noLink.fundSubscriptionWithNative{value: DEFAULT_NATIVE_FUNDING}(subId_noLink);
    v1Coordinator_noLink.addConsumer(subId_noLink, address(testConsumer_noLink));

    // subscription exists in V1 coordinator before migration
    (
      uint96 balance,
      uint96 nativeBalance,
      uint64 reqCount,
      address owner,
      address[] memory consumers
    ) = v1Coordinator_noLink.getSubscription(subId_noLink);
    assertEq(balance, 0);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(owner, address(OWNER));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(testConsumer_noLink));

    assertEq(v1Coordinator_noLink.s_totalBalance(), 0);
    assertEq(v1Coordinator_noLink.s_totalNativeBalance(), DEFAULT_NATIVE_FUNDING);

    // Update consumer to point to the new coordinator
    vm.expectEmit(
      false, // no first indexed field
      false, // no second indexed field
      false, // no third indexed field
      true // check data fields
    );
    emit MigrationCompleted(v2CoordinatorAddr_noLink, subId_noLink);
    v1Coordinator_noLink.migrate(subId_noLink, v2CoordinatorAddr_noLink);

    // subscription no longer exists in v1 coordinator after migration
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    v1Coordinator_noLink.getSubscription(subId);
    assertEq(v1Coordinator_noLink.s_totalBalance(), 0);
    assertEq(v1Coordinator_noLink.s_totalNativeBalance(), 0);
    assertEq(linkToken.balanceOf(v1CoordinatorAddr_noLink), 0);
    assertEq(v1CoordinatorAddr_noLink.balance, 0);

    // subscription exists in v2 coordinator
    (balance, nativeBalance, reqCount, owner, consumers) = v2Coordinator_noLink.getSubscription(subId_noLink);
    assertEq(owner, address(OWNER));
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(testConsumer_noLink));
    assertEq(reqCount, 0);
    assertEq(balance, 0);
    assertEq(nativeBalance, DEFAULT_NATIVE_FUNDING);
    assertEq(v2Coordinator_noLink.s_totalLinkBalance(), 0);
    assertEq(v2Coordinator_noLink.s_totalNativeBalance(), DEFAULT_NATIVE_FUNDING);
    assertEq(linkToken.balanceOf(v2CoordinatorAddr_noLink), 0);
    assertEq(v2CoordinatorAddr_noLink.balance, DEFAULT_NATIVE_FUNDING);

    // calling migrate again on V1 coordinator should fail
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    v1Coordinator_noLink.migrate(subId_noLink, v2CoordinatorAddr_noLink);

    // test request still works after migration
    testConsumer_noLink.requestRandomWords(
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_NUM_WORDS,
      KEY_HASH,
      false
    );
    assertEq(testConsumer_noLink.s_recentRequestId(), 1);

    v2Coordinator_noLink.fulfillRandomWords(testConsumer_noLink.s_recentRequestId());
    assertEq(
      testConsumer_noLink.getRandomness(testConsumer_noLink.s_recentRequestId(), 0),
      v2Coordinator_noLink.generateFakeRandomness(testConsumer_noLink.s_recentRequestId())[0]
    );
  }

  function testMigrateRevertsWhenInvalidCoordinator() external {
    address invalidCoordinator = makeAddr("invalidCoordinator");

    vm.expectRevert(
      abi.encodeWithSelector(VRFCoordinatorV2_5.CoordinatorNotRegistered.selector, address(invalidCoordinator))
    );
    v1Coordinator.migrate(subId, invalidCoordinator);
  }

  function testMigrateRevertsWhenInvalidCaller() external {
    changePrank(makeAddr("invalidCaller"));
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.MustBeSubOwner.selector, OWNER));
    v1Coordinator.migrate(subId, v2CoordinatorAddr);
  }

  function testMigrateRevertsWhenPendingFulfillment() external {
    v1Coordinator.addConsumer(subId, address(testConsumer));
    testConsumer.setSubId(subId);
    testConsumer.requestRandomWords(
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_NUM_WORDS,
      KEY_HASH,
      false
    );

    vm.expectRevert(SubscriptionAPI.PendingRequestExists.selector);
    v1Coordinator.migrate(subId, v2CoordinatorAddr);
  }

  function testMigrateRevertsWhenReentrant() public {
    // deploy malicious contracts, subscriptions
    address maliciousUser = makeAddr("maliciousUser");
    changePrank(maliciousUser);
    uint256 maliciousSubId = v1Coordinator.createSubscription();
    VRFV2PlusMaliciousMigrator prankster = new VRFV2PlusMaliciousMigrator(address(v1Coordinator));
    v1Coordinator.addConsumer(maliciousSubId, address(prankster));

    // try to migrate malicious subscription, should fail
    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.Reentrant.selector));
    v1Coordinator.migrate(maliciousSubId, v2CoordinatorAddr);
  }

  function registerProvingKey() public {
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(UNCOMPRESSED_PUBLIC_KEY);
    v1Coordinator.registerProvingKey(uncompressedKeyParts, GAS_LANE_MAX_GAS);
    v1Coordinator_noLink.registerProvingKey(uncompressedKeyParts, GAS_LANE_MAX_GAS);
  }

  // note: Call this function via this.getProvingKeyParts to be able to pass memory as calldata and
  // index over the byte array.
  function getProvingKeyParts(bytes calldata uncompressedKey) public pure returns (uint256[2] memory) {
    uint256 keyPart1 = uint256(bytes32(uncompressedKey[0:32]));
    uint256 keyPart2 = uint256(bytes32(uncompressedKey[32:64]));
    return [keyPart1, keyPart2];
  }
}
