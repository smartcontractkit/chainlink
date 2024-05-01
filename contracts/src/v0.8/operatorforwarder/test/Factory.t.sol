// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "forge-std/Test.sol";

import "./testhelpers/Deployer.t.sol";

contract FactoryTest is Deployer {
  function setUp() public {
    _setUp();

    vm.startPrank(s_alice);
  }

  /**
   * @dev Test the deployment of a new operator.
   */
  function test_DeployNewOperator() public {
    // Deploy a new operator using the factory.
    address newOperator = s_factory.deployNewOperator();
    // Assert that the new operator was indeed created by the factory.
    assertTrue(s_factory.created(newOperator));
    // Ensure that Alice is the owner of the newly deployed operator.
    require(Operator(newOperator).owner() == s_alice);
  }

  /**
   * @dev Test the deployment of a new operator and a new forwarder.
   */
  function test_DeployNewOperatorAndForwarder() public {
    // Deploy both a new operator and a new forwarder using the factory.
    (address newOperator, address newForwarder) = s_factory.deployNewOperatorAndForwarder();

    // Assert that the new operator and the new forwarder were indeed created by the factory.
    assertTrue(s_factory.created(newOperator));
    assertTrue(s_factory.created(newForwarder));
    // Ensure that Alice is the owner of the newly deployed operator.
    require(Operator(newOperator).owner() == s_alice);

    //Operator to accept ownership
    vm.startPrank(newOperator);
    AuthorizedForwarder(newForwarder).acceptOwnership();

    // Ensure that the newly deployed operator is the owner of the newly deployed forwarder.
    require(AuthorizedForwarder(newForwarder).owner() == newOperator, "operator is not the owner");
  }

  /**
   * @dev Test the deployment of a new forwarder.
   */
  function test_DeployNewForwarder() public {
    // Deploy a new forwarder using the factory.
    address newForwarder = s_factory.deployNewForwarder();
    // Assert that the new forwarder was indeed created by the factory.
    assertTrue(s_factory.created(newForwarder));
    // Ensure that Alice is the owner of the newly deployed forwarder.
    require(AuthorizedForwarder(newForwarder).owner() == s_alice);
  }

  /**
   * @dev Test the deployment of a new forwarder and then transfer its ownership.
   */
  function test_DeployNewForwarderAndTransferOwnership() public {
    // Deploy a new forwarder with a proposal to transfer its ownership to Bob.
    address newForwarder = s_factory.deployNewForwarderAndTransferOwnership(s_bob, new bytes(0));
    // Assert that the new forwarder was indeed created by the factory.
    assertTrue(s_factory.created(newForwarder));
    // Ensure that Alice is still the current owner of the newly deployed forwarder.
    require(AuthorizedForwarder(newForwarder).owner() == s_alice);

    // Only proposed owner can call acceptOwnership()
    vm.expectRevert("Must be proposed owner");
    AuthorizedForwarder(newForwarder).acceptOwnership();

    vm.startPrank(s_bob);
    // Let Bob accept the ownership.
    AuthorizedForwarder(newForwarder).acceptOwnership();
    // Ensure that Bob is now the owner of the forwarder after the transfer.
    require(AuthorizedForwarder(newForwarder).owner() == s_bob);
  }
}
