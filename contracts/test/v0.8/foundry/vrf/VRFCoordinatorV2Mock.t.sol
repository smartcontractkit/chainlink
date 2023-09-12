pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {VRFCoordinatorV2Mock} from "../../../../src/v0.8/mocks/VRFCoordinatorV2Mock.sol";
import {VRFConsumerV2} from "../../../../src/v0.8/vrf/testhelpers/VRFConsumerV2.sol";
import {VRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/VRFCoordinatorV2Plus.sol";
import {console} from "forge-std/console.sol";

contract VRFCoordinatorV2MockTest is BaseTest {
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkEthFeed;
  VRFCoordinatorV2Mock s_vrfCoordinatorV2Mock;
  VRFConsumerV2 s_vrfConsumerV2;
  uint256 s_subId;
  address s_subOwner = address(1234);

  // VRF KeyV2 generated from a node; not sensitive information.
  // The secret key used to generate this key is: 10.
  bytes internal constant UNCOMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c7893aba425419bc27a3b6c7e693a24c696f794c2ed877a1593cbee53b037368d7";
  bytes internal constant COMPRESSED_PUBLIC_KEY =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c701";
  bytes32 internal constant KEY_HASH = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";

  uint96 pointOneLink = 0.1 ether;
  uint256 oneLink = 1 ether;
  bytes32 keyHash = hex"e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0";
  address internal constant testConsumerAddress = 0x1111000000000000000000000000000000001111;
  address internal constant testConsumerAddress2 = 0x1111000000000000000000000000000000001110;

  event SubscriptionCreated(uint64 indexed subId, address owner);

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(OWNER, 10_000 ether);

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
    s_subId = s_vrfCoordinatorV2Mock.createSubscription();
    
    s_vrfCoordinatorV2Mock.setConfig();
  }

  function testCreateSubscription() public {
    vm.expectEmit(
      true, // no first indexed topic
      false, // no second indexed topic
      false, // no third indexed topic
      true // check data (target coordinator address)
    );
    emit SubscriptionCreated(2, address(OWNER));
    uint64 subId = s_vrfCoordinatorV2Mock.createSubscription();
    assertEq(subId, 2);

    (uint96 balance,
     uint64 reqCount,
     address owner,
     address[] memory consumers) = s_vrfCoordinatorV2Mock.getSubscription(subId);
    assertEq(balance, 0);
    assertEq(reqCount, 0);
    assertEq(owner, address(OWNER));
    assertEq(consumers.length, 0);
    // s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);
  }
}