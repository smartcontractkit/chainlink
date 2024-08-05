// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {IReceiver} from "../interfaces/IReceiver.sol";
import {IRouter} from "../interfaces/IRouter.sol";
import {KeystoneForwarder} from "../KeystoneForwarder.sol";
import {Receiver} from "./mocks/Receiver.sol";

contract KeystoneRouter_SetConfigTest is Test {
  address internal ADMIN = address(1);
  address internal constant STRANGER = address(2);
  address internal constant FORWARDER = address(99);
  address internal constant TRANSMITTER = address(50);
  address internal constant RECEIVER = address(51);

  bytes internal metadata = hex"01020304";
  bytes internal report = hex"9998";
  bytes32 internal id = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";

  KeystoneForwarder internal s_router;
  Receiver internal s_receiver;

  function setUp() public virtual {
    vm.prank(ADMIN);
    s_router = new KeystoneForwarder();
    s_receiver = new Receiver();
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

  function test_RemoveForwarder_Success() public {
    vm.prank(ADMIN);
    vm.expectEmit(true, false, false, false);
    emit IRouter.ForwarderRemoved(FORWARDER);
    s_router.removeForwarder(FORWARDER);
  }

  function test_Route_RevertWhen_UnauthorizedForwarder() public {
    vm.prank(STRANGER);
    vm.expectRevert(IRouter.UnauthorizedForwarder.selector);
    s_router.route(id, TRANSMITTER, RECEIVER, metadata, report);
  }

  function test_Route_Success() public {
    assertEq(s_router.isForwarder(FORWARDER), false);

    vm.prank(ADMIN);
    s_router.addForwarder(FORWARDER);
    assertEq(s_router.isForwarder(FORWARDER), true);

    vm.prank(FORWARDER);
    vm.mockCall(address(s_receiver), abi.encodeCall(IReceiver.onReport, (metadata, report)), abi.encode());
    vm.expectCall(address(s_receiver), abi.encodeCall(IReceiver.onReport, (metadata, report)));
    s_router.route(id, TRANSMITTER, address(s_receiver), metadata, report);
  }
}
