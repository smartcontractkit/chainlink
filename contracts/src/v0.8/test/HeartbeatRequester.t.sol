// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import 'forge-std/Test.sol';
import {HeartbeatRequester} from "../HeartbeatRequester.sol";
import '../mocks/MockAggregator.sol';
import '../mocks/MockAggregatorProxy.sol';
import '../mocks/MockOffchainAggregator.sol';

    // from contracts directory,
    // forge test --match-path src/v0.8/test/HeartbeatRequester.t.sol

contract HeartbeatRequesterTest is Test {
    HeartbeatRequester public heartbeatRequester;
    MockAggregator public aggregator;
    IAggregatorProxy public aggregatorProxy;

    // MockAggregator public aggregator2;
    // IAggregatorProxy public aggregatorProxy2;

    address internal constant USER = address(2);

    event HeartbeatPermitted(address indexed permittedCaller, address newProxy, address oldProxy);
    event HeartbeatRemoved(address indexed permittedCaller, address removedProxy);
    error HeartbeatNotPermitted();

    function setUp() public {
        heartbeatRequester = new HeartbeatRequester();
        aggregator = new MockAggregator();
        aggregatorProxy = IAggregatorProxy(new MockAggregatorProxy(address(aggregator)));
    }

    function test_permitHeartbeat() public {
        // Allow address(this) to request heartbeats from aggregatorProxy
        vm.expectEmit(true, true, true, true);
        emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
        heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

        // Allow USER to request heartbeats from aggregatorProxy
        vm.expectEmit(true, true, true, true);
        emit HeartbeatPermitted(USER, address(aggregatorProxy), address(0));
        heartbeatRequester.permitHeartbeat(USER, aggregatorProxy);

    }

    function test_permitHeartbeat2() public {
        // Try adding permit from non-owner, fail is expected behavior
        vm.expectRevert(bytes("Only callable by owner"));
        vm.prank(USER);
        heartbeatRequester.permitHeartbeat(USER, aggregatorProxy);
    }

    function test_removeHeartbeat() public {
        // Allow address(this) to request heartbeats from aggregatorProxy
        // Remove address(this) from allowing to request heartbeats from its assigned Proxy
        vm.expectEmit(true, true, true, true);
        emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
        heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

        vm.expectEmit(true, true, true, true);
        emit HeartbeatRemoved(address(this), address(aggregatorProxy));
        heartbeatRequester.removeHeartbeat(address(this));
    }

    function test_removeHeartbeat2() public {
        // Try removing proxy from address(this) when address(this) has no permits
        vm.expectEmit(true, true, true, true);
        emit HeartbeatRemoved(address(this), address(0));
        heartbeatRequester.removeHeartbeat(address(this));
    }

    // REVERTING on getAggregatorAndRequestHeartbeat()'s call to aggregator.requestNewRound()
    // WHY? Need to fix
    function test_getAggregatorRequestHeartbeat() public {
        // getAggregatorRequestHeartbeat() has no return type, so how to test?

        vm.expectEmit(true, true, true, true);
        emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
        heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

        heartbeatRequester.getAggregatorAndRequestHeartbeat(address(aggregatorProxy));
    }

    function testRevert_getAggregatorRequestHeartbeat() public {
        // getAggregatorRequestHeartbeat() has no return type, so how to test?

        vm.expectRevert();
        heartbeatRequester.getAggregatorAndRequestHeartbeat(address(aggregatorProxy)); 
    }

}