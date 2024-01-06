// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {MockScrollCrossDomainMessenger} from "../../mocks/scroll/MockScrollCrossDomainMessenger.sol";
import {ScrollCrossDomainGovernor} from "../../../dev/scroll/ScrollCrossDomainGovernor.sol";
import {MultiSend} from "../../../../../v0.8/vendor/MultiSend.sol";
import {Greeter} from "../../../../../v0.8/tests/Greeter.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/scroll/ScrollCrossDomainGovernor.t.sol
//
contract ScrollCrossDomainGovernorTest is L2EPTest {
  /// Helper variables
  address internal s_strangerAddr = vm.addr(0x1);
  address internal s_l1OwnerAddr = vm.addr(0x2);

  /// Contracts
  MockScrollCrossDomainMessenger internal s_mockScrollCrossDomainMessenger;
  ScrollCrossDomainGovernor internal s_scrollCrossDomainGovernor;
  MultiSend internal s_multiSend;
  Greeter internal s_greeter;

  /// Events
  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    s_mockScrollCrossDomainMessenger = new MockScrollCrossDomainMessenger(s_l1OwnerAddr);
    s_scrollCrossDomainGovernor = new ScrollCrossDomainGovernor(s_mockScrollCrossDomainMessenger, s_l1OwnerAddr);
    s_greeter = new Greeter(address(s_scrollCrossDomainGovernor));
    s_multiSend = new MultiSend();
    vm.stopPrank();
  }
}

contract Constructor is ScrollCrossDomainGovernorTest {
  /// @notice it should set the owner correctly
  function test_Owner() public {
    assertEq(s_scrollCrossDomainGovernor.owner(), s_l1OwnerAddr);
  }

  /// @notice it should set the l1Owner correctly
  function test_L1Owner() public {
    assertEq(s_scrollCrossDomainGovernor.l1Owner(), s_l1OwnerAddr);
  }

  /// @notice it should set the crossdomain messenger correctly
  function test_CrossDomainMessenger() public {
    assertEq(s_scrollCrossDomainGovernor.crossDomainMessenger(), address(s_mockScrollCrossDomainMessenger));
  }

  /// @notice it should set the typeAndVersion correctly
  function test_TypeAndVersion() public {
    assertEq(s_scrollCrossDomainGovernor.typeAndVersion(), "ScrollCrossDomainGovernor 1.0.0");
  }
}

contract Forward is ScrollCrossDomainGovernorTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger or owner");
    s_scrollCrossDomainGovernor.forward(address(s_greeter), abi.encode(""));
    vm.stopPrank();
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_Forward() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      encodeCrossDomainSetGreetingMsg(s_scrollCrossDomainGovernor.forward.selector, address(s_greeter), greeting), // message
      0 // gas limit
    );

    // Checks that the greeter got the message
    assertEq(s_greeter.greeting(), greeting);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert when contract call reverts
  function test_ForwardRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends an invalid message
    vm.expectRevert("Invalid greeting length");
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      encodeCrossDomainSetGreetingMsg(s_scrollCrossDomainGovernor.forward.selector, address(s_greeter), ""), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by L2 owner
  function test_CallableByL2Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_scrollCrossDomainGovernor.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, greeting)
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), greeting);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract ForwardDelegate is ScrollCrossDomainGovernorTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger or owner");
    s_scrollCrossDomainGovernor.forwardDelegate(address(s_greeter), abi.encode(""));
    vm.stopPrank();
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_CallableByCrossDomainMessengerAddressOrL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      encodeCrossDomainMultiSendMsg(
        s_scrollCrossDomainGovernor.forwardDelegate.selector,
        address(s_multiSend),
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), "bar"))
      ), // message
      0 // gas limit
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), "bar");

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by L2 owner
  function test_CallableByL2Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      encodeCrossDomainMultiSendMsg(
        s_scrollCrossDomainGovernor.forwardDelegate.selector,
        address(s_multiSend),
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), "bar"))
      ), // message
      0 // gas limit
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), "bar");

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert batch when one call fails
  function test_RevertsBatchWhenOneCallFails() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends an invalid message (empty transaction data is not allowed)
    vm.expectRevert("Governor delegatecall reverted");
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      encodeCrossDomainMultiSendMsg(
        s_scrollCrossDomainGovernor.forwardDelegate.selector,
        address(s_multiSend),
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), ""))
      ), // message
      0 // gas limit
    );

    // Checks that the greeter message is unchanged
    assertEq(s_greeter.greeting(), "");

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should bubble up revert when contract call reverts
  function test_BubbleUpRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends an invalid message (empty transaction data is not allowed)
    vm.expectRevert("Greeter: revert triggered");
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      abi.encodeWithSelector(
        ScrollCrossDomainGovernor.forwardDelegate.selector,
        address(s_greeter),
        abi.encodeWithSelector(Greeter.triggerRevert.selector)
      ), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract TransferL1Ownership is ScrollCrossDomainGovernorTest {
  /// @notice it should not be callable by non-owners
  function test_NotCallableByNonOwners() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_scrollCrossDomainGovernor.transferL1Ownership(s_strangerAddr);
    vm.stopPrank();
  }

  /// @notice it should not be callable by L2 owner
  function test_NotCallableByL2Owner() public {
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    assertEq(s_scrollCrossDomainGovernor.owner(), s_l1OwnerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_scrollCrossDomainGovernor.transferL1Ownership(s_strangerAddr);
    vm.stopPrank();
  }

  /// @notice it should be callable by current L1 owner
  function test_CallableByL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferRequested(s_scrollCrossDomainGovernor.l1Owner(), s_strangerAddr);

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainGovernor.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by current L1 owner to zero address
  function test_CallableByL1OwnerOrZeroAddress() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Defines the cross domain message to send
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferRequested(s_scrollCrossDomainGovernor.l1Owner(), address(0));

    // Sends the message
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainGovernor.transferL1Ownership.selector, address(0)), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AcceptL1Ownership is ScrollCrossDomainGovernorTest {
  /// @notice it should not be callable by non pending-owners
  function test_NotCallableByNonPendingOwners() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Sends the message
    vm.expectRevert("Must be proposed L1 owner");
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainGovernor.acceptL1Ownership.selector), // message
      0 // gas limit
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by pending L1 owner
  function test_CallableByPendingL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Request ownership transfer
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainGovernor.transferL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Sets a mock message sender
    s_mockScrollCrossDomainMessenger._setMockMessageSender(s_strangerAddr);

    // Prepares expected event payload
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferred(s_l1OwnerAddr, s_strangerAddr);

    // Accepts ownership transfer request
    s_mockScrollCrossDomainMessenger.sendMessage(
      address(s_scrollCrossDomainGovernor), // target
      0, // value
      abi.encodeWithSelector(s_scrollCrossDomainGovernor.acceptL1Ownership.selector, s_strangerAddr), // message
      0 // gas limit
    );

    // Asserts that the ownership was actually transferred
    assertEq(s_scrollCrossDomainGovernor.l1Owner(), s_strangerAddr);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}
