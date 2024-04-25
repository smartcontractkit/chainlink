// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";

import "../helpers/Deployer.sol";

contract ForwarderTest is Deployer {

    function setUp() public {
        _setUp();

        vm.prank(alice);
        forwarder = AuthorizedForwarder(factory.deployNewForwarder());
    }

    /**
     * @dev Tests the functionality of setting authorized senders.
     */
    function testSetAuthorizedSenders() public {
        address[] memory senders;

        // Expect a revert when trying to set an empty list of authorized senders
        vm.expectRevert("Cannot set authorized senders");
        forwarder.setAuthorizedSenders(senders);

        vm.prank(alice);
        // Expect a revert because the sender list is empty
        vm.expectRevert("Must have at least 1 sender");
        forwarder.setAuthorizedSenders(senders);

        // Create a list with two identical sender addresses
        senders = new address[](2);
        senders[0] = sender1;
        senders[1] = sender1;

        vm.prank(alice);
        // Expect a revert because the sender list has duplicates
        vm.expectRevert("Must not have duplicate senders");
        forwarder.setAuthorizedSenders(senders);

        // Set the second sender to a different address
        senders[1] = sender2;
        
        vm.prank(alice);
        // Update the authorized senders list
        forwarder.setAuthorizedSenders(senders);

        // Check if both sender1 and sender2 are now authorized
        assertTrue(forwarder.isAuthorizedSender(sender1));
        assertTrue(forwarder.isAuthorizedSender(sender2));

        // Fetch the authorized senders and verify they match the set addresses
        address[] memory returnedSenders = forwarder.getAuthorizedSenders();
        require(returnedSenders[0] == senders[0]);
        require(returnedSenders[1] == senders[1]);

        // Create a new list with only sender3
        senders = new address[](1);
        senders[0] = sender3;

        // Prank 'alice' and update the authorized senders to just sender3
        vm.prank(alice);
        forwarder.setAuthorizedSenders(senders);

        // Ensure sender1 and sender2 are no longer authorized
        assertFalse(forwarder.isAuthorizedSender(sender1));
        assertFalse(forwarder.isAuthorizedSender(sender2));

        // Check that sender3 is now the only authorized sender
        assertTrue(forwarder.isAuthorizedSender(sender3));
        returnedSenders = forwarder.getAuthorizedSenders();
        require(returnedSenders[0] == senders[0]);
    }

    /**
     * @dev Tests the behavior of single forward
     */
    function testForward(uint256 _value) public {
        _addSenders();

        vm.expectRevert("Not authorized sender");
        forwarder.forward(address(0), new bytes(0));

        vm.prank(sender1);
        vm.expectRevert("Cannot forward to Link token");
        forwarder.forward(address(link), new bytes(0));

        vm.prank(sender1);
        vm.expectRevert("Must forward to a contract");
        forwarder.forward(address(0), new bytes(0));

        vm.prank(sender1);
        vm.expectRevert("Forwarded call reverted without reason");
        forwarder.forward(address(mockReceiver), new bytes(0));

        vm.prank(sender1);
        vm.expectRevert("test revert message");
        forwarder.forward(address(mockReceiver), abi.encodeWithSignature("revertMessage()"));

        vm.prank(sender1);
        forwarder.forward(address(mockReceiver), abi.encodeWithSignature("receiveData(uint256)", _value));

        require(mockReceiver.value() == _value);
    }

    function testMultiForward(uint256 _value1, uint256 _value2) public {
        _addSenders();

        address[] memory tos;
        bytes[] memory datas;

        vm.expectRevert("Not authorized sender");
        forwarder.multiForward(tos, datas);

        tos = new address[](2);
        datas = new bytes[](1);

        vm.prank(sender1);
        vm.expectRevert("Arrays must have the same length");
        forwarder.multiForward(tos, datas);

        datas = new bytes[](2);

        vm.prank(sender1);
        vm.expectRevert("Must forward to a contract");
        forwarder.multiForward(tos, datas);

        tos[0] = address(mockReceiver);
        tos[1] = address(link);

        vm.prank(sender1);
        vm.expectRevert("Forwarded call reverted without reason");
        forwarder.multiForward(tos, datas);
        
        datas[0] = abi.encodeWithSignature("receiveData(uint256)", _value1);
        datas[1] = abi.encodeWithSignature("receiveData(uint256)", _value2);

        vm.prank(sender1);
        vm.expectRevert("Cannot forward to Link token");
        forwarder.multiForward(tos, datas);
        
        tos[1] = address(mockReceiver);
        
        vm.prank(sender1);
        forwarder.multiForward(tos, datas);

        require(mockReceiver.value() == _value2);
    }

    /**
     * @dev tests the difference between ownerForward and forward
     * specifically owner can forward to link token
     */
    function testOwnerForward() public {
        vm.expectRevert("Only callable by owner");
        forwarder.ownerForward(address(0), new bytes(0));

        vm.prank(alice);
        vm.expectRevert("Forwarded call reverted without reason");
        forwarder.ownerForward(address(link), new bytes(0));

        vm.prank(alice);
        forwarder.ownerForward(address(link), abi.encodeWithSignature("balanceOf(address)", address(0)));
    }

    /**
     * @dev Tests the behavior of transfer and accept ownership of the contract.
     */
    function testTransferOwnershipWithMessage() public {
        vm.prank(bob);
        vm.expectRevert("Only callable by owner");
        forwarder.transferOwnershipWithMessage(bob, new bytes(0));

        vm.prank(alice);
        forwarder.transferOwnershipWithMessage(bob, new bytes(0));

        vm.expectRevert("Must be proposed owner");
        forwarder.acceptOwnership();

        vm.prank(bob);
        forwarder.acceptOwnership();

        require(forwarder.owner() == bob);
    }

    /**
     * @dev Helper function to setup senders 
     */
    function _addSenders() internal {
        address[] memory senders = new address[](3);
        senders[0] = sender1;
        senders[1] = sender2;
        senders[2] = sender3;

        vm.prank(alice);
        forwarder.setAuthorizedSenders(senders);
    }

}