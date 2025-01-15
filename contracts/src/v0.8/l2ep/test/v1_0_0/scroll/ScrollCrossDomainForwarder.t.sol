// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockScrollCrossDomainMessenger} from "../../mocks/scroll/MockScrollCrossDomainMessenger.sol";
import {ScrollCrossDomainForwarder} from "../../../scroll/ScrollCrossDomainForwarder.sol";
import {Greeter} from "../../Greeter.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ScrollCrossDomainForwarderTest is L2EPTest {
  /// Contracts
  MockScrollCrossDomainMessenger internal s_mockScrollCrossDomainMessenger;
  ScrollCrossDomainForwarder internal s_scrollCrossDomainForwarder;
  Greeter internal s_greeter;

  /// Events
  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    vm.startPrank(s_l1OwnerAddr);
    s_mockScrollCrossDomainMessenger = new MockScrollCrossDomainMessenger(s_l1OwnerAddr);
    s_scrollCrossDomainForwarder = new ScrollCrossDomainForwarder(s_mockScrollCrossDomainMessenger, s_l1OwnerAddr);
    s_greeter = new Greeter(address(s_scrollCrossDomainForwarder));
    vm.stopPrank();
  }
}

contract ScrollCrossDomainForwarder_Constructor is ScrollCrossDomainForwarderTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialState() public {
    // it should set the owner correctly
    assertEq(s_scrollCrossDomainForwarder.owner(), s_l1OwnerAddr);

    // it should set the l1Owner correctly
    assertEq(s_scrollCrossDomainForwarder.l1Owner(), s_l1OwnerAddr);

    // it should set the crossdomain messenger correctly
    assertEq(s_scrollCrossDomainForwarder.crossDomainMessenger(), address(s_mockScrollCrossDomainMessenger));

    // it should set the typeAndVersion correctly
    assertEq(s_scrollCrossDomainForwarder.typeAndVersion(), "ScrollCrossDomainForwarder 1.0.0");
  }
}

contract ScrollCrossDomainForwarder_Forward is ScrollCrossDomainForwarderTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_scrollCrossDomainForwarder.forward(address(s_greeter), abi.encode(""));
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_Forward() public {
    vm.startPrank(s_strangerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      encodeCrossDomainSetGreetingMsg(s_scrollCrossDomainForwarder.forward.selector, address(s_greeter), greeting), // message
      0 // gas limit
    );

    // Checks that the greeter got the message
    assertEq(s_greeter.greeting(), greeting);
  }

  /// @notice it should revert when contract call reverts
  function test_ForwardRevert() public {
    vm.startPrank(s_strangerAddr);

    // Sends an invalid message
    vm.expectRevert("Invalid greeting length");
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      encodeCrossDomainSetGreetingMsg(s_scrollCrossDomainForwarder.forward.selector, address(s_greeter), ""), // message
      0 // gas limit
    );
  }
}

contract ScrollCrossDomainForwarder_TransferL1Ownership is ScrollCrossDomainForwarderTest {
  /// @notice it should not be callable by non-owners
  function test_NotCallableByNonOwners() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_scrollCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should not be callable by L2 owner
  function test_NotCallableByL2Owner() public {
    vm.startPrank(s_l1OwnerAddr);
    assertEq(s_scrollCrossDomainForwarder.owner(), s_l1OwnerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_scrollCrossDomainForwarder.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should be callable by current L1 owner
  function test_CallableByL1Owner() public {
    vm.startPrank(s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit();
    emit L1OwnershipTransferRequested(s_scrollCrossDomainForwarder.l1Owner(), s_strangerAddr);

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainForwarder.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );
  }

  /// @notice it should be callable by current L1 owner to zero address
  function test_CallableByL1OwnerOrZeroAddress() public {
    vm.startPrank(s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit();
    emit L1OwnershipTransferRequested(s_scrollCrossDomainForwarder.l1Owner(), address(0));

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainForwarder.transferL1Ownership.selector, address(0)), // message
      0 // gas limit
    );
  }
}

contract ScrollCrossDomainForwarder_AcceptL1Ownership is ScrollCrossDomainForwarderTest {
  /// @notice it should not be callable by non pending-owners
  function test_NotCallableByNonPendingOwners() public {
    vm.startPrank(s_strangerAddr);

    // Sends the message
    vm.expectRevert("Must be proposed L1 owner");
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainForwarder.acceptL1Ownership.selector), // message
      0 // gas limit
    );
  }

  /// @notice it should be callable by pending L1 owner
  function test_CallableByPendingL1Owner() public {
    vm.startPrank(s_strangerAddr);

    // Request ownership transfer
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainForwarder.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Sets a mock message sender
    s_mockScrollCrossDomainMessenger._setMockMessageSender(s_strangerAddr);

    // Prepares expected event payload
    vm.expectEmit();
    emit L1OwnershipTransferred(s_l1OwnerAddr, s_strangerAddr);

    // Accepts ownership transfer request
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainForwarder), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainForwarder.acceptL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Asserts that the ownership was actually transferred
    assertEq(s_scrollCrossDomainForwarder.l1Owner(), s_strangerAddr);
  }
}
