// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {Deployer} from "./testhelpers/Deployer.t.sol";
import {AuthorizedForwarder} from "../AuthorizedForwarder.sol";

contract ForwarderTest is Deployer {
  AuthorizedForwarder internal s_forwarder;

  function setUp() public {
    _setUp();

    vm.prank(ALICE);
    s_forwarder = AuthorizedForwarder(s_factory.deployNewForwarder());
  }

  function test_SetAuthorizedSenders_Success() public {
    address[] memory senders;

    // Expect a revert when trying to set an empty list of authorized senders
    vm.expectRevert("Cannot set authorized senders");
    s_forwarder.setAuthorizedSenders(senders);

    vm.prank(ALICE);
    // Expect a revert because the sender list is empty
    vm.expectRevert("Must have at least 1 sender");
    s_forwarder.setAuthorizedSenders(senders);

    // Create a list with two identical sender addresses
    senders = new address[](2);
    senders[0] = SENDER_1;
    senders[1] = SENDER_1;

    vm.prank(ALICE);
    // Expect a revert because the sender list has duplicates
    vm.expectRevert("Must not have duplicate senders");
    s_forwarder.setAuthorizedSenders(senders);

    // Set the second sender to a different address
    senders[1] = SENDER_2;

    vm.prank(ALICE);
    // Update the authorized senders list
    s_forwarder.setAuthorizedSenders(senders);

    // Check if both SENDER_1 and SENDER_2 are now authorized
    assertTrue(s_forwarder.isAuthorizedSender(SENDER_1));
    assertTrue(s_forwarder.isAuthorizedSender(SENDER_2));

    // Fetch the authorized senders and verify they match the set addresses
    address[] memory returnedSenders = s_forwarder.getAuthorizedSenders();
    require(returnedSenders[0] == senders[0]);
    require(returnedSenders[1] == senders[1]);

    // Create a new list with only SENDER_3
    senders = new address[](1);
    senders[0] = SENDER_3;

    // Prank 'alice' and update the authorized senders to just SENDER_3
    vm.prank(ALICE);
    s_forwarder.setAuthorizedSenders(senders);

    // Ensure SENDER_1 and SENDER_2 are no longer authorized
    assertFalse(s_forwarder.isAuthorizedSender(SENDER_1));
    assertFalse(s_forwarder.isAuthorizedSender(SENDER_2));

    // Check that SENDER_3 is now the only authorized sender
    assertTrue(s_forwarder.isAuthorizedSender(SENDER_3));
    returnedSenders = s_forwarder.getAuthorizedSenders();
    require(returnedSenders[0] == senders[0]);
  }

  function testFuzz_Forward_Success(uint256 _value) public {
    _addSenders();

    vm.expectRevert("Not authorized sender");
    s_forwarder.forward(address(0), new bytes(0));

    vm.prank(SENDER_1);
    vm.expectRevert("Cannot forward to Link token");
    s_forwarder.forward(address(s_link), new bytes(0));

    vm.prank(SENDER_1);
    vm.expectRevert("Must forward to a contract");
    s_forwarder.forward(address(0), new bytes(0));

    vm.prank(SENDER_1);
    vm.expectRevert("Forwarded call reverted without reason");
    s_forwarder.forward(address(s_mockReceiver), new bytes(0));

    vm.prank(SENDER_1);
    vm.expectRevert("test revert message");
    s_forwarder.forward(address(s_mockReceiver), abi.encodeWithSignature("revertMessage()"));

    vm.prank(SENDER_1);
    s_forwarder.forward(address(s_mockReceiver), abi.encodeWithSignature("receiveData(uint256)", _value));

    require(s_mockReceiver.getValue() == _value);
  }

  function testFuzz_MultiForward_Success(uint256 _value1, uint256 _value2) public {
    _addSenders();

    address[] memory tos;
    bytes[] memory datas;

    vm.expectRevert("Not authorized sender");
    s_forwarder.multiForward(tos, datas);

    tos = new address[](2);
    datas = new bytes[](1);

    vm.prank(SENDER_1);
    vm.expectRevert("Arrays must have the same length");
    s_forwarder.multiForward(tos, datas);

    datas = new bytes[](2);

    vm.prank(SENDER_1);
    vm.expectRevert("Must forward to a contract");
    s_forwarder.multiForward(tos, datas);

    tos[0] = address(s_mockReceiver);
    tos[1] = address(s_link);

    vm.prank(SENDER_1);
    vm.expectRevert("Forwarded call reverted without reason");
    s_forwarder.multiForward(tos, datas);

    datas[0] = abi.encodeWithSignature("receiveData(uint256)", _value1);
    datas[1] = abi.encodeWithSignature("receiveData(uint256)", _value2);

    vm.prank(SENDER_1);
    vm.expectRevert("Cannot forward to Link token");
    s_forwarder.multiForward(tos, datas);

    tos[1] = address(s_mockReceiver);

    vm.prank(SENDER_1);
    s_forwarder.multiForward(tos, datas);

    require(s_mockReceiver.getValue() == _value2);
  }

  function test_OwnerForward_Success() public {
    vm.expectRevert("Only callable by owner");
    s_forwarder.ownerForward(address(0), new bytes(0));

    vm.prank(ALICE);
    vm.expectRevert("Forwarded call reverted without reason");
    s_forwarder.ownerForward(address(s_link), new bytes(0));

    vm.prank(ALICE);
    s_forwarder.ownerForward(address(s_link), abi.encodeWithSignature("balanceOf(address)", address(0)));
  }

  function test_TransferOwnershipWithMessage_Success() public {
    vm.prank(BOB);
    vm.expectRevert("Only callable by owner");
    s_forwarder.transferOwnershipWithMessage(BOB, new bytes(0));

    vm.prank(ALICE);
    s_forwarder.transferOwnershipWithMessage(BOB, new bytes(0));

    vm.expectRevert("Must be proposed owner");
    s_forwarder.acceptOwnership();

    vm.prank(BOB);
    s_forwarder.acceptOwnership();

    require(s_forwarder.owner() == BOB);
  }

  function _addSenders() internal {
    address[] memory senders = new address[](3);
    senders[0] = SENDER_1;
    senders[1] = SENDER_2;
    senders[2] = SENDER_3;

    vm.prank(ALICE);
    s_forwarder.setAuthorizedSenders(senders);
  }
}
