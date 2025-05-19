// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {MockScrollL1CrossDomainMessenger} from "../../mocks/scroll/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/scroll/MockScrollL2CrossDomainMessenger.sol";
import {ScrollSequencerUptimeFeed} from "../../../scroll/ScrollSequencerUptimeFeed.sol";
import {BaseSequencerUptimeFeed} from "../../../base/BaseSequencerUptimeFeed.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ScrollSequencerUptimeFeedTestWrapper is ScrollSequencerUptimeFeed {
  constructor(
    address l1SenderAddress,
    address l2CrossDomainMessengerAddr,
    bool initialStatus
  ) ScrollSequencerUptimeFeed(l1SenderAddress, l2CrossDomainMessengerAddr, initialStatus) {}

  /// @notice It exposes the internal _validateSender function for testing
  function validateSenderTestWrapper(address l1Sender) external view {
    super._validateSender(l1Sender);
  }
}

contract ScrollSequencerUptimeFeed_Setup is L2EPTest {
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

contract ScrollSequencerUptimeFeed_Constructor is ScrollSequencerUptimeFeed_Setup {
  /// @notice Reverts when L2 Cross Domain Messenger address is invalid
  function test_Constructor_RevertWhen_InvalidL2XDomainMessenger() public {
    // L2 cross domain messenger address must not be the zero address
    vm.expectRevert(ScrollSequencerUptimeFeed.ZeroAddress.selector);
    new ScrollSequencerUptimeFeed(s_l1OwnerAddr, address(0), false);

    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Checks L1 sender
    address actualL1Addr = s_scrollSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_scrollSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }

  /// @notice Tests initial state with valid L2 Cross Domain Messenger
  function test_Constructor_InitialState_WhenValidL2XDomainMessenger() public {
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

contract ScrollSequencerUptimeFeed_ValidateSender is ScrollSequencerUptimeFeed_Setup {
  /// @notice Reverts when sender is not L2 Cross Domain Messenger address
  function test_ValidateSender_RevertWhen_SenderIsNotL2CrossDomainMessengerAddr() public {
    vm.startPrank(s_strangerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_scrollSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }

  /// @notice Reverts when L1 Cross Domain Messenger address is not L1 sender address
  function test_ValidateSender_RevertWhen_L1CrossDomainMessengerAddrIsNotL1SenderAddr() public {
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr);

    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_scrollSequencerUptimeFeed.validateSenderTestWrapper(s_strangerAddr);
  }

  /// @notice Updates status when status changes and incoming timestamp is the same as latest
  function test_ValidateSender_UpdateStatusWhen_StatusChangeAndNoTimeChange() public {
    address l2MessengerAddr = address(s_mockScrollL2CrossDomainMessenger);
    vm.startPrank(l2MessengerAddr);

    s_scrollSequencerUptimeFeed.validateSenderTestWrapper(s_l1OwnerAddr);
  }
}
