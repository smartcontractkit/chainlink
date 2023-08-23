// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../BaseTest.t.sol";
import "../../../../src/v0.8/dev/vrf/VRFV2PlusPriceRegistry.sol";
import "../../../../src/v0.8/tests/MockV3Aggregator.sol";

contract VRFV2PlusPriceRegistryTest is BaseTest {
  event LinkEthFeedSet(address oldFeed, address newFeed);
  event LinkUSDFeedSet(address oldFeed, address newFeed);
  event EthUSDFeedSet(address oldFeed, address newFeed);

  VRFV2PlusPriceRegistry public s_registry;

  function setUp() public override {
    BaseTest.setUp();

    // 1 LINK == 0.01 eth
    MockV3Aggregator linkEthAggregator = new MockV3Aggregator(
      18, 1e16);
    // 1 LINK == 10 USD
    MockV3Aggregator linkUSDAggregator = new MockV3Aggregator(
      8, 10e8);
    // 1 ETH == 1000 USD
    MockV3Aggregator ethUSDAggregator = new MockV3Aggregator(
      8, 1000e8);

    s_registry = new VRFV2PlusPriceRegistry(
      address(linkEthAggregator),
      address(linkUSDAggregator),
      address(ethUSDAggregator));
  }

  function testSetLinkEthFeed() public {
    address oldFeed = address(s_registry.s_linkEthFeed());
    MockV3Aggregator newFeed = new MockV3Aggregator(18, 1e16);
    vm.expectEmit(false, false, false, true /* data */);
    emit LinkEthFeedSet(oldFeed, address(newFeed));
    s_registry.setLINKETHFeed(address(newFeed));
  }
}
