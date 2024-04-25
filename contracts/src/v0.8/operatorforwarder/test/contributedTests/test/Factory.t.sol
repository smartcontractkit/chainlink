// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";

import "../helpers/Deployer.sol";

contract FactoryTest is Deployer {

    function setUp() public {
        _setUp();
    }

    /**
     * @dev Test the deployment of a new operator.
     */
    function testDeployNewOperator() public {
        vm.prank(alice);
        // Deploy a new operator using the factory.
        address newOperator = factory.deployNewOperator();
        // Assert that the new operator was indeed created by the factory.
        assertTrue(factory.created(newOperator));
        // Ensure that Alice is the owner of the newly deployed operator.
        require(Operator(newOperator).owner() == alice);
    }

    /**
     * @dev Test the deployment of a new operator and a new forwarder.
     */
    function testDeployNewOperatorAndForwarder() public {
        vm.prank(alice);
        // Deploy both a new operator and a new forwarder using the factory.
        (address newOperator, address newForwarder) = factory.deployNewOperatorAndForwarder();

        // Assert that the new operator and the new forwarder were indeed created by the factory.
        assertTrue(factory.created(newOperator));
        assertTrue(factory.created(newForwarder));
        // Ensure that Alice is the owner of the newly deployed operator.
        require(Operator(newOperator).owner() == alice);

        //Operator to accept ownership
        vm.prank(newOperator);
        AuthorizedForwarder(newForwarder).acceptOwnership();

        // Ensure that the newly deployed operator is the owner of the newly deployed forwarder.
        require(AuthorizedForwarder(newForwarder).owner() == newOperator, "operator is not the owner");
    }

    /**
     * @dev Test the deployment of a new forwarder.
     */
    function testDeployNewForwarder() public {
        vm.prank(alice);
        // Deploy a new forwarder using the factory.
        address newForwarder = factory.deployNewForwarder();
        // Assert that the new forwarder was indeed created by the factory.
        assertTrue(factory.created(newForwarder));
        // Ensure that Alice is the owner of the newly deployed forwarder.
        require(AuthorizedForwarder(newForwarder).owner() == alice);
    }

    /**
     * @dev Test the deployment of a new forwarder and then transfer its ownership.
     */
    function testDeployNewForwarderAndTransferOwnership() public {
        vm.prank(alice);
        // Deploy a new forwarder with a proposal to transfer its ownership to Bob.
        address newForwarder = factory.deployNewForwarderAndTransferOwnership(bob, new bytes(0));
        // Assert that the new forwarder was indeed created by the factory.
        assertTrue(factory.created(newForwarder));
        // Ensure that Alice is still the current owner of the newly deployed forwarder.
        require(AuthorizedForwarder(newForwarder).owner() == alice);

        // Only proposed owner can call acceptOwnership()
        vm.expectRevert("Must be proposed owner");
        AuthorizedForwarder(newForwarder).acceptOwnership();
        
        vm.prank(bob);
        // Let Bob accept the ownership.
        AuthorizedForwarder(newForwarder).acceptOwnership();
        // Ensure that Bob is now the owner of the forwarder after the transfer.
        require(AuthorizedForwarder(newForwarder).owner() == bob);
    }

}