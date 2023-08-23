// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../BaseTest.t.sol";
import "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import "../../../../src/v0.8/dev/vrf/testhelpers/ExposedPriceRegistry.sol";

contract VRFV2PlusPriceRegistryTest is BaseTest {
  event LinkEthFeedSet(address oldFeed, address newFeed);
  event LinkUSDFeedSet(address oldFeed, address newFeed);
  event EthUSDFeedSet(address oldFeed, address newFeed);
  event ConfigSet(
    uint32 stalenessSeconds,
    int256 fallbackWeiPerUnitLink,
    int256 fallbackUSDPerUnitEth,
    int256 fallbackUSDPerUnitLink,
    uint40 fulfillmentFlatFeeLinkUSD,
    uint40 fulfillmentFlatFeeEthUSD
  );

  ExposedPriceRegistry public s_registry;
  MockV3Aggregator s_linkEthAggregator;
  MockV3Aggregator s_linkUSDAggregator;
  MockV3Aggregator s_ethUSDAggregator;

  function setUp() public override {
    BaseTest.setUp();

    // 1 LINK == 0.01 eth
    s_linkEthAggregator = new MockV3Aggregator(
      18, 1e16);
    // 1 LINK == 10 USD
    s_linkUSDAggregator = new MockV3Aggregator(
      8, 10e8);
    // 1 ETH == 1000 USD
    s_ethUSDAggregator = new MockV3Aggregator(
      8, 1000e8);

    s_registry = new ExposedPriceRegistry(
      address(s_linkEthAggregator),
      address(s_linkUSDAggregator),
      address(s_ethUSDAggregator));
  }

  function testSetLinkEthFeed() public {
    address oldFeed = address(s_registry.s_linkETHFeed());
    MockV3Aggregator newFeed = new MockV3Aggregator(18, 1e16);
    vm.expectEmit(false, false, false, true /* data */);
    emit LinkEthFeedSet(oldFeed, address(newFeed));
    s_registry.setLINKETHFeed(address(newFeed));
  }

  function testSetLinkUsdFeed() public {
    address oldFeed = address(s_registry.s_linkUSDFeed());
    MockV3Aggregator newFeed = new MockV3Aggregator(8, 10e8);
    vm.expectEmit(false, false, false, true /* data */);
    emit LinkUSDFeedSet(oldFeed, address(newFeed));
    s_registry.setLINKUSDFeed(address(newFeed));
  }

  function testSetEthUsdFeed() public {
    address oldFeed = address(s_registry.s_ethUSDFeed());
    MockV3Aggregator newFeed = new MockV3Aggregator(8, 1000e8);
    vm.expectEmit(false, false, false, true /* data */);
    emit EthUSDFeedSet(oldFeed, address(newFeed));
    s_registry.setETHUSDFeed(address(newFeed));
  }

  function testSetConfig() public {
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 1e16; // 0.01 ether
    int256 fallbackUSDPerUnitEth = 1000e8; // 1000 USD
    int256 fallbackUSDPerUnitLink = 10e8; // 10 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    vm.expectEmit(false, false, false, true /* data */);
    emit ConfigSet(
      stalenessSeconds,
      fallbackWeiPerUnitLink,
      fallbackUSDPerUnitEth,
      fallbackUSDPerUnitLink,
      fulfillmentFlatFeeLinkUSD,
      fulfillmentFlatFeeEthUSD);
    s_registry.setConfig(
      stalenessSeconds,
      gasAfterPaymentCalculation,
      fallbackWeiPerUnitLink,
      fallbackUSDPerUnitEth,
      fallbackUSDPerUnitLink,
      fulfillmentFlatFeeLinkUSD,
      fulfillmentFlatFeeEthUSD);
    // check that config was indeed updated
    (uint32 cfgStalenessSeconds,
      uint32 cfgGasAfterPaymentCalculation,
      uint40 cfgFulfillmentFlatFeeLinkUSD,
      uint40 cfgFulfillmentFlatFeeEthUSD) = s_registry.s_config();
    assertEq(cfgStalenessSeconds, stalenessSeconds);
    assertEq(cfgGasAfterPaymentCalculation, gasAfterPaymentCalculation);
    assertEq(cfgFulfillmentFlatFeeEthUSD, fulfillmentFlatFeeEthUSD);
    assertEq(cfgFulfillmentFlatFeeLinkUSD, fulfillmentFlatFeeLinkUSD);
    // check that fallback prices were updated
    assertEq(s_registry.s_fallbackWeiPerUnitLink(), fallbackWeiPerUnitLink);
    assertEq(s_registry.s_fallbackUSDPerUnitEth(), fallbackUSDPerUnitEth);
    assertEq(s_registry.s_fallbackUSDPerUnitLink(), fallbackUSDPerUnitLink);
  }

  function testGetLinkEthFeedData() public {
    // set config so that the fallback is different than the price of the feed
    // this is so that we can tell if the fallback was fallen back to or not
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // LINK/ETH
    {
      (, int256 expectedAnswer, , , ) = s_linkEthAggregator.latestRoundData();
      int256 answer = s_registry.getLINKEthFeedDataExternal();
      assertEq(answer, expectedAnswer);
    }
  }

  function testGetLinkEthFeedDataStale() public {
    // set config so that the fallback is different than the price of the feed
    // this is so that we can tell if the fallback was fallen back to or not
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // warp block timestamp so that we get staleness being triggered
    vm.warp(block.timestamp + stalenessSeconds + 1);

    // LINK/ETH
    {
      int256 expectedAnswer = fallbackWeiPerUnitLink;
      int256 answer = s_registry.getLINKEthFeedDataExternal();
      assertEq(answer, expectedAnswer);
    }
  }

  function testGetUSDFeedData() public {
    // set config so that the fallback is different than the price of the feed
    // this is so that we can tell if the fallback was fallen back to or not
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // LINK/USD
    {
      uint8 expectedDecimals = s_linkUSDAggregator.decimals();
      (, int256 expectedAnswer, , , ) = s_linkUSDAggregator.latestRoundData();
      (int256 answer, uint8 decimals) = s_registry.getUSDFeedDataExternal(s_registry.s_linkUSDFeed());
      assertEq(decimals, expectedDecimals);
      assertEq(answer, expectedAnswer);
    }

    // ETH/USD
    {
      uint8 expectedDecimals = s_ethUSDAggregator.decimals();
      (, int256 expectedAnswer, , , ) = s_ethUSDAggregator.latestRoundData();
      (int256 answer, uint8 decimals) = s_registry.getUSDFeedDataExternal(s_registry.s_ethUSDFeed());
      assertEq(decimals, expectedDecimals);
      assertEq(answer, expectedAnswer);
    }

  }

  function testGetUSDFeedDataStale() public {
    // set config so that the fallback is different than the price of the feed
    // this is so that we can tell if the fallback was fallen back to or not
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // warp block timestamp so that we get staleness being triggered
    vm.warp(block.timestamp + stalenessSeconds + 1);

    // LINK/USD
    {
      uint8 expectedDecimals = s_linkUSDAggregator.decimals();
      int256 expectedAnswer = fallbackUSDPerUnitLink;
      (int256 answer, uint8 decimals) = s_registry.getUSDFeedDataExternal(s_registry.s_linkUSDFeed());
      assertEq(decimals, expectedDecimals);
      assertEq(answer, expectedAnswer, "fallbackUSDPerUnitLink");
    }

    // ETH/USD
    {
      uint8 expectedDecimals = s_ethUSDAggregator.decimals();
      int256 expectedAnswer = fallbackUSDPerUnitEth;
      (int256 answer, uint8 decimals) = s_registry.getUSDFeedDataExternal(s_registry.s_ethUSDFeed());
      assertEq(decimals, expectedDecimals);
      assertEq(answer, expectedAnswer, "fallbackUSDPerUnitEth");
    }
  }

  function testGetUSDFeedDataReverts() public {
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // any feed other than LINK/USD or ETH/USD should revert
    AggregatorV3Interface invalidFeed = s_registry.s_linkETHFeed();
    vm.expectRevert();
    s_registry.getUSDFeedDataExternal(invalidFeed);
  }

  function testCalculateFlatFeeFromUSDZeroPremium() public {
    // no set config needed for this test due to early check and return

    // LINK/USD
    {
      uint256 expectedFee = 0;
      uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(0, s_registry.s_linkUSDFeed());
      assertEq(fee, expectedFee);
    }

    // ETH/USD
    {
      uint256 expectedFee = 0;
      uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(0, s_registry.s_ethUSDFeed());
      assertEq(fee, expectedFee);
    }
  }

  function testCalculateFlatFeeFromUSDForLINKEqualDecimals() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // since 10 cents is the flat fee for LINK, and the conversion rate
    // in the aggregator is 1 LINK == 10 USD, the fee should be
    // 10 cents / 10 USD = 0.01 LINK, or 1e16 juels
    uint256 expectedFee = 1e16;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeLinkUSD, s_registry.s_linkUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForLINKEqualDecimalsMoreComplexPrice() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    MockV3Aggregator linkUsdAggregator = new MockV3Aggregator(8, 1521591000); // 15.21591 USD
    s_registry.setLINKUSDFeed(address(linkUsdAggregator));

    // since 10 cents is the flat fee for LINK, and the conversion rate
    // in the aggregator is 1 LINK == 15.21591 USD, the fee should be
    // 10 cents / 15.21591 USD = 6572068315335724 juels
    // or 0.006572068315335724 LINK in fixed point
    uint256 expectedFee = 6572068315335724;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeLinkUSD, s_registry.s_linkUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForEthEqualDecimalsMoreComplexPrice() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    MockV3Aggregator ethUsdAggregator = new MockV3Aggregator(8, 154354550000); // 1543.5455 USD
    s_registry.setETHUSDFeed(address(ethUsdAggregator));

    // since 10 cents is the flat fee for ETH, and the conversion rate
    // in the aggregator is 1 ETH == 1543.5455 USD, the fee should be
    // 10 cents / 1543.5455 USD = 64785910101127 wei
    // or 0.000064785910101127 ETH in fixed point
    uint256 expectedFee = 64785910101127;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeEthUSD, s_registry.s_ethUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForEthEqualDecimals() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // since 10 cents is the flat fee for ETH, and the conversion rate
    // in the aggregator is 1 ETH == 1000 USD, the fee should be
    // 10 cents / 1000 USD = 0.0001 ETH, or 1e13 wei
    uint256 expectedFee = 1e14;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeEthUSD, s_registry.s_ethUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForLINKFeedGreaterDecimals() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // set a LINK/USD feed with 18 decimals (for example)
    // to test out the case where feed decimals is greater than 8,
    // which are the configured decimals in the price registry contract
    MockV3Aggregator linkUsdAggregator = new MockV3Aggregator(18, 15215910000000000000); // 15.21591 USD
    s_registry.setLINKUSDFeed(address(linkUsdAggregator));

    // since 10 cents is the flat fee for LINK, and the conversion rate
    // in the aggregator is 1 LINK == 15.21591 USD, the fee should be
    // 10 cents / 15.21591 USD = 6572068315335724 juels
    // or 0.006572068315335724 LINK in fixed point
    // note that the number of decimals in the feed should not change this amount
    uint256 expectedFee = 6572068315335724;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeLinkUSD, s_registry.s_linkUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForEthFeedGreaterDecimals() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // set a ETH/USD feed with 18 decimals (for example)
    // to test out the case where feed decimals is greater than 8,
    // which are the configured decimals in the price registry contract
    MockV3Aggregator ethUsdAggregator = new MockV3Aggregator(18, 1543545500000000000000); // 1543.5455 USD
    s_registry.setETHUSDFeed(address(ethUsdAggregator));

    // since 10 cents is the flat fee for ETH, and the conversion rate
    // in the aggregator is 1 ETH == 1543.5455 USD, the fee should be
    // 10 cents / 1543.5455 USD = 64785910101127 wei
    // or 0.000064785910101127 ETH in fixed point
    // note that the number of decimals in the feed should not change this amount
    uint256 expectedFee = 64785910101127;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeEthUSD, s_registry.s_ethUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForLINKFeedLessDecimals() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // set a LINK/USD feed with 6 decimals (for example)
    // to test out the case where feed decimals is less than 8,
    // which are the configured decimals in the price registry contract
    MockV3Aggregator linkUsdAggregator = new MockV3Aggregator(6, 15215910); // 15.21591 USD
    s_registry.setLINKUSDFeed(address(linkUsdAggregator));

    // since 10 cents is the flat fee for LINK, and the conversion rate
    // in the aggregator is 1 LINK == 15.21591 USD, the fee should be
    // 10 cents / 15.21591 USD = 6572068315335724 juels
    // or 0.006572068315335724 LINK in fixed point
    // note that the number of decimals in the feed should not change this amount
    uint256 expectedFee = 6572068315335724;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeLinkUSD, s_registry.s_linkUSDFeed());
    assertEq(fee, expectedFee);
  }

  function testCalculateFlatFeeFromUSDForEthFeedLessDecimals() public {
    // set config to enable premium calculation
    // use a staleness seconds of 30 days
    uint32 stalenessSeconds = 30 * 24 * 60 * 60;
    uint32 gasAfterPaymentCalculation = 33_285;
    int256 fallbackWeiPerUnitLink = 2e16; // 0.02 ether
    int256 fallbackUSDPerUnitEth = 500e8; // 500 USD
    int256 fallbackUSDPerUnitLink = 5e8; // 5 USD
    uint40 fulfillmentFlatFeeLinkUSD = 1e7; // 10 cents
    uint40 fulfillmentFlatFeeEthUSD = 1e7; // 10 cents
    s_registry.setConfig(
      stalenessSeconds, // stalenessSeconds
      gasAfterPaymentCalculation, // gasAfterPaymentCalculation
      fallbackWeiPerUnitLink, // fallbackWeiPerUnitLink
      fallbackUSDPerUnitEth, // fallbackUSDPerUnitEth 500 USD
      fallbackUSDPerUnitLink, // fallbackUSDPerUnitLink 5 USD
      fulfillmentFlatFeeLinkUSD, // fulfillmentFlatFeeLinkUSD
      fulfillmentFlatFeeEthUSD // fulfillmentFlatFeeEthUSD
    );

    // set a ETH/USD feed with 6 decimals (for example)
    // to test out the case where feed decimals is less than 8,
    // which are the configured decimals in the price registry contract
    MockV3Aggregator ethUsdAggregator = new MockV3Aggregator(6, 1543545500); // 1543.5455 USD
    s_registry.setETHUSDFeed(address(ethUsdAggregator));

    // since 10 cents is the flat fee for ETH, and the conversion rate
    // in the aggregator is 1 ETH == 1543.5455 USD, the fee should be
    // 10 cents / 1543.5455 USD = 64785910101127 wei
    // or 0.000064785910101127 ETH in fixed point
    // note that the number of decimals in the feed should not change this amount
    uint256 expectedFee = 64785910101127;
    uint256 fee = s_registry.calculateFlatFeeFromUSDExternal(fulfillmentFlatFeeEthUSD, s_registry.s_ethUSDFeed());
    assertEq(fee, expectedFee);
  }
}
