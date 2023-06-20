// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "forge-std/Test.sol";
import {HeartbeatRequester, IAggregatorProxy, IOffchainAggregator} from "../HeartbeatRequester.sol";
import {MockAggregator} from "../mocks/MockAggregator.sol";
import {MockAggregatorProxy} from "../mocks/MockAggregatorProxy.sol";
import {MockOffchainAggregator} from "../mocks/MockOffchainAggregator.sol";

// from contracts directory,
// forge test --match-path src/v0.8/test/HeartbeatRequester.t.sol

contract HeartbeatRequesterSetUp is Test {
  HeartbeatRequester internal heartbeatRequester;
  MockAggregator internal aggregator;
  IAggregatorProxy internal aggregatorProxy;
  MockAggregator internal aggregator2;
  IAggregatorProxy internal aggregatorProxy2;
  address internal OWNER;
  address internal constant STRANGER = address(999);

  event HeartbeatPermitted(address indexed permittedCaller, address newProxy, address oldProxy);
  event HeartbeatRemoved(address indexed permittedCaller, address removedProxy);
  error HeartbeatNotPermitted();

  function setUp() public {
    OWNER = address(this);
    deal(OWNER, 1e20);
    vm.startPrank(OWNER);
    heartbeatRequester = new HeartbeatRequester();
    aggregator = new MockAggregator();
    aggregatorProxy = IAggregatorProxy(new MockAggregatorProxy(address(aggregator)));
    aggregator2 = new MockAggregator();
    aggregatorProxy2 = IAggregatorProxy(new MockAggregatorProxy(address(aggregator2)));
  }
}

contract HeartbeatRequester_permitHeartbeat is HeartbeatRequesterSetUp {
  function testBasicSuccess() public {
    vm.expectEmit();
    emit HeartbeatPermitted(STRANGER, address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(STRANGER, aggregatorProxy);

    vm.expectEmit();
    emit HeartbeatPermitted(STRANGER, address(aggregatorProxy2), address(aggregatorProxy));
    heartbeatRequester.permitHeartbeat(STRANGER, aggregatorProxy2);
  }

  function testBasicDeployerSuccess() public {
    vm.expectEmit();
    emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

    vm.expectEmit();
    emit HeartbeatPermitted(address(this), address(aggregatorProxy2), address(aggregatorProxy));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy2);
  }

  function testOnlyCallableByOwnerReverts() public {
    vm.expectRevert(bytes("Only callable by owner"));
    changePrank(STRANGER);
    heartbeatRequester.permitHeartbeat(STRANGER, aggregatorProxy);
  }
}

contract HeartbeatRequester_removeHeartbeat is HeartbeatRequesterSetUp {
  function testBasicSuccess() public {
    vm.expectEmit();
    emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);

    vm.expectEmit();
    emit HeartbeatRemoved(address(this), address(aggregatorProxy));
    heartbeatRequester.removeHeartbeat(address(this));
  }

  function testRemoveNoPermitsSuccess() public {
    vm.expectEmit();
    emit HeartbeatRemoved(address(this), address(0));
    heartbeatRequester.removeHeartbeat(address(this));
  }

  function testOnlyCallableByOwnerReverts() public {
    vm.expectRevert(bytes("Only callable by owner"));
    changePrank(STRANGER);
    heartbeatRequester.removeHeartbeat(address(this));
  }
}

contract HeartbeatRequester_getAggregatorRequestHeartbeat is HeartbeatRequesterSetUp {
  function testBasicSuxxess() public {
    vm.expectEmit();
    emit HeartbeatPermitted(address(this), address(aggregatorProxy), address(0));
    heartbeatRequester.permitHeartbeat(address(this), aggregatorProxy);
    heartbeatRequester.getAggregatorAndRequestHeartbeat(address(aggregatorProxy));
    // getter for newRoundCalled value
    bool val = aggregator.newRoundCalled();
    assertEq(val, true);
  }

  function testHeartbeatNotPermittedReverts() public {
    bytes32 hashedReason = keccak256(abi.encodePacked("HeartbeatNotPermitted()"));
    bytes memory revertMessage = bytes32ToBytes(hashedReason);
    vm.expectRevert(revertMessage);
    heartbeatRequester.getAggregatorAndRequestHeartbeat(address(aggregatorProxy));
    bool val = aggregator.newRoundCalled();
    assertFalse(val);
  }

  function bytes32ToBytes(bytes32 _bytes32) public pure returns (bytes memory) {
    bytes memory bytesArray = new bytes(4);
    for (uint256 i; i < 4; ++i) {
      bytesArray[i] = _bytes32[i];
    }
    return bytesArray;
  }
}