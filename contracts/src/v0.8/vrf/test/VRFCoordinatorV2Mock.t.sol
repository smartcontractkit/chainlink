pragma solidity 0.8.6;

import "./BaseTest.t.sol";
import {VRF} from "../VRF.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../tests/MockV3Aggregator.sol";
import {VRFCoordinatorV2Mock} from "../mocks/VRFCoordinatorV2Mock.sol";
import {VRFConsumerV2} from "../testhelpers/VRFConsumerV2.sol";

contract VRFCoordinatorV2MockTest is BaseTest {
  MockLinkToken internal s_linkToken;
  MockV3Aggregator internal s_linkEthFeed;
  VRFCoordinatorV2Mock internal s_vrfCoordinatorV2Mock;
  VRFConsumerV2 internal s_vrfConsumerV2;
  address internal s_subOwner = address(1234);
  address internal s_randomOwner = address(4567);

  // VRF KeyV2 generated from a node; not sensitive information.
  // The secret key used to generate this key is: 10.
  bytes internal constant UNCOMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c7893aba425419bc27a3b6c7e693a24c696f794c2ed877a1593cbee53b037368d7";
  bytes internal constant COMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c701";
  bytes32 internal constant KEY_HASH = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";

  uint32 internal constant DEFAULT_CALLBACK_GAS_LIMIT = 500_000;
  uint16 internal constant DEFAULT_REQUEST_CONFIRMATIONS = 3;
  uint32 internal constant DEFAULT_NUM_WORDS = 1;

  uint96 pointOneLink = 0.1 ether;
  uint96 oneLink = 1 ether;

  event SubscriptionCreated(uint64 indexed subId, address owner);
  event SubscriptionFunded(uint64 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionCanceled(uint64 indexed subId, address to, uint256 amount);
  event ConsumerAdded(uint64 indexed subId, address consumer);
  event ConsumerRemoved(uint64 indexed subId, address consumer);
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint64 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    address indexed sender
  );
  event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bool success);

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(OWNER, 10_000 ether);
    vm.deal(s_subOwner, 20 ether);

    // Deploy link token and link/eth feed.
    s_linkToken = new MockLinkToken();
    s_linkEthFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Deploy coordinator and consumer.
    s_vrfCoordinatorV2Mock = new VRFCoordinatorV2Mock(
      pointOneLink,
      1_000_000_000 // 0.000000001 LINK per gas
    );
    address coordinatorAddr = address(s_vrfCoordinatorV2Mock);
    s_vrfConsumerV2 = new VRFConsumerV2(coordinatorAddr, address(s_linkToken));

    s_vrfCoordinatorV2Mock.setConfig();
  }

  function testCreateSubscription() public {
    vm.startPrank(s_subOwner);
    vm.expectEmit(
      true, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    emit SubscriptionCreated(1, s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();
    assertEq(subId, 1);

    (uint96 balance, uint64 reqCount, address owner, address[] memory consumers) = s_vrfCoordinatorV2Mock
      .getSubscription(subId);
    assertEq(balance, 0);
    assertEq(reqCount, 0);
    assertEq(owner, s_subOwner);
    assertEq(consumers.length, 0);
    // s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);

    // Test if subId increments
    vm.expectEmit(true, false, false, true);
    emit SubscriptionCreated(2, s_subOwner);
    subId = s_vrfCoordinatorV2Mock.createSubscription();
    assertEq(subId, 2);
    vm.stopPrank();
  }

  function testAddConsumer() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();
    vm.expectEmit(true, false, false, true);
    emit ConsumerAdded(subId, address(s_vrfConsumerV2));
    s_vrfCoordinatorV2Mock.addConsumer(subId, address(s_vrfConsumerV2));

    (uint96 balance, uint64 reqCount, address owner, address[] memory consumers) = s_vrfCoordinatorV2Mock
      .getSubscription(subId);
    assertEq(balance, 0);
    assertEq(reqCount, 0);
    assertEq(owner, s_subOwner);
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_vrfConsumerV2));
    vm.stopPrank();
  }

  // cannot add a consumer to a nonexistent subscription
  function testAddConsumerToInvalidSub() public {
    vm.startPrank(s_subOwner);
    bytes4 reason = bytes4(keccak256("InvalidSubscription()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.addConsumer(1, address(s_vrfConsumerV2));
    vm.stopPrank();
  }

  // cannot add more than the consumer maximum
  function testAddMaxConsumers() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();
    // Add 100 consumers
    for (uint64 i = 101; i <= 200; ++i) {
      s_vrfCoordinatorV2Mock.addConsumer(subId, address(bytes20(keccak256(abi.encodePacked(i)))));
    }
    // Adding 101th consumer should revert
    bytes4 reason = bytes4(keccak256("TooManyConsumers()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.addConsumer(subId, address(s_vrfConsumerV2));
    vm.stopPrank();
  }

  // can remove a consumer from a subscription
  function testRemoveConsumerFromSub() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    s_vrfCoordinatorV2Mock.addConsumer(subId, address(s_vrfConsumerV2));

    (, , , address[] memory consumers) = s_vrfCoordinatorV2Mock.getSubscription(subId);
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_vrfConsumerV2));

    vm.expectEmit(true, false, false, true);
    emit ConsumerRemoved(subId, address(s_vrfConsumerV2));
    s_vrfCoordinatorV2Mock.removeConsumer(subId, address(s_vrfConsumerV2));

    vm.stopPrank();
  }

  // cannot remove a consumer from a nonexistent subscription
  function testRemoveConsumerFromInvalidSub() public {
    vm.startPrank(s_subOwner);
    bytes4 reason = bytes4(keccak256("InvalidSubscription()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.removeConsumer(1, address(s_vrfConsumerV2));
    vm.stopPrank();
  }

  // cannot remove a consumer after it is already removed
  function testRemoveConsumerAgain() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    s_vrfCoordinatorV2Mock.addConsumer(subId, address(s_vrfConsumerV2));

    (, , , address[] memory consumers) = s_vrfCoordinatorV2Mock.getSubscription(subId);
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_vrfConsumerV2));

    vm.expectEmit(true, false, false, true);
    emit ConsumerRemoved(subId, address(s_vrfConsumerV2));
    s_vrfCoordinatorV2Mock.removeConsumer(subId, address(s_vrfConsumerV2));

    // Removing consumer again should revert with InvalidConsumer
    bytes4 reason = bytes4(keccak256("InvalidConsumer()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.removeConsumer(subId, address(s_vrfConsumerV2));
    vm.stopPrank();
  }

  // can fund a subscription
  function testFundSubscription() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    vm.expectEmit(true, false, false, true);
    emit SubscriptionFunded(subId, 0, oneLink);
    s_vrfCoordinatorV2Mock.fundSubscription(subId, oneLink);

    (uint96 balance, , , address[] memory consumers) = s_vrfCoordinatorV2Mock.getSubscription(subId);
    assertEq(balance, oneLink);
    assertEq(consumers.length, 0);
    vm.stopPrank();
  }

  // cannot fund a nonexistent subscription
  function testFundInvalidSubscription() public {
    vm.startPrank(s_subOwner);

    // Removing consumer again should revert with InvalidConsumer
    bytes4 reason = bytes4(keccak256("InvalidSubscription()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.removeConsumer(1, address(s_vrfConsumerV2));

    vm.stopPrank();
  }

  // can cancel a subscription
  function testCancelSubscription() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    s_vrfCoordinatorV2Mock.fundSubscription(subId, oneLink);

    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, s_subOwner, oneLink);
    s_vrfCoordinatorV2Mock.cancelSubscription(subId, s_subOwner);

    bytes4 reason = bytes4(keccak256("InvalidSubscription()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.getSubscription(subId);

    vm.stopPrank();
  }

  // fails to fulfill without being a valid consumer
  function testRequestRandomWordsInvalidConsumer() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    s_vrfCoordinatorV2Mock.fundSubscription(subId, oneLink);

    bytes4 reason = bytes4(keccak256("InvalidConsumer()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.requestRandomWords(
      KEY_HASH,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS
    );
    vm.stopPrank();
  }

  // fails to fulfill with insufficient funds
  function testRequestRandomWordsInsufficientFunds() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    address consumerAddr = address(s_vrfConsumerV2);
    s_vrfCoordinatorV2Mock.addConsumer(subId, address(s_vrfConsumerV2));

    vm.stopPrank();

    vm.startPrank(consumerAddr);

    vm.expectEmit(true, false, false, true);
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS,
      address(s_subOwner)
    );
    uint256 reqId = s_vrfCoordinatorV2Mock.requestRandomWords(
      KEY_HASH,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS
    );

    bytes4 reason = bytes4(keccak256("InsufficientBalance()"));
    vm.expectRevert(toBytes(reason));
    s_vrfCoordinatorV2Mock.fulfillRandomWords(reqId, consumerAddr);

    vm.stopPrank();
  }

  // can request and fulfill [ @skip-coverage ]
  function testRequestRandomWordsHappyPath() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    s_vrfCoordinatorV2Mock.fundSubscription(subId, oneLink);

    address consumerAddr = address(s_vrfConsumerV2);
    s_vrfCoordinatorV2Mock.addConsumer(subId, consumerAddr);

    vm.expectEmit(true, false, false, true);
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS,
      address(s_subOwner)
    );
    uint256 reqId = s_vrfConsumerV2.requestRandomness(
      KEY_HASH,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS
    );

    vm.expectEmit(true, false, false, true);
    emit RandomWordsFulfilled(reqId, 1, 100090236000000000, true);
    s_vrfCoordinatorV2Mock.fulfillRandomWords(reqId, consumerAddr);

    vm.stopPrank();
  }

  // Correctly allows for user override of fulfillRandomWords [ @skip-coverage ]
  function testRequestRandomWordsUserOverride() public {
    vm.startPrank(s_subOwner);
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();

    s_vrfCoordinatorV2Mock.fundSubscription(subId, oneLink);

    address consumerAddr = address(s_vrfConsumerV2);
    s_vrfCoordinatorV2Mock.addConsumer(subId, consumerAddr);

    vm.expectEmit(true, false, false, true);
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      2,
      address(s_subOwner)
    );
    uint256 reqId = s_vrfConsumerV2.requestRandomness(
      KEY_HASH,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      2
    );

    bytes4 reason = bytes4(keccak256("InvalidRandomWords()"));
    vm.expectRevert(toBytes(reason));
    uint256[] memory words1 = new uint256[](5);
    words1[0] = 1;
    words1[1] = 2;
    words1[2] = 3;
    words1[3] = 4;
    words1[4] = 5;
    s_vrfCoordinatorV2Mock.fulfillRandomWordsWithOverride(reqId, consumerAddr, uint256[](words1));

    vm.expectEmit(true, false, false, true);
    uint256[] memory words2 = new uint256[](2);
    words1[0] = 2533;
    words1[1] = 1768;
    emit RandomWordsFulfilled(reqId, 1, 100072314000000000, true);
    s_vrfCoordinatorV2Mock.fulfillRandomWordsWithOverride(reqId, consumerAddr, words2);

    vm.stopPrank();
  }

  function toBytes(bytes4 _data) public pure returns (bytes memory) {
    return abi.encodePacked(_data);
  }
}
