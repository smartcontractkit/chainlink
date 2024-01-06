// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ArbitrumCrossDomainGovernor} from "../../../dev/arbitrum/ArbitrumCrossDomainGovernor.sol";
import {MultiSend} from "../../../../../v0.8/vendor/MultiSend.sol";
import {Greeter} from "../../../../../v0.8/tests/Greeter.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

// Use this command from the /contracts directory to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --match-path ./src/v0.8/l2ep/test/v1_0_0/arbitrum/ArbitrumCrossDomainGovernor.t.sol
//
contract ArbitrumCrossDomainGovernorTest is L2EPTest {
  /// Helper variables
  address internal s_strangerAddr = vm.addr(0x1);
  address internal s_l1OwnerAddr = vm.addr(0x2);
  address internal s_crossDomainMessengerAddr = toArbitrumL2AliasAddress(s_l1OwnerAddr);
  address internal s_newOwnerCrossDomainMessengerAddr = toArbitrumL2AliasAddress(s_strangerAddr);

  /// Contracts
  ArbitrumCrossDomainGovernor internal s_arbitrumCrossDomainGovernor;
  MultiSend internal s_multiSend;
  Greeter internal s_greeter;

  /// Events
  event L1OwnershipTransferRequested(address indexed from, address indexed to);
  event L1OwnershipTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    s_arbitrumCrossDomainGovernor = new ArbitrumCrossDomainGovernor(s_l1OwnerAddr);
    s_greeter = new Greeter(address(s_arbitrumCrossDomainGovernor));
    s_multiSend = new MultiSend();
    vm.stopPrank();
  }
}

contract Constructor is ArbitrumCrossDomainGovernorTest {
  /// @notice it should set the owner correctly
  function test_Owner() public {
    assertEq(s_arbitrumCrossDomainGovernor.owner(), s_l1OwnerAddr);
  }

  /// @notice it should set the l1Owner correctly
  function test_L1Owner() public {
    assertEq(s_arbitrumCrossDomainGovernor.l1Owner(), s_l1OwnerAddr);
  }

  /// @notice it should set the crossdomain messenger correctly
  function test_CrossDomainMessenger() public {
    assertEq(s_arbitrumCrossDomainGovernor.crossDomainMessenger(), s_crossDomainMessengerAddr);
  }

  /// @notice it should set the typeAndVersion correctly
  function test_TypeAndVersion() public {
    assertEq(s_arbitrumCrossDomainGovernor.typeAndVersion(), "ArbitrumCrossDomainGovernor 1.0.0");
  }
}

contract Forward is ArbitrumCrossDomainGovernorTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger or owner");
    s_arbitrumCrossDomainGovernor.forward(address(s_greeter), abi.encode(""));
    vm.stopPrank();
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_Forward() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Defines the cross domain message to send
    string memory greeting = "hello";

    // Sends the message
    s_arbitrumCrossDomainGovernor.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, greeting)
    );

    // Checks that the greeter got the message
    assertEq(s_greeter.greeting(), greeting);

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
    s_arbitrumCrossDomainGovernor.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, greeting)
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), greeting);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert when contract call reverts
  function test_ForwardRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Sends an invalid message
    vm.expectRevert("Invalid greeting length");
    s_arbitrumCrossDomainGovernor.forward(
      address(s_greeter),
      abi.encodeWithSelector(s_greeter.setGreeting.selector, "")
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract ForwardDelegate is ArbitrumCrossDomainGovernorTest {
  /// @notice it should not be callable by unknown address
  function test_NotCallableByUnknownAddress() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger or owner");
    s_arbitrumCrossDomainGovernor.forwardDelegate(address(s_multiSend), abi.encode(""));
    vm.stopPrank();
  }

  /// @notice it should be callable by crossdomain messenger address / L1 owner
  function test_CallableByCrossDomainMessengerAddressOrL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Sends the message
    s_arbitrumCrossDomainGovernor.forwardDelegate(
      address(s_multiSend),
      abi.encodeWithSelector(
        MultiSend.multiSend.selector,
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), "bar"))
      )
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
    s_arbitrumCrossDomainGovernor.forwardDelegate(
      address(s_multiSend),
      abi.encodeWithSelector(
        MultiSend.multiSend.selector,
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), "bar"))
      )
    );

    // Checks that the greeter message was updated
    assertEq(s_greeter.greeting(), "bar");

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should revert batch when one call fails
  function test_RevertsBatchWhenOneCallFails() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Sends an invalid message (empty transaction data is not allowed)
    vm.expectRevert("Governor delegatecall reverted");
    s_arbitrumCrossDomainGovernor.forwardDelegate(
      address(s_multiSend),
      abi.encodeWithSelector(
        MultiSend.multiSend.selector,
        abi.encodePacked(encodeMultiSendTx(address(s_greeter), "foo"), encodeMultiSendTx(address(s_greeter), ""))
      )
    );

    // Checks that the greeter message is unchanged
    assertEq(s_greeter.greeting(), "");

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should bubble up revert when contract call reverts
  function test_BubbleUpRevert() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Sends an invalid message (empty transaction data is not allowed)
    vm.expectRevert("Greeter: revert triggered");
    s_arbitrumCrossDomainGovernor.forwardDelegate(
      address(s_greeter),
      abi.encodeWithSelector(Greeter.triggerRevert.selector)
    );

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract TransferL1Ownership is ArbitrumCrossDomainGovernorTest {
  /// @notice it should not be callable by non-owners
  function test_NotCallableByNonOwners() public {
    vm.startPrank(s_strangerAddr, s_strangerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_arbitrumCrossDomainGovernor.transferL1Ownership(s_strangerAddr);
    vm.stopPrank();
  }

  /// @notice it should not be callable by L2 owner
  function test_NotCallableByL2Owner() public {
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    assertEq(s_arbitrumCrossDomainGovernor.owner(), s_l1OwnerAddr);
    vm.expectRevert("Sender is not the L2 messenger");
    s_arbitrumCrossDomainGovernor.transferL1Ownership(s_strangerAddr);
    vm.stopPrank();
  }

  /// @notice it should be callable by current L1 owner
  function test_CallableByL1Owner() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Defines the cross domain message to send
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferRequested(s_arbitrumCrossDomainGovernor.l1Owner(), s_strangerAddr);

    // Sends the message
    s_arbitrumCrossDomainGovernor.transferL1Ownership(s_strangerAddr);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by current L1 owner to zero address
  function test_CallableByL1OwnerOrZeroAddress() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Defines the cross domain message to send
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferRequested(s_arbitrumCrossDomainGovernor.l1Owner(), address(0));

    // Sends the message
    s_arbitrumCrossDomainGovernor.transferL1Ownership(address(0));

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}

contract AcceptL1Ownership is ArbitrumCrossDomainGovernorTest {
  /// @notice it should not be callable by non pending-owners
  function test_NotCallableByNonPendingOwners() public {
    // Sets msg.sender and tx.origin
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);

    // Sends the message
    vm.expectRevert("Must be proposed L1 owner");
    s_arbitrumCrossDomainGovernor.acceptL1Ownership();

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }

  /// @notice it should be callable by pending L1 owner
  function test_CallableByPendingL1Owner() public {
    // Request ownership transfer
    vm.startPrank(s_crossDomainMessengerAddr, s_crossDomainMessengerAddr);
    s_arbitrumCrossDomainGovernor.transferL1Ownership(s_strangerAddr);

    // Prepares expected event payload
    vm.expectEmit(false, false, false, true);
    emit L1OwnershipTransferred(s_l1OwnerAddr, s_strangerAddr);

    // Accepts ownership transfer request
    vm.startPrank(s_newOwnerCrossDomainMessengerAddr, s_newOwnerCrossDomainMessengerAddr);
    s_arbitrumCrossDomainGovernor.acceptL1Ownership();

    // Asserts that the ownership was actually transferred
    assertEq(s_arbitrumCrossDomainGovernor.l1Owner(), s_strangerAddr);

    // Resets msg.sender and tx.origin
    vm.stopPrank();
  }
}
