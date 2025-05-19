// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {OptimismCrossDomainGovernor} from "../../../optimism/OptimismCrossDomainGovernor.sol";
import {MockOVMCrossDomainMessenger} from "../../mocks/optimism/MockOVMCrossDomainMessenger.sol";
import {Greeter} from "../../Greeter.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

import {MultiSend} from "../../../../vendor/MultiSend.sol";

contract OptimismCrossDomainGovernorTest is L2EPTest {
  /// Contracts
  MockOVMCrossDomainMessenger internal s_mockOptimismCrossDomainMessenger;
  OptimismCrossDomainGovernor internal s_optimismCrossDomainGovernor;
  MultiSend internal s_multiSend;
  Greeter internal s_greeter;

  /// Events
  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    vm.startPrank(s_l1OwnerAddr);
    s_mockOptimismCrossDomainMessenger = new MockOVMCrossDomainMessenger(s_l1OwnerAddr);
    s_optimismCrossDomainGovernor = new OptimismCrossDomainGovernor(s_mockOptimismCrossDomainMessenger, s_l1OwnerAddr);
    s_greeter = new Greeter(address(s_optimismCrossDomainGovernor));
    s_multiSend = new MultiSend();
    vm.stopPrank();
  }
}

contract OptimismCrossDomainGovernor_Constructor is OptimismCrossDomainGovernorTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialState() public {
    // it should set the owner correctly
    assertEq(s_optimismCrossDomainGovernor.owner(), s_l1OwnerAddr);

    // it should set the l1Owner correctly
    assertEq(s_optimismCrossDomainGovernor.l1Owner(), s_l1OwnerAddr);

    // it should set the crossdomain messenger correctly
    assertEq(s_optimismCrossDomainGovernor.crossDomainMessenger(), address(s_mockOptimismCrossDomainMessenger));

    // it should set the typeAndVersion correctly
    assertEq(s_optimismCrossDomainGovernor.typeAndVersion(), "OptimismCrossDomainGovernor 1.0.0");
  }
}

contract OptimismCrossDomainGovernor_Forward is OptimismCrossDomainGovernorTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger or owner");
    s_optimismCrossDomainGovernor.forward(address(s_greeter), abi.encode(""));
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_Forward() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      encodeCrossDomainSetGreetingMsg(s_optimismCrossDomainGovernor.forward.selector, address(s_greeter), greeting), // message
      0 // gas limit
    );

    // Checks that the greeter got the message
    assertEq(s_greeter.greeting(), greeting);
  }

  /// @notice it should revert when contract call reverts
  function test_ForwardRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Sends an invalid message
    vm.expectRevert("Invalid greeting length");
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      encodeCrossDomainSetGreetingMsg(s_optimismCrossDomainGovernor.forward.selector, address(s_greeter), ""), // message
      0 // gas limit
    );
  }

  /// @notice it should be callable by L2 owner
  function test_CallableByL2Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_l1OwnerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_optimismCrossDomainGovernor.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, greeting)
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), greeting);
  }
}

contract OptimismCrossDomainGovernor_ForwardDelegate is OptimismCrossDomainGovernorTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger or owner");
    s_optimismCrossDomainGovernor.forwardDelegate(address(s_greeter), abi.encode(""));
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_CallableByCrossDomainMessengerAddressOrL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      encodeCrossDomainMultiSendMsg(
        s_optimismCrossDomainGovernor.forwardDelegate.selector,
        address(s_multiSend),
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), "bar"))
      ), // message
      0 // gas limit
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), "bar");
  }

  /// @notice it should be callable by L2 owner
  function test_CallableByL2Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_l1OwnerAddr);

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      encodeCrossDomainMultiSendMsg(
        s_optimismCrossDomainGovernor.forwardDelegate.selector,
        address(s_multiSend),
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), "bar"))
      ), // message
      0 // gas limit
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), "bar");
  }

  /// @notice it should revert batch when one call fails
  function test_RevertsBatchWhenOneCallFails() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Sends an invalid message (empty transaction data is not allowed)
    vm.expectRevert("Governor delegatecall reverted");
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      encodeCrossDomainMultiSendMsg(
        s_optimismCrossDomainGovernor.forwardDelegate.selector,
        address(s_multiSend),
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), ""))
      ), // message
      0 // gas limit
    );

    // Checks that the greeter message is unchanged
    assertEq(s_greeter.greeting(), "");
  }

  /// @notice it should bubble up revert when contract call reverts
  function test_BubbleUpRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Sends an invalid message (empty transaction data is not allowed)
    vm.expectRevert("Greeter: revert triggered");
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      abi.encodeWithSelector(
        OptimismCrossDomainGovernor.forwardDelegate.selector,
        address(s_greeter),
        abi.encodeWithSelector(Greeter.triggerRevert.selector)
      ), // message
      0 // gas limit
    );
  }
}

contract OptimismCrossDomainGovernor_TransferL1Ownership is OptimismCrossDomainGovernorTest {
  /// @notice it should not be callable by non-owners
  function test_NotCallableByNonOwners() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_optimismCrossDomainGovernor.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should not be callable by L2 owner
  function test_NotCallableByL2Owner() public {
    vm.startPrank(s_l1OwnerAddr);
    assertEq(s_optimismCrossDomainGovernor.owner(), s_l1OwnerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_optimismCrossDomainGovernor.transferL1Ownership(s_strangerAddr);
  }

  /// @notice it should be callable by current L1 owner
  function test_CallableByL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit();
    emit L1OwnershipTransferRequested(s_optimismCrossDomainGovernor.l1Owner(), s_strangerAddr);

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      abi.encodeWithSelector(s_optimismCrossDomainGovernor.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );
  }

  /// @notice it should be callable by current L1 owner to zero address
  function test_CallableByL1OwnerOrZeroAddress() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit();
    emit L1OwnershipTransferRequested(s_optimismCrossDomainGovernor.l1Owner(), address(0));

    // Sends the message
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      abi.encodeWithSelector(s_optimismCrossDomainGovernor.transferL1Ownership.selector, address(0)), // message
      0 // gas limit
    );
  }
}

contract OptimismCrossDomainGovernor_AcceptL1Ownership is OptimismCrossDomainGovernorTest {
  /// @notice it should not be callable by non pending-owners
  function test_NotCallableByNonPendingOwners() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Sends the message
    vm.expectRevert("Must be proposed L1 owner");
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      abi.encodeWithSelector(s_optimismCrossDomainGovernor.acceptL1Ownership.selector), // message
      0 // gas limit
    );
  }

  /// @notice it should be callable by pending L1 owner
  function test_CallableByPendingL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr);

    // Request ownership transfer
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      abi.encodeWithSelector(s_optimismCrossDomainGovernor.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Sets a mock message sender
    s_mockOptimismCrossDomainMessenger._setMockMessageSender(s_strangerAddr);

    // Prepares expected event payload
    vm.expectEmit();
    emit L1OwnershipTransferred(s_l1OwnerAddr, s_strangerAddr);

    // Accepts ownership transfer request
    s_mockOptimismCrossDomainMessenger.sendMessage(
      address(s_optimismCrossDomainGovernor), // target
      abi.encodeWithSelector(s_optimismCrossDomainGovernor.acceptL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Asserts that the ownership was actually transferred
    assertEq(s_optimismCrossDomainGovernor.l1Owner(), s_strangerAddr);
  }
}
