pragma solidity ^0.8.0;

import {Test} from "forge-std/Test.sol";
import "../../../../src/v0.8/dev/HeartbeatRequester.sol";

contract HeartbeatRequesterSetup is Test {
  HeartbeatRequester s_requester;
  address internal constant OWNER = 0x00107e64E1Fb0c487F25dD6d3601FF6Af8d32E4e;
  address internal constant CRON_UPKEEP_ADDRESS = 0x1B4Ef6fAC5f6Be758508605Bd87356776BFbfba0;
  address internal constant INVALID_ADDRESS = 0x2B4Ef6FAC5F6BE758508605bd87356776BfbfbA0;
  address internal constant AGGREGATOR_PROXY = 0x132477835BC60C777321E1E09F1fE6cFb82E111c;
  address internal constant AGGREGATOR = 0x232477835BC60c777321E1e09F1fe6cFB82E111C;

  function setUp() public virtual {
    changePrank(OWNER);

    s_requester = new HeartbeatRequester();
  }
}

contract HeartbeatRequester_permitHeartbeat is HeartbeatRequesterSetup {
  event HeartbeatPermitted(address indexed permittedCaller, address indexed proxy);

  function testPermitHeartbeat() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.expectEmit(true, true, false, true);
    emit HeartbeatPermitted(CRON_UPKEEP_ADDRESS, AGGREGATOR_PROXY);
    s_requester.permitHeartbeat(CRON_UPKEEP_ADDRESS, IAggregatorProxy(AGGREGATOR_PROXY));
  }
}

contract HeartbeatRequester_removeHeartbeat is HeartbeatRequesterSetup {
  event HeartbeatRemoved(address indexed permittedCaller);

  function testRemoveHeartbeat() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.expectEmit(true, false, false, true);
    emit HeartbeatRemoved(CRON_UPKEEP_ADDRESS);
    s_requester.removeHeartbeat(CRON_UPKEEP_ADDRESS);
  }
}

contract HeartbeatRequester_getAggregatorAndRequestHeartbeat is HeartbeatRequesterSetup {

  function testInvalidCaller() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    s_requester.permitHeartbeat(CRON_UPKEEP_ADDRESS, IAggregatorProxy(AGGREGATOR_PROXY));
    vm.stopPrank();
    vm.startPrank(INVALID_ADDRESS);
    vm.expectRevert(HeartbeatRequester.HeartbeatNotPermitted.selector);
    s_requester.getAggregatorAndRequestHeartbeat(AGGREGATOR_PROXY);
  }

  function testGetAggregatorAndRequestHeartbeat() public {
    vm.stopPrank();
    vm.startPrank(OWNER);
    vm.mockCall(
      AGGREGATOR_PROXY,
      abi.encodeWithSelector(IAggregatorProxy.aggregator.selector),
      abi.encode(AGGREGATOR)
    );
    vm.mockCall(
      AGGREGATOR,
      abi.encodeWithSelector(IOffchainAggregator.requestNewRound.selector),
      abi.encode(1)
    );

    s_requester.permitHeartbeat(CRON_UPKEEP_ADDRESS, IAggregatorProxy(AGGREGATOR_PROXY));
    vm.stopPrank();
    vm.startPrank(CRON_UPKEEP_ADDRESS);
    s_requester.getAggregatorAndRequestHeartbeat(AGGREGATOR_PROXY);
  }
}
