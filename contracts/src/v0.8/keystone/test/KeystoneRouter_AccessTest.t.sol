// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {IReceiver} from "../interfaces/IReceiver.sol";
import {KeystoneRouter} from "../KeystoneRouter.sol";

contract KeystoneRouter_SetConfigTest is Test {
  address internal ADMIN = address(1);
  address internal constant STRANGER = address(2);
  address internal constant FORWARDER = address(99);
  address internal constant TRANSMITTER = address(50);
  address internal constant RECEIVER = address(51);

  bytes internal metadata = hex"01020304";
  bytes internal report = hex"9998";
  bytes32 internal id = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";

  KeystoneRouter internal s_router;

  function setUp() public virtual {
    vm.prank(ADMIN);
    s_router = new KeystoneRouter();
  }

  function test_AddForwarder_RevertWhen_NotOwner() public {
    vm.prank(STRANGER);
    vm.expectRevert();
    s_router.addForwarder(FORWARDER);
  }

  function test_RemoveForwarder_RevertWhen_NotOwner() public {
    vm.prank(STRANGER);
    vm.expectRevert();
    s_router.removeForwarder(FORWARDER);
  }

  function test_Route_RevertWhen_Unauthorized() public {
    vm.prank(STRANGER);
    vm.expectRevert(KeystoneRouter.Unauthorized.selector);
    s_router.route(id, TRANSMITTER, RECEIVER, metadata, report);
  }

  function test_Route_Success() public {
    vm.prank(ADMIN);
    s_router.addForwarder(FORWARDER);

    vm.prank(FORWARDER);
    vm.mockCall(RECEIVER, abi.encodeCall(IReceiver.onReport, (metadata, report)), abi.encode());
    vm.expectCall(RECEIVER, abi.encodeCall(IReceiver.onReport, (metadata, report)));
    s_router.route(id, TRANSMITTER, RECEIVER, metadata, report);
  }
}
