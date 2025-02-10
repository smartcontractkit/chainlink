// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {BundleAggregatorProxy} from "../BundleAggregatorProxy.sol";
import {DataFeedsCache} from "../DataFeedsCache.sol";
import {BaseTest} from "./BaseTest.t.sol";

contract BundleAggregatorProxyTest is BaseTest {
  BundleAggregatorProxy internal proxy;
  DataFeedsCache internal aggregator;

  function setUp() public override {
    super.setUp();
    aggregator = new DataFeedsCache();
    proxy = new BundleAggregatorProxy(address(aggregator), OWNER);

    bytes16[] memory datIds = new bytes16[](1);
    datIds[0] = bytes16("1");

    address[] memory proxies = new address[](1);
    proxies[0] = address(proxy);

    aggregator.setFeedAdmin(OWNER, true);
    aggregator.updateDataIdMappingsForProxies(proxies, datIds);
  }

  function test_aggregator() public {
    assertEq(proxy.aggregator(), address(aggregator));
  }

  function test_version() public {
    assertEq(proxy.version(), 7);
  }

  function test_description() public {
    assertEq(proxy.description(), "");
  }

  function test_latestBundle() public {
    bytes memory bundle = proxy.latestBundle();
    assertEq(bundle.length, 0);
  }

  function test_latestBundleTimestamp() public {
    assertEq(proxy.latestBundleTimestamp(), 0);
  }

  function test_bundleDecimals() public {
    uint8[] memory decimals = proxy.bundleDecimals();
    assertEq(decimals.length, 0);
  }

  function test_proposeAggregator() public {
    address newAggregator = address(123);
    vm.expectEmit();
    emit BundleAggregatorProxy.AggregatorProposed({current: address(aggregator), proposed: newAggregator});
    proxy.proposeAggregator(newAggregator);

    assertEq(proxy.proposedAggregator(), newAggregator);
  }

  function test_confirmAggregatorRevertNotProposed() public {
    address newAggregator = address(123);
    vm.expectRevert(abi.encodeWithSelector(BundleAggregatorProxy.AggregatorNotProposed.selector, newAggregator));
    proxy.confirmAggregator(newAggregator);
  }

  function test_confirmAggregatorSuccess() public {
    address newAggregator = address(123);
    proxy.proposeAggregator(newAggregator);
    vm.expectEmit();
    emit BundleAggregatorProxy.AggregatorConfirmed({previous: address(aggregator), latest: newAggregator});
    proxy.confirmAggregator(newAggregator);
  }
}
