pragma solidity 0.8.19;

import "./BaseTest.t.sol";
import {VRFV2PlusClient} from "../dev/libraries/VRFV2PlusClient.sol";
import {SubscriptionAPI} from "../dev/SubscriptionAPI.sol";
import {VRFCoordinatorV2_5Mock} from "../mocks/VRFCoordinatorV2_5Mock.sol";
import {VRFConsumerV2Plus} from "../testhelpers/VRFConsumerV2Plus.sol";
import {MockLinkToken} from "../../mocks/MockLinkToken.sol";

contract VRFCoordinatorV2_5MockTest is BaseTest {
  MockLinkToken internal s_linkToken;
  VRFCoordinatorV2_5Mock internal s_vrfCoordinatorV2_5Mock;
  VRFConsumerV2Plus internal s_vrfConsumerV2Plus;
  address internal s_subOwner = address(1234);

  bytes32 internal constant KEY_HASH = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";

  uint32 internal constant DEFAULT_CALLBACK_GAS_LIMIT = 500_000;
  uint16 internal constant DEFAULT_REQUEST_CONFIRMATIONS = 3;
  uint32 internal constant DEFAULT_NUM_WORDS = 1;

  uint96 internal constant oneNative = 1 ether;
  uint96 internal constant twoLink = 2 ether;

  event SubscriptionCreated(uint256 indexed subId, address owner);
  event SubscriptionFunded(uint256 indexed subId, uint256 oldBalance, uint256 newBalance);
  event SubscriptionFundedWithNative(uint256 indexed subId, uint256 oldNativeBalance, uint256 newNativeBalance);
  event SubscriptionConsumerAdded(uint256 indexed subId, address consumer);
  event SubscriptionConsumerRemoved(uint256 indexed subId, address consumer);
  event SubscriptionCanceled(uint256 indexed subId, address to, uint256 amountLink, uint256 amountNative);

  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint256 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bytes extraArgs,
    address indexed sender
  );
  event RandomWordsFulfilled(
    uint256 indexed requestId,
    uint256 outputSeed,
    uint256 indexed subId,
    uint96 payment,
    bool nativePayment,
    bool success,
    bool onlyPremium
  );

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(OWNER, 10_000 ether);
    vm.deal(s_subOwner, 20 ether);

    // Deploy link token and link/eth feed.
    s_linkToken = new MockLinkToken();

    // Deploy coordinator and consumer.
    s_vrfCoordinatorV2_5Mock = new VRFCoordinatorV2_5Mock(0.002 ether, 40 gwei, 0.004 ether);
    address coordinatorAddr = address(s_vrfCoordinatorV2_5Mock);
    s_vrfConsumerV2Plus = new VRFConsumerV2Plus(coordinatorAddr, address(s_linkToken));

    s_vrfCoordinatorV2_5Mock.setConfig();
  }

  function test_CreateSubscription() public {
    vm.startPrank(s_subOwner);
    uint256 expectedSubId = uint256(
      keccak256(
        abi.encodePacked(
          s_subOwner,
          blockhash(block.number - 1),
          address(s_vrfCoordinatorV2_5Mock),
          s_vrfCoordinatorV2_5Mock.s_currentSubNonce()
        )
      )
    );
    vm.expectEmit(
      true, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    emit SubscriptionCreated(expectedSubId, s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();
    assertEq(subId, expectedSubId);

    (
      uint96 balance,
      uint96 nativeBalance,
      uint64 reqCount,
      address owner,
      address[] memory consumers
    ) = s_vrfCoordinatorV2_5Mock.getSubscription(subId);
    assertEq(balance, 0);
    assertEq(nativeBalance, 0);
    assertEq(reqCount, 0);
    assertEq(owner, s_subOwner);
    assertEq(consumers.length, 0);
    vm.stopPrank();
  }

  function test_AddConsumer() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();
    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerAdded(subId, address(s_vrfConsumerV2Plus));
    s_vrfCoordinatorV2_5Mock.addConsumer(subId, address(s_vrfConsumerV2Plus));

    (uint96 balance, , uint64 reqCount, address owner, address[] memory consumers) = s_vrfCoordinatorV2_5Mock
      .getSubscription(subId);
    assertEq(balance, 0);
    assertEq(reqCount, 0);
    assertEq(owner, s_subOwner);
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_vrfConsumerV2Plus));
    vm.stopPrank();
  }

  // cannot add a consumer to a nonexistent subscription
  function test_AddConsumerToInvalidSub() public {
    vm.startPrank(s_subOwner);
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_vrfCoordinatorV2_5Mock.addConsumer(1, address(s_vrfConsumerV2Plus));
    vm.stopPrank();
  }

  // cannot add more than the consumer maximum
  function test_AddMaxConsumers() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();
    // Add 100 consumers
    for (uint64 i = 101; i <= 200; ++i) {
      s_vrfCoordinatorV2_5Mock.addConsumer(subId, address(bytes20(keccak256(abi.encodePacked(i)))));
    }
    // Adding 101th consumer should revert
    vm.expectRevert(SubscriptionAPI.TooManyConsumers.selector);
    s_vrfCoordinatorV2_5Mock.addConsumer(subId, address(s_vrfConsumerV2Plus));
    vm.stopPrank();
  }

  // can remove a consumer from a subscription
  function test_RemoveConsumerFromSub() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.addConsumer(subId, address(s_vrfConsumerV2Plus));

    (, , , , address[] memory consumers) = s_vrfCoordinatorV2_5Mock.getSubscription(subId);
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_vrfConsumerV2Plus));

    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerRemoved(subId, address(s_vrfConsumerV2Plus));
    s_vrfCoordinatorV2_5Mock.removeConsumer(subId, address(s_vrfConsumerV2Plus));

    // Removing consumer again should revert with InvalidConsumer
    vm.expectRevert(
      abi.encodeWithSelector(SubscriptionAPI.InvalidConsumer.selector, subId, address(s_vrfConsumerV2Plus))
    );
    s_vrfCoordinatorV2_5Mock.removeConsumer(subId, address(s_vrfConsumerV2Plus));

    vm.stopPrank();
  }

  // cannot remove a consumer from a nonexistent subscription
  function test_RemoveConsumerFromInvalidSub() public {
    vm.startPrank(s_subOwner);
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_vrfCoordinatorV2_5Mock.removeConsumer(1, address(s_vrfConsumerV2Plus));
    vm.stopPrank();
  }

  // cannot remove a consumer after it is already removed
  function test_RemoveConsumerAgain() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.addConsumer(subId, address(s_vrfConsumerV2Plus));

    (, , , , address[] memory consumers) = s_vrfCoordinatorV2_5Mock.getSubscription(subId);
    assertEq(consumers.length, 1);
    assertEq(consumers[0], address(s_vrfConsumerV2Plus));

    vm.expectEmit(true, false, false, true);
    emit SubscriptionConsumerRemoved(subId, address(s_vrfConsumerV2Plus));
    s_vrfCoordinatorV2_5Mock.removeConsumer(subId, address(s_vrfConsumerV2Plus));

    // Removing consumer again should revert with InvalidConsumer
    vm.expectRevert(
      abi.encodeWithSelector(SubscriptionAPI.InvalidConsumer.selector, subId, address(s_vrfConsumerV2Plus))
    );
    s_vrfCoordinatorV2_5Mock.removeConsumer(subId, address(s_vrfConsumerV2Plus));
    vm.stopPrank();
  }

  // can fund a subscription
  function test_FundSubscription() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    vm.expectEmit(true, false, false, true);
    emit SubscriptionFunded(subId, 0, twoLink);
    s_vrfCoordinatorV2_5Mock.fundSubscription(subId, twoLink);

    (uint96 balance, , , , address[] memory consumers) = s_vrfCoordinatorV2_5Mock.getSubscription(subId);
    assertEq(balance, twoLink);
    assertEq(consumers.length, 0);

    assertEq(s_vrfCoordinatorV2_5Mock.s_totalBalance(), twoLink);

    vm.stopPrank();
  }

  // cannot fund a nonexistent subscription
  function testFuzz_FundSubscription_RevertIfInvalidSubscription(uint256 subId) public {
    vm.startPrank(s_subOwner);

    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_vrfCoordinatorV2_5Mock.fundSubscription(subId, twoLink);

    vm.stopPrank();
  }

  // can fund a subscription with native
  function test_FundSubscriptionWithNative() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    vm.expectEmit(true, false, false, true);
    emit SubscriptionFundedWithNative(subId, 0, oneNative);
    s_vrfCoordinatorV2_5Mock.fundSubscriptionWithNative{value: oneNative}(subId);

    (, uint256 nativeBalance, , , address[] memory consumers) = s_vrfCoordinatorV2_5Mock.getSubscription(subId);
    assertEq(nativeBalance, oneNative);
    assertEq(consumers.length, 0);

    assertEq(s_vrfCoordinatorV2_5Mock.s_totalNativeBalance(), oneNative);

    vm.stopPrank();
  }

  // cannot fund a nonexistent subscription
  function testFuzz_FundSubscriptionWithNative_RevertIfInvalidSubscription(uint256 subId) public {
    vm.startPrank(s_subOwner);

    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_vrfCoordinatorV2_5Mock.fundSubscriptionWithNative{value: oneNative}(subId);

    vm.stopPrank();
  }

  // can cancel a subscription
  function test_CancelSubscription_Link() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.fundSubscription(subId, twoLink);

    uint256 totalBalance = s_vrfCoordinatorV2_5Mock.s_totalBalance();

    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, s_subOwner, twoLink, 0);
    s_vrfCoordinatorV2_5Mock.cancelSubscription(subId, s_subOwner);

    // check coordinator balance decreased
    assertEq(s_vrfCoordinatorV2_5Mock.s_totalBalance(), totalBalance - twoLink);

    // sub owner balance did not increase as no actual token is involved

    // check subscription removed
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_vrfCoordinatorV2_5Mock.getSubscription(subId);

    vm.stopPrank();
  }

  // can cancel a subscription
  function test_CancelSubscription_Native() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.fundSubscriptionWithNative{value: oneNative}(subId);

    uint256 balance = address(s_subOwner).balance;
    uint256 totalNativeBalance = s_vrfCoordinatorV2_5Mock.s_totalNativeBalance();

    vm.expectEmit(true, false, false, true);
    emit SubscriptionCanceled(subId, s_subOwner, 0, oneNative);
    s_vrfCoordinatorV2_5Mock.cancelSubscription(subId, s_subOwner);

    // check coordinator balance decreased
    assertEq(s_vrfCoordinatorV2_5Mock.s_totalNativeBalance(), totalNativeBalance - oneNative);

    // check sub owner balance increased
    assertEq(address(s_subOwner).balance, balance + oneNative);

    // check subscription removed
    vm.expectRevert(SubscriptionAPI.InvalidSubscription.selector);
    s_vrfCoordinatorV2_5Mock.getSubscription(subId);

    vm.stopPrank();
  }

  // fails to fulfill without being a valid consumer
  function testFuzz_RequestRandomWords_RevertIfInvalidConsumer(bool nativePayment) public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.fundSubscription(subId, twoLink);

    vm.expectRevert(abi.encodeWithSelector(SubscriptionAPI.InvalidConsumer.selector, subId, address(s_subOwner)));
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: KEY_HASH,
      subId: subId,
      requestConfirmations: DEFAULT_REQUEST_CONFIRMATIONS,
      callbackGasLimit: DEFAULT_CALLBACK_GAS_LIMIT,
      numWords: DEFAULT_NUM_WORDS,
      extraArgs: VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: nativePayment}))
    });
    s_vrfCoordinatorV2_5Mock.requestRandomWords(req);
    vm.stopPrank();
  }

  // fails to fulfill with insufficient funds
  function testFuzz_RequestRandomWords_RevertIfInsufficientFunds(bool nativePayment) public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    address consumerAddr = address(s_vrfConsumerV2Plus);
    s_vrfCoordinatorV2_5Mock.addConsumer(subId, address(s_vrfConsumerV2Plus));

    vm.stopPrank();

    vm.startPrank(consumerAddr);

    vm.expectEmit(true, false, false, true);
    bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: nativePayment}));
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS,
      extraArgs,
      address(s_subOwner)
    );
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: KEY_HASH,
      subId: subId,
      requestConfirmations: DEFAULT_REQUEST_CONFIRMATIONS,
      callbackGasLimit: DEFAULT_CALLBACK_GAS_LIMIT,
      numWords: DEFAULT_NUM_WORDS,
      extraArgs: extraArgs
    });
    uint256 reqId = s_vrfCoordinatorV2_5Mock.requestRandomWords(req);

    vm.expectRevert(SubscriptionAPI.InsufficientBalance.selector);
    s_vrfCoordinatorV2_5Mock.fulfillRandomWords(reqId, consumerAddr);

    vm.stopPrank();
  }

  // can request and fulfill [ @skip-coverage ]
  function test_RequestRandomWords_Link_HappyPath() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.fundSubscription(subId, twoLink);

    address consumerAddr = address(s_vrfConsumerV2Plus);
    s_vrfCoordinatorV2_5Mock.addConsumer(subId, consumerAddr);

    vm.expectEmit(true, false, false, true);
    bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: false}));
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS,
      extraArgs,
      address(s_subOwner)
    );
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: KEY_HASH,
      subId: subId,
      requestConfirmations: DEFAULT_REQUEST_CONFIRMATIONS,
      callbackGasLimit: DEFAULT_CALLBACK_GAS_LIMIT,
      numWords: DEFAULT_NUM_WORDS,
      extraArgs: extraArgs
    });
    uint256 reqId = s_vrfConsumerV2Plus.requestRandomness(req);

    vm.expectEmit(true, false, false, true);
    emit RandomWordsFulfilled(reqId, 1, subId, 1432960000000000000, false, true, false);
    s_vrfCoordinatorV2_5Mock.fulfillRandomWords(reqId, consumerAddr);

    vm.stopPrank();
  }

  // can request and fulfill [ @skip-coverage ]
  function test_RequestRandomWords_Native_HappyPath() public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    s_vrfCoordinatorV2_5Mock.fundSubscriptionWithNative{value: oneNative}(subId);

    address consumerAddr = address(s_vrfConsumerV2Plus);
    s_vrfCoordinatorV2_5Mock.addConsumer(subId, consumerAddr);

    vm.expectEmit(true, false, false, true);
    bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: true}));
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      DEFAULT_NUM_WORDS,
      extraArgs,
      address(s_subOwner)
    );
    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: KEY_HASH,
      subId: subId,
      requestConfirmations: DEFAULT_REQUEST_CONFIRMATIONS,
      callbackGasLimit: DEFAULT_CALLBACK_GAS_LIMIT,
      numWords: DEFAULT_NUM_WORDS,
      extraArgs: extraArgs
    });
    uint256 reqId = s_vrfConsumerV2Plus.requestRandomness(req);

    vm.expectEmit(true, false, false, true);
    emit RandomWordsFulfilled(reqId, 1, subId, 5731840000000000, true, true, false);
    s_vrfCoordinatorV2_5Mock.fulfillRandomWords(reqId, consumerAddr);

    vm.stopPrank();
  }

  // Correctly allows for user override of fulfillRandomWords [ @skip-coverage ]
  function testFuzz_RequestRandomWordsUserOverride(bool nativePayment) public {
    vm.startPrank(s_subOwner);
    uint256 subId = s_vrfCoordinatorV2_5Mock.createSubscription();

    uint96 expectedPayment;
    if (nativePayment) {
      expectedPayment = 5011440000000000;
      s_vrfCoordinatorV2_5Mock.fundSubscriptionWithNative{value: oneNative}(subId);
    } else {
      expectedPayment = 1252860000000000000;
      s_vrfCoordinatorV2_5Mock.fundSubscription(subId, twoLink);
    }

    address consumerAddr = address(s_vrfConsumerV2Plus);
    s_vrfCoordinatorV2_5Mock.addConsumer(subId, consumerAddr);

    bytes memory extraArgs = VRFV2PlusClient._argsToBytes(VRFV2PlusClient.ExtraArgsV1({nativePayment: nativePayment}));
    vm.expectEmit(true, false, false, true);
    emit RandomWordsRequested(
      KEY_HASH,
      1,
      100,
      subId,
      DEFAULT_REQUEST_CONFIRMATIONS,
      DEFAULT_CALLBACK_GAS_LIMIT,
      2,
      extraArgs,
      address(s_subOwner)
    );

    VRFV2PlusClient.RandomWordsRequest memory req = VRFV2PlusClient.RandomWordsRequest({
      keyHash: KEY_HASH,
      subId: subId,
      requestConfirmations: DEFAULT_REQUEST_CONFIRMATIONS,
      callbackGasLimit: DEFAULT_CALLBACK_GAS_LIMIT,
      numWords: 2,
      extraArgs: extraArgs
    });
    uint256 reqId = s_vrfConsumerV2Plus.requestRandomness(req);

    vm.expectRevert(VRFCoordinatorV2_5Mock.InvalidRandomWords.selector);
    uint256[] memory words1 = new uint256[](5);
    words1[0] = 1;
    words1[1] = 2;
    words1[2] = 3;
    words1[3] = 4;
    words1[4] = 5;
    s_vrfCoordinatorV2_5Mock.fulfillRandomWordsWithOverride(reqId, consumerAddr, uint256[](words1));

    vm.expectEmit(true, false, false, true);
    uint256[] memory words2 = new uint256[](2);
    words1[0] = 2533;
    words1[1] = 1768;
    emit RandomWordsFulfilled(reqId, 1, subId, expectedPayment, nativePayment, true, false);
    s_vrfCoordinatorV2_5Mock.fulfillRandomWordsWithOverride(reqId, consumerAddr, words2);

    vm.stopPrank();
  }
}
