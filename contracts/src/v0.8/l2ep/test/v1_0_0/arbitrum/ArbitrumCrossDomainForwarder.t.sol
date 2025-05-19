// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ArbitrumCrossDomainForwarder} from "../../../arbitrum/ArbitrumCrossDomainForwarder.sol";
import {Greeter} from "../../Greeter.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ArbitrumCrossDomainForwarderTest is L2EPTest {
  /// Helper variable(s)
  address internal s_crossDomainMessengerAddr = toArbitrumL2AliasAddress(s_l1OwnerAddr);
  address internal s_newOwnerCrossDomainMessengerAddr = toArbitrumL2AliasAddress(s_strangerAddr);

  /// Contracts
  ArbitrumCrossDomainForwarder internal s_arbitrumCrossDomainForwarder;
  Greeter internal s_greeter;

  /// Events
  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    vm.startPrank(s_l1OwnerAddr);
    s_arbitrumCrossDomainForwarder = new ArbitrumCrossDomainForwarder(s_l1OwnerAddr);
    s_greeter = new Greeter(address(s_arbitrumCrossDomainForwarder));
    vm.stopPrank();
  }
}

contract ArbitrumCrossDomainForwarder_Constructor is ArbitrumCrossDomainForwarderTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialState() public {
    // it should set the owner correctly
    assertEq(s_arbitrumCrossDomainForwarder.owner(), s_l1OwnerAddr);

    // it should set the l1Owner correctly
    assertEq(s_arbitrumCrossDomainForwarder.l1Owner(), s_l1OwnerAddr);

    // it should set the crossdomain messenger correctly
    assertEq(s_arbitrumCrossDomainForwarder.crossDomainMessenger(), s_crossDomainMessengerAddr);

    // it should set the typeAndVersion correctly
    assertEq(s_arbitrumCrossDomainForwarder.typeAndVersion(), "ArbitrumCrossDomainForwarder 1.0.0");
  }
}

contract ArbitrumCrossDomainForwarder_Forward is ArbitrumCrossDomainForwarderTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_arbitrumCrossDomainForwarder.forward(address(s_greeter), abi.encode(""));
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_Forward() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_arbitrumCrossDomainForwarder.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, greeting)
    );

    // Checks that the greeter got the message
    assertEq(s_greeter.greeting(), greeting);
  }

  /// @notice it should revert when contract call reverts
  function test_ForwardRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr);

    // Sends an invalid message
    vm.expectRevert("Invalid greeting length");
    s_arbitrumCrossDomainForwarder.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, "")
    );
  }
}

contract ArbitrumCrossDomainForwarder_TransferL1Ownership is ArbitrumCrossDomainForwarderTest {
  /// @notice it should not be callable by non-owners
  function test_NotCallableByNonOwners() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_arbitrumCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should not be callable by L2 owner
  function test_NotCallableByL2Owner() public {
    vm.startPrank(s_l1OwnerAddr);
    assertEq(s_arbitrumCrossDomainForwarder.owner(), s_l1OwnerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_arbitrumCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should be callable by current L1 owner
  function test_CallableByL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr);

    // Defines the cross domain message to send
    vm.expectEmit();
    emit L1OwnershipTransferRequested(s_arbitrumCrossDomainForwarder.l1Owner(), s_strangerAddr);

    // Sends the message
    s_arbitrumCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should be callable by current L1 owner to zero address
  function test_CallableByL1OwnerOrZeroAddress() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr);

    // Defines the cross domain message to send
    vm.expectEmit();
    emit L1OwnershipTransferRequested(s_arbitrumCrossDomainForwarder.l1Owner(), address(0));

    // Sends the message
    s_arbitrumCrossDomainForwarder.transferL1Ownership(address(0));
  }
}

contract ArbitrumCrossDomainForwarder_AcceptL1Ownership is ArbitrumCrossDomainForwarderTest {
  /// @notice it should not be callable by non pending-owners
  function test_NotCallableByNonPendingOwners() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr);

    // Sends the message
    vm.expectRevert("Must be proposed L1 owner");
    s_arbitrumCrossDomainForwarder.acceptL1Ownership();
  }

  /// @notice it should be callable by pending L1 owner
  function test_CallableByPendingL1Owner() public {
    // Request ownership transfer
    vm.startPrank(s_crossDomainMessengerAddr);
    s_arbitrumCrossDomainForwarder.transferL1Ownership(s_strangerAddr);

    // Prepares expected event payload
    vm.expectEmit();
    emit L1OwnershipTransferred(s_l1OwnerAddr, s_strangerAddr);

    // Accepts ownership transfer request
    vm.startPrank(s_newOwnerCrossDomainMessengerAddr);
    s_arbitrumCrossDomainForwarder.acceptL1Ownership();

    // Asserts that the ownership was actually transferred
    assertEq(s_arbitrumCrossDomainForwarder.l1Owner(), s_strangerAddr);
  }
}
