// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {BundleAggregatorProxy} from "../BundleAggregatorProxy.sol";
import {DataFeedsCache} from "../DataFeedsCache.sol";
import {BaseTest} from "./BaseTest.t.sol";

contract BundleAggregatorProxyTest is BaseTest {
  BundleAggregatorProxy internal s_proxy;
  DataFeedsCache internal s_aggregator;

  function setUp() public override {
    super.setUp();
    s_aggregator = new DataFeedsCache();
    s_proxy = new BundleAggregatorProxy(address(s_aggregator), OWNER);

    bytes16[] memory datIds = new bytes16[](1);
    datIds[0] = bytes16("1");

    address[] memory proxies = new address[](1);
    proxies[0] = address(s_proxy);

    s_aggregator.setFeedAdmin(OWNER, true);
    s_aggregator.updateDataIdMappingsForProxies(proxies, datIds);
  }

  function test_aggregator() public {
    assertEq(s_proxy.aggregator(), address(s_aggregator));
  }

  function test_version() public {
    assertEq(s_proxy.version(), 7);
  }

  function test_description() public {
    assertEq(s_proxy.description(), "");
  }

  function test_latestBundle() public {
    bytes memory bundle = s_proxy.latestBundle();
    assertEq(bundle.length, 0);
  }

  function test_latestBundleTimestamp() public {
    assertEq(s_proxy.latestBundleTimestamp(), 0);
  }

  function test_bundleDecimals() public {
    uint8[] memory decimals = s_proxy.bundleDecimals();
    assertEq(decimals.length, 0);
  }

  function test_proposeAggregator() public {
    address newAggregator = address(123);
    vm.expectEmit();
    emit BundleAggregatorProxy.AggregatorProposed({current: address(s_aggregator), proposed: newAggregator});
    s_proxy.proposeAggregator(newAggregator);

    assertEq(s_proxy.proposedAggregator(), newAggregator);
  }

  function test_confirmAggregatorRevertNotProposed() public {
    address newAggregator = address(123);
    vm.expectRevert(abi.encodeWithSelector(BundleAggregatorProxy.AggregatorNotProposed.selector, newAggregator));
    s_proxy.confirmAggregator(newAggregator);
  }

  function test_confirmAggregatorSuccess() public {
    address newAggregator = address(123);
    s_proxy.proposeAggregator(newAggregator);
    vm.expectEmit();
    emit BundleAggregatorProxy.AggregatorConfirmed({previous: address(s_aggregator), latest: newAggregator});
    s_proxy.confirmAggregator(newAggregator);
  }
}
