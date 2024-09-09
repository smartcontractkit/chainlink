// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockScrollL1CrossDomainMessenger} from "../../mocks/scroll/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/scroll/MockScrollL2CrossDomainMessenger.sol";
import {ScrollSequencerUptimeFeed} from "../../../dev/scroll/ScrollSequencerUptimeFeed.sol";
import {BaseSequencerUptimeFeed} from "../../../dev/shared/BaseSequencerUptimeFeed.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ScrollSequencerUptimeFeedTestWrapper is ScrollSequencerUptimeFeed {
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) ScrollSequencerUptimeFeed(l1SenderAddress, l2CrossDomainMessengerAddr, initialStatus) {}

  /// @notice it exposes the internal _validateSender function for testing
  function validateSenderTestWrapper(address l1Sender) external view {
    super._validateSender(l1Sender);
  }
}

contract ScrollSequencerUptimeFeedTest is L2EPTest {
  /// Constants
  uint256 internal constant GAS_USED_DEVIATION = 100;

  /// L2EP contracts
  MockScrollL1CrossDomainMessenger internal s_mockScrollL1CrossDomainMessenger;
  MockScrollL2CrossDomainMessenger internal s_mockScrollL2CrossDomainMessenger;
  ScrollSequencerUptimeFeedTestWrapper internal s_scrollSequencerUptimeFeed;

  /// Events
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);
  event RoundUpdated(int256 status, uint64 updatedAt);

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_mockScrollL1CrossDomainMessenger = new MockScrollL1CrossDomainMessenger();
    s_mockScrollL2CrossDomainMessenger = new MockScrollL2CrossDomainMessenger();
    s_scrollSequencerUptimeFeed = new ScrollSequencerUptimeFeedTestWrapper(
      s_l1OwnerAddr,
      address(s_mockScrollL2CrossDomainMessenger),
      false
    );

    // Sets mock sender in mock L2 messenger contract
    s_mockScrollL2CrossDomainMessenger.setSender(s_l1OwnerAddr);
  }
}

contract ScrollSequencerUptimeFeed_Constructor is ScrollSequencerUptimeFeedTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialStateWithInvalidL2XDomainManager() public {
    // L2 cross domain messenger address must not be the zero address
    vm.expectRevert(ScrollSequencerUptimeFeed.ZeroAddress.selector);
    new ScrollSequencerUptimeFeed(s_l1OwnerAddr, address(0), false);

    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Checks L1 sender
    address actualL1Addr = s_scrollSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_scrollSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }

  function test_InitialStateWithValidL2XDomainManager() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);
    ScrollSequencerUptimeFeed scrollSequencerUptimeFeed = new ScrollSequencerUptimeFeed(
      s_l1OwnerAddr,
      address(s_mockScrollL2CrossDomainMessenger),
      false
    );

    // Checks L1 sender
    address actualL1Addr = scrollSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = scrollSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract ScrollSequencerUptimeFeed_ValidateSender is ScrollSequencerUptimeFeedTest {
  /// @notice it should revert if called by an address that is not the L2 Cross Domain Messenger
  function test_RevertIfSenderIsNotL2CrossDomainMessengerAddr() public {
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    // Sets msg.sender to a different address
    vm.startPrank(s_strangerAddr, l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_scrollSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }

  /// @notice it should revert if the L1 sender address is not the L1 Cross Domain Messenger Sender
  function test_RevertIfL1CrossDomainMessengerAddrIsNotL1SenderAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_scrollSequencerUptimeFeed.validateSenderTestWrapper(s_strangerAddr);
  }

  /// @notice it should update status when status has changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenStatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr, l2MessengerAddr);

    s_scrollSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }
}
