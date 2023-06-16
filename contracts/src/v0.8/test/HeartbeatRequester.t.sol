// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "forge-std/Test.sol";
import {HeartbeatRequester} from "../HeartbeatRequester.sol";
import "../mocks/MockAggregator.sol";
import "../mocks/MockAggregatorProxy.sol";
import "../mocks/MockOffchainAggregator.sol";

// from contracts directory,
// forge test --match-path src/v0.8/test/HeartbeatRequester.t.sol

contract HeartbeatRequesterTest is Test {
  HeartbeatRequester public heartbeatRequester;
  MockAggregator public aggregator;
  IAggregatorProxy public aggregatorProxy;
  MockAggregator public aggregator2;
  IAggregatorProxy public aggregatorProxy2;

  address internal constant USER = address(2);

  event HeartbeatPermitted(address indexed permittedCaller, address newProxy, address oldProxy);
  event HeartbeatRemoved(address indexed permittedCaller, address removedProxy);
  error HeartbeatNotPermitted();

  function setUp() public {
    heartbeatRequester = new HeartbeatRequester();
    aggregator = new MockAggregator();
    aggregatorProxy = IAggregatorProxy(new MockAggregatorProxy(address(aggregator)));
    aggregator2 = new MockAggregator();
    aggregatorProxy2 = IAggregatorProxy(new MockAggregatorProxy(address(aggregator2)));
  }

  function test_permitHeartbeat_ArbitraryUser() public {
    vm.expectEmit(true, true, true, true);
    emit HeartbeatPermitted(USER, address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(USER, aggregatorProxy);

    vm.expectEmit(true, true, true, true);
    emit HeartbeatPermitted(USER, address(aggregatorProxy2), address(aggregatorProxy));
    heartbeatRequester.permitHeartbeat(USER, aggregatorProxy2);
  }

  function test_permitHeartbeat_Deployer() public {
    vm.expectEmit(true, true, true, true);
    emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

    vm.expectEmit(true, true, true, true);
    emit HeartbeatPermitted(address(this), address(aggregatorProxy2), address(aggregatorProxy));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy2);
  }

  function test_permitHeartbeat_NotFromOwner() public {
    vm.expectRevert(bytes("Only callable by owner"));
    vm.prank(USER);
    heartbeatRequester.permitHeartbeat(USER, aggregatorProxy);
  }

  function test_removeHeartbeat() public {
    vm.expectEmit(true, true, true, true);
    emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

    vm.expectEmit(true, true, true, true);
    emit HeartbeatRemoved(address(this), address(aggregatorProxy));
    heartbeatRequester.removeHeartbeat(address(this));
  }

  function test_removeHeartbeat_WithNoPermits() public {
    vm.expectEmit(true, true, true, true);
    emit HeartbeatRemoved(address(this), address(0));
    heartbeatRequester.removeHeartbeat(address(this));
  }

  function test_removeHeartbeat_NotFromOwner() public {
    vm.expectRevert(bytes("Only callable by owner"));
    vm.prank(USER);
    heartbeatRequester.removeHeartbeat(address(this));
  }

  function test_getAggregatorRequestHeartbeat() public {
    vm.expectEmit(true, true, true, true);
    emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);
    heartbeatRequester.getAggregatorAndRequestHeartbeat(address(aggregatorProxy));
    // getter for newRoundCalled value
    bool val = aggregator.newRoundCalled();
    assertEq(val, true);
  }

  function testRevert_getAggregatorRequestHeartbeat() public {
    bytes32 hashedReason = keccak256(abi.encodePacked("HeartbeatNotPermitted()"));
    bytes memory revertMessage = bytes32ToBytes(hashedReason);
    vm.expectRevert(revertMessage);
    heartbeatRequester.getAggregatorAndRequestHeartbeat(address(aggregatorProxy));
    bool val = aggregator.newRoundCalled();
    assertFalse(val);
  }

  function bytes32ToBytes(bytes32 _bytes32) public pure returns (bytes memory) {
    bytes memory bytesArray = new bytes(4);
    for (uint256 i; i < 4; i++) {
      bytesArray[i] = _bytes32[i];
    }
    return bytesArray;
  }
}
