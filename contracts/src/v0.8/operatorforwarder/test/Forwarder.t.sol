// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";

import "./testhelpers/Deployer.t.sol";
import {AuthorizedForwarder} from "../AuthorizedForwarder.sol";

contract ForwarderTest is Deployer {
  AuthorizedForwarder internal s_forwarder;

  function setUp() public {
    _setUp();

    vm.prank(s_alice);
    s_forwarder = AuthorizedForwarder(s_factory.deployNewForwarder());
  }

  /**
   * @dev Tests the functionality of setting authorized senders.
   */
  function test_SetAuthorizedSenders() public {
    address[] memory senders;

    // Expect a revert when trying to set an empty list of authorized senders
    vm.expectRevert("Cannot set authorized senders");
    s_forwarder.setAuthorizedSenders(senders);

    vm.prank(s_alice);
    // Expect a revert because the sender list is empty
    vm.expectRevert("Must have at least 1 sender");
    s_forwarder.setAuthorizedSenders(senders);

    // Create a list with two identical sender addresses
    senders = new address[](2);
    senders[0] = s_sender1;
    senders[1] = s_sender1;

    vm.prank(s_alice);
    // Expect a revert because the sender list has duplicates
    vm.expectRevert("Must not have duplicate senders");
    s_forwarder.setAuthorizedSenders(senders);

    // Set the second sender to a different address
    senders[1] = s_sender2;

    vm.prank(s_alice);
    // Update the authorized senders list
    s_forwarder.setAuthorizedSenders(senders);

    // Check if both s_sender1 and s_sender2 are now authorized
    assertTrue(s_forwarder.isAuthorizedSender(s_sender1));
    assertTrue(s_forwarder.isAuthorizedSender(s_sender2));

    // Fetch the authorized senders and verify they match the set addresses
    address[] memory returnedSenders = s_forwarder.getAuthorizedSenders();
    require(returnedSenders[0] == senders[0]);
    require(returnedSenders[1] == senders[1]);

    // Create a new list with only s_sender3
    senders = new address[](1);
    senders[0] = s_sender3;

    // Prank 'alice' and update the authorized senders to just s_sender3
    vm.prank(s_alice);
    s_forwarder.setAuthorizedSenders(senders);

    // Ensure s_sender1 and s_sender2 are no longer authorized
    assertFalse(s_forwarder.isAuthorizedSender(s_sender1));
    assertFalse(s_forwarder.isAuthorizedSender(s_sender2));

    // Check that s_sender3 is now the only authorized sender
    assertTrue(s_forwarder.isAuthorizedSender(s_sender3));
    returnedSenders = s_forwarder.getAuthorizedSenders();
    require(returnedSenders[0] == senders[0]);
  }

  /**
   * @dev Tests the behavior of single forward
   */
  function test_Forward(uint256 _value) public {
    _addSenders();

    vm.expectRevert("Not authorized sender");
    s_forwarder.forward(address(0), new bytes(0));

    vm.prank(s_sender1);
    vm.expectRevert("Cannot forward to Link token");
    s_forwarder.forward(address(s_link), new bytes(0));

    vm.prank(s_sender1);
    vm.expectRevert("Must forward to a contract");
    s_forwarder.forward(address(0), new bytes(0));

    vm.prank(s_sender1);
    vm.expectRevert("Forwarded call reverted without reason");
    s_forwarder.forward(address(s_mockReceiver), new bytes(0));

    vm.prank(s_sender1);
    vm.expectRevert("test revert message");
    s_forwarder.forward(address(s_mockReceiver), abi.encodeWithSignature("revertMessage()"));

    vm.prank(s_sender1);
    s_forwarder.forward(address(s_mockReceiver), abi.encodeWithSignature("receiveData(uint256)", _value));

    require(s_mockReceiver.getValue() == _value);
  }

  function test_MultiForward(uint256 _value1, uint256 _value2) public {
    _addSenders();

    address[] memory tos;
    bytes[] memory datas;

    vm.expectRevert("Not authorized sender");
    s_forwarder.multiForward(tos, datas);

    tos = new address[](2);
    datas = new bytes[](1);

    vm.prank(s_sender1);
    vm.expectRevert("Arrays must have the same length");
    s_forwarder.multiForward(tos, datas);

    datas = new bytes[](2);

    vm.prank(s_sender1);
    vm.expectRevert("Must forward to a contract");
    s_forwarder.multiForward(tos, datas);

    tos[0] = address(s_mockReceiver);
    tos[1] = address(s_link);

    vm.prank(s_sender1);
    vm.expectRevert("Forwarded call reverted without reason");
    s_forwarder.multiForward(tos, datas);

    datas[0] = abi.encodeWithSignature("receiveData(uint256)", _value1);
    datas[1] = abi.encodeWithSignature("receiveData(uint256)", _value2);

    vm.prank(s_sender1);
    vm.expectRevert("Cannot forward to Link token");
    s_forwarder.multiForward(tos, datas);

    tos[1] = address(s_mockReceiver);

    vm.prank(s_sender1);
    s_forwarder.multiForward(tos, datas);

    require(s_mockReceiver.getValue() == _value2);
  }

  /**
   * @dev tests the difference between ownerForward and forward
   * specifically owner can forward to link token
   */
  function test_OwnerForward() public {
    vm.expectRevert("Only callable by owner");
    s_forwarder.ownerForward(address(0), new bytes(0));

    vm.prank(s_alice);
    vm.expectRevert("Forwarded call reverted without reason");
    s_forwarder.ownerForward(address(s_link), new bytes(0));

    vm.prank(s_alice);
    s_forwarder.ownerForward(address(s_link), abi.encodeWithSignature("balanceOf(address)", address(0)));
  }

  /**
   * @dev Tests the behavior of transfer and accept ownership of the contract.
   */
  function test_TransferOwnershipWithMessage() public {
    vm.prank(s_bob);
    vm.expectRevert("Only callable by owner");
    s_forwarder.transferOwnershipWithMessage(s_bob, new bytes(0));

    vm.prank(s_alice);
    s_forwarder.transferOwnershipWithMessage(s_bob, new bytes(0));

    vm.expectRevert("Must be proposed owner");
    s_forwarder.acceptOwnership();

    vm.prank(s_bob);
    s_forwarder.acceptOwnership();

    require(s_forwarder.owner() == s_bob);
  }

  /**
   * @dev Helper function to setup senders
   */
  function _addSenders() internal {
    address[] memory senders = new address[](3);
    senders[0] = s_sender1;
    senders[1] = s_sender2;
    senders[2] = s_sender3;

    vm.prank(s_alice);
    s_forwarder.setAuthorizedSenders(senders);
  }
}
